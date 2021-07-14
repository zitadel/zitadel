package cmds

import (
	"errors"
	"fmt"
	"sort"

	"github.com/caos/zitadel/pkg/databases"
	"github.com/spf13/cobra"
)

func BackupListCommand(getRv GetRootValues) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "backuplist",
			Short: "Get a list of all backups",
			Long:  "Get a list of all backups",
		}
	)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rv, err := getRv("backuplist", nil, "")
		if err != nil {
			return err
		}
		defer func() {
			err = rv.ErrFunc(err)
		}()

		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		gitClient := rv.GitClient

		if !rv.Gitops {
			return errors.New("backuplist command is only supported with the --gitops flag yet")
		}

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			monitor.Error(err)
			return nil
		}

		if err := gitClient.Clone(); err != nil {
			monitor.Error(err)
			return nil
		}

		backups, err := databases.ListBackups(monitor, gitClient)
		if err != nil {
			monitor.Error(err)
			return nil
		}

		sort.Slice(backups, func(i, j int) bool {
			return backups[i] > backups[j]
		})
		for _, backup := range backups {
			fmt.Println(backup)
		}
		return nil
	}
	return cmd
}
