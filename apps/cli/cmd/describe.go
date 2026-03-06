package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/gen"
)

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe [group] [command ...]",
		Short: "Describe CLI commands as machine-readable JSON schema",
		Long: `Dump the schema of CLI commands as JSON for agent introspection.

Examples:
  zitadel describe                     # list all command groups and commands
  zitadel describe users               # list all commands in the users group
  zitadel describe users create human  # describe a specific subcommand with flags`,
		RunE: func(cmd *cobra.Command, args []string) error {
			allMeta := gen.AllMeta()

			switch len(args) {
			case 0:
				return describeAll(allMeta)
			case 1:
				return describeGroup(allMeta, args[0])
			default:
				return describeCommand(allMeta, args[0], strings.Join(args[1:], " "))
			}
		},
		SilenceUsage: true,
	}
}

// describeAllOutput is the top-level describe response.
type describeAllOutput struct {
	GlobalFlags []gen.FlagMeta      `json:"global_flags"`
	Groups      map[string][]string `json:"groups"`
}

func describeAll(meta []gen.CommandMeta) error {
	groups := make(map[string][]string)
	for _, m := range meta {
		name := commandPath(m.Name)
		groups[m.Group] = appendUnique(groups[m.Group], m.Group+" "+name)
	}
	for group := range groups {
		sort.Strings(groups[group])
	}
	return writeJSON(describeAllOutput{
		GlobalFlags: globalFlagsMeta(),
		Groups:      groups,
	})
}

func globalFlagsMeta() []gen.FlagMeta {
	return []gen.FlagMeta{
		{Name: "from-json", Type: "bool", Help: "Read request body as JSON from stdin. When set, required flags are not enforced."},
		{Name: "request-json", Type: "string", Help: "Provide request body as inline JSON string. Alternative to --from-json with stdin."},
		{Name: "dry-run", Type: "bool", Help: "Print the request as JSON without calling the API."},
		{Name: "output", Type: "string", Help: "Output format: table or json (auto-detected from TTY).", EnumValues: []string{"table", "json"}},
		{Name: "context", Type: "string", Help: "Override the active context."},
	}
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
	sort.Slice(commands, func(i, j int) bool {
		left := commandPath(commands[i].Name)
		right := commandPath(commands[j].Name)
		if left == right {
			return commands[i].Name < commands[j].Name
		}
		return left < right
	})
	return writeJSON(commands)
}

func describeCommand(meta []gen.CommandMeta, group, command string) error {
	var matches []gen.CommandMeta
	for _, m := range meta {
		if m.Group != group {
			continue
		}
		if commandPath(m.Name) == command {
			matches = append(matches, m)
		}
	}

	if len(matches) == 1 {
		return writeJSON(matches[0])
	}
	if len(matches) > 1 {
		var options []string
		for _, m := range matches {
			options = appendUnique(options, commandPath(m.Name))
		}
		sort.Strings(options)
		return fmt.Errorf("ambiguous command %q in group %q, choose one of: %s", command, group, strings.Join(options, ", "))
	}

	// Backward-compatible fallback for callers that still pass only the base verb.
	var fallback []gen.CommandMeta
	for _, m := range meta {
		if m.Group != group {
			continue
		}
		path := commandPath(m.Name)
		if idx := strings.IndexByte(path, ' '); idx >= 0 {
			path = path[:idx]
		}
		if path == command {
			fallback = append(fallback, m)
		}
	}
	if len(fallback) == 1 {
		return writeJSON(fallback[0])
	}
	if len(fallback) > 1 {
		var options []string
		for _, m := range fallback {
			options = appendUnique(options, commandPath(m.Name))
		}
		sort.Strings(options)
		return fmt.Errorf("ambiguous command %q in group %q, choose one of: %s", command, group, strings.Join(options, ", "))
	}

	var available []string
	for _, m := range meta {
		if m.Group == group {
			available = appendUnique(available, commandPath(m.Name))
		}
	}
	if len(available) > 0 {
		sort.Strings(available)
		return fmt.Errorf("unknown command %q in group %q. Available commands: %s", command, group, strings.Join(available, ", "))
	}
	return fmt.Errorf("unknown command group %q", group)
}

func commandPath(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return ""
	}
	path := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			break
		}
		path = append(path, part)
	}
	if len(path) == 0 {
		return name
	}
	return strings.Join(path, " ")
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
