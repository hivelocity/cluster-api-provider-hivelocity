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

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	// DeviceReadyCondition reports on current status of the device. Ready indicates the device is in a Running state.
	DeviceReadyCondition clusterv1.ConditionType = "DeviceReady"

	// DeviceTerminatedReason device is in a terminated state.
	DeviceTerminatedReason = "DeviceTerminated"

	// DeviceOffReason device is off.
	DeviceOffReason = "DeviceOff"

	// DeviceNotFoundReason (Severity=Error) documents a HivelocityMachine controller detecting
	// the underlying device has been deleted unexpectedly.
	DeviceNotFoundReason = "DeviceNotFound"
)

const (
	// MachineBootstrapReadyCondition reports on current status of the machine. BootstrapReady indicates the bootstrap is ready.
	MachineBootstrapReadyCondition clusterv1.ConditionType = "MachineBootstrapReady"
	// MachineBootstrapNotReadyReason bootstrap not ready yet.
	MachineBootstrapNotReadyReason = "MachineBootstrapNotReady"
)

const (
	// NetworkAttached reports on whether there is a network attached to the cluster.
	NetworkAttached clusterv1.ConditionType = "NetworkAttached"
	// NetworkDisabledReason indicates that network is disabled.
	NetworkDisabledReason = "NetworkDisabled"
	// NetworkUnreachableReason indicates that network is unreachable.
	NetworkUnreachableReason = "NetworkUnreachable"
)

const (
	// HivelocityClusterReady reports on whether the Hivelocity cluster is in ready state.
	HivelocityClusterReady clusterv1.ConditionType = "HivelocityClusterReady"
	// HivelocitySecretUnreachableReason indicates that Hivelocity secret is unreachable.
	HivelocitySecretUnreachableReason = "HivelocitySecretUnreachable" // #nosec
	// HivelocityCredentialsInvalidReason indicates that credentials for Hivelocity are invalid.
	HivelocityCredentialsInvalidReason = "HivelocityCredentialsInvalid" // #nosec
)

const (
	// RateLimitExceeded reports whether the rate limit has been reached.
	RateLimitExceeded clusterv1.ConditionType = "RateLimitExceeded"
	// RateLimitNotReachedReason indicates that the rate limit is not reached yet.
	RateLimitNotReachedReason = "RateLimitNotReached"
)

const (
	// HivelocityBareMetalHostReady reports on whether the Hivelocity cluster is in ready state.
	HivelocityBareMetalHostReady clusterv1.ConditionType = "HivelocityBareMetalHostReady"
)
