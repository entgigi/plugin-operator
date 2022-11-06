/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	v1alpha1 "github.com/entgigi/plugin-operator.git/api/v1alpha1"
	"github.com/entgigi/plugin-operator.git/common"
	"github.com/go-logr/logr"
)

const (
	controllerIngressLogName = "EntandoINgressV2 Controller"
)

// EntandoIngressV2Reconciler reconciles a EntandoIngressV2 object
type EntandoIngressV2Reconciler struct {
	Base   common.BaseK8sStructure
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=plugin.entando.org,resources=entandoingressv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=plugin.entando.org,resources=entandoingressv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=plugin.entando.org,resources=entandoingressv2s/finalizers,verbs=update

func NewEntandoIngressV2Reconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *EntandoIngressV2Reconciler {
	return &EntandoIngressV2Reconciler{
		Base:   common.BaseK8sStructure{Client: client, Log: log},
		Scheme: scheme,
	}
}

func (r *EntandoIngressV2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := r.Base.Log.WithName(controllerIngressLogName)
	log.Info("Start reconciling EntandoIngressV2 custom resources")

	cr := &v1alpha1.EntandoIngressV2{}
	err := r.Base.Get(ctx, req.NamespacedName, cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the EntandoApp instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isEntandoIngressV2MarkedToBeDeleted := cr.GetDeletionTimestamp() != nil
	if isEntandoIngressV2MarkedToBeDeleted {
		if err := r.removeFinalizer(ctx, cr, log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	err = r.addFinalizer(ctx, cr)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Reconciled EntandoIngressV2 custom resources")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntandoIngressV2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.EntandoIngressV2{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}). //solo modifiche a spec
		Complete(r)
}

// =====================================================================
// Add the cleanup steps that the operator
// needs to do before the CR can be deleted. Examples
// of finalizers include performing backups and deleting
// resources that are not owned by this CR, like a PVC.
// =====================================================================
func (r *EntandoIngressV2Reconciler) finalizeEntandoApp(log logr.Logger, m *v1alpha1.EntandoIngressV2) error {
	log.Info("Successfully finalized entandoApp")
	return nil
}

func (r *EntandoIngressV2Reconciler) addFinalizer(ctx context.Context, cr *v1alpha1.EntandoIngressV2) error {
	if !controllerutil.ContainsFinalizer(cr, entandoPluginFinalizer) {
		controllerutil.AddFinalizer(cr, entandoPluginFinalizer)
		return r.Base.Update(ctx, cr)
	}
	return nil
}

func (r *EntandoIngressV2Reconciler) removeFinalizer(ctx context.Context, cr *v1alpha1.EntandoIngressV2, log logr.Logger) error {
	if controllerutil.ContainsFinalizer(cr, entandoPluginFinalizer) {
		// Run finalization logic for entandoAppFinalizer. If the
		// finalization logic fails, don't remove the finalizer so
		// that we can retry during the next reconciliation.
		if err := r.finalizeEntandoApp(log, cr); err != nil {
			return err
		}

		// Remove entandoAppFinalizer. Once all finalizers have been
		// removed, the object will be deleted.
		controllerutil.RemoveFinalizer(cr, entandoPluginFinalizer)
		err := r.Base.Update(ctx, cr)
		if err != nil {
			return err
		}
	}
	return nil
}
