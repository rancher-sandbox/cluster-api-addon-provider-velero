//go:build !ignore_autogenerated

/*
Copyright 2024.

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
	runtime "k8s.io/apimachinery/pkg/runtime"
	apiv1alpha1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupStatus) DeepCopyInto(out *BackupStatus) {
	*out = *in
	in.BackupStatus.DeepCopyInto(&out.BackupStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupStatus.
func (in *BackupStatus) DeepCopy() *BackupStatus {
	if in == nil {
		return nil
	}
	out := new(BackupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in BackupStatuses) DeepCopyInto(out *BackupStatuses) {
	{
		in := &in
		*out = make(BackupStatuses, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupStatuses.
func (in BackupStatuses) DeepCopy() BackupStatuses {
	if in == nil {
		return nil
	}
	out := new(BackupStatuses)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RestoreStatus) DeepCopyInto(out *RestoreStatus) {
	*out = *in
	in.RestoreStatus.DeepCopyInto(&out.RestoreStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreStatus.
func (in *RestoreStatus) DeepCopy() *RestoreStatus {
	if in == nil {
		return nil
	}
	out := new(RestoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in RestoreStatuses) DeepCopyInto(out *RestoreStatuses) {
	{
		in := &in
		*out = make(RestoreStatuses, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RestoreStatuses.
func (in RestoreStatuses) DeepCopy() RestoreStatuses {
	if in == nil {
		return nil
	}
	out := new(RestoreStatuses)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduleStatus) DeepCopyInto(out *ScheduleStatus) {
	*out = *in
	in.ScheduleStatus.DeepCopyInto(&out.ScheduleStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleStatus.
func (in *ScheduleStatus) DeepCopy() *ScheduleStatus {
	if in == nil {
		return nil
	}
	out := new(ScheduleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in ScheduleStatuses) DeepCopyInto(out *ScheduleStatuses) {
	{
		in := &in
		*out = make(ScheduleStatuses, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduleStatuses.
func (in ScheduleStatuses) DeepCopy() ScheduleStatuses {
	if in == nil {
		return nil
	}
	out := new(ScheduleStatuses)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroBackup) DeepCopyInto(out *VeleroBackup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroBackup.
func (in *VeleroBackup) DeepCopy() *VeleroBackup {
	if in == nil {
		return nil
	}
	out := new(VeleroBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroBackup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroBackupList) DeepCopyInto(out *VeleroBackupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VeleroBackup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroBackupList.
func (in *VeleroBackupList) DeepCopy() *VeleroBackupList {
	if in == nil {
		return nil
	}
	out := new(VeleroBackupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroBackupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroBackupSpec) DeepCopyInto(out *VeleroBackupSpec) {
	*out = *in
	in.BackupSpec.DeepCopyInto(&out.BackupSpec)
	in.Installation.DeepCopyInto(&out.Installation)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroBackupSpec.
func (in *VeleroBackupSpec) DeepCopy() *VeleroBackupSpec {
	if in == nil {
		return nil
	}
	out := new(VeleroBackupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroBackupStatus) DeepCopyInto(out *VeleroBackupStatus) {
	*out = *in
	if in.Statuses != nil {
		in, out := &in.Statuses, &out.Statuses
		*out = make(BackupStatuses, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroBackupStatus.
func (in *VeleroBackupStatus) DeepCopy() *VeleroBackupStatus {
	if in == nil {
		return nil
	}
	out := new(VeleroBackupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroInstallation) DeepCopyInto(out *VeleroInstallation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroInstallation.
func (in *VeleroInstallation) DeepCopy() *VeleroInstallation {
	if in == nil {
		return nil
	}
	out := new(VeleroInstallation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroInstallation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroInstallationList) DeepCopyInto(out *VeleroInstallationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VeleroInstallation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroInstallationList.
func (in *VeleroInstallationList) DeepCopy() *VeleroInstallationList {
	if in == nil {
		return nil
	}
	out := new(VeleroInstallationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroInstallationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroInstallationSpec) DeepCopyInto(out *VeleroInstallationSpec) {
	*out = *in
	if in.HelmChartProxySpec != nil {
		in, out := &in.HelmChartProxySpec, &out.HelmChartProxySpec
		*out = new(apiv1alpha1.HelmChartProxySpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroInstallationSpec.
func (in *VeleroInstallationSpec) DeepCopy() *VeleroInstallationSpec {
	if in == nil {
		return nil
	}
	out := new(VeleroInstallationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroInstallationStatus) DeepCopyInto(out *VeleroInstallationStatus) {
	*out = *in
	in.HelmChartProxyStatus.DeepCopyInto(&out.HelmChartProxyStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroInstallationStatus.
func (in *VeleroInstallationStatus) DeepCopy() *VeleroInstallationStatus {
	if in == nil {
		return nil
	}
	out := new(VeleroInstallationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroRestore) DeepCopyInto(out *VeleroRestore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroRestore.
func (in *VeleroRestore) DeepCopy() *VeleroRestore {
	if in == nil {
		return nil
	}
	out := new(VeleroRestore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroRestore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroRestoreList) DeepCopyInto(out *VeleroRestoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VeleroRestore, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroRestoreList.
func (in *VeleroRestoreList) DeepCopy() *VeleroRestoreList {
	if in == nil {
		return nil
	}
	out := new(VeleroRestoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroRestoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroRestoreSpec) DeepCopyInto(out *VeleroRestoreSpec) {
	*out = *in
	in.RestoreSpec.DeepCopyInto(&out.RestoreSpec)
	in.Installation.DeepCopyInto(&out.Installation)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroRestoreSpec.
func (in *VeleroRestoreSpec) DeepCopy() *VeleroRestoreSpec {
	if in == nil {
		return nil
	}
	out := new(VeleroRestoreSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroRestoreStatus) DeepCopyInto(out *VeleroRestoreStatus) {
	*out = *in
	if in.Statuses != nil {
		in, out := &in.Statuses, &out.Statuses
		*out = make(RestoreStatuses, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroRestoreStatus.
func (in *VeleroRestoreStatus) DeepCopy() *VeleroRestoreStatus {
	if in == nil {
		return nil
	}
	out := new(VeleroRestoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroSchedule) DeepCopyInto(out *VeleroSchedule) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroSchedule.
func (in *VeleroSchedule) DeepCopy() *VeleroSchedule {
	if in == nil {
		return nil
	}
	out := new(VeleroSchedule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroSchedule) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroScheduleList) DeepCopyInto(out *VeleroScheduleList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VeleroSchedule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroScheduleList.
func (in *VeleroScheduleList) DeepCopy() *VeleroScheduleList {
	if in == nil {
		return nil
	}
	out := new(VeleroScheduleList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VeleroScheduleList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroScheduleSpec) DeepCopyInto(out *VeleroScheduleSpec) {
	*out = *in
	in.ScheduleSpec.DeepCopyInto(&out.ScheduleSpec)
	in.Installation.DeepCopyInto(&out.Installation)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroScheduleSpec.
func (in *VeleroScheduleSpec) DeepCopy() *VeleroScheduleSpec {
	if in == nil {
		return nil
	}
	out := new(VeleroScheduleSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VeleroScheduleStatus) DeepCopyInto(out *VeleroScheduleStatus) {
	*out = *in
	if in.Statuses != nil {
		in, out := &in.Statuses, &out.Statuses
		*out = make(ScheduleStatuses, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VeleroScheduleStatus.
func (in *VeleroScheduleStatus) DeepCopy() *VeleroScheduleStatus {
	if in == nil {
		return nil
	}
	out := new(VeleroScheduleStatus)
	in.DeepCopyInto(out)
	return out
}
