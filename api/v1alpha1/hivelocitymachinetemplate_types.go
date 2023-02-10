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
)

// HivelocityMachineTemplateSpec defines the desired state of HivelocityMachineTemplate.
type HivelocityMachineTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HivelocityMachineTemplate. Edit hivelocitymachinetemplate_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// HivelocityMachineTemplateStatus defines the observed state of HivelocityMachineTemplate.
type HivelocityMachineTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HivelocityMachineTemplate is the Schema for the hivelocitymachinetemplates API.
type HivelocityMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HivelocityMachineTemplateSpec   `json:"spec,omitempty"`
	Status HivelocityMachineTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HivelocityMachineTemplateList contains a list of HivelocityMachineTemplate.
type HivelocityMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HivelocityMachineTemplate{}, &HivelocityMachineTemplateList{})
}
