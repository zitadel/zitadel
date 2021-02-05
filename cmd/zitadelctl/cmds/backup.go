package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/start"
	"github.com/spf13/cobra"
	"io/ioutil"
)

func BackupCommand(getRv GetRootValues) *cobra.Command {
	var (
		kubeconfig string
		backup     string
		cmd        = &cobra.Command{
			Use:   "backup",
			Short: "Instant backup",
			Long:  "Instant backup",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "Kubeconfig of cluster where the backup should be done")
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

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			return err
		}

		if err := gitClient.Clone(); err != nil {
			return err
		}

		found, err := api.ExistsDatabaseYml(gitClient)
		if err != nil {
			return err
		}
		if found {

			value, err := ioutil.ReadFile(kubeconfig)
			if err != nil {
				monitor.Error(err)
				return nil
			}
			kubeconfigStr := string(value)

			k8sClient := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
			if k8sClient.Available() {
				if err := start.Backup(
					monitor,
					orbConfig.Path,
					k8sClient,
					backup,
					&version,
				); err != nil {
					return err
				}
			}

		}
		return nil
	}
	return cmd
}
