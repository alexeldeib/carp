/*
Copyright 2020 Juan-Lee Pang.

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

// nolint: dupl
package controllers

import (
	"context"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrastructurev1alpha1 "github.com/juan-lee/carp/api/v1alpha1"
)

// WorkerReconciler reconciles a Worker object
type WorkerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=workers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=workers/status,verbs=get;update;patch

func (r *WorkerReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()
	log := r.Log.WithValues("worker", req.NamespacedName)

	var worker infrastructurev1alpha1.Worker
	if err := r.Get(ctx, req.NamespacedName, &worker); err != nil {
		log.Error(err, "unable to fetch worker")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	defer func() {
		if err := r.Status().Update(ctx, &worker); err != nil && reterr == nil {
			log.Error(err, "failed to update worker status")
			reterr = err
		}
	}()

	if worker.Status.AvailableCapacity == nil {
		worker.Status.AvailableCapacity = &worker.Spec.Capacity
		worker.Status.LastScheduledTime = metav1.Now()
	}

	worker.Status.Phase = infrastructurev1alpha1.WorkerRunning

	return ctrl.Result{}, nil
}

func (r *WorkerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.Worker{}).
		Complete(r)
}
