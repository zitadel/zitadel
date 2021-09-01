package main

import (
	"github.com/caos/zitadel/cmd/backupctl/cmds/backup"
	"github.com/caos/zitadel/cmd/backupctl/cmds/restore"
	"os"

	"github.com/caos/orbos/mntr"

	"github.com/caos/zitadel/cmd/backupctl/cmds"
)

func main() {
	monitor := mntr.Monitor{
		OnInfo:         mntr.LogMessage,
		OnChange:       mntr.LogMessage,
		OnError:        mntr.LogError,
		OnRecoverPanic: mntr.LogPanic,
	}

	defer func() { monitor.RecoverPanic(recover()) }()

	rootCmd := cmds.RootCommand()
	backupCmd := cmds.BackupCommand(monitor)
	backupCmd.AddCommand(
		backup.S3Command(monitor),
		backup.GCSCommand(monitor),
	)
	rootCmd.AddCommand(backupCmd)
	restoreCmd := cmds.RestoreCommand(monitor)
	restoreCmd.AddCommand(
		restore.S3Command(monitor),
		restore.GCSCommand(monitor),
	)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(cmds.RestoreCommand(monitor))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
