package main

import (
	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
	httpserver "github.com/nikitavaulin/metrics/internal/server"
)

func main() {
	serverCfg, err := serverconfig.New()
	if err != nil {
		panic(err)
	}

	httpServer := httpserver.New(serverCfg)

	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
