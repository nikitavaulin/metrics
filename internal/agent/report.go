package agent

import (
	"fmt"
	"io"
	"log"
	"net/http"

	models "github.com/nikitavaulin/metrics/internal/model"
)

const (
	contentType string = "text/plain"
)

func (a *Agent) Report() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// runtime
	for name, value := range a.metrics {
		url := a.getPath(models.Gauge, name, fmt.Sprint(value))
		resp, err := http.Post(url, contentType, nil)
		logResponse(resp)
		if err != nil {
			return fmt.Errorf("failed to POST metric: %w", err)
		}
	}

	// random
	url := a.getPath(models.Gauge, randomValueName, fmt.Sprint(a.randomValue))
	resp, err := http.Post(url, contentType, nil)
	logResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to POST random value metric: %w", err)
	}

	// counter
	url = a.getPath(models.Counter, pollCountName, fmt.Sprint(a.pollCount))
	resp, err = http.Post(url, contentType, nil)
	logResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to POST poll count: %w", err)
	}

	return nil
}

func (a *Agent) getPath(mtype, name, value string) string {
	return fmt.Sprintf("%s/update/%s/%s/%s", a.serverAddr, mtype, name, value)
}

func logResponse(resp *http.Response) {
	if resp != nil {
		var errMsg string = "null"
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
