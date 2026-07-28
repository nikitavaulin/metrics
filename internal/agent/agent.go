package agent

import (
	"log"
	"sync"
	"time"

	"github.com/nikitavaulin/metrics/internal/domain"
)

const (
	pollInterval    time.Duration = 2 * time.Second
	reportInterval  time.Duration = 10 * time.Second
	randomValueName string        = "RandomValue"
	pollCountName   string        = "PollCount"
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
		serverAddr: serverAddr,
	}
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
