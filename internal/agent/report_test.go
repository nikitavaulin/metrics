package agent

import (
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nikitavaulin/metrics/internal/domain"
	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendJSONHash(t *testing.T) {
	key, emptyKey := "secret", ""
	for _, tt := range []struct {
		name string
		key  *string
	}{
		{name: "without key"},
		{name: "with key", key: &key},
		{name: "with empty key", key: &emptyKey},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reader, err := gzip.NewReader(r.Body)
				if !assert.NoError(t, err) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				defer reader.Close()
				body, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"metric":"Alloc"}`, string(body))
				if tt.key == nil {
					assert.Nil(t, r.Header.Values("HashSHA256"))
				} else {
					h := hmac.New(sha256.New, []byte(*tt.key))
					h.Write(body)
					assert.Equal(t, hex.EncodeToString(h.Sum(nil)), r.Header.Get("HashSHA256"))
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			a := &Agent{hashingKey: tt.key}
			require.NoError(t, a.sendJSON(context.Background(), server.URL, map[string]string{"metric": "Alloc"}, server.Client()))
		})
	}
}

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
		metricsRuntime: gaugeMetrics{
			"Alloc":     123.45,
			"HeapAlloc": 456.78,
			"NumGC":     10.0,
		},
		randomValue: 0.987,
		pollCount:   5,
		serverAddr:  server.URL,
	}

	err := agent.Report(context.Background())
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
	require.ErrorContains(t, agent.Report(context.Background()), "500")
}

func TestAgentReportBatches(t *testing.T) {
	for _, count := range []int{batchSize - 1, batchSize, batchSize + 1, 2*batchSize + 3} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			var sizes []int
			var received []models.Metrics
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reader, err := gzip.NewReader(r.Body)
				if !assert.NoError(t, err) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				defer reader.Close()
				var batch []models.Metrics
				if !assert.NoError(t, json.NewDecoder(reader).Decode(&batch)) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				sizes = append(sizes, len(batch))
				received = append(received, batch...)
			}))
			defer server.Close()
			a := &Agent{serverAddr: server.URL, metricsGopsutil: make(gaugeMetrics)}
			for i := 0; i < count-2; i++ {
				a.metricsGopsutil[fmt.Sprintf("metric%d", i)] = domain.Gauge(i)
			}
			expected := a.snapshotMetrics()
			require.NoError(t, a.Report(context.Background()))
			require.Len(t, sizes, (count+batchSize-1)/batchSize)
			for i, size := range sizes {
				assert.Equal(t, min(batchSize, count-i*batchSize), size)
			}
			assert.ElementsMatch(t, expected, received)
		})
	}
}

func TestRunWorkerLimit(t *testing.T) {
	for _, limit := range []int{1, 3} {
		t.Run(fmt.Sprint(limit), func(t *testing.T) {
			entered := make(chan struct{}, 100)
			release := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				entered <- struct{}{}
				<-release
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			defer close(release)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			a := &Agent{
				serverAddr: server.URL, rateLimit: limit,
				pollInterval: time.Hour, reportInterval: time.Millisecond,
				metricsRuntime: make(gaugeMetrics),
			}
			for i := 0; i < (limit+1)*batchSize; i++ {
				a.metricsRuntime[fmt.Sprintf("metric%d", i)] = domain.Gauge(i)
			}
			done := make(chan struct{})
			go func() {
				a.Run(ctx)
				close(done)
			}()
			for i := 0; i < limit; i++ {
				select {
				case <-entered:
				case <-time.After(3 * time.Second):
					t.Fatal("workers did not start concurrent requests")
				}
			}
			select {
			case <-entered:
				t.Fatal("concurrent requests exceeded rate limit")
			case <-time.After(50 * time.Millisecond):
			}
			cancel()
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fatal("Run did not stop after cancellation")
			}
			require.Empty(t, entered)
		})
	}
}
