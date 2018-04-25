package k8s

import (
	apiv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

type IngressManager struct {
	Manager   `inject:""`
	Validator *ServiceValidator `inject:""`
}

func (m *IngressManager) getIngressClient() v1beta1.IngressInterface {
	return m.Client.ExtensionsV1beta1().Ingresses(m.Namespace)
}

func (m *IngressManager) List() (ingresses []apiv1beta1.Ingress, err error) {
	list, err := m.getIngressClient().List(metav1.ListOptions{})
	if err != nil {
		return
	}
	return list.Items, nil
}

func (m *IngressManager) Create(ingress *apiv1beta1.Ingress) (err error) {
	err = m.Validator.Validate(ingress)
	if err != nil {
		return
	}

	_, err = m.getIngressClient().Create(ingress)
	return
}

func (m *IngressManager) Get(name string) (ingress *apiv1beta1.Ingress, err error) {
	return m.getIngressClient().Get(name, metav1.GetOptions{})
}

func (m *IngressManager) Update(ingress *apiv1beta1.Ingress) (err error) {
	_, err = m.getIngressClient().Update(ingress)
	return
}

func (m *IngressManager) Remove(name string) (err error) {
	propagation := metav1.DeletePropagationForeground
	return m.getIngressClient().Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &propagation,
	})
}
