package cmds

import (
	"github.com/caos/orbos/mntr"
	"github.com/spf13/cobra"
)

func BackupCommand(monitor mntr.Monitor) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "Backup",
			Short: "Backup to storage",
			Long:  "Backup to storage",
		}
	)

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		monitor.Info("Please select to which storage you want to backup to")
		return nil
	}
	return cmd
}
