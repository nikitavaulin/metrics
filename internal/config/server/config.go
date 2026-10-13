package serverconfig

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/nikitavaulin/metrics/internal/validation"
)

const (
	defaultAddress         = "localhost:8080"
	defaultLogLvl          = "Info"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "metrics.json"
	defaultIsNeedRestore   = false

	emptyDBConnStr = "empty"
)

type Config struct {
	Address         string  `env:"ADDRESS"`
	LogLevel        string  `env:"LOG_LEVEL"`
	StoreInterval   *int    `env:"STORE_INTERVAL"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH"`
	IsNeedRestore   *bool   `env:"RESTORE"`
	DatabaseDsn     *string `env:"DATABASE_DSN"`
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
	var (
		address         string
		logLevel        string
		storeInterval   int
		fileStoragePath string
		isNeedRestore   bool
		dbConnString    string
	)
	flag.StringVar(&address, "a", defaultAddress, "address to run a server")
	flag.StringVar(&logLevel, "l", defaultLogLvl, "logger level")
	flag.IntVar(&storeInterval, "i", defaultStoreInterval, "metrics store interval (int seconds)")
	flag.StringVar(&fileStoragePath, "f", defaultFileStoragePath, "file storage destination (path)")
	flag.BoolVar(&isNeedRestore, "r", defaultIsNeedRestore, "restore metrics after start")
	flag.StringVar(&dbConnString, "d", emptyDBConnStr, "postgres db connection string")

	flag.Parse()

	if cfg.Address == "" {
		cfg.Address = address
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = logLevel
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = &storeInterval
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = fileStoragePath
	}
	if cfg.IsNeedRestore == nil {
		cfg.IsNeedRestore = &isNeedRestore
	}
	if cfg.DatabaseDsn == nil && dbConnString != emptyDBConnStr {
		cfg.DatabaseDsn = &dbConnString
	}
}
