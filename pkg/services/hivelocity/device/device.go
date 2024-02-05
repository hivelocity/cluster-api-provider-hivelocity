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
	_ "embed"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	hvlabels "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/labels"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/utils"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
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
	errSSHKeyNotFound = fmt.Errorf("ssh key not found")

	errWrongMachineTag = fmt.Errorf("machine has wrong machine tag")

	errWrongClusterTag = fmt.Errorf("machine has wrong cluster tag")
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
		return reconcile.Result{}, fmt.Errorf("action %q failed: %w", initialState, err)
	}

	// TODO: Verify that patching the object is fine. Alternative would be to update it if that is better since we update the Spec as well
	return result, nil
}

// actionAssociateDevice claims an unused HV device by settings tags and returns it.
func (s *Service) actionAssociateDevice(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionAssociateDevice")
	log.V(1).Info("Started function")

	device, reason, err := GetFirstFreeDevice(ctx, s.scope.HVClient, s.scope.HivelocityMachine.Spec, s.scope.HivelocityCluster)
	if err != nil {
		s.handleRateLimitExceeded(err, "ListDevices")
		return actionError{err: fmt.Errorf("failed to find available device: %w (%s)", err, reason)}
	}

	if device == nil {
		msg := fmt.Sprintf("no available device found (selector: %+v) (%s)", s.scope.HivelocityMachine.Spec.DeviceSelector, reason)
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceAssociateSucceededCondition,
			infrav1.NoAvailableDeviceReason,
			clusterv1.ConditionSeverityWarning,
			msg,
		)
		record.Warn(s.scope.HivelocityMachine, "NoFreeDeviceFound", msg)
		return actionContinue{delay: 30 * time.Second}
	}
	conditions.Delete(s.scope.HivelocityMachine, infrav1.DeviceAssociateSucceededCondition)

	// associate this device with the machine object by setting tags
	device.Tags = append(device.Tags,
		s.scope.HivelocityCluster.DeviceTagOwned().ToString(),
		s.scope.HivelocityCluster.DeviceTag().ToString(),
		s.scope.HivelocityMachine.DeviceTag().ToString(),
		s.scope.DeviceTagMachineType().ToString(),
	)

	if err := s.scope.HVClient.SetDeviceTags(ctx, device.DeviceId, device.Tags); err != nil {
		s.handleRateLimitExceeded(err, "SetDeviceTags")
		return actionError{err: fmt.Errorf("failed to set tags on device %v: %w", device.DeviceId, err)}
	}

	// set providerID on machine object which is based on deviceID
	s.scope.HivelocityMachine.SetProviderID(device.DeviceId)

	log.V(1).Info("Completed function")
	return actionComplete{}
}

// GetFirstFreeDevice finds the first free matching device. It returns nil if no device is found.
func GetFirstFreeDevice(ctx context.Context, hvclient hvclient.Client, hvMachineSpec infrav1.HivelocityMachineSpec, hvCluster *infrav1.HivelocityCluster) (*hv.BareMetalDevice, string, error) {
	// list all devices
	devices, err := hvclient.ListDevices(ctx)
	if err != nil {
		return nil, "", err
	}

	// Since we don't have a LoadBalancer we use the IP of the first ControlPlane
	// for hvCluster.Spec.ControlPlaneEndpoint. When the cluster gets created,
	// and there are several available devices, we need to get a stable result.
	// Later, if we have LoadBalancers, then we want to opposite behaviour.
	// Then, we want list of devices to get randomized, to reduce conflicts
	// during concurrent associating of devices.
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].DeviceId < devices[j].DeviceId
	})

	device, reason, err := findAvailableDeviceFromList(devices, hvMachineSpec.DeviceSelector, hvCluster.Name)
	if err != nil {
		return nil, "", err
	}
	log := ctrl.LoggerFrom(ctx)
	if device != nil {
		log.Info(fmt.Sprintf("GetFirstFreeDevice hvMachineSpec.DeviceSelector %+v device.Tags: %+v",
			hvMachineSpec.DeviceSelector,
			device.Tags))
	}
	return device, reason, nil
}

func findAvailableDeviceFromList(devices []hv.BareMetalDevice, deviceSelector infrav1.DeviceSelector, clusterName string) (
	*hv.BareMetalDevice, string, error,
) {
	labelSelector, err := getLabelSelector(deviceSelector)
	if err != nil {
		return nil, "", fmt.Errorf("getLabelSelector failed: %w", err)
	}
	mapOfSkipReasons := make(map[string]int)

	for _, device := range devices {
		// Skip if caphv-permanent-error exists
		_, err := hvtag.PermanentErrorTagFromList(device.Tags)
		if err == nil {
			// The tag exists. Skip this Device
			mapOfSkipReasons["permanent-error"]++
			continue
		}

		if !hvtag.DeviceUsableByCAPI(device.Tags) {
			// not allowed to use the device
			mapOfSkipReasons["caphv-use-allow-is-missing"]++
			continue
		}

		// Ignore if associated to other cluster
		clusterTag, err := hvtag.ClusterTagFromList(device.Tags)
		if err != nil && !errors.Is(err, hvtag.ErrDeviceTagNotFound) {
			// unexpected error - continue
			mapOfSkipReasons["unexpected-error-while-get-cluster-tag"]++
			continue
		}
		if clusterTag.Value != "" && clusterTag.Value != clusterName {
			// associated to another cluster - continue
			mapOfSkipReasons["associated-to-other-cluster"]++
			continue
		}

		// Ignore if associated already
		machineTag, err := hvtag.MachineTagFromList(device.Tags)
		if err != nil && !errors.Is(err, hvtag.ErrDeviceTagNotFound) {
			// unexpected error - continue
			mapOfSkipReasons["unexpected-error-while-get-machine-tag"]++
			continue
		}
		if machineTag.Value != "" {
			// associated to other machine - continue
			mapOfSkipReasons["machine-already-associated"]++
			continue
		}

		if !labelSelector.Matches(hvlabels.Tags(device.Tags)) {
			mapOfSkipReasons["label-selector-does-not-match"]++
			continue
		}

		return &device, "", nil
	}
	reasons := make([]string, 0, len(mapOfSkipReasons))
	keys := maps.Keys(mapOfSkipReasons)
	slices.Sort(keys)
	for _, key := range keys {
		value := mapOfSkipReasons[key]
		if value == 0 {
			continue
		}
		reasons = append(reasons, fmt.Sprintf("%s: %d", key, value))
	}
	return nil, strings.Join(reasons, ", "), nil
}

func getLabelSelector(deviceSelector infrav1.DeviceSelector) (labels.Selector, error) {
	labelSelector := labels.NewSelector()
	var reqs labels.Requirements

	for labelKey, labelVal := range deviceSelector.MatchLabels {
		r, err := labels.NewRequirement(labelKey, selection.Equals, []string{labelVal})
		if err != nil {
			return labelSelector, err
		}
		reqs = append(reqs, *r)
	}
	for _, req := range deviceSelector.MatchExpressions {
		lowercaseOperator := selection.Operator(strings.ToLower(string(req.Operator)))
		r, err := labels.NewRequirement(req.Key, lowercaseOperator, req.Values)
		if err != nil {
			return labelSelector, err
		}
		reqs = append(reqs, *r)
	}

	return labelSelector.Add(reqs...), nil
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
		s.handleRateLimitExceeded(err, "GetDevice")
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
		record.Eventf(s.scope.HivelocityMachine, "SuccessfulAssociateDevice", "Device %d was associated with cluster %q", deviceID,
			s.scope.HivelocityCluster.Name)
		conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.DeviceAssociateSucceededCondition)
		return actionComplete{}
	}

	// Tags are not properly set or another machine also set its tags.
	// Remove cluster and machine tags and associate a new device.
	newTagList, updatedTags1 := s.scope.HivelocityCluster.DeviceTag().RemoveFromList(device.Tags)
	newTagList, updatedTags2 := s.scope.HivelocityMachine.DeviceTag().RemoveFromList(newTagList)
	if updatedTags1 || updatedTags2 {
		if err := s.scope.HVClient.SetDeviceTags(ctx, deviceID, newTagList); err != nil {
			s.handleRateLimitExceeded(err, "SetDeviceTags")
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

// actionVerifyShutdown makes sure that the device is shut down.
func (s *Service) actionVerifyShutdown(ctx context.Context) actionResult {
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("[actionVerifyShutdown] ProviderIDToDeviceID failed: %w", err)}
	}

	isReloading, isPoweredOn, err := s.getPowerAndReloadingState(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("[actionVerifyShutdown] getPowerAndReloadingState failed: %w", err)}
	}

	// if device is powered off and not reloading, then we are done and can start provisioning
	if !isPoweredOn && !isReloading {
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceProvisioningSucceededCondition,
			infrav1.DeviceShutDownReason,
			clusterv1.ConditionSeverityInfo,
			"device is shut down and will be provisioned",
		)
		return actionComplete{}
	}

	provisionCondition := conditions.Get(s.scope.HivelocityMachine, infrav1.DeviceProvisioningSucceededCondition)

	// handle reloading state

	if isReloading {
		if s.isReloadingTooLong(provisionCondition, isPoweredOn) {
			return s.setReloadingTooLongTag(ctx, deviceID, provisionCondition.LastTransitionTime)
		}

		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceProvisioningSucceededCondition,
			infrav1.DeviceReloadingReason,
			clusterv1.ConditionSeverityWarning,
			fmt.Sprintf("device %d is reloading", deviceID),
		)
		return actionContinue{delay: 1 * time.Minute}
	}

	// handle powered on state

	// if shutdown has been called in the past two minutes already, do not call it again and wait
	if provisionCondition != nil && provisionCondition.Reason == infrav1.DeviceShutdownCalledReason && !hasTimedOut(&provisionCondition.LastTransitionTime, 2*time.Minute) {
		return actionContinue{delay: 30 * time.Second}
	}

	// remove condition to reset the timer - we set the condition anyway again
	conditions.Delete(s.scope.HivelocityMachine, infrav1.DeviceProvisioningSucceededCondition)

	err = s.scope.HVClient.ShutdownDevice(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("[actionVerifyShutdown] ShutdownDevice failed: %w", err)}
	}

	record.Eventf(s.scope.HivelocityMachine, "SuccessfulShutdownDevice", "Called ShutdownDevice API for %d", deviceID)

	conditions.MarkFalse(
		s.scope.HivelocityMachine,
		infrav1.DeviceProvisioningSucceededCondition,
		infrav1.DeviceShutdownCalledReason,
		clusterv1.ConditionSeverityInfo,
		"device shut down has been triggered",
	)
	return actionContinue{delay: 30 * time.Second}
}

func (s *Service) isReloadingTooLong(condition *clusterv1.Condition, isPowerOn bool) bool {
	if condition == nil {
		return false
	}
	if condition.Reason != infrav1.DeviceReloadingReason {
		return false
	}
	timeout := 5 * time.Minute
	if isPowerOn {
		// the device is "reloading" during provisioning, which can take longer.
		timeout = 25 * time.Minute
	}
	if !hasTimedOut(&condition.LastTransitionTime, timeout) {
		return false
	}
	return true
}

// Set permanent error with an appropriate label,
// and go back to the state associate to associate with another device via actionGoBack.
// This method is used for provisioning and deprovisioning.
func (s *Service) setReloadingTooLongTag(ctx context.Context, deviceID int32, lastTransitionTime metav1.Time) actionResult {
	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err, "GetDevice")
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			msg := fmt.Sprintf("Hivelocity device %d not found", deviceID)
			conditions.MarkFalse(
				s.scope.HivelocityMachine,
				infrav1.DeviceReadyCondition,
				infrav1.DeviceNotFoundReason,
				clusterv1.ConditionSeverityError,
				msg,
			)
			record.Warnf(s.scope.HivelocityMachine, "DeviceNotFound", msg)
			s.scope.HivelocityMachine.SetFailure(capierrors.UpdateMachineError, infrav1.FailureMessageDeviceNotFound)
			return actionComplete{}
		}
		return actionError{err: fmt.Errorf("failed to get associated device: %w", err)}
	}
	tags := hvtag.RemoveEphemeralTags(device.Tags)
	_, err = hvtag.PermanentErrorTagFromList(tags)
	if errors.Is(err, hvtag.ErrDeviceTagNotFound) {
		tags = append(tags, fmt.Sprintf("%s=reloading-since-%s",
			hvtag.DeviceTagKeyPermanentError,
			lastTransitionTime.Format(time.RFC3339)))
	} else if err != nil {
		return actionError{err: fmt.Errorf("[setReloadingTooLongTag] PermanentErrorTagFromList failed: %w", err)}
	}

	err = s.scope.HVClient.SetDeviceTags(ctx, device.DeviceId, tags)
	if err != nil {
		return actionError{err: fmt.Errorf("[setReloadingTooLongTag] SetDeviceTags failed: %w", err)}
	}

	msg := fmt.Sprintf("device %d reloading too long. Tag %q was set. Trying next device.", device.DeviceId,
		hvtag.DeviceTagKeyPermanentError)

	conditions.MarkFalse(
		s.scope.HivelocityMachine,
		infrav1.DeviceReadyCondition,
		infrav1.DeviceReloadingTooLongReason,
		clusterv1.ConditionSeverityError,
		msg,
	)
	record.Warnf(s.scope.HivelocityMachine, "DeviceReloadingTooLong", msg)
	return actionGoBack{nextState: infrav1.StateAssociateDevice}
}

func (s *Service) getPowerAndReloadingState(ctx context.Context, deviceID int32) (
	isReloading bool, isPoweredOn bool, err error,
) {
	dump, err := s.scope.HVClient.GetDeviceDump(ctx, deviceID)
	if err != nil {
		return false, false, fmt.Errorf("[getPowerAndReloadingState] GetDeviceDump failed: %d %w", deviceID, err)
	}
	power, ok := dump.PowerStatus.(string)
	if !ok {
		return false, false, fmt.Errorf("[getPowerAndReloadingState] dump.PowerStatus failed: %d %+v %w", deviceID, dump.PowerStatus, err)
	}
	switch power {
	case "ON":
		isPoweredOn = true
	case "OFF":
		isPoweredOn = false
	default:
		return false, false, fmt.Errorf("[getPowerAndReloadingState] dump.PowerStatus unknown: %d %s %w", deviceID, power, err)
	}
	return dump.IsReload, isPoweredOn, nil
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
		s.handleRateLimitExceeded(err, "GetDevice")
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// if device cannot be found, we associate a new one
			log.Info("Device to provision not found. Go back to StateAssociateDevice")
			record.Warnf(s.scope.HivelocityMachine, "DeviceNotFound", "Hivelocity device not found. Associate new one")
			return actionGoBack{nextState: infrav1.StateAssociateDevice}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	isReloading, isPoweredOn, err := s.getPowerAndReloadingState(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("[actionProvisionDevice] getPowerAndReloadingState failed: %s", err)}
	}
	if isReloading {
		msg := fmt.Sprintf("unexpected: device %d is reloading.", deviceID)
		record.Warnf(s.scope.HivelocityMachine, "ProvisionDeviceIsReloading", msg)
		return actionError{err: errors.New(msg)}
	}
	if isPoweredOn {
		msg := fmt.Sprintf("device %d has power on. Waiting until it is shut-down.", deviceID) // xbug: here every 5 min.
		record.Eventf(s.scope.HivelocityMachine, "ProvisionDeviceIsPoweredOn", msg)
		return actionContinue{delay: 4 * time.Second}
	}

	userData, err := s.scope.GetRawBootstrapData(ctx)
	if err != nil {
		// This is common while starting a new cluster.
		record.Eventf(s.scope.HivelocityMachine, "FailedGetBootstrapData", "device %d: %s", deviceID, err.Error())
		return actionContinue{delay: 10 * time.Second}
	}

	image, err := s.getDeviceImage(ctx)
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get device image: %w", err)}
	}

	opts := hv.BareMetalDeviceUpdate{
		Hostname:    fmt.Sprintf("%s.example.com", s.scope.Name()), // TODO: HV API requires a FQDN.
		Tags:        device.Tags,
		Script:      "#cloud-config\n" + string(userData), // cloud-init script
		OsName:      image,
		ForceReload: true,
	}

	if s.scope.HivelocityCluster.Spec.SSHKey != nil {
		// find ssh key in Hivelocity API based on the name specified in the HVCluster spec
		sshKeyName := s.scope.HivelocityCluster.Spec.SSHKey.Name
		keys, err := s.scope.HVClient.ListSSHKeys(ctx)
		if err != nil {
			s.handleRateLimitExceeded(err, "ListSSHKeys")
			return actionError{err: fmt.Errorf("failed to list ssh keys: %w", err)}
		}
		sshKeyID, err := findSSHKey(keys, sshKeyName)
		if err != nil {
			if errors.Is(err, errSSHKeyNotFound) {
				// do not return an error in the reconcile loop as we cannot do anything about this without the intervention
				// of the user. Only after the SSH key has been uploaded correctly, the provisioning can continue.
				// This is why we wait for 5m and then reconcile again to see whether the SSH key exists then.
				msg := fmt.Sprintf("ssh key %q could not be found", sshKeyName)
				conditions.MarkFalse(s.scope.HivelocityCluster, infrav1.CredentialsAvailableCondition, infrav1.HivelocitySSHKeyNotFoundReason, clusterv1.ConditionSeverityWarning, msg)
				record.Warnf(s.scope.HivelocityCluster, "SSHKeyNotFound", msg)
				return actionFailed{}
			}
			return actionError{err: fmt.Errorf("error with ssh keys: %w", err)}
		}
		conditions.MarkTrue(s.scope.HivelocityCluster, infrav1.CredentialsAvailableCondition)
		opts.PublicSshKeyId = sshKeyID
	}

	// Provision the device
	if _, err := s.scope.HVClient.ProvisionDevice(ctx, deviceID, opts); err != nil {
		s.handleRateLimitExceeded(err, "ProvisionDevice")
		record.Warnf(s.scope.HivelocityMachine, "FailedProvisionDevice", "Failed to provision device %d: %s", deviceID, err)
		return actionContinue{delay: 30 * time.Second}
	}

	record.Eventf(s.scope.HivelocityMachine, "SuccessfulStartedProvisionDevice", "Successfully started ProvisionDevice: %d", deviceID)

	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.DeviceProvisioningSucceededCondition)

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
func (s *Service) getDeviceImage(_ context.Context) (string, error) {
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
		s.handleRateLimitExceeded(err, "GetDevice")
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
		msg := fmt.Sprintf("verifyAssociatedDevice failed for device %d: %s", device.DeviceId, err.Error())
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceReadyCondition,
			infrav1.DeviceTagsInvalidReason,
			clusterv1.ConditionSeverityError,
			msg,
		)
		record.Warnf(s.scope.HivelocityMachine, "DeviceTagsInvalid", msg)
		s.scope.HivelocityMachine.SetFailure(capierrors.UpdateMachineError, infrav1.FailureMessageDeviceTagsInvalid)
		return actionComplete{}
	}

	isReloading, _, err := s.getPowerAndReloadingState(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("[actionDeviceProvisioned] getPowerAndReloadingState failed: %w", err)}
	}
	if isReloading {
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceProvisioningSucceededCondition,
			infrav1.DeviceReloadingReason,
			clusterv1.ConditionSeverityWarning,
			fmt.Sprintf("Provisioned device %d is reloading", deviceID),
		)
		return actionContinue{delay: 15 * time.Second}
	}
	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.DeviceProvisioningSucceededCondition)

	// update machine object with infos from device
	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.DeviceReadyCondition)
	s.scope.HivelocityMachine.SetMachineStatus(device)
	if device.PowerStatus == hvclient.PowerStatusOff {
		conditions.MarkFalse(s.scope.HivelocityMachine, infrav1.HivelocityMachineReadyCondition, infrav1.DevicePowerOffReason, clusterv1.ConditionSeverityError, "the device is in power off state")
		s.scope.HivelocityMachine.Status.Ready = false
		return actionContinue{delay: 20 * time.Second}
	}
	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.HivelocityMachineReadyCondition)
	s.scope.HivelocityMachine.Status.Ready = true

	log.V(1).Info("Completed function. This is the final state. The machine is provisioned.",
		"DeviceId", device.DeviceId,
		"PowerStatus", device.PowerStatus,
		"script", utils.FirstN(device.Script, 50))

	return actionComplete{}
}

func (s *Service) verifyAssociatedDevice(device *hv.BareMetalDevice) error {
	deviceTag, err := hvtag.ClusterTagFromList(device.Tags)
	if err != nil {
		return err
	}
	if s.scope.HivelocityCluster.Name != deviceTag.Value {
		return fmt.Errorf("expected %q got %q: %w", s.scope.HivelocityCluster.Name, deviceTag.Value, errWrongClusterTag)
	}

	machineTag, err := hvtag.MachineTagFromList(device.Tags)
	if err != nil {
		return err
	}
	if s.scope.HivelocityMachine.Name != machineTag.Value {
		return fmt.Errorf("expected %q got %q: %w", s.scope.HivelocityMachine.Name, machineTag.Value, errWrongMachineTag)
	}
	return nil
}

// actionDeleteDeviceDeProvision re-provisions a device to remove it from cluster.
func (s *Service) actionDeleteDeviceDeProvision(ctx context.Context) (ar actionResult) {
	log := s.scope.Logger.WithValues("function", "actionDeleteDeviceDeProvision")
	log.V(1).Info("Started function")

	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err, "GetDevice")
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// Nothing to do if device is not found
			s.scope.Info("Unable to locate Hivelocity device by ID or tags")
			record.Warnf(s.scope.HivelocityMachine, "NoDeviceFound", "Unable to find matching Hivelocity device for %s", s.scope.Name())
			return actionComplete{}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	isReloading, isPoweredOn, err := s.getPowerAndReloadingState(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("actionDeleteDeviceDeProvision] getPowerAndReloadingState failed: %w", err)}
	}

	if !isReloading && !isPoweredOn {
		return s.actionDeleteDeviceDeProvisionPowerIsOff(ctx, device)
	}

	deprovisionCondition := conditions.Get(s.scope.HivelocityMachine, infrav1.DeviceDeProvisioningSucceededCondition)

	// handle reloading state
	if isReloading {
		if s.isReloadingTooLong(deprovisionCondition, isPoweredOn) {
			return s.setReloadingTooLongTag(ctx, deviceID, deprovisionCondition.LastTransitionTime)
		}

		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceDeProvisioningSucceededCondition,
			infrav1.DeviceReloadingReason,
			clusterv1.ConditionSeverityWarning,
			fmt.Sprintf("device %d is reloading", deviceID),
		)
		return actionContinue{delay: 1 * time.Minute}
	}

	if !strings.Contains(device.Script, "cloud-init") {
		// This is the dummy OS
		return actionComplete{}
	}

	// handle powered on state

	// if shutdown has been called in the past two minutes already, do not call it again and wait
	if deprovisionCondition != nil && deprovisionCondition.Reason == infrav1.DeviceShutdownCalledReason && !hasTimedOut(&deprovisionCondition.LastTransitionTime, 2*time.Minute) {
		return actionContinue{delay: 30 * time.Second}
	}

	// remove condition to reset the timer - we set the condition anyway again
	conditions.Delete(s.scope.HivelocityMachine, infrav1.DeviceDeProvisioningSucceededCondition)

	err = s.scope.HVClient.ShutdownDevice(ctx, deviceID)
	if err != nil {
		return actionError{err: fmt.Errorf("[actionDeleteDeviceDeProvision] ShutdownDevice failed: %w", err)}
	}

	record.Eventf(s.scope.HivelocityMachine, "SuccessfulShutdownDevice", "Called ShutdownDevice API for %d", deviceID)

	conditions.MarkFalse(
		s.scope.HivelocityMachine,
		infrav1.DeviceDeProvisioningSucceededCondition,
		infrav1.DeviceShutdownCalledReason,
		clusterv1.ConditionSeverityInfo,
		"device shut down has been triggered",
	)
	return actionContinue{delay: 30 * time.Second}
}

func (s *Service) actionDeleteDeviceDeProvisionPowerIsOff(ctx context.Context, device hv.BareMetalDevice) actionResult {
	log := s.scope.Logger.WithValues("function", "actionDeleteDeviceDeProvisionPowerIsOff")
	deviceID := device.DeviceId

	opts := hv.BareMetalDeviceUpdate{
		Hostname:    s.scope.Name() + "-deleted.example.com",
		OsName:      defaultImageName,
		ForceReload: true,
		Script:      "",
		Tags:        device.Tags,
	}

	// Deprovision the device with default image.
	if _, err := s.scope.HVClient.ProvisionDevice(ctx, deviceID, opts); err != nil {
		// TODO: Handle error that machine is not shut down
		s.handleRateLimitExceeded(err, "ProvisionDevice")
		record.Warnf(s.scope.HivelocityMachine, "FailedCallProvisionToDeprovision", "Failed to call provision to deprovision device %d: %s", deviceID, err)
		return actionError{err: fmt.Errorf("failed to de-provision device %d: %s", deviceID, err)}
	}
	msg := fmt.Sprintf("Successfully called provision to deprovision %d with %s",
		deviceID, opts.OsName)
	record.Eventf(s.scope.HivelocityMachine, "SuccessfulCallProvisionToDeprovision", msg)

	log.V(1).Info("Completed function")
	return actionContinue{
		delay: time.Minute,
	}
}

// actionDeleteDeviceDissociate ensures that the device has no tags of machine.
func (s *Service) actionDeleteDeviceDissociate(ctx context.Context) actionResult {
	log := s.scope.Logger.WithValues("function", "actionDeleteDeviceDissociate")
	log.V(1).Info("Started function")

	if s.scope.HivelocityMachine.Spec.ProviderID == nil || *(s.scope.HivelocityMachine.Spec.ProviderID) == "" {
		log.V(1).Info("No ProviderID, no need to dissociate device: actionComplete")
		return actionComplete{}
	}
	deviceID, err := s.scope.HivelocityMachine.DeviceIDFromProviderID()
	if err != nil {
		return actionError{err: fmt.Errorf("failed to get deviceID from providerID: %w", err)}
	}

	device, err := s.scope.HVClient.GetDevice(ctx, deviceID)
	if err != nil {
		s.handleRateLimitExceeded(err, "GetDevice")
		if errors.Is(err, hvclient.ErrDeviceNotFound) {
			// Nothing to do if device is not found
			msg := fmt.Sprintf("[actionDeleteDeviceDissociate] Unable to find matching Hivelocity device %d", deviceID)
			s.scope.Info(msg)
			record.Warnf(s.scope.HivelocityMachine, "NoDeviceFound", msg)
			return actionComplete{}
		}
		return actionError{err: fmt.Errorf("failed to get device: %w", err)}
	}

	if device.PowerStatus != hvclient.PowerStatusOff {
		err = s.scope.HVClient.ShutdownDevice(ctx, deviceID)
		if err != nil {
			s.handleRateLimitExceeded(err, "ShutdownDevice")
			return actionError{err: fmt.Errorf("[actionDeleteDeviceDissociate] failed to shutdown device: %w", err)}
		}
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.DeviceDeProvisioningSucceededCondition,
			infrav1.DeviceShutdownCalledReason,
			clusterv1.ConditionSeverityInfo,
			fmt.Sprintf("shut down for device %d was called", deviceID),
		)
		return actionContinue{delay: 30 * time.Second}
	}

	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.DeviceDeProvisioningSucceededCondition)

	newTags, updated1 := s.scope.HivelocityCluster.DeviceTag().RemoveFromList(device.Tags)
	newTags, updated2 := s.scope.HivelocityMachine.DeviceTag().RemoveFromList(newTags)
	newTags, updated3 := s.scope.DeviceTagMachineType().RemoveFromList(newTags)

	if updated1 || updated2 || updated3 {
		if err := s.scope.HVClient.SetDeviceTags(ctx, device.DeviceId, newTags); err != nil {
			s.handleRateLimitExceeded(err, "SetDeviceTags")
			return actionError{err: fmt.Errorf("failed to set tags: %w", err)}
		}
	}

	log.V(1).Info("Completed function")
	return actionComplete{}
}

func (s *Service) handleRateLimitExceeded(err error, functionName string) {
	if errors.Is(err, hvclient.ErrRateLimitExceeded) {
		msg := fmt.Sprintf("exceeded hivelocity rate limit with calling function: %q", functionName)
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.HivelocityAPIReachableCondition,
			infrav1.RateLimitExceededReason,
			clusterv1.ConditionSeverityWarning,
			msg,
		)
		record.Warnf(s.scope.HivelocityMachine, "RateLimitExceeded", msg)
	}
}
