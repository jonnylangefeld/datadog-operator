/*
Copyright 2020 jonnylangefeld.

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

package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	datadogv1 "github.com/jonnylangefeld/datadog-operator/apis/datadog/v1"
	"github.com/jonnylangefeld/datadog-operator/pkg/datadog"
)

const (
	monitorFinalizer = "monitor.finalizers.datadog.jonnylangefeld.com"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	DatadogClient datadog.ClientInterface
}

// +kubebuilder:rbac:groups=datadog.jonnylangefeld.com,resources=monitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=datadog.jonnylangefeld.com,resources=monitors/status,verbs=get;update;patch

func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("monitor", req.NamespacedName)

	m := &datadogv1.Monitor{}
	if err := r.Get(ctx, req.NamespacedName, m); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch monitor")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !m.DeletionTimestamp.IsZero() {
		return Delete(ctx, r.Client, m, m.Status.ID, monitorFinalizer, r.DatadogClient.DeleteMonitor)
	}
	controllerutil.AddFinalizer(m, monitorFinalizer)

	if err := r.CreateOrUpdate(ctx, m); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MonitorReconciler) CreateOrUpdate(ctx context.Context, m *datadogv1.Monitor) error {
	transformed, err := m.Transform()
	if err != nil {
		return err
	}

	err = r.DatadogClient.UpdateMonitor(transformed)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			monitor, err := r.DatadogClient.CreateMonitor(transformed)
			if err != nil {
				return err
			}
			m.Status.ID = *monitor.Id
		} else {
			return err
		}
	}

	return r.Update(ctx, m)
}

func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datadogv1.Monitor{}).
		Complete(r)
}
