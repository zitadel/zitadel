// Package runtime provides a reflection-based generic command runner for the ZITADEL CLI.
// Instead of generating per-method glue code, it uses protoreflect + dynamicpb to
// build Cobra commands, map flags, and call ConnectRPC dynamically.
package runtime

// CommandSpec is the declarative metadata for a single CLI command.
// It is produced by the code generator (or hand-written) and contains
// no executable logic — only the data needed to wire up a Cobra command.
type CommandSpec struct {
	// Group is the CLI noun (e.g. "users", "orgs").
	Group string
	// GroupDescription is shown in --help for the group command.
	GroupDescription string
	// Verb is the CLI verb with optional suffix (e.g. "create-human", "list", "get").
	Verb string
	// FullMethodName is "package.Service/Method" (e.g. "zitadel.user.v2.UserService/ListUsers").
	FullMethodName string
	// Short is the one-line description.
	Short string
	// Long is the multi-line description (optional).
	Long string
	// Example is a ready-to-run example (optional).
	Example string
	// Deprecated marks the command as deprecated with a warning.
	Deprecated bool

	// PositionalArgs lists args that come before flags (e.g. <user_id>).
	PositionalArgs []PosArg
	// HasOneofSubcmds is true when the request has a top-level oneof that
	// is expanded into variant subcommands (e.g. "create human", "create machine").
	HasOneofSubcmds bool
	// OneofGroups lists variant subcommands for the top-level oneof.
	OneofGroups []OneofGroupSpec

	// TableColumns defines how to render the response in table mode.
	// Each column specifies a proto field path and display header.
	TableColumns []ColumnSpec
	// IsListMethod is true for List* RPCs that return repeated results.
	IsListMethod bool
	// ListFieldName is the proto field name of the repeated result field in the response.
	ListFieldName string
	// ResponseUnwrapField is the proto field name to unwrap for single-resource
	// responses (e.g. "user" in GetUserByIDResponse). Empty means use top-level.
	ResponseUnwrapField string

	// FilterConvenience holds convenience flag specs for common filter patterns.
	FilterConvenience []FilterConvenienceSpec
}

// PosArg describes a positional argument.
type PosArg struct {
	// ProtoFieldName is the proto field name (e.g. "user_id").
	ProtoFieldName string
	// Required is true if the arg must be provided.
	Required bool
	// IsOptional is true for proto3 optional fields (Go *string).
	IsOptional bool
}

// OneofGroupSpec describes one oneof field expanded into variant subcommands.
type OneofGroupSpec struct {
	// ProtoOneofName is the proto oneof name (e.g. "type").
	ProtoOneofName string
	// Variants lists the available subcommands for this oneof.
	Variants []VariantSpec
}

// VariantSpec describes one variant within a oneof group.
type VariantSpec struct {
	// CliName is the kebab-case subcommand name (e.g. "human", "machine").
	CliName string
	// ProtoFieldName is the proto field name within the oneof (e.g. "human").
	ProtoFieldName string
}

// ColumnSpec describes a single table column.
type ColumnSpec struct {
	// Header is the uppercase column header (e.g. "ID", "ORGANIZATION ID").
	Header string
	// FieldPath is a dot-separated proto field path (e.g. "details.resource_owner").
	// For list methods, the path is relative to each element.
	FieldPath string
	// IsTimestamp renders the value as RFC3339.
	IsTimestamp bool
	// IsEnum renders the value using the enum name instead of the number.
	IsEnum bool
}

// FilterConvenienceSpec describes a convenience flag that maps to a search filter.
type FilterConvenienceSpec struct {
	// FlagName is the kebab-case flag name (e.g. "user-id").
	FlagName string
	// Help is the flag help text.
	Help string
	// FilterListFieldName is the proto field name of the repeated filter list.
	FilterListFieldName string
	// FilterOneofFieldName is the proto field name within the filter's oneof.
	FilterOneofFieldName string
}
