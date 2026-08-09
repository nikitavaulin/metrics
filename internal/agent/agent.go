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

var (
	pollInterval   time.Duration = 2 * time.Second
	reportInterval time.Duration = 10 * time.Second
)

type runtimeMetrics map[string]domain.Gauge

type Agent struct {
	mu          sync.Mutex
	metrics     runtimeMetrics
	randomValue domain.Gauge
	pollCount   domain.Counter
	serverAddr  string
}

func New(serverAddr string) *Agent {
	return &Agent{
		serverAddr: "http://" + serverAddr,
	}
}

func (a *Agent) SetIntervals(pollIvl, reportIvl time.Duration) {
	pollInterval = pollIvl
	reportInterval = reportIvl
}

func (a *Agent) Run() {
	var wg sync.WaitGroup

	wg.Go(func() {
		for {
			time.Sleep(pollInterval)
			a.Poll()
		}
	})

	wg.Go(func() {
		for {
			time.Sleep(reportInterval)
			if err := a.Report(); err != nil {
				log.Printf("report error: %v", err)
			}
		}
	})

	wg.Wait()
}
