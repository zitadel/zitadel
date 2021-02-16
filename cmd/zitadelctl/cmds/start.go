package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/start"
	"github.com/spf13/cobra"
)

func StartOperator(rv RootValues) *cobra.Command {
	var (
		kubeconfig string
		cmd        = &cobra.Command{
			Use:   "operator",
			Short: "Launch a ZITADEL operator",
			Long:  "Ensures a desired state of ZITADEL",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Kubeconfig for ZITADEL operator deployment")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, monitor, orbConfig, _, version, errFunc, err := rv()
		if err != nil {
			return err
		}
		defer func() {
			err = errFunc(err)
		}()

		kubeconfig = helpers.PruneHome(kubeconfig)

		k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, kubeconfig)
		if err != nil {
			monitor.Error(err)
			return nil
		}

		if k8sClient.Available() {
			if err := start.Operator(monitor, orbConfig.Path, k8sClient, &version); err != nil {
				monitor.Error(err)
				return nil
			}
		}
		return nil
	}
	return cmd
}

func StartDatabase(rv RootValues) *cobra.Command {
	var (
		kubeconfig string
		cmd        = &cobra.Command{
			Use:   "database",
			Short: "Launch a database operator",
			Long:  "Ensures a desired state of the database",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig used by zitadel operator")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		_, monitor, orbConfig, _, version, errFunc, err := rv()
		if err != nil {
			return err
		}
		defer func() {
			err = errFunc(err)
		}()

		k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, kubeconfig)
		if err != nil {
			return err
		}

		if k8sClient.Available() {
			return start.Database(monitor, orbConfig.Path, k8sClient, &version)
		}
		return nil
	}
	return cmd
}
