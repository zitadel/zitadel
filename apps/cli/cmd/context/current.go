package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

func newCurrentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the active context",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			ctx, name, err := config.ActiveCtx(cfg)
			if err != nil {
				return err
			}

			if name == "" {
				fmt.Println("Using environment variables (no named context)")
				fmt.Printf("  Instance: %s\n", ctx.Instance)
			} else {
				fmt.Printf("Context:  %s\n", name)
				fmt.Printf("Instance: %s\n", ctx.Instance)
				fmt.Printf("Auth:     %s\n", ctx.AuthMethod)
			}
			return nil
		},
	}
}
