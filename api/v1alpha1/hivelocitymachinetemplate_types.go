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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// HivelocityMachineTemplateSpec defines the desired state of HivelocityMachineTemplate.
type HivelocityMachineTemplateSpec struct {
	Template HivelocityMachineTemplateResource `json:"template"`
}

// HivelocityMachineTemplateStatus defines the observed state of HivelocityMachineTemplate.
type HivelocityMachineTemplateStatus struct {
	// Capacity defines the resource capacity for this machine.
	// This value is used for autoscaling from zero operations as defined in:
	// https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
	// +optional
	Capacity corev1.ResourceList `json:"capacity,omitempty"`

	// Conditions defines current service state of the HivelocityMachineTemplate.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hivelocitymachinetemplates,scope=Namespaced,categories=cluster-api,shortName=capihvcmt
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.template.spec.imageName",description="Image name"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.template.spec.type",description="Server type"
// +kubebuilder:storageversion
// +k8s:defaulter-gen=true

// HivelocityMachineTemplate is the Schema for the hivelocitymachinetemplates API.
type HivelocityMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HivelocityMachineTemplateSpec   `json:"spec,omitempty"`
	Status HivelocityMachineTemplateStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the HivelocityMachine resource.
func (r *HivelocityMachineTemplate) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the HivelocityMachine to the predescribed clusterv1.Conditions.
func (r *HivelocityMachineTemplate) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// HivelocityMachineTemplateList contains a list of HivelocityMachineTemplate.
type HivelocityMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityMachineTemplate `json:"items"`
}

// HivelocityMachineTemplateResource describes the data needed to create am HivelocityMachine from a template.
type HivelocityMachineTemplateResource struct {
	// Standard object's metadata.
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the machine.
	Spec HivelocityMachineSpec `json:"spec"`
}

func init() {
	SchemeBuilder.Register(&HivelocityMachineTemplate{}, &HivelocityMachineTemplateList{})
}
