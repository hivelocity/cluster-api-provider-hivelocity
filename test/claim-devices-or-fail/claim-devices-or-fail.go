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

	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/slices"
)

type controlPlaneDeviceID struct {
	ip string
}

var controlPlaneIDs = []int32{
	14730, // 66.165.243.74
	15335, // 66.206.8.178
	15336, // 66.206.8.186
}

var workerNodeIDs = []int32{
	15337, // 66.206.8.194
	// 15338 is used for test/hvapi/main.go
}

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
		releaseOldMachines(ctx, apiClient, deviceType, allDevices)
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
		devicesWithTag = append(devicesWithTag, device)
	}
	fmt.Printf("resetting all devices which have %s=%s. Found %d devices\n",
		hvtag.DeviceTagKeyDeviceType, deviceType, len(devicesWithTag))
	for _, device := range devicesWithTag {
		fmt.Printf("    resetting device %d\n", device.DeviceId)
		_, _, err := apiClient.DeviceApi.PutDeviceTagIdResource(ctx, device.DeviceId, hv.DeviceTag{
			Tags: []string{tag}}, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
