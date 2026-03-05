package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/auth"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

func newLoginCmd() *cobra.Command {
	var (
		instance string
		clientID string
		ctxName  string
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a ZITADEL instance using browser-based OIDC PKCE login",
		RunE: func(cmd *cobra.Command, args []string) error {
			if token := os.Getenv("ZITADEL_TOKEN"); token != "" {
				fmt.Fprintln(os.Stderr, "ZITADEL_TOKEN is set — PAT mode is active. Browser login skipped.")
				return nil
			}

			if instance == "" {
				return fmt.Errorf("--instance is required")
			}
			if clientID == "" {
				return fmt.Errorf("--client-id is required (create a Native application in your ZITADEL project)")
			}

			if ctxName == "" {
				ctxName = instance
			}

			token, err := auth.Login(cmd.Context(), instance, clientID)
			if err != nil {
				return fmt.Errorf("login failed: %w", err)
			}

			ctx := config.Context{
				Instance:   instance,
				AuthMethod: "interactive",
				ClientID:   clientID,
				Token:      token.AccessToken,
			}
			cfg.Contexts[ctxName] = ctx
			cfg.ActiveContext = ctxName

			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(os.Stderr, "✓ Logged in to %s (context: %s)\n", instance, ctxName)
			return nil
		},
	}

	cmd.Flags().StringVar(&instance, "instance", "", "ZITADEL instance host (e.g. mycompany.zitadel.cloud)")
	cmd.Flags().StringVar(&clientID, "client-id", "", "OIDC client ID of a native application")
	cmd.Flags().StringVar(&ctxName, "context", "", "name for this context (defaults to instance host)")

	return cmd
}
