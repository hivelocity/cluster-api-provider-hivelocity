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

package controllers

import (
	"sync"
	"testing"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/test/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubectl/pkg/scheme"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

const (
	defaultPodNamespace = "caphv-system"
	timeout             = time.Second * 5
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var (
	testEnv *helpers.TestEnvironment
	ctx     = ctrl.SetupSignalHandler()
	wg      sync.WaitGroup

	defaultFailureDomain = "LAX2"
)

var _ = BeforeSuite(func() {
	utilruntime.Must(infrav1.AddToScheme(scheme.Scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme.Scheme))

	testEnv = helpers.NewTestEnvironment()

	wg.Add(1)

	Expect((&HivelocityClusterReconciler{
		Client:           testEnv.Manager.GetClient(),
		APIReader:        testEnv.Manager.GetAPIReader(),
		HVClientFactory:  testEnv.HVClientFactory,
		WatchFilterValue: "",
	}).SetupWithManager(ctx, testEnv.Manager, controller.Options{})).To(Succeed())

	Expect((&HivelocityMachineReconciler{
		Client:           testEnv.Manager.GetClient(),
		APIReader:        testEnv.Manager.GetAPIReader(),
		HVClientFactory:  testEnv.HVClientFactory,
		WatchFilterValue: "",
	}).SetupWithManager(ctx, testEnv.Manager, controller.Options{})).To(Succeed())

	Expect((&HivelocityMachineTemplateReconciler{
		Client:           testEnv.Manager.GetClient(),
		APIReader:        testEnv.Manager.GetAPIReader(),
		HVClientFactory:  testEnv.HVClientFactory,
		WatchFilterValue: "",
	}).SetupWithManager(ctx, testEnv.Manager, controller.Options{})).To(Succeed())

	go func() {
		defer GinkgoRecover()
		Expect(testEnv.StartManager(ctx)).To(Succeed())
	}()

	<-testEnv.Manager.Elected()

	// wait for webhook port to be open prior to running tests
	testEnv.WaitForWebhooks()

	// create manager pod namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: defaultPodNamespace,
		},
	}

	Expect(testEnv.Create(ctx, ns)).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(testEnv.Stop()).To(Succeed())
	wg.Done() // Main manager has been stopped
	wg.Wait() // Wait for target cluster manager
})

func getDefaultHivelocityClusterSpec() infrav1.HivelocityClusterSpec {
	return infrav1.HivelocityClusterSpec{
		ControlPlaneEndpoint: &clusterv1.APIEndpoint{},
		HivelocitySecret: infrav1.HivelocitySecretRef{
			Key:  "HIVELOCITY_API_KEY",
			Name: "hv-secret",
		},
		SSHKey:             &infrav1.SSHKey{Name: "testsshkey"},
		ControlPlaneRegion: "LAX2",
	}
}

func getDefaultHivelocitySecret(namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hv-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"HIVELOCITY_API_KEY": []byte("my-api-key"),
		},
	}
}

func getDefaultBootstrapSecret(namespace string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bootstrap-secret",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"value": []byte("my-bootstrap"),
		},
	}
}

func isPresentAndFalseWithReason(key types.NamespacedName, getter conditions.Getter, condition clusterv1.ConditionType, reason string) bool {
	ExpectWithOffset(1, testEnv.Get(ctx, key, getter)).To(Succeed())
	if !conditions.Has(getter, condition) {
		return false
	}
	objectCondition := conditions.Get(getter, condition)
	return objectCondition.Status == corev1.ConditionFalse &&
		objectCondition.Reason == reason
}

func isPresentAndTrue(key types.NamespacedName, getter conditions.Getter, condition clusterv1.ConditionType) bool {
	ExpectWithOffset(1, testEnv.Get(ctx, key, getter)).To(Succeed())
	if !conditions.Has(getter, condition) {
		return false
	}
	objectCondition := conditions.Get(getter, condition)
	return objectCondition.Status == corev1.ConditionTrue
}
