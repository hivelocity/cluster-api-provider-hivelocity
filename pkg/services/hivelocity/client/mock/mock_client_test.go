/*
Copyright 2022 The Kubernetes Authors.

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

package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SetTags(t *testing.T) {
	client := NewMockedHVClientFactory().NewClient("dummy-key")
	ctx := context.Background()
	err := client.SetTags(ctx, FreeDeviceID, []string{"tag1", "tag2"})
	require.NoError(t, err)
	device, err := client.GetDevice(ctx, FreeDeviceID)
	require.NoError(t, err)
	require.ElementsMatch(t, device.Tags, []string{"tag1", "tag2"})
}

func Test_GetDevice(t *testing.T) {
	client := NewMockedHVClientFactory().NewClient("dummy-key")
	ctx := context.Background()
	device, err := client.GetDevice(ctx, FreeDeviceID)
	require.NoError(t, err)
	require.Equal(t, int32(FreeDeviceID), device.DeviceId)

	device, err = client.GetDevice(ctx, -1)
	require.Error(t, err)
}
func Test_NewMockedHVClientFactory(t *testing.T) {
	factory := NewMockedHVClientFactory()
	client := factory.NewClient("dummy-key")
	ctx := context.Background()
	device, err := client.GetDevice(ctx, FreeDeviceID)
	require.NoError(t, err)
	require.ElementsMatch(t, device.Tags, []string{"caphv-device-type=hvCustom"})
	err = client.SetTags(ctx, FreeDeviceID, []string{"new-tag"})
	require.NoError(t, err)

	device, err = client.GetDevice(ctx, FreeDeviceID)
	require.NoError(t, err)
	require.ElementsMatch(t, device.Tags, []string{"new-tag"})

	client2 := factory.NewClient("dummy-key")
	device, err = client2.GetDevice(ctx, FreeDeviceID)
	require.NoError(t, err)
	require.ElementsMatch(t, device.Tags, []string{"new-tag"})

	factoryNewF := NewMockedHVClientFactory()
	clientNewF := factoryNewF.NewClient("dummy-key")
	device, err = clientNewF.GetDevice(ctx, FreeDeviceID)
	require.NoError(t, err)
	require.ElementsMatch(t, device.Tags, []string{"caphv-device-type=hvCustom"})
}
