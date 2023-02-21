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
	// InstanceReadyCondition reports on current status of the instance. Ready indicates the instance is in a Running state.
	InstanceReadyCondition clusterv1.ConditionType = "InstanceReady"
	// InstanceTerminatedReason instance is in a terminated state.
	InstanceTerminatedReason = "InstanceTerminated"
	// DeviceOffReason instance is off.
	DeviceOffReason = "DeviceOff"
	// InstanceAsControlPlaneUnreachableReason control plane is (not yet) reachable.
	InstanceAsControlPlaneUnreachableReason = "InstanceAsControlPlaneUnreachable"
)

const (
	// InstanceBootstrapReadyCondition reports on current status of the instance. BootstrapReady indicates the bootstrap is ready.
	InstanceBootstrapReadyCondition clusterv1.ConditionType = "InstanceBootstrapReady"
	// InstanceBootstrapNotReadyReason bootstrap not ready yet.
	InstanceBootstrapNotReadyReason = "InstanceBootstrapNotReady"
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
