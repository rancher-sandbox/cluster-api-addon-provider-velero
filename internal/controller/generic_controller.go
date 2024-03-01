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
	"errors"
	"reflect"
	"slices"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"sigs.k8s.io/cluster-api/controllers/remote"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	veleroaddonv1 "addons.cluster.x-k8s.io/cluster-api-addon-provider-velero/api/v1alpha1"
)

const (
	proxyKeyAnnotation = "addons.cluster.x-k8s.io/velero-proxy"
	finalizer          = "addons.cluster.x-k8s.io/velero"
)

// GenericReconciler is a generic interface for velero objects reconciler
type GenericReconciler[P veleroaddonv1.VeleroProxy[V], V veleroaddonv1.VeleroOrigin] interface {
	client.Client
	GetTracker() *remote.ClusterCacheTracker
}

type Reconciler[P veleroaddonv1.VeleroProxy[V], V veleroaddonv1.VeleroOrigin] struct {
	client.Client
	Tracker *remote.ClusterCacheTracker

	Controller controller.Controller
}

type VeleroReconciler[P veleroaddonv1.VeleroProxy[V], V veleroaddonv1.VeleroOrigin] interface {
	UpdateRemote(ctx context.Context, i *veleroaddonv1.VeleroInstallation, proxy P) (reconcile.Result, error)
	CleanupRemote(ctx context.Context, i *veleroaddonv1.VeleroInstallation, proxy P) (reconcile.Result, error)
	ReconcileProxy(ctx context.Context, i *veleroaddonv1.VeleroInstallation, proxy P) (reconcile.Result, error)
}

// AsVeleroReconciler creates a Reconciler based on the given ObjectReconciler.
func AsVeleroReconciler[P veleroaddonv1.VeleroProxy[V], V veleroaddonv1.VeleroOrigin](client client.Client, rec VeleroReconciler[P, V]) reconcile.ObjectReconciler[P] {
	return &Adapter[P, V]{
		objReconciler: rec,
		client:        client,
	}
}

type Adapter[P veleroaddonv1.VeleroProxy[V], V veleroaddonv1.VeleroOrigin] struct {
	objReconciler VeleroReconciler[P, V]
	client        client.Client
}

// ReconcileProxy implements VeleroReconciler.
func (a *Adapter[P, V]) ReconcileProxy(ctx context.Context, i *veleroaddonv1.VeleroInstallation, proxy P) (res reconcile.Result, err error) {
	if proxy.GetDeletionTimestamp() != nil {
		return a.objReconciler.CleanupRemote(ctx, i, proxy)
	}

	if res, err := a.objReconciler.ReconcileProxy(ctx, i, proxy); err != nil || res.Requeue || res.RequeueAfter > 0 {
		return res, err
	}

	return a.objReconciler.UpdateRemote(ctx, i, proxy)
}

// Reconcile implements Reconciler.
func (a *Adapter[P, V]) Reconcile(ctx context.Context, proxy P) (reconcile.Result, error) {
	ref := proxy.GetInstallRef()

	i := &veleroaddonv1.VeleroInstallation{}
	if err := a.client.Get(ctx, veleroaddonv1.RefToNamespaceName(ref).ObjectKey(), i); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	return a.ReconcileProxy(ctx, i, proxy)
}

func (r *Reconciler[P, V]) GetTracker() *remote.ClusterCacheTracker {
	return r.Tracker
}

func (r *Reconciler[P, V]) UpdateRemote(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, proxy P, obj V) (res reconcile.Result, reterr error) {
	var errs []error

	finalizers := proxy.GetFinalizers()
	if finalizers == nil {
		finalizers = []string{}
	}

	hasOwnerReferences := util.HasOwner(proxy.GetOwnerReferences(), veleroaddonv1.GroupVersion.String(), []string{installation.Kind})

	for _, clusterRef := range installation.Status.MatchingClusters {
		clusterKey := veleroaddonv1.RefToNamespaceName(&clusterRef).ObjectKey()
		if cl, err := r.GetTracker().GetClient(ctx, clusterKey); errors.Is(err, remote.ErrClusterLocked) {
			res.Requeue = true

			continue
		} else if err != nil {
			errs = append(errs, err)
		} else if err := cl.Get(ctx, client.ObjectKeyFromObject(obj), obj); apierrors.IsNotFound(err) {
			errs = append(errs, cl.Create(ctx, obj))
		} else if err != nil {
			errs = append(errs, err)
		} else if err := cl.Update(ctx, obj); err != nil {
			errs = append(errs, err)
		}

		if kerrors.NewAggregate(errs) == nil {
			proxy.SetClusterStatus(veleroaddonv1.NamespaceName(clusterKey.String()), obj)
			if err := r.Status().Update(ctx, proxy); !apierrors.IsConflict(err) {
				errs = append(errs, err)
			} else {
				res.Requeue = true
			}
		}

		errs = append(errs, controllerutil.SetOwnerReference(installation, proxy, r.Client.Scheme()))

		if err := r.GetTracker().Watch(ctx, remote.WatchInput{
			Name:         "remote-" + obj.GetObjectKind().GroupVersionKind().GroupKind().String(),
			Cluster:      clusterKey,
			Watcher:      r.Controller,
			Kind:         reflect.New(reflect.TypeOf(*new(V)).Elem()).Interface().(V),
			EventHandler: handler.EnqueueRequestsFromMapFunc(r.remoteToProxy),
		}); !errors.Is(err, remote.ErrClusterLocked) {
			errs = append(errs, err)
		}
	}

	if kerrors.NewAggregate(errs) == nil {
		if !slices.Contains(finalizers, finalizer) {
			proxy.SetFinalizers(append(finalizers, finalizer))
		}

		if !hasOwnerReferences {
			errs = append(errs, controllerutil.SetOwnerReference(installation, proxy, r.Client.Scheme()))
			if err := r.Update(ctx, proxy); !apierrors.IsConflict(err) {
				errs = append(errs, err)
			} else {
				res.Requeue = true
			}
		}
	}

	return res, kerrors.NewAggregate(errs)
}

func (r *Reconciler[P, V]) CleanupRemote(ctx context.Context, installation *veleroaddonv1.VeleroInstallation, proxy P, obj V) (res reconcile.Result, reterr error) {
	var errs []error

	for _, clusterRef := range installation.Status.MatchingClusters {
		clusterKey := veleroaddonv1.RefToNamespaceName(&clusterRef).ObjectKey()
		if cl, err := r.GetTracker().GetClient(ctx, clusterKey); errors.Is(err, remote.ErrClusterLocked) {
			res.Requeue = true

			continue
		} else if err != nil {
			errs = append(errs, err)
		} else if err := cl.Delete(ctx, obj); client.IgnoreNotFound(err) != nil {
			errs = append(errs, err)
		}
	}

	if removed := controllerutil.RemoveFinalizer(proxy, finalizer); removed && kerrors.NewAggregate(errs) == nil {
		errs = append(errs, r.Update(ctx, proxy))
	}

	return res, kerrors.NewAggregate(errs)
}

func (r *Reconciler[P, V]) remoteToProxy(ctx context.Context, remote client.Object) (req []reconcile.Request) {
	var annotation veleroaddonv1.NamespaceName

	if value, found := remote.GetAnnotations()[proxyKeyAnnotation]; !found {
		return
	} else {
		annotation = veleroaddonv1.NamespaceName(value)
	}

	proxy := reflect.New(reflect.TypeOf(*new(P)).Elem()).Interface().(P)
	if err := r.Client.Get(ctx, annotation.ObjectKey(), proxy); err != nil {
		return
	}

	req = append(req, reconcile.Request{
		NamespacedName: annotation.ObjectKey(),
	})

	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler[P, V]) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) *builder.Builder {
	return ctrl.NewControllerManagedBy(mgr).
		For(reflect.New(reflect.TypeOf(*new(P)).Elem()).Interface().(P), builder.WithPredicates(
			predicate.NewPredicateFuncs(r.onlyWithInstallation(ctx)))).
		WithOptions(options)
}

func (r *Reconciler[P, V]) onlyWithInstallation(ctx context.Context) func(client.Object) bool {
	return func(obj client.Object) bool {
		ref := obj.(P).GetInstallRef()
		installation := reflect.New(reflect.TypeOf(*new(P)).Elem()).Interface().(P)

		if ref == nil {
			return false
		}

		return r.Get(ctx, veleroaddonv1.RefToNamespaceName(ref).ObjectKey(), installation) != nil
	}
}
