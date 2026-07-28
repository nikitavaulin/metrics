package main

import "github.com/nikitavaulin/metrics/internal/agent"

func main() {
	agent := agent.New("http://localhost:8080") // TODO: refactor
	agent.Run()
}
