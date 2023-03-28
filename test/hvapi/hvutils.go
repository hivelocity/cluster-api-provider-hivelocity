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
	"fmt"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/slices"
)

var errMultipleDevicesFound = fmt.Errorf(
	"multiple devices found while trying to find a single device",
)

// findDeviceByTags returns the device with the given clusterTag and machineTag.
// Returns nil if no device was found.
func findDeviceByTags(
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
