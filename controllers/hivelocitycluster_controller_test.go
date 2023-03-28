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
	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

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
