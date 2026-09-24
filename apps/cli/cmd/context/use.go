package context

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

func newUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use [name]",
		Short: "Set the active context. If no name is provided, an interactive menu is shown.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if len(cfg.Contexts) == 0 {
				return fmt.Errorf("no contexts configured. Run 'zitadel-cli login' to create one.")
			}

			var name string
			if len(args) == 1 {
				name = args[0]
			} else {
				// Interactive picker
				var options []string
				for k := range cfg.Contexts {
					options = append(options, k)
				}
				
				err := survey.AskOne(&survey.Select{
					Message: "Choose a context to use:",
					Options: options,
					Default: cfg.ActiveContext,
				}, &name)
				if err != nil {
					return err
				}
			}

			if _, ok := cfg.Contexts[name]; !ok {
				return fmt.Errorf("context %q not found", name)
			}

			cfg.ActiveContext = name
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Switched to context %q\n", name)
			return nil
		},
	}
}
