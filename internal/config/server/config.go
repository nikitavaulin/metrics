package serverconfig

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/nikitavaulin/metrics/internal/validation"
)

const (
	defaultAddress = "localhost:8080"
	defaultLogLvl  = "Info"
)

type Config struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	cfg.parseFlags()

	if err := validation.ValidateServerAddress(cfg.Address); err != nil {
		return nil, fmt.Errorf("invalid server address: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) parseFlags() {
	var flagCfg Config
	flag.StringVar(&flagCfg.Address, "a", defaultAddress, "address to run a server")
	flag.StringVar(&flagCfg.LogLevel, "l", defaultLogLvl, "logger level")
	flag.Parse()

	if cfg.Address == "" {
		cfg.Address = flagCfg.Address
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = flagCfg.LogLevel
	}
}
