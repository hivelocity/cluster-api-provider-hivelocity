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
	"fmt"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	hvutils "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/hvutils"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const maxShutDownTime = 2 * time.Minute
const deviceOffTimeout = 10 * time.Minute

// Service defines struct with machine scope to reconcile Hivelocity machines.
type Service struct {
	scope *scope.MachineScope
}

// NewService outs a new service with machine scope.
func NewService(scope *scope.MachineScope) *Service {
	return &Service{
		scope: scope,
	}
}

// Reconcile implements reconcilement of Hivelocity machines.
func (s *Service) Reconcile(ctx context.Context) (_ *ctrl.Result, err error) {
	if s.scope.HivelocityCluster.Spec.HivelocitySecret.Key == "" {
		record.Eventf(s.scope.HivelocityMachine, corev1.EventTypeWarning, "NoAPIKey", "No Hivelocity API Key found")
		return nil, fmt.Errorf("no Hivelocity API Key provided - cannot reconcile Hivelocity device")
	}

	// detect failure domain. question: two names "failure domain" and "region". One name would be better.
	failureDomain, err := s.scope.GetFailureDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to get failure domain: %w", err)
	}
	s.scope.HivelocityMachine.Status.Region = infrav1.Region(failureDomain)

	// Waiting for bootstrap data to be ready
	if !s.scope.IsBootstrapDataReady(ctx) {
		s.scope.Info("Bootstrap not ready - requeuing")
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.InstanceBootstrapReadyCondition,
			infrav1.InstanceBootstrapNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"bootstrap not ready yet",
		)
		return &ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.InstanceBootstrapReadyCondition,
	)

	// Try to find the associate device.
	device, err := s.findAssociateDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	// If no device is found we have to create one
	if device == nil {
		device, err = s.createDevice(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create device: %w", err)
		}
		record.Eventf(
			s.scope.HivelocityMachine,
			"SuccessfulCreate",
			"Created new device with id %d",
			device.DeviceId,
		)
	}

	c := s.scope.HivelocityMachine.Status.Conditions.DeepCopy()
	s.scope.HivelocityMachine.Status = setStatusFromAPI(device)
	s.scope.HivelocityMachine.Status.Conditions = c

	switch device.PowerStatus {
	case hvclient.PowerStatusOff:
		return s.handleDeviceStatusOff(ctx, device)
	/* todo: up to now we only get ON/OFF
	case hv.BareMetalDeviceStatusStarting:
		// Requeue here so that device does not switch back and forth between off and starting.
		// If we don't return here, the condition InstanceReady would get marked as true in this
		// case. However, if the device is stuck and does not power on, we should not mark the
		// condition InstanceReady as true to be able to remediate the device after a timeout.
		return &reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	*/
	case hvclient.PowerStatusOn: // Do nothing
	default:
		s.scope.HivelocityMachine.Status.Ready = false
		s.scope.V(1).Info("device not in running state",
			"device", device.Hostname,
			"powerStatus", device.PowerStatus)
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	providerID := fmt.Sprintf("hivelocity://%d", device.DeviceId)

	if !s.scope.IsControlPlane() {
		s.scope.HivelocityMachine.Spec.ProviderID = &providerID
		s.scope.HivelocityMachine.Status.Ready = true
		conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition)
		return nil, nil
	}

	s.scope.HivelocityMachine.Spec.ProviderID = &providerID
	s.scope.HivelocityMachine.Status.Ready = true
	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition)

	return nil, nil
}

func (s *Service) handleDeviceStatusOff(ctx context.Context, device *hv.BareMetalDevice) (*reconcile.Result, error) {
	// Check if device is in DeviceStatusOff and turn it on. This is to avoid a bug of Hivelocity where
	// sometimes machines are created and not turned on

	condition := conditions.Get(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition)
	if condition != nil &&
		condition.Status == corev1.ConditionFalse &&
		condition.Reason == infrav1.DeviceOffReason {
		if time.Now().Before(condition.LastTransitionTime.Time.Add(deviceOffTimeout)) {
			// Not yet timed out, try again to power on
			if err := s.scope.HVClient.PowerOnDevice(ctx, device.DeviceId); err != nil {
				return nil, fmt.Errorf("failed to power on device: %w", err)
			}
		} else {
			// Timed out. Set failure reason
			s.scope.SetError("reached timeout of waiting for machines that are switched off", capierrors.CreateMachineError)
			return nil, nil
		}
	} else { // todo: too complicated. Is there a way to avoid the "else"?
		// No condition set yet. Try to power device on.
		if err := s.scope.HVClient.PowerOnDevice(ctx, device.DeviceId); err != nil {
			if hvclient.IsRateLimitExceededError(err) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine,
					"RateLimitExceeded",
					"exceeded rate limit with calling hivelocity function PowerOnDevice",
				)
			}
			return nil, fmt.Errorf("failed to power on device: %w", err)
		}
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.InstanceReadyCondition,
			infrav1.DeviceOffReason,
			clusterv1.ConditionSeverityInfo,
			"device is switched off",
		)
	}

	// Try again in 30 sec.
	return &reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

func (s *Service) createDevice(ctx context.Context) (*hv.BareMetalDevice, error) {
	// get userData
	userData, err := s.scope.GetRawBootstrapData(ctx)
	if err != nil {
		record.Warnf(
			s.scope.HivelocityMachine,
			"FailedGetBootstrapData",
			err.Error(),
		)
		return nil, fmt.Errorf("failed to get raw bootstrap data: %s", err)
	}

	image, err := s.getDeviceImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get device image: %w", err)
	}

	devices, err := s.scope.HVClient.ListDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("[Service.createDevice] ListDevices() failed. cluster %q: %w",
			s.scope.Name(), err)
	}
	unusedDevice, err := hvutils.FindUnusedDevice(devices, s.scope.Name(), "question instance-type")
	if err != nil {
		return nil, fmt.Errorf("[Service.createDevice] FindUnusedDevice() failed. cluster %q: %w",
			s.scope.Name(), err)
	}
	if unusedDevice == nil {
		return nil, fmt.Errorf("[Service.createDevice] FindUnusedDevice() found no unused device. cluster %q: %w",
			s.scope.Name(), err)
	}
	opts := hv.BareMetalDeviceUpdate{
		Hostname: s.scope.Name(),
		Tags:     createTags(s.scope.ClusterScope.Name(), s.scope.Name(), s.scope.IsControlPlane()),
		Script:   string(userData), // cloud-init script
		OsName:   image,
	}

	if s.scope.HivelocityCluster.Spec.SSHKey != nil {
		sshKeyID, err := s.getSSHKeyIDFromSSHKeyName(ctx, s.scope.HivelocityCluster.Spec.SSHKey)
		if err != nil {
			return nil, fmt.Errorf("error with ssh keys: %w", err)
		}

		opts.PublicSshKeyId = sshKeyID
	}

	// Create the device
	device, err := s.scope.HVClient.CreateDevice(ctx, unusedDevice.DeviceId, opts)
	if err != nil {
		if hvclient.IsRateLimitExceededError(err) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine,
				"RateLimitExceeded",
				"exceeded rate limit with calling Hivelocity function CreateDevice",
			)
		}
		record.Warnf(s.scope.HivelocityMachine,
			"FailedCreateHivelocityDevice",
			"Failed to create Hivelocity device %s: %s",
			s.scope.Name(),
			err,
		)
		return nil, fmt.Errorf("error while creating Hivelocity device %s: %s", s.scope.HivelocityMachine.Name, err)
	}

	return &device, nil
}

const defaultImageName = "Ubuntu 20.x"

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
	device, err := s.findAssociateDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	// If no device has been found then nothing can be deleted
	if device == nil {
		s.scope.V(2).Info("Unable to locate Hivelocity device by ID or tags")
		record.Warnf(s.scope.HivelocityMachine, "NoInstanceFound", "Unable to find matching Hivelocity device for %s", s.scope.Name())
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
	if conditions.IsTrue(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition) ||
		conditions.IsFalse(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition) &&
			conditions.GetReason(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition) == infrav1.InstanceTerminatedReason &&
			time.Now().Before(conditions.GetLastTransitionTime(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition).Time.Add(maxShutDownTime)) {
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
			infrav1.InstanceReadyCondition,
			infrav1.InstanceTerminatedReason,
			clusterv1.ConditionSeverityInfo,
			"Instance has been shut down")
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

func setStatusFromAPI(device *hv.BareMetalDevice) infrav1.HivelocityMachineStatus {
	var status infrav1.HivelocityMachineStatus
	// todo: HV does not have a detailed status for their devices. Only ON or OFF.
	return status
}

// We write the machine name in the labels, so that all labels are or should be unique.
func (s *Service) findAssociateDevice(ctx context.Context) (*hv.BareMetalDevice, error) {
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

func createTags(clusterName string, machineName string, isControlPlane bool) []string {
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
