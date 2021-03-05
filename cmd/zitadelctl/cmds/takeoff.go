package cmds

import (
	"io/ioutil"

	"github.com/caos/orbos/pkg/orb"

	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/api"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/spf13/cobra"
)

func TakeoffCommand(getRv GetRootValues) *cobra.Command {
	var (
		gitOpsZitadel  bool
		gitOpsDatabase bool
		cmd            = &cobra.Command{
			Use:   "takeoff",
			Short: "Launch a ZITADEL operator on the orb",
			Long:  "Ensures a desired state of the resources on the orb",
		}
	)

	flags := cmd.Flags()
	flags.BoolVar(&gitOpsZitadel, "gitops-zitadel", false, "defines if the zitadel operator should run in gitops mode")
	flags.BoolVar(&gitOpsDatabase, "gitops-database", false, "defines if the database operator should run in gitops mode")

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

		orbConfigIsIncompleteErr := orb.IsComplete(orbConfig)
		if orbConfigIsIncompleteErr != nil && (rv.Gitops || gitOpsZitadel || gitOpsDatabase) {
			return err
		}
		if orbConfigIsIncompleteErr == nil {
			if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
				monitor.Error(err)
				return nil
			}

			if err := gitClient.Clone(); err != nil {
				monitor.Error(err)
				return nil
			}
		}

		value, err := ioutil.ReadFile(rv.Kubeconfig)
		if err != nil {
			// print help
			return err
		}
		kubeconfigStr := string(value)

		if err := deployOperator(
			monitor,
			gitClient,
			&kubeconfigStr,
			rv.Version,
			rv.Gitops || gitOpsZitadel,
		); err != nil {
			monitor.Error(err)
		}

		if err := deployDatabase(
			monitor,
			gitClient,
			&kubeconfigStr,
			rv.Version,
			rv.Gitops || gitOpsDatabase,
		); err != nil {
			monitor.Error(err)
		}
		return nil
	}
	return cmd
}

func deployOperator(monitor mntr.Monitor, gitClient *git.Client, kubeconfig *string, version string, gitops bool) error {
	if gitops {
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
				desired, err := orbzit.ParseDesiredV0(desiredTree)
				if err != nil {
					return err
				}
				spec := desired.Spec
				spec.GitOps = gitops

				// at takeoff the artifacts have to be applied
				spec.SelfReconciling = true
				if err := orbzit.Reconcile(monitor, spec)(k8sClient); err != nil {
					return err
				}
			}
		}
	} else {
		k8sClient := kubernetes.NewK8sClient(monitor, kubeconfig)

		if k8sClient.Available() {
			// at takeoff the artifacts have to be applied
			spec := &orbzit.Spec{
				Version:         version,
				SelfReconciling: true,
				GitOps:          gitops,
			}

			if err := orbzit.Reconcile(monitor, spec)(k8sClient); err != nil {
				return err
			}
		} else {
			monitor.Info("Failed to connect to k8s")
		}
	}

	return nil
}

func deployDatabase(monitor mntr.Monitor, gitClient *git.Client, kubeconfig *string, version string, gitops bool) error {
	if gitops {
		found, err := api.ExistsDatabaseYml(gitClient)
		if err != nil {
			return err
		}
		if found {
			k8sClient := kubernetes.NewK8sClient(monitor, kubeconfig)

			if k8sClient.Available() {
				desiredTree, err := api.ReadDatabaseYml(gitClient)
				if err != nil {
					return err
				}
				desired, err := orbdb.ParseDesiredV0(desiredTree)
				if err != nil {
					return err
				}
				spec := desired.Spec
				spec.GitOps = gitops

				// at takeoff the artifacts have to be applied
				spec.SelfReconciling = true
				if err := orbdb.Reconcile(
					monitor,
					spec)(k8sClient); err != nil {
					return err
				}
			} else {
				monitor.Info("Failed to connect to k8s")
			}
		}
	} else {
		k8sClient := kubernetes.NewK8sClient(monitor, kubeconfig)

		if k8sClient.Available() {
			// at takeoff the artifacts have to be applied
			spec := &orbdb.Spec{
				Version:         version,
				SelfReconciling: true,
				GitOps:          gitops,
			}

			if err := orbdb.Reconcile(
				monitor,
				spec)(k8sClient); err != nil {
				return err
			}
		} else {
			monitor.Info("Failed to connect to k8s")
		}
	}
	return nil
}
