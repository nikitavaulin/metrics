package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	models "github.com/nikitavaulin/metrics/internal/model"
)

const (
	contentType     string = "text/plain"
	contentTypeJSON string = "application/json"
)

func (a *Agent) Report() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/update/", a.serverAddr)

	for name, value := range a.metrics {
		val := float64(value)

		metric := models.Metrics{
			ID:    name,
			Value: &val,
			MType: models.Gauge,
		}

		if err := a.sendJSONMetric(url, metric, client); err != nil {
			return fmt.Errorf("failed to send metric %q: %w", name, err)
		}
	}

	randomValue := float64(a.randomValue)
	randomMetric := models.Metrics{
		ID:    randomValueName,
		Value: &randomValue,
		MType: models.Gauge,
	}

	if err := a.sendJSONMetric(url, randomMetric, client); err != nil {
		return fmt.Errorf("failed to send random value metric: %w", err)
	}

	pollCount := int64(a.pollCount)
	pollCountMetric := models.Metrics{
		ID:    pollCountName,
		Delta: &pollCount,
		MType: models.Counter,
	}

	if err := a.sendJSONMetric(url, pollCountMetric, client); err != nil {
		return fmt.Errorf("failed to send poll count metric: %w", err)
	}

	return nil
}

func (a *Agent) sendJSONMetric(url string, metric models.Metrics, client *http.Client) error {
	body, err := json.Marshal(&metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}

	compressedBody, err := compress(body)
	if err != nil {
		return fmt.Errorf("failed to compress json metrics: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(compressedBody))
	if err != nil {
		return fmt.Errorf("failed to create post request: %w", err)
	}
	req.Header.Add("Content-Type", contentTypeJSON)
	req.Header.Add("Content-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to POST metric: %w", err)
	}
	defer resp.Body.Close()

	logResponse(resp)

	return nil
}

func (a *Agent) getPath(mtype, name, value string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s", a.serverAddr, mtype, name, value)
}

func logResponse(resp *http.Response) {
	if resp != nil {
		errMsg := "null"
		if resp.StatusCode != http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("failed to read req body: %v\n", err)
			}
			errMsg = string(body)
		}
		log.Printf(
			"DEBUG: method: %s url: %s status: %s error: %s\n",
			resp.Request.Method,
			resp.Request.URL.Path,
			resp.Status,
			errMsg,
		)
		resp.Body.Close()
	}
}
