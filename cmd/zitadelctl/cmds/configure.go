package cmds

import (
	"errors"

	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/orbos/pkg/cfg"
	"github.com/caos/orbos/pkg/git"

	"github.com/caos/orbos/pkg/kubernetes/cli"
	"github.com/caos/orbos/pkg/orb"
	"github.com/spf13/cobra"

	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	orbzit "github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

func ConfigCommand(getRv GetRootValues, ghClientID, ghClientSecret string) *cobra.Command {

	var (
		newMasterKey string
		newRepoURL   string
		cmd          = &cobra.Command{
			Use:     "configure",
			Short:   "Configures and reconfigures an orb",
			Long:    "Generates missing secrets where it makes sense",
			Aliases: []string{"reconfigure", "config", "reconfig"},
		}
	)

	flags := cmd.Flags()
	flags.StringVar(&newMasterKey, "masterkey", "", "Reencrypts all secrets")
	flags.StringVar(&newRepoURL, "repourl", "", "Configures the repository URL")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {

		rv, _ := getRv()
		defer func() {
			err = rv.ErrFunc(err)
		}()

		if !rv.Gitops {
			return errors.New("configure command is only supported with the --gitops flag")
		}

		if err := orb.Reconfigure(rv.Ctx, rv.Monitor, rv.OrbConfig, newRepoURL, newMasterKey, rv.GitClient, ghClientID, ghClientSecret); err != nil {
			return err
		}

		k8sClient, err := cli.Client(rv.Monitor, rv.OrbConfig, rv.GitClient, rv.Kubeconfig, rv.Gitops)
		if err != nil {
			// ignore
			err = nil
		}

		if err := cfg.ApplyOrbconfigSecret(
			rv.OrbConfig,
			k8sClient,
			rv.Monitor,
		); err != nil {
			return err
		}

		return cfg.ConfigureOperators(
			rv.GitClient,
			rv.OrbConfig.Masterkey,
			append(cfg.ORBOSConfigurers(
				rv.Monitor,
				rv.OrbConfig,
				rv.GitClient,
			), cfg.OperatorConfigurer(
				git.DatabaseFile,
				rv.Monitor,
				rv.GitClient,
				func() (*tree.Tree, interface{}, error) {
					tree, err := rv.GitClient.ReadTree(git.DatabaseFile)
					if err != nil {
						return nil, nil, err
					}

					parsed, err := orbdb.ParseDesiredV0(tree)
					return tree, parsed, err
				},
			), cfg.OperatorConfigurer(
				git.ZitadelFile,
				rv.Monitor,
				rv.GitClient,
				func() (*tree.Tree, interface{}, error) {
					tree, err := rv.GitClient.ReadTree(git.ZitadelFile)
					if err != nil {
						return nil, nil, err
					}

					parsed, err := orbzit.ParseDesiredV0(tree)
					return tree, parsed, err
				},
			)))
	}
	return cmd
}
