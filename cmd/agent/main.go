package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikitavaulin/metrics/internal/agent"
	agentconfig "github.com/nikitavaulin/metrics/internal/config/agent"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	agentCfg, err := agentconfig.New()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get agent config: %w", err))
	}

	agent := agent.New(agentCfg)
	agent.Run(ctx)
}
