package setup

import (
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "setup ZITADEL instance",
		Long: `sets up data to start ZITADEL.
Requirements:
- cockroachdb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.Log("SETUP-e88M6").Info("hello world")
			return nil
		},
	}
}
