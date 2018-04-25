package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func initLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		if os.Getenv("VERBOSE") != "" {
			logger.SetLevel(logrus.DebugLevel)
		}
	}

	return logger
}

func injectLogger(graph *inject.Graph) {
	logger := initLogger()

	graph.Provide(
		&inject.Object{Name: "service logger", Value: logger.WithField("source", "service")},
		&inject.Object{Name: "syncer logger", Value: logger.WithField("source", "syncer")},
	)
}
