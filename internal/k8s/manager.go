package k8s

import "k8s.io/client-go/kubernetes"

type Manager struct {
	Client    kubernetes.Interface `inject:"k8s"`
	Namespace string               `inject:"namespace"`
}
