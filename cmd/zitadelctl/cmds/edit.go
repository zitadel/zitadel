package cmds

import (
	"errors"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/cli"
	"github.com/spf13/cobra"
)

func EditCommand(getRv GetRootValues) *cobra.Command {
	return &cobra.Command{
		Use:     "edit <path>",
		Short:   "Edit the file in your favorite text editor",
		Args:    cobra.ExactArgs(1),
		Example: `zitadelctl file edit desired.yml`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			rv := getRv("edit", map[string]interface{}{"file": args[0]}, "")
			defer func() {
				err = rv.ErrFunc(err)
			}()

			orbConfig := rv.OrbConfig
			gitClient := rv.GitClient

			if !rv.Gitops {
				return mntr.ToUserError(errors.New("edit command is only supported with the --gitops flag"))
			}

			if err := cli.InitRepo(orbConfig, gitClient); err != nil {
				return err
			}

			return cli.Edit(rv.GitClient, args[0])
		},
	}
}
