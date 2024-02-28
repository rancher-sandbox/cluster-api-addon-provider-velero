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

	"k8s.io/apimachinery/pkg/runtime"
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
	Scheme *runtime.Scheme
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

	for _, plugin := range installation.Spec.State.Plugins {
		switch plugin {
		case veleroaddonv1.AWS:
			installation.Spec.State.InitContainers = []corev1.Container{{
				Name:            "velero-plugin-for-aws",
				Image:           "velero/velero-plugin-for-aws:v1.9.0",
				ImagePullPolicy: corev1.PullIfNotPresent,
				VolumeMounts: []corev1.VolumeMount{{
					Name:      "plugins",
					MountPath: "/target",
				}},
			}}
		}

	}

	spec, err := yaml.Marshal(installation.Spec.State)
	if err != nil {
		return ctrl.Result{}, err
	}

	helmSpec := installation.Spec.HelmChartProxySpec
	helmSpec.ValuesTemplate = cmp.Or(helmSpec.ValuesTemplate, string(spec))

	helmProxy := templateHelmChartProxy(installation, helmSpec)
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

func templateHelmChartProxy(installation *veleroaddonv1.VeleroInstallation, helmSpec helmv1.HelmChartProxySpec) *helmv1.HelmChartProxy {
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
			ClusterSelector: metav1.LabelSelector{},
			RepoURL:         helmSpec.RepoURL,
			ChartName:       helmSpec.ChartName,
			ValuesTemplate:  helmSpec.ValuesTemplate,
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
