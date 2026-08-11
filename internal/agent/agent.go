package agent

import (
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

func (a *Agent) Run() {
	go func() {
		for {
			time.Sleep(a.pollInterval)
			a.Poll()
		}
	}()

	for {
		time.Sleep(a.reportInterval)
		if err := a.Report(); err != nil {
			log.Printf("report error: %v", err)
		}
	}
}
