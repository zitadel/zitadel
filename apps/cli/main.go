package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/zitadel/zitadel/apps/cli/cmd"
)

var version = "dev" // overridden at build time via -ldflags

func main() {
	cmd.SetVersion(version)
	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		execCmd, _, _ := root.Find(os.Args[1:])
		tip := root.Name()
		if execCmd != nil && execCmd != root {
			tip = execCmd.CommandPath()
		}

		outputFlag, _ := root.PersistentFlags().GetString("output")
		if outputFlag == "json" || isStdoutPiped() {
			je := struct {
				Error string `json:"error"`
				Code  string `json:"code,omitempty"`
				Hint  string `json:"hint"`
			}{
				Error: err.Error(),
				Hint:  fmt.Sprintf("Run '%s -h' for help.", tip),
			}
			// Extract connect error code from "[CODE] message" format.
			if msg := err.Error(); strings.HasPrefix(msg, "[") {
				if idx := strings.Index(msg, "] "); idx > 0 {
					je.Code = msg[1:idx]
					je.Error = msg[idx+2:]
				}
			}
			data, _ := json.MarshalIndent(je, "", "  ")
			fmt.Fprintln(os.Stderr, string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n\nRun '%s -h' for help.\n", err, tip)
		}
		os.Exit(1)
	}
}

func isStdoutPiped() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}
