package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/zitadel/zitadel/apps/cli/cmd"
	"github.com/zitadel/zitadel/apps/cli/internal/client"
	"github.com/zitadel/zitadel/apps/cli/internal/output"
)

var version = "dev" // overridden at build time via -ldflags

func main() {
	cmd.SetVersion(version)
	client.UserAgent = fmt.Sprintf("zitadel-cli/%s (%s/%s)", version, runtime.GOOS, runtime.GOARCH)
	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		execCmd, _, _ := root.Find(os.Args[1:])
		tip := root.Name()
		if execCmd != nil && execCmd != root {
			tip = execCmd.CommandPath()
		}

		outputFlag, _ := root.PersistentFlags().GetString("output")
		output.HandleError(err, outputFlag, tip, os.Stderr)
		os.Exit(1)
	}
}
