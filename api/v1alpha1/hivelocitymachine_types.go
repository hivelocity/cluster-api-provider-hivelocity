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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

const (
	// MachineFinalizer allows ReconcileHivelocityMachine to clean up Hivelocity
	// resources associated with HivelocityMachine before removing it from the
	// apiserver.
	MachineFinalizer = "hivelocitymachine.infrastructure.cluster.x-k8s.io"
)

const (
	// FailureMessageDeviceNotFound indicates that the associated device could not be found.
	FailureMessageDeviceNotFound = "device not found"

	// FailureMessageDeviceTagsInvalid indicates that the associated device has invalid tags.
	// This is probably due to a user changing device tags on his own.
	FailureMessageDeviceTagsInvalid = "device tags invalid"
)

var (
	// ErrEmptyProviderID indicates an empty providerID.
	ErrEmptyProviderID = fmt.Errorf("providerID is empty")
	// ErrInvalidProviderID indicates an invalid providerID.
	ErrInvalidProviderID = fmt.Errorf("providerID is invalid")
	// ErrInvalidDeviceID indicates an invalid deviceID.
	ErrInvalidDeviceID = fmt.Errorf("deviceID is invalid")
)

// ProvisioningState defines the states the provisioner will report the host has having.
type ProvisioningState string

const (
	// StateNone means the state is unknown.
	StateNone ProvisioningState = ""

	// StateAssociateDevice .
	StateAssociateDevice ProvisioningState = "associate-device"

	// StateVerifyAssociate .
	StateVerifyAssociate ProvisioningState = "verify-associate"

	// StateVerifyShutdown .
	StateVerifyShutdown ProvisioningState = "verify-shutdown"

	// StateProvisionDevice .
	StateProvisionDevice ProvisioningState = "provision-device"

	// StateDeviceProvisioned .
	StateDeviceProvisioned ProvisioningState = "provisioned"

	// StateDeleteDeviceDeProvision .
	StateDeleteDeviceDeProvision ProvisioningState = "delete-deprovision"

	// StateDeleteDeviceDissociate .
	StateDeleteDeviceDissociate ProvisioningState = "delete-dissociate"

	// StateDeleteDevice .
	StateDeleteDevice ProvisioningState = "delete"
)

// HivelocityMachineSpec defines the desired state of HivelocityMachine.
type HivelocityMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// DeviceSelector can be used to limit the set of devices that this HivelocityMachine can claim.
	// +optional
	DeviceSelector DeviceSelector `json:"deviceSelector,omitempty"`

	// ImageName is the reference to the Machine Image from which to create the device.
	// +kubebuilder:validation:MinLength=1
	ImageName string `json:"imageName"`

	// Status contains all status information of the controller. Do not edit these values!
	// +optional
	Status ControllerGeneratedStatus `json:"status,omitempty"`
}

// DeviceSelector specifies matching criteria for tags on devices.
// This is used to target a specific set of devices that can be claimed by the HivelocityMachine.
type DeviceSelector struct {
	// Key/value pairs of labels that must exist on a chosen Device
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`

	// MatchExpressions match expressions that must be true on a chosen Device
	// +optional
	MatchExpressions []DeviceSelectorRequirement `json:"matchExpressions,omitempty"`
}

// Validate validates the deviceSelector.
func (deviceSelector *DeviceSelector) Validate() error {
	_, err := deviceSelector.GetLabelSelector()
	return err
}

// GetLabelSelector returns the labels.Selector for the deviceSelector.
func (deviceSelector *DeviceSelector) GetLabelSelector() (labels.Selector, error) {
	labelSelector := labels.NewSelector()
	var reqs labels.Requirements

	var errs []error
	for labelKey, labelVal := range deviceSelector.MatchLabels {
		r, err := labels.NewRequirement(labelKey, selection.Equals, []string{labelVal})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		reqs = append(reqs, *r)
	}
	for _, req := range deviceSelector.MatchExpressions {
		lowercaseOperator := selection.Operator(strings.ToLower(string(req.Operator)))
		r, err := labels.NewRequirement(req.Key, lowercaseOperator, req.Values)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		reqs = append(reqs, *r)
	}

	return labelSelector.Add(reqs...), errors.Join(errs...)
}

// DeviceSelectorRequirement defines a requirement used for MatchExpressions to select device.
type DeviceSelectorRequirement struct {
	Key      string             `json:"key"`
	Operator selection.Operator `json:"operator"`
	Values   []string           `json:"values"`
}

// ControllerGeneratedStatus contains all status information which is important to persist.
type ControllerGeneratedStatus struct {
	// Information tracked by the provisioner.
	// +optional
	ProvisioningState ProvisioningState `json:"provisioningState"`

	// Time stamp of last update of status.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`
}

// HivelocityDeviceType defines the Hivelocity device type.
// +kubebuilder:validation:Enum=pool;hvCustom;hvControlPlane;hvWorker;e2eControlPlane;e2eWorker
// hvlabel:foo=bar
type HivelocityDeviceType string // TODO: this should not be an enum. Rename to HVLabel, and make a label selector.

// HivelocityMachineStatus defines the observed state of HivelocityMachine.
type HivelocityMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the machine's associated addresses.
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// Region contains the name of the Hivelocity location the device is running.
	Region Region `json:"region,omitempty"`

	// PowerState is the power state of the device for this machine (ON|OFF).
	// +optional
	PowerState string `json:"powerState,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	// +optional
	FailureReason *capierrors.MachineStatusError `json:"failureReason,omitempty"`

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
// +kubebuilder:resource:path=hivelocitymachines,scope=Namespaced,categories=cluster-api,shortName=capihvm
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this HivelocityMachine belongs"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.imageName",description="Image name"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type",description="Device type"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.powerState",description="Hivelocity device state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Machine ready status"
// +kubebuilder:printcolumn:name="ProviderID",type="string",JSONPath=".spec.providerID",description="ProviderID of machine object"
// +kubebuilder:printcolumn:name="Machine",type="string",JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object which owns with this HivelocityMachine"
// +kubebuilder:printcolumn:name="Prov.State",type="string",JSONPath=".spec.status.provisioningState"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
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

// SetFailure sets a failure reason and message.
func (r *HivelocityMachine) SetFailure(reason capierrors.MachineStatusError, message string) {
	r.Status.FailureReason = &reason
	r.Status.FailureMessage = &message
}

// SetProviderID sets the providerID based on a deviceID.
func (r *HivelocityMachine) SetProviderID(deviceID int32) {
	providerID := providerIDFromDeviceID(deviceID)
	r.Spec.ProviderID = &providerID
}

// SetMachineStatus sets the providerID based on a deviceID.
func (r *HivelocityMachine) SetMachineStatus(device hv.BareMetalDevice) {
	r.Status.Addresses = []clusterv1.MachineAddress{
		{
			Type:    clusterv1.MachineHostName,
			Address: device.Hostname,
		},
		{
			Type:    clusterv1.MachineInternalIP,
			Address: device.PrimaryIp,
		},
		{
			Type:    clusterv1.MachineExternalIP,
			Address: device.PrimaryIp,
		},
	}
	r.Status.PowerState = device.PowerStatus
	r.Status.Region = Region(device.LocationName)
}

// providerIDFromDeviceID converts a deviceID to ProviderID.
func providerIDFromDeviceID(deviceID int32) string {
	return fmt.Sprintf("hivelocity://%d", deviceID)
}

// DeviceTag returns a DeviceTag object for the machine tag.
func (r *HivelocityMachine) DeviceTag() hvtag.DeviceTag {
	return hvtag.DeviceTag{
		Key:   hvtag.DeviceTagKeyMachine,
		Value: r.Name,
	}
}

// DeviceIDFromProviderID converts the ProviderID (hivelocity://NNNN) to the DeviceID.
func (r *HivelocityMachine) DeviceIDFromProviderID() (int32, error) {
	if r.Spec.ProviderID == nil || r.Spec.ProviderID != nil && *r.Spec.ProviderID == "" {
		return 0, ErrEmptyProviderID
	}
	prefix := "hivelocity://"
	if !strings.HasPrefix(*r.Spec.ProviderID, prefix) {
		return 0, ErrInvalidProviderID
	}

	deviceID, err := strconv.ParseInt(strings.TrimPrefix(*r.Spec.ProviderID, prefix), 10, 32)
	if err != nil {
		return 0, ErrInvalidDeviceID
	}
	return int32(deviceID), nil
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
