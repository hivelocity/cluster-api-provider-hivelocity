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

// SSHKey defines the SSHKey for Hivelocity.
type SSHKey struct {
	// Name of SSH key
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Fingerprint of SSH key - added by controller // question: by controller? I thought by command-line-tool
	// +optional
	Fingerprint string `json:"fingerprint,omitempty"`
}

// HivelocityMachineType defines the Hivelocity Machine type.
type HivelocityMachineType string // question: rename to HivelocityDeviceType ?

// ResourceLifecycle configures the lifecycle of a resource.
type ResourceLifecycle string

// HivelocitySecretRef defines the name of the Secret and the relevant key in the secret to access the Hivelocity API.
type HivelocitySecretRef struct {
	// +optional
	// +kubebuilder:default=hivelocity
	Name string `json:"name,omitempty"`

	// +optional
	// +kubebuilder:default=HIVELOCITY_API_KEY
	Key string `json:"key,omitempty"`
}

// PublicNetworkSpec contains specs about public network spec of an Hivelocity device.
type PublicNetworkSpec struct {
	// +optional
	// +kubebuilder:default=true
	EnableIPv4 bool `json:"enableIPv4"`
	// +optional
	// +kubebuilder:default=true
	EnableIPv6 bool `json:"enableIPv6"`
}

// Region is a Hivelocity Location
// +kubebuilder:validation:Enum=TODO1;TODO2
type Region string
