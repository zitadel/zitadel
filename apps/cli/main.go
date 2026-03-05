package main

import (
	"fmt"
	"os"

	"github.com/zitadel/zitadel/apps/cli/cmd"
)

var version = "dev" // overridden at build time via -ldflags

func main() {
	cmd.SetVersion(version)
	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		// Determine which (sub)command was attempted so the tip is specific.
		execCmd, _, _ := root.Find(os.Args[1:])
		tip := root.Name()
		if execCmd != nil && execCmd != root {
			tip = execCmd.CommandPath()
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n\nRun '%s -h' for help.\n", err, tip)
		os.Exit(1)
	}
}
