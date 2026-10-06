package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"syscall"
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
	metrics = append(
		metrics,
		models.Metrics{ID: randomValueName, MType: models.Gauge, Value: &randomValue},
		models.Metrics{ID: pollCountName, MType: models.Counter, Delta: &pollCount},
	)
	a.mu.Unlock()

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/updates/", a.serverAddr)

	if err := a.sendRetriable(url, metrics, client); err != nil {
		return fmt.Errorf("failed to send metrics batch: %w", err)
	}
	return nil
}

func (a *Agent) sendRetriable(url string, data any, client *http.Client) error {
	err := a.sendJSON(url, data, client)
	if err != nil {
		try := 1
		duration := 1
		for err != nil && isRetriableError(err) && try <= 3 {
			time.Sleep(time.Duration(duration) * time.Second)
			err = a.sendJSON(url, data, client)
			duration += 2
			try++
		}
		if err != nil {
			return fmt.Errorf("failed to send json data: %w", err)
		}
	}

	return nil
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
		_, _ = io.Copy(io.Discard, resp.Body)
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

func isRetriableError(err error) bool {
	if err == nil {
		return false
	}
	// Ошибка отмены контекста пользователем
	if errors.Is(err, context.Canceled) {
		return false
	}
	// Превышен таймаут или дедлайн контекста
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	// Соединение отклонено (сервер упал/перезагружается)
	if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	// Временная ошибка сети или таймаут сокета
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}
