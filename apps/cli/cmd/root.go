package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	ctxcmd "github.com/zitadel/zitadel/apps/cli/cmd/context"
	"github.com/zitadel/zitadel/apps/cli/gen"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

var (
	version    = "dev"
	flagCtx    string
	flagOutput string
	cfg        *config.Config
)

// SetVersion sets the version string displayed by --version.
func SetVersion(v string) {
	version = v
}

// NewRootCmd creates the top-level cobra command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "zitadel",
		Short:   "ZITADEL CLI — manage your ZITADEL instances from the command line",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config loading for completion commands
			if cmd.Name() == "completion" || cmd.Name() == "help" {
				return nil
			}
			var err error
			cfg, err = config.Load()
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			// Override active context if --context flag is set
			if flagCtx != "" {
				cfg.ActiveContext = flagCtx
			}
			return nil
		},
		SilenceUsage: true,
	}

	root.PersistentFlags().StringVar(&flagCtx, "context", "", "override the active context")
	root.PersistentFlags().StringVarP(&flagOutput, "output", "o", "table", "output format: table or json")

	root.AddCommand(newLoginCmd())
	root.AddCommand(newLogoutCmd())
	root.AddCommand(ctxcmd.NewCmd())
	gen.RegisterAll(root, func() *config.Config { return cfg }, func() string { return flagOutput })

	return root
}

// activeContext resolves the current context from loaded config.
func activeContext() (*config.Context, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("config not loaded")
	}
	return config.ActiveCtx(cfg)
}
