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
	"strings"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	hvutils "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/hvutils"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
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
	// errTooManyTagsFound gets returned, if there are multiple tags with the same key,
	// and the key should be unique.
	errTooManyTagsFound = fmt.Errorf("too many tags found")

	// errNoMatchingTagFound gets returned, if no matching tag was found.
	errNoMatchingTagFound = fmt.Errorf("no matching tag found")

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
			infrav1.DeviceBootstrapReadyCondition,
			infrav1.DeviceBootstrapNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"bootstrap not ready yet",
		)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.DeviceBootstrapReadyCondition,
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

// Delete implements delete method of the HV device.
func (s *Service) Delete(ctx context.Context) (_ *ctrl.Result, err error) {
	// find current device
	device, err := s.getAssociatedDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	// If no device has been found then nothing can be deleted
	if device == nil {
		s.scope.V(2).Info("Unable to locate Hivelocity device by ID or tags")
		record.Warnf(s.scope.HivelocityMachine, "NoDeviceFound", "Unable to find matching Hivelocity device for %s", s.scope.Name())
		return nil, nil
	}

	// First shut the device down, then delete it
	switch status := device.PowerStatus; status {
	case hvclient.PowerStatusOn:
		return s.handleDeleteDeviceStatusRunning(ctx, device)
	case hvclient.PowerStatusOff:
		return s.handleDeleteDeviceStatusOff(ctx, device)
	default:
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}
}

func (s *Service) handleDeleteDeviceStatusRunning(ctx context.Context, device *hv.BareMetalDevice) (*ctrl.Result, error) {
	// Check if the device has been tried to shut down already and if so,
	// if time of last condition change + maxWaitTime is already in the past.
	// With one of these two conditions true, delete device immediately. Otherwise, shut it down and requeue.
	if conditions.IsTrue(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition) ||
		conditions.IsFalse(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition) &&
			conditions.GetReason(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition) == infrav1.DeviceTerminatedReason &&
			time.Now().Before(conditions.GetLastTransitionTime(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition).Time.Add(maxShutDownTime)) {
		if err := s.scope.HVClient.ShutdownDevice(ctx, device.DeviceId); err != nil {
			if hvclient.IsRateLimitExceededError(err) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine,
					"RateLimitExceeded",
					"exceeded rate limit with calling Hivelocity function ShutdownDevice",
				)
			}
			return &reconcile.Result{}, fmt.Errorf("failed to shutdown device: %w", err)
		}
		conditions.MarkFalse(s.scope.HivelocityMachine,
			infrav1.DeviceReadyCondition,
			infrav1.DeviceTerminatedReason,
			clusterv1.ConditionSeverityInfo,
			"Device has been shut down")
		return &ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}
	if err := s.scope.HVClient.DeleteDevice(ctx, device.DeviceId); err != nil {
		if hvclient.IsRateLimitExceededError(err) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine,
				"RateLimitExceeded",
				"exceeded rate limit with calling Hivelocity function DeleteDevice",
			)
		}
		record.Warnf(s.scope.HivelocityMachine, "FailedDeleteHivelocityDevice", "Failed to delete Hivelocity device %s", s.scope.Name())
		return &reconcile.Result{}, fmt.Errorf("failed to delete device: %w", err)
	}

	record.Eventf(
		s.scope.HivelocityMachine,
		"HivelocityDeviceDeleted",
		"Hivelocity device %s deleted",
		s.scope.Name(),
	)
	return nil, nil
}

func (s *Service) handleDeleteDeviceStatusOff(ctx context.Context, device *hv.BareMetalDevice) (*ctrl.Result, error) {
	return nil, fmt.Errorf("todo: handleDeleteDeviceStatusOff()")
}

// We write the machine name in the labels, so that all labels are or should be unique.
func (s *Service) getAssociatedDevice(ctx context.Context) (*hv.BareMetalDevice, error) {
	clusterTag := hvclient.GetClusterTag(s.scope.ClusterScope.Name())
	machineTag := hvclient.GetMachineTag(s.scope.Name())
	devices, err := s.scope.HVClient.ListDevices(ctx)
	if err != nil {
		if hvclient.IsRateLimitExceededError(err) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine,
				"RateLimitExceeded",
				"exceeded rate limit with calling ListDevices",
			)
		}
		return nil, err
	}
	return hvutils.FindDeviceByTags(clusterTag, machineTag, devices)
}

func createTags(clusterName, machineName string, isControlPlane bool) []string {
	var machineType string
	if isControlPlane {
		machineType = "control_plane"
	} else {
		machineType = "worker"
	}
	return []string{
		hvclient.GetClusterTag(clusterName),
		hvclient.GetMachineTag(machineName),
		"caphv-machinetype=" + machineType,
	}
}

// setMachineAddress gets the address from the device and sets it on the Machine object.
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

// associateDevice claims an unused HV device by settings tags and returns it.
func (s *Service) associateDevice(ctx context.Context) (*hv.BareMetalDevice, error) {
	device, err := s.chooseDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to choose device: %w", err)
	}

	device.Tags = append(device.Tags, clusterAndMachineTag(s.scope.HivelocityCluster.Name, s.scope.Name())...)
	if err := s.scope.HVClient.SetTags(ctx, device.DeviceId, device.Tags); err != nil {
		return nil, fmt.Errorf("failed to set tags on machine %s :%w", s.scope.Name(), err)
	}
	return device, nil
}

func clusterAndMachineTag(clusterName, machineName string) []string {
	return []string{
		hvclient.TagKeyClusterName + "=" + clusterName,
		hvclient.TagKeyMachineName + "=" + machineName,
	}
}

// chooseDevice searches for an unused device.
func (s *Service) chooseDevice(ctx context.Context) (*hv.BareMetalDevice, error) {
	devices, err := s.scope.HVClient.ListDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("[chooseDevice] ListDevices() failed. machine %q: %w",
			s.scope.Name(), err)
	}
	return chooseAvailableFromList(ctx, devices, s.scope.HivelocityMachine.Spec.Type, s.scope.Cluster.Name)
}

// deviceExists returns true if the device exists.
func (s *Service) deviceExists(ctx context.Context, deviceID int32) (bool, error) {
	// question: should we check if the device is in the current cluster first?
	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if errors.Is(err, hvclient.ErrDeviceNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if device.PrimaryIp != "" {
		return true, nil
	}
	return false, nil
}

func chooseAvailableFromList(ctx context.Context, devices []*hv.BareMetalDevice, deviceType infrav1.HivelocityDeviceType, clusterName string) (*hv.BareMetalDevice, error) {
	log := ctrl.LoggerFrom(ctx)
	for _, device := range devices {
		// Ignore if associated already
		if isAssociated(device) {
			continue
		}

		dt, err := getDeviceType(device)
		if err != nil {
			if errors.Is(err, errNoMatchingTagFound) {
				continue
			}
			// don't return on err, otherwise a single broken device would block this method
			log.Error(err, "[chooseAvailableFromList] getDeviceType() failed")
			continue
		}

		// Ignore if has wrong device type
		if dt != string(deviceType) {
			continue
		}

		// Ignore if associated to other cluster
		cn, found := findAssociatedCluster(device)
		if found && clusterName != cn {
			continue
		}
		return device, nil
	}
	return nil, errNoDeviceAvailable
}

// getDeviceType returns the device-type of this BareMetalDevice.
func getDeviceType(device *hv.BareMetalDevice) (string, error) {
	deviceType, err := findValueForKeyInTags(device.Tags, hvclient.TagKeyDeviceType)
	if err != nil {
		return "", fmt.Errorf("[getDeviceType] findValueForKeyInTags() failed: %w", err)
	}
	return deviceType, nil
}

// isAssociated returns true if the device is associated with a HivelocityMachine.
func isAssociated(device *hv.BareMetalDevice) bool {
	deviceType, _ := findValueForKeyInTags(device.Tags, hvclient.TagKeyMachineName)
	return deviceType != ""
}

// getAssociatedMachines returns the associated HivelocityMachine.
func getAssociatedMachines(device *hv.BareMetalDevice) ([]string, error) {
	machines, err := findAllValuesForKeyInTags(device.Tags, hvclient.TagKeyMachineName)
	if err != nil {
		return nil, fmt.Errorf("[getAssociatedMachine] findAllValuesForKeyInTags() failed: %w", err)
	}

	return machines, nil
}

// findAssociatedCluster tries to find a cluster name in the tags. If it finds one, it returns the name and true
// If not, it returns an empty string and false.
func findAssociatedCluster(device *hv.BareMetalDevice) (string, bool) {
	cluster, _ := findValueForKeyInTags(device.Tags, hvclient.TagKeyClusterName)
	return cluster, cluster != ""
}

// findValueForKeyInTags returns the value of a TagKey of a device.
// Example: If a device has the tag "foo=bar", then getValueOfTag
// will return "bar".
// If there is no such tag, or there are two tags, then an error gets returned.
func findValueForKeyInTags(tagList []string, key string) (string, error) {
	prefix := key + "="
	found := 0
	var value string

	for _, tag := range tagList {
		if !strings.HasPrefix(tag, prefix) {
			continue
		}
		if found > 1 {
			return "", errTooManyTagsFound
		}
		found++
		value = tag[len(prefix):]
	}
	if found == 0 {
		return "", errNoMatchingTagFound
	}
	return value, nil
}

// findAllValuesForKeyInTags returns all values of a TagKey of a device.
// Example: If a device has the tag "foo=bar", then getValueOfTag
// will return "bar".
// If there is no such tag, then an error gets returned.
func findAllValuesForKeyInTags(tagList []string, key string) ([]string, error) {
	prefix := key + "="

	// expect one value only
	values := make([]string, 0, 1)

	var found int
	for _, tag := range tagList {
		if !strings.HasPrefix(tag, prefix) {
			continue
		}
		found++
		values = append(values, tag[len(prefix):])
	}
	if found == 0 {
		return nil, errNoMatchingTagFound
	}
	return values, nil
}

// findAllValuesForKeyInTags returns all values of a TagKey of a device.
// Example: If a device has the tag "foo=bar", then getValueOfTag
// will return "bar".
// If there is no such tag, then an error gets returned.
func removeAllValuesForKeyFromTags(tagList []string, key string) (newTagList []string, updatedTags bool) {
	prefix := key + "="
	newTagList = make([]string, 0, len(tagList))
	for _, tag := range tagList {
		if strings.HasPrefix(tag, prefix) {
			updatedTags = true
		} else {
			newTagList = append(newTagList, tag)
		}
	}
	return newTagList, updatedTags
}

// actionAssociateDevice claims an unused HV device by settings tags and returns it.
func (s *Service) actionAssociateDevice(ctx context.Context) actionResult {
	device, err := s.chooseDevice(ctx)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to choose device: %w", err)}
	}

	device.Tags = append(device.Tags, clusterAndMachineTag(s.scope.HivelocityCluster.Name, s.scope.Name())...)
	if err := s.scope.HVClient.SetTags(ctx, device.DeviceId, device.Tags); err != nil {
		return actionError{err: fmt.Errorf("failed to set tags: %w", err)}
	}
	providerID := hvutils.DeviceIDToProviderID(device.DeviceId)
	s.scope.HivelocityMachine.Spec.ProviderID = &providerID
	return actionComplete{}
}

// actionVerifyAssociate verifies that the HV device has actually been associated to this machine and only this.
// Checking whether there are other machines also associated avoids situations where machines are selected at the same time.
func (s *Service) actionVerifyAssociate(ctx context.Context) actionResult {
	// TODO: We should be able to control this value from outside. Do we do it with flags?
	// wait for 3 seconds at least before checking again
	const waitFor = 100 * time.Millisecond

	// if the waiting time has not yet passed, then we reconcile again without changing state
	if !hasTimedOut(s.scope.HivelocityMachine.Spec.Status.LastUpdated, waitFor) {
		return actionContinue{delay: 100 * time.Millisecond}
	}

	// if waiting time is over, we check the server for tags
	deviceID, err := hvutils.ProviderIDToDeviceID(*s.scope.HivelocityMachine.Spec.ProviderID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	machineNames, err := getAssociatedMachines(&device)
	if err != nil {
		// if no associated machine is found, we need to associate new machine
		if errors.Is(err, errNoMatchingTagFound) {
			return actionError{err: errGoToPreviousState}
		}
		return actionError{err: fmt.Errorf("failed to get associated machine: %w", err)}
	}

	// if only one machine is associated and it is the correct one, then complete.
	if len(machineNames) == 1 && machineNames[0] == s.scope.Name() {
		return actionComplete{}
	}

	// Something unexpected happened. Remove label of current machine from device.
	newTagList, updatedTags := removeAllValuesForKeyFromTags(device.Tags, hvclient.TagKeyMachineName)
	if updatedTags {
		if err := s.scope.HVClient.SetTags(ctx, deviceID, newTagList); err != nil {
			return actionError{err: fmt.Errorf("failed to remove associated machine from labels: %w", err)}
		}
	}

	return actionError{err: errGoToPreviousState}
}

func hasTimedOut(lastUpdated *metav1.Time, timeout time.Duration) bool {
	if lastUpdated == nil {
		return false
	}
	now := metav1.Now()
	return lastUpdated.Add(timeout).Before(now.Time)
}

// actionShutDownDevice shuts down the device.
func (s *Service) actionShutDownDevice(ctx context.Context) actionResult {
	deviceID, err := hvutils.ProviderIDToDeviceID(*s.scope.HivelocityMachine.Spec.ProviderID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	if err := s.scope.HVClient.ShutdownDevice(ctx, deviceID); err != nil {
		return actionError{err: fmt.Errorf("failed to shut device down: %w", err)}
	}
	return actionComplete{}
}

// actionEnsureDeviceShutDown ensures that the device is shut down.
func (s *Service) actionEnsureDeviceShutDown(ctx context.Context) actionResult {
	const timeout = 10 * time.Minute

	deviceID, err := hvutils.ProviderIDToDeviceID(*s.scope.HivelocityMachine.Spec.ProviderID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	_, err = s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	// TODO: Ping server and check whether it can be reached
	pingSuccessful := true
	if pingSuccessful {
		return actionComplete{}
	}
	// Device is not shut down yet, check the timeout. If we wait for too long, trigger the shutdown again.
	if hasTimedOut(s.scope.HivelocityMachine.Spec.Status.LastUpdated, timeout) {
		return actionError{err: errGoToPreviousState}
	}

	// wait for another minute
	return actionContinue{delay: time.Minute}
}

// actionProvisionDevice provisions the device.
func (s *Service) actionProvisionDevice(ctx context.Context) actionResult {
	deviceID, err := hvutils.ProviderIDToDeviceID(*s.scope.HivelocityMachine.Spec.ProviderID)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
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
	opts := hv.BareMetalDeviceUpdate{
		Hostname:    s.scope.Name(),
		Tags:        createTags(s.scope.ClusterScope.Name(), s.scope.Name(), s.scope.IsControlPlane()),
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
	s.scope.Info("[updateDevice] ProvisionDevice was called", "err", err, "deviceID", deviceID)
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
	deviceID, err := hvutils.ProviderIDToDeviceID(*s.scope.HivelocityMachine.Spec.ProviderID)
	if err != nil {
		return actionError{err: fmt.Errorf("[updateDevice] ProviderIDToDeviceID failed: %w", err)}
	}

	// FIXME: we already get the device
	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		// TODO: Filter out device not found error. Replace the below statement.
		if errors.Is(err, errNoDeviceAvailable) {
			conditions.MarkFalse(s.scope.HivelocityMachine,
				infrav1.DeviceReadyCondition,
				infrav1.DeviceDeletedReason,
				clusterv1.ConditionSeverityError,
				fmt.Sprintf("device %d does not exists anymore", device.DeviceId))
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition)
	setMachineAddress(s.scope.HivelocityMachine, &device)
	s.scope.HivelocityMachine.Status.Ready = true

	return actionComplete{}
}
