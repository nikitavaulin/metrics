package agent

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentGetPath(t *testing.T) {
	agent := &Agent{serverAddr: "http://localhost:8080"}

	tests := []struct {
		name     string
		mtype    string
		metric   string
		value    string
		expected string
	}{
		{
			name:     "gauge metric",
			mtype:    models.Gauge,
			metric:   "Alloc",
			value:    "123.45",
			expected: "http://localhost:8080/update/gauge/Alloc/123.45",
		},
		{
			name:     "counter metric",
			mtype:    models.Counter,
			metric:   "PollCount",
			value:    "42",
			expected: "http://localhost:8080/update/counter/PollCount/42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.getPath(tt.mtype, tt.metric, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAgentReport(t *testing.T) {
	var requestsCount int
	var received []models.Metrics

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsCount++
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/updates/", r.URL.Path)
		assert.Equal(t, contentTypeJSON, r.Header.Get("Content-Type"))
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"))
		reader, err := gzip.NewReader(r.Body)
		if !assert.NoError(t, err) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer reader.Close()
		if err := json.NewDecoder(reader).Decode(&received); !assert.NoError(t, err) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	agent := &Agent{
		metrics: runtimeMetrics{
			"Alloc":     123.45,
			"HeapAlloc": 456.78,
			"NumGC":     10.0,
		},
		randomValue: 0.987,
		pollCount:   5,
		serverAddr:  server.URL,
	}

	err := agent.Report()
	assert.NoError(t, err)
	require.Equal(t, 1, requestsCount)
	require.Len(t, received, 5)
	got := make(map[string]models.Metrics, len(received))
	for _, metric := range received {
		got[metric.ID] = metric
	}
	for name, expected := range map[string]float64{
		"Alloc": 123.45, "HeapAlloc": 456.78, "NumGC": 10, randomValueName: 0.987,
	} {
		metric := got[name]
		assert.Equal(t, models.Gauge, metric.MType)
		require.NotNil(t, metric.Value)
		assert.Equal(t, expected, *metric.Value)
		assert.Nil(t, metric.Delta)
	}
	metric := got[pollCountName]
	assert.Equal(t, models.Counter, metric.MType)
	require.NotNil(t, metric.Delta)
	assert.Equal(t, int64(5), *metric.Delta)
	assert.Nil(t, metric.Value)
}

func TestAgentReportServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	agent := &Agent{serverAddr: server.URL}
	require.ErrorContains(t, agent.Report(), "500")
}
