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

// Package device implements functions to manage the lifecycle of Hivelocity devices.
package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Service defines struct with machine scope to reconcile Hivelocity machines.
type Service struct {
	scope *scope.MachineScope
}

const (
	maxShutDownTime  = 2 * time.Minute
	deviceOffTimeout = 10 * time.Minute
	defaultImageName = "Ubuntu 20.x"
)

var (
	errMachineTagNotFound = fmt.Errorf("machine tag not found")
	errClusterTagNotFound = fmt.Errorf("cluster tag not found")

	errNoDeviceAvailable = fmt.Errorf("no available device found")
)

// NewService outs a new service with machine scope.
func NewService(scope *scope.MachineScope) *Service {
	return &Service{
		scope: scope,
	}
}

// Reconcile implements reconcilement of Hivelocity machines.
func (s *Service) Reconcile(ctx context.Context) (_ ctrl.Result, err error) {
	log := s.scope.Logger

	// detect failure domain
	failureDomain, err := s.scope.GetFailureDomain()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get failure domain: %w", err)
	}
	s.scope.HivelocityMachine.Status.Region = infrav1.Region(failureDomain)

	// Waiting for bootstrap data to be ready
	if !s.scope.IsBootstrapDataReady(ctx) {
		log.Info("Bootstrap not ready - requeuing")
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.MachineBootstrapReadyCondition,
			infrav1.MachineBootstrapNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"bootstrap not ready yet",
		)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.MachineBootstrapReadyCondition,
	)

	initialState := s.scope.HivelocityMachine.Spec.Status.ProvisioningState

	hostStateMachine := newStateMachine(s.scope.HivelocityMachine, s)
	actResult := hostStateMachine.ReconcileState(ctx)
	result, err := actResult.Result()
	if err != nil {
		err = fmt.Errorf("action %q failed: %w", initialState, err)
		return ctrl.Result{Requeue: true}, err
	}

	return result, nil
}

// TODO: Implement logic to add multiple images.
func (s *Service) getDeviceImage(ctx context.Context) (string, error) {
	return defaultImageName, nil
}

func (s *Service) getSSHKeyIDFromSSHKeyName(ctx context.Context, sshKey *infrav1.SSHKey) (int32, error) {
	if sshKey == nil {
		return 0, fmt.Errorf("SSHKey is nil")
	}
	keys, err := s.scope.HVClient.ListSSHKeys(ctx)
	if err != nil {
		return 0, fmt.Errorf("[getSSHKeyIDFromSSHKeyName] ListSSHKeys() failed. Name %q: %w", sshKey.Name, err)
	}
	for _, key := range keys {
		if key.Name == sshKey.Name {
			return key.SshKeyId, nil
		}
	}
	return 0, fmt.Errorf("[getSSHKeyIDFromSSHKeyName] no corresponding ssh-key found in HV API. Name %q",
		sshKey.Name)
}

// Delete implements delete method of the HivelocityMachine.
func (s *Service) Delete(ctx context.Context) (_ *ctrl.Result, err error) {
	// If no providerID is set, there is nothing to do
	if s.scope.HivelocityMachine.Spec.ProviderID == nil {
		return nil, nil
	}

	// Check whether device still exists
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return nil, fmt.Errorf("[Delete] ProviderIDToDeviceID failed: %w", err)
	}

	_, err = s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// Nothing to do if device is not found
			s.scope.Info("Unable to locate Hivelocity device by ID or tags")
			record.Warnf(s.scope.HivelocityMachine, "NoDeviceFound", "Unable to find matching Hivelocity device for %s", s.scope.Name())
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// TODO: Finishing up Delete
	return nil, nil
}

// setMachineAddress gets the address from the device and sets it on the HivelocityMachine object.
func setMachineAddress(hvMachine *infrav1.HivelocityMachine, hvDevice *hv.BareMetalDevice) {
	hvMachine.Status.Addresses = []clusterv1.MachineAddress{
		{
			Type:    clusterv1.MachineHostName,
			Address: hvDevice.Hostname,
		},
		{
			Type:    clusterv1.MachineInternalIP,
			Address: hvDevice.PrimaryIp,
		},
		{
			Type:    clusterv1.MachineExternalIP,
			Address: hvDevice.PrimaryIp,
		},
	}
}

// chooseDevice searches for an unused device.
func (s *Service) chooseDevice(ctx context.Context) (hv.BareMetalDevice, error) {
	devices, err := s.scope.HVClient.ListDevices(ctx)
	if err != nil {
		return hv.BareMetalDevice{}, fmt.Errorf("[chooseDevice] ListDevices() failed. machine %q: %w",
			s.scope.Name(), err)
	}
	return chooseAvailableFromList(devices, s.scope.HivelocityMachine.Spec.Type, s.scope.HivelocityCluster.Name, s.scope.Name())
}

func chooseAvailableFromList(devices []hv.BareMetalDevice, deviceType infrav1.HivelocityDeviceType, clusterName, machineName string) (hv.BareMetalDevice, error) {
	for _, device := range devices {
		// Ignore if associated already
		machineTag, err := hvtag.MachineTagFromList(device.Tags)
		if err != nil && !errors.Is(err, hvtag.ErrDeviceTagNotFound) {
			// unexpected error - continue
			continue
		}
		if machineTag.Value == machineName {
			// associated to this machine - return
			return device, nil
		}
		if machineTag.Value != "" {
			// associated to other machine - continue
			continue
		}

		deviceTypeTag, err := hvtag.DeviceTypeTagFromList(device.Tags)
		if err != nil {
			// TODO: What if no device type has been set? Is the device available or not?
			continue
		}

		// Ignore if has wrong device type
		if deviceTypeTag.Value != string(deviceType) {
			continue
		}

		// Ignore if associated to other cluster
		clusterTag, err := hvtag.ClusterTagFromList(device.Tags)
		if err != nil && !errors.Is(err, hvtag.ErrDeviceTagNotFound) {
			// unexpected error - continue
			continue
		}
		if clusterTag.Value != "" && clusterTag.Value != clusterName {
			// associated to another cluster - continue
			continue
		}

		return device, nil
	}
	return hv.BareMetalDevice{}, errNoDeviceAvailable
}

// actionAssociateDevice claims an unused HV device by settings tags and returns it.
func (s *Service) actionAssociateDevice(ctx context.Context) actionResult {
	s.scope.Info("Start actionAssociateDevice")
	device, err := s.chooseDevice(ctx)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to choose device: %w", err)}
	}

	device.Tags = append(device.Tags,
		s.scope.HivelocityCluster.DeviceTag().ToString(),
		s.scope.HivelocityCluster.DeviceTagOwned().ToString(),
		s.scope.HivelocityMachine.DeviceTag().ToString(),
	)
	if err := s.scope.HVClient.SetDeviceTags(ctx, device.DeviceId, device.Tags); err != nil {
		return actionError{err: fmt.Errorf("failed to set tags: %w", err)}
	}
	providerID := providerIDFromDeviceID(device.DeviceId)
	s.scope.HivelocityMachine.Spec.ProviderID = &providerID
	return actionComplete{}
}

// actionVerifyAssociate verifies that the HV device has actually been associated to this machine and only this.
// Checking whether there are other machines also associated avoids situations where machines are selected at the same time.
func (s *Service) actionVerifyAssociate(ctx context.Context) actionResult {
	s.scope.Info("Start actionVerifyAssociate")
	// TODO: We should be able to control this value from outside. Do we do it with flags?
	// wait for 3 seconds at least before checking again
	const waitFor = 100 * time.Millisecond

	// if the waiting time has not yet passed, then we reconcile again without changing state
	if !hasTimedOut(s.scope.HivelocityMachine.Spec.Status.LastUpdated, waitFor) {
		return actionContinue{delay: 100 * time.Millisecond}
	}

	// if waiting time is over, we check the server for tags
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	// check if cluster and machine tags are properly set
	_, clusterTagErr := hvtag.ClusterTagFromList(device.Tags)
	_, machineTagErr := hvtag.MachineTagFromList(device.Tags)

	if clusterTagErr == nil && machineTagErr == nil {
		return actionComplete{}
	}

	// something is wrong - remove cluster and machine tags and associate a new device
	newTagList, updatedTags1 := s.scope.HivelocityCluster.DeviceTag().RemoveFromList(device.Tags)
	newTagList, updatedTags2 := s.scope.HivelocityMachine.DeviceTag().RemoveFromList(newTagList)
	if updatedTags1 || updatedTags2 {
		if err := s.scope.HVClient.SetDeviceTags(ctx, deviceID, newTagList); err != nil {
			return actionError{err: fmt.Errorf("failed to remove associated machine from tags: %w", err)}
		}
	}
	s.scope.HivelocityMachine.Spec.ProviderID = nil
	return actionError{err: errGoToPreviousState}
}

func hasTimedOut(lastUpdated *metav1.Time, timeout time.Duration) bool {
	if lastUpdated == nil {
		return false
	}
	now := metav1.Now()
	return lastUpdated.Add(timeout).Before(now.Time)
}

// actionEnsureDeviceShutDown ensures that the device is shut down.
func (s *Service) actionEnsureDeviceShutDown(ctx context.Context) actionResult {
	s.scope.Info("Start actionEnsureDeviceShutDown")

	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	err = s.scope.HVClient.ShutdownDevice(ctx, deviceID)
	// if device is already shut down, we can go to the next step
	if errors.Is(err, hvclient.ErrDeviceShutDownAlready) {
		return actionComplete{}
	}
	s.scope.Info("Device is not shut down yet", "error from ShutDown endpoint", err)

	// wait for another 10 seconds
	return actionContinue{delay: 10 * time.Second}
}

// actionProvisionDevice provisions the device.
func (s *Service) actionProvisionDevice(ctx context.Context) actionResult {
	s.scope.Info("Start actionProvisionDevice")
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	userData, err := s.scope.GetRawBootstrapData(ctx)
	if err != nil {
		record.Warnf(
			s.scope.HivelocityMachine,
			"FailedGetBootstrapData",
			err.Error(),
		)
		return actionError{err: fmt.Errorf("failed to get raw bootstrap data: %s", err)}
	}

	image, err := s.getDeviceImage(ctx)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get device image: %w", err)}
	}

	tags := []string{
		s.scope.HivelocityCluster.DeviceTag().ToString(),
		s.scope.HivelocityMachine.DeviceTag().ToString(),
		s.scope.DeviceTagMachineType().ToString(),
	}

	device.Tags = append(device.Tags, tags...)

	opts := hv.BareMetalDeviceUpdate{
		Hostname:    s.scope.Name(),
		Tags:        device.Tags,
		Script:      string(userData), // cloud-init script
		OsName:      image,
		ForceReload: true,
	}

	if s.scope.HivelocityCluster.Spec.SSHKey != nil {
		sshKeyID, err := s.getSSHKeyIDFromSSHKeyName(ctx, s.scope.HivelocityCluster.Spec.SSHKey)
		if err != nil {
			return actionError{err: fmt.Errorf("error with ssh keys: %w", err)}
		}
		opts.PublicSshKeyId = sshKeyID
	}

	// Provision the device
	provisionedDevice, err := s.scope.HVClient.ProvisionDevice(ctx, deviceID, opts)
	s.scope.Info("[actionProvisionDevice] ProvisionDevice was called", "err", err, "deviceID", deviceID)
	if err != nil {
		// TODO: Handle error that machine is not shut down
		if hvclient.IsRateLimitExceededError(err) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine,
				"RateLimitExceeded",
				"exceeded rate limit with calling Hivelocity function ProvisionDevice",
			)
		}
		record.Warnf(s.scope.HivelocityMachine,
			"FailedProvisionHivelocityDevice",
			"Failed to provision Hivelocity device %s: %s",
			s.scope.Name(),
			err,
		)
		return actionError{err: fmt.Errorf("failed to provision device %d: %s", deviceID, err)}
	}
	setMachineAddress(s.scope.HivelocityMachine, &provisionedDevice)

	s.scope.HivelocityMachine.Status.Ready = true
	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition)
	return actionComplete{}
}

// actionDeviceProvisioned reconciles a provisioned device.
func (s *Service) actionDeviceProvisioned(ctx context.Context) actionResult {
	// Check whether device still exists
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("[actionDeviceProvisioned] ProviderIDToDeviceID failed: %w", err)}
	}

	// FIXME: we already get the device
	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			conditions.MarkFalse(s.scope.HivelocityMachine,
				infrav1.DeviceReadyCondition,
				infrav1.DeviceNotFoundReason,
				clusterv1.ConditionSeverityError,
				fmt.Sprintf("device %d not found", device.DeviceId))
			// TODO: Return fatal error
		}
		return actionError{err: fmt.Errorf("failed to get associated device: %w", err)}
	}

	if err := s.verifyAssociatedDevice(&device); err != nil {
		// TODO: Fatal error
		return actionError{err: fmt.Errorf("associated device could not be verified")}
	}

	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition)
	setMachineAddress(s.scope.HivelocityMachine, &device)
	s.scope.HivelocityMachine.Status.Ready = true

	return actionComplete{}
}

func (s *Service) verifyAssociatedDevice(device *hv.BareMetalDevice) error {
	if !s.scope.HivelocityCluster.DeviceTag().IsInStringList(device.Tags) {
		return errClusterTagNotFound
	}
	if !s.scope.HivelocityMachine.DeviceTag().IsInStringList(device.Tags) {
		return errMachineTagNotFound
	}
	return nil
}

// providerIDFromDeviceID converts a deviceID to ProviderID.
func providerIDFromDeviceID(deviceID int32) string {
	return fmt.Sprintf("hivelocity://%d", deviceID)
}
