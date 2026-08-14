package agent

import (
	"context"
	"log"
	"sync"
	"time"

	agentconfig "github.com/nikitavaulin/metrics/internal/config/agent"
	"github.com/nikitavaulin/metrics/internal/domain"
)

const (
	randomValueName string = "RandomValue"
	pollCountName   string = "PollCount"
)

type runtimeMetrics map[string]domain.Gauge

type Agent struct {
	mu             sync.Mutex
	metrics        runtimeMetrics
	randomValue    domain.Gauge
	pollCount      domain.Counter
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
}

func New(cfg *agentconfig.Config) *Agent {
	return &Agent{
		serverAddr:     "http://" + cfg.TargetServerAddr,
		pollInterval:   time.Duration(cfg.PollInterval) * time.Second,
		reportInterval: time.Duration(cfg.ReportInterval) * time.Second,
	}
}

func (a *Agent) SetSecondsIntervals(pollIvl, reportIvl int) {
	a.pollInterval = time.Duration(pollIvl) * time.Second
	a.reportInterval = time.Duration(reportIvl) * time.Second
}

func (a *Agent) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(a.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.Poll()
			case <-ctx.Done():
				return
			}
		}
	}()

	ticker := time.NewTicker(a.reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := a.Report(); err != nil {
				log.Printf("report error: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
