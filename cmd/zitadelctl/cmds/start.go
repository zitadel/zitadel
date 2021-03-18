package cmds

import (
	"flag"

	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/crtlcrd"
	"github.com/caos/zitadel/operator/crtlgitops"
	"github.com/spf13/cobra"
)

func StartOperator(getRv GetRootValues) *cobra.Command {
	var (
		metricsAddr string
		gitOpsMode  bool
		cmd         = &cobra.Command{
			Use:   "operator",
			Short: "Launch a ZITADEL operator",
			Long:  "Ensures a desired state of ZITADEL",
		}
	)
	flags := cmd.Flags()
	flags.BoolVar(&gitOpsMode, "gitops", false, "defines if the ZITADEL operator should run in gitops mode")
	flag.StringVar(&metricsAddr, "metrics-addr", "", "The address the metric endpoint binds to.")

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

		if gitOpsMode {
			k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, rv.Kubeconfig)
			if err != nil {
				monitor.Error(err)
				return nil
			}

			if k8sClient.Available() {
				if err := crtlgitops.Operator(monitor, orbConfig.Path, k8sClient, &version, gitOpsMode); err != nil {
					monitor.Error(err)
					return nil
				}
			}
		} else {
			if err := crtlcrd.Start(monitor, version, metricsAddr, crtlcrd.Zitadel); err != nil {
				return err
			}
		}

		return nil
	}
	return cmd
}

func StartDatabase(getRv GetRootValues) *cobra.Command {
	var (
		kubeconfig  string
		gitOpsMode  bool
		metricsAddr string
		cmd         = &cobra.Command{
			Use:   "database",
			Short: "Launch a database operator",
			Long:  "Ensures a desired state of the database",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig used by database operator")
	flags.BoolVar(&gitOpsMode, "gitops", false, "defines if the database operator should run in gitops mode")
	flags.StringVar(&metricsAddr, "metrics-addr", "", "The address the metric endpoint binds to.")

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

		if gitOpsMode {
			k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, kubeconfig)
			if err != nil {
				return err
			}

			if k8sClient.Available() {
				return crtlgitops.Database(monitor, orbConfig.Path, k8sClient, &version, gitOpsMode)
			}
		} else {
			if err := crtlcrd.Start(monitor, version, metricsAddr, crtlcrd.Database); err != nil {
				return err
			}
		}

		return nil
	}
	return cmd
}
