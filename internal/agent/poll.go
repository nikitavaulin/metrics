package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/nikitavaulin/metrics/internal/domain"
)

func (a *Agent) Poll() {
	a.mu.Lock()
	a.metrics = pollRuntime()
	a.randomValue = getRandom()
	a.pollCount++
	a.mu.Unlock()
}

func getRandom() domain.Gauge {
	value := rand.Float64()
	return domain.Gauge(value)
}

func pollRuntime() runtimeMetrics {
	metrics := make(runtimeMetrics)

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics["Alloc"] = domain.Gauge(stats.Alloc)
	metrics["BuckHashSys"] = domain.Gauge(stats.BuckHashSys)
	metrics["Frees"] = domain.Gauge(stats.Frees)
	metrics["GCCPUFraction"] = domain.Gauge(stats.GCCPUFraction)
	metrics["GCSys"] = domain.Gauge(stats.GCSys)
	metrics["HeapAlloc"] = domain.Gauge(stats.HeapAlloc)
	metrics["HeapIdle"] = domain.Gauge(stats.HeapIdle)
	metrics["HeapInuse"] = domain.Gauge(stats.HeapInuse)
	metrics["HeapObjects"] = domain.Gauge(stats.HeapObjects)
	metrics["HeapReleased"] = domain.Gauge(stats.HeapReleased)
	metrics["HeapSys"] = domain.Gauge(stats.HeapSys)
	metrics["LastGC"] = domain.Gauge(stats.LastGC)
	metrics["Lookups"] = domain.Gauge(stats.Lookups)
	metrics["MCacheInuse"] = domain.Gauge(stats.MCacheInuse)
	metrics["MCacheSys"] = domain.Gauge(stats.MCacheSys)
	metrics["MSpanInuse"] = domain.Gauge(stats.MSpanInuse)
	metrics["MSpanSys"] = domain.Gauge(stats.MSpanSys)
	metrics["Mallocs"] = domain.Gauge(stats.Mallocs)
	metrics["NextGC"] = domain.Gauge(stats.NextGC)
	metrics["NumForcedGC"] = domain.Gauge(stats.NumForcedGC)
	metrics["NumGC"] = domain.Gauge(stats.NumGC)
	metrics["OtherSys"] = domain.Gauge(stats.OtherSys)
	metrics["PauseTotalNs"] = domain.Gauge(stats.PauseTotalNs)
	metrics["StackInuse"] = domain.Gauge(stats.StackInuse)
	metrics["StackSys"] = domain.Gauge(stats.StackSys)
	metrics["Sys"] = domain.Gauge(stats.Sys)
	metrics["TotalAlloc"] = domain.Gauge(stats.TotalAlloc)

	return metrics
}
