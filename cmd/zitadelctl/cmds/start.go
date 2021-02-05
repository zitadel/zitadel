package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/controller"
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/start"
	"github.com/spf13/cobra"
)

func StartOperator(getRv GetRootValues) *cobra.Command {
	var (
		kubeconfig string
		gitOpsMode bool
		cmd        = &cobra.Command{
			Use:   "operator",
			Short: "Launch a ZITADEL operator",
			Long:  "Ensures a desired state of ZITADEL",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "Kubeconfig for ZITADEL operator deployment")
	flags.BoolVar(&gitOpsMode, "gitops", false, "defines if the ZITADEL operator should run in gitops mode")

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
		version := rv.Version
		metricsAddr := rv.MetricsAddr

		kubeconfig = helpers.PruneHome(kubeconfig)

		if gitOpsMode {
			k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, kubeconfig)
			if err != nil {
				monitor.Error(err)
				return nil
			}

			if k8sClient.Available() {
				if err := start.Operator(monitor, orbConfig.Path, k8sClient, &version, gitOpsMode); err != nil {
					monitor.Error(err)
					return nil
				}
			}
		} else {
			if err := controller.Start(monitor, version, metricsAddr, controller.Zitadel); err != nil {
				return err
			}
		}

		return nil
	}
	return cmd
}

func StartDatabase(getRv GetRootValues) *cobra.Command {
	var (
		kubeconfig string
		gitOpsMode bool
		cmd        = &cobra.Command{
			Use:   "database",
			Short: "Launch a database operator",
			Long:  "Ensures a desired state of the database",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig used by database operator")
	flags.BoolVar(&gitOpsMode, "gitops", false, "defines if the database operator should run in gitops mode")

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
		version := rv.Version
		metricsAddr := rv.MetricsAddr

		if gitOpsMode {
			k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, kubeconfig)
			if err != nil {
				return err
			}

			if k8sClient.Available() {
				return start.Database(monitor, orbConfig.Path, k8sClient, &version, gitOpsMode)
			}
		} else {
			if err := controller.Start(monitor, version, metricsAddr, controller.Database); err != nil {
				return err
			}
		}

		return nil
	}
	return cmd
}
