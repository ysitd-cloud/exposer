package internal

import (
	"fmt"
	"strings"

	"code.ysitd.cloud/component/exposer/internal/k8s"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ingressName(hostname string) string {
	return fmt.Sprintf("expose-%s", strings.Replace(hostname, ".", "-", -1))
}

type Syncer struct {
	Manager  *k8s.IngressManager `inject:""`
	Listener chan Mutation
}

func (s *Syncer) Connect(hostname, service string, port int) error {
	ingress := s.createIngress(hostname, service, port)
	return s.Manager.Create(ingress)
}

func (s *Syncer) createIngress(hostname, service string, port int) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressName(hostname),
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: service,
										ServicePort: intstr.FromInt(port),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s *Syncer) Migrate(hostname, service string, port int) (err error) {
	ingress, err := s.Manager.Get(ingressName(hostname))
	if err != nil {
		return
	}

	ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName = service
	ingress.Spec.Rules[0].HTTP.Paths[0].Backend.ServicePort = intstr.FromInt(port)

	return s.Manager.Update(ingress)
}

func (s *Syncer) Disconnect(hostname string) (err error) {
	return s.Manager.Remove(ingressName(hostname))
}
