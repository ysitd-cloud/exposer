package bootstrap

import (
	"code.ysitd.cloud/component/exposer/internal/http"
	"code.ysitd.cloud/component/exposer/internal/manager"
	"github.com/facebookgo/inject"
)

var service *http.Service
var syncer *manager.Syncer

type injector func(graph *inject.Graph)

func init() {
	var graph inject.Graph

	logger := initLogger()
	injectionLogger := logger.WithField("source", "inject")

	graph.Logger = injectionLogger

	injectors := []injector{
		injectLogger,
		injectVodka,
		injectK8s,
	}

	var s http.Service
	var sync manager.Syncer

	graph.Provide(
		&inject.Object{Value: &s},
		&inject.Object{Value: &sync},
	)

	for _, fn := range injectors {
		fn(&graph)
	}

	if err := graph.Populate(); err != nil {
		injectionLogger.Errorln(err)
		panic(err)
	}
	service = &s
	syncer = &sync
}

func GetService() *http.Service {
	return service
}

func GetSyncer() *manager.Syncer {
	return syncer
}
