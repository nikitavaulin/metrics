package agent

import (
	"testing"

	"github.com/nikitavaulin/metrics/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPollRuntime(t *testing.T) {
	metrics := pollRuntime()

	assert.NotNil(t, metrics)

	expectedMetrics := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc",
	}

	for _, name := range expectedMetrics {
		value, exists := metrics[name]
		assert.True(t, exists, "Metric %s should exist", name)
		assert.IsType(t, domain.Gauge(0), value, "Metric %s should be of type Gauge", name)
	}

	assert.GreaterOrEqual(t, float64(metrics["Alloc"]), 0.0)
	assert.GreaterOrEqual(t, float64(metrics["NumGC"]), 0.0)
	assert.GreaterOrEqual(t, float64(metrics["TotalAlloc"]), 0.0)
}

func TestAgentPoll(t *testing.T) {
	agent := &Agent{
		metrics:     make(runtimeMetrics),
		randomValue: 0,
		pollCount:   0,
	}

	assert.Equal(t, 0, int(agent.pollCount))
	assert.Equal(t, domain.Gauge(0), agent.randomValue)
	assert.Empty(t, agent.metrics)

	agent.Poll()

	assert.NotEmpty(t, agent.metrics)

	assert.Equal(t, 1, int(agent.pollCount))

	prevPollCount := agent.pollCount
	agent.Poll()

	assert.Equal(t, prevPollCount+1, agent.pollCount)

	assert.NotEmpty(t, agent.metrics)
}
