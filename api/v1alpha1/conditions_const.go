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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// DeviceReadyCondition reports on current status of the device. Ready indicates the device is in a Running state.
	DeviceReadyCondition clusterv1.ConditionType = "DeviceReady"
	// DeviceNotFoundReason (Severity=Error) documents a HivelocityMachine controller detecting
	// the underlying device cannot be found anymore.
	DeviceNotFoundReason = "DeviceNotFound"
	// DeviceTagsInvalidReason documents a HivelocityMachine controller detecting invalid device tags.
	DeviceTagsInvalidReason = "DeviceTagsInvalid"
)

const (
	// MachineBootstrapReadyCondition reports on current status of the machine. BootstrapReady indicates the bootstrap is ready.
	MachineBootstrapReadyCondition clusterv1.ConditionType = "MachineBootstrapReady"
	// MachineBootstrapNotReadyReason bootstrap not ready yet.
	MachineBootstrapNotReadyReason = "MachineBootstrapNotReady"
)

const (
	// HivelocityAPIReachableCondition reports whether the Hivelocity APIs are reachable.
	HivelocityAPIReachableCondition clusterv1.ConditionType = "RateLimitExceeded"
	// RateLimitExceededReason indicates that a rate limit has been exceeded.
	RateLimitExceededReason = "RateLimitExceeded"
)

const (
	// CredentialsAvailableCondition reports on whether the Hivelocity cluster is in ready state.
	CredentialsAvailableCondition clusterv1.ConditionType = "CredentialsAvailable"

	// HivelocitySSHKeyNotFoundReason indicates that ssh for Hivelocity not found.
	HivelocitySSHKeyNotFoundReason = "HivelocitySSHKeyNotFound"

	// HivelocityWrongAPIKeyReason indicates that API for Hivelocity is wrong.
	HivelocityWrongAPIKeyReason = "HivelocityWrongAPIKey"

	// HivelocitySecretUnreachableReason indicates that Hivelocity secret is unreachable.
	HivelocitySecretUnreachableReason = "HivelocitySecretUnreachable" // #nosec

	// HivelocityCredentialsInvalidReason indicates that credentials for Hivelocity are invalid.
	HivelocityCredentialsInvalidReason = "HivelocityCredentialsInvalid" // #nosec
)

const (
	// HivelocityMachineReadyCondition reports on whether the Hivelocity machine is in ready state.
	HivelocityMachineReadyCondition clusterv1.ConditionType = "HivelocityMachineReady"

	// DevicePowerOffReason indicates that device is power off.
	DevicePowerOffReason = "DevicePowerOff"
)

const (
	// DeviceAssociateSucceededCondition indicates that a device has been associated.
	DeviceAssociateSucceededCondition clusterv1.ConditionType = "DeviceAssociateSucceeded"

	// NoAvailableDeviceReason indicates that there is no available device.
	NoAvailableDeviceReason = "NoAvailableDevice"
)

const (
	// RateLimitExceeded reports whether the rate limit has been reached.
	RateLimitExceeded clusterv1.ConditionType = "RateLimitExceeded"
	// RateLimitNotReachedReason indicates that the rate limit is not reached yet.
	RateLimitNotReachedReason = "RateLimitNotReached"
)

const (
	// TargetClusterSecretReadyCondition reports on whether the hivelocity secret in the target cluster is ready.
	TargetClusterSecretReadyCondition clusterv1.ConditionType = "TargetClusterSecretReady"

	// TargetSecretSyncFailedReason indicates that the target secret could not be synced.
	TargetSecretSyncFailedReason = "TargetSecretSyncFailed"
)

const (
	// TargetClusterReadyCondition reports on whether the kubeconfig in the target cluster is ready.
	TargetClusterReadyCondition clusterv1.ConditionType = "TargetClusterReady"

	// TargetClusterControlPlaneNotReadyReason indicates that the target cluster's control plane is not ready yet.
	TargetClusterControlPlaneNotReadyReason = "TargetClusterControlPlaneNotReady"

	// TargetClusterCreateFailedReason indicates that the target cluster could not be created.
	TargetClusterCreateFailedReason = "TargetClusterCreateFailed"

	// KubeConfigNotFoundReason indicates that the Kubeconfig could not be found.
	KubeConfigNotFoundReason = "KubeConfigNotFound"

	// KubeAPIServerNotRespondingReason indicates that the api server cannot be reached.
	KubeAPIServerNotRespondingReason = "KubeAPIServerNotResponding"
)
