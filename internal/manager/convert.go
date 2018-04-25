package manager

import (
	"fmt"
	"k8s.io/api/extensions/v1beta1"
	"strings"
)

func ingressName(hostname string) string {
	return fmt.Sprintf("expose-%s", strings.Replace(hostname, ".", "-", -1))
}

func ingressToConnection(ingress *v1beta1.Ingress) *Connection {
	rule := ingress.Spec.Rules[0]
	return &Connection{
		Hostname: rule.Host,
		Service:  rule.HTTP.Paths[0].Backend.ServiceName,
		Port:     int(rule.HTTP.Paths[0].Backend.ServicePort.IntVal),
	}
}
