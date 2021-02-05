package cmds

import (
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"io/ioutil"

	"github.com/caos/zitadel/operator/helpers"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/spf13/cobra"
)

func TakeoffCommand(rv RootValues) *cobra.Command {
	var (
		kubeconfig string
		cmd        = &cobra.Command{
			Use:   "takeoff",
			Short: "Launch a ZITADEL operator on the orb",
			Long:  "Ensures a desired state of the resources on the orb",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "Kubeconfig for ZITADEL operator deployment")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, monitor, orbConfig, gitClient, _, errFunc, err := rv()
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

		if err := deployOperator(
			monitor,
			gitClient,
			&kubeconfigStr,
		); err != nil {
			monitor.Error(err)
		}

		if err := deployDatabase(
			monitor,
			gitClient,
			&kubeconfigStr,
		); err != nil {
			monitor.Error(err)
		}
		return nil
	}
	return cmd
}

func deployOperator(monitor mntr.Monitor, gitClient *git.Client, kubeconfig *string) error {
	found, err := api.ExistsZitadelYml(gitClient)
	if err != nil {
		return err
	}
	if !found {
		monitor.Info("No ZITADEL operator deployed as no zitadel.yml present")
		return nil
	}

	if found {
		k8sClient := kubernetes.NewK8sClient(monitor, kubeconfig)

		if k8sClient.Available() {
			desiredTree, err := api.ReadZitadelYml(gitClient)
			if err != nil {
				return err
			}
			if err := orb.Reconcile(monitor, desiredTree, true)(k8sClient); err != nil {
				return err
			}
		}
	}
	return nil
}

func deployDatabase(monitor mntr.Monitor, gitClient *git.Client, kubeconfig *string) error {
	found, err := api.ExistsDatabaseYml(gitClient)
	if err != nil {
		return err
	}
	if found {
		k8sClient := kubernetes.NewK8sClient(monitor, kubeconfig)

		if k8sClient.Available() {
			tree, err := api.ReadDatabaseYml(gitClient)
			if err != nil {
				return err
			}

			if err := orbdb.Reconcile(
				monitor,
				tree)(k8sClient); err != nil {
				return err
			}
		} else {
			monitor.Info("Failed to connect to k8s")
		}
	}
	return nil
}
