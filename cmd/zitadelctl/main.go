package main

import (
	"fmt"
	"os"

	"github.com/caos/orbos/mntr"

	"github.com/caos/zitadel/cmd/zitadelctl/cmds"
)

var (
	Version            = "unknown"
	githubClientID     = "none"
	githubClientSecret = "none"
)

func main() {
	monitor := mntr.Monitor{
		OnInfo:         mntr.LogMessage,
		OnChange:       mntr.LogMessage,
		OnError:        mntr.LogError,
		OnRecoverPanic: mntr.LogPanic,
	}

	defer monitor.RecoverPanic()

	rootCmd, rootValues := cmds.RootCommand(Version, monitor)
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
		cmds.ConfigCommand(rootValues, githubClientID, githubClientSecret),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
