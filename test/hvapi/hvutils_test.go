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

package main

import (
	"testing"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"
)

func Test_findDeviceByTags(t *testing.T) {
	type args struct {
		clusterTag string
		machineTag string
		devices    []hv.BareMetalDevice
	}
	tests := []struct {
		name    string
		args    args
		want    hv.BareMetalDevice
		wantErr error
	}{
		{
			name: "no tags, no devices, no result, no error",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				devices:    []hv.BareMetalDevice{},
			},
			want:    hv.BareMetalDevice{},
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
			want: hv.BareMetalDevice{
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
			want:    hv.BareMetalDevice{},
			wantErr: errMultipleDevicesFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findDeviceByTags(tt.args.clusterTag, tt.args.machineTag, tt.args.devices)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
