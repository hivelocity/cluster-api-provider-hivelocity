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
	"time"

	"github.com/antihax/optional"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/utils"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PowerStatusOff is "OFF".
const PowerStatusOff = "OFF"

// PowerStatusOn is "ON".
const PowerStatusOn = "ON"

// Client collects all methods used by the controller in the Hivelocity API.
type Client interface {
	PowerOnDevice(ctx context.Context, deviceID int32) error
	ProvisionDevice(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error)
	ListDevices(context.Context) ([]hv.BareMetalDevice, error)
	ShutdownDevice(ctx context.Context, deviceID int32) error
	ListImages(ctx context.Context, productID int32) ([]string, error)
	ListSSHKeys(context.Context) ([]hv.SshKeyResponse, error)

	// GetDevice return the device. If the device is not found ErrDeviceNotFound is returned.
	GetDevice(ctx context.Context, deviceID int32) (hv.BareMetalDevice, error)

	// SetDeviceTags sets the tags to the given list.
	SetDeviceTags(ctx context.Context, deviceID int32, tags []string) error
}

// Factory is the interface for creating new Client objects.
type Factory interface {
	NewClient(hvAPIKey string) Client
}

// HivelocityFactory implements the Factory interface.
type HivelocityFactory struct{}

var (
	// ErrDeviceNotFound gets returned if no matching device was found.
	ErrDeviceNotFound = fmt.Errorf("device was not found")

	// ErrDeviceShutDownAlready indicates that the device is shut down already.
	ErrDeviceShutDownAlready = fmt.Errorf("device is shut down already")

	// ErrDeviceTurnedOnAlready indicates that the device turned on already.
	ErrDeviceTurnedOnAlready = fmt.Errorf("device is turned on already")

	// ErrRateLimitExceeded indicates that the device turned on already.
	ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded")
)

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

func (c *realClient) GetDevice(ctx context.Context, deviceID int32) (hv.BareMetalDevice, error) {
	// https://developers.hivelocity.net/reference/get_bare_metal_device_id_resource
	device, _, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceIdResource(ctx, deviceID, nil) //nolint:bodyclose // Close() gets done in client
	if err == nil {
		return device, nil
	}
	var swaggerErr hv.GenericSwaggerError
	if errors.As(err, &swaggerErr) {
		if strings.HasPrefix(swaggerErr.Error(), fmt.Sprint(http.StatusNotFound)) {
			return device, ErrDeviceNotFound
		}
		log := log.FromContext(ctx)
		log.Info("GetDevice() failed", "body", string(swaggerErr.Body()))
	}
	return device, checkRateLimit(err)
}

func (c *realClient) SetDeviceTags(ctx context.Context, deviceID int32, tags []string) error {
	// https://developers.hivelocity.net/reference/put_device_tag_id_resource
	// Existing Tags will be removed by the HV API.
	deviceTags := hv.DeviceTag{
		Tags: tags,
	}
	_, _, err := c.client.DeviceApi.PutDeviceTagIdResource(ctx, deviceID, deviceTags, nil) //nolint:bodyclose // Close() gets done in client
	return checkRateLimit(err)
}

func (c *realClient) PowerOnDevice(_ context.Context, _ int32) error {
	return nil // todo
}

func (c *realClient) ProvisionDevice(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error) {
	log := log.FromContext(ctx)
	var swaggerErr hv.GenericSwaggerError

	power, _, err := c.client.DeviceApi.GetPowerResource(ctx, deviceID, nil) //nolint:bodyclose // Close() gets done in client
	if errors.As(err, &swaggerErr) {
		body := string(swaggerErr.Body())
		log.Info("ProvisionDevice() failed (GetPowerResource)", "DeviceID", deviceID, "body", body)
	}

	if power.PowerStatus == PowerStatusOn {
		// First we need to send "shutdown".
		// https://developers.hivelocity.net/reference/post_power_resource
		_, _, err := c.client.DeviceApi.PostPowerResource(ctx, deviceID, "shutdown", nil) //nolint:bodyclose // Close() gets done in client
		if errors.As(err, &swaggerErr) {
			body := string(swaggerErr.Body())
			log.Info("ProvisionDevice() failed (PostPowerResource)", "DeviceID", deviceID, "body", body)
		}
		log.Info("ProvisionDevice() called PostPowerResource shutdown", "DeviceID", deviceID)
		time.Sleep(30 * time.Second)
	}

	log.Info("calling ProvisionDevice()", "DeviceID", deviceID, "hostname", opts.Hostname, "OsName", opts.OsName,
		"script", utils.FirstN(opts.Script, 50),
		"ForceReload", opts.ForceReload)

	// https://developers.hivelocity.net/reference/put_bare_metal_device_id_resource
	localVars := hv.BareMetalDevicesApiPutBareMetalDeviceIdResourceOpts{
		SkipPowerCheck: optional.NewBool(true),
	}

	device, _, err := c.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(ctx, deviceID, opts, &localVars) //nolint:bodyclose // Close() gets done in client
	if errors.As(err, &swaggerErr) {
		body := string(swaggerErr.Body())
		log.Info("ProvisionDevice() failed (PutBareMetalDeviceIdResource)", "DeviceID", deviceID, "body", body)
		err = fmt.Errorf("%s: %w", body, swaggerErr)
	}
	if err == nil {
		log.Info("ProvisionDevice() was successful (PutBareMetalDeviceIdResource)", "DeviceID", deviceID)
	}
	return device, checkRateLimit(err)
}

func (c *realClient) ListDevices(ctx context.Context) ([]hv.BareMetalDevice, error) {
	devices, _, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceResource(ctx, nil) //nolint:bodyclose // Close() gets done in client
	return devices, checkRateLimit(err)
}

func (c *realClient) ShutdownDevice(ctx context.Context, deviceID int32) error {
	_, _, err := c.client.DeviceApi.PostPowerResource(ctx, deviceID, "shutdown", nil) //nolint:bodyclose // Close() gets done in client
	if err != nil {
		swaggerErr, ok := err.(hv.GenericSwaggerError)
		if ok {
			body := string(swaggerErr.Body())
			if strings.Contains(body, "Can't do this while server is powered off.") {
				return ErrDeviceShutDownAlready

			}
			err = fmt.Errorf("%s: %w", body, err)
		}
		return checkRateLimit(err)
	}
	return nil
}

func (c *realClient) ListImages(ctx context.Context, productID int32) ([]string, error) {
	// https://developers.hivelocity.net/reference/get_product_operating_systems_resource
	opts, _, err := c.client.ProductApi.GetProductOperatingSystemsResource(ctx, productID, nil) //nolint:bodyclose // Close() gets done in client
	ret := make([]string, 0, len(opts))
	if err != nil {
		return []string{}, checkRateLimit(err)
	}
	for i := range opts {
		ret = append(ret, opts[i].Name)
	}
	return ret, nil
}

func (c *realClient) ListSSHKeys(ctx context.Context) ([]hv.SshKeyResponse, error) {
	// https://developers.hivelocity.net/reference/get_ssh_key_resource
	sshKeys, _, err := c.client.SshKeyApi.GetSshKeyResource(ctx, nil) //nolint:bodyclose // Close() gets done in client
	return sshKeys, checkRateLimit(err)
}

// checkRateLimit returns true, if the Hivelocity rate limit was reached.
func checkRateLimit(err error) error {
	if err == nil {
		return nil
	}

	var swaggerErr hv.GenericSwaggerError
	if !errors.As(err, &swaggerErr) {
		return err
	}

	if strings.HasPrefix(swaggerErr.Error(), fmt.Sprint(http.StatusTooManyRequests)) {
		return ErrRateLimitExceeded
	}
	return err
}
