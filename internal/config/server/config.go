package serverconfig

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/nikitavaulin/metrics/internal/validation"
)

const defaultAddress = "localhost:8080"

type Config struct {
	Address string `env:"ADDRESS"`
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
	flag.Parse()

	if cfg.Address == "" {
		cfg.Address = flagCfg.Address
	}
}
