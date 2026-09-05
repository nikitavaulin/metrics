package metricshandler

import (
	"net/http"

	models "github.com/nikitavaulin/metrics/internal/model"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
	"github.com/nikitavaulin/metrics/internal/validation"
)

func (h *MetricsHandler) UpdateMetricByJSON(rw http.ResponseWriter, r *http.Request) {
	var metric models.Metrics

	if err := httpserver.DecodeRequestBody(r, &metric); err != nil {
		httpserver.ErrorResponse(rw, err, "failed to decode metric", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateMetricName(metric.ID); err != nil {
		httpserver.ErrorResponse(rw, err, err.Error(), http.StatusNotFound)
		return
	}

	if err := validation.ValidateMetric(metric); err != nil {
		httpserver.ErrorResponse(rw, err, "invalid metric", http.StatusBadRequest)
		return
	}

	if err := h.metricsService.Add(metric.ID, metric); err != nil {
		httpserver.ErrorResponse(rw, err, "unexpected error", http.StatusInternalServerError)
		return
	}

	h.metricsService.SetMValueFromStorage(&metric)
	httpserver.JSONResponse(rw, &metric, http.StatusOK)
}
