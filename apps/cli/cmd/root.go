package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	ctxcmd "github.com/zitadel/zitadel/apps/cli/cmd/context"
	_ "github.com/zitadel/zitadel/apps/cli/gen" // keep gen import so init() runs and registers specs
	"github.com/zitadel/zitadel/apps/cli/internal/client"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"github.com/zitadel/zitadel/apps/cli/internal/output"
	"github.com/zitadel/zitadel/apps/cli/internal/runtime"
)

var (
	version         = "dev"
	flagCtx         string
	flagOutput      string
	flagFromJSON    bool
	flagDryRun      bool
	flagDebug       bool
	flagYes         bool
	flagQuiet       bool
	flagRequestJSON string
	cfg             *config.Config
)

// SetVersion sets the version string displayed by --version.
func SetVersion(v string) {
	version = v
}

// NewRootCmd creates the top-level cobra command.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "zitadel-cli",
		Short: "ZITADEL CLI — manage your ZITADEL instances from the command line",
		Long: `zitadel-cli — ZITADEL Management CLI

⚠️  EXPERIMENTAL: This CLI is under active development. Commands, flags, and
    output formats may change without notice between releases. Do not use in
    production scripts without pinning to a specific version.

    Feedback and bug reports: https://github.com/zitadel/zitadel/issues`,
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Setup logging
			logLevel := slog.LevelInfo
			if flagDebug {
				logLevel = slog.LevelDebug
				client.EnableDebug = true
			}
			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: logLevel,
			}))
			slog.SetDefault(logger)

			// Skip config loading for completion and describe commands
			if cmd.Name() == "completion" || cmd.Name() == "help" || cmd.Name() == "describe" {
				return nil
			}

			// Print experimental warning to stderr unless suppressed.
			if os.Getenv("ZITADEL_CLI_NO_WARN") == "" {
				fmt.Fprintln(os.Stderr, "⚠️  zitadel-cli is EXPERIMENTAL — commands and output may change without notice. Set ZITADEL_CLI_NO_WARN=1 to suppress.")
			}

			// Auto-detect output format: auto for TTY (table if columns, describe otherwise), JSON for pipes/redirects.
			if !cmd.Flags().Changed("output") {
				if output.IsStdoutPiped() {
					flagOutput = "json"
				} else {
					flagOutput = ""
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
	root.PersistentFlags().StringVarP(&flagOutput, "output", "o", "table", "output format: table, json, or describe (default: describe for TTY, json for pipes)")
	root.PersistentFlags().BoolVar(&flagFromJSON, "from-json", false, "read request body as JSON from stdin")
	root.PersistentFlags().StringVar(&flagRequestJSON, "request-json", "", "provide request body as inline JSON string (alternative to --from-json with stdin)")
	root.PersistentFlags().BoolVar(&flagDryRun, "dry-run", false, "print the request as JSON without calling the API")
	root.PersistentFlags().BoolVarP(&flagYes, "yes", "y", false, "skip confirmation prompts for destructive operations")
	root.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "suppress output on success (exit code only)")
	root.PersistentFlags().BoolVar(&flagDebug, "debug", false, "enable debug logging")

	root.AddCommand(newLoginCmd())
	root.AddCommand(newLogoutCmd())
	root.AddCommand(ctxcmd.NewCmd())
	root.AddCommand(newDescribeCmd())
	root.AddCommand(newSkillsCmd())
	root.AddCommand(newMCPCmd(func() *config.Config { return cfg }))
	runtime.BuildCommands(root, func() *config.Config { return cfg }, func() string { return flagOutput })

	return root
}

// activeContext resolves the current context from loaded config.
func activeContext() (*config.Context, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("config not loaded")
	}
	return config.ActiveCtx(cfg)
}
