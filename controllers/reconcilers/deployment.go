package reconcilers

import (
	"context"

	"github.com/entgigi/plugin-operator.git/api/v1alpha1"
	"github.com/entgigi/plugin-operator.git/utility"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
)

const labelKey = "app"

func (d *DeployManager) isUpgrade(ctx context.Context, cr *v1alpha1.EntandoPluginV2, deployment *appsv1.Deployment) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: makeDeploymentName(cr), Namespace: cr.GetNamespace()}, deployment)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *DeployManager) buildDeployment(cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) *appsv1.Deployment {
	replicas := cr.Spec.Replicas
	deploymentName := makeDeploymentName(cr)
	labels := map[string]string{labelKey: deploymentName}
	port := int32(cr.Spec.Port)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: cr.GetNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{StrVal: "25%", Type: intstr.String},
					MaxSurge:       &intstr.IntOrString{StrVal: "25%", Type: intstr.String},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           cr.Spec.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            "mycontainer",
						Ports: []corev1.ContainerPort{{
							ContainerPort: port,
						}},
						Env: cr.Spec.EnvironmentVariables,
						ReadinessProbe: &v1.Probe{
							ProbeHandler: v1.ProbeHandler{
								HTTPGet: &v1.HTTPGetAction{Path: cr.Spec.HealthCheckPath, Port: intstr.IntOrString{
									IntVal: port,
								}},
							},
							InitialDelaySeconds: 10,
						},
						LivenessProbe: &v1.Probe{
							ProbeHandler: v1.ProbeHandler{
								HTTPGet: &v1.HTTPGetAction{Path: cr.Spec.HealthCheckPath, Port: intstr.IntOrString{
									IntVal: port,
								}},
							},
							InitialDelaySeconds: 10,
						},
						StartupProbe: &v1.Probe{
							ProbeHandler: v1.ProbeHandler{
								HTTPGet: &v1.HTTPGetAction{Path: cr.Spec.HealthCheckPath, Port: intstr.IntOrString{
									IntVal: port,
								}},
							},
							InitialDelaySeconds: 20,
						},
					}},
				},
			},
		},
	}
	// set owner
	ctrl.SetControllerReference(cr, deployment, scheme)
	return deployment
}

func makeDeploymentName(cr *v1alpha1.EntandoPluginV2) string {
	return "plugin-" + utility.TruncateString(cr.GetName(), 200) + "-deployment"
}
