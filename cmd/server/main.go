package main

import (
	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
	metricshandler "github.com/nikitavaulin/metrics/internal/handler/metrics"
	"github.com/nikitavaulin/metrics/internal/repository/memstorage"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
	metricsservice "github.com/nikitavaulin/metrics/internal/service/metrics"
)

func main() {
	serverCfg, err := serverconfig.New()
	if err != nil {
		panic(err)
	}

	mStorage := memstorage.New()
	mService := metricsservice.New(mStorage)
	mHandler := metricshandler.New(mService)

	router := httpserver.NewRouter()
	router.RegisterRoutes(mHandler.Routes())

	httpServer := httpserver.New(serverCfg)
	httpServer.RegisterRouter(*router)

	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
