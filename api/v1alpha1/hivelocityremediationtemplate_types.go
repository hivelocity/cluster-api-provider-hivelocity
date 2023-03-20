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

// HivelocityRemediationTemplateSpec defines the desired state of HivelocityRemediationTemplate.
type HivelocityRemediationTemplateSpec struct {
	Template HivelocityRemediationTemplateResource `json:"template"`
}

// HivelocityRemediationTemplateResource describes the data needed to create a HivelocityRemediation from a template.
type HivelocityRemediationTemplateResource struct {
	// Spec is the specification of the desired behavior of the HivelocityRemediation.
	Spec HivelocityRemediationSpec `json:"spec"`
}

// HivelocityRemediationTemplateStatus defines the observed state of HivelocityRemediationTemplate.
type HivelocityRemediationTemplateStatus struct {
	// HivelocityRemediationStatus defines the observed state of HivelocityRemediation
	Status HivelocityRemediationStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path= hivelocityremediationtemplates,scope=Namespaced,categories=cluster-api,shortName=hvrt;hvremediationtemplate;hvremediationtemplates; hivelocityrt; hivelocityremediationtemplate
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Strategy",type=string,JSONPath=".spec.template.spec.strategy.type",description="Type of the remediation strategy"
// +kubebuilder:printcolumn:name="Retry limit",type=string,JSONPath=".spec.template.spec.strategy.retryLimit",description="How many times remediation controller should attempt to remediate the host"
// +kubebuilder:printcolumn:name="Timeout",type=string,JSONPath=".spec.template.spec.strategy.timeout",description="Timeout for the remediation"

// HivelocityRemediationTemplate is the Schema for the  hivelocityremediationtemplates API.
type HivelocityRemediationTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec HivelocityRemediationTemplateSpec `json:"spec,omitempty"`
	// +optional
	Status HivelocityRemediationTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HivelocityRemediationTemplateList contains a list of HivelocityRemediationTemplate.
type HivelocityRemediationTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityRemediationTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HivelocityRemediationTemplate{}, &HivelocityRemediationTemplateList{})
}
