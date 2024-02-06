/*
Copyright 2022 The Kubernetes Authors.

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

// Package scope defines cluster and machine scope as well as a repository for the Hivelocity API.
package scope

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams defines the input parameters used to create a new scope.
type ClusterScopeParams struct {
	Client            client.Client
	APIReader         client.Reader
	Logger            logr.Logger
	HVClient          hvclient.Client
	Cluster           *clusterv1.Cluster
	HivelocityCluster *infrav1.HivelocityCluster
}

// NewClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewClusterScope(_ context.Context, params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.HivelocityCluster == nil {
		return nil, errors.New("failed to generate new scope from nil HivelocityCluster")
	}
	if params.HVClient == nil {
		return nil, errors.New("failed to generate new scope from nil HVClient")
	}

	helper, err := patch.NewHelper(params.HivelocityCluster, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}

	return &ClusterScope{
		Logger:            params.Logger,
		Client:            params.Client,
		APIReader:         params.APIReader,
		Cluster:           params.Cluster,
		HivelocityCluster: params.HivelocityCluster,
		HVClient:          params.HVClient,
		patchHelper:       helper,
	}, nil
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	logr.Logger
	Client      client.Client
	APIReader   client.Reader
	patchHelper *patch.Helper
	HVClient    hvclient.Client

	Cluster           *clusterv1.Cluster
	HivelocityCluster *infrav1.HivelocityCluster
}

// Name returns the HivelocityCluster name.
func (s *ClusterScope) Name() string {
	return s.HivelocityCluster.Name
}

// Namespace returns the namespace name.
func (s *ClusterScope) Namespace() string {
	return s.HivelocityCluster.Namespace
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close(ctx context.Context) error {
	return s.patchHelper.Patch(ctx, s.HivelocityCluster)
}

// PatchObject persists the machine spec and status.
func (s *ClusterScope) PatchObject(ctx context.Context) error {
	return s.patchHelper.Patch(ctx, s.HivelocityCluster)
}

// ClientConfig return a kusecretbernetes client config for the cluster context.
func (s *ClusterScope) ClientConfig(ctx context.Context) (clientcmd.ClientConfig, error) {
	cluster := client.ObjectKey{
		Name:      fmt.Sprintf("%s-%s", s.Cluster.Name, secret.Kubeconfig),
		Namespace: s.Cluster.Namespace,
	}

	var kubeConfigSecret corev1.Secret
	if err := s.Client.Get(ctx, cluster, &kubeConfigSecret); err != nil {
		return nil, fmt.Errorf("failed to find kube config secret: %w", err)
	}

	kubeconfigBytes, ok := kubeConfigSecret.Data[secret.KubeconfigDataName]
	if !ok {
		return nil, errors.Errorf("missing key %q in secret data", secret.KubeconfigDataName)
	}
	return clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
}

// ClientConfigWithAPIEndpoint returns a client config.
func (s *ClusterScope) ClientConfigWithAPIEndpoint(ctx context.Context, endpoint clusterv1.APIEndpoint) (clientcmd.ClientConfig, error) {
	c, err := s.ClientConfig(ctx)
	if err != nil {
		return nil, err
	}

	raw, err := c.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("error retrieving rawConfig from clientConfig: %w", err)
	}
	// update cluster endpint in config
	for key := range raw.Clusters {
		raw.Clusters[key].Server = fmt.Sprintf("https://%s:%d", endpoint.Host, endpoint.Port)
	}

	return clientcmd.NewDefaultClientConfig(raw, &clientcmd.ConfigOverrides{}), nil
}

// ListMachines returns HivelocityMachines.
func (s *ClusterScope) ListMachines(ctx context.Context) ([]*clusterv1.Machine, []*infrav1.HivelocityMachine, error) {
	// get and index Machines by HivelocityMachine name
	var machineListRaw clusterv1.MachineList
	machineByHivelocityMachineName := make(map[string]*clusterv1.Machine)
	if err := s.Client.List(ctx, &machineListRaw, client.InNamespace(s.Namespace())); err != nil {
		return nil, nil, err
	}
	expectedGK := infrav1.GroupVersion.WithKind("HivelocityMachine").GroupKind()
	for pos := range machineListRaw.Items {
		m := &machineListRaw.Items[pos]
		actualGK := m.Spec.InfrastructureRef.GroupVersionKind().GroupKind()
		if m.Spec.ClusterName != s.Cluster.Name ||
			actualGK.String() != expectedGK.String() {
			continue
		}
		machineByHivelocityMachineName[m.Spec.InfrastructureRef.Name] = m
	}

	// match HivelocityMachines to Machines
	var hivelocityMachineListRaw infrav1.HivelocityMachineList
	if err := s.Client.List(ctx, &hivelocityMachineListRaw, client.InNamespace(s.Namespace())); err != nil {
		return nil, nil, err
	}

	machineList := make([]*clusterv1.Machine, 0, len(hivelocityMachineListRaw.Items))
	hivelocityMachineList := make([]*infrav1.HivelocityMachine, 0, len(hivelocityMachineListRaw.Items))

	for pos := range hivelocityMachineListRaw.Items {
		hm := &hivelocityMachineListRaw.Items[pos]
		m, ok := machineByHivelocityMachineName[hm.Name]
		if !ok {
			continue
		}

		machineList = append(machineList, m)
		hivelocityMachineList = append(hivelocityMachineList, hm)
	}

	return machineList, hivelocityMachineList, nil
}

// ErrWorkloadControlPlaneNotReady indicates that the control plane is not ready (or not reachable).
var ErrWorkloadControlPlaneNotReady = errors.New("Workload ControlPlane not Ready (or not reachable)")

// IsControlPlaneReady returns nil if the control plane is ready.
func IsControlPlaneReady(ctx context.Context, c clientcmd.ClientConfig) error {
	restConfig, err := c.ClientConfig()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkloadControlPlaneNotReady, err)
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkloadControlPlaneNotReady, err)
	}

	_, err = clientSet.Discovery().RESTClient().Get().AbsPath("/readyz").DoRaw(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkloadControlPlaneNotReady, err)
	}
	return nil
}
