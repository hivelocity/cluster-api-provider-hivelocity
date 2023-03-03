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

// Package hvclient provides the interfaces to communicate with
// the API of Hivelocity.
// We use interfaces to make mocking easier.
package hvclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	hv "github.com/hivelocity/hivelocity-client-go/client"
)

// PowerStatusOff is "OFF".
const PowerStatusOff = "OFF"

// PowerStatusOn is "ON".
const PowerStatusOn = "ON"

// Client collects all methods used by the controller in the Hivelocity API.
type Client interface {
	Close()
	PowerOnDevice(ctx context.Context, deviceID int32) error
	ProvisionDevice(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error)
	ListDevices(context.Context) ([]*hv.BareMetalDevice, error)
	ShutdownDevice(ctx context.Context, deviceID int32) error
	DeleteDevice(ctx context.Context, deviceID int32) error
	ListImages(ctx context.Context, productID int32) ([]string, error)
	ListSSHKeys(context.Context) ([]hv.SshKeyResponse, error)
	GetDevice(ctx context.Context, deviceID int32) (hv.BareMetalDevice, error)

	// SetTags sets the tags to the given list.
	SetTags(ctx context.Context, deviceID int32, tags []string) error
}

// Factory is the interface for creating new Client objects.
type Factory interface {
	NewClient(hvAPIKey string) Client
}

// HivelocityFactory implements the Factory interface.
type HivelocityFactory struct{}

// ErrDeviceNotFound gets returned if no matching device was found.
var ErrDeviceNotFound = fmt.Errorf("device was not found")

var _ Factory = &HivelocityFactory{}

// NewClient creates new Hivelocity clients.
func (f *HivelocityFactory) NewClient(hvAPIKey string) Client {
	config := hv.NewConfiguration()
	config.AddDefaultHeader("X-API-KEY", hvAPIKey)
	apiClient := hv.NewAPIClient(config)
	return &realClient{
		client: apiClient,
	}
}

type realClient struct {
	client *hv.APIClient
}

var _ Client = &realClient{}

// Close implements the Close method of the HVClient interface.
func (c *realClient) Close() {}

func (c *realClient) GetDevice(ctx context.Context, deviceID int32) (hv.BareMetalDevice, error) {
	// https://developers.hivelocity.net/reference/get_bare_metal_device_id_resource
	device, _, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceID, nil) //nolint:bodyclose // Close() gets done in client
	return device, err
}

func (c *realClient) SetTags(ctx context.Context, deviceID int32, tags []string) error {
	// https://developers.hivelocity.net/reference/put_device_tag_id_resource
	// Existing Tags will be removed by the HV API.
	deviceTags := hv.DeviceTag{
		Tags: tags,
	}
	_, _, err := c.client.DeviceApi.PutDeviceTagIdResource(ctx, deviceID, deviceTags, nil) //nolint:bodyclose // Close() gets done in client
	return err
}

func (c *realClient) PowerOnDevice(ctx context.Context, deviceID int32) error {
	return nil // todo
}

func (c *realClient) ProvisionDevice(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error) {
	device, err := c.GetDevice(ctx, deviceID)
	if err != nil {
		return hv.BareMetalDevice{}, fmt.Errorf("[ProvisionDevice] GetDevice() failed. deviceID %d: %w",
			deviceID, err)
	}

	// FIXME: state-machine: After powering off, we need to wait until the device is off.
	if device.PowerStatus == PowerStatusOn {
		_, _, err := c.client.DeviceApi.PostPowerResource(ctx, deviceID, "shutdown", nil) //nolint:bodyclose // Close() gets done in client
		if err != nil {
			return hv.BareMetalDevice{}, fmt.Errorf("[ProvisionDevice] PostPowerResource() failed. deviceID %d: %w",
				deviceID, err)
		}
	}

	// https://developers.hivelocity.net/reference/put_bare_metal_device_id_resource
	device, _, err = c.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(ctx, deviceID, opts, nil) //nolint:bodyclose // Close() gets done in client
	return device, err
}

func (c *realClient) ListDevices(ctx context.Context) ([]*hv.BareMetalDevice, error) {
	devices, _, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceResource(ctx, nil) //nolint:bodyclose // Close() gets done in client
	ret := make([]*hv.BareMetalDevice, 0, len(devices))
	for i := range devices {
		ret = append(ret, &devices[i])
	}
	return ret, err
}

func (c *realClient) DeleteDevice(ctx context.Context, deviceID int32) error {
	return fmt.Errorf("todo DeleteDevice")
}

func (c *realClient) ShutdownDevice(ctx context.Context, deviceID int32) error {
	return fmt.Errorf("todo ShutdownDevice")
}

func (c *realClient) ListImages(ctx context.Context, productID int32) ([]string, error) {
	// https://developers.hivelocity.net/reference/get_product_operating_systems_resource
	opts, _, err := c.client.ProductApi.GetProductOperatingSystemsResource(ctx, productID, nil) //nolint:bodyclose // Close() gets done in client
	ret := make([]string, 0, len(opts))
	if err != nil {
		return []string{}, err
	}
	for i := range opts {
		ret = append(ret, opts[i].Name)
	}
	return ret, nil
}

func (c *realClient) ListSSHKeys(ctx context.Context) ([]hv.SshKeyResponse, error) {
	// https://developers.hivelocity.net/reference/get_ssh_key_resource
	sshKeys, _, err := c.client.SshKeyApi.GetSshKeyResource(ctx, nil) //nolint:bodyclose // Close() gets done in client
	return sshKeys, err
}

// DeviceStatus specifies a device's status.
type DeviceStatus string

const (
	// DeviceStatusInitializing is the status when a device is initializing.
	DeviceStatusInitializing DeviceStatus = "initializing" // TODO AFAIK HV does not provide these detailed infos

	// DeviceStatusOff is the status when a device is off.
	DeviceStatusOff DeviceStatus = "off"

	// DeviceStatusRunning is the status when a device is running.
	DeviceStatusRunning DeviceStatus = "running"

	// DeviceStatusStarting is the status when a device is being started.
	DeviceStatusStarting DeviceStatus = "starting"

	// DeviceStatusStopping is the status when a device is being stopped.
	DeviceStatusStopping DeviceStatus = "stopping"

	// DeviceStatusMigrating is the status when a device is being migrated.
	DeviceStatusMigrating DeviceStatus = "migrating"

	// DeviceStatusRebuilding is the status when a device is being rebuilt.
	DeviceStatusRebuilding DeviceStatus = "rebuilding"

	// DeviceStatusDeleting is the status when a device is being deleted.
	DeviceStatusDeleting DeviceStatus = "deleting"

	// DeviceStatusUnknown is the status when a device's state is unknown.
	DeviceStatusUnknown DeviceStatus = "unknown"

	// TagKeyMachineName is the prefix for HV tags for machine names.
	TagKeyMachineName = "caphv-machine-name"

	// TagKeyClusterName is the prefix for HV tags for cluster names.
	TagKeyClusterName = "caphv-cluster-name"

	// TagKeyDeviceType is the prefix for HV tags for device types.
	TagKeyDeviceType = "caphv-device-type"
)

// GetMachineTag create tag for HV API. Example: "mymachine" --> "caphv-machine-name=mymachine".
func GetMachineTag(machineName string) string {
	return TagKeyMachineName + "=" + machineName
}

// GetClusterTag create tag for HV API. Example: "mycluster" --> "caphv-cluster-name=mycluster".
func GetClusterTag(clusterName string) string {
	return TagKeyClusterName + "=" + clusterName
}

// IsRateLimitExceededError returns true, if the Hivelocity rate limit was reached.
func IsRateLimitExceededError(err error) bool {
	var swaggerErr hv.GenericSwaggerError
	if !errors.As(err, &swaggerErr) {
		return false
	}
	if strings.HasPrefix(swaggerErr.Error(), fmt.Sprint(http.StatusTooManyRequests)) {
		return true
	}
	return false
}
