package http

import (
	"code.ysitd.cloud/component/exposer/internal/manager"
	"github.com/tonyhhyip/vodka"
	"net/http"
)

type Service struct {
	Server  *vodka.Server   `inject:""`
	Syncer  *manager.Syncer `inject:""`
	handler http.Handler
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.handler == nil {
		s.handler = s.ReadyHandler()
	}

	s.handler.ServeHTTP(w, r)
}

func (s *Service) ReadyHandler() http.Handler {
	router := s.createRouter()
	return s.Server.WrapHandler(router.Handler())
}

func (s *Service) createRouter() (router *vodka.Router) {
	router = vodka.NewRouter()

	router.POST("/:hostname", s.createBond)

	return
}
