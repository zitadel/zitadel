package context

import (
	"github.com/spf13/cobra"
)

// NewCmd creates the "context" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "context",
		Aliases: []string{"ctx"},
		Short:   "Manage CLI contexts (configured ZITADEL instances)",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newUseCmd())
	cmd.AddCommand(newCurrentCmd())

	return cmd
}
