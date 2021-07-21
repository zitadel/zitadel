package cmds

import (
	"errors"
	"fmt"
	"sort"

	"github.com/caos/orbos/mntr"

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

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		rv := getRv("backuplist", nil, "")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		gitClient := rv.GitClient

		if !rv.Gitops {
			return mntr.ToUserError(errors.New("backuplist command is only supported with the --gitops flag yet"))
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
