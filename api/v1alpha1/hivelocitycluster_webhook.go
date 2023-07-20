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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var hivelocityclusterlog = logf.Log.WithName("hivelocitycluster-resource")

func (r *HivelocityCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocitycluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters,verbs=create;update,versions=v1alpha1,name=mhivelocitycluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HivelocityCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *HivelocityCluster) Default() {}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocitycluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclusters,verbs=create;update,versions=v1alpha1,name=vhivelocitycluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HivelocityCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityCluster) ValidateCreate() error {
	hivelocityclusterlog.V(1).Info("validate create", "name", r.Name)
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityCluster) ValidateUpdate(_ runtime.Object) error {
	hivelocityclusterlog.V(1).Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityCluster) ValidateDelete() error {
	hivelocityclusterlog.V(1).Info("validate delete", "name", r.Name)
	return nil
}
