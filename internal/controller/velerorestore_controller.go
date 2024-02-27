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

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// VeleroRestoreReconciler reconciles a VeleroRestore object
type VeleroRestoreReconciler struct {
	Reconciler[*veleroaddonv1.VeleroRestore, *velerov1.Restore]
	Scheme  *runtime.Scheme
	Restore *velerov1.Restore
}

// SetupWithManager sets up the controller with the Manager.
func (r *VeleroRestoreReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return r.Reconciler.SetupWithManager(ctx, mgr, options).Complete(
		reconcile.AsReconciler(r.Client, AsVeleroReconciler(r.Client, r)))
}

//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=velerorestores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=velerorestores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=velerorestores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VeleroRestore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *VeleroRestoreReconciler) ReconcileProxy(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, restore *veleroaddonv1.VeleroRestore) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	r.Installation = installation
	r.Restore = &velerov1.Restore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      restore.Name,
			Namespace: "default",
		},
		Spec: restore.Spec.RestoreSpec,
	}
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *VeleroRestoreReconciler) GetObject() client.Object {
	return &veleroaddonv1.VeleroRestore{}
}

func (r *VeleroRestoreReconciler) UpdateRemote(ctx context.Context, restore *veleroaddonv1.VeleroRestore) error {
	return r.Reconciler.UpdateRemote(ctx, r.Installation, restore, r.Restore)
}
