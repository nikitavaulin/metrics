package httpserver

import (
	"log"
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
	log.Printf("Running server on %s...\n", s.address)
	return http.ListenAndServe(s.address, s.mux)
}
