package cmds

import (
	"errors"

	"github.com/caos/orbos/pkg/kubernetes/cli"

	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/crtlgitops"
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
		version := rv.Version

		if !rv.Gitops {
			return errors.New("backup command is only supported with the --gitops flag yet")
		}

		k8sClient, _, err := cli.Client(monitor, orbConfig, gitClient, rv.Kubeconfig, rv.Gitops)
		if err != nil && !rv.Gitops {
			return err
		}

		found, err := api.ExistsDatabaseYml(gitClient)
		if err != nil {
			return err
		}
		if found {

			if err := crtlgitops.Backup(
				monitor,
				orbConfig.Path,
				k8sClient,
				backup,
				&version,
			); err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}
