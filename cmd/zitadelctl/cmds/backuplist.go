package cmds

import (
	"fmt"

	"github.com/caos/orbos/pkg/databases"
	"github.com/spf13/cobra"
)

func BackupListCommand(rv RootValues) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "backuplist",
			Short: "Get a list of all backups",
			Long:  "Get a list of all backups",
		}
	)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, monitor, orbConfig, gitClient, _, errFunc := rv()
		if errFunc != nil {
			return errFunc(cmd)
		}

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			return err
		}

		if err := gitClient.Clone(); err != nil {
			return err
		}

		backups, err := databases.ListBackups(monitor, gitClient)
		if err != nil {
			return err
		}

		for _, backup := range backups {
			fmt.Println(backup)
		}
		return nil
	}
	return cmd
}
