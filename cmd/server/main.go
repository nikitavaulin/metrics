package main

import (
	"context"
	"fmt"
	"log"

	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
	metricshandler "github.com/nikitavaulin/metrics/internal/handler/metrics"
	"github.com/nikitavaulin/metrics/internal/logger"
	"github.com/nikitavaulin/metrics/internal/repository/filestorage"
	"github.com/nikitavaulin/metrics/internal/repository/memstorage"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
	metricsservice "github.com/nikitavaulin/metrics/internal/service/metrics"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverCfg, err := serverconfig.New()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get server cfg: %w", err))
	}

	if err := logger.Initialize(serverCfg.LogLevel); err != nil {
		log.Fatal(fmt.Errorf("failed to initialize logger: %w", err))
	}

	fStorage := filestorage.New()

	mStorage := memstorage.New()
	mService := metricsservice.New(mStorage, fStorage)
	mHandler := metricshandler.New(mService)

	if *serverCfg.IsNeedRestore {
		mService.LoadMetrics(serverCfg.FileStoragePath)
	}
	mService.SaveToFile(ctx, serverCfg.FileStoragePath, *serverCfg.StoreInterval)

	router := httpserver.NewRouter()
	router.RegisterRoutes(mHandler.Routes())

	httpServer := httpserver.New(serverCfg)
	httpServer.RegisterRouter(*router)

	logger.Log.Info("Running HTTP server", zap.String("address", serverCfg.Address))
	if err := httpServer.Run(); err != nil {
		logger.Log.Fatal("Running server error", zap.Error(err))
	}
}
