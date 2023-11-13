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
		log.Fatalln("please provide one more device types (like caphvlabel:deviceType=hvControlPlane)")
	}
	apiClient := hv.NewAPIClient(hv.NewConfiguration())
	allDevices, _, err := apiClient.BareMetalDevicesApi.GetBareMetalDeviceResource(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	var done []string
	for i := 1; i < len(os.Args); i++ {
		tag := os.Args[i]
		if slices.Contains(done, tag) {
			continue
		}
		done = append(done, tag)
		if !strings.HasPrefix(tag, "caphvlabel:deviceType=") {
			log.Fatalln("tag must start with caphvlabel:deviceType=")
		}
		err := releaseOldMachines(ctx, apiClient, tag, allDevices)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func releaseOldMachines(ctx context.Context, apiClient *hv.APIClient, tag string,
	allDevices []hv.BareMetalDevice,
) error {
	devicesWithTag := make([]hv.BareMetalDevice, 0)

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
		return fmt.Errorf("no device found with %q, Please make a corresponding device available. For example, by giving a machine the appropriate label via the HV web UI", tag)
	}
	fmt.Printf("resetting labels of all devices which have %s. Found %d devices\n", tag, len(devicesWithTag))
	for _, device := range devicesWithTag {
		fmt.Printf("    resetting labels of device %d\n", device.DeviceId)
		newTags := hvtag.RemoveEphemeralTags(device.Tags)
		_, _, err := apiClient.DeviceApi.PutDeviceTagIdResource(ctx, device.DeviceId, hv.DeviceTag{
			Tags: newTags,
		}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
