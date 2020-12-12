package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/start"
	"github.com/spf13/cobra"
)

func StartOperator(rv RootValues) *cobra.Command {
	var (
		kubeconfig     string
		migrationsPath string
		cmd            = &cobra.Command{
			Use:   "operator",
			Short: "Launch a ZITADEL operator",
			Long:  "Ensures a desired state of ZITADEL",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Kubeconfig for ZITADEL operator deployment")
	flags.StringVar(&migrationsPath, "migrations", "./migrations/", "Path to the migration files")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		_, monitor, orbConfig, _, version, errFunc := rv()
		if errFunc != nil {
			return errFunc(cmd)
		}

		k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, kubeconfig)
		if err != nil {
			return err
		}

		if k8sClient.Available() {
			return start.Operator(monitor, orbConfig.Path, k8sClient, migrationsPath, &version)
		}
		return nil
	}
	return cmd
}
