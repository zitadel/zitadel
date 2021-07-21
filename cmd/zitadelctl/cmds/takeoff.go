package cmds

import (
	"github.com/caos/orbos/pkg/kubernetes/cli"
	"gopkg.in/yaml.v3"

	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TakeoffCommand(getRv GetRootValues) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "takeoff",
			Short: "Launch a ZITADEL operator and database operator on the orb",
			Long:  "Ensures a desired state of the resources on the orb",
		}
	)

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		rv := getRv("takeoff", nil, "")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		monitor := rv.Monitor
		orbConfig := rv.OrbConfig
		gitClient := rv.GitClient

		k8sClient, err := cli.Client(
			monitor,
			orbConfig,
			gitClient,
			rv.Kubeconfig,
			rv.Gitops,
			true,
		)
		if err != nil {
			return err
		}

		if err := kubernetes.EnsureCaosSystemNamespace(monitor, k8sClient); err != nil {
			return err
		}

		if rv.Gitops {

			orbConfigBytes, err := yaml.Marshal(orbConfig)
			if err != nil {
				return err
			}

			if err := kubernetes.EnsureOrbconfigSecret(monitor, k8sClient, orbConfigBytes); err != nil {
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
			return err
		}

		return deployDatabase(
			monitor,
			gitClient,
			k8sClient,
			rv.Version,
			rv.Gitops,
		)
	}
	return cmd
}

func deployOperator(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, version string, gitops bool) error {
	if !gitops {

		// at takeoff the artifacts have to be applied
		spec := &orbzit.Spec{
			Version:         version,
			SelfReconciling: true,
		}

		return orbzit.Reconcile(monitor, spec, gitops)(k8sClient)
	}

	if !gitClient.Exists(git.ZitadelFile) {
		monitor.WithField("file", git.ZitadelFile).Info("File not found in git, skipping deployment")
		return nil
	}

	desiredTree, err := gitClient.ReadTree(git.ZitadelFile)
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
	return orbzit.Reconcile(monitor, spec, gitops)(k8sClient)
}

func deployDatabase(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, version string, gitops bool) error {
	if !gitops {

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

	if !gitClient.Exists(git.DatabaseFile) {
		monitor.WithField("file", git.DatabaseFile).Info("File not found in git, skipping deployment")
		return nil
	}
	desiredTree, err := gitClient.ReadTree(git.DatabaseFile)
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
	return orbdb.Reconcile(
		monitor,
		spec,
		gitops,
	)(k8sClient)
}
