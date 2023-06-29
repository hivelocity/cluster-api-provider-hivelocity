/*
Copyright 2022 The Kubernetes Authors.

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

package scope

import (
	"context"
	"fmt"
	"hash/crc32"
	"sort"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
)

// MachineScopeParams defines the input parameters used to create a new Scope.
type MachineScopeParams struct {
	ClusterScopeParams
	Machine           *clusterv1.Machine
	HivelocityMachine *infrav1.HivelocityMachine
}

// ErrBootstrapDataNotReady return an error if no bootstrap data is ready.
var ErrBootstrapDataNotReady = errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")

// ErrFailureDomainNotFound returns an error if no region is found.
var ErrFailureDomainNotFound = errors.New("error no failure domain available")

// NewMachineScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewMachineScope(ctx context.Context, params MachineScopeParams) (*MachineScope, error) {
	if params.Machine == nil {
		return nil, errors.New("failed to generate new scope from nil Machine")
	}
	if params.HivelocityMachine == nil {
		return nil, errors.New("failed to generate new scope from nil HivelocityMachine")
	}

	cs, err := NewClusterScope(ctx, params.ClusterScopeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err) // question: does the err message fit?
	}

	cs.patchHelper, err = patch.NewHelper(params.HivelocityMachine, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}

	return &MachineScope{
		ClusterScope:      *cs,
		Machine:           params.Machine,
		HivelocityMachine: params.HivelocityMachine,
	}, nil
}

// MachineScope defines the basic context for an actuator to operate upon.
type MachineScope struct {
	ClusterScope
	Machine           *clusterv1.Machine
	HivelocityMachine *infrav1.HivelocityMachine
}

// Close closes the current scope persisting the cluster configuration and status.
func (m *MachineScope) Close(ctx context.Context) error {
	return m.patchHelper.Patch(ctx, m.HivelocityMachine)
}

// IsControlPlane returns true if the machine is a control plane.
func (m *MachineScope) IsControlPlane() bool {
	return util.IsControlPlaneMachine(m.Machine)
}

// Name returns the HivelocityMachine name.
func (m *MachineScope) Name() string {
	return m.HivelocityMachine.Name
}

// Namespace returns the namespace name.
func (m *MachineScope) Namespace() string {
	return m.HivelocityMachine.Namespace
}

// PatchObject persists the machine spec and status.
func (m *MachineScope) PatchObject(ctx context.Context) error {
	return m.patchHelper.Patch(ctx, m.HivelocityMachine)
}

// DeviceTagMachineType returns a DeviceTag object for the cluster tag.
func (m *MachineScope) DeviceTagMachineType() hvtag.DeviceTag {
	var value string
	if m.IsControlPlane() {
		value = "control_plane"
	} else {
		value = "worker"
	}
	return hvtag.DeviceTag{
		Key:   hvtag.DeviceTagKeyMachineType,
		Value: value,
	}
}

// SetError sets the ErrorMessage and ErrorReason fields on the machine and logs
// the message. It assumes the reason is invalid configuration, since that is
// currently the only relevant MachineStatusError choice.
func (m *MachineScope) SetError(message string, reason capierrors.MachineStatusError) {
	m.HivelocityMachine.Status.FailureMessage = &message
	m.HivelocityMachine.Status.FailureReason = &reason
}

// IsBootstrapDataReady checks the readiness of a capi machine's bootstrap data.
func (m *MachineScope) IsBootstrapDataReady(ctx context.Context) bool {
	return m.Machine.Spec.Bootstrap.DataSecretName != nil
}

// GetFailureDomain returns the machine's failure domain or a default one based on a hash.
func (m *MachineScope) GetFailureDomain() (string, error) {
	if m.Machine.Spec.FailureDomain != nil {
		return *m.Machine.Spec.FailureDomain, nil
	}

	failureDomainNames := make([]string, 0, len(m.Cluster.Status.FailureDomains))
	for fdName, fd := range m.Cluster.Status.FailureDomains {
		// filter out zones if we are a control plane and the cluster object
		// wants to avoid control planes in that zone
		if m.IsControlPlane() && !fd.ControlPlane {
			continue
		}
		failureDomainNames = append(failureDomainNames, fdName)
	}

	if len(failureDomainNames) == 0 {
		return "", ErrFailureDomainNotFound
	}
	if len(failureDomainNames) == 1 {
		return failureDomainNames[0], nil
	}

	sort.Strings(failureDomainNames)

	// assign the node a zone based on a hash
	pos := int(crc32.ChecksumIEEE([]byte(m.HivelocityMachine.Name))) % len(failureDomainNames)

	return failureDomainNames[pos], nil
}

// GetRawBootstrapData returns the bootstrap data from the secret in the Machine's bootstrap.dataSecretName.
func (m *MachineScope) GetRawBootstrapData(ctx context.Context) ([]byte, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil, ErrBootstrapDataNotReady
	}

	key := types.NamespacedName{Namespace: m.Namespace(), Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	// Look for secret in the filtered cache
	var secret corev1.Secret
	if err := m.Client.Get(ctx, key, &secret); err != nil {
		return nil, fmt.Errorf("failed to find bootstrap secret %+v: %w", key, err)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return nil, errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return value, nil
}
