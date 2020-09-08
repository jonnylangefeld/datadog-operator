package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func HasFinalizer(finalizers []string, finalizer string) bool {
	for _, f := range finalizers {
		if f == finalizer {
			return true
		}
	}
	return false
}

func Delete(ctx context.Context, c client.Client, object interface{}, id int, finalizer string, f func(int) error) (ctrl.Result, error) {
	if HasFinalizer(object.(metav1.Object).GetFinalizers(), finalizer) && id != 0 {
		if err := f(id); err != nil {
			return ctrl.Result{}, err
		}
		controllerutil.RemoveFinalizer(object.(metav1.Object), finalizer)
		if err := c.Update(ctx, object.(runtime.Object)); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}
