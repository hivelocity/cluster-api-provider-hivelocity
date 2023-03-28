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

// Package device implements functions to manage the lifecycle of Hivelocity devices.
package device

import (
	"testing"

	mockclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client/mock"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"github.com/stretchr/testify/require"
)

func Test_chooseAvailableFromList(t *testing.T) {
	devices := []hv.BareMetalDevice{
		mockclient.NoTagsDevice,
		mockclient.FreeDevice,
	}
	_, err := chooseAvailableFromList(devices, "fooDeviceType", "my-cluster", "my-machine")
	require.ErrorIs(t, err, errNoDeviceAvailable)
}
