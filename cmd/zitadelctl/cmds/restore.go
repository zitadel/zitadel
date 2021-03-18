package cmds

import (
	"errors"
	"io/ioutil"

	"github.com/caos/zitadel/operator/crtlgitops"
	"github.com/caos/zitadel/operator/helpers"

	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/pkg/databases"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func RestoreCommand(getRv GetRootValues) *cobra.Command {
	var (
		backup     string
		kubeconfig string
		gitOpsMode bool
		cmd        = &cobra.Command{
			Use:   "restore",
			Short: "Restore from backup",
			Long:  "Restore from backup",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backup, "backup", "", "Backup used for db restore")
	flags.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "Kubeconfig for ZITADEL operator deployment")
	flags.BoolVar(&gitOpsMode, "gitops", false, "defines if the operator should run in gitops mode")

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
		version := rv.Version

		kubeconfig = helpers.PruneHome(kubeconfig)

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			monitor.Error(err)
			return nil
		}

		if err := gitClient.Clone(); err != nil {
			monitor.Error(err)
			return nil
		}

		value, err := ioutil.ReadFile(kubeconfig)
		if err != nil {
			monitor.Error(err)
			return nil
		}
		kubeconfigStr := string(value)

		k8sClient := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
		if k8sClient.Available() {
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

			if err := crtlgitops.Restore(monitor, gitClient, orbConfig, k8sClient, backup, gitOpsMode, &version); err != nil {
				monitor.Error(err)
			}
			return nil
		}
		return nil
	}
	return cmd
}
