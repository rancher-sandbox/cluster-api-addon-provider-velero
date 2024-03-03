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
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/controllers/remote"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
	"addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/internal/plugin"
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

	switch {
	case installation.Spec.Provider.AWS != nil:
		return (&plugin.PluginReconciler[*veleroaddonv1.AWS]{
			Client:  r.Client,
			Scheme:  r.Scheme,
			Tracker: r.Tracker,
		}).Reconcile(ctx, installation, &plugin.AWSPlugin{}, installation.Spec.Provider.AWS)
	case installation.Spec.Provider.Azure != nil:
		return (&plugin.PluginReconciler[*veleroaddonv1.Azure]{
			Client:  r.Client,
			Scheme:  r.Scheme,
			Tracker: r.Tracker,
		}).Reconcile(ctx, installation, &plugin.AzurePlugin{}, installation.Spec.Provider.Azure)
	case installation.Spec.Provider.GCP != nil:
		return (&plugin.PluginReconciler[*veleroaddonv1.GCP]{
			Client:  r.Client,
			Scheme:  r.Scheme,
			Tracker: r.Tracker,
		}).Reconcile(ctx, installation, &plugin.GCPPlugin{}, installation.Spec.Provider.GCP)
	default:
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *VeleroInstallationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&veleroaddonv1.VeleroInstallation{}).
		Owns(&helmv1.HelmChartProxy{}).
		Complete(reconcile.AsReconciler(r.Client, r))
}
