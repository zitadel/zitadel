package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes/cli"
	"gopkg.in/yaml.v3"

	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/api"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TakeoffCommand(getRv GetRootValues) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "takeoff",
			Short: "Launch a ZITADEL operator on the orb",
			Long:  "Ensures a desired state of the resources on the orb",
		}
	)

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
		gitClient := rv.GitClient

		k8sClient, _, err := cli.Client(
			monitor,
			orbConfig,
			gitClient,
			rv.Kubeconfig,
			rv.Gitops,
		)
		if err != nil {
			return err
		}

		if err := kubernetes.EnsureCaosSystemNamespace(monitor, k8sClient); err != nil {
			monitor.Info("failed to apply common resources into k8s-cluster")
			return err
		}

		if rv.Gitops {

			orbConfigBytes, err := yaml.Marshal(orbConfig)
			if err != nil {
				return err
			}

			if err := kubernetes.EnsureOrbconfigSecret(monitor, k8sClient, orbConfigBytes); err != nil {
				monitor.Info("failed to apply configuration resources into k8s-cluster")
				return err
			}
		}

		if err := deployOperator(
			monitor,
			gitClient,
			k8sClient,
			rv.Version,
			rv.Gitops,
		); err != nil {
			monitor.Error(err)
		}

		if err := deployDatabase(
			monitor,
			gitClient,
			k8sClient,
			rv.Version,
			rv.Gitops,
		); err != nil {
			monitor.Error(err)
		}
		return nil
	}
	return cmd
}

func deployOperator(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, version string, gitops bool) error {
	if gitops {
		found, err := api.ExistsZitadelYml(gitClient)
		if err != nil {
			return err
		}
		if found {

			desiredTree, err := api.ReadZitadelYml(gitClient)
			if err != nil {
				return err
			}
			desired, err := orbzit.ParseDesiredV0(desiredTree)
			if err != nil {
				return err
			}
			spec := desired.Spec

			// at takeoff the artifacts have to be applied
			spec.SelfReconciling = true
			rec, _ := orbzit.Reconcile(monitor, spec, gitops)
			if err := rec(k8sClient); err != nil {
				return err
			}
		}
	} else {
		// at takeoff the artifacts have to be applied
		spec := &orbzit.Spec{
			Version:         version,
			SelfReconciling: true,
		}

		rec, _ := orbzit.Reconcile(monitor, spec, gitops)
		if err := rec(k8sClient); err != nil {
			return err
		}
	}

	return nil
}

func deployDatabase(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, version string, gitops bool) error {
	if gitops {
		found, err := api.ExistsDatabaseYml(gitClient)
		if err != nil {
			return err
		}
		if found {
			desiredTree, err := api.ReadDatabaseYml(gitClient)
			if err != nil {
				return err
			}
			desired, err := orbdb.ParseDesiredV0(desiredTree)
			if err != nil {
				return err
			}
			spec := desired.Spec

			// at takeoff the artifacts have to be applied
			spec.SelfReconciling = true
			if err := orbdb.Reconcile(
				monitor,
				spec,
				gitops,
			)(k8sClient); err != nil {
				return err
			}
		}
	} else {
		// at takeoff the artifacts have to be applied
		spec := &orbdb.Spec{
			Version:         version,
			SelfReconciling: true,
		}

		if err := orbdb.Reconcile(
			monitor,
			spec,
			gitops,
		)(k8sClient); err != nil {
			return err
		}
	}
	return nil
}
