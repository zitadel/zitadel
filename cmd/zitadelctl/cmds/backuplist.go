package cmds

import (
	"fmt"
	"github.com/caos/orbos/pkg/kubernetes/cli"
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
		rv, err := getRv()
		if err != nil {
			return err
		}
		defer func() {
			err = rv.ErrFunc(err)
		}()

		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		gitClient := rv.GitClient

		backups := make([]string, 0)
		if rv.Gitops {
			if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
				monitor.Error(err)
				return nil
			}

			if err := gitClient.Clone(); err != nil {
				monitor.Error(err)
				return nil
			}

			backupsT, err := databases.GitOpsListBackups(monitor, gitClient)
			if err != nil {
				monitor.Error(err)
				return nil
			}
			backups = backupsT
		} else {
			k8sClient, _, err := cli.Client(monitor, orbConfig, rv.GitClient, rv.Kubeconfig, rv.Gitops)
			if err != nil {
				return err
			}

			backupsT, err := databases.CrdListBackups(monitor, k8sClient)
			if err != nil {
				monitor.Error(err)
				return nil
			}
			backups = backupsT
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
