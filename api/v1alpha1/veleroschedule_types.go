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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VeleroScheduleSpec defines the desired state of VeleroSchedule
type VeleroScheduleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of VeleroSchedule. Edit veleroschedule_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// VeleroScheduleStatus defines the observed state of VeleroSchedule
type VeleroScheduleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
