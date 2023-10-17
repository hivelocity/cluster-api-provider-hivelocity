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

// Package hvtag contains objects and utility functions to handle tags of Hivelocity devices.
package hvtag

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// DeviceTagKey defines the key of the key-value pair that is stored as tag of Hivelocity devices.
type DeviceTagKey string

const (
	// DeviceTagKeyMachine is the key for the name of the associated HivelocityMachine object.
	DeviceTagKeyMachine DeviceTagKey = "caphv-machine-name"

	// DeviceTagKeyCluster is the key for the name of the associated HivelocityCluster object.
	DeviceTagKeyCluster DeviceTagKey = "caphv-cluster-name"

	// DeviceTagKeyDeviceType is the key for the device type that users can reference in HivelocityMachine.Spec.
	DeviceTagKeyDeviceType DeviceTagKey = "caphv-device-type"

	// DeviceTagKeyMachineType is the key for the machine type, i.e. worker, control_plane.
	DeviceTagKeyMachineType DeviceTagKey = "caphv-machine-type"

	// DeviceTagKeyPermanentError is the key for machines which need a manual reset by a Hivelocity admin.
	DeviceTagKeyPermanentError DeviceTagKey = "caphv-permanent-error"
)

// Prefix returns the prefix based on this DeviceTagKey used in Hivelocity tag strings.
func (key DeviceTagKey) Prefix() string {
	return fmt.Sprintf("%s=", key)
}

var (
	// ErrDeviceTagKeyInvalid indicates that the device tag key is invalid.
	ErrDeviceTagKeyInvalid = fmt.Errorf("invalid device tag key")
	// ErrDeviceTagInvalidFormat indicates that the device tag has an invalid format.
	ErrDeviceTagInvalidFormat = fmt.Errorf("invalid format of device tag")
	// ErrDeviceTagEmptyKey indicates that the device tag has an empty key.
	ErrDeviceTagEmptyKey = fmt.Errorf("device tag has empty key")
	// ErrDeviceTagEmptyValue indicates that the device tag has an empty value.
	ErrDeviceTagEmptyValue = fmt.Errorf("device tag has empty value")
	// ErrMultipleDeviceTagsFound indicates that multiple device tags have been found in a list.
	ErrMultipleDeviceTagsFound = fmt.Errorf("found multiple device tags")
	// ErrDeviceTagNotFound indicates that no device tag has been found.
	ErrDeviceTagNotFound = fmt.Errorf("no device tag found")
)

// IsValid checks whether a DeviceTagKey is valid.
func (key DeviceTagKey) IsValid() bool {
	return key == DeviceTagKeyMachine ||
		key == DeviceTagKeyCluster ||
		key == DeviceTagKeyDeviceType ||
		key == DeviceTagKeyMachineType
}

// DeviceTag defines the object that represents a key-value pair that is stored as tag of Hivelocity devices.
type DeviceTag struct {
	Key   DeviceTagKey
	Value string
}

// DeviceTagFromList takes the tag of a HV device and returns a DeviceTag or an error if it is invalid.
func DeviceTagFromList(key DeviceTagKey, tagList []string) (DeviceTag, error) {
	var found bool
	var deviceTag DeviceTag
	var err error

	for _, tagString := range tagList {
		// filter out irrelevant tagStrings quickly
		if !strings.HasPrefix(tagString, key.Prefix()) {
			continue
		}

		// get DeviceTag from tagString
		deviceTag, err = deviceTagFromString(tagString)
		if err != nil {
			continue
		}

		// additional check if key is correct - probably not necessary due to HasPrefix check
		if deviceTag.Key != key {
			continue
		}

		// Check whether a correct DeviceTag has been found already. If so, return with error.
		if found {
			return DeviceTag{}, ErrMultipleDeviceTagsFound
		}
		found = true
	}

	if !found {
		return DeviceTag{}, ErrDeviceTagNotFound
	}

	return deviceTag, nil
}

// MachineTagFromList returns the machine tag from a list of tag strings.
func MachineTagFromList(tagList []string) (DeviceTag, error) {
	return DeviceTagFromList(DeviceTagKeyMachine, tagList)
}

// ClusterTagFromList returns the cluster tag from a list of tag strings.
func ClusterTagFromList(tagList []string) (DeviceTag, error) {
	return DeviceTagFromList(DeviceTagKeyCluster, tagList)
}

// DeviceTypeTagFromList returns the device type tag from a list of tag strings.
func DeviceTypeTagFromList(tagList []string) (DeviceTag, error) {
	return DeviceTagFromList(DeviceTagKeyDeviceType, tagList)
}

// PermanentErrorTagFromList returns the permanent error tag from a list of tag strings.
func PermanentErrorTagFromList(tagList []string) (DeviceTag, error) {
	return DeviceTagFromList(DeviceTagKeyPermanentError, tagList)
}

// deviceTagFromString takes the tag of a HV device and returns a DeviceTag or an error if it is invalid.
func deviceTagFromString(tagString string) (DeviceTag, error) {
	tagElements := strings.Split(tagString, "=")
	if len(tagElements) != 2 {
		return DeviceTag{}, ErrDeviceTagInvalidFormat
	}

	key := DeviceTagKey(tagElements[0])
	value := tagElements[1]

	if key == "" {
		return DeviceTag{}, ErrDeviceTagEmptyKey
	}
	if value == "" {
		return DeviceTag{}, ErrDeviceTagEmptyValue
	}

	if !key.IsValid() {
		return DeviceTag{}, ErrDeviceTagKeyInvalid
	}

	return DeviceTag{key, value}, nil
}

// ToString returns the string representation of a DeviceTag object.
func (deviceTag DeviceTag) ToString() string {
	return string(deviceTag.Key) + "=" + deviceTag.Value
}

// IsInStringList checks whether a DeviceTag object can be found in a list of tag strings.
func (deviceTag DeviceTag) IsInStringList(tagList []string) bool {
	return slices.Contains(tagList, deviceTag.ToString())
}

// RemoveFromList removes all tag strings from a list that equal the string representation of DeviceTag.
func (deviceTag DeviceTag) RemoveFromList(tagList []string) (newTagList []string, updated bool) {
	newTagList = make([]string, 0, len(tagList))
	for _, tagString := range tagList {
		// append all tag strings to newTagList which do not equal this device tag
		if tagString == deviceTag.ToString() {
			updated = true
		} else {
			newTagList = append(newTagList, tagString)
		}
	}
	return newTagList, updated
}
