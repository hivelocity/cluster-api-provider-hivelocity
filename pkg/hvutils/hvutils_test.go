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
	"testing"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"
)

func Test_findServerByTags(t *testing.T) {
	type args struct {
		clusterTag string
		machineTag string
		servers    []hv.BareMetalDevice
	}
	tests := []struct {
		name    string
		args    args
		want    *hv.BareMetalDevice
		wantErr error
	}{
		{
			name: "no tags, no servers, no result, no error",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				servers:    []hv.BareMetalDevice{},
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "matching tags, found server",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				servers: []hv.BareMetalDevice{
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
			name: "matching tags, but found two servers",
			args: args{
				clusterTag: "ct=foo",
				machineTag: "mt=bar",
				servers: []hv.BareMetalDevice{
					{
						Tags: []string{"ct=foo", "mt=bar"},
					},
					{
						Tags: []string{"ct=foo", "mt=bar"},
					},
				},
			},
			want:    nil,
			wantErr: errMultipleServerFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindServerByTags(tt.args.clusterTag, tt.args.machineTag, ToPointers(tt.args.servers))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestServerHasTagKey(t *testing.T) {
	type args struct {
		server *hv.BareMetalDevice
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
				server: &hv.BareMetalDevice{},
				tagKey: "machine-name",
			},
			want: false,
		},
		{
			name: "machine with tags",
			args: args{
				server: &hv.BareMetalDevice{
					Tags: []string{
						"machine-name=foo",
						"cluster-name=bar",
					}},
				tagKey: "machine-name",
			},
			want: true,
		},
		{
			name: "machine with tag without equal sign",
			args: args{
				server: &hv.BareMetalDevice{
					Tags: []string{
						"machine-name",
					}},
				tagKey: "machine-name",
			},
			want: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ServerHasTagKey(tt.args.server, tt.args.tagKey); got != tt.want {
				t.Errorf("ServerHasTagKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ToPointers[T any](s []T) []*T {
	ret := make([]*T, 0, len(s))
	for i := range s {
		ret = append(ret, &s[i])
	}
	return ret
}
