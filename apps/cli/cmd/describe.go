package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/gen"
)

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe [group] [command]",
		Short: "Describe CLI commands as machine-readable JSON schema",
		Long: `Dump the schema of CLI commands as JSON for agent introspection.

Examples:
  zitadel describe                     # list all command groups and commands
  zitadel describe users               # list all commands in the users group
  zitadel describe users list-passkeys # describe a specific command with flags`,
		RunE: func(cmd *cobra.Command, args []string) error {
			allMeta := gen.AllMeta()

			switch len(args) {
			case 0:
				return describeAll(allMeta)
			case 1:
				return describeGroup(allMeta, args[0])
			default:
				return describeCommand(allMeta, args[0], args[1])
			}
		},
		SilenceUsage: true,
	}
}

func describeAll(meta []gen.CommandMeta) error {
	groups := make(map[string][]string)
	for _, m := range meta {
		name := m.Name
		if idx := strings.IndexByte(name, ' '); idx >= 0 {
			name = name[:idx]
		}
		groups[m.Group] = appendUnique(groups[m.Group], name)
	}
	return writeJSON(groups)
}

func describeGroup(meta []gen.CommandMeta, group string) error {
	var commands []gen.CommandMeta
	for _, m := range meta {
		if m.Group == group {
			commands = append(commands, m)
		}
	}
	if len(commands) == 0 {
		return fmt.Errorf("unknown command group %q", group)
	}
	return writeJSON(commands)
}

func describeCommand(meta []gen.CommandMeta, group, command string) error {
	for _, m := range meta {
		if m.Group != group {
			continue
		}
		name := m.Name
		if idx := strings.IndexByte(name, ' '); idx >= 0 {
			name = name[:idx]
		}
		if name == command {
			return writeJSON(m)
		}
	}
	return fmt.Errorf("unknown command %q in group %q", command, group)
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}
