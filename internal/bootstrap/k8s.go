package bootstrap

import (
	"github.com/facebookgo/inject"
	"code.ysitd.cloud/k8s/utils/go"
)

func injectK8s(graph *inject.Graph) {
	client, err := utils.AutoConnect()
	if err != nil {
		panic(err)
	}
	graph.Provide(
		&inject.Object{Name: "k8s", Value: client},
	)
}
