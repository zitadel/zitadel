package main

import (
	"fmt"
	"github.com/caos/zitadel/cmd/zitadelctl/cmds"
	"os"
)

var (
	Version = "none"
)

func main() {
	rootCmd, rootValues := cmds.RootCommand(Version)
	rootCmd.Version = fmt.Sprintf("%s\n", Version)

	startCmd := cmds.StartOperator(rootValues)
	takeoffCmd := cmds.TakeoffCommand(rootValues)
	backuplistCmd := cmds.BackupListCommand(rootValues)
	restoreCmd := cmds.RestoreCommand(rootValues)
	readsecretCmd := cmds.ReadSecretCommand(rootValues)
	writesecretCmd := cmds.WriteSecretCommand(rootValues)

	rootCmd.AddCommand(
		startCmd,
		takeoffCmd,
		backuplistCmd,
		restoreCmd,
		readsecretCmd,
		writesecretCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
