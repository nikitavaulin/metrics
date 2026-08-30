package metricshandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nikitavaulin/metrics/internal/commonerrors"
	models "github.com/nikitavaulin/metrics/internal/model"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
	"github.com/nikitavaulin/metrics/internal/validation"
)

func (h *MetricsHandler) GetJSONMetrics(rw http.ResponseWriter, r *http.Request) {
	var metric models.Metrics
	if err := httpserver.DecodeRequestBody(r, &metric); err != nil {
		http.Error(rw, "failed to decode metric", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateMetricType(metric.MType); err != nil {
		http.Error(rw, fmt.Sprintf("invalid mtype: %v", err), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateMetricName(metric.ID); err != nil {
		http.Error(rw, fmt.Sprintf("invalid mname: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.metricsService.SetMValueFromStorage(&metric); err != nil {
		if errors.Is(err, commonerrors.ErrNotFound) {
			http.Error(rw, fmt.Sprintf("metric %s not found", metric.ID), http.StatusNotFound)
		} else {
			http.Error(rw, fmt.Sprintf("failed to get metric: %s", err), http.StatusInternalServerError)
		}
		return
	}

	httpserver.JSONResponse(rw, &metric, http.StatusOK)
}
