package controllers

import (
	"context"

	"github.com/go-logr/logr"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	installv1alpha1 "kore3lab.io/kore/operator/api/v1alpha1"
	"kore3lab.io/kore/pkg/helmreconciler"
)

const (
	finalizer           = "finalizer.install.kore3lab.io"
	finalizerMaxRetries = 1
)

// KoreOperatorReconciler reconciles a KoreOperator object
type KoreOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger logr.Logger
}

// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *KoreOperatorReconciler) Reconcile(ctx context.Context, req reconcile.Request) (ctrl.Result, error) {

	// parsing a KoreOperator
	operator := &installv1alpha1.KoreOperator{}
	if err := r.Get(context.TODO(), req.NamespacedName, operator); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("KoreOperator is not found", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, nil
		}
		r.Logger.Error(err, "Fail to get a KoreOperator", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	}
	if operator.Spec == nil {
		operator.Spec = &installv1alpha1.KoreOperatorSpec{
			Components: map[string]installv1alpha1.KoreOperatorComponent{},
		}
	}

	// executing finalizer if operator deleted
	deleted := operator.GetDeletionTimestamp() != nil
	if deleted == true {
		if controllerutil.ContainsFinalizer(operator, finalizer) {
			if reconciler, err := helmreconciler.NewHelmReconciler(operator, nil, r.Client); err != nil {
				r.Logger.Error(err, "Delete to finalizer by KoreOperator")
			} else if err := reconciler.Finalize(); err == nil {
				controllerutil.RemoveFinalizer(operator, finalizer)
				err := r.Update(context.TODO(), operator)
				if err != nil {
					r.Logger.Error(err, "Fail to update a KoreOperator/finalizer")
					return ctrl.Result{}, err
				}
			} else {
				r.Logger.Error(err, "Fail to delete objects", "name", req.Name, "namespace", req.Namespace)
				return ctrl.Result{}, err
			}
			r.Logger.Info("Finalizer was successfully removed", "name", req.Name, "namespace", req.Namespace)
		}
		r.Logger.Info("Successfully deleting", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

	// regist finalizer if it not contains
	if !controllerutil.ContainsFinalizer(operator, finalizer) {
		controllerutil.AddFinalizer(operator, finalizer)
		if err := r.Update(context.TODO(), operator); err != nil {
			r.Logger.Error(err, "Update to finalizer by KoreOperator")
			return ctrl.Result{}, err
		}
		r.Logger.Info("Finalizer was successfully added")
	}

	if reconciler, err := helmreconciler.NewHelmReconciler(operator, nil, r.Client); err != nil {
		r.Logger.Error(err, "Fail to create a instance of Reconciler", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, err
	} else {
		if err := reconciler.Begin(); err != nil {
			r.Logger.Error(err, "Fail to reconcile begin", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		status := reconciler.Reconcile()
		if status.Status != installv1alpha1.STATUS_COMPLETE {
			r.Logger.Info("Fail to reconcile", "name", req.Name, "namespace", req.Namespace, "status", status)
		}

		if err := reconciler.End(status); err != nil {
			r.Logger.Error(err, "Fail to update a reconciling status", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, err
		}

		r.Logger.Info("Successfully reconciling", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, nil
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *KoreOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetConfig()
	return ctrl.NewControllerManagedBy(mgr).
		For(&installv1alpha1.KoreOperator{}).
		Complete(r)
}

func (r *KoreOperatorReconciler) createNamespace(ctx context.Context, namespacedName types.NamespacedName, spec *installv1alpha1.KoreOperatorSpec, ns *v1.Namespace) error {

	ns.ObjectMeta = metaV1.ObjectMeta{Name: namespacedName.Name, Namespace: namespacedName.Namespace}
	if err := r.Create(ctx, ns); err != nil {
		r.Logger.Error(err, "Create a namespace", "namespace", namespacedName.Name)
		return err
	}
	r.Logger.Info("create a namespace", "namespace", namespacedName.Name)
	return nil
}
