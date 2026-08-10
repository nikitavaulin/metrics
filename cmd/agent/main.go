package main

import (
	"fmt"
	"log"

	"github.com/nikitavaulin/metrics/internal/agent"
	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
)

func main() {
	parseFlags()

	if err := serverconfig.ValidateServerAddress(flagServerAddr); err != nil {
		log.Fatal(fmt.Errorf("invalid server address: %w", err))
	}

	agent := agent.New(flagServerAddr)
	agent.SetSecondsIntervals(flagPollInterval, flagReportInterval)
	agent.Run()
}
