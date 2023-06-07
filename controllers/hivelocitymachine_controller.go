/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	secretutil "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/secrets"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/device"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/predicates"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// HivelocityMachineReconciler reconciles a HivelocityMachine object.
type HivelocityMachineReconciler struct {
	client.Client
	APIReader        client.Reader
	HVClientFactory  hvclient.Factory
	WatchFilterValue string
}

//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachines/finalizers,verbs=update

// Reconcile manages the lifecycle of an Hivelocity machine object.
func (r *HivelocityMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	log.Info("Reconcile HivelocityMachine")

	// Fetch the HivelocityMachine.
	hivelocityMachine := &infrav1.HivelocityMachine{}
	err := r.Get(ctx, req.NamespacedName, hivelocityMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	log = log.WithValues("HivelocityMachine", klog.KObj(hivelocityMachine))

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, hivelocityMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("Machine", klog.KObj(machine))

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist", "error", err)
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, hivelocityMachine) {
		log.Info("HivelocityMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))

	hvCluster := &infrav1.HivelocityCluster{}

	hvClusterName := client.ObjectKey{
		Namespace: hivelocityMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, hvClusterName, hvCluster); err != nil {
		log.Info("HivelocityCluster is not available yet", "error", err)
		return reconcile.Result{}, nil
	}

	log = log.WithValues("HivelocityCluster", klog.KObj(hvCluster))
	ctx = ctrl.LoggerInto(ctx, log)

	// Create the scope.
	secretManager := secretutil.NewSecretManager(log, r.Client, r.APIReader)
	hvAPIKey, _, err := getAndValidateHivelocityAPIKey(ctx, req.Namespace, hvCluster, secretManager)
	if err != nil {
		return hvAPIKeyErrorResult(ctx, err, hivelocityMachine, infrav1.DeviceReadyCondition, r.Client)
	}

	hvClient := r.HVClientFactory.NewClient(hvAPIKey)

	machineScope, err := scope.NewMachineScope(ctx, scope.MachineScopeParams{
		ClusterScopeParams: scope.ClusterScopeParams{
			Client:            r.Client,
			Logger:            log,
			Cluster:           cluster,
			HivelocityCluster: hvCluster,
			HVClient:          hvClient,
			APIReader:         r.APIReader,
		},
		Machine:           machine,
		HivelocityMachine: hivelocityMachine,
	})
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Always close the scope when exiting this function so we can persist any HivelocityMachine changes.
	defer func() {
		if err := machineScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// check whether rate limit has been reached and if so, then wait.
	if wait := reconcileRateLimit(hivelocityMachine); wait {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	if !hivelocityMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		baseMsg := fmt.Sprintf("hivelocityMachine has DeletionTimestamp %s ProvisioningState=%s",
			hivelocityMachine.ObjectMeta.DeletionTimestamp.Time.Format("2006-01-02 15:04:05"),
			hivelocityMachine.Spec.Status.ProvisioningState)
		switch hivelocityMachine.Spec.Status.ProvisioningState {
		case infrav1.StateDeleteDevice:
			// Device has been removed from cluster - machine can be deleted.
			log.Info(baseMsg + " do: RemoveFinalizer")
			controllerutil.RemoveFinalizer(machineScope.HivelocityMachine, infrav1.MachineFinalizer)
			return reconcile.Result{}, nil
		case infrav1.StateNone, infrav1.StateAssociateDevice, infrav1.StateVerifyAssociate:
			// if device is not yet provisioned, we can just dissociate the device from the machine by deleting the tags.
			log.Info(baseMsg + " do: ProvisioningState = StateDeleteDeviceDissociate")
			hivelocityMachine.Spec.Status.ProvisioningState = infrav1.StateDeleteDeviceDissociate
			return ctrl.Result{}, nil
		case infrav1.StateDeviceProvisioned, infrav1.StateProvisionDevice:
			log.Info(baseMsg + " do: ProvisioningState = StateDeleteDeviceDeProvision")
			hivelocityMachine.Spec.Status.ProvisioningState = infrav1.StateDeleteDeviceDeProvision
			return ctrl.Result{}, nil
		case infrav1.StateDeleteDeviceDeProvision:
			log.Info(baseMsg + " waiting for deprovisioning to finish.")
		default:
			log.Info(baseMsg + " do: nothing?")
		}
	}

	return r.reconcile(ctx, machineScope)
}

func (r *HivelocityMachineReconciler) reconcile(ctx context.Context, machineScope *scope.MachineScope) (reconcile.Result, error) {
	machineScope.Info("Reconciling HivelocityMachine")
	hivelocityMachine := machineScope.HivelocityMachine

	// If the HivelocityMachine doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(machineScope.HivelocityMachine, infrav1.MachineFinalizer)

	// Register the finalizer immediately to avoid orphaning Hivelocity resources on delete
	if err := machineScope.PatchObject(ctx); err != nil {
		return ctrl.Result{}, err
	}

	// reconcile device
	result, err := device.NewService(machineScope).Reconcile(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to reconcile device for HivelocityMachine %s/%s: %w", hivelocityMachine.Namespace, hivelocityMachine.Name, err)
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HivelocityMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	c, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.HivelocityMachine{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		Watches(
			&source.Kind{Type: &clusterv1.Machine{}},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("HivelocityMachine"))),
		).
		Watches(
			&source.Kind{Type: &infrav1.HivelocityCluster{}},
			handler.EnqueueRequestsFromMapFunc(r.HivelocityClusterToHivelocityMachines(ctx, log)),
		).
		Build(r)
	if err != nil {
		return fmt.Errorf("error creating controller: %w", err)
	}

	clusterToObjectFunc, err := util.ClusterToObjectsMapper(r.Client, &infrav1.HivelocityMachineList{}, mgr.GetScheme())
	if err != nil {
		return fmt.Errorf("failed to create mapper for Cluster to HivelocityMachines: %w", err)
	}

	// Add a watch on clusterv1.Cluster object for unpause & ready notifications.
	if err := c.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		handler.EnqueueRequestsFromMapFunc(clusterToObjectFunc),
		predicates.ClusterUnpausedAndInfrastructureReady(log),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	return nil
}

// HivelocityClusterToHivelocityMachines is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of HivelocityMachines.
func (r *HivelocityMachineReconciler) HivelocityClusterToHivelocityMachines(ctx context.Context, log logr.Logger) handler.MapFunc {
	return func(o client.Object) []ctrl.Request {
		result := []ctrl.Request{}

		c, ok := o.(*infrav1.HivelocityCluster)
		if !ok {
			log.Error(errors.Errorf("expected a HivelocityCluster but got a %T", o), "failed to get HivelocityMachine for HivelocityCluster")
			return nil
		}

		log = log.WithValues("objectMapper", "hvClusterToHivelocityMachine", "namespace", c.Namespace, "hvCluster", c.Name)

		// Don't handle deleted HivelocityCluster
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(1).Info("HivelocityCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, c.ObjectMeta)
		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			log.V(1).Info("Cluster for HivelocityCluster not found, skipping mapping.")
			return result
		case err != nil:
			log.Error(err, "failed to get owning cluster, skipping mapping.")
			return result
		}

		labels := map[string]string{clusterv1.ClusterNameLabel: cluster.Name}
		machineList := &clusterv1.MachineList{}
		if err := r.List(ctx, machineList, client.InNamespace(c.Namespace), client.MatchingLabels(labels)); err != nil {
			log.Error(err, "failed to list Machines, skipping mapping.")
			return nil
		}
		for _, m := range machineList.Items {
			log = log.WithValues("machine", m.Name)
			if m.Spec.InfrastructureRef.GroupVersionKind().Kind != "HivelocityMachine" {
				log.V(1).Info("Machine has an InfrastructureRef for a different type, will not add to reconciliation request.")
				continue
			}
			if m.Spec.InfrastructureRef.Name == "" {
				continue
			}
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}
			log = log.WithValues("hivelocityMachine", name.Name)
			log.V(1).Info("Adding HivelocityMachine to reconciliation request.")
			result = append(result, ctrl.Request{NamespacedName: name})
		}

		return result
	}
}
