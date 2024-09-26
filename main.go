package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/cmd"
)

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, args)
	cobra.CheckErr(rootCmd.Execute())
}
