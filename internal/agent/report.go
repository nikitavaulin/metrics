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
	metrics := make([]models.Metrics, 0, len(a.metrics)+2)
	for name, value := range a.metrics {
		val := float64(value)
		metrics = append(metrics, models.Metrics{
			ID: name, MType: models.Gauge, Value: &val,
		})
	}
	randomValue := float64(a.randomValue)
	pollCount := int64(a.pollCount)
	metrics = append(metrics,
		models.Metrics{ID: randomValueName, MType: models.Gauge, Value: &randomValue},
		models.Metrics{ID: pollCountName, MType: models.Counter, Delta: &pollCount},
	)
	a.mu.Unlock()

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/updates/", a.serverAddr)

	if err := a.sendJSON(url, metrics, client); err != nil {
		return fmt.Errorf("failed to send metrics batch: %w", err)
	}
	return nil
}

func (a *Agent) sendJSONMetric(url string, metric models.Metrics, client *http.Client) error {
	return a.sendJSON(url, metric, client)
}

func (a *Agent) sendJSON(url string, data any, client *http.Client) error {
	body, err := json.Marshal(data)
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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %s", resp.Status)
	}

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
