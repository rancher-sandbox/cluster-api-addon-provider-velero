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
	"context"
	"slices"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api/controllers/remote"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helmv1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
)

type PluginReconciler[P Plugin] struct {
	client.Client
	Scheme  *runtime.Scheme
	Tracker *remote.ClusterCacheTracker
}

func (r *PluginReconciler[P]) Reconcile(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, veleroPlugin VeleroPlugin[P], provider P) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	locations := installation.Spec.State.Configuration.BackupStorageLocations
	snapshotLocations := installation.Spec.State.Configuration.VolumeSnapshotLocations
	index, snapshotIndex := -1, -1

	prov := installation.Spec.Provider
	location := veleroaddonv1.BackupStorageLocation{
		Provider: prov.Name(),
		Bucket:   installation.Spec.Bucket,
		Prefix:   ptr.To("{{ .Cluster.metadata.name }}"),
	}
	snapshotLocation := veleroaddonv1.VolumeSnapshotLocation{
		Provider: prov.Name(),
	}
	if index = slices.IndexFunc(locations, func(l veleroaddonv1.BackupStorageLocation) bool {
		return l.Provider == prov.Name()
	}); index > -1 {
		location = locations[index]
	} else {
		index = len(locations)
		locations = append(locations, location)
	}

	if snapshotIndex = slices.IndexFunc(snapshotLocations, func(l veleroaddonv1.VolumeSnapshotLocation) bool {
		return l.Provider == prov.Name()
	}); snapshotIndex > -1 {
		snapshotLocation = snapshotLocations[snapshotIndex]
	} else {
		snapshotIndex = len(snapshotLocations)
		snapshotLocations = append(snapshotLocations, snapshotLocation)
	}

	veleroPlugin.Plugin(installation, provider)
	locations[index] = veleroPlugin.BackupStorageLocation(location, provider)
	snapshotLocations[snapshotIndex] = veleroPlugin.VolumeSnapshotLocation(snapshotLocation, provider)

	installation.Spec.State.Configuration.VolumeSnapshotLocations = snapshotLocations
	installation.Spec.State.Configuration.BackupStorageLocations = locations

	from := veleroPlugin.Secret(provider)
	if from.Name != "" {
		from.Namespace = cmp.Or(from.Namespace, installation.Namespace)
		to := types.NamespacedName{
			Name:      locations[index].CredentialKey.Name,
			Namespace: cmp.Or(installation.Spec.HelmSpec.ReleaseNamespace, installation.Spec.Namespace, "velero"),
		}
		if err := r.syncSecret(ctx, installation, from, to); err != nil {
			return ctrl.Result{}, err
		}
	}

	spec, err := yaml.Marshal(installation.Spec.State)
	if err != nil {
		return ctrl.Result{}, err
	}

	helmProxy := templateHelmChartProxy(installation, installation.Spec.HelmSpec, spec)
	if err := controllerutil.SetOwnerReference(installation, helmProxy, r.Client.Scheme()); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Client.Patch(
		ctx, helmProxy,
		client.Apply, client.ForceOwnership, client.FieldOwner("velero-addon")); err != nil {
		return ctrl.Result{}, err
	}

	installation.Status.HelmChartProxyStatus = helmProxy.Status

	return ctrl.Result{}, r.Client.Status().Update(ctx, installation)
}

func (r *PluginReconciler[P]) syncSecret(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, from, to client.ObjectKey) error {
	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, from, secret); err != nil {
		return err
	}

	var errs []error
	for _, cluster := range installation.Status.MatchingClusters {
		cl, err := r.Tracker.GetClient(ctx, veleroaddonv1.RefToNamespaceName(&cluster).ObjectKey())
		if err != nil {
			errs = append(errs, client.IgnoreNotFound(err))

			continue
		}

		newSecret := &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Secret",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      to.Name,
				Namespace: to.Namespace,
			},
			Data: secret.Data,
		}

		errs = append(errs, cl.Patch(ctx, newSecret, client.Apply, client.ForceOwnership, client.FieldOwner("velero-addon")))
	}

	return kerrors.NewAggregate(errs)
}

func templateHelmChartProxy(installation *veleroaddonv1.VeleroInstallation, helmSpec helmv1.HelmChartProxySpec, values []byte) *helmv1.HelmChartProxy {
	clusterSelector := helmSpec.ClusterSelector
	if installation.Spec.ClusterSelector.MatchExpressions != nil || installation.Spec.ClusterSelector.MatchLabels != nil {
		clusterSelector = installation.Spec.ClusterSelector
	}

	options := cmp.Or(helmSpec.Options, helmv1.HelmOptions{
		Install: helmv1.HelmInstallOptions{
			CreateNamespace: true,
		},
	})

	return &helmv1.HelmChartProxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: helmv1.GroupVersion.String(),
			Kind:       "HelmChartProxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      installation.Name,
			Namespace: installation.Namespace,
		},
		Spec: helmv1.HelmChartProxySpec{
			ClusterSelector:  clusterSelector,
			ReleaseNamespace: cmp.Or(installation.Spec.Namespace, "velero"),
			RepoURL:          cmp.Or(helmSpec.RepoURL, "https://vmware-tanzu.github.io/helm-charts"),
			ChartName:        cmp.Or(helmSpec.ChartName, "velero"),
			ValuesTemplate:   cmp.Or(helmSpec.ValuesTemplate, string(values)),
			Options:          options,
		},
	}
}
