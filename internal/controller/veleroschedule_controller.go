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

// VeleroScheduleReconciler reconciles a VeleroSchedule object
type VeleroScheduleReconciler struct {
	Reconciler[*veleroaddonv1.VeleroSchedule, *velerov1.Schedule]
	Scheme   *runtime.Scheme
	Schedule *velerov1.Schedule
}

// SetupWithManager sets up the controller with the Manager.
func (r *VeleroScheduleReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return r.Reconciler.SetupWithManager(ctx, mgr, options).Complete(
		reconcile.AsReconciler(r.Client, AsVeleroReconciler(r.Client, r)))
}

//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=veleroschedules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=veleroschedules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=addons.cluster.x-k8s.io,resources=veleroschedules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VeleroSchedule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *VeleroScheduleReconciler) ReconcileProxy(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, schedule *veleroaddonv1.VeleroSchedule) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	r.Installation = installation
	r.Schedule = &velerov1.Schedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      schedule.Name,
			Namespace: "default",
		},
		Spec: schedule.Spec.ScheduleSpec,
	}
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

func (r *VeleroScheduleReconciler) GetObject() client.Object {
	return &veleroaddonv1.VeleroSchedule{}
}

func (r *VeleroScheduleReconciler) UpdateRemote(ctx context.Context, schedule *veleroaddonv1.VeleroSchedule) error {
	return r.Reconciler.UpdateRemote(ctx, r.Installation, schedule, r.Schedule)
}
