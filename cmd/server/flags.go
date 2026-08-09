package main

import "flag"

var flagServerAddr string

func parseFlags() {
	flag.StringVar(&flagServerAddr, "a", "localhost:8080", "address to run a server")
	flag.Parse()
}
