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

// Package controllers provides the controllers for CAPHV.
package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	secretutil "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/secrets"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	secretErrorRetryDelay = time.Second * 10
	rateLimitWaitTime     = 5 * time.Minute
)

// HivelocityClusterReconciler reconciles a HivelocityCluster object.
type HivelocityClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	WatchFilterValue string

	APIReader       client.Reader
	HVClientFactory hvclient.Factory
}

//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters/finalizers,verbs=update

// Reconcile aims to move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *HivelocityClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the HivelocityCluster
	hvCluster := &infrav1.HivelocityCluster{}
	err := r.Get(ctx, req.NamespacedName, hvCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	log = log.WithValues("HivelocityCluster", klog.KObj(hvCluster))

	log.Info("Starting reconciling cluster")

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, hvCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get owner cluster: %w", err)
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{
			RequeueAfter: 2 * time.Second,
		}, nil
	}

	log = log.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, log)

	if annotations.IsPaused(cluster, hvCluster) {
		log.Info("HivelocityCluster or linked Cluster is marked as paused. Won't reconcile")
		return reconcile.Result{}, nil
	}

	log.V(1).Info("Creating cluster scope")

	// Create the scope.
	secretManager := secretutil.NewSecretManager(log, r.Client, r.APIReader)
	apiKey, hvSecret, err := getAndValidateHivelocityAPIKey(ctx, req.Namespace, hvCluster, secretManager)
	if err != nil {
		return hvAPIKeyErrorResult(ctx, err, hvCluster, infrav1.HivelocityClusterReady, r.Client)
	}

	hvClient := r.HVClientFactory.NewClient(apiKey)

	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:            r.Client,
		Logger:            log,
		Cluster:           cluster,
		HivelocityCluster: hvCluster,
		HVClient:          hvClient,
	})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}
	// Always close the scope when exiting this function so we can persist any HivelocityCluster changes.
	defer func() {
		if err := clusterScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// check whether rate limit has been reached and if so, then wait.
	if wait := reconcileRateLimit(hvCluster); wait {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// Handle deleted clusters
	if !hvCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope, hvSecret)
	}

	return r.reconcileNormal(ctx, clusterScope)
}

func (r *HivelocityClusterReconciler) reconcileNormal(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	log := clusterScope.Logger
	log.V(1).Info("Reconciling HivelocityCluster oooooooooooooooooooooooooooooooooooooooooo")

	hvCluster := clusterScope.HivelocityCluster

	// If the HivelocityCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(hvCluster, infrav1.ClusterFinalizer)
	if err := clusterScope.PatchObject(ctx); err != nil {
		return ctrl.Result{}, err
	}

	/* question clusterScope.SetStatusFailureDomain
	// set failure domains in status using information in spec
	clusterScope.SetStatusFailureDomain(clusterScope.GetSpecRegion())
	*/

	hvCluster.Spec.ControlPlaneEndpoint = &clusterv1.APIEndpoint{
		Host: "66.165.243.74",
		Port: 443,
	}
	hvCluster.Status.Ready = true

	if err := r.reconcileTargetClusterManager(ctx, clusterScope); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile target cluster manager: %w", err)
	}

	log.V(1).Info("Reconciling finished")
	return reconcile.Result{}, nil
}

func (r *HivelocityClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope, hvSecret *corev1.Secret) (reconcile.Result, error) {
	log := clusterScope.Logger

	log.Info("Reconciling HivelocityCluster delete")

	hvCluster := clusterScope.HivelocityCluster

	// wait for all HivelocityMachines to be deleted
	machines, _, err := clusterScope.ListMachines(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list machines for HivelocityCluster %s/%s: %w",
			hvCluster.Namespace, hvCluster.Name, err)
	}
	if len(machines) > 0 {
		names := make([]string, len(machines))
		for i, m := range machines {
			names[i] = fmt.Sprintf("machine/%s", m.Name)
		}
		record.Eventf(
			hvCluster,
			"WaitingForMachineDeletion",
			"Machines %s still running, waiting with deletion of HivelocityCluster",
			strings.Join(names, ", "),
		)
		return reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	}

	secretManager := secretutil.NewSecretManager(log, r.Client, r.APIReader)
	// Remove finalizer of secret
	if err := secretManager.ReleaseSecret(ctx, hvSecret); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to release HivelocitySecret: %w", err)
	}

	/* todo, later.
	// Stop CSR manager
	r.targetClusterManagersLock.Lock()
	defer r.targetClusterManagersLock.Unlock()

	key := types.NamespacedName{
		Namespace: clusterScope.HivelocityCluster.Namespace,
		Name:      clusterScope.HivelocityCluster.Name,
	}
	if stopCh, ok := r.targetClusterManagersStopCh[key]; ok {
		close(stopCh)
		delete(r.targetClusterManagersStopCh, key)
	}
	*/

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(clusterScope.HivelocityCluster, infrav1.ClusterFinalizer)

	return reconcile.Result{}, nil
}

// reconcileRateLimit checks whether a rate limit has been reached and returns whether
// the controller should wait a bit more.
func reconcileRateLimit(setter conditions.Setter) bool {
	condition := conditions.Get(setter, infrav1.RateLimitExceeded)
	if condition != nil && condition.Status == corev1.ConditionTrue {
		if time.Now().Before(condition.LastTransitionTime.Time.Add(rateLimitWaitTime)) {
			// Not yet timed out, reconcile again after timeout
			// Don't give a more precise requeueAfter value to not reconcile too many
			// objects at the same time
			return true
		}
		// Wait time is over, we continue
		conditions.MarkFalse(
			setter,
			infrav1.RateLimitExceeded,
			infrav1.RateLimitNotReachedReason,
			clusterv1.ConditionSeverityInfo,
			"wait time is over. Try reconciling again",
		)
	}
	return false
}

func getAndValidateHivelocityAPIKey(ctx context.Context, namespace string, hvCluster *infrav1.HivelocityCluster, secretManager *secretutil.SecretManager) (string, *corev1.Secret, error) {
	// retrieve Hivelocity secret
	secretNamspacedName := types.NamespacedName{Namespace: namespace, Name: hvCluster.Spec.HivelocitySecret.Name}

	hvSecret, err := secretManager.AcquireSecret(
		ctx,
		secretNamspacedName,
		hvCluster,
		false,
		hvCluster.DeletionTimestamp.IsZero(),
	)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return "", nil, &secretutil.ResolveSecretRefError{
				Message: fmt.Sprintf("The Hivelocity secret %s does not exist", secretNamspacedName),
			}
		}
		return "", nil, err
	}

	apiKey := string(hvSecret.Data[hvCluster.Spec.HivelocitySecret.Key])

	// Validate apiKey
	if apiKey == "" {
		return "", nil, &secretutil.HivelocityAPIKeyValidationError{}
	}

	return apiKey, hvSecret, nil
}

func hvAPIKeyErrorResult(
	ctx context.Context,
	err error,
	setter conditions.Setter,
	conditionType clusterv1.ConditionType,
	client client.Client,
) (res ctrl.Result, reterr error) {
	switch err.(type) {
	// In the event that the reference to the secret is defined, but we cannot find it
	// we requeue the host as we will not know if they create the secret
	// at some point in the future.
	case *secretutil.ResolveSecretRefError:
		conditions.MarkFalse(setter,
			conditionType,
			infrav1.HivelocitySecretUnreachableReason,
			clusterv1.ConditionSeverityError,
			"could not find HivelocitySecret",
		)
		res = ctrl.Result{Requeue: true, RequeueAfter: secretErrorRetryDelay}

	// No need to reconcile again, as it will be triggered as soon as the secret is updated.
	case *secretutil.HivelocityAPIKeyValidationError:
		conditions.MarkFalse(setter,
			conditionType,
			infrav1.HivelocityCredentialsInvalidReason,
			clusterv1.ConditionSeverityError,
			"invalid or not specified credentials for Hivelocity in secret",
		)

	default:
		return ctrl.Result{}, fmt.Errorf("an unhandled failure occurred with the Hivelocity secret: %w", err)
	}

	if err := client.Status().Update(ctx, setter); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to update: %w", err)
	}

	return res, err
}

func (r *HivelocityClusterReconciler) reconcileTargetClusterManager(ctx context.Context, clusterScope *scope.ClusterScope) error {
	/* question targetClusterManagersLock.Lock()
	r.targetClusterManagersLock.Lock()
	defer r.targetClusterManagersLock.Unlock()

	key := types.NamespacedName{
		Namespace: clusterScope.HivelocityCluster.Namespace,
		Name:      clusterScope.HivelocityCluster.Name,
	}

	if _, ok := r.targetClusterManagersStopCh[key]; !ok {
		// create a new cluster manager
		m, err := r.newTargetClusterManager(ctx, clusterScope)
		if err != nil {
			return fmt.Errorf("failed to create a clusterManager for HivelocityCluster %s/%s: %w",
				clusterScope.HivelocityCluster.Namespace,
				clusterScope.HivelocityCluster.Name,
				err,
		    )
		}
		r.targetClusterManagersStopCh[key] = make(chan struct{})

		ctx, cancel := context.WithCancel(ctx)

		r.TargetClusterManagersWaitGroup.Add(1)

		// Start manager
		go func() {
			defer r.TargetClusterManagersWaitGroup.Done()

			if err := m.Start(ctx); err != nil {
				clusterScope.Error(err, "failed to start a targetClusterManager")
			} else {
				clusterScope.Info("stop targetClusterManager")
			}
			r.targetClusterManagersLock.Lock()
			defer r.targetClusterManagersLock.Unlock()
			delete(r.targetClusterManagersStopCh, key)
		}()

		// Cancel when stop channel received input
		go func() {
			<-r.targetClusterManagersStopCh[key]
			cancel()
		}()
	}
	*/
	return nil
}

/* question ManagementCluster, is this the target-cluster? Then the name is misleading.

var _ ManagementCluster = &managementCluster{}

type managementCluster struct {
	client.Client
	hvCluster *infrav1.HivelocityCluster
}

func (c *managementCluster) Namespace() string {
	return c.hvCluster.Namespace
}

func (r *HivelocityClusterReconciler) newTargetClusterManager(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Manager, error) {
	hvCluster := clusterScope.HivelocityCluster

	clientConfig, err := clusterScope.ClientConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get a clientConfig for the API of HivelocityCluster %s/%s: %w", hvCluster.Namespace, hvCluster.Name, err)
	}
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get a restConfig for the API of HivelocityCluster %s/%s: %w", hvCluster.Namespace, hvCluster.Name, err)
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get a clientSet for the API of HivelocityCluster %s/%s: %w", hvCluster.Namespace, hvCluster.Name, err)
	}

	scheme := runtime.NewScheme()
	_ = certificatesv1.AddToScheme(scheme)
	_ = infrav1.AddToScheme(scheme)

	clusterMgr, err := ctrl.NewManager(
		restConfig,
		ctrl.Options{
			Scheme:             scheme,
			MetricsBindAddress: "0",
			LeaderElection:     false,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup guest cluster manager: %w", err)
	}

	gr := &GuestCSRReconciler{
		Client: clusterMgr.GetClient(),
		mCluster: &managementCluster{
			Client:    r.Client,
			hvCluster: hvCluster,
		},
		WatchFilterValue: r.WatchFilterValue,
		clientSet:        clientSet,
	}

	if err := gr.SetupWithManager(ctx, clusterMgr, controller.Options{}); err != nil {
		return nil, fmt.Errorf("failed to setup CSR controller: %w", err)
	}

	return clusterMgr, nil
}.
*/

// SetupWithManager sets up the controller with the Manager.
func (r *HivelocityClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := log.FromContext(ctx)

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.HivelocityCluster{}).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		WithEventFilter(predicates.ResourceIsNotExternallyManaged(log)).
		Owns(&corev1.Secret{}).
		Build(r)
	if err != nil {
		return fmt.Errorf("error creating controller: %w", err)
	}

	return controller.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		handler.EnqueueRequestsFromMapFunc(func(o client.Object) []reconcile.Request {
			c, ok := o.(*clusterv1.Cluster)
			if !ok {
				panic(fmt.Sprintf("Expected a Cluster but got a %T", o))
			}

			log = log.WithValues("objectMapper", "clusterToHivelocityCluster", "namespace", c.Namespace, "cluster", c.Name)

			// Don't handle deleted clusters
			if !c.ObjectMeta.DeletionTimestamp.IsZero() {
				log.V(1).Info("Cluster has a deletion timestamp, skipping mapping.")
				return nil
			}

			// Make sure the ref is set
			if c.Spec.InfrastructureRef == nil {
				log.V(1).Info("Cluster does not have an InfrastructureRef, skipping mapping.")
				return nil
			}

			if c.Spec.InfrastructureRef.GroupVersionKind().Kind != "HivelocityCluster" {
				log.V(1).Info("Cluster has an InfrastructureRef for a different type, skipping mapping.")
				return nil
			}

			hvCluster := &infrav1.HivelocityCluster{}
			key := types.NamespacedName{Namespace: c.Spec.InfrastructureRef.Namespace, Name: c.Spec.InfrastructureRef.Name}

			if err := r.Get(ctx, key, hvCluster); err != nil {
				log.V(1).Error(err, "Failed to get HivelocityCluster")
				return nil
			}

			if annotations.IsExternallyManaged(hvCluster) {
				log.V(1).Info("HivelocityCluster is externally managed, skipping mapping.")
				return nil
			}

			log.V(1).Info("Adding request.", "hivelocityCluster", c.Spec.InfrastructureRef.Name)
			return []ctrl.Request{
				{
					NamespacedName: client.ObjectKey{Namespace: c.Namespace, Name: c.Spec.InfrastructureRef.Name},
				},
			}
		}),
	)
}
