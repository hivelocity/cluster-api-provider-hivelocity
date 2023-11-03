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

package labels

import (
	"strings"
)

// Tags is a slice of string. It implements https://pkg.go.dev/k8s.io/apimachinery/pkg/labels#Labels interface.
type Tags []string

// Has returns whether the provided label exists in the slice.
// It implements the Has function of the Labels interface.
func (ls Tags) Has(label string) bool {
	for _, tag := range ls {
		if !strings.HasPrefix(tag, "caphvlabel:") {
			continue
		}

		trimmedTag := strings.TrimPrefix(tag, "caphvlabel:")
		splittedTrimTag := strings.Split(trimmedTag, "=")

		if len(splittedTrimTag) == 2 && splittedTrimTag[0] == label {
			return true
		}
	}

	return false
}

// Get returns the value in the slice for the provided label.
// It implements the Has function of the Labels interface.
func (ls Tags) Get(label string) string {
	for _, tag := range ls {
		if !strings.HasPrefix(tag, "caphvlabel:") {
			continue
		}

		trimmedTag := strings.TrimPrefix(tag, "caphvlabel:")
		splittedTrimTag := strings.Split(trimmedTag, "=")

		if len(splittedTrimTag) == 2 && splittedTrimTag[0] == label {
			return splittedTrimTag[1]
		}
	}

	return ""
}
