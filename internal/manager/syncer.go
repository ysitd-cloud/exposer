package manager

import (
	"code.ysitd.cloud/component/exposer/internal/k8s"

	"github.com/sirupsen/logrus"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Syncer struct {
	Manager  *k8s.IngressManager `inject:""`
	listener chan Mutation
	Logger   logrus.FieldLogger `inject:"syncer logger"`
}

func (s *Syncer) Run() {
	if s.listener == nil {
		s.listener = make(chan Mutation, 5)
	}

	for {
		mutation := <-s.listener
		var err error
		switch mutation.Action {
		case Create:
			err = s.Connect(mutation.Hostname, mutation.ServiceName, mutation.Port)
			break
		case Update:
			err = s.Migrate(mutation.Hostname, mutation.ServiceName, mutation.Port)
			break
		case Delete:
			err = s.Disconnect(mutation.Hostname)
			break
		}

		if err != nil {
			s.Logger.Errorln(err)
		}
	}
}

func (s *Syncer) GetChannel() chan<- Mutation {
	if s.listener == nil {
		s.listener = make(chan Mutation, 5)
	}
	return s.listener
}

func (s *Syncer) GetBonding(hostname string) (conn *Connection, err error) {
	ingress, err := s.Manager.Get(ingressName(hostname))
	if errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return
	}

	conn = ingressToConnection(ingress)
	return
}

func (s *Syncer) ListBonding() (connections []*Connection, err error) {
	ingresses, err := s.Manager.List("creator=exposer")
	if err != nil {
		return
	}

	for _, ingress := range ingresses {
		connections = append(connections, ingressToConnection(&ingress))
	}
	return
}

func (s *Syncer) Connect(hostname, service string, port int) error {
	ingress := s.createIngress(hostname, service, port)
	return s.Manager.Create(ingress)
}

func (s *Syncer) createIngress(hostname, service string, port int) *v1beta1.Ingress {
	return &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: ingressName(hostname),
			Labels: map[string]string{
				"creator":  "exposer",
				"hostname": hostname,
				"service":  service,
			},
			Annotations: map[string]string{
				"code.ysitd.cloud/exposer/version":           "v1alpha1",
				"code.ysitd.cloud/exposer/v1alpha1/hostname": hostname,
				"code.ysitd.cloud/exposer/v1alpha1/service":  service,
			},
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

	backend := ingress.Spec.Rules[0].HTTP.Paths[0].Backend

	backend.ServiceName = service
	backend.ServicePort = intstr.FromInt(port)

	return s.Manager.Update(ingress)
}

func (s *Syncer) Disconnect(hostname string) (err error) {
	return s.Manager.Remove(ingressName(hostname))
}
