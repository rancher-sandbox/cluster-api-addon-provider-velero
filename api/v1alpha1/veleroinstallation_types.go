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
	helmv1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// VeleroInstallationSpec defines the desired state of VeleroInstallation
type VeleroInstallationSpec struct {
	// HelmSpec is a Helm chart proxy installation spec
	// +optional
	HelmSpec helmv1.HelmChartProxySpec `json:"helmSpec,omitempty"`

	// ClusterSelector selects Clusters in the same namespace with a label that matches the specified label selector. The Helm
	// chart will be installed on all selected Clusters. If a Cluster is no longer selected, the Helm release will be uninstalled.
	// +optional
	ClusterSelector metav1.LabelSelector `json:"clusterSelector,omitempty"`

	// +optional
	Namespace string `json:"namespace,omitempty"`

	State VeleroHelmState `json:"state,omitempty"`

	Bucket string `json:"bucket"`

	Provider Provider `json:"provider,omitempty"`
}

type Provider struct {
	AWS   *AWS   `json:"aws,omitempty"`
	Azure *Azure `json:"azure,omitempty"`
	GCP   *GCP   `json:"gcp,omitempty"`
}

type AWS struct {
	// +optional
	PluginURL string `json:"pluginURL"`

	// +optional
	PluginTag string `json:"pluginTag"`

	CredentialMap CredentialMap `json:"credentialMap,omitempty"`

	// +optional
	Config AWSConfig `json:"config,omitempty"`
}

type Azure struct {
	// +optional
	PluginURL string `json:"pluginURL"`

	// +optional
	PluginTag string `json:"pluginTag"`

	CredentialMap CredentialMap `json:"credentialMap,omitempty"`

	// +optional
	Config AzureConfig `json:"config,omitempty"`
}

type GCP struct {
	// +optional
	PluginURL string `json:"pluginURL"`

	// +optional
	PluginTag string `json:"pluginTag"`

	CredentialMap CredentialMap `json:"credentialMap,omitempty"`

	Config GCPConfig `json:"config,omitempty"`
}

type AWSConfig struct {
	// +optional
	Region string `json:"region,omitempty"`

	// +optional
	S3Url string `json:"s3Url,omitempty"`
}

type AzureConfig struct {
	// AZURE_BACKUP_RESOURCE_GROUP
	ResourceGroup string `json:"resourceGroup"`

	// AZURE_STORAGE_ACCOUNT_ID
	StorageAccount string `json:"storageAccount"`

	// AZURE_STORAGE_ACCOUNT_ACCESS_KEY
	// +optional
	StorageAccountKeyEnvVar string `json:"storageAccountKeyEnvVar"`

	// AZURE_BACKUP_SUBSCRIPTION_ID
	// +optional
	SubscriptionId string `json:"subscriptionId"`
}

type GCPConfig struct {
	// Name of the GCP service account to use for this backup storage location. Specify the
	// service account here if you want to use workload identity instead of providing the key file.
	//
	// Optional (defaults to "false").
	// +optional
	ServiceAccount string `json:"serviceAccount"`

	// Name of the Cloud KMS key to use to encrypt backups stored in this location, in the form
	// "projects/P/locations/L/keyRings/R/cryptoKeys/K". See customer-managed Cloud KMS keys
	// (https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for details.
	// +optional
	KMSKeyName string `json:"kmsKeyName"`

	// The GCP location where snapshots should be stored. See the GCP documentation
	// (https://cloud.google.com/storage/docs/locations#available_locations) for the
	// full list. If not specified, snapshots are stored in the default location
	// (https://cloud.google.com/compute/docs/disks/create-snapshots#default_location).
	//
	// Example: us-central1
	// +optional
	SnapshotLocation string `json:"snapshotLocation,omitempty"`

	// The project ID where existing snapshots should be retrieved from during restores, if
	// different than the project that your IAM account is in. This field has no effect on
	// where new snapshots are created; it is only useful for restoring existing snapshots
	// from a different project.
	//
	// Optional (defaults to the project that the GCP IAM account is in).
	// Example: my-alternate-project
	Project string `json:"project,omitempty"`
}

type VeleroHelmState struct {
	DeployNodeAgent bool `json:"deployNodeAgent"`
	CleanUpCRDs     bool `json:"cleanUpCRDs"`

	// Configuration is a bucket configuration
	// +optional
	Configuration Configuration `json:"configuration,omitempty"`

	// Info about the secret to be used by the Velero deployment, which
	// should contain credentials for the cloud provider IAM account you've
	// set up for Velero.
	// +optional
	Credentials Credentials `json:"credentials,omitempty"`

	//+optional
	InitContainers []corev1.Container `json:"initContainers,omitempty"`
}

type Configuration struct {
	BackupStorageLocations  []BackupStorageLocation  `json:"backupStorageLocation"`
	VolumeSnapshotLocations []VolumeSnapshotLocation `json:"volumeSnapshotLocation"`
}

type VolumeSnapshotLocation struct {
	// Name of this backup storage location. If unspecified, use "default".
	// +optional
	Name *string `json:"name,omitempty"`

	// The name for the backup storage provider.
	Provider string `json:"provider"`

	CredentialKey CredentialKey `json:"credential,omitempty"`

	// Config containe additional provider-specific configuration. See link above
	// for details of required/optional fields for your provider.
	Config map[string]string `json:"config,omitempty"`
}

type BackupStorageLocation struct {
	// Name of this backup storage location. If unspecified, use "default".
	// +optional
	Name *string `json:"name,omitempty"`

	// The name for the backup storage provider.
	Provider string `json:"provider"`

	// The name or ID of the bucket to store backups in. Required.
	Bucket string `json:"bucket"`

	// Base64 encoded CA bundle used when verifying TLS connections to the provider.
	// +optional
	CAcert *string `json:"caCert,omitempty"`

	// Directory under which all Velero data will be stored within the bucket. Optional.
	// +optional
	Prefix *string `json:"prefix,omitempty"`

	// Flag to indicate if this is the default backup storage location (used as fallback if no other location is specified). Optional.
	Default bool `json:"default,omitempty"`

	// Frequency at which Velero should perform validation checks on this location. Optional.
	ValidationFrequency int64 `json:"validationFrequency,omitempty"`

	// Access mode for this backup storage location. Defaults to ReadWrite.
	AccessMode AccessMode `json:"accessMode,omitempty"`

	CredentialKey CredentialKey `json:"credential,omitempty"`

	// Config containe additional provider-specific configuration. See link above
	// for details of required/optional fields for your provider.
	Config map[string]string `json:"config,omitempty"`
}

type CredentialKey struct {
	// Name of the secret used by this backupStorageLocation.
	Name string `json:"name,omitempty"`

	// Key that contains the secret data to be used.
	// +optional
	Key string `json:"key,omitempty"`
}

type CredentialMap struct {
	// +optional
	From string `json:"from,omitempty"`

	// +optional
	NamespaceName CredentialNamespaceName `json:"namespaceName,omitempty"`

	// +optional
	To string `json:"to,omitempty"`
}

type CredentialNamespaceName struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type AccessMode string

const (
	ReadWrite AccessMode = "ReadWrite"
	ReadOnly  AccessMode = "ReadOnly"
)

type Credentials struct {
	// Set to false if not using a secret for credentials (i.e., use KIAM or WID)
	UseSecret bool `json:"useSecret,omitempty"`

	// If set, name of pre-existing Velero secret to be used in case of 'useSecret=true' and empty 'existingSecret'.
	Name string `json:"name,omitempty"`

	// Pre-existing secret name in the Velero namespace, if any.
	ExistingSecret *string `json:"existingSecret,omitempty"`

	// Map storing secret contents (key: "<cloud>", value: entire IAM credential file).
	Contents map[string]string `json:"contents,omitempty"`

	// Environment variables from the secret to be loaded into velero/node-agent.
	ExtraEnvVars map[string]string `json:"extraEnvVars,omitempty"`

	// Reference to existing secrets with environment variable format.
	ExtraSecretRef *string `json:"extraSecretRef,omitempty"`
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
