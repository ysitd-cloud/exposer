package http

import (
	"net/http"

	"code.ysitd.cloud/component/exposer/internal/manager"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/tonyhhyip/vodka"
)

type Service struct {
	Server  *vodka.Server   `inject:""`
	Syncer  *manager.Syncer `inject:""`
	Logger  *logrus.Entry   `inject:"service logger"`
	handler http.Handler
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.handler == nil {
		s.handler = s.ReadyHandler()
	}

	s.handler.ServeHTTP(w, r)
}

func (s *Service) ReadyHandler() (handler http.Handler) {
	router := s.createRouter()
	handler = s.Server.WrapHandler(router.Handler())
	handler = handlers.CombinedLoggingHandler(s.Logger.Writer(), handler)
	handler = handlers.RecoveryHandler(handlers.RecoveryLogger(s.Logger), handlers.PrintRecoveryStack(true))(handler)
	return
}

func (s *Service) createRouter() (router *vodka.Router) {
	router = vodka.NewRouter()

	router.GET("/", s.listBonds)
	router.GET("/:hostname", s.getBond)
	router.POST("/:hostname", s.createBond)
	router.PUT("/:hostname", s.updateBond)
	router.DELETE("/:hostname", s.removeBond)

	return
}
