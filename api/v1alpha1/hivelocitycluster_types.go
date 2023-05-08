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

package v1alpha1

import (
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ClusterFinalizer allows ReconcileHivelocityCluster to clean up Hivelocity
	// resources associated with HivelocityCluster before removing it from the
	// apiserver.
	ClusterFinalizer = "hivelocitycluster.infrastructure.cluster.x-k8s.io"
)

// HivelocityClusterSpec defines the desired state of HivelocityCluster.
type HivelocityClusterSpec struct {
	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint *clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ControlPlaneRegion is a Hivelocity Region (LAX2, ...).
	ControlPlaneRegion Region `json:"controlPlaneRegion"`

	// HivelocitySecret is a reference to a Kubernetes Secret.
	HivelocitySecret HivelocitySecretRef `json:"hivelocitySecretRef"`

	// SSHKey is cluster wide. Valid value is a valid SSH key name.
	// +optional
	SSHKey *SSHKey `json:"sshKey,omitempty"`
}

// HivelocitySecretRef defines the name of the Secret and the relevant key in the secret to access the Hivelocity API.
type HivelocitySecretRef struct {
	// +optional
	// +kubebuilder:default=hivelocity
	Name string `json:"name,omitempty"`

	// +optional
	// +kubebuilder:default=HIVELOCITY_API_KEY
	Key string `json:"key,omitempty"`
}

// SSHKey defines the SSHKey for Hivelocity.
type SSHKey struct {
	// Name of SSH key.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// HivelocityClusterStatus defines the observed state of HivelocityCluster.
type HivelocityClusterStatus struct {
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=hivelocityclusters,scope=Namespaced,categories=cluster-api,shortName=capihvc
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this HivelocityCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for Nodes"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Region",type="string",JSONPath=".spec.controlPlaneRegion",description="Control plane region"
// +k8s:defaulter-gen=true

// HivelocityCluster is the Schema for the hivelocityclusters API.
type HivelocityCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HivelocityClusterSpec   `json:"spec,omitempty"`
	Status HivelocityClusterStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the HivelocityCluster resource.
func (r *HivelocityCluster) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the HivelocityCluster to the predescribed clusterv1.Conditions.
func (r *HivelocityCluster) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

// DeviceTag returns a DeviceTag object for the cluster tag.
func (r *HivelocityCluster) DeviceTag() hvtag.DeviceTag {
	return hvtag.DeviceTag{
		Key:   hvtag.DeviceTagKeyCluster,
		Value: r.Name,
	}
}

// DeviceTagOwned returns a DeviceTag object for the ResourceLifeCycle tag.
func (r *HivelocityCluster) DeviceTagOwned() hvtag.DeviceTag {
	return hvtag.DeviceTag{
		Key:   hvtag.DeviceTagKey(ClusterTagKey(r.Name)),
		Value: string(ResourceLifecycleOwned),
	}
}

// SetStatusFailureDomain sets the region for the status.
func (r *HivelocityCluster) SetStatusFailureDomain(region Region) {
	if r.Status.FailureDomains == nil {
		r.Status.FailureDomains = make(clusterv1.FailureDomains, 1)
	}

	r.Status.FailureDomains[string(region)] = clusterv1.FailureDomainSpec{
		ControlPlane: true,
	}
}

// +kubebuilder:object:root=true
// +k8s:defaulter-gen=true

// HivelocityClusterList contains a list of HivelocityCluster.
type HivelocityClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HivelocityCluster{}, &HivelocityClusterList{})
}
