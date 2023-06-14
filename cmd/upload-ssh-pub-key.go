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
	"context"
	"fmt"
	"os"

	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/spf13/cobra"
)

var uploadSshKey = &cobra.Command{
	Use:   "upload-ssh-pub-key [flags] SSH_PUB_KEY_NAME SSH_PUB_KEY_FILE",
	Short: "Uploads a ssh pub-key to Hivelocity",
	Run:   runUploadSshKey,
	Args:  cobra.ExactArgs(2),
}

func runUploadSshKey(cmd *cobra.Command, args []string) {
	sshPubKeyName := args[0]
	sshPubKeyContent, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	apiKey := os.Getenv("HIVELOCITY_API_KEY")
	if apiKey == "" {
		fmt.Println("Missing environment variable HIVELOCITY_API_KEY")
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
			fmt.Printf("Key %q already exists. ID: %d\n", sshPubKeyName, sshKey.SshKeyId)
			return
		}
	}

	// key does not exist yet
	ssKey, _, err := apiClient.SshKeyApi.PostSshKeyResource(ctx, hv.SshKey{
		Name:      sshPubKeyName,
		PublicKey: string(sshPubKeyContent),
	}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Key %q created. ID: %d\n", sshPubKeyName, ssKey.SshKeyId)
}
