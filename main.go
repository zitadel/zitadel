package main

import (
	"os"

	"github.com/zitadel/zitadel/cmd"
)

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, args, nil)
	if err := rootCmd.Execute(); err != nil {
		// error is logged by the command itself
		os.Exit(1)
	}
}
