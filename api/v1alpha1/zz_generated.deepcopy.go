//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControllerGeneratedStatus) DeepCopyInto(out *ControllerGeneratedStatus) {
	*out = *in
	if in.LastUpdated != nil {
		in, out := &in.LastUpdated, &out.LastUpdated
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControllerGeneratedStatus.
func (in *ControllerGeneratedStatus) DeepCopy() *ControllerGeneratedStatus {
	if in == nil {
		return nil
	}
	out := new(ControllerGeneratedStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceSelector) DeepCopyInto(out *DeviceSelector) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.MatchExpressions != nil {
		in, out := &in.MatchExpressions, &out.MatchExpressions
		*out = make([]DeviceSelectorRequirement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceSelector.
func (in *DeviceSelector) DeepCopy() *DeviceSelector {
	if in == nil {
		return nil
	}
	out := new(DeviceSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeviceSelectorRequirement) DeepCopyInto(out *DeviceSelectorRequirement) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeviceSelectorRequirement.
func (in *DeviceSelectorRequirement) DeepCopy() *DeviceSelectorRequirement {
	if in == nil {
		return nil
	}
	out := new(DeviceSelectorRequirement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityCluster) DeepCopyInto(out *HivelocityCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityCluster.
func (in *HivelocityCluster) DeepCopy() *HivelocityCluster {
	if in == nil {
		return nil
	}
	out := new(HivelocityCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterList) DeepCopyInto(out *HivelocityClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HivelocityCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterList.
func (in *HivelocityClusterList) DeepCopy() *HivelocityClusterList {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterSpec) DeepCopyInto(out *HivelocityClusterSpec) {
	*out = *in
	if in.ControlPlaneEndpoint != nil {
		in, out := &in.ControlPlaneEndpoint, &out.ControlPlaneEndpoint
		*out = new(v1beta1.APIEndpoint)
		**out = **in
	}
	out.HivelocitySecret = in.HivelocitySecret
	if in.SSHKey != nil {
		in, out := &in.SSHKey, &out.SSHKey
		*out = new(SSHKey)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterSpec.
func (in *HivelocityClusterSpec) DeepCopy() *HivelocityClusterSpec {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterStatus) DeepCopyInto(out *HivelocityClusterStatus) {
	*out = *in
	if in.FailureDomains != nil {
		in, out := &in.FailureDomains, &out.FailureDomains
		*out = make(v1beta1.FailureDomains, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(v1beta1.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterStatus.
func (in *HivelocityClusterStatus) DeepCopy() *HivelocityClusterStatus {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterTemplate) DeepCopyInto(out *HivelocityClusterTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterTemplate.
func (in *HivelocityClusterTemplate) DeepCopy() *HivelocityClusterTemplate {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityClusterTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterTemplateList) DeepCopyInto(out *HivelocityClusterTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HivelocityClusterTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterTemplateList.
func (in *HivelocityClusterTemplateList) DeepCopy() *HivelocityClusterTemplateList {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityClusterTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterTemplateResource) DeepCopyInto(out *HivelocityClusterTemplateResource) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterTemplateResource.
func (in *HivelocityClusterTemplateResource) DeepCopy() *HivelocityClusterTemplateResource {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityClusterTemplateSpec) DeepCopyInto(out *HivelocityClusterTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterTemplateSpec.
func (in *HivelocityClusterTemplateSpec) DeepCopy() *HivelocityClusterTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachine) DeepCopyInto(out *HivelocityMachine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachine.
func (in *HivelocityMachine) DeepCopy() *HivelocityMachine {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineList) DeepCopyInto(out *HivelocityMachineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HivelocityMachine, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineList.
func (in *HivelocityMachineList) DeepCopy() *HivelocityMachineList {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityMachineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineSpec) DeepCopyInto(out *HivelocityMachineSpec) {
	*out = *in
	if in.ProviderID != nil {
		in, out := &in.ProviderID, &out.ProviderID
		*out = new(string)
		**out = **in
	}
	in.DeviceSelector.DeepCopyInto(&out.DeviceSelector)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineSpec.
func (in *HivelocityMachineSpec) DeepCopy() *HivelocityMachineSpec {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineStatus) DeepCopyInto(out *HivelocityMachineStatus) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]v1beta1.MachineAddress, len(*in))
		copy(*out, *in)
	}
	if in.FailureReason != nil {
		in, out := &in.FailureReason, &out.FailureReason
		*out = new(errors.MachineStatusError)
		**out = **in
	}
	if in.FailureMessage != nil {
		in, out := &in.FailureMessage, &out.FailureMessage
		*out = new(string)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(v1beta1.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineStatus.
func (in *HivelocityMachineStatus) DeepCopy() *HivelocityMachineStatus {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineTemplate) DeepCopyInto(out *HivelocityMachineTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineTemplate.
func (in *HivelocityMachineTemplate) DeepCopy() *HivelocityMachineTemplate {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityMachineTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineTemplateList) DeepCopyInto(out *HivelocityMachineTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HivelocityMachineTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineTemplateList.
func (in *HivelocityMachineTemplateList) DeepCopy() *HivelocityMachineTemplateList {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityMachineTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineTemplateResource) DeepCopyInto(out *HivelocityMachineTemplateResource) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineTemplateResource.
func (in *HivelocityMachineTemplateResource) DeepCopy() *HivelocityMachineTemplateResource {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineTemplateSpec) DeepCopyInto(out *HivelocityMachineTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineTemplateSpec.
func (in *HivelocityMachineTemplateSpec) DeepCopy() *HivelocityMachineTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachineTemplateStatus) DeepCopyInto(out *HivelocityMachineTemplateStatus) {
	*out = *in
	if in.Capacity != nil {
		in, out := &in.Capacity, &out.Capacity
		*out = make(v1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(v1beta1.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityMachineTemplateStatus.
func (in *HivelocityMachineTemplateStatus) DeepCopy() *HivelocityMachineTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(HivelocityMachineTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediation) DeepCopyInto(out *HivelocityRemediation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediation.
func (in *HivelocityRemediation) DeepCopy() *HivelocityRemediation {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityRemediation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationList) DeepCopyInto(out *HivelocityRemediationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HivelocityRemediation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationList.
func (in *HivelocityRemediationList) DeepCopy() *HivelocityRemediationList {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityRemediationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationSpec) DeepCopyInto(out *HivelocityRemediationSpec) {
	*out = *in
	if in.Strategy != nil {
		in, out := &in.Strategy, &out.Strategy
		*out = new(RemediationStrategy)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationSpec.
func (in *HivelocityRemediationSpec) DeepCopy() *HivelocityRemediationSpec {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationStatus) DeepCopyInto(out *HivelocityRemediationStatus) {
	*out = *in
	if in.LastRemediated != nil {
		in, out := &in.LastRemediated, &out.LastRemediated
		*out = (*in).DeepCopy()
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationStatus.
func (in *HivelocityRemediationStatus) DeepCopy() *HivelocityRemediationStatus {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationTemplate) DeepCopyInto(out *HivelocityRemediationTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationTemplate.
func (in *HivelocityRemediationTemplate) DeepCopy() *HivelocityRemediationTemplate {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityRemediationTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationTemplateList) DeepCopyInto(out *HivelocityRemediationTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HivelocityRemediationTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationTemplateList.
func (in *HivelocityRemediationTemplateList) DeepCopy() *HivelocityRemediationTemplateList {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HivelocityRemediationTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationTemplateResource) DeepCopyInto(out *HivelocityRemediationTemplateResource) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationTemplateResource.
func (in *HivelocityRemediationTemplateResource) DeepCopy() *HivelocityRemediationTemplateResource {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationTemplateSpec) DeepCopyInto(out *HivelocityRemediationTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationTemplateSpec.
func (in *HivelocityRemediationTemplateSpec) DeepCopy() *HivelocityRemediationTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityRemediationTemplateStatus) DeepCopyInto(out *HivelocityRemediationTemplateStatus) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityRemediationTemplateStatus.
func (in *HivelocityRemediationTemplateStatus) DeepCopy() *HivelocityRemediationTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(HivelocityRemediationTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocitySecretRef) DeepCopyInto(out *HivelocitySecretRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocitySecretRef.
func (in *HivelocitySecretRef) DeepCopy() *HivelocitySecretRef {
	if in == nil {
		return nil
	}
	out := new(HivelocitySecretRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RemediationStrategy) DeepCopyInto(out *RemediationStrategy) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(metav1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RemediationStrategy.
func (in *RemediationStrategy) DeepCopy() *RemediationStrategy {
	if in == nil {
		return nil
	}
	out := new(RemediationStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SSHKey) DeepCopyInto(out *SSHKey) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SSHKey.
func (in *SSHKey) DeepCopy() *SSHKey {
	if in == nil {
		return nil
	}
	out := new(SSHKey)
	in.DeepCopyInto(out)
	return out
}
