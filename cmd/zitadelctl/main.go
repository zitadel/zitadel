package main

import (
	"fmt"
	"github.com/caos/zitadel/cmd/zitadelctl/cmds"
	"os"
)

var (
	Version = "unknown"
)

func main() {
	rootCmd, rootValues := cmds.RootCommand(Version)
	rootCmd.Version = fmt.Sprintf("%s\n", Version)

	rootCmd.AddCommand(
		cmds.StartOperator(rootValues),
		cmds.TakeoffCommand(rootValues),
		cmds.BackupListCommand(rootValues),
		cmds.RestoreCommand(rootValues),
		cmds.ReadSecretCommand(rootValues),
		cmds.WriteSecretCommand(rootValues),
		cmds.BackupCommand(rootValues),
		cmds.StartDatabase(rootValues),
		cmds.TeardownCommand(rootValues),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
