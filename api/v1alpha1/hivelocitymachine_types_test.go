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
	hv "github.com/hivelocity/hivelocity-client-go/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
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

var _ = Describe("Test SetFailure", func() {
	hvMachine := HivelocityMachine{}
	newFailureMessage := "bad error"
	newFailureReason := capierrors.CreateMachineError

	It("sets new failure on the machine with existing failure", func() {
		failureMessage := "first message"
		failureReason := capierrors.MachineStatusError("first error")
		hvMachine.Status.FailureMessage = &failureMessage
		hvMachine.Status.FailureReason = &failureReason

		hvMachine.SetFailure(newFailureReason, newFailureMessage)

		Expect(hvMachine.Status.FailureMessage).ToNot(BeNil())
		Expect(hvMachine.Status.FailureReason).ToNot(BeNil())
		Expect(*hvMachine.Status.FailureMessage).To(Equal(newFailureMessage))
		Expect(*hvMachine.Status.FailureReason).To(Equal(newFailureReason))
	})

	It("sets new failure on the machine without existing failure", func() {
		hvMachine.SetFailure(newFailureReason, newFailureMessage)

		Expect(hvMachine.Status.FailureMessage).ToNot(BeNil())
		Expect(hvMachine.Status.FailureReason).ToNot(BeNil())
		Expect(*hvMachine.Status.FailureMessage).To(Equal(newFailureMessage))
		Expect(*hvMachine.Status.FailureReason).To(Equal(newFailureReason))
	})
})

var _ = Describe("Test SetMachineStatus", func() {
	type testCaseSetMachineStatus struct {
		existingStatus HivelocityMachineStatus
		device         hv.BareMetalDevice
		expectStatus   HivelocityMachineStatus
	}

	DescribeTable("Test SetMachineStatus",
		func(tc testCaseSetMachineStatus) {
			hvMachine := HivelocityMachine{}
			hvMachine.Status = tc.existingStatus
			hvMachine.SetMachineStatus(tc.device)

			Expect(hvMachine.Status).Should(Equal(tc.expectStatus))
		},
		Entry("existing status", testCaseSetMachineStatus{
			existingStatus: HivelocityMachineStatus{
				Addresses: []clusterv1.MachineAddress{
					{
						Type:    clusterv1.MachineHostName,
						Address: "hostname",
					},
				},
				Region:     Region("testregion"),
				PowerState: "ON",
			},
			device: hv.BareMetalDevice{
				Hostname:     "device-hostname",
				PrimaryIp:    "127.0.0.1",
				LocationName: "LAX2",
				PowerStatus:  "OFF",
			},
			expectStatus: HivelocityMachineStatus{
				Addresses: []clusterv1.MachineAddress{
					{
						Type:    clusterv1.MachineHostName,
						Address: "device-hostname",
					},
					{
						Type:    clusterv1.MachineInternalIP,
						Address: "127.0.0.1",
					},
					{
						Type:    clusterv1.MachineExternalIP,
						Address: "127.0.0.1",
					},
				},
				Region:     Region("LAX2"),
				PowerState: "OFF",
			},
		}),
		Entry("no existing status", testCaseSetMachineStatus{
			existingStatus: HivelocityMachineStatus{},
			device: hv.BareMetalDevice{
				Hostname:     "device-hostname",
				PrimaryIp:    "127.0.0.1",
				LocationName: "LAX2",
				PowerStatus:  "OFF",
			},
			expectStatus: HivelocityMachineStatus{
				Addresses: []clusterv1.MachineAddress{
					{
						Type:    clusterv1.MachineHostName,
						Address: "device-hostname",
					},
					{
						Type:    clusterv1.MachineInternalIP,
						Address: "127.0.0.1",
					},
					{
						Type:    clusterv1.MachineExternalIP,
						Address: "127.0.0.1",
					},
				},
				Region:     Region("LAX2"),
				PowerState: "OFF",
			},
		}),
	)
})

var _ = Describe("Test providerIDFromDeviceID", func() {
	Expect(providerIDFromDeviceID(42)).To(Equal("hivelocity://42"))
})

var _ = Describe("Test SetProviderID", func() {
	type testCaseSetProviderID struct {
		existingProviderID *string
		deviceID           int32
		expectProviderID   *string
	}

	DescribeTable("Test SetProviderID",
		func(tc testCaseSetProviderID) {
			hvMachine := HivelocityMachine{}
			hvMachine.Spec.ProviderID = tc.existingProviderID
			hvMachine.SetProviderID(tc.deviceID)

			Expect(hvMachine.Spec.ProviderID).Should(Equal(tc.expectProviderID))
		},
		Entry("existing providerID", testCaseSetProviderID{
			existingProviderID: pointer.String("hivelocity://42"),
			deviceID:           1,
			expectProviderID:   pointer.String("hivelocity://1"),
		}),
		Entry("no existing status", testCaseSetProviderID{
			existingProviderID: nil,
			deviceID:           1,
			expectProviderID:   pointer.String("hivelocity://1"),
		}),
	)
})
