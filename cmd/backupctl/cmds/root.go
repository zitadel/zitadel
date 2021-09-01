package cmds

import (
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "backupctl [flags]",
		Short:   "Interact with backups ",
		Long:    `Interact with backups`,
		Example: ``,
	}
}
