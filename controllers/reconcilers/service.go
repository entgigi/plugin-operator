package reconcilers

import (
	"context"

	"github.com/entgigi/plugin-operator.git/api/v1alpha1"
	"github.com/entgigi/plugin-operator.git/utility"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (d *DeployManager) isServiceUpgrade(ctx context.Context, cr *v1alpha1.EntandoPluginV2, service *corev1.Service) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: makeServiceName(cr), Namespace: cr.GetNamespace()}, service)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *DeployManager) buildService(cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) *corev1.Service {
	serviceName := makeServiceName(cr)
	labels := map[string]string{labelKey: makeContainerName(cr)}
	port := int32(cr.Spec.Port)

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: cr.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{{
				Name:       "server-port",
				Port:       port,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.IntOrString{IntVal: port},
			}},
			Selector: labels,
		},
	}
	// set owner
	ctrl.SetControllerReference(cr, service, scheme)
	return service
}

func makeServiceName(cr *v1alpha1.EntandoPluginV2) string {
	return "plugin-" + utility.TruncateString(cr.GetName(), 200) + "-service"
}

func (d *DeployManager) ApplyService(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {
	baseService := d.buildService(cr, scheme)
	service := &corev1.Service{}

	err, isUpgrade := d.isServiceUpgrade(ctx, cr, service)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		service.Spec = baseService.Spec
		applyError = d.Base.Client.Update(ctx, service)

	} else {
		applyError = d.Base.Client.Create(ctx, baseService)
	}

	if applyError != nil {
		return applyError
	}
	return nil
}