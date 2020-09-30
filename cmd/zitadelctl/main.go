package main

import (
	"fmt"
	"github.com/caos/zitadel/cmd/zitadelctl/cmds"
	"os"
)

var (
	version = "none"
)

func main() {
	rootCmd, rootValues := cmds.RootCommand(version)
	rootCmd.Version = fmt.Sprintf("%s\n", version)

	startCmd := cmds.StartOperator(rootValues)
	takeoffCmd := cmds.TakeoffCommand(rootValues)

	rootCmd.AddCommand(
		startCmd,
		takeoffCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
