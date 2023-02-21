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
	"errors"
	"fmt"
	"strings"

	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/slices"
)

var errMultipleDevicesFound = fmt.Errorf(
	"multiple devices found while trying to find a single device",
)

// FindDeviceByTags returns the device with the given clusterTag and machineTag.
// Returns nil if no device was found.
func FindDeviceByTags(
	clusterTag string,
	machineTag string,
	devices []*hv.BareMetalDevice,
) (*hv.BareMetalDevice, error) {
	var device *hv.BareMetalDevice
	found := 0
	for i := range devices {
		if slices.Contains(devices[i].Tags, clusterTag) &&
			slices.Contains(devices[i].Tags, machineTag) {
			device = devices[i]
			found++
		}
	}
	if found > 1 {
		return nil, fmt.Errorf("found %v devices with tags %s and %s. Expected one: %w",
			found, clusterTag, machineTag, errMultipleDevicesFound)
	} else if found == 0 {
		return nil, nil
	}
	return device, nil
}

// FindUnusedDevice returns an unused device. Returns nil if no device was found.
func FindUnusedDevice(devices []*hv.BareMetalDevice, clusterName, deviceType string) (*hv.BareMetalDevice, error) {
	for i := range devices {
		device := devices[i]
		it, err := GetDeviceType(device)
		if err != nil {
			return nil, fmt.Errorf("[FindUnusedDevice] GetDeviceType() failed: %w", err)
		}
		if it != deviceType {
			continue
		}
		if DeviceHasTagKey(device, hvclient.TagKeyMachineName) {
			continue
		}
		cn, err := DeviceGetTagValue(device, hvclient.TagKeyClusterName)
		if errors.Is(err, ErrTooManyTagsFound) {
			continue
		}
		if errors.Is(err, ErrNoMatchingTagFound) {
			// this could lead to a race-condition, if two controllers of two clusters
			// try to fetch an unused device.
			// TODO: re-check after N seconds if there is a second tag from a second controller.
			return device, nil
		}
		if err != nil {
			return nil, err
		}
		if cn != clusterName {
			continue
		}
		return device, nil
	}
	return nil, nil
}

// DeviceHasTagKey returns true if the device has the tagKey set.
// Example: Your can check if a machine has already a name by using tagKey="machine-name".
func DeviceHasTagKey(device *hv.BareMetalDevice, tagKey string) bool {
	prefix := tagKey + "="
	for i := range device.Tags {
		if strings.HasPrefix(device.Tags[i], prefix) {
			return true
		}
	}
	return false
}

// ErrTooManyTagsFound gets returned, if there are multiple tags with the same key,
// and the key should be unique.
var ErrTooManyTagsFound = fmt.Errorf("too many tags found")

// ErrNoMatchingTagFound gets returned, if no matching tag was found.
var ErrNoMatchingTagFound = fmt.Errorf("no matching tag found")

// DeviceGetTagValue returns the value of a TagKey of a device.
// Example: If a device has the tag "foo=bar", then DeviceGetTagValue
// will return "bar".
// If there is no such tag, or there are two tags, then an error gets returned.
func DeviceGetTagValue(device *hv.BareMetalDevice, tagKey string) (string, error) {
	prefix := tagKey + "="
	found := 0
	value := ""
	for i := range device.Tags {
		if !strings.HasPrefix(device.Tags[i], prefix) {
			continue
		}
		if found > 0 {
			return "", fmt.Errorf("[DeviceGetTagValue] device %q, tagKey %q: %w",
				device.Hostname, tagKey, ErrTooManyTagsFound)
		}
		found++
		value = device.Tags[i][len(prefix):]
	}
	if found == 0 {
		return "", fmt.Errorf("[DeviceGetTagValue] device %q, tagKey %q: %w",
			device.Hostname, tagKey, ErrNoMatchingTagFound)
	}
	return value, nil
}

// GetDeviceType returns the device-type of this BareMetalDevice.
func GetDeviceType(device *hv.BareMetalDevice) (string, error) {
	deviceType, err := DeviceGetTagValue(device, hvclient.TagKeyDeviceType)
	if err != nil {
		return "", fmt.Errorf("[GetDeviceType] DeviceGetTagValue() failed: %w", err)
	}
	return deviceType, nil
}
