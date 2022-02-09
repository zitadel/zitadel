package initialise

import (
	_ "embed"

	"github.com/caos/logging"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialize ZITADEL instance",
		Long: `init sets up the minimum requirements to start ZITADEL.
Prereqesits:
- cockroachdb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logging.New().Info("hello world")
			return nil
		},
	}
}
