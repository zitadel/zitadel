package cmds

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/kinds/orb"
	"github.com/spf13/cobra"
)

func TakeoffCommand(rv RootValues) *cobra.Command {
	var (
		kubeconfig string
		cmd        = &cobra.Command{
			Use:   "takeoff",
			Short: "Launch an ZITADEL operator",
			Long:  "Ensures a desired state",
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Kubeconfig for boom deployment")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, monitor, orbConfig, gitClient, version, errFunc := rv()
		if errFunc != nil {
			return errFunc(cmd)
		}

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			return err
		}

		k8sClient := kubernetes.NewK8sClient(monitor, &kubeconfig)
		if k8sClient.Available() {
			return deployOperator(
				monitor,
				gitClient,
				&kubeconfig,
				version,
			)
		}
		return nil
	}
	return cmd
}

func deployOperator(monitor mntr.Monitor, gitClient *git.Client, kubeconfig *string, version string) error {
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
			if err := orb.Reconcile(monitor, desiredTree, version)(k8sClient); err != nil {
				return err
			}
		}
	}
	return nil
}
