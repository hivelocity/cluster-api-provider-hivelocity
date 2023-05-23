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
	"errors"
	"fmt"
	"log"
	"os"

	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
)

const manualDeviceID = 15335

func main() {
	manualTests()
}

func cliTestListDevices(ctx context.Context, client hvclient.Client) {
	allDevices, err := client.ListDevices(ctx)
	if err != nil {
		panic(err.Error())
	}
	for _, device := range allDevices {
		fmt.Printf("device: %+v\n", device)
	}
}

func cliTestSetTags(ctx context.Context, client hvclient.Client) {
	err := client.SetDeviceTags(ctx, manualDeviceID, []string{"foo"})
	if err != nil {
		panic(err)
	}
}

func cliTestGetDevice(ctx context.Context, client hvclient.Client) {
	device, err := client.GetDevice(ctx, manualDeviceID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("device: %+v\n", device)
}

func cliTestListImages(ctx context.Context, client hvclient.Client) {
	images, err := client.ListImages(ctx, 504)
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Println(image)
	}
}

func cliTestListSSHKeys(ctx context.Context, client hvclient.Client) {
	keys, err := client.ListSSHKeys(ctx)
	if err != nil {
		panic(err)
	}
	for _, sshKey := range keys {
		fmt.Println(sshKey)
	}
}

func cliTestProvisionDevice(ctx context.Context, client hvclient.Client) {
	script := `#cloud-config
write_files:
- content: |
		a & b && c <foo>
	path: /opt/test.txt
	`
	opts := hv.BareMetalDeviceUpdate{
		Hostname:       "my-host-name.example.com",
		OsName:         "Ubuntu 20.x",
		PublicSshKeyId: 918,
		Script:         script,
	}
	device, err := client.ProvisionDevice(ctx, manualDeviceID, opts)
	if err != nil {
		var swaggerErr hv.GenericSwaggerError
		if errors.As(err, &swaggerErr) {
			fmt.Println(string(swaggerErr.Body()))
		}
		log.Fatalln(err)
	}
	fmt.Printf("device: %+v\n", device)
}

func manualTests() {
	factory := hvclient.HivelocityFactory{}
	client := factory.NewClient(os.Getenv("HIVELOCITY_API_KEY"))
	ctx := context.Background()
	if len(os.Args) == 1 {
		log.Fatalln("see code for possible options")
	}
	switch arg := os.Args[1]; arg {
	case "ListDevices":
		cliTestListDevices(ctx, client)
		os.Exit(0)
	case "ListSSHKeys":
		cliTestListSSHKeys(ctx, client)
		os.Exit(0)
	case "ListImages":
		cliTestListImages(ctx, client)
		os.Exit(0)
	case "SetTags":
		cliTestSetTags(ctx, client)
		os.Exit(0)
	case "GetDevice":
		cliTestGetDevice(ctx, client)
		os.Exit(0)
	case "ProvisionDevice":
		cliTestProvisionDevice(ctx, client)
		os.Exit(0)
	}
}
