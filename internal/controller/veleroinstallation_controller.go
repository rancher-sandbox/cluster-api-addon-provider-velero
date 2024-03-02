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

package controller

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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/yaml"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	helmv1 "sigs.k8s.io/cluster-api-addon-provider-helm/api/v1alpha1"
)

// VeleroInstallationReconciler reconciles a VeleroInstallation object
type VeleroInstallationReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Tracker *remote.ClusterCacheTracker
}

//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=veleroinstallations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=veleroinstallations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=veleroinstallations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VeleroInstallation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *VeleroInstallationReconciler) Reconcile(ctx context.Context, installation *veleroaddonv1.VeleroInstallation) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	locations := installation.Spec.State.Configuration.BackupStorageLocations
	snapshotLocations := installation.Spec.State.Configuration.VolumeSnapshotLocations
	index, snapshotIndex := -1, -1
	from, fromNamespace := "", ""

	provider := installation.Spec.Provider
	location := veleroaddonv1.BackupStorageLocation{
		Provider: provider.Name(),
	}
	snapshotLocation := veleroaddonv1.VolumeSnapshotLocation{
		Provider: provider.Name(),
	}
	if index = slices.IndexFunc(locations, func(l veleroaddonv1.BackupStorageLocation) bool {
		return l.Name == ptr.To(provider.Name())
	}); index > -1 {
		location = locations[index]
	}

	if snapshotIndex = slices.IndexFunc(snapshotLocations, func(l veleroaddonv1.VolumeSnapshotLocation) bool {
		return l.Name == ptr.To(provider.Name())
	}); snapshotIndex > -1 {
		snapshotLocation = snapshotLocations[snapshotIndex]
	}
	switch {
	case provider.AWS != nil:
		location.Bucket = installation.Spec.Bucket
		location.Config = map[string]string{
			"s3Url":  provider.AWS.Config.S3Url,
			"region": provider.AWS.Config.Region,
		}

		snapshotLocation.Config = map[string]string{
			"region": provider.AWS.Config.Region,
		}

		image := cmp.Or(installation.Spec.Provider.AWS.PluginURL, "velero/velero-plugin-for-aws")
		tag := cmp.Or(installation.Spec.Provider.AWS.PluginTag, "latest")
		installation.Spec.State.InitContainers = []corev1.Container{{
			Name:            "velero-plugin-for-aws",
			Image:           image + ":" + tag,
			ImagePullPolicy: corev1.PullIfNotPresent,
			VolumeMounts: []corev1.VolumeMount{{
				Name:      "plugins",
				MountPath: "/target",
			}},
		}}

		from = cmp.Or(provider.AWS.CredentialMap.NamespaceName.Name, provider.AWS.CredentialMap.From)
		fromNamespace = cmp.Or(provider.AWS.CredentialMap.NamespaceName.Namespace, installation.Namespace)
		location.CredentialKey = veleroaddonv1.CredentialKey{
			Name: cmp.Or(provider.AWS.CredentialMap.To, from),
			Key:  provider.Name(),
		}
		snapshotLocation.CredentialKey = veleroaddonv1.CredentialKey{
			Name: cmp.Or(provider.AWS.CredentialMap.To, from),
			Key:  provider.Name(),
		}

	case provider.Azure != nil:
		location.Bucket = installation.Spec.Bucket
		location.Config = map[string]string{
			"resourceGroup":           provider.Azure.Config.ResourceGroup,
			"storageAccount":          provider.Azure.Config.StorageAccount,
			"storageAccountKeyEnvVar": cmp.Or(provider.Azure.Config.StorageAccountKeyEnvVar, "AZURE_STORAGE_ACCOUNT_ACCESS_KEY"),
		}

		snapshotLocation.Config = map[string]string{}

		image := cmp.Or(installation.Spec.Provider.Azure.PluginURL, "velero/velero-plugin-for-microsoft-azure")
		tag := cmp.Or(installation.Spec.Provider.Azure.PluginTag, "latest")
		installation.Spec.State.InitContainers = []corev1.Container{{
			Name:            "velero-plugin-for-microsoft-azure",
			Image:           image + ":" + tag,
			ImagePullPolicy: corev1.PullIfNotPresent,
			VolumeMounts: []corev1.VolumeMount{{
				Name:      "plugins",
				MountPath: "/target",
			}},
		}}

		from = cmp.Or(provider.Azure.CredentialMap.NamespaceName.Name, provider.Azure.CredentialMap.From)
		fromNamespace = cmp.Or(provider.Azure.CredentialMap.NamespaceName.Namespace, installation.Namespace)
		location.CredentialKey = veleroaddonv1.CredentialKey{
			Name: cmp.Or(provider.Azure.CredentialMap.To, from),
			Key:  provider.Name(),
		}
		snapshotLocation.CredentialKey = veleroaddonv1.CredentialKey{
			Name: cmp.Or(provider.Azure.CredentialMap.To, from),
			Key:  provider.Name(),
		}
	}

	// Plugins / values
	if index > -1 {
		locations[index] = location
	} else {
		locations = append(locations, location)
	}

	if snapshotIndex > -1 {
		snapshotLocations[snapshotIndex] = snapshotLocation
	} else {
		snapshotLocations = append(snapshotLocations, snapshotLocation)
	}

	installation.Spec.State.Configuration.BackupStorageLocations = locations
	installation.Spec.State.Configuration.VolumeSnapshotLocations = snapshotLocations

	// Secret sync
	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name:      from,
		Namespace: fromNamespace,
	}, secret); err != nil {
		return ctrl.Result{}, err
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
				Name:      location.CredentialKey.Name,
				Namespace: cmp.Or(installation.Spec.HelmSpec.ReleaseNamespace, installation.Spec.Namespace, "velero"),
			},
			Data: secret.Data,
		}

		errs = append(errs, cl.Patch(ctx, newSecret, client.Apply, client.ForceOwnership, client.FieldOwner("velero-addon")))
	}

	if err := kerrors.NewAggregate(errs); err != nil {
		return ctrl.Result{}, err
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

// SetupWithManager sets up the controller with the Manager.
func (r *VeleroInstallationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&veleroaddonv1.VeleroInstallation{}).
		Owns(&helmv1.HelmChartProxy{}).
		Complete(reconcile.AsReconciler(r.Client, r))
}
