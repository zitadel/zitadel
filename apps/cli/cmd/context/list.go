package context

import (
	"sort"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"github.com/zitadel/zitadel/apps/cli/internal/output"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			if len(cfg.Contexts) == 0 {
				cmd.Println("No contexts configured. Use 'zitadel login' to add one.")
				return nil
			}

			// Sort names for stable output
			names := make([]string, 0, len(cfg.Contexts))
			for name := range cfg.Contexts {
				names = append(names, name)
			}
			sort.Strings(names)

			header := []string{"", "NAME", "INSTANCE", "AUTH METHOD"}
			var rows [][]string
			for _, name := range names {
				ctx := cfg.Contexts[name]
				active := ""
				if name == cfg.ActiveContext {
					active = "*"
				}
				rows = append(rows, []string{active, name, ctx.Instance, ctx.AuthMethod})
			}

			output.Table(header, rows)
			return nil
		},
	}
}
