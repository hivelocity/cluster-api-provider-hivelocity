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
	"fmt"
	"strconv"
	"strings"

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

// ProviderIDToDeviceID converts the ProviderID (hivelocity://NNNN) to the DeviceID.
func ProviderIDToDeviceID(providerID string) (int32, error) {
	if providerID == "" {
		return 0, fmt.Errorf("[ProviderIDToDeviceID] providerID is empty")
	}
	prefix := "hivelocity://"
	if !strings.HasPrefix(providerID, prefix) {
		return 0, fmt.Errorf("[ProviderIDToDeviceID] missing prefix %q in providerID %q",
			prefix, providerID)
	}
	deviceID, err := strconv.ParseInt(
		strings.TrimPrefix(providerID, prefix),
		10,
		32,
	)
	if err != nil {
		return 0, fmt.Errorf("[ProviderIDToDeviceID] failed to convert providerID %q: %w",
			providerID, err)
	}
	return int32(deviceID), nil
}

// DeviceIDToProviderID converts a deviceID to ProviderID.
func DeviceIDToProviderID(deviceID int32) string {
	return fmt.Sprintf("hivelocity://%d", deviceID)
}
