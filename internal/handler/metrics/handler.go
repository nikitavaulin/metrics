package metricshandler

import (
	"net/http"

	httpserver "github.com/nikitavaulin/metrics/internal/server"
	"github.com/nikitavaulin/metrics/internal/service"
)

type MetricsHandler struct {
	metricsService service.MetricsServie
}

func New(metricsService service.MetricsServie) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

func (h *MetricsHandler) Routes() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodPost,
			Path:    "/update/{type}/{name}/{value}",
			Handler: h.UpdateMetrics,
		},
	}
}
