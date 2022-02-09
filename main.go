package main

import (
	"os"

	"github.com/caos/zitadel/cmd"
	"github.com/spf13/cobra"
)

func main() {
	args := os.Args[1:]
	rootCmd := cmd.New(os.Stdout, os.Stdin /*, int(os.Stdin.Fd())*/, args)
	cobra.CheckErr(rootCmd.Execute())
}
