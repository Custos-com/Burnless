package main

import (
	"os"

	"github.com/burnless/burnless/internal/cli"
)

// version is set by GoReleaser at build time via ldflags.
var version = "dev"

func main() {
	cli.SetVersion(version)
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
