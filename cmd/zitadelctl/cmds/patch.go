package cmds

import (
	"errors"
	"github.com/caos/orbos/pkg/cli"
	"strings"

	"github.com/spf13/cobra"
)

func PatchCommand(getRv GetRootValues) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch <filepath> [yamlpath]",
		Short: "Patch a yaml property",
		Args:  cobra.MinimumNArgs(1),
		Example: `Overwiting a file: zitadelctl file patch zitadel.yml --exact
Patching an edge property interactively: zitadelctl file patch zitadel.yml
Patching a node property non-interactively: zitadelctl file path zitadel.yml iam --exact --file /path/to/my/iam/definition.yml`,
	}
	flags := cmd.Flags()
	var (
		value        string
		file         string
		stdin, exact bool
	)
	flags.StringVar(&value, "value", "", "Content value")
	flags.StringVar(&file, "file", "", "File containing the content value")
	flags.BoolVar(&stdin, "stdin", false, "Read content value by stdin")
	flags.BoolVar(&exact, "exact", false, "Write the content exactly at the path given without further prompting")

	cmd.RunE = func(cmd *cobra.Command, args []string) (err error) {

		var path []string
		if len(args) > 1 {
			path = strings.Split(args[1], ".")
		}

		filePath := args[0]

		rv := getRv("patch", map[string]interface{}{"value": value, "filePath": filePath, "valuePath": file, "stdin": stdin, "exact": exact}, "")
		defer func() {
			err = rv.ErrFunc(err)
		}()

		if !rv.Gitops {
			return errors.New("patch command is only supported with the --gitops flag")
		}

		contentStr, err := cli.Content(value, file, stdin)
		if err != nil {
			return err
		}

		return cli.PatchFile(rv.OrbConfig, rv.GitClient, path, contentStr, exact, filePath)
	}

	return cmd
}
