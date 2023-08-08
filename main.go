package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/cmd"
)

func main() {
	arg := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin, arg, nil)
	cobra.CheckErr(rootCmd.Execute())
}
