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
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var hivelocityclustertemplatelog = logf.Log.WithName("hivelocityclustertemplate-resource")

func (r *HivelocityClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocityclustertemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclustertemplates,verbs=create;update,versions=v1alpha1,name=mhivelocityclustertemplate.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HivelocityClusterTemplate{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *HivelocityClusterTemplate) Default() {}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocityclustertemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocityclustertemplates,verbs=create;update,versions=v1alpha1,name=vhivelocityclustertemplate.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HivelocityClusterTemplate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityClusterTemplate) ValidateCreate() (admission.Warnings, error) {
	hivelocityclustertemplatelog.V(1).Info("validate create", "name", r.Name)
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityClusterTemplate) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	hivelocityclustertemplatelog.V(1).Info("validate update", "name", r.Name)
	old, ok := oldRaw.(*HivelocityClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an HivelocityClusterTemplate but got a %T", oldRaw))
	}

	if !reflect.DeepEqual(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("HivelocityClusterTemplate.Spec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityClusterTemplate) ValidateDelete() (admission.Warnings, error) {
	hivelocityclustertemplatelog.V(1).Info("validate delete", "name", r.Name)
	return nil, nil
}
