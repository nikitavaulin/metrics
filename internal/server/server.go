package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
)

type HTTPServer struct {
	mux        *chi.Mux
	address    string
	hashingKey *string
}

func New(config *serverconfig.Config) *HTTPServer {
	return &HTTPServer{
		mux:        chi.NewRouter(),
		address:    config.Address,
		hashingKey: config.Key,
	}
}

func (s *HTTPServer) RegisterRouter(router Router) {
	s.mux.Mount("/", router)
}

func (s *HTTPServer) Run() error {
	return http.ListenAndServe(s.address, LogRequest(s.handler()))
}

func (s *HTTPServer) handler() http.Handler {
	return GzipResponse(HashResponse(GzipRequest(HashRequest(s.mux, s.hashingKey)), s.hashingKey))
}
