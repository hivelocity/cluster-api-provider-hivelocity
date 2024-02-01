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
	"strconv"
	"strings"
	"sync"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	secretutil "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/secrets"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/device"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
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

	targetClusterManagersStopCh    map[types.NamespacedName]chan struct{}
	targetClusterManagersLock      sync.Mutex
	TargetClusterManagersWaitGroup *sync.WaitGroup
}

//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters/finalizers,verbs=update

// Reconcile aims to move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *HivelocityClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	logger := ctrl.LoggerFrom(ctx)

	// Fetch the HivelocityCluster
	hvCluster := &infrav1.HivelocityCluster{}
	err := r.Get(ctx, req.NamespacedName, hvCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{RequeueAfter: 1 * time.Second}, err
	}

	logger = logger.WithValues("HivelocityCluster", klog.KObj(hvCluster))

	logger.Info("Starting reconciling cluster")

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, hvCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get owner cluster: %w", err)
	}

	if cluster == nil {
		logger.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{
			RequeueAfter: 2 * time.Second,
		}, nil
	}

	logger = logger.WithValues("Cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, logger)

	if annotations.IsPaused(cluster, hvCluster) {
		logger.Info("HivelocityCluster or linked Cluster is marked as paused. Won't reconcile")
		return reconcile.Result{}, nil
	}

	logger.V(1).Info("Creating cluster scope")

	// Create the scope.
	secretManager := secretutil.NewSecretManager(logger, r.Client, r.APIReader)
	apiKey, hvSecret, err := getAndValidateHivelocityAPIKey(ctx, req.Namespace, hvCluster, secretManager)
	if err != nil {
		return hvAPIKeyErrorResult(ctx, err, hvCluster, infrav1.CredentialsAvailableCondition, r.Client)
	}

	hvClient := r.HVClientFactory.NewClient(apiKey)

	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		Client:            r.Client,
		APIReader:         r.APIReader,
		Logger:            logger,
		Cluster:           cluster,
		HivelocityCluster: hvCluster,
		HVClient:          hvClient,
	})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}
	// Always close the scope when exiting this function so we can persist any HivelocityCluster changes.
	defer func() {
		conditions.SetSummary(hvCluster)

		if err := clusterScope.Close(ctx); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// check whether rate limit has been reached and if so, then wait.
	if wait := reconcileRateLimit(hvCluster); wait {
		// don't wait too long. Otherwise: context canceled
		return ctrl.Result{RequeueAfter: 20 * time.Second}, nil
	}

	// Handle deleted clusters
	if !hvCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope, hvSecret)
	}

	return r.reconcileNormal(ctx, clusterScope)
}

func (r *HivelocityClusterReconciler) reconcileNormal(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	logger := clusterScope.Logger
	logger.V(1).Info("Reconciling HivelocityCluster")

	hvCluster := clusterScope.HivelocityCluster

	// If the HivelocityCluster doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(hvCluster, infrav1.ClusterFinalizer)
	if err := clusterScope.PatchObject(ctx); err != nil {
		return ctrl.Result{}, err
	}

	// set failure domains in status using information in spec
	hvCluster.SetStatusFailureDomain(hvCluster.Spec.ControlPlaneRegion)

	// dirty hack. Loadbalancer are not supported yet.
	if hvCluster.Spec.ControlPlaneEndpoint.Host == "" {
		hmt := infrav1.HivelocityMachineTemplate{}
		name := hvCluster.Name + "-control-plane"

		err := r.Client.Get(ctx, client.ObjectKey{
			Namespace: hvCluster.ObjectMeta.Namespace,
			Name:      name,
		}, &hmt)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to get HivelocityMachineTemplate %q: %w", name, err)
		}

		hvDevice, reason, err := device.GetFirstFreeDevice(ctx, clusterScope.HVClient, hmt.Spec.Template.Spec, hvCluster)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("device.GetFirstFreeDevice() failed: %w (%s)", err, reason)
		}
		if hvDevice == nil {
			return ctrl.Result{}, fmt.Errorf("device.GetFirstFreeDevice() found no device: %+v", hmt.Spec.Template.Spec.DeviceSelector)
		}
		logger.Info(fmt.Sprintf("Setting hvCluster.Spec.ControlPlaneEndpoint.Host to %q", hvDevice.PrimaryIp))

		hvCluster.Spec.ControlPlaneEndpoint.Host = hvDevice.PrimaryIp
		hvCluster.Spec.ControlPlaneEndpoint.Port = 6443
	}

	emptyResult := reconcile.Result{}

	hvCluster.Status.Ready = true

	result, err := r.reconcileTargetClusterManager(ctx, clusterScope)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to reconcile target cluster manager: %w", err)
	}
	if result != emptyResult {
		return result, nil
	}

	conditions.MarkTrue(hvCluster, infrav1.TargetClusterReadyCondition)

	if err = reconcileTargetSecret(ctx, clusterScope); err != nil {
		reterr := fmt.Errorf("failed to reconcile target secret: %w", err)
		conditions.MarkFalse(
			clusterScope.HivelocityCluster,
			infrav1.TargetClusterSecretReadyCondition,
			infrav1.TargetSecretSyncFailedReason,
			clusterv1.ConditionSeverityError,
			reterr.Error(),
		)
		return reconcile.Result{}, reterr
	}

	// target cluster secret is ready
	conditions.MarkTrue(hvCluster, infrav1.TargetClusterSecretReadyCondition)

	logger.V(1).Info("Reconciling finished")
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

	// Cluster is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(clusterScope.HivelocityCluster, infrav1.ClusterFinalizer)

	return reconcile.Result{}, nil
}

// reconcileRateLimit checks whether a rate limit has been reached and returns whether
// the controller should wait a bit more.
func reconcileRateLimit(setter conditions.Setter) bool {
	condition := conditions.Get(setter, infrav1.HivelocityAPIReachableCondition)
	if condition != nil && condition.Status == corev1.ConditionFalse {
		if time.Now().Before(condition.LastTransitionTime.Time.Add(rateLimitWaitTime)) {
			// Not yet timed out, reconcile again after timeout
			// Don't give a more precise requeueAfter value to not reconcile too many
			// objects at the same time
			return true
		}
		// Wait time is over, we continue
		conditions.MarkTrue(setter, infrav1.HivelocityAPIReachableCondition)
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

func reconcileTargetSecret(ctx context.Context, clusterScope *scope.ClusterScope) error {
	clientConfig, err := clusterScope.ClientConfig(ctx)
	if err != nil {
		clusterScope.V(1).Info("failed to get clientconfig with api endpoint")
		return err
	}

	if err := scope.IsControlPlaneReady(ctx, clientConfig); err != nil {
		conditions.MarkFalse(
			clusterScope.HivelocityCluster,
			infrav1.TargetClusterSecretReadyCondition,
			infrav1.TargetClusterControlPlaneNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"target cluster not ready",
		)
		return nil //nolint:nilerr
	}

	// Workload Control plane ready, so we can check if the secret exists already

	// getting client set
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to get rest config: %w", err)
	}

	targetClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to get client set: %w", err)
	}

	secretName := clusterScope.HivelocityCluster.Spec.HivelocitySecret.Name
	targetNS := metav1.NamespaceSystem
	sourceNS := clusterScope.HivelocityCluster.Namespace

	_, err = targetClientSet.CoreV1().Secrets(targetNS).Get(
		ctx,
		secretName,
		metav1.GetOptions{},
	)

	if err == nil {
		// Secret exists. Nothing to do.
		return nil
	}

	if !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	apiKeySecretName := types.NamespacedName{
		Namespace: sourceNS,
		Name:      secretName,
	}
	secretManager := secretutil.NewSecretManager(clusterScope.Logger, clusterScope.Client, clusterScope.APIReader)
	apiKeySecret, err := secretManager.AcquireSecret(ctx, apiKeySecretName, clusterScope.HivelocityCluster, false, clusterScope.HivelocityCluster.DeletionTimestamp.IsZero())
	if err != nil {
		return fmt.Errorf("failed to acquire secret: %w", err)
	}

	key := clusterScope.HivelocityCluster.Spec.HivelocitySecret.Key

	apiKey, keyExists := apiKeySecret.Data[key]
	if !keyExists {
		return fmt.Errorf(
			"error key %s does not exist in secret/%s: %w",
			key,
			apiKeySecretName,
			err,
		)
	}

	var immutable bool
	data := make(map[string][]byte)
	data[key] = apiKey

	// Save api server information
	data["apiserver-host"] = []byte(clusterScope.HivelocityCluster.Spec.ControlPlaneEndpoint.Host)
	data["apiserver-port"] = []byte(strconv.Itoa(int(clusterScope.HivelocityCluster.Spec.ControlPlaneEndpoint.Port)))

	newSecret := corev1.Secret{
		Immutable: &immutable,
		Data:      data,
		TypeMeta:  metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: targetNS,
		},
	}

	// create secret in cluster
	if _, err := targetClientSet.CoreV1().Secrets(targetNS).Create(ctx, &newSecret, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

func (r *HivelocityClusterReconciler) reconcileTargetClusterManager(ctx context.Context, clusterScope *scope.ClusterScope) (res reconcile.Result, err error) {
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
			conditions.MarkFalse(
				clusterScope.HivelocityCluster,
				infrav1.TargetClusterReadyCondition,
				infrav1.TargetClusterCreateFailedReason,
				clusterv1.ConditionSeverityError,
				err.Error(),
			)
			return res, fmt.Errorf("failed to create a clusterManager for HivelocityCluster %s/%s: %w",
				clusterScope.HivelocityCluster.Namespace,
				clusterScope.HivelocityCluster.Name,
				err,
			)
		}

		// manager could not be created yet - reconcile again
		if m == nil {
			return reconcile.Result{Requeue: true}, nil
		}

		r.targetClusterManagersStopCh[key] = make(chan struct{})

		ctx, cancel := context.WithCancel(ctx)

		r.TargetClusterManagersWaitGroup.Add(1)

		// Start manager
		go func() {
			defer r.TargetClusterManagersWaitGroup.Done()

			if err := m.Start(ctx); err != nil {
				clusterScope.Error(err, "failed to start a targetClusterManager")
				conditions.MarkFalse(
					clusterScope.HivelocityCluster,
					infrav1.TargetClusterReadyCondition,
					infrav1.TargetClusterCreateFailedReason,
					clusterv1.ConditionSeverityError,
					"failed to start a targetClusterManager: %s", err.Error(),
				)
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
	return res, nil
}

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
		if apierrors.IsNotFound(err) {
			conditions.MarkFalse(
				hvCluster,
				infrav1.TargetClusterReadyCondition,
				infrav1.KubeConfigNotFoundReason,
				clusterv1.ConditionSeverityInfo,
				"kubeconfig not found (yet)",
			)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get a clientConfig for the API of HivelocityCluster %s/%s: %w", hvCluster.Namespace, hvCluster.Name, err)
	}

	if err := scope.IsControlPlaneReady(ctx, clientConfig); err != nil {
		conditions.MarkFalse(
			hvCluster,
			infrav1.TargetClusterReadyCondition,
			infrav1.TargetClusterControlPlaneNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"target cluster not ready",
		)
		return nil, nil //nolint:nilerr
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

	httpClient, err := rest.HTTPClientFor(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get an HTTP client for the API of HivelocityCluster %s/%s: %w", hvCluster.Namespace, hvCluster.Name, err)
	}

	// Check whether kubeapi server responds
	if _, err := apiutil.NewDynamicRESTMapper(restConfig, httpClient); err != nil {
		conditions.MarkFalse(
			hvCluster,
			infrav1.TargetClusterReadyCondition,
			infrav1.KubeAPIServerNotRespondingReason,
			clusterv1.ConditionSeverityInfo,
			"kubeapi server not responding (yet)",
		)
		return nil, nil //nolint:nilerr
	}

	clusterMgr, err := ctrl.NewManager(
		restConfig,
		ctrl.Options{
			Scheme:         scheme,
			LeaderElection: false,
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
}

// SetupWithManager sets up the controller with the Manager.
func (r *HivelocityClusterReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := log.FromContext(ctx)

	if r.targetClusterManagersStopCh == nil {
		r.targetClusterManagersStopCh = make(map[types.NamespacedName]chan struct{})
	}

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
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, o client.Object) []reconcile.Request {
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
