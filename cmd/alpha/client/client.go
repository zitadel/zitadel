package client

import (
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "client calls zitadel",
	}

	cmd.AddCommand(
		admin.AdminServiceCmd,
	)

	return cmd
}
