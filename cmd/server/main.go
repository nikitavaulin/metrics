package main

import (
	"context"
	"fmt"
	"log"

	dbconfig "github.com/nikitavaulin/metrics/internal/config/db"
	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
	metricshandler "github.com/nikitavaulin/metrics/internal/handler/metrics"
	pinghandler "github.com/nikitavaulin/metrics/internal/handler/ping"
	"github.com/nikitavaulin/metrics/internal/logger"
	"github.com/nikitavaulin/metrics/internal/repository"
	"github.com/nikitavaulin/metrics/internal/repository/filestorage"
	"github.com/nikitavaulin/metrics/internal/repository/memstorage"
	postgresstorage "github.com/nikitavaulin/metrics/internal/repository/postgres"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
	metricsservice "github.com/nikitavaulin/metrics/internal/service/metrics"
	pingservice "github.com/nikitavaulin/metrics/internal/service/ping"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverCfg, err := serverconfig.New()
	if err != nil {
		return fmt.Errorf("failed to get server cfg: %w", err)
	}

	if err := logger.Initialize(serverCfg.LogLevel); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	var mStorage repository.Storage

	if serverCfg.DatabaseDsn != nil {
		dbPool, err := dbconfig.New(ctx, *serverCfg.DatabaseDsn)
		if err != nil {
			return fmt.Errorf("failed to configure database pool: %w", err)
		}
		defer dbPool.Close()

		dbStorage := postgresstorage.New(dbPool)
		mStorage = dbStorage

	} else {
		memStorage := memstorage.New()
		mStorage = memStorage
	}

	fStorage := filestorage.New()

	mService := metricsservice.New(mStorage, fStorage)
	mHandler := metricshandler.New(mService)

	router := httpserver.NewRouter()
	router.RegisterRoutes(mHandler.Routes())

	dbStorage, _ := mStorage.(pingservice.Pinger)
	pHandler := pinghandler.New(pingservice.New(dbStorage))
	router.RegisterRoutes(pHandler.Routes())

	if _, ok := mStorage.(*memstorage.MemStorage); ok {
		if *serverCfg.IsNeedRestore {
			mService.LoadMetrics(ctx, serverCfg.FileStoragePath)
		}
		mService.SaveToFile(ctx, serverCfg.FileStoragePath, *serverCfg.StoreInterval)
	}

	httpServer := httpserver.New(serverCfg)
	httpServer.RegisterRouter(*router)

	logger.Log.Info("Running HTTP server", zap.String("address", serverCfg.Address))
	if err := httpServer.Run(); err != nil {
		return fmt.Errorf("running server: %w", err)
	}

	return nil
}
