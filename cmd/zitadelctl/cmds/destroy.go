package cmds

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/cli"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/crtlcrd"
	"github.com/caos/zitadel/operator/crtlgitops"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/spf13/cobra"
)

func TeardownCommand(getRv GetRootValues) *cobra.Command {

	var (
		cmd = &cobra.Command{
			Use:   "teardown",
			Short: "Tear down an Orb",
			Long:  "Destroys a whole Orb",
			Aliases: []string{
				"shoot",
				"destroy",
				"devastate",
				"annihilate",
				"crush",
				"bulldoze",
				"total",
				"smash",
				"decimate",
				"kill",
				"trash",
				"wipe-off-the-map",
				"pulverize",
				"take-apart",
				"destruct",
				"obliterate",
				"disassemble",
				"explode",
				"blow-up",
			},
		}
	)

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
		gitClient := rv.GitClient
		version := rv.Version

		if err := gitClient.Configure(orbConfig.URL, []byte(orbConfig.Repokey)); err != nil {
			return err
		}

		if err := gitClient.Clone(); err != nil {
			return err
		}

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

		monitor.WithFields(map[string]interface{}{
			"version": version,
			"repoURL": orbConfig.URL,
		}).Info("Destroying Orb")

		if err := destroyOperator(monitor, gitClient, k8sClient, rv.Gitops); err != nil {
			return err
		}

		if err := destroyDatabase(monitor, gitClient, k8sClient, rv.Gitops); err != nil {
			return err
		}

		if rv.Gitops {
			k8sClient, _, err := cli.Client(monitor, orbConfig, rv.GitClient, rv.Kubeconfig, rv.Gitops)
			if err != nil {
				return err
			}

			if err := crtlgitops.DestroyOperator(monitor, orbConfig.Path, k8sClient, &version, rv.Gitops); err != nil {
				return err
			}

			if err := crtlgitops.DestroyDatabase(monitor, orbConfig.Path, k8sClient, &version, rv.Gitops); err != nil {
				return err
			}
		} else {
			if err := crtlcrd.Destroy(monitor, k8sClient, "zitadel", "database"); err != nil {
				return err
			}
		}

		return nil
	}
	return cmd
}

func destroyOperator(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, gitops bool) error {
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
			_, del := orbzit.Reconcile(monitor, spec, gitops)
			if err := del(k8sClient); err != nil {
				return err
			}
		}
	} else {
		_, del := orbzit.Reconcile(monitor, &orbzit.Spec{}, gitops)
		if err := del(k8sClient); err != nil {
			return err
		}
	}
	return nil
}

func destroyDatabase(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, gitops bool) error {
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
			_, del := orbdb.Reconcile(monitor, spec, gitops)
			if err := del(k8sClient); err != nil {
				return err
			}
		}
	} else {
		_, del := orbdb.Reconcile(monitor, &orbdb.Spec{}, gitops)
		if err := del(k8sClient); err != nil {
			return err
		}
	}
	return nil
}
