package cmds

import (
	"errors"
	"github.com/caos/orbos/pkg/cli"
	"github.com/spf13/cobra"
	"strings"
)

func RemoveCommand(getRv GetRootValues) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "remove <filepath>",
		Short:   "Remove file from git repository",
		Long:    "If the file doesn't exist, the command completes successfully",
		Args:    cobra.MinimumNArgs(1),
		Example: `zitadelctl file remove caos-internal/orbiter/current.yml`,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {

		filesStr := strings.Join(args, ",")

		rv := getRv("remove", map[string]interface{}{"files": filesStr}, "")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		if !rv.Gitops {
			return errors.New("remove command is only supported with the --gitops flag")
		}

		if err := cli.InitRepo(rv.OrbConfig, rv.GitClient); err != nil {
			return err
		}

		return cli.Remove(rv.GitClient, args)
	}

	return cmd
}
