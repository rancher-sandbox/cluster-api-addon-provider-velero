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

package v1alpha1

import (
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VeleroScheduleSpec defines the desired state of VeleroSchedule
type VeleroScheduleSpec struct {
	// Schedule is representing velero Schedule spec
	velerov1.ScheduleSpec `json:",inline"`

	// Installation is a helm chart proxy reference to use
	Installation clusterv1.LocalObjectTemplate `json:"installation,omitempty"`
}

// VeleroScheduleStatus defines the observed state of Velero Schedule resource
type VeleroScheduleStatus struct {
	// Status is representing velero backup status
	Statuses ScheduleStatuses `json:"statuses,omitempty"`
}

// ScheduleStatuses is a mapping of the cluster name to schedule status
type ScheduleStatuses map[ClusterName]ScheduleStatus

// ScheduleStatus is representing status of an individual Schedule resouce
type ScheduleStatus struct {
	velerov1.ScheduleStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VeleroSchedule is the Schema for the veleroschedules API
type VeleroSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VeleroScheduleSpec   `json:"spec,omitempty"`
	Status VeleroScheduleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VeleroScheduleList contains a list of VeleroSchedule
type VeleroScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VeleroSchedule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VeleroSchedule{}, &VeleroScheduleList{})
}

func (v *VeleroSchedule) GetInstallRef() *corev1.ObjectReference {
	return v.Spec.Installation.Ref
}

func (v *VeleroSchedule) SetClusterStatus(key ClusterName, schedule *velerov1.Schedule) {
	if v.Status.Statuses == nil {
		v.Status.Statuses = ScheduleStatuses{}
	}
	v.Status.Statuses[key] = ScheduleStatus{
		ScheduleStatus: schedule.Status,
	}
}
