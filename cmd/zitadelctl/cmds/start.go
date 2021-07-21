package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/crtlcrd"
	"github.com/caos/zitadel/operator/crtlgitops"
	"github.com/spf13/cobra"
)

func StartOperator(getRv GetRootValues) *cobra.Command {
	var (
		metricsAddr string
		cmd         = &cobra.Command{
			Use:   "operator",
			Short: "Launch a ZITADEL operator",
			Long:  "Ensures a desired state of ZITADEL",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&metricsAddr, "metrics-addr", "", "The address the metric endpoint binds to.")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		rv := getRv("operator", nil, "zitadel-operator", "zitadel")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		version := rv.Version

		if rv.Gitops {
			k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, rv.Kubeconfig)
			if err != nil {
				return err
			}

			return crtlgitops.Operator(monitor, orbConfig.Path, k8sClient, &version, rv.Gitops)
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
		metricsAddr string
		cmd         = &cobra.Command{
			Use:   "database",
			Short: "Launch a database operator",
			Long:  "Ensures a desired state of the database",
		}
	)
	flags := cmd.Flags()
	flags.StringVar(&metricsAddr, "metrics-addr", "", "The address the metric endpoint binds to.")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		rv := getRv("database", nil, "database-operator")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		version := rv.Version

		if rv.Gitops {
			k8sClient, err := kubernetes.NewK8sClientWithPath(monitor, rv.Kubeconfig)
			if err != nil {
				return err
			}
			return crtlgitops.Database(monitor, orbConfig.Path, k8sClient, &version, rv.Gitops)
		} else {
			if err := crtlcrd.Start(monitor, version, metricsAddr, crtlcrd.Database); err != nil {
				return err
			}
		}

		return nil
	}
	return cmd
}
