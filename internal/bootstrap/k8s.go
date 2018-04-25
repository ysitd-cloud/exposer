package bootstrap

import (
	"code.ysitd.cloud/k8s/utils/go"
	"github.com/facebookgo/inject"
	"os"
)

func injectK8s(graph *inject.Graph) {
	client, err := utils.AutoConnect()
	if err != nil {
		panic(err)
	}
	graph.Provide(
		&inject.Object{Name: "k8s", Value: client},
		&inject.Object{Name: "namespace", Value: os.Getenv("NAMESPACE")},
	)
}
