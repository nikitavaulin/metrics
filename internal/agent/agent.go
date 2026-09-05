package agent

import (
	"context"
	"log"
	"sync"
	"time"

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

func New(serverAddr string) *Agent {
	return &Agent{
		serverAddr:     "http://" + serverAddr,
		pollInterval:   2 * time.Second,
		reportInterval: 10 * time.Second,
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
