package cmds

import (
	"errors"
	"io/ioutil"

	"github.com/caos/zitadel/operator/helpers"

	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/start"
	"github.com/caos/zitadel/pkg/databases"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func RestoreCommand(rv RootValues) *cobra.Command {
	var (
		backup     string
		kubeconfig string
		cmd        = &cobra.Command{
			Use:   "restore",
			Short: "Restore from backup",
			Long:  "Restore from backup",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&backup, "backup", "", "Backup used for db restore")
	flags.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "Kubeconfig for ZITADEL operator deployment")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, monitor, orbConfig, gitClient, version, errFunc, err := rv()
		if err != nil {
			return err
		}
		defer func() {
			err = errFunc(err)
		}()

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

			if err := start.Restore(monitor, gitClient, orbConfig, k8sClient, backup, &version); err != nil {
				monitor.Error(err)
			}
			return nil
		}
		return nil
	}
	return cmd
}
