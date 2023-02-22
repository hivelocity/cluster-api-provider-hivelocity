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

// FreeDeviceID is a deviceID which references a device which is not associated with a node.
const FreeDeviceID = 1

// FreeDevice is  a device which is not associated with a node.
var FreeDevice = hv.BareMetalDevice{
	Hostname:    "host-FreeDevice",
	Tags:        []string{hvclient.TagKeyDeviceType + "=hvCustom"},
	DeviceId:    FreeDeviceID,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// OtherClusterDeviceID is a deviceID which references a device which is from an other cluster.
const OtherClusterDeviceID = 2

// OtherClusterDevice is a device which is from an other cluster.
var OtherClusterDevice = hv.BareMetalDevice{
	Hostname: "host2-OtherClusterDevice",
	Tags: []string{
		hvclient.GetClusterTag("other-cluster"),
	},
	DeviceId:    OtherClusterDeviceID,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// NoTagsDeviceID is a deviceID which references a device which has no tags.
const NoTagsDeviceID = 3

// NoTagsDevice is a device which has no tags.
var NoTagsDevice = hv.BareMetalDevice{
	Hostname:    "host3-unused",
	Tags:        []string{},
	DeviceId:    NoTagsDeviceID,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// WithPrimaryIPDeviceID is a deviceID which references a device which has a PrimaryIp.
const WithPrimaryIPDeviceID = 4

// WithPrimaryIPDevice is a device which has a PrimaryIp.
var WithPrimaryIPDevice = hv.BareMetalDevice{
	Hostname:    "host4-with-ip",
	Tags:        []string{},
	DeviceId:    WithPrimaryIPDeviceID,
	PowerStatus: "ON",
	OsName:      defaultImage,
	PrimaryIp:   "127.0.0,1",
}

type mockedHVClient struct {
	store deviceStore
}

var _ hvclient.Client = &mockedHVClient{}

// NewClient gives reference to the mock client using the in memory store.
func (f *mockedHVClientFactory) NewClient(hvAPIKey string) hvclient.Client {
	var store deviceStore
	devices := []hv.BareMetalDevice{
		FreeDevice,
		OtherClusterDevice,
		NoTagsDevice,
		WithPrimaryIPDevice,
	}
	store.idMap = make(map[int32]*hv.BareMetalDevice, len(devices))
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
	c.store = deviceStore{
		idMap: make(map[int32]*hv.BareMetalDevice),
	}
}

type mockedHVClientFactory struct{}

// NewHVClientFactory creates new mock Hivelocity client factories using the in memory store.
func NewHVClientFactory() hvclient.Factory {
	return &mockedHVClientFactory{}
}

var _ = hvclient.Factory(&mockedHVClientFactory{})

// deviceStore is an in memory store for the state for the mocked client.
type deviceStore struct {
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

func (c *mockedHVClient) ProvisionDevice(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error) {
	device := hv.BareMetalDevice{
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

	// Add device to store
	c.store.idMap[device.DeviceId] = &device
	return device, nil
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

func (c *mockedHVClient) SetTags(ctx context.Context, deviceID int32, tags []string) error {
	device := c.store.idMap[deviceID]
	device.Tags = append([]string(nil), tags...)
	return nil
}

func (c *mockedHVClient) GetDevice(ctx context.Context, deviceID int32) (hv.BareMetalDevice, error) {
	device, ok := c.store.idMap[deviceID]
	if !ok {
		return hv.BareMetalDevice{}, fmt.Errorf("no such device: %d", deviceID)
	}
	return *device, nil
}
