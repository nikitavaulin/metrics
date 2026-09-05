package metricshandler

import (
	"fmt"
	"net/http"
)

func (h *MetricsHandler) GetMetricsList(rw http.ResponseWriter, r *http.Request) {
	metrics, err := h.metricsService.GetList()
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get metrics list: %v", err), http.StatusInternalServerError)
		return
	}
	html := getMetricsHTML(metrics)
	rw.Header().Set("Content-Type", "text/html")
	rw.Write([]byte(html))
	rw.WriteHeader(http.StatusOK)
}

func getMetricsHTML(metrics map[string]any) string {
	html := `<html><head><title>Metrics</title><body><ul>`
	for name, value := range metrics {
		html += fmt.Sprintf("<li>%s: %v</li>", name, value)
	}
	html += `</ul></body></html>`
	return html
}
