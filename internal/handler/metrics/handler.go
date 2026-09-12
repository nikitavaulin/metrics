package metricshandler

import (
	"net/http"

	httpserver "github.com/nikitavaulin/metrics/internal/server"
	"github.com/nikitavaulin/metrics/internal/service"
)

type MetricsHandler struct {
	metricsService service.MetricsService
}

func New(metricsService service.MetricsService) *MetricsHandler {
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
		{
			Method:  http.MethodPost,
			Path:    "/update/",
			Handler: h.UpdateMetricByJSON,
		},
		{
			Method:  http.MethodGet,
			Path:    "/value/{type}/{name}",
			Handler: h.GetMetric,
		},
		{
			Method:  http.MethodPost,
			Path:    "/value/",
			Handler: h.GetJSONMetrics,
		},
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: h.GetMetricsList,
		},
	}
}
