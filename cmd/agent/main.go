package main

import (
	"os"

	"github.com/burnless/burnless/internal/agent"
)

var version = "dev"

func main() {
	if err := agent.NewAgentCmd(version).Execute(); err != nil {
		os.Exit(1)
	}
}
