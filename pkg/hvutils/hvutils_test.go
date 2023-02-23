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

// Package hvutils implements helper functions for the HV API.

package hvutils

import (
	"context"
	"testing"

	mockclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client/mock"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"
)

func Test_FindDeviceByTags(t *testing.T) {
	type args struct {
		clusterTag string
		machineTag string
		devices    []hv.BareMetalDevice
	}
	tests := []struct {
		name    string
		args    args
		want    *hv.BareMetalDevice
		wantErr error
	}{
		{
			name: "no tags, no devices, no result, no error",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				devices:    []hv.BareMetalDevice{},
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "matching tags, found device",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				devices: []hv.BareMetalDevice{
					{
						Tags: []string{"ct=foo", "mt=bar"},
					},
					{
						Tags: []string{"ct=other", "mt=bar"},
					},
				},
			},
			want: &hv.BareMetalDevice{
				Tags: []string{"ct=foo", "mt=bar"},
			},
			wantErr: nil,
		},
		{
			name: "matching tags, but found two devices",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				devices: []hv.BareMetalDevice{
					{
						Tags: []string{"ct=foo", "mt=bar"},
					},
					{
						Tags: []string{"ct=foo", "mt=bar"},
					},
				},
			},
			want:    nil,
			wantErr: errMultipleDevicesFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindDeviceByTags(tt.args.clusterTag, tt.args.machineTag, toPointers(tt.args.devices))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDeviceHasTagKey(t *testing.T) {
	type args struct {
		device *hv.BareMetalDevice
		tagKey string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "machine without tags",
			args: args{
				device: &hv.BareMetalDevice{},
				tagKey: "machine-name",
			},
			want: false,
		},
		{
			name: "machine with tags",
			args: args{
				device: &hv.BareMetalDevice{
					Tags: []string{
						"machine-name=foo",
						"cluster-name=bar",
					},
				},
				tagKey: "machine-name",
			},
			want: true,
		},
		{
			name: "machine with tag without equal sign",
			args: args{
				device: &hv.BareMetalDevice{
					Tags: []string{
						"machine-name",
					},
				},
				tagKey: "machine-name",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeviceHasTagKey(tt.args.device, tt.args.tagKey); got != tt.want {
				t.Errorf("DeviceHasTagKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func toPointers[T any](s []T) []*T {
	ret := make([]*T, 0, len(s))
	for i := range s {
		ret = append(ret, &s[i])
	}
	return ret
}

func TestAddTags(t *testing.T) {
	client := mockclient.NewHVClientFactory().NewClient("my-api-key")
	ctx := context.Background()
	device, err := client.GetDevice(ctx, mockclient.NoTagsDeviceID)
	require.NoError(t, err)
	require.Equal(t, []string{}, device.Tags)

	err = AddTags(ctx, client, &device, []string{"foo"})
	require.NoError(t, err)
	device, err = client.GetDevice(ctx, mockclient.NoTagsDeviceID)
	require.NoError(t, err)
	require.Equal(t, []string{"foo"}, device.Tags)

	err = AddTags(ctx, client, &device, []string{"bar"})
	require.NoError(t, err)
	device, err = client.GetDevice(ctx, mockclient.NoTagsDeviceID)
	require.NoError(t, err)
	require.Equal(t, []string{"bar", "foo"}, device.Tags)

	// don't create duplicates
	err = AddTags(ctx, client, &device, []string{"bar"})
	require.NoError(t, err)
	device, err = client.GetDevice(ctx, mockclient.NoTagsDeviceID)
	require.NoError(t, err)
	require.Equal(t, []string{"bar", "foo"}, device.Tags)

	err = AddTags(ctx, client, &device, []string{})
	require.NoError(t, err)
	device, err = client.GetDevice(ctx, mockclient.NoTagsDeviceID)
	require.NoError(t, err)
	require.Equal(t, []string{"bar", "foo"}, device.Tags)

}

func Test_AssociateDevice(t *testing.T) {
	client := mockclient.NewHVClientFactory().NewClient("my-api-key")
	ctx := context.Background()
	device, err := client.GetDevice(ctx, mockclient.FreeDeviceID)
	require.NoError(t, err)
	err = AssociateDevice(ctx, client, &device, "my-cluster", "my-machine")
	require.NoError(t, err)
	device, err = client.GetDevice(ctx, mockclient.FreeDeviceID)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"caphv-cluster-name=my-cluster",
		"caphv-device-type=hvCustom",
		"caphv-machine-name=my-machine",
	}, device.Tags)
}

func Test_FindAndAssociateDevice(t *testing.T) {
	client := mockclient.NewHVClientFactory().NewClient("my-api-key")
	ctx := context.Background()
	device, err := FindAndAssociateDevice(ctx, client, "my-cluster", "my-machine")
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"caphv-cluster-name=my-cluster",
		"caphv-device-type=hvCustom",
		"caphv-machine-name=my-machine",
	}, device.Tags)
}
