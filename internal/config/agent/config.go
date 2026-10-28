package agentconfig

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/nikitavaulin/metrics/internal/validation"
)

const (
	defaultTargetServerAddr string = "localhost:8080"
	defaultReportInterval   int    = 10
	defaultPollInterval     int    = 2
)

type Config struct {
	Key              *string `env:"KEY"`
	TargetServerAddr string  `env:"ADDRESS"`
	ReportInterval   *int    `env:"REPORT_INTERVAL"`
	PollInterval     *int    `env:"POLL_INTERVAL"`
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
		key              string
	)

	flag.StringVar(&targetServerAddr, "a", defaultTargetServerAddr, "address to run a server")
	flag.IntVar(&reportInterval, "r", defaultReportInterval, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", defaultPollInterval, "poll interval in seconds")
	flag.StringVar(&key, "k", "", "key for signing metrics")
	flag.Parse()

	if envKey, ok := os.LookupEnv("KEY"); ok {
		cfg.Key = &envKey
	}
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "k" && cfg.Key == nil {
			cfg.Key = &key
		}
	})
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
