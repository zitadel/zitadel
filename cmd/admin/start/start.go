package start

import (
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "starts ZITADEL instance",
		Long: `starts ZITADEL.
Requirements:
- cockroachdb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.Info("hello world")
			logging.WithFields("field", 1).Info("hello world")
			return nil
		},
	}
}
