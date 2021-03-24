package cmds

import (
	"errors"

	"github.com/caos/orbos/pkg/kubernetes/cli"

	"github.com/caos/zitadel/operator/crtlgitops"
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

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rv, err := getRv()
		if err != nil {
			return err
		}
		defer func() {
			err = rv.ErrFunc(err)
		}()

		// TODO: Why?
		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		gitClient := rv.GitClient
		version := rv.Version

		if !rv.Gitops {
			return errors.New("restore command is only supported with the --gitops flag yet")
		}

		k8sClient, _, err := cli.Client(monitor, orbConfig, gitClient, rv.Kubeconfig, rv.Gitops)
		if err != nil && !rv.Gitops {
			return err
		}

		list, err := databases.ListBackups(monitor, gitClient)
		if err != nil {
			monitor.Error(err)
			return nil
		}

		if backup == "" {
			prompt := promptui.Select{
				Label: "Select backup to restore",
				Items: list,
			}

			_, result, err := prompt.Run()
			if err != nil {
				monitor.Error(err)
				return nil
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
			monitor.Error(errors.New("chosen backup is not existing"))
			return nil
		}

		if err := crtlgitops.Restore(monitor, gitClient, orbConfig, k8sClient, backup, rv.Gitops, &version); err != nil {
			monitor.Error(err)
		}
		return nil
	}
	return cmd
}
