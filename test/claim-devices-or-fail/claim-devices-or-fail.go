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

// Package main contains functions to test the Hivelocity API.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/slices"
)

func main() {
	apiKey := os.Getenv("HIVELOCITY_API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing environment variable HIVELOCITY_API_KEY")
	}
	ctx := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: apiKey,
	})
	if len(os.Args) < 2 {
		log.Fatalln("please provide one more device types (like hvControlPlane)")
	}
	apiClient := hv.NewAPIClient(hv.NewConfiguration())
	allDevices, _, err := apiClient.BareMetalDevicesApi.GetBareMetalDeviceResource(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 1; i < len(os.Args); i++ {
		deviceType := os.Args[i]
		err := releaseOldMachines(ctx, apiClient, deviceType, allDevices)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func releaseOldMachines(ctx context.Context, apiClient *hv.APIClient, deviceType string,
	allDevices []hv.BareMetalDevice) error {
	devicesWithTag := make([]hv.BareMetalDevice, 0)
	tag := fmt.Sprintf("%s=%s", hvtag.DeviceTagKeyDeviceType, deviceType)

	for _, device := range allDevices {
		if !slices.Contains(device.Tags, tag) {
			continue
		}
		skip := false
		for _, tag := range device.Tags {
			if strings.HasPrefix(tag, string(hvtag.DeviceTagKeyPermanentError)) {
				fmt.Printf("skipping device %d because %s\n", device.DeviceId, tag)
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		devicesWithTag = append(devicesWithTag, device)
	}
	if len(devicesWithTag) == 0 {
		return fmt.Errorf("no device found with %s=%s", hvtag.DeviceTagKeyDeviceType, deviceType)
	}
	fmt.Printf("resetting labels of all devices which have %s=%s. Found %d devices\n",
		hvtag.DeviceTagKeyDeviceType, deviceType, len(devicesWithTag))
	for _, device := range devicesWithTag {
		err := resetTags(ctx, device, apiClient)
		if err != nil {
			return err
		}
	}

	return nil
}

// resetTags: Remove our labels, but keep DeviceTagKeyPermanentError and DeviceTagKeyDeviceType.
// And keep other labels.
func resetTags(ctx context.Context, device hv.BareMetalDevice, apiClient *hv.APIClient) error {
	fmt.Printf("    resetting labels of device %d\n", device.DeviceId)
	var newTags []string
	for _, tag := range device.Tags {
		if removeTag(tag) {
			continue
		}
		newTags = append(newTags, tag)
	}
	_, _, err := apiClient.DeviceApi.PutDeviceTagIdResource(ctx, device.DeviceId, hv.DeviceTag{
		Tags: newTags}, nil)
	if err != nil {
		return err
	}
	return nil
}

// tag: Something like caphv-cluster-name=hv-guettli
func removeTag(tag string) bool {
	if !strings.HasPrefix(tag, "caphv-") {
		return false
	}
	for _, keepPrefix := range []string{
		string(hvtag.DeviceTagKeyPermanentError),
		string(hvtag.DeviceTagKeyDeviceType),
	} {
		if strings.HasPrefix(tag, keepPrefix) {
			return false
		}
	}
	return true
}
