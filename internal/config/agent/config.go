package agentconfig

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/nikitavaulin/metrics/internal/validation"
)

const (
	defaultTargetServerAddr string = "localhost:8080"
	defaultReportInterval   int    = 10
	defaultPollInterval     int    = 2
)

type Config struct {
	TargetServerAddr string `env:"ADDRESS"`
	ReportInterval   *int   `env:"REPORT_INTERVAL"`
	PollInterval     *int   `env:"POLL_INTERVAL"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env var: %w", err)
	}

	cfg.parseFlags()

	if err := validation.ValidateServerAddress(cfg.TargetServerAddr); err != nil {
		return nil, fmt.Errorf("invalid server address: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) parseFlags() {
	var (
		targetServerAddr string
		reportInterval   int
		pollInterval     int
	)

	flag.StringVar(&targetServerAddr, "a", defaultTargetServerAddr, "address to run a server")
	flag.IntVar(&reportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.Parse()

	if cfg.TargetServerAddr == "" {
		cfg.TargetServerAddr = targetServerAddr
	}
	if cfg.ReportInterval == nil {
		cfg.ReportInterval = &reportInterval
	}
	if cfg.PollInterval == nil {
		cfg.PollInterval = &pollInterval
	}
}
