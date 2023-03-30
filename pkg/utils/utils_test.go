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

package utils_test

import (
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringInList", func() {
	DescribeTable("Test string in list",
		func(list []string, str string, expectedOutcome bool) {
			out := utils.StringInList(list, str)
			Expect(out).To(Equal(expectedOutcome))
		},
		Entry("entry1", []string{"a", "b", "c"}, "a", true),
		Entry("entry2", []string{"a", "b", "c"}, "d", false))
})

var _ = Describe("FilterStringFromList", func() {
	DescribeTable("Test filter string from list",
		func(list []string, str string, expectedOutcome []string) {
			out := utils.FilterStringFromList(list, str)
			Expect(out).To(Equal(expectedOutcome))
		},
		Entry("entry1", []string{"a", "b", "c"}, "a", []string{"b", "c"}),
		Entry("entry2", []string{"a", "b", "c"}, "d", []string{"a", "b", "c"}))
})
