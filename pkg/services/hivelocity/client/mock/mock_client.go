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

// Package mock implements mocks for important interfaces like the Hivelocity api.
package mock

import (
	"context"
	"fmt"

	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/maps"
)

// DefaultCPUCores defines the default cpu cores for Hivelocity machines' capacities.
const DefaultCPUCores = 1

// DefaultMemoryInGB defines the default memory in GB for Hivelocity machines' capacities.
const DefaultMemoryInGB = float32(4)

type mockedHVClient struct {
	store serverStore
}

var _ hvclient.Client = &mockedHVClient{}

// NewClient gives reference to the mock client using the in memory store.
func (f *mockedHVClientFactory) NewClient(hvAPIKey string) hvclient.Client {
	var store serverStore
	store.idMap = make(map[int32]*hv.BareMetalDevice)
	devices := []hv.BareMetalDevice{
		{
			Hostname:    "host1-unused",
			Tags:        []string{},
			DeviceId:    1,
			PowerStatus: "ON",
			OsName:      defaultImage,
		},
		{
			Hostname: "host2-other-cluster",
			Tags: []string{
				hvclient.GetClusterTag("other-cluster"),
			},
			DeviceId:    2,
			PowerStatus: "ON",
			OsName:      defaultImage,
		},
		{
			Hostname:    "host3-unused",
			Tags:        []string{},
			DeviceId:    3,
			PowerStatus: "ON",
			OsName:      defaultImage,
		},
	}
	for i := range devices {
		store.idMap[devices[i].DeviceId] = &devices[i]
	}
	client := &mockedHVClient{
		store: store,
	}
	return client
}

// Close implements Close method of HV client interface.
func (c *mockedHVClient) Close() {
	c.store = serverStore{
		idMap: make(map[int32]*hv.BareMetalDevice),
	}
}

type mockedHVClientFactory struct{}

// NewHVClientFactory creates new mock Hivelocity client factories using the in memory store.
func NewHVClientFactory() hvclient.Factory {
	return &mockedHVClientFactory{}
}

var _ = hvclient.Factory(&mockedHVClientFactory{})

// serverStore is an in memory store for the state for the mocked client.
type serverStore struct {
	idMap map[int32]*hv.BareMetalDevice
}

var defaultSSHKey = hv.SshKeyResponse{
	Name:      "testsshkey",
	PublicKey: "AAAAB3NzaC1yc2EAAAADAQABAAABg...",
	SshKeyId:  0,
}

var defaultImage = "Ubuntu 20.x"

func (c *mockedHVClient) ListImages(ctx context.Context, productID int32) ([]string, error) {
	return []string{defaultImage}, nil
}

func (c *mockedHVClient) CreateDevice(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error) {
	if _, found := c.store.idMap[deviceID]; found {
		return hv.BareMetalDevice{}, fmt.Errorf("already exists")
	}

	server := hv.BareMetalDevice{
		Hostname:                 "",
		PrimaryIp:                "",
		Tags:                     []string{},
		CustomIPXEScriptURL:      "",
		LocationName:             "",
		ServiceId:                0,
		DeviceId:                 deviceID,
		ProductName:              "",
		VlanId:                   0,
		Period:                   "",
		PublicSshKeyId:           0,
		Script:                   "",
		PowerStatus:              "",
		CustomIPXEScriptContents: "",
		OrderId:                  0,
		OsName:                   "",
		ProductId:                0,
	}

	// Add server to store
	c.store.idMap[server.DeviceId] = &server
	return server, nil
}

func (c *mockedHVClient) ListDevices(ctx context.Context) ([]*hv.BareMetalDevice, error) {
	return maps.Values(c.store.idMap), nil
}

func (c *mockedHVClient) ShutdownDevice(ctx context.Context, deviceID int32) error {
	if _, found := c.store.idMap[deviceID]; !found {
		return fmt.Errorf("[ShutdownDevice] deviceID %d: %w", deviceID, hvclient.ErrDeviceNotFound)
	}
	c.store.idMap[deviceID].PowerStatus = hvclient.PowerStatusOff
	return nil
}

func (c *mockedHVClient) PowerOnDevice(ctx context.Context, deviceID int32) error {
	if _, found := c.store.idMap[deviceID]; !found {
		return fmt.Errorf("[PowerOnDevice] deviceID %d: %w", deviceID, hvclient.ErrDeviceNotFound)
	}
	c.store.idMap[deviceID].PowerStatus = hvclient.PowerStatusOn
	return nil
}

func (c *mockedHVClient) DeleteDevice(ctx context.Context, deviceID int32) error {
	if _, found := c.store.idMap[deviceID]; !found {
		return fmt.Errorf("[DeleteDevice] deviceID %d: %w", deviceID, hvclient.ErrDeviceNotFound)
	}
	delete(c.store.idMap, deviceID)
	return nil
}

func (c *mockedHVClient) ListSSHKeys(ctx context.Context) ([]hv.SshKeyResponse, error) {
	return []hv.SshKeyResponse{defaultSSHKey}, nil
}
