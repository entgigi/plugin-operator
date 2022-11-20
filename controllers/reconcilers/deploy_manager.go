package reconcilers

import (
	"context"
	"time"

	"github.com/entgigi/plugin-operator.git/api/v1alpha1"

	"github.com/entgigi/plugin-operator.git/common"
	"github.com/entgigi/plugin-operator.git/controllers/services"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type DeployManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewDeployManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *DeployManager {
	return &DeployManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *DeployManager) IsDeployApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsDeployApplied(ctx, cr)
}

func (d *DeployManager) IsDeployReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsDeployReady(ctx, cr)
}

func (d *DeployManager) ApplyDeploy(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {
	time.Sleep(time.Second * 10)
	baseDeployment := d.buildDeployment(cr, scheme)
	deployment := &appsv1.Deployment{}

	err, isUpgrade := d.isUpgrade(ctx, cr, deployment)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		deployment.Spec = baseDeployment.Spec
		applyError = d.Base.Client.Update(ctx, deployment)

	} else {
		applyError = d.Base.Client.Create(ctx, baseDeployment)
	}

	if applyError != nil {
		return applyError
	}

	return d.Conditions.SetConditionDeployApplied(ctx, cr)
}

func (d *DeployManager) CheckDeploy(ctx context.Context, cr *v1alpha1.EntandoPluginV2) (bool, error) {
	time.Sleep(time.Second * 10)
	ready := true

	if ready {
		return ready, d.Conditions.SetConditionDeployReady(ctx, cr)
	}

	return ready, nil

}
