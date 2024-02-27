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
	helmv1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VeleroInstallationSpec defines the desired state of VeleroInstallation
type VeleroInstallationSpec struct {
	// Proxy is a Helm chart proxy installation
	// +optional
	*helmv1.HelmChartProxySpec `json:",inline"`
}

// VeleroInstallationStatus defines the observed state of VeleroInstallation
type VeleroInstallationStatus struct {
	helmv1.HelmChartProxyStatus `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VeleroInstallation is the Schema for the veleroinstallations API
type VeleroInstallation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VeleroInstallationSpec   `json:"spec,omitempty"`
	Status VeleroInstallationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VeleroInstallationList contains a list of VeleroInstallation
type VeleroInstallationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VeleroInstallation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VeleroInstallation{}, &VeleroInstallationList{})
}
