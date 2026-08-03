package metricshandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nikitavaulin/metrics/internal/commonerrors"
	models "github.com/nikitavaulin/metrics/internal/model"
)

func (h *MetricsHandler) GetMetric(rw http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")

	if err := validateMType(mtype); err != nil {
		http.Error(rw, fmt.Sprintf("invalid mtype: %v", err), http.StatusBadRequest)
		return
	}

	if err := validateMName(name); err != nil {
		http.Error(rw, fmt.Sprintf("invalid mname: %v", err), http.StatusBadRequest)
		return
	}

	value, err := h.metricsService.Get(name)
	if err != nil {
		if errors.Is(err, commonerrors.ErrNotFound) {
			http.Error(rw, fmt.Sprintf("metric %s not found", name), http.StatusNotFound)
		} else {
			http.Error(rw, fmt.Sprintf("failed to get metric: %s", err), http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(rw, value)
	rw.WriteHeader(http.StatusOK)
}

func validateMType(mtype string) error {
	if mtype == "" {
		return fmt.Errorf("mtype is empty")
	}
	if mtype != models.Counter && mtype != models.Gauge {
		return fmt.Errorf("unknown metric type: %s", mtype)
	}
	return nil
}

func validateMName(name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}
	return nil
}
