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
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityCluster) DeepCopyInto(out *HivelocityCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
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
	out.Spec = in.Spec
	out.Status = in.Status
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
func (in *HivelocityClusterTemplateSpec) DeepCopyInto(out *HivelocityClusterTemplateSpec) {
	*out = *in
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
func (in *HivelocityClusterTemplateStatus) DeepCopyInto(out *HivelocityClusterTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HivelocityClusterTemplateStatus.
func (in *HivelocityClusterTemplateStatus) DeepCopy() *HivelocityClusterTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(HivelocityClusterTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HivelocityMachine) DeepCopyInto(out *HivelocityMachine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
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
	out.Spec = in.Spec
	out.Status = in.Status
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
func (in *HivelocityMachineTemplateSpec) DeepCopyInto(out *HivelocityMachineTemplateSpec) {
	*out = *in
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
	out.Spec = in.Spec
	out.Status = in.Status
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
	out.Spec = in.Spec
	out.Status = in.Status
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
func (in *HivelocityRemediationTemplateSpec) DeepCopyInto(out *HivelocityRemediationTemplateSpec) {
	*out = *in
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
