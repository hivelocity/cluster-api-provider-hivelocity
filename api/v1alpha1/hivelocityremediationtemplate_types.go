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

// HivelocityRemediationTemplateSpec defines the desired state of HivelocityRemediationTemplate
type HivelocityRemediationTemplateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of HivelocityRemediationTemplate. Edit hivelocityremediationtemplate_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// HivelocityRemediationTemplateStatus defines the observed state of HivelocityRemediationTemplate
type HivelocityRemediationTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HivelocityRemediationTemplate is the Schema for the hivelocityremediationtemplates API
type HivelocityRemediationTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HivelocityRemediationTemplateSpec   `json:"spec,omitempty"`
	Status HivelocityRemediationTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HivelocityRemediationTemplateList contains a list of HivelocityRemediationTemplate
type HivelocityRemediationTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityRemediationTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HivelocityRemediationTemplate{}, &HivelocityRemediationTemplateList{})
}
