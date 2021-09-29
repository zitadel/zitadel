package cmds

import (
	"errors"
	"fmt"
	"github.com/caos/orbos/pkg/cli"
	"github.com/spf13/cobra"
)

func PrintCommand(getRv GetRootValues) *cobra.Command {

	return &cobra.Command{
		Use:     "print <path>",
		Short:   "Print the files contents to stdout",
		Args:    cobra.ExactArgs(1),
		Example: `zitadelctl file print orbiter.yml`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			file := args[0]

			rv := getRv("print", map[string]interface{}{"file": file}, "")
			defer func() {
				err = rv.ErrFunc(err)
			}()

			if !rv.Gitops {
				return errors.New("print command is only supported with the --gitops flag")
			}

			if err := cli.InitRepo(rv.OrbConfig, rv.GitClient); err != nil {
				return err
			}

			fmt.Print(string(rv.GitClient.Read(file)))
			return nil
		},
	}
}
