package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/stretchr/testify/assert"
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsCount++
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
	assert.Equal(t, 5, requestsCount)
}
