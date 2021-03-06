package router

import (
	"server/interactor"
	mid "server/presenter/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	Route *chi.Mux
}

func NewRouter() *Server {
	return &Server{
		Route: chi.NewRouter(),
	}
}

func (s *Server) Router(h *interactor.Handler, m mid.Middleware) {
	s.Route.Use(middleware.Logger)
	s.Route.Use(middleware.Recoverer)
	s.Route.Route("/", func(r chi.Router) {
		r.Get("/", h.SampleHandler.SampleHTML)
	})
	s.Route.Route("/api/v1", func(r chi.Router) {
		r.Get("/", h.SampleHandler.SampleIndex)
		// TODO
	})
}