package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	ctxcmd "github.com/zitadel/zitadel/apps/cli/cmd/context"
	"github.com/zitadel/zitadel/apps/cli/gen"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

var (
	version        = "dev"
	flagCtx        string
	flagOutput     string
	flagFromJSON   bool
	flagDryRun     bool
	flagRequestJSON string
	cfg            *config.Config
)

// SetVersion sets the version string displayed by --version.
func SetVersion(v string) {
	version = v
}

// NewRootCmd creates the top-level cobra command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "zitadel-cli",
		Short:   "ZITADEL CLI — manage your ZITADEL instances from the command line",
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip config loading for completion and describe commands
			if cmd.Name() == "completion" || cmd.Name() == "help" || cmd.Name() == "describe" {
				return nil
			}

			// Auto-detect output format: JSON when stdout is not a TTY
			if !cmd.Flags().Changed("output") {
				if fi, err := os.Stdout.Stat(); err == nil && fi.Mode()&os.ModeCharDevice == 0 {
					flagOutput = "json"
				}
			}

			// Skip config loading for dry-run if no config exists
			var err error
			cfg, err = config.Load()
			if err != nil {
				if flagDryRun {
					cfg = &config.Config{Contexts: make(map[string]config.Context)}
				} else {
					return fmt.Errorf("loading config: %w", err)
				}
			}
			// Override active context if --context flag is set
			if flagCtx != "" {
				cfg.ActiveContext = flagCtx
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&flagCtx, "context", "", "override the active context")
	root.PersistentFlags().StringVarP(&flagOutput, "output", "o", "table", "output format: table or json (auto-detected from TTY)")
	root.PersistentFlags().BoolVar(&flagFromJSON, "from-json", false, "read request body as JSON from stdin")
	root.PersistentFlags().StringVar(&flagRequestJSON, "request-json", "", "provide request body as inline JSON string (alternative to --from-json with stdin)")
	root.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "print the request as JSON without calling the API")

	root.AddCommand(newLoginCmd())
	root.AddCommand(newLogoutCmd())
	root.AddCommand(ctxcmd.NewCmd())
	root.AddCommand(newDescribeCmd())
	gen.RegisterAll(root, func() *config.Config { return cfg }, func() string { return flagOutput })

	return root
}

// FlagFromJSON returns whether the --from-json flag was set.
func FlagFromJSON() bool { return flagFromJSON }

// FlagDryRun returns whether the --dry-run flag was set.
func FlagDryRun() bool { return flagDryRun }

// activeContext resolves the current context from loaded config.
func activeContext() (*config.Context, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("config not loaded")
	}
	return config.ActiveCtx(cfg)
}
