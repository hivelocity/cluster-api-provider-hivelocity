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
	"testing"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	mockclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client/mock"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_findAvailableDeviceFromList(t *testing.T) {
	tests := []struct {
		description string
		devices     []hv.BareMetalDevice
		deviceType  infrav1.DeviceSelector
		wantMachine *hv.BareMetalDevice
		shouldNil   bool
	}{
		{
			description: "checks that no device is selected if no DeviceSelector matches",
			devices: []hv.BareMetalDevice{
				mockclient.NoTagsDevice,
				mockclient.FreeDevice,
			},
			deviceType: infrav1.DeviceSelector{
				MatchLabels: map[string]string{
					"foo1": "bar1",
				},
			},
			shouldNil: true,
		},
		{
			description: "check no device selected if device has no caphv-use=allow tag",
			devices: []hv.BareMetalDevice{
				mockclient.CaphNotAllowDevice,
			},
			deviceType: infrav1.DeviceSelector{
				MatchLabels: map[string]string{
					"deviceType": "hvCustom",
				},
			},
			shouldNil: true,
		},
		{
			description: "selects device which has all the tags",
			devices: []hv.BareMetalDevice{
				mockclient.MultiLabelsDevice,
				mockclient.FreeDevice,
			},
			deviceType: infrav1.DeviceSelector{
				MatchLabels: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
			},
			wantMachine: &mockclient.MultiLabelsDevice,
			shouldNil:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			device := findAvailableDeviceFromList(test.devices, test.deviceType, "my-cluster")

			if test.shouldNil {
				require.Nil(t, device)
			} else {
				require.NotNil(t, device)
				require.Equal(t, test.wantMachine, device)
			}
		})
	}
}

func TestService_verifyAssociatedDevice(t *testing.T) {
	service := Service{
		scope: &scope.MachineScope{
			ClusterScope: scope.ClusterScope{HivelocityCluster: &infrav1.HivelocityCluster{ObjectMeta: metav1.ObjectMeta{Name: "dummy-cluster"}}},
			HivelocityMachine: &infrav1.HivelocityMachine{
				ObjectMeta: metav1.ObjectMeta{Name: "dummy-machine"},
			},
		},
	}
	// test hv-labels are ok
	device := hv.BareMetalDevice{
		Tags: []string{
			string(hvtag.DeviceTagKeyCluster) + "=dummy-cluster",
			string(hvtag.DeviceTagKeyMachine) + "=dummy-machine",
		},
	}
	err := service.verifyAssociatedDevice(&device)
	require.NoError(t, err)

	// wrong cluster
	device = hv.BareMetalDevice{
		Tags: []string{
			string(hvtag.DeviceTagKeyCluster) + "=other-cluster",
			string(hvtag.DeviceTagKeyMachine) + "=dummy-machine",
		},
	}
	err = service.verifyAssociatedDevice(&device)
	require.Error(t, err)
	require.Equal(t, `expected "dummy-cluster" got "other-cluster": machine has wrong cluster tag`, err.Error())

	// wrong machine
	device = hv.BareMetalDevice{
		Tags: []string{
			string(hvtag.DeviceTagKeyCluster) + "=dummy-cluster",
			string(hvtag.DeviceTagKeyMachine) + "=other-machine",
		},
	}
	err = service.verifyAssociatedDevice(&device)
	require.Error(t, err)
	require.Equal(t, `expected "dummy-machine" got "other-machine": machine has wrong machine tag`, err.Error())

	// missing tags
	device = hv.BareMetalDevice{
		Tags: []string{},
	}
	err = service.verifyAssociatedDevice(&device)
	require.ErrorIs(t, err, hvtag.ErrDeviceTagNotFound)
}
