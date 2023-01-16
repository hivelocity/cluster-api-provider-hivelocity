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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ClusterFinalizer allows ReconcileHivelocityCluster to clean up Hivelocity
	// resources associated with HivelocityCluster before removing it from the
	// apiserver.
	ClusterFinalizer = "hivelocitycluster.infrastructure.cluster.x-k8s.io"
)

// HivelocityClusterSpec defines the desired state of HivelocityCluster
type HivelocityClusterSpec struct {
	// Important: Run "make generate" to regenerate code after modifying this file

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint *clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// HivelocitySecret is a reference to a Kubernetes Secret.
	HivelocitySecret HivelocitySecretRef `json:"hivelocitySecretRef"`

	// SSHKey is cluster wide. Valid value is a valid SSH key name.
	SSHKey *SSHKey `json:"sshKey"`
}

// HivelocityClusterStatus defines the observed state of HivelocityCluster
type HivelocityClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HivelocityCluster is the Schema for the hivelocityclusters API
type HivelocityCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HivelocityClusterSpec   `json:"spec,omitempty"`
	Status HivelocityClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HivelocityClusterList contains a list of HivelocityCluster
type HivelocityClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityCluster `json:"items"`
}

// GetConditions returns the observations of the operational state of the HivelocityCluster resource.
func (r *HivelocityCluster) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the HivelocityCluster to the predescribed clusterv1.Conditions.
func (r *HivelocityCluster) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&HivelocityCluster{}, &HivelocityClusterList{})
}
