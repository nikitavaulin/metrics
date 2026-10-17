package pinghandler

import (
	"net/http"

	httpserver "github.com/nikitavaulin/metrics/internal/server"
	"github.com/nikitavaulin/metrics/internal/service"
)

type PingHandler struct {
	pingService service.PingService
}

func New(pingService service.PingService) *PingHandler {
	return &PingHandler{pingService: pingService}
}

func (h *PingHandler) Routes() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodGet,
			Path:    "/ping",
			Handler: h.PingDB,
		},
	}
}
