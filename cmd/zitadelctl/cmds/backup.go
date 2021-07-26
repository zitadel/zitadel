package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes/cli"
	"github.com/caos/zitadel/pkg/databases"

	"github.com/spf13/cobra"
)

func BackupCommand(getRv GetRootValues) *cobra.Command {
	var (
		backup string
		cmd    = &cobra.Command{
			Use:   "backup",
			Short: "Instant backup",
			Long:  "Instant backup",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backup, "backup", "", "Name used for backup folder")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
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

		k8sClient, err := cli.Client(monitor, orbConfig, gitClient, rv.Kubeconfig, rv.Gitops, true)
		if err != nil {
			return err
		}

		if rv.Gitops {
			if err := databases.GitOpsInstantBackup(
				monitor,
				k8sClient,
				gitClient,
				backup,
			); err != nil {
				return err
			}

		} else {
			if err := databases.CrdInstantBackup(
				monitor,
				k8sClient,
				backup,
			); err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
