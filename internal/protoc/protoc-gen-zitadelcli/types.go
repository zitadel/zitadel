package main

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

// serviceData holds all data for generating commands for one proto service.
type serviceData struct {
	// ServiceName is the Go name of the service (e.g. "OrganizationService").
	ServiceName string
	// ResourceName is the kebab-case short name for the CLI group (e.g. "orgs").
	ResourceName string
	// ResourceDescription is the human-readable name (e.g. "organizations").
	ResourceDescription string
	// GoImportPath is the Go import path for the proto types (e.g. "github.com/zitadel/zitadel/pkg/grpc/org/v2").
	GoImportPath string
	// GoImportAlias is a short alias for the import (e.g. "orgpb").
	GoImportAlias string
	// ConnectImportPath is the Go import path for the connect client.
	ConnectImportPath string
	// ConnectImportAlias is a short alias (e.g. "orgconnect").
	ConnectImportAlias string
	// ConnectClientConstructor is the constructor name (e.g. "NewOrganizationServiceClient").
	ConnectClientConstructor string
	// Methods is the list of RPC methods to generate commands for.
	Methods []methodData
	// ExtraImports holds additional Go imports needed by generated code (e.g., filter types).
	ExtraImports []extraImport
}

// extraImport represents an additional Go import for the generated file.
type extraImport struct {
	Alias string
	Path  string
}

// methodData holds data for generating a single CLI command from an RPC method.
type methodData struct {
	// RPCName is the proto RPC name (e.g. "AddOrganization").
	RPCName string
	// Verb is the CLI verb (e.g. "create", "list", "get", "update", "delete").
	Verb string
	// Use is the cobra Use string (e.g. "create --name NAME", "get <id>").
	Use string
	// Short is the short description.
	Short string
	// Long is the long description, populated from proto RPC comments.
	Long string
	// Example is a ready-to-run example invocation.
	Example string
	// RequestType is the Go type name for the request (e.g. "AddOrganizationRequest").
	RequestType string
	// ResponseType is the Go type name for the response.
	ResponseType string
	// FullMethodName is the fully-qualified gRPC method name (e.g. "zitadel.user.v2.UserService/GetUserByID").
	FullMethodName string
	// Flags is the list of flags derived from top-level request fields (no oneof groups).
	Flags []flagDef
	// OneofGroups holds the per-variant subcommand data for each oneof field in the request.
	OneofGroups []inlinedOneofGroup
	// HasOneofSubcmds is true when this method uses the subcommand pattern (len(OneofGroups)==1).
	// When true, the parent command has no RunE except for --from-json.
	HasOneofSubcmds bool
	// IDArg is non-empty if the method takes a positional ID argument (the proto field name).
	IDArg string
	// IDArgGoName is the Go getter name for the ID field.
	IDArgGoName string
	// IDArgIsOptional is true when the ID field is a proto3 optional (Go *string).
	IDArgIsOptional bool
	// ResponseColumns defines table columns for rendering responses.
	ResponseColumns []columnDef
	// IsListMethod is true if this is a List* method with repeated results.
	IsListMethod bool
	// ListFieldGoName is the Go name of the repeated field in the response (e.g. "Result", "Users", "Projects").
	ListFieldGoName string
	// ListFieldProtoName is the proto field name of the repeated result (e.g. "result", "users").
	ListFieldProtoName string
	// ResponseUnwrapField is the proto field name to unwrap for single-resource responses
	// (e.g. "user" in GetUserByIDResponse). Empty means use top-level.
	ResponseUnwrapField string
	// JSONTemplate is a JSON template string showing all fields of the request.
	JSONTemplate string
	// FilterConvenience holds convenience flags for common filter patterns (e.g., --user-id for IDFilter).
	FilterConvenience []filterConvenienceFlag
	// Deprecated is true when the RPC is marked deprecated in its OpenAPI annotation.
	Deprecated bool
}

// filterConvenienceFlag describes a convenience flag that maps to a search filter.
type filterConvenienceFlag struct {
	FlagName            string
	Help                string
	OneofFieldName      string
	FilterListField     string
	FilterMsgType       string
	WrapperType         string
	OneofGoFieldName    string
	IDFilterImportAlias string
	IDFilterType        string
}

// inlinedOneofGroup represents a proto oneof field whose variants have been expanded into
// individual, variant-prefixed flags instead of a compound selector+JSON-blob pair.
type inlinedOneofGroup struct {
	// GoName is the Go field name for the oneof on the request (e.g. "UserType").
	GoName string
	// ProtoName is the proto oneof name (e.g. "type" or "user_type").
	ProtoName string
	// KebabName is the kebab-case version (e.g. "user-type").
	KebabName string
	// Variants is the list of oneof message variants.
	Variants []inlinedVariant
}

// inlinedVariant represents one alternative within an inlinedOneofGroup.
type inlinedVariant struct {
	// VariantName is the proto field name, used as the subcommand name (e.g. "human", "machine").
	VariantName string
	// ProtoFieldName is the proto field name within the oneof (may differ from VariantName if kebab-converted).
	ProtoFieldName string
	// VarPrefix is the Go-style capitalized prefix for variable names (e.g. "Human", "Machine").
	VarPrefix string
	// GoMsgType is the Go type of the variant message (e.g. "CreateUserRequest_Human").
	GoMsgType string
	// GoWrapperType is the Go oneof wrapper type (e.g. "CreateUserRequest_Human_").
	GoWrapperType string
	// GoFieldName is the field on the wrapper struct (e.g. "Human").
	GoFieldName string
	// Flags are the unprefixed flags for this variant subcommand.
	Flags []variantFlagDef
	// IsScalarBoolVariant is true when the oneof field is a plain bool (not a message).
	// The subcommand takes no flags; invoking it sets the bool field to true.
	IsScalarBoolVariant bool
	// ScalarGoFieldName is the Go field name set to true on the request for scalar bool variants.
	ScalarGoFieldName string
	// IsScalarStringVariant is true when the oneof field is a plain string.
	// The subcommand takes the string as a single positional argument.
	IsScalarStringVariant bool
	// ScalarStringGoFieldName is the Go wrapper field name set on the request for string variants.
	ScalarStringGoFieldName string
	// JSONTemplate is a JSON template for this specific variant (filtered by the oneof choice).
	JSONTemplate string
	// NestedOneofs holds metadata about nested oneof fields expanded within this variant.
	NestedOneofs []nestedOneofMeta
}

// nestedOneofMeta describes a nested oneof within a variant that has been expanded to flags.
type nestedOneofMeta struct {
	GoName   string
	Variants []string
}

// variantFlagDef describes one flag on a variant subcommand.
// Flag names are unprefixed (e.g. "given-name", not "human-given-name") because
// the subcommand itself provides the type context.
type variantFlagDef struct {
	// FlagKebabName is the flag name without variant prefix (e.g. "given-name").
	FlagKebabName string
	// GoVarSuffix is the Go variable suffix within the variant's scope (e.g. "GivenName").
	GoVarSuffix string
	// ChildGoField is the Go field name to set on the message (e.g. "GivenName").
	ChildGoField string
	// ParentGoField is non-empty for depth-1 fields; the intermediate Go field name (e.g. "Profile").
	ParentGoField string
	// ParentGoType is the Go type for lazy-init of the parent message (e.g. "HumanProfile").
	ParentGoType string
	// Help is the flag description.
	Help string
	// FlagType is the Go type (e.g. "string", "bool", "int32").
	FlagType string
	// FlagFunc is the cobra flag function (e.g. "StringVar").
	FlagFunc string
	// DefaultValue is the zero value as a Go literal (e.g. `""`, "false", "0").
	DefaultValue string
	// IsOptionalScalar is true for proto3 optional scalars (need pointer assignment).
	IsOptionalScalar bool
	// NeedChanged is true when we must use cmd.Flags().Changed() to detect explicit setting
	// (required for optional bool fields where false == default).
	NeedChanged bool
	// IsEnum is true for proto enum fields.
	IsEnum bool
	// EnumGoType is the Go enum type name (e.g. "Gender").
	EnumGoType string
	// EnumValues is the list of valid enum value names.
	EnumValues []string
	// Required is true when the field has (google.api.field_behavior) = REQUIRED.
	Required bool
	// NestedOneofGroup is the Go name of the nested oneof this flag belongs to (empty for normal flags).
	NestedOneofGroup string
	// NestedOneofVariant is the Go field name of the chosen variant within the nested oneof.
	NestedOneofVariant string
	// NestedOneofWrapperType is the Go oneof wrapper type for the nested oneof variant.
	NestedOneofWrapperType string
	// NestedOneofMsgType is the Go message type inside the wrapper (empty for scalar variants).
	NestedOneofMsgType string
}

// flagDef describes a CLI flag derived from a top-level (non-oneof) proto field.
type flagDef struct {
	// Name is the kebab-case flag name (e.g. "organization-id").
	Name string
	// GoName is the Go field name (e.g. "OrganizationId").
	GoName string
	// FlagType is the Go type for the flag variable (e.g. "string", "bool", "int32", "[]string").
	FlagType string
	// FlagFunc is the cobra flag registration function (e.g. "StringVar", "BoolVar", "Int32Var", "StringSliceVar").
	FlagFunc string
	// Required is true if the field has REQUIRED behavior.
	Required bool
	// Help is the flag help text (from proto field comments).
	Help string
	// IsOptionalScalar is true for proto optional scalars (need pointer wrappers).
	IsOptionalScalar bool
	// ProtoFieldNumber is the proto field number.
	ProtoFieldNumber int
	// Deprecated is true for deprecated fields.
	Deprecated bool
	// DefaultValue is the default for the flag type ("", false, 0, etc.).
	DefaultValue string
	// IsEnum is true if this field is a proto enum.
	IsEnum bool
	// EnumGoType is the Go type of the enum (e.g. "OIDCApplicationType").
	EnumGoType string
	// EnumValues lists the valid enum value names.
	EnumValues []string
	// ParentGoField is non-empty for expanded message fields; the parent Go field name.
	ParentGoField string
	// ParentGoType is the Go type for lazy-init of the parent message.
	ParentGoType string
}

// columnDef describes a table column for rendering responses.
type columnDef struct {
	// Header is the uppercase column name for table output.
	Header string
	// GoAccessor is the Go getter chain (e.g. "GetOrganizationId()").
	GoAccessor string
	// FieldPath is a dot-separated proto field path (e.g. "details.resource_owner").
	// Used by the dynamic runtime for table rendering.
	FieldPath string
	// IsTimestamp is true for google.protobuf.Timestamp columns.
	IsTimestamp bool
	// IsEnum is true for enum columns.
	IsEnum bool
	// IsOneofType is true for oneof discriminator columns that show which variant is set.
	IsOneofType bool
	// OneofVariants maps Go wrapper type names to human-readable variant names (e.g. "AuthFactor_Otp" → "otp").
	OneofVariants []oneofVariantColumn
}

// oneofVariantColumn maps a Go oneof wrapper type to a display name for table output.
type oneofVariantColumn struct {
	GoWrapperType string // e.g. "*userpb.AuthFactor_Otp"
	DisplayName   string // e.g. "otp"
}

// serviceConfig holds basic naming config for a proto service CLI group.
type serviceConfig struct {
	resourceName string
	resourceDesc string
}

// v2ServiceFilter is the list of v2 proto packages we want to generate CLI commands for.
var v2ServiceFilter = map[string]serviceConfig{
	"zitadel.action.v2": {
		resourceName: "actions",
		resourceDesc: "actions and executions",
	},
	"zitadel.application.v2": {
		resourceName: "apps",
		resourceDesc: "applications",
	},
	"zitadel.authorization.v2": {
		resourceName: "authorizations",
		resourceDesc: "user authorizations",
	},
	"zitadel.feature.v2": {
		resourceName: "features",
		resourceDesc: "instance and organization features",
	},
	"zitadel.group.v2": {
		resourceName: "groups",
		resourceDesc: "user groups",
	},
	"zitadel.idp.v2": {
		resourceName: "idps",
		resourceDesc: "identity providers",
	},
	"zitadel.instance.v2": {
		resourceName: "instances",
		resourceDesc: "ZITADEL instances",
	},
	"zitadel.oidc.v2": {
		resourceName: "oidc",
		resourceDesc: "OIDC introspection and token exchange",
	},
	"zitadel.org.v2": {
		resourceName: "orgs",
		resourceDesc: "organizations",
	},
	"zitadel.project.v2": {
		resourceName: "projects",
		resourceDesc: "projects",
	},
	"zitadel.saml.v2": {
		resourceName: "saml",
		resourceDesc: "SAML service provider metadata",
	},
	"zitadel.session.v2": {
		resourceName: "sessions",
		resourceDesc: "user sessions",
	},
	"zitadel.settings.v2": {
		resourceName: "settings",
		resourceDesc: "instance and organization settings",
	},
	"zitadel.user.v2": {
		resourceName: "users",
		resourceDesc: "users",
	},
	"zitadel.webkey.v2": {
		resourceName: "webkeys",
		resourceDesc: "web keys for OIDC/SAML signing",
	},
}

// expandableMessages lists well-known message types that should be expanded into flat flags.
var expandableMessages = map[protoreflect.FullName]bool{
	"zitadel.object.v2.ListQuery": true,
}
