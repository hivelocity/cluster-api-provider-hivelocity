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
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var hivelocitymachinelog = logf.Log.WithName("hivelocitymachine-resource")

func (r *HivelocityMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocitymachine,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachines,verbs=create;update,versions=v1alpha1,name=mhivelocitymachine.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HivelocityMachine{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *HivelocityMachine) Default() {}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocitymachine,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachines,verbs=create;update,versions=v1alpha1,name=vhivelocitymachine.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HivelocityMachine{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityMachine) ValidateCreate() (admission.Warnings, error) {
	hivelocitymachinelog.V(1).Info("validate create", "name", r.Name)
	return nil, r.Spec.DeviceSelector.Validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityMachine) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	hivelocitymachinelog.V(1).Info("validate update", "name", r.Name)
	old, ok := oldRaw.(*HivelocityMachine)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an HivelocityMachine but got a %T", oldRaw))
	}

	var allErrs field.ErrorList

	// DeviceSelector is immutable
	if !reflect.DeepEqual(old.Spec.DeviceSelector, r.Spec.DeviceSelector) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "DeviceSelector"), r.Spec.DeviceSelector, "field is immutable"),
		)
	}

	// ImageName is immutable
	if !reflect.DeepEqual(old.Spec.ImageName, r.Spec.ImageName) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "imageName"), r.Spec.ImageName, "field is immutable"),
		)
	}

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *HivelocityMachine) ValidateDelete() (admission.Warnings, error) {
	hivelocitymachinelog.V(1).Info("validate delete", "name", r.Name)
	return nil, nil
}
