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

// HivelocityClusterTemplateSpec defines the desired state of HivelocityClusterTemplate.
type HivelocityClusterTemplateSpec struct {
	Template HivelocityClusterTemplateResource `json:"template"`
}

// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=hivelocityclustertemplates,scope=Namespaced,categories=cluster-api,shortName=capihvct
// +k8s:defaulter-gen=true

// HivelocityClusterTemplate is the Schema for the hivelocityclustertemplates API.
type HivelocityClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HivelocityClusterTemplateSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// HivelocityClusterTemplateList contains a list of HivelocityClusterTemplate.
type HivelocityClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HivelocityClusterTemplate{}, &HivelocityClusterTemplateList{})
}

// HivelocityClusterTemplateResource contains spec for HivelocityClusterSpec.
type HivelocityClusterTemplateResource struct {
	// +optional
	ObjectMeta clusterv1.ObjectMeta  `json:"metadata,omitempty"`
	Spec       HivelocityClusterSpec `json:"spec"`
}
