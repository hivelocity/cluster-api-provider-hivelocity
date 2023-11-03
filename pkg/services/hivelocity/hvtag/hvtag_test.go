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

package hvtag

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHVTag(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HVTag Suite")
}

var _ = Describe("Test DeviceTagKey.Prefix", func() {
	type testCaseDeviceTagKeyPrefix struct {
		key          DeviceTagKey
		expectPrefix string
	}

	DescribeTable("Test DeviceTagKey.Prefix",
		func(tc testCaseDeviceTagKeyPrefix) {
			Expect(tc.key.Prefix()).Should(Equal(tc.expectPrefix))
		},
		Entry("cluster key", testCaseDeviceTagKeyPrefix{
			key:          DeviceTagKeyCluster,
			expectPrefix: "caphv-cluster-name=",
		}),
		Entry("machine key", testCaseDeviceTagKeyPrefix{
			key:          DeviceTagKeyMachine,
			expectPrefix: "caphv-machine-name=",
		}),
	)
})

var _ = Describe("Test DeviceTagKey.IsValid", func() {
	type testCaseDeviceTagKeyIsValid struct {
		key           DeviceTagKey
		expectIsValid bool
	}

	DescribeTable("Test DeviceTagKey.IsValid",
		func(tc testCaseDeviceTagKeyIsValid) {
			Expect(tc.key.IsValid()).Should(Equal(tc.expectIsValid))
		},
		Entry("cluster key", testCaseDeviceTagKeyIsValid{
			key:           DeviceTagKeyCluster,
			expectIsValid: true,
		}),
		Entry("machine key", testCaseDeviceTagKeyIsValid{
			key:           DeviceTagKeyMachine,
			expectIsValid: true,
		}),
		Entry("machine type key", testCaseDeviceTagKeyIsValid{
			key:           DeviceTagKeyMachineType,
			expectIsValid: true,
		}),
		Entry("other key", testCaseDeviceTagKeyIsValid{
			key:           "caphv-other",
			expectIsValid: false,
		}),
	)
})

var _ = Describe("TestDeviceTagFromList", func() {
	type testCaseDeviceTagFromList struct {
		key             DeviceTagKey
		tagList         []string
		expectDeviceTag DeviceTag
		expectError     error
	}

	DescribeTable("TestDeviceTagFromList",
		func(tc testCaseDeviceTagFromList) {
			deviceTag, err := DeviceTagFromList(tc.key, tc.tagList)
			if tc.expectError != nil {
				Expect(err).ToNot(BeNil())
				Expect(err).Should(Equal(tc.expectError))
			}
			Expect(deviceTag).Should(Equal(tc.expectDeviceTag))
		},
		Entry("valid and unique device tag - only one tag in list", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyCluster.Prefix() + "value"},
			expectDeviceTag: DeviceTag{DeviceTagKeyCluster, "value"},
			expectError:     nil,
		}),
		Entry("valid and unique device tag - multiple tags in list", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyCluster.Prefix() + "value", DeviceTagKeyMachine.Prefix() + "othervalue", "othertag"},
			expectDeviceTag: DeviceTag{DeviceTagKeyCluster, "value"},
			expectError:     nil,
		}),
		Entry("valid but not unique device tag", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyCluster.Prefix() + "value", DeviceTagKeyCluster.Prefix() + "othervalue"},
			expectDeviceTag: DeviceTag{},
			expectError:     ErrMultipleDeviceTagsFound,
		}),
		Entry("valid device tag existing twice", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyCluster.Prefix() + "value", DeviceTagKeyCluster.Prefix() + "value"},
			expectDeviceTag: DeviceTag{},
			expectError:     ErrMultipleDeviceTagsFound,
		}),
		Entry("valid but not unique device tag - key exists twice", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyCluster.Prefix() + "value", string(DeviceTagKeyCluster) + "value"},
			expectDeviceTag: DeviceTag{DeviceTagKeyCluster, "value"},
			expectError:     nil,
		}),
		Entry("device tag does not exist", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyMachine.Prefix() + "value", string(DeviceTagKeyCluster) + "value"},
			expectDeviceTag: DeviceTag{},
			expectError:     ErrDeviceTagNotFound,
		}),
		Entry("device tag has empty value", testCaseDeviceTagFromList{
			key:             DeviceTagKeyCluster,
			tagList:         []string{DeviceTagKeyCluster.Prefix(), DeviceTagKeyMachine.Prefix() + "othervalue"},
			expectDeviceTag: DeviceTag{},
			expectError:     ErrDeviceTagNotFound,
		}),
	)
})

var _ = Describe("TestMachineTagFromList", func() {
	It("returns a DeviceTag with DeviceTagKeyMachine", func() {
		deviceTag, err := MachineTagFromList([]string{DeviceTagKeyMachine.Prefix() + "value", string(DeviceTagKeyCluster) + "value"})
		Expect(err).To(BeNil())
		Expect(deviceTag).To(Equal(DeviceTag{Key: DeviceTagKeyMachine, Value: "value"}))
	})
})

var _ = Describe("TestClusterTagFromList", func() {
	It("returns a DeviceTag with DeviceTagKeyCluster", func() {
		deviceTag, err := ClusterTagFromList([]string{DeviceTagKeyCluster.Prefix() + "value", string(DeviceTagKeyCluster) + "value"})
		Expect(err).To(BeNil())
		Expect(deviceTag).To(Equal(DeviceTag{Key: DeviceTagKeyCluster, Value: "value"}))
	})
})

var _ = Describe("deviceTagFromString", func() {
	type testCaseDeviceTagFromString struct {
		tagString       string
		expectDeviceTag DeviceTag
		expectError     error
	}

	DescribeTable("deviceTagFromString",
		func(tc testCaseDeviceTagFromString) {
			deviceTag, err := deviceTagFromString(tc.tagString)
			if tc.expectError != nil {
				Expect(err).ToNot(BeNil())
				Expect(err).Should(Equal(tc.expectError))
			}
			Expect(deviceTag).Should(Equal(tc.expectDeviceTag))
		},
		Entry("valid device tag - cluster", testCaseDeviceTagFromString{
			tagString:       "caphv-cluster-name=mycluster",
			expectDeviceTag: DeviceTag{DeviceTagKeyCluster, "mycluster"},
			expectError:     nil,
		}),
		Entry("valid device tag - machine", testCaseDeviceTagFromString{
			tagString:       "caphv-machine-name=mymachine",
			expectDeviceTag: DeviceTag{DeviceTagKeyMachine, "mymachine"},
			expectError:     nil,
		}),
		Entry("device tag key does not exist", testCaseDeviceTagFromString{
			tagString:       "wrongkey=value",
			expectDeviceTag: DeviceTag{},
			expectError:     ErrDeviceTagKeyInvalid,
		}),
		Entry("device tag has empty value", testCaseDeviceTagFromString{
			tagString:       "caphv-cluster-name=",
			expectDeviceTag: DeviceTag{},
			expectError:     ErrDeviceTagEmptyValue,
		}),
		Entry("device tag has empty key", testCaseDeviceTagFromString{
			tagString:       "=value",
			expectDeviceTag: DeviceTag{},
			expectError:     ErrDeviceTagEmptyKey,
		}),
		Entry("device tag has invalid format", testCaseDeviceTagFromString{
			tagString:       "key1=key2=value",
			expectDeviceTag: DeviceTag{},
			expectError:     ErrDeviceTagInvalidFormat,
		}),
	)
})

var _ = Describe("Test DeviceTag.ToString", func() {
	type testCaseDeviceTagToString struct {
		deviceTag    DeviceTag
		expectString string
	}

	DescribeTable("Test DeviceTag.ToString",
		func(tc testCaseDeviceTagToString) {
			Expect(tc.deviceTag.ToString()).Should(Equal(tc.expectString))
		},
		Entry("cluster key", testCaseDeviceTagToString{
			deviceTag:    DeviceTag{DeviceTagKeyCluster, "mycluster"},
			expectString: "caphv-cluster-name=mycluster",
		}),
		Entry("machine key", testCaseDeviceTagToString{
			deviceTag:    DeviceTag{DeviceTagKeyMachine, "mymachine"},
			expectString: "caphv-machine-name=mymachine",
		}),
	)
})

var _ = Describe("Test DeviceTag.IsInStringList", func() {
	type testCaseDeviceTagIsInStringList struct {
		deviceTag  DeviceTag
		tagList    []string
		expectBool bool
	}

	DescribeTable("Test DeviceTag.IsInStringList",
		func(tc testCaseDeviceTagIsInStringList) {
			Expect(tc.deviceTag.IsInStringList(tc.tagList)).Should(Equal(tc.expectBool))
		},
		Entry("deviceTag exists among multiple", testCaseDeviceTagIsInStringList{
			deviceTag:  DeviceTag{DeviceTagKeyCluster, "mycluster"},
			tagList:    []string{"caphv-cluster-name=mycluster", "othertag", "caphv-machine-name=mymachine"},
			expectBool: true,
		}),
		Entry("deviceTag does not exists among multiple", testCaseDeviceTagIsInStringList{
			deviceTag:  DeviceTag{DeviceTagKeyCluster, "mycluster"},
			tagList:    []string{"caphv-cluster-type=mytype", "othertag", "caphv-machine-name=mymachine"},
			expectBool: false,
		}),
		Entry("deviceTag exists among multiple - machine tag", testCaseDeviceTagIsInStringList{
			deviceTag:  DeviceTag{DeviceTagKeyMachine, "mymachine"},
			tagList:    []string{"caphv-machine-name=mymachine", "othertag"},
			expectBool: true,
		}),
	)
})

var _ = Describe("Test DeviceTag.RemoveFromList", func() {
	type testCaseDeviceTagRemoveFromList struct {
		deviceTag     DeviceTag
		tagList       []string
		expectTagList []string
		expectUpdated bool
	}

	DescribeTable("Test DeviceTag.RemoveFromList",
		func(tc testCaseDeviceTagRemoveFromList) {
			tagList, updated := tc.deviceTag.RemoveFromList(tc.tagList)
			Expect(tagList).Should(Equal(tc.expectTagList))
			Expect(updated).Should(Equal(tc.expectUpdated))
		},
		Entry("deviceTag exists among multiple", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyCluster, "mycluster"},
			tagList:       []string{"caphv-cluster-name=mycluster", "othertag", "caphv-machine-name=mymachine"},
			expectTagList: []string{"othertag", "caphv-machine-name=mymachine"},
			expectUpdated: true,
		}),
		Entry("deviceTag exists and key exists twice", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyCluster, "mycluster"},
			tagList:       []string{"caphv-cluster-name=mycluster", "caphv-cluster-name=othercluster", "caphv-machine-name=mymachine"},
			expectTagList: []string{"caphv-cluster-name=othercluster", "caphv-machine-name=mymachine"},
			expectUpdated: true,
		}),
		Entry("deviceTag exists twice", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyCluster, "mycluster"},
			tagList:       []string{"caphv-cluster-name=mycluster", "caphv-cluster-name=mycluster", "caphv-machine-name=mymachine"},
			expectTagList: []string{"caphv-machine-name=mymachine"},
			expectUpdated: true,
		}),
		Entry("deviceTag exists among multiple - machine tag", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyMachine, "mymachine"},
			tagList:       []string{"caphv-cluster-type=mytype", "othertag", "caphv-machine-name=mymachine"},
			expectTagList: []string{"caphv-cluster-type=mytype", "othertag"},
			expectUpdated: true,
		}),
		Entry("deviceTag does not exist among multiple", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyMachine, "mynewmachine"},
			tagList:       []string{"caphv-cluster-type=mytype", "othertag", "caphv-machine-name=mymachine"},
			expectTagList: []string{"caphv-cluster-type=mytype", "othertag", "caphv-machine-name=mymachine"},
			expectUpdated: false,
		}),
		Entry("empty tag list", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyMachine, "mynewmachine"},
			tagList:       []string{},
			expectTagList: []string{},
			expectUpdated: false,
		}),
		Entry("nil tag list", testCaseDeviceTagRemoveFromList{
			deviceTag:     DeviceTag{DeviceTagKeyMachine, "mynewmachine"},
			tagList:       nil,
			expectTagList: []string{},
			expectUpdated: false,
		}),
	)
})

var _ = Describe("RemoveEphemeralTags", func() {
	It("removes ephemeral tags, but keeps permanent tags", func() {
		newTags := RemoveEphemeralTags([]string{
			// non-ephemeral (keep)
			DeviceTagKeyPermanentError.Prefix() + "my-permantent-error",
			DeviceTagKeyCAPHVUseAllowed.Prefix() + "allow",

			// remove these:
			DeviceTagKeyCluster.Prefix() + "my-cluster",
			DeviceTagKeyMachine.Prefix() + "my-machine",
			"some-other-tag",
		})
		Expect(newTags).To(Equal([]string{
			"caphv-permanent-error=my-permantent-error",
			"caphv-use=allow",
			"some-other-tag"}))
	})
})

var _ = Describe("PermanentErrorTagFromList", func() {
	It("return permantent error from list", func() {
		tag, err := PermanentErrorTagFromList([]string{
			DeviceTagKeyPermanentError.Prefix() + "my-permantent-error",
			DeviceTagKeyCluster.Prefix() + "my-cluster",
			DeviceTagKeyMachine.Prefix() + "my-machine",
			"some-other-tag",
		})
		Expect(err).To(BeNil())
		Expect(tag).To(Equal(DeviceTag{
			Key:   "caphv-permanent-error",
			Value: "my-permantent-error",
		}))
	})
})
