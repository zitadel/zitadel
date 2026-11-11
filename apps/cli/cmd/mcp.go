package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	_ "github.com/zitadel/zitadel/apps/cli/gen"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
	"github.com/zitadel/zitadel/apps/cli/internal/runtime"
)

func newMCPCmd(getCfg func() *config.Config) *cobra.Command {
	var services string

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start an MCP (Model Context Protocol) tool server over stdio",
		Long: `Run ZITADEL CLI as an MCP server, exposing API commands as tools
for AI agents via JSON-RPC 2.0 over stdin/stdout.

Use --services to limit which command groups are exposed as tools.

Example agent configuration (Claude Desktop):
  {
    "mcpServers": {
      "zitadel": {
        "command": "zitadel-cli",
        "args": ["mcp", "--services", "users,orgs"]
      }
    }
  }`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var filterGroups []string
			if services != "" {
				for _, s := range strings.Split(services, ",") {
					s = strings.TrimSpace(s)
					if s != "" {
						filterGroups = append(filterGroups, s)
					}
				}
			}
			server := runtime.NewMCPServer(getCfg, filterGroups)
			return server.Run()
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&services, "services", "", "Comma-separated list of service groups to expose (e.g., users,orgs)")

	return cmd
}
