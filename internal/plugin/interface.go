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
	"sigs.k8s.io/controller-runtime/pkg/client"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
)

type Plugin interface {
	*veleroaddonv1.AWS | *veleroaddonv1.Azure | *veleroaddonv1.GCP
}

type VeleroPlugin[P Plugin] interface {
	Plugin(installation *veleroaddonv1.VeleroInstallation, provider P)
	BackupStorageLocation(location veleroaddonv1.BackupStorageLocation, provider P) veleroaddonv1.BackupStorageLocation
	VolumeSnapshotLocation(snapshotLocation veleroaddonv1.VolumeSnapshotLocation, provider P) veleroaddonv1.VolumeSnapshotLocation
	Secret(provider P) client.ObjectKey
}
