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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VeleroRestoreSpec defines the desired state of VeleroRestore
type VeleroRestoreSpec struct {
	// Restore is representing velero restore spec
	Restore velerov1.Restore `json:",inline"`

	// Installation is a helm chart proxy reference to use
	Installation clusterv1.LocalObjectTemplate `json:"installation,omitempty"`
}

// VeleroRestoreStatus defines the observed state of VeleroRestore
type VeleroRestoreStatus struct {
	// Status is representing velero Restore status
	Statuses RestoreStatuses `json:"statuses,omitempty"`
}

// Statuses is a mapping of the cluster name to Restore status
type RestoreStatuses map[ClusterName]RestoreStatus

// RestoreStatus is representing status of an individual Restore resouce
type RestoreStatus struct {
	Status velerov1.BackupStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VeleroRestore is the Schema for the velerorestores API
type VeleroRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VeleroRestoreSpec   `json:"spec,omitempty"`
	Status VeleroRestoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VeleroRestoreList contains a list of VeleroRestore
type VeleroRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VeleroRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VeleroRestore{}, &VeleroRestoreList{})
}

func (v *VeleroRestore) GetInstallRef() *corev1.ObjectReference {
	return v.Spec.Installation.Ref
}
