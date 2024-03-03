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

package plugin

import (
	"cmp"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AzurePlugin struct{}

func (p *AzurePlugin) Plugin(installation *veleroaddonv1.VeleroInstallation, provider *veleroaddonv1.Azure) {
	installation.Spec.State.InitContainers = []corev1.Container{{
		Name:            "velero-plugin-for-microsoft-azure",
		Image:           cmp.Or(provider.PluginURL, "velero/velero-plugin-for-microsoft-azure") + ":" + cmp.Or(provider.PluginTag, "latest"),
		ImagePullPolicy: corev1.PullIfNotPresent,
		VolumeMounts: []corev1.VolumeMount{{
			Name:      "plugins",
			MountPath: "/target",
		}},
	}}
}

func (p *AzurePlugin) BackupStorageLocation(location veleroaddonv1.BackupStorageLocation, provider *veleroaddonv1.Azure) veleroaddonv1.BackupStorageLocation {
	location.CredentialKey = veleroaddonv1.CredentialKey{
		Name: cmp.Or(provider.CredentialMap.To, p.Secret(provider).Name),
		Key:  veleroaddonv1.Provider{Azure: provider}.Name(),
	}
	location.Config = map[string]string{
		"resourceGroup":           provider.Config.ResourceGroup,
		"storageAccount":          provider.Config.StorageAccount,
		"storageAccountKeyEnvVar": cmp.Or(provider.Config.StorageAccountKeyEnvVar, "AZURE_STORAGE_ACCOUNT_ACCESS_KEY"),
	}

	return location
}

func (p *AzurePlugin) VolumeSnapshotLocation(snapshotLocation veleroaddonv1.VolumeSnapshotLocation, provider *veleroaddonv1.Azure) veleroaddonv1.VolumeSnapshotLocation {
	snapshotLocation.CredentialKey = veleroaddonv1.CredentialKey{
		Name: cmp.Or(provider.CredentialMap.To, p.Secret(provider).Name),
		Key:  veleroaddonv1.Provider{Azure: provider}.Name(),
	}

	return snapshotLocation
}

func (p *AzurePlugin) Secret(provider *veleroaddonv1.Azure) client.ObjectKey {
	return types.NamespacedName{
		Name:      cmp.Or(provider.CredentialMap.NamespaceName.Name, provider.CredentialMap.From),
		Namespace: provider.CredentialMap.NamespaceName.Namespace,
	}
}
