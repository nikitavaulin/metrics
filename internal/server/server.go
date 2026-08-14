package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
)

type HTTPServer struct {
	mux     *chi.Mux
	address string
}

func New(config *serverconfig.Config) *HTTPServer {
	return &HTTPServer{
		mux:     chi.NewRouter(),
		address: config.Address,
	}
}

func (s *HTTPServer) RegisterRouter(router Router) {
	s.mux.Mount("/", router)
}

func (s *HTTPServer) Run() error {
	return http.ListenAndServe(s.address, LogRequest(s.mux))
}
