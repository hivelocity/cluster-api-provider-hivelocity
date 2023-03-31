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
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Service defines struct with machine scope to reconcile Hivelocity devices.
type Service struct {
	scope *scope.MachineScope
}

const (
	defaultImageName = "Ubuntu 20.x"
)

var (
	errMachineTagNotFound = fmt.Errorf("machine tag not found")
	errClusterTagNotFound = fmt.Errorf("cluster tag not found")

	errSSHKeyNotFound = fmt.Errorf("ssh key not found")

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
	// detect failure domain
	failureDomain, err := s.scope.GetFailureDomain()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get failure domain: %w", err)
	}
	s.scope.HivelocityMachine.Status.Region = infrav1.Region(failureDomain)

	// Waiting for bootstrap data to be ready
	if !s.scope.IsBootstrapDataReady(ctx) {
		s.scope.Info("Bootstrap not ready - requeuing")
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.MachineBootstrapReadyCondition,
			infrav1.MachineBootstrapNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"bootstrap not ready yet",
		)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.MachineBootstrapReadyCondition)

	initialState := s.scope.HivelocityMachine.Spec.Status.ProvisioningState

	hostStateMachine := newStateMachine(s.scope.HivelocityMachine, s)
	actResult := hostStateMachine.ReconcileState(ctx)
	result, err := actResult.Result()
	if err != nil {
		return ctrl.Result{Requeue: true}, fmt.Errorf("action %q failed: %w", initialState, err)
	}

	// TODO: Verify that patching the object is fine. Alternative would be to update it if that is better since we update the Spec as well
	return result, nil
}

// actionAssociateDevice claims an unused HV device by settings tags and returns it.
func (s *Service) actionAssociateDevice(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionAssociateDevice")
	log.V(1).Info("Started function")

	// list all devices
	devices, err := s.scope.HVClient.ListDevices(ctx)
	if err != nil {
		if errors.Is(err, hvclient.ErrRateLimitExceeded) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function ListDevices")
		}
		return actionError{err: fmt.Errorf("failed to list devices: %w", err)}
	}

	// find available device
	device, err := findAvailableDeviceFromList(devices, s.scope.HivelocityMachine.Spec.Type, s.scope.HivelocityCluster.Name, s.scope.Name())
	if err != nil {
		return actionError{err: fmt.Errorf("failed to find available device: %w", err)}
	}

	// associate this device with the machine object by setting tags
	device.Tags = append(device.Tags,
		s.scope.HivelocityCluster.DeviceTag().ToString(),
		s.scope.HivelocityCluster.DeviceTagOwned().ToString(),
		s.scope.HivelocityMachine.DeviceTag().ToString(),
	)

	if err := s.scope.HVClient.SetDeviceTags(ctx, device.DeviceId, device.Tags); err != nil {
		if errors.Is(err, hvclient.ErrRateLimitExceeded) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function SetDeviceTags")
		}
		return actionError{err: fmt.Errorf("failed to set tags on device %v: %w", device.DeviceId, err)}
	}

	// set providerID on machine object which is based on deviceID
	s.scope.HivelocityMachine.SetProviderID(device.DeviceId)

	log.V(1).Info("Completed function")
	return actionComplete{}
}

func findAvailableDeviceFromList(devices []hv.BareMetalDevice, deviceType infrav1.HivelocityDeviceType, clusterName, machineName string) (hv.BareMetalDevice, error) {
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

// actionVerifyAssociate verifies that the HV device has actually been associated to this machine and only this.
// Checking whether there are other machines also associated avoids situations where machines are selected at the same time.
func (s *Service) actionVerifyAssociate(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionVerifyAssociate")
	log.V(1).Info("Started function")

	// TODO: We should be able to control this value from outside. Do we do it with flags?
	// wait for 3 seconds at least before checking again
	const waitFor = 100 * time.Millisecond

	// if the waiting time has not yet passed, then we reconcile again without changing state
	if !hasTimedOut(s.scope.HivelocityMachine.Spec.Status.LastUpdated, waitFor) {
		return actionContinue{delay: 100 * time.Millisecond}
	}

	// if waiting time is over, we get the device and check its tags
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err)
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// if device cannot be found, we associate a new one
			log.Info("Device not found. Go back to StateAssociateDevice")
			record.Warnf(s.scope.HivelocityMachine, "DeviceNotFound", "Hivelocity device not found. Associate new one")
			return actionGoBack{nextState: infrav1.StateAssociateDevice}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	// check if cluster and machine tags are properly set
	_, clusterTagErr := hvtag.ClusterTagFromList(device.Tags)
	_, machineTagErr := hvtag.MachineTagFromList(device.Tags)

	if clusterTagErr == nil && machineTagErr == nil {
		log.V(1).Info("Completed function")
		return actionComplete{}
	}

	// Tags are not properly set or another machine also set its tags.
	// Remove cluster and machine tags and associate a new device.
	newTagList, updatedTags1 := s.scope.HivelocityCluster.DeviceTag().RemoveFromList(device.Tags)
	newTagList, updatedTags2 := s.scope.HivelocityMachine.DeviceTag().RemoveFromList(newTagList)
	if updatedTags1 || updatedTags2 {
		if err := s.scope.HVClient.SetDeviceTags(ctx, deviceID, newTagList); err != nil {
			if errors.Is(err, hvclient.ErrRateLimitExceeded) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function SetDeviceTags")
			}
			return actionError{err: fmt.Errorf("failed to remove associated machine from tags: %w", err)}
		}
	}
	s.scope.HivelocityMachine.Spec.ProviderID = nil

	log.Info("Device has been dissociated. Go back to StateAssociateDevice")
	return actionGoBack{nextState: infrav1.StateAssociateDevice}
}

func hasTimedOut(lastUpdated *metav1.Time, timeout time.Duration) bool {
	if lastUpdated == nil {
		return false
	}
	now := metav1.Now()
	return lastUpdated.Add(timeout).Before(now.Time)
}

// actionEnsureDeviceShutDown ensures that the device is shut down.
// This happens through repeatedly calling ShutdownDevice, as the potential error messages of that endpoint
// (e.g. 'device shut down already') are the most reliable source.
func (s *Service) actionEnsureDeviceShutDown(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionEnsureDeviceShutDown")
	log.V(1).Info("Started function")

	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	err = s.scope.HVClient.ShutdownDevice(ctx, deviceID)
	if errors.Is(err, hvclient.ErrRateLimitExceeded) {
		conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
		record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function ShutdownDevice")
	}
	// if device is already shut down, we can go to the next step
	if errors.Is(err, hvclient.ErrDeviceShutDownAlready) {
		log.V(1).Info("Completed function")
		return actionComplete{}
	}

	// device is not shut down yet - wait and call the function again
	log.Info("Device is not shut down yet", "error from ShutDown endpoint", err)

	// wait for another 10 seconds
	// TODO: Make this flexible for unit tests that don't want to wait here and real-world cases where we
	// have to wait longer for a device to be shut down
	return actionContinue{delay: 1 * time.Second}
}

// actionProvisionDevice provisions the device.
func (s *Service) actionProvisionDevice(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionProvisionDevice")
	log.V(1).Info("Started function")

	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err)
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// if device cannot be found, we associate a new one
			log.Info("Device to provision not found. Go back to StateAssociateDevice")
			record.Warnf(s.scope.HivelocityMachine, "DeviceNotFound", "Hivelocity device not found. Associate new one")
			return actionGoBack{nextState: infrav1.StateAssociateDevice}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	userData, err := s.scope.GetRawBootstrapData(ctx)
	if err != nil {
		record.Warnf(s.scope.HivelocityMachine, "FailedGetBootstrapData", err.Error())
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
		// find ssh key in Hivelocity API based on the name specified in the HVCluster spec
		sshKeyName := s.scope.HivelocityCluster.Spec.SSHKey.Name
		keys, err := s.scope.HVClient.ListSSHKeys(ctx)
		if err != nil {
			if errors.Is(err, hvclient.ErrRateLimitExceeded) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function ListSSHKeys")
			}
			return actionError{err: fmt.Errorf("failed to list ssh keys: %w", err)}
		}
		sshKeyID, err := findSSHKey(keys, sshKeyName)
		if err != nil {
			if errors.Is(err, errSSHKeyNotFound) {
				// do not return an error in the reconcile loop as we cannot do anything about this without the intervention
				// of the user. Only after the SSH key has been uploaded correctly, the provisioning can continue.
				// This is why we wait for 5m and then reconcile again to see whether the SSH key exists then.
				record.Warnf(s.scope.HivelocityCluster, "SSHKeyNotFound", "ssh key %s could not be found", sshKeyName)
				return actionFailed{}
			}
			return actionError{err: fmt.Errorf("error with ssh keys: %w", err)}
		}
		opts.PublicSshKeyId = sshKeyID
	}

	// Provision the device
	if _, err := s.scope.HVClient.ProvisionDevice(ctx, deviceID, opts); err != nil {
		// TODO: Handle error that machine is not shut down
		if errors.Is(err, hvclient.ErrRateLimitExceeded) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function ProvisionDevice")
		}
		record.Warnf(s.scope.HivelocityMachine, "FailedProvisionDevice", "Failed to provision device %v: %s", deviceID, err)
		return actionError{err: fmt.Errorf("failed to provision device %d: %s", deviceID, err)}
	}

	log.V(1).Info("Completed function")
	return actionComplete{}
}

func findSSHKey(sshKeysInAPI []hv.SshKeyResponse, sshKeyName string) (int32, error) {
	for _, key := range sshKeysInAPI {
		if key.Name == sshKeyName {
			return key.SshKeyId, nil
		}
	}
	return 0, errSSHKeyNotFound
}

// TODO: Implement logic to add multiple images.
func (s *Service) getDeviceImage(ctx context.Context) (string, error) {
	return defaultImageName, nil
}

// actionDeviceProvisioned reconciles a provisioned device.
func (s *Service) actionDeviceProvisioned(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionDeviceProvisioned")
	log.V(1).Info("Started function")

	// get device
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("[actionDeviceProvisioned] ProviderIDToDeviceID failed: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err)
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// fatal error when device was not found
			conditions.MarkFalse(
				s.scope.HivelocityMachine,
				infrav1.DeviceReadyCondition,
				infrav1.DeviceNotFoundReason,
				clusterv1.ConditionSeverityError,
				fmt.Sprintf("device %d not found", device.DeviceId),
			)
			record.Warnf(s.scope.HivelocityMachine, "DeviceNotFound", "Hivelocity device not found")
			s.scope.HivelocityMachine.SetFailure(capierrors.UpdateMachineError, infrav1.FailureMessageDeviceNotFound)
			return actionComplete{}
		}
		return actionError{err: fmt.Errorf("failed to get associated device: %w", err)}
	}

	// verify device
	if err := s.verifyAssociatedDevice(&device); err != nil {
		// fatal error when device could not be verified
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceReadyCondition,
			infrav1.DeviceTagsInvalidReason,
			clusterv1.ConditionSeverityError,
			fmt.Sprintf("device %d has invalid tags", device.DeviceId),
		)
		record.Warnf(s.scope.HivelocityMachine, "DeviceTagsInvalid", "Hivelocity device not found.")
		s.scope.HivelocityMachine.SetFailure(capierrors.UpdateMachineError, infrav1.FailureMessageDeviceTagsInvalid)
		return actionComplete{}
	}

	// update machine object with infos from device
	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition)
	s.scope.HivelocityMachine.SetMachineStatus(device)
	s.scope.HivelocityMachine.Status.Ready = true

	log.V(1).Info("Completed function")
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

// actionDeleteDeviceDeProvision re-provisions a device to remove it from cluster.
func (s *Service) actionDeleteDeviceDeProvision(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionDeleteDeviceDeProvision")
	log.V(1).Info("Started function")

	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err)
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// Nothing to do if device is not found
			s.scope.Info("Unable to locate Hivelocity device by ID or tags")
			record.Warnf(s.scope.HivelocityMachine, "NoDeviceFound", "Unable to find matching Hivelocity device for %s", s.scope.Name())
			return actionComplete{}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	newTags, _ := s.scope.HivelocityCluster.DeviceTag().RemoveFromList(device.Tags)
	newTags, _ = s.scope.HivelocityMachine.DeviceTag().RemoveFromList(newTags)
	newTags, _ = s.scope.DeviceTagMachineType().RemoveFromList(newTags)

	opts := hv.BareMetalDeviceUpdate{
		Hostname:    s.scope.Name(),
		Tags:        newTags,
		OsName:      defaultImageName,
		ForceReload: true,
	}

	// Provision the device
	if _, err := s.scope.HVClient.ProvisionDevice(ctx, deviceID, opts); err != nil {
		// TODO: Handle error that machine is not shut down
		if errors.Is(err, hvclient.ErrRateLimitExceeded) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function ProvisionDevice")
		}
		record.Warnf(s.scope.HivelocityMachine, "FailedDeProvisionDevice", "Failed to de-provision device %s: %s", deviceID, err)
		return actionError{err: fmt.Errorf("failed to de-provision device %d: %s", deviceID, err)}
	}
	log.V(1).Info("Completed function")
	return actionComplete{}
}

// actionDeleteDeviceDissociate ensures that the device has no tags of machine.
func (s *Service) actionDeleteDeviceDissociate(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionDeleteDeviceDissociate")
	log.V(1).Info("Started function")

	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err)
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// Nothing to do if device is not found
			s.scope.Info("Unable to locate Hivelocity device by ID or tags")
			record.Warnf(s.scope.HivelocityMachine, "NoDeviceFound", "Unable to find matching Hivelocity device for %s", s.scope.Name())
			return actionComplete{}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	newTags, updated1 := s.scope.HivelocityCluster.DeviceTag().RemoveFromList(device.Tags)
	newTags, updated2 := s.scope.HivelocityMachine.DeviceTag().RemoveFromList(newTags)
	newTags, updated3 := s.scope.DeviceTagMachineType().RemoveFromList(newTags)

	if updated1 || updated2 || updated3 {
		if err := s.scope.HVClient.SetDeviceTags(ctx, device.DeviceId, newTags); err != nil {
			if errors.Is(err, hvclient.ErrRateLimitExceeded) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function SetDeviceTags")
			}
			return actionError{err: fmt.Errorf("failed to set tags: %w", err)}
		}
	}

	log.V(1).Info("Completed function")
	return actionComplete{}
}

func (s *Service) handleRateLimitExceeded(err error) {
	if errors.Is(err, hvclient.ErrRateLimitExceeded) {
		conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
		record.Event(s.scope.HivelocityMachine, "RateLimitExceeded", "exceeded rate limit with calling Hivelocity function ListSSHKeys")
	}
}
