package httpserver

import (
	"net/http"

	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
)

type HTTPServer struct {
	mux     *http.ServeMux
	address string
}

func New(config *serverconfig.HTTPServerConfig) *HTTPServer {
	return &HTTPServer{
		mux:     http.NewServeMux(),
		address: config.Address,
	}
}

func (s *HTTPServer) Run() error {
	return http.ListenAndServe(s.address, s.mux)
}
