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
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/maps"
)

// DefaultCPUCores defines the default cpu cores for Hivelocity machines' capacities.
const DefaultCPUCores = 1

// DefaultMemoryInGB defines the default memory in GB for Hivelocity machines' capacities.
const DefaultMemoryInGB = float32(4)

const (
	// FreeDeviceID is a deviceID which references a device which is not associated with a node.
	FreeDeviceID = 1
)

// FreeDevice is  a device which is not associated with a node.
var FreeDevice = hv.BareMetalDevice{
	Hostname:    "host-FreeDevice",
	Tags:        []string{"caphvlabel:deviceType=hvCustom", "caphv-use=allow"},
	DeviceId:    FreeDeviceID,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// FreeDevicePool1 is  a device which is not associated with a node.
var FreeDevicePool1 = hv.BareMetalDevice{
	Hostname:    "host-FreeDevice",
	Tags:        []string{"caphvlabel:deviceType=pool", "caphv-use=allow"},
	DeviceId:    51,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// FreeDevicePool2 is  a device which is not associated with a node.
var FreeDevicePool2 = hv.BareMetalDevice{
	Hostname:    "host-FreeDevice",
	Tags:        []string{"caphvlabel:deviceType=pool", "caphv-use=allow"},
	DeviceId:    52,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// FreeDevicePool3 is  a device which is not associated with a node.
var FreeDevicePool3 = hv.BareMetalDevice{
	Hostname:    "host-FreeDevice",
	Tags:        []string{"caphvlabel:deviceType=pool", "caphv-use=allow"},
	DeviceId:    53,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// OtherClusterDeviceID is a deviceID which references a device which is from an other cluster.
const OtherClusterDeviceID = 2

// OtherClusterDevice is a device which is from an other cluster.
var OtherClusterDevice = hv.BareMetalDevice{
	Hostname: "host2-OtherClusterDevice",
	Tags: []string{
		hvtag.DeviceTag{Key: hvtag.DeviceTagKeyCluster, Value: "other-cluster"}.ToString(),
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

// CaphvNotAllowDevice is a device which is has no "caphv-use=allow" tag.
var CaphvNotAllowDevice = hv.BareMetalDevice{
	Hostname:    "host-FreeDevice",
	Tags:        []string{"caphvlabel:deviceType=hvCustom"},
	DeviceId:    FreeDeviceID,
	PowerStatus: "ON",
	OsName:      defaultImage,
}

// MultiLabelsDevice is a device which has multiple tags.
var MultiLabelsDevice = hv.BareMetalDevice{
	Hostname: "host3-unused",
	Tags: []string{
		"caphvlabel:foo1=bar1",
		"caphvlabel:foo2=bar2",
		"caphv-use=allow",
	},
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
	store *deviceStore
}

var _ hvclient.Client = &mockedHVClient{}

// NewClient gives reference to the mock client using the in memory store.
func (f *mockedHVClientFactory) NewClient(_ string) hvclient.Client {
	return &mockedHVClient{
		store: f.store,
	}
}

type mockedHVClientFactory struct {
	store *deviceStore
}

// NewMockedHVClientFactory creates new mock Hivelocity client factories using the in memory store.
// We re-use the client, so that changes done in Reconcile() are visible in the
// tests.
func NewMockedHVClientFactory() hvclient.Factory {
	var store deviceStore
	devices := []hv.BareMetalDevice{
		FreeDevice,
		FreeDevicePool1,
		FreeDevicePool2,
		FreeDevicePool3,
		OtherClusterDevice,
		NoTagsDevice,
		WithPrimaryIPDevice,
	}
	store.idMap = make(map[int32]hv.BareMetalDevice, len(devices))
	for i := range devices {
		store.idMap[devices[i].DeviceId] = devices[i]
	}
	return &mockedHVClientFactory{store: &store}
}

var _ = hvclient.Factory(&mockedHVClientFactory{})

// deviceStore is an in memory store for the state for the mocked client.
type deviceStore struct {
	idMap map[int32]hv.BareMetalDevice
}

var defaultSSHKey = hv.SshKeyResponse{
	Name:      "testsshkey",
	PublicKey: "AAAAB3NzaC1yc2EAAAADAQABAAABg...",
	SshKeyId:  0,
}

var defaultImage = "Ubuntu 20.x"

func (c *mockedHVClient) ListImages(_ context.Context, _ int32) ([]string, error) {
	return []string{defaultImage}, nil
}

func (c *mockedHVClient) ProvisionDevice(_ context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error) {
	device, ok := c.store.idMap[deviceID]
	if !ok {
		return hv.BareMetalDevice{}, fmt.Errorf("[ProvisionDevice] deviceID %d unknown", deviceID)
	}
	device.Tags = opts.Tags
	device.PowerStatus = hvclient.PowerStatusOn
	c.store.idMap[deviceID] = device
	return device, nil
}

func (c *mockedHVClient) ListDevices(_ context.Context) ([]hv.BareMetalDevice, error) {
	return maps.Values(c.store.idMap), nil
}

func (c *mockedHVClient) ShutdownDevice(_ context.Context, deviceID int32) error {
	device, found := c.store.idMap[deviceID]
	if !found {
		return fmt.Errorf("[ShutdownDevice] deviceID %d: %w", deviceID, hvclient.ErrDeviceNotFound)
	}
	if device.PowerStatus == hvclient.PowerStatusOff {
		return hvclient.ErrDeviceShutDownAlready
	}

	device.PowerStatus = hvclient.PowerStatusOff
	c.store.idMap[deviceID] = device
	return nil
}

func (c *mockedHVClient) PowerOnDevice(_ context.Context, deviceID int32) error {
	device, found := c.store.idMap[deviceID]
	if !found {
		return fmt.Errorf("[PowerOnDevice] deviceID %d: %w", deviceID, hvclient.ErrDeviceNotFound)
	}
	if device.PowerStatus == hvclient.PowerStatusOn {
		return hvclient.ErrDeviceTurnedOnAlready
	}

	device.PowerStatus = hvclient.PowerStatusOn
	c.store.idMap[deviceID] = device
	return nil
}

func (c *mockedHVClient) ListSSHKeys(_ context.Context) ([]hv.SshKeyResponse, error) {
	return []hv.SshKeyResponse{defaultSSHKey}, nil
}

func (c *mockedHVClient) SetDeviceTags(_ context.Context, deviceID int32, tags []string) error {
	device := c.store.idMap[deviceID]
	device.Tags = append([]string(nil), tags...)
	c.store.idMap[deviceID] = device
	return nil
}

func (c *mockedHVClient) GetDevice(_ context.Context, deviceID int32) (hv.BareMetalDevice, error) {
	device, ok := c.store.idMap[deviceID]
	if !ok {
		return hv.BareMetalDevice{}, hvclient.ErrDeviceNotFound
	}
	return device, nil
}

func (c *mockedHVClient) GetDeviceDump(ctx context.Context, deviceID int32) (hv.DeviceDump, error) {
	device, _ := c.GetDevice(ctx, deviceID)
	return hv.DeviceDump{
		DeviceId:           deviceID,
		Name:               "",
		Status:             "",
		DeviceType:         "",
		DeviceTypeGroup:    "",
		PowerStatus:        device.PowerStatus,
		HasCancellation:    false,
		IsManaged:          false,
		IsReload:           false,
		MonitorsUp:         0,
		MonitorsTotal:      0,
		ManagedAlertsTotal: 0,
		Ports:              []interface{}{},
		Hostname:           "",
		IpmiEnabled:        false,
		DisplayedTags:      []interface{}{},
		Tags:               []string{},
		Location:           nil,
		NetworkAutomation:  nil,
		PrimaryIp:          "",
		IpmiAddress:        nil,
		ServiceMonitors:    []string{},
		BillingInfo:        nil,
		ServicePlan:        0,
		LastInvoiceId:      0,
		SelfProvisioning:   false,
		Metadata:           nil,
		SpsStatus:          "",
	}, nil
}
