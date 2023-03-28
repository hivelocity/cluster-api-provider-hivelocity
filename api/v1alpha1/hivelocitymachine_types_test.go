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

package v1alpha1

import (
	"testing"

	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/hvtag"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMachineDeviceTag(t *testing.T) {
	hvMachine := HivelocityMachine{}
	hvMachine.Name = "hvmachinename"
	deviceTag := hvMachine.DeviceTag()
	expectDeviceTag := hvtag.DeviceTag{Key: hvtag.DeviceTagKeyMachine, Value: "hvmachinename"}
	if deviceTag != expectDeviceTag {
		t.Fatalf("wrong device tag. Expect %+v, got %+v", expectDeviceTag, deviceTag)
	}
}

var _ = Describe("Test DeviceIDFromProviderID", func() {
	It("gives error on nil providerID", func() {
		hvMachine := HivelocityMachine{}
		deviceID, err := hvMachine.DeviceIDFromProviderID()
		Expect(err).ToNot(BeNil())
		Expect(err).To(MatchError(ErrEmptyProviderID))
		Expect(deviceID).To(Equal(int32(0)))
	})
	type testCaseDeviceIDFromProviderID struct {
		providerID     string
		expectDeviceID int32
		expectError    error
	}

	DescribeTable("Test DeviceIDFromProviderID",
		func(tc testCaseDeviceIDFromProviderID) {
			hvMachine := HivelocityMachine{}
			hvMachine.Spec.ProviderID = &tc.providerID
			deviceID, err := hvMachine.DeviceIDFromProviderID()
			if tc.expectError != nil {
				Expect(err).To(MatchError(tc.expectError))
			} else {
				Expect(err).To(BeNil())
			}
			Expect(deviceID).Should(Equal(tc.expectDeviceID))
		},
		Entry("empty providerID", testCaseDeviceIDFromProviderID{
			providerID:     "",
			expectDeviceID: 0,
			expectError:    ErrEmptyProviderID,
		}),
		Entry("wrong prefix", testCaseDeviceIDFromProviderID{
			providerID:     "hivelocit://42",
			expectDeviceID: 0,
			expectError:    ErrInvalidProviderID,
		}),
		Entry("no prefix", testCaseDeviceIDFromProviderID{
			providerID:     "42",
			expectDeviceID: 0,
			expectError:    ErrInvalidProviderID,
		}),
		Entry("no deviceID", testCaseDeviceIDFromProviderID{
			providerID:     "hivelocity://",
			expectDeviceID: 0,
			expectError:    ErrInvalidDeviceID,
		}),
		Entry("invalid deviceID - no int", testCaseDeviceIDFromProviderID{
			providerID:     "hivelocity://deviceID",
			expectDeviceID: 0,
			expectError:    ErrInvalidDeviceID,
		}),
		Entry("invalid deviceID - too long", testCaseDeviceIDFromProviderID{
			providerID:     "hivelocity://9999999999999999999999999999999",
			expectDeviceID: 0,
			expectError:    ErrInvalidDeviceID,
		}),
		Entry("correct providerID", testCaseDeviceIDFromProviderID{
			providerID:     "hivelocity://42",
			expectDeviceID: 42,
			expectError:    nil,
		}),
	)
})
