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
	"fmt"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/utils"
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

var _ = Describe("HivelocityMachineReconciler", func() {
	var (
		capiCluster *clusterv1.Cluster
		capiMachine *clusterv1.Machine

		infraCluster *infrav1.HivelocityCluster
		infraMachine *infrav1.HivelocityMachine

		testNs *corev1.Namespace

		hvSecret        *corev1.Secret
		bootstrapSecret *corev1.Secret

		key client.ObjectKey

		hvClient hvclient.Client
	)

	BeforeEach(func() {
		var err error
		testNs, err = testEnv.CreateNamespace(ctx, "hivelocitymachine-reconciler")
		Expect(err).NotTo(HaveOccurred())

		capiCluster = &clusterv1.Cluster{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test1-",
				Namespace:    testNs.Name,
				Finalizers:   []string{clusterv1.ClusterFinalizer},
			},
			Spec: clusterv1.ClusterSpec{
				InfrastructureRef: &corev1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:       "HivelocityCluster",
					Name:       "hv-test1",
					Namespace:  testNs.Name,
				},
			},
			Status: clusterv1.ClusterStatus{
				InfrastructureReady: true,
			},
		}
		Expect(testEnv.Create(ctx, capiCluster)).To(Succeed())

		infraCluster = &infrav1.HivelocityCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hv-test1",
				Namespace: testNs.Name,
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

		hvSecret = getDefaultHivelocitySecret(testNs.Name)
		Expect(testEnv.Create(ctx, hvSecret)).To(Succeed())

		bootstrapSecret = getDefaultBootstrapSecret(testNs.Name)
		Expect(testEnv.Create(ctx, bootstrapSecret)).To(Succeed())

		hvClient = testEnv.HVClientFactory.NewClient("my-api-key")
	})

	AfterEach(func() {
		Expect(testEnv.Cleanup(ctx, testNs, capiCluster, infraCluster, capiMachine,
			infraMachine, hvSecret, bootstrapSecret)).To(Succeed())
	})

	Context("Basic test", func() {

		BeforeEach(func() {
			hivelocityMachineName := utils.GenerateName(nil, "hv-machine-")

			capiMachine = &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "capi-machine-",
					Namespace:    testNs.Name,
					Finalizers:   []string{clusterv1.MachineFinalizer},
					Labels: map[string]string{
						clusterv1.ClusterLabelName: capiCluster.Name,
					},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: capiCluster.Name,
					InfrastructureRef: corev1.ObjectReference{
						APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
						Kind:       "HivelocityMachine",
						Name:       hivelocityMachineName,
					},
					FailureDomain: &defaultFailureDomain,
				},
			}
			Expect(testEnv.Create(ctx, capiMachine)).To(Succeed())

			infraMachine = &infrav1.HivelocityMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      hivelocityMachineName,
					Namespace: testNs.Name,
					Labels: map[string]string{
						clusterv1.ClusterLabelName:             capiCluster.Name,
						clusterv1.MachineControlPlaneLabelName: "",
					},
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
					ImageName: "fedora-control-plane",
					Type:      "hvCustom",
				},
			}

			Expect(testEnv.Create(ctx, infraMachine)).To(Succeed())
			Expect(testEnv.Create(ctx, infraCluster)).To(Succeed())

			key = client.ObjectKey{Namespace: testNs.Name, Name: infraMachine.Name}
		})

		It("creates the infra machine", func() {
			Eventually(func() bool {
				if err := testEnv.Get(ctx, key, infraMachine); err != nil {
					return false
				}
				return true
			}, timeout).Should(BeTrue())
		})

		It("creates the Hivelocity machine in Hivelocity", func() {
			// Check that there is no machine yet.
			Eventually(func() bool {
				servers, err := hvClient.ListServers(ctx)
				if err != nil {
					return false
				}
				return len(servers) > 0
			}).Should(BeTrue())

			// Check whether bootstrap condition is not ready
			Eventually(func() bool {
				if err := testEnv.Get(ctx, key, infraMachine); err != nil {
					return false
				}
				return isPresentAndFalseWithReason(key, infraMachine, infrav1.InstanceBootstrapReadyCondition, infrav1.InstanceBootstrapNotReadyReason)
			}, timeout, time.Second).Should(BeTrue())

			By("setting the bootstrap data")
			Eventually(func() error {
				ph, err := patch.NewHelper(capiMachine, testEnv)
				Expect(err).ShouldNot(HaveOccurred())
				capiMachine.Spec.Bootstrap = clusterv1.Bootstrap{
					DataSecretName: pointer.String("bootstrap-secret"),
				}
				return ph.Patch(ctx, capiMachine, patch.WithStatusObservedGeneration{})
			}, timeout, time.Second).Should(BeNil())

			// Check whether bootstrap condition is ready
			Eventually(func() bool {
				if err := testEnv.Get(ctx, key, infraMachine); err != nil {
					return false
				}
				objectCondition := conditions.Get(infraMachine, infrav1.InstanceBootstrapReadyCondition)
				fmt.Println(objectCondition)
				return isPresentAndTrue(key, infraMachine, infrav1.InstanceBootstrapReadyCondition)
			}, timeout, time.Second).Should(BeTrue())

			Eventually(func() int {
				servers, err := hvClient.ListServers(ctx)
				if err != nil {
					return 0
				}
				return len(servers)
			}, timeout, time.Second).Should(BeNumerically(">", 0))
		})
	})

	Context("various specs", func() {

		BeforeEach(func() {
			hivelocityMachineName := utils.GenerateName(nil, "hv-machine-")

			capiMachine = &clusterv1.Machine{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "capi-machine-",
					Namespace:    testNs.Name,
					Finalizers:   []string{clusterv1.MachineFinalizer},
					Labels: map[string]string{
						clusterv1.ClusterLabelName: capiCluster.Name,
					},
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: capiCluster.Name,
					InfrastructureRef: corev1.ObjectReference{
						APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
						Kind:       "HivelocityMachine",
						Name:       hivelocityMachineName,
					},
					FailureDomain: &defaultFailureDomain,
					Bootstrap: clusterv1.Bootstrap{
						DataSecretName: pointer.String("bootstrap-secret"),
					},
				},
			}
			Expect(testEnv.Create(ctx, capiMachine)).To(Succeed())

			infraMachine = &infrav1.HivelocityMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      hivelocityMachineName,
					Namespace: testNs.Name,
					Labels: map[string]string{
						clusterv1.ClusterLabelName:             capiCluster.Name,
						clusterv1.MachineControlPlaneLabelName: "",
					},
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
					ImageName: "fedora-control-plane",
					Type:      "hvCustom",
				},
			}

			key = client.ObjectKey{Namespace: testNs.Name, Name: infraMachine.Name}
		})

		Context("with public network specs", func() {
			BeforeEach(func() {
				Expect(testEnv.Create(ctx, infraCluster)).To(Succeed())
				Expect(testEnv.Create(ctx, infraMachine)).To(Succeed())
			})

			It("creates the Hivelocity machine in Hivelocity", func() {
				Eventually(func() int {
					servers, err := hvClient.ListServers(ctx)
					if err != nil {
						return 0
					}
					return len(servers)
				}, timeout, time.Second).Should(BeNumerically(">", 0))
			})
		})
	})
})

var _ = Describe("Hivelocity secret", func() {
	var (
		hvCluster         *infrav1.HivelocityCluster
		capiCluster       *clusterv1.Cluster
		hivelocityMachine *infrav1.HivelocityMachine
		capiMachine       *clusterv1.Machine
		key               client.ObjectKey
		hvSecret          *corev1.Secret
		hvClusterName     string
	)

	BeforeEach(func() {
		var err error
		Expect(err).NotTo(HaveOccurred())

		hvClusterName = utils.GenerateName(nil, "hv-cluster-test")
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
					Name:       hvClusterName,
					Namespace:  "default",
				},
			},
		}
		Expect(testEnv.Create(ctx, capiCluster)).To(Succeed())

		hvCluster = &infrav1.HivelocityCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      hvClusterName,
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
		Expect(testEnv.Create(ctx, hvCluster)).To(Succeed())

		hivelocityMachineName := utils.GenerateName(nil, "hv-machine-")

		capiMachine = &clusterv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "capi-machine-",
				Namespace:    "default",
				Finalizers:   []string{clusterv1.MachineFinalizer},
				Labels: map[string]string{
					clusterv1.ClusterLabelName: capiCluster.Name,
				},
			},
			Spec: clusterv1.MachineSpec{
				ClusterName: capiCluster.Name,
				InfrastructureRef: corev1.ObjectReference{
					APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
					Kind:       "HivelocityMachine",
					Name:       hivelocityMachineName,
				},
				FailureDomain: &defaultFailureDomain,
				Bootstrap: clusterv1.Bootstrap{
					DataSecretName: pointer.String("bootstrap-secret"),
				},
			},
		}
		Expect(testEnv.Create(ctx, capiMachine)).To(Succeed())

		hivelocityMachine = &infrav1.HivelocityMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      hivelocityMachineName,
				Namespace: "default",
				Labels: map[string]string{
					clusterv1.ClusterLabelName:             capiCluster.Name,
					clusterv1.MachineControlPlaneLabelName: "",
				},
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
				ImageName: "fedora-control-plane",
				Type:      "hvCustom",
			},
		}
		Expect(testEnv.Create(ctx, hivelocityMachine)).To(Succeed())
		key = client.ObjectKey{Namespace: "default", Name: hivelocityMachine.Name}
	})

	AfterEach(func() {
		Expect(testEnv.Cleanup(ctx, hvCluster, capiCluster, hivelocityMachine, capiMachine, hvSecret)).To(Succeed())

		Eventually(func() bool {
			if err := testEnv.Get(ctx, client.ObjectKey{Namespace: hvSecret.Namespace, Name: hvSecret.Name}, hvSecret); err != nil && apierrors.IsNotFound(err) {
				return true
			} else if err != nil {
				return false
			}
			// Secret still there, so the finalizers have not been removed. Patch to remove them.
			ph, err := patch.NewHelper(hvSecret, testEnv)
			Expect(err).ShouldNot(HaveOccurred())
			hvSecret.Finalizers = nil
			Expect(ph.Patch(ctx, hvSecret, patch.WithStatusObservedGeneration{})).To(Succeed())
			// Should delete secret
			if err := testEnv.Delete(ctx, hvSecret); err != nil && apierrors.IsNotFound(err) {
				// Has been deleted already
				return true
			}
			return false
		}, time.Second, time.Second).Should(BeTrue())
	})

	DescribeTable("test different hv secret",
		func(secret corev1.Secret, expectedReason string) {
			hvSecret = &secret
			Expect(testEnv.Create(ctx, hvSecret)).To(Succeed())

			Eventually(func() bool {
				if err := testEnv.Get(ctx, key, hivelocityMachine); err != nil {
					return false
				}
				return isPresentAndFalseWithReason(key, hivelocityMachine, infrav1.InstanceReadyCondition, expectedReason)
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
		Entry("empty API key", corev1.Secret{
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

var _ = Describe("HivelocityMachine validation", func() {
	var (
		infraMachine *infrav1.HivelocityMachine
		testNs       *corev1.Namespace
	)

	BeforeEach(func() {
		var err error
		testNs, err = testEnv.CreateNamespace(ctx, "hivelocitymachine-validation")
		Expect(err).NotTo(HaveOccurred())

		infraMachine = &infrav1.HivelocityMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hv-validation-machine",
				Namespace: testNs.Name,
			},
			Spec: infrav1.HivelocityMachineSpec{
				ImageName: "fedora-control-plane",
				Type:      "hvCustom",
			},
		}
	})

	AfterEach(func() {
		Expect(testEnv.Cleanup(ctx, testNs, infraMachine)).To(Succeed())
	})

	It("should fail with wrong type", func() {
		infraMachine.Spec.Type = "wrong-type"
		Expect(testEnv.Create(ctx, infraMachine)).ToNot(Succeed())
	})

	It("should fail without imageName", func() {
		infraMachine.Spec.ImageName = ""
		Expect(testEnv.Create(ctx, infraMachine)).ToNot(Succeed())
	})
})