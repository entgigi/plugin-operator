package reconcilers

import (
	"context"
	"time"

	"github.com/entgigi/plugin-operator.git/api/v1alpha1"
	"github.com/entgigi/plugin-operator.git/utility"

	"github.com/entgigi/plugin-operator.git/common"
	"github.com/entgigi/plugin-operator.git/controllers/services"

	"k8s.io/apimachinery/pkg/runtime"
)

const labelKey = "app"

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
	applyError := d.ApplyDeployment(ctx, cr, scheme)
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

func makeContainerName(cr *v1alpha1.EntandoPluginV2) string {
	return "plugin-" + utility.TruncateString(cr.GetName(), 200) + "-container"
}
