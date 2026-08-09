package main

import (
	"log"

	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
	metricshandler "github.com/nikitavaulin/metrics/internal/handler/metrics"
	"github.com/nikitavaulin/metrics/internal/repository/memstorage"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
	metricsservice "github.com/nikitavaulin/metrics/internal/service/metrics"
)

func main() {
	parseFlags()

	serverCfg, err := serverconfig.New(flagServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	mStorage := memstorage.New()
	mService := metricsservice.New(mStorage)
	mHandler := metricshandler.New(mService)

	router := httpserver.NewRouter()
	router.RegisterRoutes(mHandler.Routes())

	httpServer := httpserver.New(serverCfg)
	httpServer.RegisterRouter(*router)

	if err := httpServer.Run(); err != nil {
		log.Fatal(err)
	}
}
