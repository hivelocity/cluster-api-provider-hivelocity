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

package controllers

import (
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/utils"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Hivelocity ClusterReconciler", func() {
	It("should create a basic cluster", func() {
		// Create the secret
		hivelocitySecret := getDefaultHivelocitySecret("default")
		Expect(testEnv.Create(ctx, hivelocitySecret)).To(Succeed())
		defer func() {
			Expect(testEnv.Cleanup(ctx, hivelocitySecret)).To(Succeed())
		}()

		// Create the HivelocityCluster object
		instance := &infrav1.HivelocityCluster{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "hivelocity-test1",
				Namespace:    "default",
			},
			Spec: getDefaultHivelocityClusterSpec(),
		}
		Expect(testEnv.Create(ctx, instance)).To(Succeed())
		defer func() {
			Expect(testEnv.Delete(ctx, instance)).To(Succeed())
		}()

		key := client.ObjectKey{Namespace: instance.Namespace, Name: instance.Name}

		// Create capi cluster
		capiCluster := &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test1-",
				Namespace:    "default",
				Finalizers:   []string{clusterv1.ClusterFinalizer},
			},
			Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "HivelocityCluster",
					Name:       instance.Name,
				},
			},
		}
		Expect(testEnv.Create(ctx, capiCluster)).To(Succeed())
		defer func() {
			Expect(testEnv.Cleanup(ctx, capiCluster)).To(Succeed())
		}()

		// Make sure the HivelocityCluster exists.
		Eventually(func() error {
			return testEnv.Get(ctx, key, instance)
		}, timeout, 100*time.Millisecond).Should(BeNil())

		By("setting the OwnerRef on the HivelocityCluster")
		// Set owner reference to Hivelocity cluster
		Eventually(func() error {
			ph, err := patch.NewHelper(instance, testEnv)
			Expect(err).ShouldNot(HaveOccurred())
			instance.OwnerReferences = append(instance.OwnerReferences, metav1.OwnerReference{
				Kind:       "Cluster",
				APIVersion: clusterv1.GroupVersion.String(),
				Name:       capiCluster.Name,
				UID:        capiCluster.UID,
			})
			return ph.Patch(ctx, instance, patch.WithStatusObservedGeneration{})
		}, timeout, 100*time.Millisecond).Should(BeNil())

		// Check whether finalizer has been set for HivelocityCluster
		Eventually(func() bool {
			if err := testEnv.Get(ctx, key, instance); err != nil {
				return false
			}
			return len(instance.Finalizers) > 0
		}, timeout, 100*time.Millisecond).Should(BeTrue())
	})

	Context("For HivelocityMachines belonging to the cluster", func() {
		var (
			namespace        string
			testNs           *corev1.Namespace
			hivelocitySecret *corev1.Secret
			bootstrapSecret  *corev1.Secret
		)

		BeforeEach(func() {
			var err error
			testNs, err = testEnv.CreateNamespace(ctx, "hivelocity-owner-ref")
			Expect(err).NotTo(HaveOccurred())
			namespace = testNs.Name

			// Create the hivelocity secret
			hivelocitySecret = getDefaultHivelocitySecret(namespace)
			Expect(testEnv.Create(ctx, hivelocitySecret)).To(Succeed())
			// Create the bootstrap secret
			bootstrapSecret = getDefaultBootstrapSecret(namespace)
			Expect(testEnv.Create(ctx, bootstrapSecret)).To(Succeed())
		})

		AfterEach(func() {
			Expect(testEnv.Delete(ctx, bootstrapSecret)).To(Succeed())
			Expect(testEnv.Delete(ctx, hivelocitySecret)).To(Succeed())
			Expect(testEnv.Delete(ctx, testNs)).To(Succeed())
		})

		It("sets owner references to those machines", func() {
			// Create the HivelocityCluster object
			hvCluster := &infrav1.HivelocityCluster{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "test1-",
					Namespace:    namespace,
				},
				Spec: getDefaultHivelocityClusterSpec(),
			}
			Expect(testEnv.Create(ctx, hvCluster)).To(Succeed())
			defer func() {
				Expect(testEnv.Cleanup(ctx, hvCluster)).To(Succeed())
			}()
			capiCluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "capi-test1-",
					Namespace:    namespace,
					Finalizers:   []string{clusterv1.ClusterFinalizer},
				},
				Spec: clusterv1.ClusterSpec{
					InfrastructureRef: &corev1.ObjectReference{
						APIVersion: infrav1.GroupVersion.String(),
						Kind:       "HivelocityCluster",
						Name:       hvCluster.Name,
						Namespace:  namespace,
					},
				},
			}
			Expect(testEnv.Create(ctx, capiCluster)).To(Succeed())
			defer func() {
				Expect(testEnv.Cleanup(ctx, capiCluster)).To(Succeed())
			}()

			// Make sure the HivelocityCluster exists.
			Eventually(func() error {
				return testEnv.Get(ctx, client.ObjectKey{Namespace: namespace, Name: hvCluster.Name}, hvCluster)
			}, timeout, 100*time.Millisecond).Should(BeNil())

			// Create machines
			machineCount := 3
			for i := 0; i < machineCount; i++ {
				hvMachineName := utils.GenerateName(nil, "hv-machine")
				capiMachine := &clusterv1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						GenerateName: "capi-machine-",
						Namespace:    namespace,
						Finalizers:   []string{clusterv1.MachineFinalizer},
						Labels: map[string]string{
							clusterv1.ClusterNameLabel: capiCluster.Name,
						},
					},
					Spec: clusterv1.MachineSpec{
						ClusterName: capiCluster.Name,
						InfrastructureRef: corev1.ObjectReference{
							APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
							Kind:       "HivelocityMachine",
							Name:       hvMachineName,
						},
						FailureDomain: &defaultFailureDomain,
						Bootstrap: clusterv1.Bootstrap{
							DataSecretName: pointer.String("bootstrap-secret"),
						},
					},
				}
				Expect(testEnv.Create(ctx, capiMachine)).To(Succeed())

				hvMachine := &infrav1.HivelocityMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      hvMachineName,
						Namespace: namespace,
						Labels:    map[string]string{clusterv1.ClusterNameLabel: capiCluster.Name},
						OwnerReferences: []metav1.OwnerReference{
							{
								APIVersion: clusterv1.GroupVersion.String(),
								Kind:       "Machine",
								Name:       capiMachine.Name,
								UID:        capiMachine.UID,
							},
						},
					},
					Spec: infrav1.HivelocityMachineSpec{
						ImageName: "Ubuntu 20.x",
						Type:      "pool",
					},
				}
				Expect(testEnv.Create(ctx, hvMachine)).To(Succeed())
			}

			// Set owner reference to HivelocityCluster
			Eventually(func() bool {
				ph, err := patch.NewHelper(hvCluster, testEnv)
				Expect(err).ShouldNot(HaveOccurred())
				hvCluster.OwnerReferences = append(hvCluster.OwnerReferences, metav1.OwnerReference{
					Kind:       "Cluster",
					APIVersion: clusterv1.GroupVersion.String(),
					Name:       capiCluster.Name,
					UID:        capiCluster.UID,
				})
				Expect(ph.Patch(ctx, hvCluster, patch.WithStatusObservedGeneration{})).ShouldNot(HaveOccurred())
				return true
			}, timeout, 100*time.Millisecond).Should(BeTrue())

			By("checking for presence of HivelocityMachine objects")

			hvClient := testEnv.HVClientFactory.NewClient("dummy-key")
			resourceOwnedTag := hvCluster.DeviceTagOwned()

			// Check if devices have been associated
			Eventually(func() int {
				devices, err := hvClient.ListDevices(ctx)
				if err != nil {
					return -1
				}

				// assume devices are associated of the resource owned tag is set
				provisionedDevices := make([]hv.BareMetalDevice, 0, len(devices))
				for _, device := range devices {
					if resourceOwnedTag.IsInStringList(device.Tags) {
						provisionedDevices = append(provisionedDevices, device)
					}
				}
				return len(provisionedDevices)
			}, timeout, 100*time.Millisecond).Should(Equal(machineCount))

			// Check if all machines have been provisioned
			Eventually(func() int {
				hvMachineList := &infrav1.HivelocityMachineList{}
				if err := testEnv.Client.List(ctx, hvMachineList, client.InNamespace(namespace)); err != nil {
					return -1
				}

				// check whether machines are in provisioned state
				provisionedMachines := make([]infrav1.HivelocityMachine, 0, len(hvMachineList.Items))
				for _, machine := range hvMachineList.Items {
					if machine.Spec.Status.ProvisioningState == infrav1.StateDeviceProvisioned {
						provisionedMachines = append(provisionedMachines, machine)
					}
				}
				return len(provisionedMachines)
			}, timeout, 100*time.Millisecond).Should(Equal(machineCount))
		})
	})
})

var _ = Describe("HivelocityCluster validation", func() {
	var (
		hvCluster *infrav1.HivelocityCluster
		testNs    *corev1.Namespace
	)

	BeforeEach(func() {
		var err error
		testNs, err = testEnv.CreateNamespace(ctx, "hivelocitycluster-validation")
		Expect(err).NotTo(HaveOccurred())

		hvCluster = &infrav1.HivelocityCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hv-validation-cluster",
				Namespace: testNs.Name,
			},
			Spec: infrav1.HivelocityClusterSpec{
				ControlPlaneEndpoint: &clusterv1.APIEndpoint{},
				ControlPlaneRegion:   infrav1.Region("LAX2"),
				HivelocitySecret:     infrav1.HivelocitySecretRef{Name: "hivelocity", Key: "HIVELOCITY_API_KEY"},
				SSHKey:               &infrav1.SSHKey{Name: "sshkey"},
			},
		}
	})

	AfterEach(func() {
		Expect(testEnv.Cleanup(ctx, testNs, hvCluster)).To(Succeed())
	})

	It("should not fail with correct spec", func() {
		Expect(testEnv.Create(ctx, hvCluster)).To(Succeed())
	})
	It("should fail with wrong region", func() {
		hvCluster.Spec.ControlPlaneRegion = infrav1.Region("not-a-valid-region")
		Expect(testEnv.Create(ctx, hvCluster)).ToNot(Succeed())
	})

	It("should fail without ssh key name", func() {
		hvCluster.Spec.SSHKey.Name = ""
		Expect(testEnv.Create(ctx, hvCluster)).ToNot(Succeed())
	})
	It("should not fail without ssh key", func() {
		hvCluster.Spec.SSHKey = nil
		Expect(testEnv.Create(ctx, hvCluster)).To(Succeed())
	})
	It("should not fail without Hivelocity secret ref", func() {
		hvCluster.Spec.HivelocitySecret = infrav1.HivelocitySecretRef{}
		Expect(testEnv.Create(ctx, hvCluster)).To(Succeed())
	})
	It("should not fail without Hivelocity secret ref key", func() {
		hvCluster.Spec.HivelocitySecret.Name = ""
		hvCluster.Spec.HivelocitySecret.Key = ""
		Expect(testEnv.Create(ctx, hvCluster)).To(Succeed())
	})
})

var _ = Describe("reconcileRateLimit", func() {
	var hvCluster *infrav1.HivelocityCluster
	BeforeEach(func() {
		hvCluster = &infrav1.HivelocityCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "rate-limit-cluster",
				Namespace: "default",
			},
			Spec: getDefaultHivelocityClusterSpec(),
		}
	})

	It("returns wait== true if rate limit condition is set and time is not over", func() {
		conditions.MarkFalse(hvCluster, infrav1.HivelocityAPIReachableCondition, infrav1.RateLimitExceededReason, clusterv1.ConditionSeverityWarning, "")
		Expect(reconcileRateLimit(hvCluster)).To(BeTrue())
	})

	It("returns wait== false if rate limit condition is set and time is over", func() {
		conditions.MarkFalse(hvCluster, infrav1.HivelocityAPIReachableCondition, infrav1.RateLimitExceededReason, clusterv1.ConditionSeverityWarning, "")
		conditionList := hvCluster.GetConditions()
		conditionList[0].LastTransitionTime = metav1.NewTime(time.Now().Add(-time.Hour))
		Expect(reconcileRateLimit(hvCluster)).To(BeFalse())
	})

	It("returns wait== false if rate limit condition is set to false", func() {
		conditions.MarkTrue(hvCluster, infrav1.HivelocityAPIReachableCondition)
		Expect(reconcileRateLimit(hvCluster)).To(BeFalse())
	})

	It("returns wait== false if rate limit condition is not set", func() {
		Expect(reconcileRateLimit(hvCluster)).To(BeFalse())
	})
})

var _ = Describe("Hivelocity secret", func() {
	var (
		hivelocityCluster     *infrav1.HivelocityCluster
		capiCluster           *clusterv1.Cluster
		key                   client.ObjectKey
		hivelocitySecret      *corev1.Secret
		hivelocityClusterName string
	)

	BeforeEach(func() {
		var err error
		Expect(err).NotTo(HaveOccurred())

		hivelocityClusterName = utils.GenerateName(nil, "hivelocity-cluster-test")
		capiCluster = &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test1-",
				Namespace:    "default",
				Finalizers:   []string{clusterv1.ClusterFinalizer},
			},
			Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{
					APIVersion: infrav1.GroupVersion.String(),
					Kind:       "HivelocityCluster",
					Name:       hivelocityClusterName,
					Namespace:  "default",
				},
			},
		}
		Expect(testEnv.Create(ctx, capiCluster)).To(Succeed())

		hivelocityCluster = &infrav1.HivelocityCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      hivelocityClusterName,
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: "cluster.x-k8s.io/v1beta1",
						Kind:       "Cluster",
						Name:       capiCluster.Name,
						UID:        capiCluster.UID,
					},
				},
			},
			Spec: getDefaultHivelocityClusterSpec(),
		}
		Expect(testEnv.Create(ctx, hivelocityCluster)).To(Succeed())

		key = client.ObjectKey{Namespace: hivelocityCluster.Namespace, Name: hivelocityCluster.Name}
	})

	AfterEach(func() {
		Expect(testEnv.Cleanup(ctx, hivelocityCluster, capiCluster, hivelocitySecret)).To(Succeed())

		Eventually(func() bool {
			if err := testEnv.Get(ctx, client.ObjectKey{Namespace: hivelocitySecret.Namespace, Name: hivelocitySecret.Name}, hivelocitySecret); err != nil && apierrors.IsNotFound(err) {
				return true
			} else if err != nil {
				return false
			}
			// Secret still there, so the finalizers have not been removed. Patch to remove them.
			ph, err := patch.NewHelper(hivelocitySecret, testEnv)
			Expect(err).ShouldNot(HaveOccurred())
			hivelocitySecret.Finalizers = nil
			Expect(ph.Patch(ctx, hivelocitySecret, patch.WithStatusObservedGeneration{})).To(Succeed())
			// Should delete secret
			if err := testEnv.Delete(ctx, hivelocitySecret); err != nil && apierrors.IsNotFound(err) {
				// Has been deleted already
				return true
			}
			return false
		}, time.Second, time.Second).Should(BeTrue())
	})

	DescribeTable("test different hivelocity secret",
		func(secret corev1.Secret, expectedReason string) {
			hivelocitySecret = &secret
			Expect(testEnv.Create(ctx, hivelocitySecret)).To(Succeed())

			Eventually(func() bool {
				if err := testEnv.Get(ctx, key, hivelocityCluster); err != nil {
					return false
				}
				return isPresentAndFalseWithReason(key, hivelocityCluster, infrav1.CredentialsAvailableCondition, expectedReason)
			}, timeout, time.Second).Should(BeTrue())
		},
		Entry("no Hivelocity secret/wrong reference", corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "wrong-name",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"HIVELOCITY_API_KEY": []byte("my-api-key"),
			},
		}, infrav1.HivelocitySecretUnreachableReason),
		Entry("empty hivelocity api key", corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hv-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"HIVELOCITY_API_KEY": []byte(""),
			},
		}, infrav1.HivelocityCredentialsInvalidReason),
		Entry("wrong key in secret", corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hv-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"wrongkey": []byte("my-api-key"),
			},
		}, infrav1.HivelocityCredentialsInvalidReason),
	)
})
