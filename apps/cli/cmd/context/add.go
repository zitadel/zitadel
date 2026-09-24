package context

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

func newAddCmd() *cobra.Command {
	var (
		instance string
		token    string
	)

	cmd := &cobra.Command{
		Use:   "add <name>",
		Short: "Add or update a context with a Personal Access Token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			if instance == "" {
				return fmt.Errorf("--instance is required")
			}
			if token == "" {
				return fmt.Errorf("--token is required")
			}

			cfg, err := config.Load()
			if err != nil {
				return err
			}

			ctx := config.Context{
				Instance:   instance,
				AuthMethod: "pat",
				PAT:        token,
			}

			cfg.Contexts[name] = ctx
			cfg.ActiveContext = name

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(os.Stderr, "✓ Added context %q (instance: %s)\n", name, instance)
			return nil
		},
	}

	cmd.Flags().StringVar(&instance, "instance", "", "ZITADEL instance host (e.g. mycompany.zitadel.cloud)")
	cmd.Flags().StringVar(&token, "token", "", "Personal Access Token (PAT)")

	return cmd
}
