package cmds

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/cli"
	"github.com/caos/zitadel/operator/crtlcrd"
	"github.com/caos/zitadel/operator/crtlgitops"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	kuberneteszit "github.com/caos/zitadel/pkg/kubernetes"
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

		monitor.WithFields(map[string]interface{}{
			"version": version,
		}).Info("Destroying Orb")

		if err := kuberneteszit.ScaleZitadelOperator(monitor, k8sClient, 0); err != nil {
			return err
		}

		if err := kuberneteszit.ScaleDatabaseOperator(monitor, k8sClient, 0); err != nil {
			return err
		}

		if rv.Gitops {
			if err := crtlgitops.DestroyOperator(monitor, orbConfig.Path, k8sClient, &version, rv.Gitops); err != nil {
				return err
			}

			if err := crtlgitops.DestroyDatabase(monitor, orbConfig.Path, k8sClient, &version, rv.Gitops); err != nil {
				return err
			}
		} else {
			if err := crtlcrd.Destroy(monitor, k8sClient, version, "zitadel", "database"); err != nil {
				return err
			}
		}

		if err := destroyOperator(monitor, gitClient, k8sClient, rv.Gitops); err != nil {
			return err
		}

		if err := destroyDatabase(monitor, gitClient, k8sClient, rv.Gitops); err != nil {
			return err
		}

		return nil
	}
	return cmd
}

func destroyOperator(monitor mntr.Monitor, gitClient *git.Client, k8sClient kubernetes.ClientInt, gitops bool) error {
	if gitops {
		if gitClient.Exists(git.ZitadelFile) {
			desiredTree, err := gitClient.ReadTree(git.ZitadelFile)
			if err != nil {
				return err
			}
			desired, err := orbzit.ParseDesiredV0(desiredTree)
			if err != nil {
				return err
			}
			spec := desired.Spec

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
		if gitClient.Exists(git.DatabaseFile) {
			desiredTree, err := gitClient.ReadTree(git.DatabaseFile)
			if err != nil {
				return err
			}
			desired, err := orbdb.ParseDesiredV0(desiredTree)
			if err != nil {
				return err
			}
			spec := desired.Spec

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
