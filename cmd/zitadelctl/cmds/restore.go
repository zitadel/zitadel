package cmds

import (
	"errors"

	"github.com/caos/zitadel/pkg/zitadel"

	"github.com/caos/orbos/mntr"

	"github.com/caos/orbos/pkg/kubernetes/cli"

	"github.com/caos/zitadel/pkg/databases"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func RestoreCommand(getRv GetRootValues) *cobra.Command {
	var (
		backup string
		cmd    = &cobra.Command{
			Use:   "restore",
			Short: "Restore from backup",
			Long:  "Restore from backup",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backup, "backup", "", "Backup used for db restore")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		rv := getRv("restore", map[string]interface{}{"backup": backup}, "")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		// TODO: Why?
		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		gitClient := rv.GitClient
		version := rv.Version

		k8sClient, err := cli.Client(monitor, orbConfig, gitClient, rv.Kubeconfig, rv.Gitops, true)
		if err != nil {
			return err
		}

		list := make([]string, 0)
		if rv.Gitops {
			listT, err := databases.GitOpsListBackups(monitor, gitClient, k8sClient)
			if err != nil {
				return err
			}
			list = listT
		} else {
			listT, err := databases.CrdListBackups(monitor, k8sClient)
			if err != nil {
				return err
			}
			list = listT
		}

		if backup == "" {
			prompt := promptui.Select{
				Label: "Select backup to restore",
				Items: list,
			}

			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			backup = result
		}
		existing := false
		for _, listedBackup := range list {
			if listedBackup == backup {
				existing = true
			}
		}

		if !existing {
			return mntr.ToUserError(errors.New("chosen backup is not existing"))
		}

		if rv.Gitops {
			if err := zitadel.GitOpsClearMigrateRestore(monitor, gitClient, orbConfig, k8sClient, backup, &version); err != nil {
				return err
			}
		} else {
			if err := zitadel.CrdClearMigrateRestore(monitor, k8sClient, backup, &version); err != nil {
				return err
			}

		}
		return nil
	}
	return cmd
}
