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
	"testing"

	"github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	mockclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client/mock"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

func newScope() Service {
	client := mockclient.NewMockedHVClientFactory().NewClient("my-api-key")
	s := Service{}
	s.scope = &scope.MachineScope{}
	s.scope.HVClient = client
	s.scope.Cluster = &clusterv1.Cluster{}
	s.scope.Cluster.Name = "capi-cluster"
	s.scope.HivelocityCluster = &v1alpha1.HivelocityCluster{}
	s.scope.HivelocityCluster.Name = "my-cluster"
	s.scope.HivelocityMachine = &v1alpha1.HivelocityMachine{}
	s.scope.HivelocityMachine.Name = "my-machine"
	s.scope.HivelocityMachine.Spec = v1alpha1.HivelocityMachineSpec{Type: "hvCustom"}
	return s
}

func Test_associateDevice(t *testing.T) {
	s := newScope()
	ctx := context.Background()
	gotDevice, err := s.associateDevice(ctx)
	require.NoError(t, err)
	device, err := s.scope.HVClient.GetDevice(ctx, mockclient.FreeDeviceID)
	require.NoError(t, err)

	// Check that device in API has tags set.
	require.ElementsMatch(t, []string{
		"caphv-cluster-name=my-cluster",
		"caphv-device-type=hvCustom",
		"caphv-machine-name=my-machine",
	}, device.Tags)

	// Check that returned device object has tags set.
	require.ElementsMatch(t, []string{
		"caphv-cluster-name=my-cluster",
		"caphv-device-type=hvCustom",
		"caphv-machine-name=my-machine",
	}, gotDevice.Tags)
}

func Test_chooseAvailableFromList(t *testing.T) {
	devices := []*hv.BareMetalDevice{
		&mockclient.NoTagsDevice,
		&mockclient.FreeDevice,
	}
	_, err := chooseAvailableFromList(context.Background(), devices, "fooDeviceType", "my-cluster")
	require.ErrorIs(t, err, errNoDeviceAvailable)
}

func Test_deviceExists(t *testing.T) {
	s := newScope()
	ctx := context.Background()
	exists, err := s.deviceExists(ctx, mockclient.FreeDeviceID)
	require.NoError(t, err)
	require.False(t, exists)

	exists, err = s.deviceExists(ctx, mockclient.WithPrimaryIPDeviceID)
	require.NoError(t, err)
	require.True(t, exists)
}
