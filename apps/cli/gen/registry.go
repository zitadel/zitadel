package gen

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

// CmdFactory creates a cobra command group for a service.
type CmdFactory func(getCfg func() *config.Config, getOutput func() string) *cobra.Command

var factories []CmdFactory

// Register adds a command factory to the registry.
// Called from init() in each generated cmd_*.go file.
func Register(f CmdFactory) {
	factories = append(factories, f)
}

// RegisterAll adds all registered service command groups to the root command.
func RegisterAll(root *cobra.Command, getCfg func() *config.Config, getOutput func() string) {
	for _, f := range factories {
		root.AddCommand(f(getCfg, getOutput))
	}
}

// ArgMeta describes a positional argument.
type ArgMeta struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

// FlagMeta describes a CLI flag for schema introspection.
type FlagMeta struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Required   bool     `json:"required,omitempty"`
	Help       string   `json:"help,omitempty"`
	EnumValues []string `json:"enum_values,omitempty"`
	Group      string   `json:"group,omitempty"` // oneof variant name, e.g. "human" or "machine"
}

// CommandMeta holds machine-readable metadata for a single generated command.
type CommandMeta struct {
	Group          string          `json:"group"`
	Name           string          `json:"name"`
	Short          string          `json:"short"`
	Long           string          `json:"long,omitempty"`
	Example        string          `json:"example,omitempty"`
	FullMethodName string          `json:"method"`
	RequestType    string          `json:"request_type"`
	ResponseType   string          `json:"response_type"`
	PositionalArgs []ArgMeta       `json:"positional_args,omitempty"`
	Flags          []FlagMeta      `json:"flags,omitempty"`
	JSONTemplate   json.RawMessage `json:"json_template,omitempty"`
}

var allMeta []CommandMeta

// RegisterMeta adds command metadata to the registry.
// Called from init() in each generated cmd_*.go file.
func RegisterMeta(m CommandMeta) {
	allMeta = append(allMeta, m)
}

// AllMeta returns all registered command metadata.
func AllMeta() []CommandMeta {
	return allMeta
}
