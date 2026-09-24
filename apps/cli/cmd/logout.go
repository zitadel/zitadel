package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear the stored token for the active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, name, err := activeContext()
			if err != nil {
				return err
			}

			ctx.Token = ""
			cfg.Contexts[name] = *ctx

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(os.Stderr, "✓ Logged out from %s (context: %s)\n", ctx.Instance, name)
			return nil
		},
	}
}
