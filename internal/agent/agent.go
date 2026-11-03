package agent

import (
	"context"
	"sync"
	"time"

	agentconfig "github.com/nikitavaulin/metrics/internal/config/agent"
	"github.com/nikitavaulin/metrics/internal/domain"
	"github.com/nikitavaulin/metrics/internal/logger"
	models "github.com/nikitavaulin/metrics/internal/model"
	"go.uber.org/zap"
)

const (
	randomValueName string = "RandomValue"
	pollCountName   string = "PollCount"
)

type gaugeMetrics map[string]domain.Gauge

type Agent struct {
	mu              sync.Mutex
	metricsRuntime  gaugeMetrics
	metricsGopsutil gaugeMetrics
	randomValue     domain.Gauge
	pollCount       domain.Counter
	serverAddr      string
	hashingKey      *string
	pollInterval    time.Duration
	reportInterval  time.Duration
	rateLimit       int
}

func New(cfg *agentconfig.Config) *Agent {
	rateLimit := 1
	if cfg.RateLimit != nil {
		rateLimit = *cfg.RateLimit
	}
	return &Agent{
		rateLimit:      rateLimit,
		hashingKey:     cfg.Key,
		serverAddr:     "http://" + cfg.TargetServerAddr,
		pollInterval:   time.Duration(*cfg.PollInterval) * time.Second,
		reportInterval: time.Duration(*cfg.ReportInterval) * time.Second,
	}
}

func (a *Agent) SetSecondsIntervals(pollIvl, reportIvl int) {
	a.pollInterval = time.Duration(pollIvl) * time.Second
	a.reportInterval = time.Duration(reportIvl) * time.Second
}

func (a *Agent) Run(ctx context.Context) {
	var wg sync.WaitGroup
	jobs := make(chan []models.Metrics)

	// start workers
	for range a.rateLimit {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.reportWorker(ctx, jobs)
		}()
	}

	defer wg.Wait()
	defer close(jobs)

	// collecting
	wg.Add(2)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(a.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.poll()
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(a.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := a.pollGopsUtil(); err != nil {
					logger.Log.Error("poll gopsutil", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// snapshot for report
	ticker := time.NewTicker(a.reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := a.enqueueMetrics(ctx, jobs); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
