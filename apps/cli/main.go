package main

import (
	"os"

	"github.com/zitadel/zitadel/apps/cli/cmd"
)

var version = "dev" // overridden at build time via -ldflags

func main() {
	cmd.SetVersion(version)
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
