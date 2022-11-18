package reconcilers

import (
	"context"
	"fmt"
	"time"

	"github.com/entgigi/plugin-operator.git/api/v1alpha1"
	"github.com/entgigi/plugin-operator.git/common"
	"github.com/entgigi/plugin-operator.git/controllers/services"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileManager struct {
	Base      *common.BaseK8sStructure
	Recorder  record.EventRecorder
	Condition *services.ConditionService
}

func NewReconcileManager(client client.Client, log logr.Logger, recorder record.EventRecorder) *ReconcileManager {
	base := &common.BaseK8sStructure{Client: client, Log: log}
	return &ReconcileManager{
		Base:      base,
		Recorder:  recorder,
		Condition: services.NewConditionService(base),
	}
}

func (r *ReconcileManager) MainReconcile(ctx context.Context, req ctrl.Request, cr *v1alpha1.EntandoPluginV2) (ctrl.Result, error) {
	log := r.Base.Log
	deployManager := NewDeployManager(r.Base, r.Condition)
	// deploy done
	applied := deployManager.IsDeployApplied(ctx, cr)

	if !applied {
		if err := deployManager.ApplyDeploy(ctx, cr); err != nil {
			log.Info("error ApplyDeploy reschedule reconcile", "error", err)
			return ctrl.Result{}, err
		}
	}
	r.Recorder.Eventf(cr, "Normal", "Updated", fmt.Sprintf("Updated deployment %s/%s", req.Namespace, req.Name))

	// deploy ready
	var err error
	ready := deployManager.IsDeployReady(ctx, cr)

	if !ready {
		if ready, err = deployManager.CheckDeploy(ctx, cr); err != nil {
			log.Info("error ApplyDeploy reschedule reconcile", "error", err)
			return ctrl.Result{}, err
		}
		if !ready {
			log.Info("Deploy not ready reschedule operator", "seconds", 10)
			r.Recorder.Eventf(cr, "Warning", "NotReady", fmt.Sprintf("Plugin deployment not ready %s/%s", req.Namespace, req.Name))
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
	}

	// ingress requested

	// ingress ready

	r.Recorder.Eventf(cr, "Normal", "Done", fmt.Sprintf("Plugin deployed %s/%s", req.Namespace, req.Name))

	return ctrl.Result{}, nil
}
