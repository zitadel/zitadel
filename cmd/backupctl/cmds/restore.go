package cmds

import (
	"github.com/caos/orbos/mntr"
	"github.com/spf13/cobra"
)

func RestoreCommand(monitor mntr.Monitor) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "Restore",
			Short: "Restore from storage",
			Long:  "Restore from storage",
		}
	)

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		monitor.Info("Please select from which storage you want to restore from")
		return nil
	}
	return cmd
}
