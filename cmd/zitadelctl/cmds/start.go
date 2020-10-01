package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/start"
	"github.com/spf13/cobra"
	"io/ioutil"
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
		_, monitor, orbConfig, _, _, errFunc := rv()
		if errFunc != nil {
			return errFunc(cmd)
		}

		value, err := ioutil.ReadFile(kubeconfig)
		if err != nil {
			monitor.Error(err)
			return err
		}
		kubeconfigStr := string(value)

		k8sClient := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
		if k8sClient.Available() {
			return start.Operator(monitor, orbConfig.Path, k8sClient)
		}
		return nil
	}
	return cmd
}
