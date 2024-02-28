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

// VeleroBackupSpec defines the desired state of VeleroBackup
type VeleroBackupSpec struct {
	// Backup is representing velero backup spec
	velerov1.BackupSpec `json:",inline"`

	// Installation is a helm chart proxy reference to use
	Installation clusterv1.LocalObjectTemplate `json:"installation,omitempty"`
}

// VeleroBackupStatus defines the observed state of VeleroBackup
type VeleroBackupStatus struct {
	// Status is representing velero backup status
	Statuses BackupStatuses `json:"statuses,omitempty"`
}

// BackupStatuses is a mapping of the cluster name to backup status
type BackupStatuses map[NamespaceName]BackupStatus

// BackupStatus is representing status of an individual Backup resource
type BackupStatus struct {
	velerov1.BackupStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VeleroBackup is the Schema for the velerobackups API
type VeleroBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VeleroBackupSpec   `json:"spec,omitempty"`
	Status VeleroBackupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VeleroBackupList contains a list of VeleroBackup
type VeleroBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VeleroBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VeleroBackup{}, &VeleroBackupList{})
}

func (v *VeleroBackup) GetInstallRef() *corev1.ObjectReference {
	return v.Spec.Installation.Ref
}

func (v *VeleroBackup) SetClusterStatus(key NamespaceName, backup *velerov1.Backup) {
	if v.Status.Statuses == nil {
		v.Status.Statuses = BackupStatuses{}
	}
	v.Status.Statuses[key] = BackupStatus{
		BackupStatus: backup.Status,
	}
}
