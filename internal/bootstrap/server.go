package bootstrap

import (
	"github.com/facebookgo/inject"
	"github.com/tonyhhyip/vodka"
)

func initServer() *vodka.Server {
	return vodka.New(":http")
}

func injectVodka(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: initServer()},
	)
}
