package cmds

import (
	"github.com/spf13/cobra"
)

func FileCommand(getRootValues GetRootValues) *cobra.Command {

	file := &cobra.Command{
		Use:     "file <path> [command]",
		Short:   "Work with an orbs remote repository file",
		Example: `orbctl file <edit|print|patch|remove> orbiter.yml `,
	}

	file.AddCommand(
		EditCommand(getRootValues),
		PrintCommand(getRootValues),
		PatchCommand(getRootValues),
		RemoveCommand(getRootValues),
	)
	return file
}
