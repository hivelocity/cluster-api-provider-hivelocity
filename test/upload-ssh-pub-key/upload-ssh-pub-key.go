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

	hv "github.com/hivelocity/hivelocity-client-go/client"
)

const sshPubKeyName = "ssh-key-hivelocity-pub"
const sshPubKeyContent = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXSatqG8kdtmo2n1VxvHUxq6HkkZ2mazjuLIcPPB5JVtyfe/EBi89mHr0/kYeY0AELk7pqiRshnvvJ9sIm7HC1xqIphE4+/XPNt/fwnFBxdU3WfPf1ho/L9xHaPcojrsC+NvLFcJKr/351EGgeYneUfb4FY7ElAKQq8oxRSMpODfKhg9BRQCDYEvTblLHR8lTUB2Vx7D33TlMThHqTSWuw2aXj8s53XeVFPBhR0KYRt7731M79oo4qPLlVMRwPUH2H2RTdoEHG50x/g+0MXYU0INO2sQMXizWt8rwV8hyiMg41+hteYkAGFOpRjZWD7ez8yMAQ2O3KcJFSgvGOs042nUmqFLB0AZFvSvptIqLBdYHrcgsQgQYVKVC1f+iLjexUQHABbp6liHJ0kSXS3twamgx+WtNsdaEvykUecLHZpzIMBLeCYXOy4S33L7ywnxWO+KOqnF8MZTpQoJP1HBZNJPBExX788UWiBlb/jKTDAksOfR43PNEHyPx8sQFPGf8= hivelocity"

func main() {
	apiKey := os.Getenv("HIVELOCITY_API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing environment variable HIVELOCITY_API_KEY")
		os.Exit(1)
	}
	ctx := context.WithValue(context.Background(), hv.ContextAPIKey, hv.APIKey{
		Key: apiKey,
	})
	apiClient := hv.NewAPIClient(hv.NewConfiguration())
	sshKeyResponses, _, err := apiClient.SshKeyApi.GetSshKeyResource(ctx, nil)
	if err != nil {
		panic(err)
	}
	for _, sshKey := range sshKeyResponses {
		if sshKey.Name == sshPubKeyName {
			fmt.Printf("Key %q already exists\n", sshPubKeyName)
			return
		}
	}

	// key does not exist yet
	_, _, err = apiClient.SshKeyApi.PostSshKeyResource(ctx, hv.SshKey{
		Name:      sshPubKeyName,
		PublicKey: sshPubKeyContent,
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Key %q created\n", sshPubKeyName)
}
