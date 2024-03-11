package alpha

import (
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/cmd/alpha/client"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "alpha state cli commands",
	}
	cmd.AddCommand(
		client.New(),
	)

	return cmd
}
