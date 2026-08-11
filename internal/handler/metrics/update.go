package metricshandler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	models "github.com/nikitavaulin/metrics/internal/model"
)

func (h *MetricsHandler) UpdateMetrics(rw http.ResponseWriter, r *http.Request) {
	metric, err := parseMetric(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(rw, "metric name is empty", http.StatusNotFound)
		return
	}

	if err := h.metricsService.Add(name, metric); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func parseMetric(r *http.Request) (models.Metrics, error) {
	mtype := chi.URLParam(r, "type")
	if mtype == "" {
		return models.Metrics{}, fmt.Errorf("metric type is empty")
	}

	value := chi.URLParam(r, "value")
	if value == "" {
		return models.Metrics{}, fmt.Errorf("value is empty")
	}

	var metric models.Metrics

	switch mtype {
	case models.Counter:
		number, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return models.Metrics{}, fmt.Errorf("failed to parse int number value: %w", err)
		}
		metric.Delta = &number

	case models.Gauge:
		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return models.Metrics{}, fmt.Errorf("failed to parse number value: %w", err)
		}
		metric.Value = &number

	default:
		return models.Metrics{}, fmt.Errorf("unknown metric type")
	}

	metric.MType = mtype
	return metric, nil
}
