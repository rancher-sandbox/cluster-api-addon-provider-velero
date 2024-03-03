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

type GCPPlugin struct{}

func (p *GCPPlugin) Plugin(installation *veleroaddonv1.VeleroInstallation, provider *veleroaddonv1.GCP) {
	installation.Spec.State.InitContainers = []corev1.Container{{
		Name:            "velero-plugin-for-gcp",
		Image:           cmp.Or(provider.PluginURL, "velero/velero-plugin-for-gcp") + ":" + cmp.Or(provider.PluginTag, "latest"),
		ImagePullPolicy: corev1.PullIfNotPresent,
		VolumeMounts: []corev1.VolumeMount{{
			Name:      "plugins",
			MountPath: "/target",
		}},
	}}
}

func (p *GCPPlugin) BackupStorageLocation(location veleroaddonv1.BackupStorageLocation, provider *veleroaddonv1.GCP) veleroaddonv1.BackupStorageLocation {
	location.Config = map[string]string{
		"serviceAccount": provider.Config.ServiceAccount,
		"kmsKeyName":     provider.Config.KMSKeyName,
	}
	location.CredentialKey = veleroaddonv1.CredentialKey{
		Name: cmp.Or(provider.CredentialMap.To, p.Secret(provider).Name),
		Key:  "gcp",
	}

	return location
}

func (p *GCPPlugin) VolumeSnapshotLocation(snapshotLocation veleroaddonv1.VolumeSnapshotLocation, provider *veleroaddonv1.GCP) veleroaddonv1.VolumeSnapshotLocation {
	snapshotLocation.Config = map[string]string{
		"snapshotLocation": provider.Config.SnapshotLocation,
		"project":          provider.Config.Project,
	}
	snapshotLocation.CredentialKey = veleroaddonv1.CredentialKey{
		Name: cmp.Or(provider.CredentialMap.To, p.Secret(provider).Name),
		Key:  "gcp",
	}

	return snapshotLocation
}

func (p *GCPPlugin) Secret(provider *veleroaddonv1.GCP) client.ObjectKey {
	return types.NamespacedName{
		Name:      cmp.Or(provider.CredentialMap.NamespaceName.Name, provider.CredentialMap.From),
		Namespace: provider.CredentialMap.NamespaceName.Namespace,
	}
}
