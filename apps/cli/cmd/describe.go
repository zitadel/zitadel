package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/runtime"
)

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe [group] [command ...]",
		Short: "Describe CLI commands as machine-readable JSON schema",
		Long: `Dump the schema of CLI commands as JSON for agent introspection.

Examples:
  zitadel-cli describe                     # list all command groups and commands
  zitadel-cli describe users               # list all commands in the users group
  zitadel-cli describe users create human  # describe a specific subcommand with full JSON schema`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				return describeAll()
			case 1:
				return describeGroup(args[0])
			default:
				return describeCommand(args[0], strings.Join(args[1:], " "))
			}
		},
		SilenceUsage: true,
	}
}

func describeAll() error {
	out := runtime.BuildDescribeAll()
	// Sort group keys for stable output.
	for group := range out.Groups {
		sort.Strings(out.Groups[group])
	}
	return writeJSON(out)
}

func describeGroup(group string) error {
	groupOut, specs := runtime.BuildDescribeGroup(group)
	if len(specs) == 0 {
		return fmt.Errorf("unknown command group %q", group)
	}
	return writeJSON(groupOut)
}

func describeCommand(group, command string) error {
	allSpecs := runtime.AllSpecs()

	// Try exact match first.
	for _, s := range allSpecs {
		if s.Group == group && s.Verb == command {
			out, err := runtime.BuildDescribeOutput(s)
			if err != nil {
				return err
			}
			return writeJSON(out)
		}
	}

	// If command contains spaces (e.g., "create human"), try matching the
	// first word as the verb. The extra words are variant/subcommand hints.
	parts := strings.Fields(command)
	if len(parts) > 1 {
		baseVerb := parts[0]
		for _, s := range allSpecs {
			if s.Group == group && s.Verb == baseVerb {
				out, err := runtime.BuildDescribeOutput(s)
				if err != nil {
					return err
				}
				return writeJSON(out)
			}
		}
	}

	// Build candidate list for error message.
	var candidates []string
	for _, s := range allSpecs {
		if s.Group == group {
			candidates = append(candidates, s.Verb)
		}
	}
	if len(candidates) > 0 {
		sort.Strings(candidates)
		return fmt.Errorf("unknown command %q in group %q. Available commands: %s", command, group, strings.Join(candidates, ", "))
	}
	return fmt.Errorf("unknown command group %q", group)
}

// writeJSON encodes v as indented JSON to stdout.
// TODO: accept io.Writer so callers can use cmd.OutOrStdout() instead of os.Stdout.
func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
