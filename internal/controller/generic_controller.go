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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/cluster-api/controllers/remote"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
)

// GenericReconciler is a generic interface for velero objects reconciler
type GenericReconciler[T veleroaddonv1.InstalledObject] interface {
	client.Client
	GetTracker() *remote.ClusterCacheTracker
}

type Reconciler[T veleroaddonv1.InstalledObject] struct {
	client.Client
	Tracker      *remote.ClusterCacheTracker
	Installation *veleroaddonv1.VeleroInstallation
}

type VeleroReconciler[T veleroaddonv1.InstalledObject] interface {
	UpdateRemote(ctx context.Context) error
	ReconcileProxy(context.Context, *veleroaddonv1.VeleroInstallation, T) (reconcile.Result, error)
}

// AsVeleroReconciler creates a Reconciler based on the given ObjectReconciler.
func AsVeleroReconciler[T veleroaddonv1.InstalledObject](client client.Client, rec VeleroReconciler[T]) reconcile.ObjectReconciler[T] {
	return &Adapter[T]{
		objReconciler: rec,
		client:        client,
	}
}

type Adapter[T veleroaddonv1.InstalledObject] struct {
	objReconciler VeleroReconciler[T]
	client        client.Client
}

// ReconcileProxy implements VeleroReconciler.
func (a *Adapter[T]) ReconcileProxy(ctx context.Context, i *veleroaddonv1.VeleroInstallation, o T) (reconcile.Result, error) {
	if res, err := a.objReconciler.ReconcileProxy(ctx, i, o); err != nil || res.Requeue || res.RequeueAfter > 0 {
		return res, err
	}

	return reconcile.Result{}, a.objReconciler.UpdateRemote(ctx)
}

// Reconcile implements Reconciler.
func (a *Adapter[T]) Reconcile(ctx context.Context, o T) (reconcile.Result, error) {
	ref := o.GetInstallRef()
	if ref == nil {
		return reconcile.Result{}, nil
	}

	i := &veleroaddonv1.VeleroInstallation{}
	if err := a.client.Get(ctx, types.NamespacedName{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}, i); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	return a.ReconcileProxy(ctx, i, o)
}

func (r *Reconciler[T]) GetTracker() *remote.ClusterCacheTracker {
	return r.Tracker
}

func (r *Reconciler[T]) UpdateRemote(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, obj client.Object) error {
	var errors []error

	for _, clusterRef := range installation.Status.MatchingClusters {
		clusterKey := types.NamespacedName{Name: clusterRef.Name, Namespace: clusterRef.Namespace}
		if cl, err := r.GetTracker().GetClient(ctx, clusterKey); err != nil {
			errors = append(errors, err)
		} else if err := cl.Get(ctx, client.ObjectKeyFromObject(obj), obj); apierrors.IsNotFound(err) {
			errors = append(errors, cl.Create(ctx, obj))
		} else if err != nil {
			errors = append(errors, err)
		} else {
			errors = append(errors, cl.Update(ctx, obj))
		}
	}

	return kerrors.NewAggregate(errors)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler[T]) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) *builder.Builder {
	return ctrl.NewControllerManagedBy(mgr).
		For(reflect.New(reflect.TypeOf(*new(T)).Elem()).Interface().(T), builder.WithPredicates(
			predicate.NewPredicateFuncs(r.onlyWithInstallation(ctx)))).
		WithOptions(options)
}

func (r *Reconciler[T]) onlyWithInstallation(ctx context.Context) func(client.Object) bool {
	return func(obj client.Object) bool {
		ref := obj.(T).GetInstallRef()
		installation := reflect.New(reflect.TypeOf(*new(T)).Elem()).Interface().(T)

		if ref == nil {
			return false
		}

		return r.Get(ctx, types.NamespacedName{Name: ref.Name, Namespace: ref.Namespace}, installation) != nil
	}
}
