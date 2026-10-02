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

	if err := h.metricsService.Add(r.Context(), metric.ID, metric); err != nil {
		httpserver.ErrorResponse(rw, err, "unexpected error", http.StatusInternalServerError)
		return
	}

	h.metricsService.SetMValueFromStorage(r.Context(), &metric)
	httpserver.JSONResponse(rw, &metric, http.StatusOK)
}

func (h *MetricsHandler) UpdateMetricsByJSON(rw http.ResponseWriter, r *http.Request) {
	metrics := make([]models.Metrics, 0)

	if err := httpserver.DecodeRequestBody(r, &metrics); err != nil {
		httpserver.ErrorResponse(rw, err, "failed to decode metrics", http.StatusBadRequest)
		return
	}

	if err := validateMetrics(metrics); err != nil {
		httpserver.ErrorResponse(rw, err, "invalid one of the metrics", http.StatusBadRequest)
		return
	}

	if err := h.metricsService.UpdateMetrics(r.Context(), metrics); err != nil {
		httpserver.ErrorResponse(rw, err, "failed to update metrics", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func validateMetrics(metrics []models.Metrics) error {
	for _, m := range metrics {
		if err := validation.ValidateMetric(m); err != nil {
			return err
		}
	}
	return nil
}
