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
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows ReconcileHivelocityMachine to clean up Hivelocity
	// resources associated with HivelocityMachine before removing it from the
	// apiserver.
	MachineFinalizer = "hivelocitymachine.infrastructure.cluster.x-k8s.io"
)

// HivelocityMachineSpec defines the desired state of HivelocityMachine.
type HivelocityMachineSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// Type is the Hivelocity Machine Type for this machine.
	// +kubebuilder:validation:Enum=hvCustom;todo-question
	Type HivelocityMachineType `json:"type"`

	// ImageName is the reference to the Machine Image from which to create the device.
	// +kubebuilder:validation:MinLength=1
	ImageName string `json:"imageName"`
}

// HivelocityMachineStatus defines the observed state of HivelocityMachine.
type HivelocityMachineStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the devices's associated addresses.
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// Region contains the name of the Hivelocity location the device is running.
	Region Region `json:"region,omitempty"`

	// DeviceState is the state of the device for this machine.
	// +optional
	DeviceState *hvclient.DeviceStatus `json:"deviceState,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
	// Conditions defines current service state of the HivelocityMachine.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=hivelocitymachines,scope=Namespaced,categories=cluster-api,shortName=capihcm
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this HivelocityMachine belongs"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.imageName",description="Image name"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type",description="Device type"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.deviceState",description="Hivelocity device state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="DeviceID",type="string",JSONPath=".spec.providerID",description="Hivelocity device ID"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this HivelocityMachine"
// +k8s:defaulter-gen=true

// HivelocityMachine is the Schema for the hivelocitymachines API.
type HivelocityMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HivelocityMachineSpec   `json:"spec,omitempty"`
	Status HivelocityMachineStatus `json:"status,omitempty"`
}

// HivelocityMachineSpec returns a DeepCopy.
func (r *HivelocityMachine) HivelocityMachineSpec() *HivelocityMachineSpec {
	return r.Spec.DeepCopy()
}

// GetConditions returns the observations of the operational state of the HivelocityMachine resource.
func (r *HivelocityMachine) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the HivelocityMachine to the predescribed clusterv1.Conditions.
func (r *HivelocityMachine) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// HivelocityMachineList contains a list of HivelocityMachine.
type HivelocityMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HivelocityMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HivelocityMachine{}, &HivelocityMachineList{})
}
