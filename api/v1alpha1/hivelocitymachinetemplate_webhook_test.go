/*
Copyright 2024 The Kubernetes Authors.

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

package v1alpha1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/selection"
)

func TestHivelocityMachineTemplateWebhook_ValidateCreate_valid(t *testing.T) {
	ctx := context.Background()
	hook := &HivelocityMachineTemplateWebhook{}
	hmt := HivelocityMachineTemplate{}
	for _, ds := range []DeviceSelector{
		{},
		{
			MatchLabels: map[string]string{"key": "value"},
			MatchExpressions: []DeviceSelectorRequirement{
				{
					Key:      "foo",
					Operator: selection.In,
					Values:   []string{"one", "two"},
				},
			},
		},
	} {
		hmt.Spec.Template.Spec.DeviceSelector = ds
		warnings, err := hook.ValidateCreate(ctx, &hmt)
		require.Nil(t, err)
		require.Len(t, warnings, 0)
	}
}

func TestHivelocityMachineTemplateWebhook_ValidateCreate_invalid(t *testing.T) {
	ctx := context.Background()
	hook := &HivelocityMachineTemplateWebhook{}
	hmt := HivelocityMachineTemplate{}
	for _, ds := range []DeviceSelector{
		{
			MatchLabels:      map[string]string{"key:invalid": "value"},
			MatchExpressions: []DeviceSelectorRequirement{},
		},
		{
			MatchLabels:      map[string]string{"key": "value:invalid"},
			MatchExpressions: []DeviceSelectorRequirement{},
		},
		{
			MatchLabels: map[string]string{"key": "value"},
			MatchExpressions: []DeviceSelectorRequirement{
				{
					Key:      "foo",
					Operator: selection.In,
					Values:   []string{"one:invalid", "two"},
				},
			},
		},
	} {
		hmt.Spec.Template.Spec.DeviceSelector = ds
		warnings, err := hook.ValidateCreate(ctx, &hmt)
		require.NotNil(t, err)
		require.Len(t, warnings, 0)

	}
}
