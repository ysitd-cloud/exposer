package k8s

import (
	"errors"

	apiCoreV1 "k8s.io/api/core/v1"
	apiv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func checkPort(actual apiCoreV1.ServicePort, expected intstr.IntOrString) bool {
	if expected.Type == intstr.Int {
		return actual.Port == expected.IntVal
	} else {
		return actual.Name == expected.StrVal
	}
}

var ErrNonListenPort = errors.New("connect to non-listen port")

type ServiceValidator struct {
	Manager `inject:""`
}

func (v *ServiceValidator) getClient() corev1.ServiceInterface {
	return v.Client.CoreV1().Services(v.Namespace)
}

func (v *ServiceValidator) validateService(serviceName string, port intstr.IntOrString) (err error) {
	svc, err := v.getClient().Get(serviceName, metav1.GetOptions{})
	if err != nil {
		return
	}

	for _, servicePort := range svc.Spec.Ports {
		if checkPort(servicePort, port) {
			return nil
		}
	}

	return ErrNonListenPort
}

func (v *ServiceValidator) Validate(ingress *apiv1beta1.Ingress) error {
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			name := path.Backend.ServiceName
			port := path.Backend.ServicePort
			if err := v.validateService(name, port); err != nil {
				return err
			}
		}
	}
	return nil
}
