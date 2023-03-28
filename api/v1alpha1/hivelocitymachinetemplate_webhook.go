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
	"context"
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api/util/topology"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var hivelocitymachinetemplatelog = logf.Log.WithName("hivelocitymachinetemplate-resource")

func (r *HivelocityMachineTemplateWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&HivelocityMachineTemplate{}).
		WithValidator(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocitymachinetemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachinetemplates,verbs=create;update,versions=v1alpha1,name=mhivelocitymachinetemplate.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HivelocityMachineTemplate{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *HivelocityMachineTemplate) Default() {}

// HivelocityMachineTemplateWebhook implements a custom validation webhook for HivelocityMachineTemplate.
// +kubebuilder:object:generate=false
type HivelocityMachineTemplateWebhook struct{}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-hivelocitymachinetemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=hivelocitymachinetemplates,verbs=create;update,versions=v1alpha1,name=vhivelocitymachinetemplate.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &HivelocityMachineTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *HivelocityMachineTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) error {
	newHivelocityMachine, ok := obj.(*HivelocityMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a HivelocityMachineTemplate but got a %T", obj))
	}

	hivelocitymachinetemplatelog.V(1).Info("validate create", "name", newHivelocityMachine)
	return nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *HivelocityMachineTemplateWebhook) ValidateUpdate(ctx context.Context, oldRaw runtime.Object, newRaw runtime.Object) error {
	newHivelocityMachineTemplate, ok := newRaw.(*HivelocityMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a HivelocityMachineTemplate but got a %T", newRaw))
	}

	hivelocitymachinetemplatelog.V(1).Info("validate update", "name", newHivelocityMachineTemplate.Name)

	oldHivelocityMachineTemplate, ok := oldRaw.(*HivelocityMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a HivelocityMachineTemplate but got a %T", oldRaw))
	}

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	var allErrs field.ErrorList

	if !topology.ShouldSkipImmutabilityChecks(req, newHivelocityMachineTemplate) && !reflect.DeepEqual(newHivelocityMachineTemplate.Spec, oldHivelocityMachineTemplate.Spec) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec"), newHivelocityMachineTemplate, "HivelocityMachineTemplate.Spec is immutable"))
	}

	return aggregateObjErrors(newHivelocityMachineTemplate.GroupVersionKind().GroupKind(), newHivelocityMachineTemplate.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *HivelocityMachineTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) error {
	return nil
}
