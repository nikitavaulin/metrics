package main

import (
	"flag"
	"time"
)

var (
	flagServerAddr     string
	flagReportInterval time.Duration
	flagPollInterval   time.Duration
)

func parseFlags() {
	flag.StringVar(&flagServerAddr, "a", "localhost:8080", "address to run a server")
	flag.DurationVar(&flagReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&flagPollInterval, "p", 2*time.Second, "poll interval")
	flag.Parse()
}
