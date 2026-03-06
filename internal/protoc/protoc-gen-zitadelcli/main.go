package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"

	annotations "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	//go:embed cmd.go.tmpl
	cmdTemplate []byte
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
}

// inlinedOneofGroup represents a proto oneof field whose variants have been expanded into
// individual, variant-prefixed flags instead of a compound selector+JSON-blob pair.
type inlinedOneofGroup struct {
	// GoName is the Go field name for the oneof on the request (e.g. "UserType").
	GoName string
	// KebabName is the kebab-case version (e.g. "user-type").
	KebabName string
	// Variants is the list of oneof message variants.
	Variants []inlinedVariant
}

// inlinedVariant represents one alternative within an inlinedOneofGroup.
type inlinedVariant struct {
	// VariantName is the proto field name, used as the subcommand name (e.g. "human", "machine").
	VariantName string
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
}

// columnDef describes a table column for rendering responses.
type columnDef struct {
	// Header is the uppercase column name for table output.
	Header string
	// GoAccessor is the Go getter chain (e.g. "GetOrganizationId()").
	GoAccessor string
	// IsTimestamp is true for google.protobuf.Timestamp columns.
	IsTimestamp bool
	// IsEnum is true for enum columns.
	IsEnum bool
}

// v2ServiceFilter is the list of v2 proto packages we want to generate CLI commands for.
var v2ServiceFilter = map[string]serviceConfig{
	"zitadel.org.v2": {
		resourceName: "orgs",
		resourceDesc: "organizations",
	},
	"zitadel.user.v2": {
		resourceName: "users",
		resourceDesc: "users",
	},
	"zitadel.project.v2": {
		resourceName: "projects",
		resourceDesc: "projects",
	},
	"zitadel.application.v2": {
		resourceName: "apps",
		resourceDesc: "applications",
	},
}

type serviceConfig struct {
	resourceName string
	resourceDesc string
}

func main() {
	input, _ := io.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(input, &req); err != nil {
		panic(err)
	}

	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	cmdTmpl := loadTemplate("cmd", cmdTemplate)

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}
		pkgName := string(file.Desc.Package())
		cfg, ok := v2ServiceFilter[pkgName]
		if !ok {
			continue
		}

		for _, service := range file.Services {
			sd := buildServiceData(service, file, cfg)
			if len(sd.Methods) == 0 {
				continue
			}

			// Generate the command file for this service.
			var buf bytes.Buffer
			if err := cmdTmpl.Execute(&buf, &sd); err != nil {
				panic(fmt.Sprintf("executing cmd template for %s: %v", sd.ServiceName, err))
			}
			outFile := plugin.NewGeneratedFile(
				fmt.Sprintf("cmd_%s.go", cfg.resourceName),
				protogen.GoImportPath(sd.GoImportPath),
			)
			outFile.Write(buf.Bytes())
		}
	}

	out, err := proto.Marshal(plugin.Response())
	if err != nil {
		panic(err)
	}
	fmt.Fprint(os.Stdout, string(out))
}

func buildServiceData(service *protogen.Service, file *protogen.File, cfg serviceConfig) serviceData {
	goImportPath := string(file.GoImportPath)
	goPkgName := string(file.GoPackageName)
	// Derive connect import path using the Go package name (from go_package option).
	// e.g. go_package "github.com/.../v2;org" → GoPackageName="org" → connectPath = ".../v2/orgconnect"
	connectPath := goImportPath + "/" + goPkgName + "connect"

	sd := serviceData{
		ServiceName:              service.GoName,
		ResourceName:             cfg.resourceName,
		ResourceDescription:      cfg.resourceDesc,
		GoImportPath:             goImportPath,
		GoImportAlias:            goPkgName + "pb",
		ConnectImportPath:        connectPath,
		ConnectImportAlias:       goPkgName + "connect",
		ConnectClientConstructor: "New" + service.GoName + "Client",
	}

	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue // Skip streaming methods.
		}

		md := buildMethodData(method, service)
		if md == nil {
			continue
		}
		md.Example = buildExample(cfg.resourceName, md)
		sd.Methods = append(sd.Methods, *md)
	}

	return sd
}

func buildMethodData(method *protogen.Method, service *protogen.Service) *methodData {
	rpcName := string(method.Desc.Name())

	// Skip deprecated methods.
	if method.Desc.Options() != nil {
		opts, ok := method.Desc.Options().(*descriptorpb.MethodOptions)
		if ok && opts != nil && opts.GetDeprecated() {
			return nil
		}
	}

	verb, suffix := rpcNameToVerbAndSuffix(rpcName, string(service.Desc.Name()))
	if verb == "" {
		return nil // Unknown method pattern — skip.
	}

	serviceName := string(service.Desc.Name())
	resourceSingular := inferResourceName(rpcName, serviceName)

	cliUse := verb
	if suffix != "" {
		cliUse = verb + "-" + suffix
	}

	md := &methodData{
		RPCName:        rpcName,
		Verb:           cliUse,
		Short:          humanizeRPCName(rpcName, resourceSingular),
		RequestType:    method.Input.GoIdent.GoName,
		ResponseType:   method.Output.GoIdent.GoName,
		FullMethodName: string(service.Desc.ParentFile().Package()) + "." + string(service.Desc.Name()) + "/" + rpcName,
	}

	// Populate Long from proto RPC comments.
	if comment := extractComment(method.Comments); comment != "" {
		md.Long = comment
	}

	// Extract flags from top-level (non-oneof) request fields.
	md.Flags = extractFlags(method.Input)

	// Extract inlined oneof groups.
	md.OneofGroups = extractOneofGroups(method.Input)
	md.HasOneofSubcmds = len(md.OneofGroups) == 1
	if len(md.OneofGroups) > 0 {
		var sb strings.Builder
		if md.Long != "" {
			sb.WriteString(md.Long)
		}
		for _, g := range md.OneofGroups {
			variantNames := make([]string, len(g.Variants))
			for i, v := range g.Variants {
				variantNames[i] = v.VariantName
			}
			sb.WriteString("\n\nChoose a sub-command: " + strings.Join(variantNames, ", ") + ".")
		}
		sb.WriteString("\nFor complex requests, pipe JSON body via: --from-json < request.json")
		md.Long = strings.TrimSpace(sb.String())
	}

	// Determine if there's a positional ID argument.
	md.IDArg, md.IDArgGoName, md.IDArgIsOptional = findIDField(method.Input, verb, suffix)
	if md.IDArg != "" && !md.HasOneofSubcmds {
		// No oneofs: ID is positional on the leaf command itself.
		md.Use = cliUse + " <" + md.IDArg + ">"
		// Remove the ID field from flags since it's positional.
		filtered := md.Flags[:0]
		for _, f := range md.Flags {
			if f.Name != toKebab(md.IDArg) {
				filtered = append(filtered, f)
			}
		}
		md.Flags = filtered
	} else if md.IDArg != "" {
		// Has oneofs: ID moves to each variant subcommand's Use.
		md.Use = cliUse
		filtered := md.Flags[:0]
		for _, f := range md.Flags {
			if f.Name != toKebab(md.IDArg) {
				filtered = append(filtered, f)
			}
		}
		md.Flags = filtered
	} else {
		md.Use = cliUse
	}

	// Extract response columns.
	md.IsListMethod, md.ListFieldGoName, md.ResponseColumns = extractResponseColumns(method.Output, verb)

	return md
}

// buildExample constructs a representative example invocation string for a command.
func buildExample(resourceName string, md *methodData) string {
	baseCmd := "  " + resourceName + " " + md.Verb
	if md.IDArg != "" {
		baseCmd += " <" + md.IDArg + ">"
	}

	var lines []string

	if len(md.OneofGroups) > 0 {
		for _, g := range md.OneofGroups {
			for _, v := range g.Variants {
				baseVariantCmd := "  " + resourceName + " " + md.Verb + " " + v.VariantName
				if md.IDArg != "" {
					baseVariantCmd += " <" + md.IDArg + ">"
				}
				if v.IsScalarBoolVariant {
					lines = append(lines, baseVariantCmd)
					continue
				}
				line := baseVariantCmd
				count := 0
				for _, f := range v.Flags {
					if count >= 3 {
						break
					}
					if f.FlagType == "bool" {
						// Skip bools in examples — they look confusing as defaults
						continue
					}
					if f.IsEnum && len(f.EnumValues) > 0 {
						line += " --" + f.FlagKebabName + " " + f.EnumValues[0]
					} else {
						line += " --" + f.FlagKebabName + " <" + f.FlagKebabName + ">"
					}
					count++
				}
				lines = append(lines, line)
			}
		}
		lines = append(lines, "  "+resourceName+" "+md.Verb+" --from-json < request.json")
	} else if len(md.Flags) > 0 {
		line := baseCmd
		count := 0
		for _, f := range md.Flags {
			if count >= 3 {
				break
			}
			if f.FlagType == "bool" || f.FlagType == "[]string" {
				continue
			}
			if f.IsEnum && len(f.EnumValues) > 0 {
				line += " --" + f.Name + " " + f.EnumValues[0]
			} else {
				line += " --" + f.Name + " <" + f.Name + ">"
			}
			count++
		}
		if count > 0 {
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}

// extractFlags returns flags for top-level, non-oneof request fields.
func extractFlags(msg *protogen.Message) []flagDef {
	var flags []flagDef
	for _, field := range msg.Fields {
		// Non-optional oneof fields are handled by extractOneofGroups instead.
		if field.Oneof != nil && !field.Desc.HasOptionalKeyword() {
			continue
		}
		fd := fieldToFlag(field, "")
		if fd == nil {
			continue
		}
		flags = append(flags, *fd)
	}
	return flags
}

// extractOneofGroups finds all non-optional oneof fields in msg and returns their inlined flag groups.
func extractOneofGroups(msg *protogen.Message) []inlinedOneofGroup {
	var groups []inlinedOneofGroup
	processed := map[string]bool{}
	for _, field := range msg.Fields {
		if field.Oneof == nil || field.Desc.HasOptionalKeyword() {
			continue
		}
		key := field.Oneof.GoName
		if processed[key] {
			continue
		}
		processed[key] = true
		g := buildOneofInlinedGroup(field.Oneof, msg)
		if g != nil {
			groups = append(groups, *g)
		}
	}
	return groups
}

// buildOneofInlinedGroup expands a proto oneof into per-variant, prefixed flags.
// Depth-0 scalar/enum fields on the variant message → --<variant>-<field>.
// Depth-0 message fields → expand their scalar/enum children → --<variant>-<child>
// (or --<variant>-<parent>-<child> if the leaf name would collide).
func buildOneofInlinedGroup(oneof *protogen.Oneof, msg *protogen.Message) *inlinedOneofGroup {
	g := &inlinedOneofGroup{
		GoName:    oneof.GoName,
		KebabName: toKebab(string(oneof.Desc.Name())),
	}

	for _, field := range oneof.Fields {
		if field.Desc.Kind() == protoreflect.BoolKind {
			// Scalar bool oneof field → zero-flag "presence" subcommand.
			// Invoking "is-verified <id>" sets the bool to true with no flags needed.
			// The Go oneof wrapper struct is e.g. SetEmailRequest_IsVerified{IsVerified: true}.
			parentGoName := string(msg.GoIdent.GoName)
			wrapperType := parentGoName + "_" + field.GoName
			variantName := toKebab(string(field.Desc.Name()))
			v := inlinedVariant{
				VariantName:         variantName,
				VarPrefix:           field.GoName,
				GoWrapperType:       wrapperType,
				GoFieldName:         oneof.GoName,
				IsScalarBoolVariant: true,
				ScalarGoFieldName:   field.GoName,
			}
			g.Variants = append(g.Variants, v)
			continue
		}
		if field.Desc.Kind() != protoreflect.MessageKind {
			continue // skip other scalar oneofs (string, int, etc.)
		}

		parentGoName := string(msg.GoIdent.GoName)
		wrapperType := parentGoName + "_" + field.GoName
		if field.Message.GoIdent.GoName == wrapperType {
			wrapperType += "_"
		}

		variantName := toKebab(string(field.Desc.Name()))
		varPrefix := field.GoName // e.g. "Human", "Machine"

		v := inlinedVariant{
			VariantName:   variantName,
			VarPrefix:     varPrefix,
			GoMsgType:     field.Message.GoIdent.GoName,
			GoWrapperType: wrapperType,
			GoFieldName:   field.GoName,
		}

		// Track used leaf names to detect collisions.
		leafNames := map[string]bool{}

		for _, varField := range field.Message.Fields {
			if isFieldDeprecated(varField) {
				continue
			}
			if varField.Desc.IsMap() || varField.Desc.IsList() {
				continue
			}
			// Skip non-optional oneof fields within the variant (too complex).
			if varField.Oneof != nil && !varField.Desc.HasOptionalKeyword() {
				continue
			}

			switch varField.Desc.Kind() {
			case protoreflect.MessageKind:
				// Depth-1: expand scalar/enum children of this sub-message.
				subMsg := varField.Message
				parentFieldGoName := varField.GoName  // e.g. "Profile"
				parentGoType := subMsg.GoIdent.GoName // e.g. "HumanProfile"

				for _, subField := range subMsg.Fields {
					if isFieldDeprecated(subField) {
						continue
					}
					if subField.Desc.IsMap() || subField.Desc.IsList() {
						continue
					}
					if subField.Oneof != nil && !subField.Desc.HasOptionalKeyword() {
						continue
					}
					if subField.Desc.Kind() == protoreflect.MessageKind {
						continue // stop at depth-1
					}

					leafKebab := toKebab(string(subField.Desc.Name()))
					flagName := leafKebab
					goVarSuffix := subField.GoName

					if leafNames[leafKebab] {
						// Collision: qualify with parent field name.
						parentKebab := toKebab(string(varField.Desc.Name()))
						flagName = parentKebab + "-" + leafKebab
						goVarSuffix = varField.GoName + subField.GoName
						leafNames[leafKebab] = true
					}

					fd := buildVariantFlag(subField, flagName, goVarSuffix, variantName, parentFieldGoName, parentGoType)
					if fd != nil {
						v.Flags = append(v.Flags, *fd)
					}
				}

			default:
				// Depth-0 scalar/enum field directly on the variant message.
				leafKebab := toKebab(string(varField.Desc.Name()))
				flagName := leafKebab
				goVarSuffix := varField.GoName

				if leafNames[leafKebab] {
					continue // collision at depth 0 (rare but skip)
				}
				leafNames[leafKebab] = true

				fd := buildVariantFlag(varField, flagName, goVarSuffix, variantName, "", "")
				if fd != nil {
					v.Flags = append(v.Flags, *fd)
				}
			}
		}

		g.Variants = append(g.Variants, v)
	}

	if len(g.Variants) == 0 {
		return nil
	}
	return g
}

// buildVariantFlag builds a variantFlagDef for a field inside an oneof variant.
// variantName is passed for context but no longer prefixes the help text (the
// subcommand name already conveys the variant context).
// parentGoField/parentGoType are non-empty for depth-1 fields.
func buildVariantFlag(field *protogen.Field, flagName, goVarSuffix, variantName, parentGoField, parentGoType string) *variantFlagDef {
	desc := field.Desc
	help := extractComment(field.Comments)
	if help == "" {
		help = strings.ReplaceAll(toKebab(string(desc.Name())), "-", " ")
	}

	isOpt := desc.HasOptionalKeyword()
	fd := &variantFlagDef{
		FlagKebabName:    flagName,
		GoVarSuffix:      goVarSuffix,
		ChildGoField:     field.GoName,
		ParentGoField:    parentGoField,
		ParentGoType:     parentGoType,
		Help:             help,
		IsOptionalScalar: isOpt,
		Required:         isRequired(field),
	}

	switch desc.Kind() {
	case protoreflect.StringKind:
		fd.FlagType = "string"
		fd.FlagFunc = "StringVar"
		fd.DefaultValue = `""`
	case protoreflect.BoolKind:
		fd.FlagType = "bool"
		fd.FlagFunc = "BoolVar"
		fd.DefaultValue = "false"
		fd.NeedChanged = isOpt
	case protoreflect.Int32Kind, protoreflect.Sint32Kind:
		fd.FlagType = "int32"
		fd.FlagFunc = "Int32Var"
		fd.DefaultValue = "0"
	case protoreflect.Uint32Kind:
		fd.FlagType = "uint32"
		fd.FlagFunc = "Uint32Var"
		fd.DefaultValue = "0"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind:
		fd.FlagType = "int64"
		fd.FlagFunc = "Int64Var"
		fd.DefaultValue = "0"
	case protoreflect.Uint64Kind:
		fd.FlagType = "uint64"
		fd.FlagFunc = "Uint64Var"
		fd.DefaultValue = "0"
	case protoreflect.EnumKind:
		fd.FlagType = "string"
		fd.FlagFunc = "StringVar"
		fd.DefaultValue = `""`
		fd.IsEnum = true
		fd.EnumGoType = string(desc.Enum().Name())
		for i := 0; i < desc.Enum().Values().Len(); i++ {
			fd.EnumValues = append(fd.EnumValues, string(desc.Enum().Values().Get(i).Name()))
		}
	default:
		return nil
	}

	return fd
}

// isFieldDeprecated returns true when the proto field is marked deprecated.
func isFieldDeprecated(field *protogen.Field) bool {
	if field.Desc.Options() == nil {
		return false
	}
	opts, ok := field.Desc.Options().(*descriptorpb.FieldOptions)
	return ok && opts != nil && opts.GetDeprecated()
}

func fieldToFlag(field *protogen.Field, prefix string) *flagDef {
	desc := field.Desc

	// Skip deprecated fields.
	if desc.Options() != nil {
		opts, ok := desc.Options().(*descriptorpb.FieldOptions)
		if ok && opts != nil && opts.GetDeprecated() {
			return nil
		}
	}

	// Skip map fields.
	if desc.IsMap() {
		return nil
	}

	// Non-optional oneof fields are handled by extractOneofGroups.
	if field.Oneof != nil && !field.Desc.HasOptionalKeyword() {
		return nil
	}

	name := string(desc.Name())
	if prefix != "" {
		name = prefix + "." + name
	}
	kebabName := toKebab(name)
	goName := field.GoName

	help := extractComment(field.Comments)
	if help == "" {
		help = "Set " + kebabName
	}

	// Check if REQUIRED (via google.api.field_behavior annotation).
	required := isRequired(field)

	fd := &flagDef{
		Name:             kebabName,
		GoName:           goName,
		Required:         required,
		Help:             help,
		ProtoFieldNumber: int(desc.Number()),
	}

	switch {
	case desc.IsList() && desc.Kind() == protoreflect.StringKind:
		fd.FlagType = "[]string"
		fd.FlagFunc = "StringSliceVar"
		fd.DefaultValue = "nil"
	case desc.IsList():
		return nil // Skip non-string repeated fields (too complex for CLI).
	case desc.Kind() == protoreflect.StringKind:
		fd.FlagType = "string"
		fd.FlagFunc = "StringVar"
		fd.DefaultValue = `""`
		fd.IsOptionalScalar = desc.HasOptionalKeyword()
	case desc.Kind() == protoreflect.BoolKind:
		fd.FlagType = "bool"
		fd.FlagFunc = "BoolVar"
		fd.DefaultValue = "false"
		fd.IsOptionalScalar = desc.HasOptionalKeyword()
	case desc.Kind() == protoreflect.Int32Kind || desc.Kind() == protoreflect.Sint32Kind:
		fd.FlagType = "int32"
		fd.FlagFunc = "Int32Var"
		fd.DefaultValue = "0"
	case desc.Kind() == protoreflect.Uint32Kind:
		fd.FlagType = "uint32"
		fd.FlagFunc = "Uint32Var"
		fd.DefaultValue = "0"
	case desc.Kind() == protoreflect.Int64Kind || desc.Kind() == protoreflect.Sint64Kind:
		fd.FlagType = "int64"
		fd.FlagFunc = "Int64Var"
		fd.DefaultValue = "0"
	case desc.Kind() == protoreflect.Uint64Kind:
		fd.FlagType = "uint64"
		fd.FlagFunc = "Uint64Var"
		fd.DefaultValue = "0"
	case desc.Kind() == protoreflect.EnumKind:
		fd.FlagType = "string"
		fd.FlagFunc = "StringVar"
		fd.DefaultValue = `""`
		fd.IsEnum = true
		fd.IsOptionalScalar = desc.HasOptionalKeyword()
		fd.EnumGoType = string(desc.Enum().Name())
		for i := 0; i < desc.Enum().Values().Len(); i++ {
			val := desc.Enum().Values().Get(i)
			fd.EnumValues = append(fd.EnumValues, string(val.Name()))
		}
	case desc.Kind() == protoreflect.MessageKind:
		// Skip complex nested messages — use --json for those.
		return nil
	default:
		return nil
	}

	return fd
}

func extractResponseColumns(msg *protogen.Message, verb string) (isList bool, listFieldGoName string, cols []columnDef) {
	// For List* methods, look for the repeated result field.
	for _, field := range msg.Fields {
		if field.Desc.IsList() && field.Desc.Kind() == protoreflect.MessageKind {
			isList = true
			listFieldGoName = field.GoName
			cols = extractMessageColumns(field.Message)
			return
		}
	}

	// For non-list methods, extract columns from top-level response fields.
	cols = extractTopLevelColumns(msg)
	if len(cols) > 0 {
		return false, "", cols
	}

	// If no scalar columns at top level, unwrap the first non-details message field
	// (e.g. GetUserByIDResponse.user, GetProjectResponse.project).
	for _, field := range msg.Fields {
		if field.Desc.Kind() == protoreflect.MessageKind && !field.Desc.IsList() && string(field.Desc.Name()) != "details" {
			innerCols := extractMessageColumns(field.Message)
			prefix := "Get" + field.GoName + "()."
			for i := range innerCols {
				innerCols[i].GoAccessor = prefix + innerCols[i].GoAccessor
			}
			return false, "", innerCols
		}
	}

	return false, "", cols
}

func extractMessageColumns(msg *protogen.Message) []columnDef {
	var cols []columnDef
	for _, field := range msg.Fields {
		col := fieldToColumn(field)
		if col != nil {
			cols = append(cols, *col)
		}
	}
	return cols
}

func extractTopLevelColumns(msg *protogen.Message) []columnDef {
	var cols []columnDef
	for _, field := range msg.Fields {
		col := fieldToColumn(field)
		if col != nil {
			cols = append(cols, *col)
		}
	}
	return cols
}

func fieldToColumn(field *protogen.Field) *columnDef {
	desc := field.Desc

	// Skip complex types.
	if desc.IsMap() || desc.IsList() {
		return nil
	}

	header := strings.ToUpper(toKebab(string(desc.Name())))
	header = strings.ReplaceAll(header, "-", " ")
	accessor := "Get" + field.GoName + "()"

	switch desc.Kind() {
	case protoreflect.StringKind, protoreflect.BoolKind,
		protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind:
		return &columnDef{Header: header, GoAccessor: accessor}
	case protoreflect.EnumKind:
		return &columnDef{Header: header, GoAccessor: accessor + ".String()", IsEnum: true}
	case protoreflect.MessageKind:
		// Check for well-known types.
		fullName := desc.Message().FullName()
		if fullName == "google.protobuf.Timestamp" {
			return &columnDef{Header: header, GoAccessor: accessor, IsTimestamp: true}
		}
		// For nested Details, extract ID if present.
		if string(desc.Name()) == "details" {
			return nil // Skip details metadata in table output.
		}
		return nil
	default:
		return nil
	}
}

// rpcNameToVerbAndSuffix maps a proto RPC name to a CLI verb and optional suffix.
// E.g. "AddOrganization" on OrganizationService → ("create", "")
//
//	"ListOrganizationDomains" on OrganizationService → ("list", "domains")
func rpcNameToVerbAndSuffix(name, serviceName string) (verb, suffix string) {
	// Strip "Service" suffix from service name for matching
	resource := strings.TrimSuffix(serviceName, "Service")

	prefixes := []struct {
		prefix string
		verb   string
	}{
		{"List", "list"},
		{"Create", "create"},
		{"Add", "create"},
		{"Get", "get"},
		{"Update", "update"},
		{"Delete", "delete"},
		{"Remove", "remove"},
		{"Deactivate", "deactivate"},
		{"Activate", "activate"},
		{"Reactivate", "reactivate"},
		{"Generate", "generate"},
		{"Set", "set"},
		{"Verify", "verify"},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p.prefix) {
			remainder := name[len(p.prefix):]
			// Strip the resource name (and its plural) from the remainder to get the sub-resource suffix.
			// E.g. "Organization" from "OrganizationDomains" → "Domains"
			// E.g. "Organization" from "Organizations" → "" (plural of resource = top-level list)
			remainder = strings.TrimPrefix(remainder, resource+"s")
			remainder = strings.TrimPrefix(remainder, resource)
			if remainder != "" {
				suffix = toKebab(remainder)
			}
			return p.verb, suffix
		}
	}
	return "", ""
}

// inferResourceName strips the verb prefix and service name to get the resource.
func inferResourceName(rpcName, serviceName string) string {
	// Remove common service suffix.
	resource := serviceName
	resource = strings.TrimSuffix(resource, "Service")
	return resource
}

// humanizeRPCName generates a short description from the RPC name.
func humanizeRPCName(rpcName, resource string) string {
	// Split CamelCase, keeping consecutive uppercase as one word (acronyms).
	var words []string
	runes := []rune(rpcName)
	start := 0
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			// If previous was also upper and next (if exists) is also upper or end, continue the acronym.
			if unicode.IsUpper(runes[i-1]) && (i+1 >= len(runes) || unicode.IsUpper(runes[i+1])) {
				continue
			}
			// If previous was upper but next is lower, the current char starts a new word.
			if unicode.IsUpper(runes[i-1]) {
				// The previous run was an acronym, split before current char.
				words = append(words, string(runes[start:i]))
				start = i
				continue
			}
			words = append(words, string(runes[start:i]))
			start = i
		}
	}
	words = append(words, string(runes[start:]))

	// Lowercase all words except known acronyms, and join.
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	result := strings.Join(words, " ")
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}
	return result
}

// findIDField looks for a field that should be a positional argument.
// For get/update/delete: the first *_id field is always positional.
// For scoped list/create (suffix != ""): the first *_id field is positional
// (e.g. list-passkeys <user_id>, create-idp-link <user_id>).
// For top-level list/create (suffix == ""): no positional ID.
func findIDField(msg *protogen.Message, verb, suffix string) (protoName, goName string, isOptional bool) {
	if (verb == "list" || verb == "create") && suffix == "" {
		return "", "", false
	}
	for _, field := range msg.Fields {
		name := string(field.Desc.Name())
		if field.Desc.Kind() == protoreflect.StringKind && strings.HasSuffix(name, "_id") {
			return name, field.GoName, field.Desc.HasOptionalKeyword()
		}
	}
	return "", "", false
}

func isRequired(field *protogen.Field) bool {
	opts := field.Desc.Options()
	if opts == nil {
		return false
	}
	behaviors := proto.GetExtension(opts, annotations.E_FieldBehavior)
	for _, b := range behaviors.([]annotations.FieldBehavior) {
		if b == annotations.FieldBehavior_REQUIRED {
			return true
		}
	}
	return false
}

// toKebab converts snake_case or CamelCase to kebab-case.
// Handles consecutive uppercase (acronyms) correctly: "OTPEmail" → "otp-email", "U2F" → "u2f".
func toKebab(s string) string {
	// First handle snake_case.
	s = strings.ReplaceAll(s, "_", "-")
	var result strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) && i > 0 {
			prev := runes[i-1]
			if prev == '-' {
				// Already have a separator.
			} else if unicode.IsUpper(prev) {
				// In a run of uppercase. Only add dash if next char is lowercase (end of acronym).
				if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result.WriteRune('-')
				}
			} else {
				result.WriteRune('-')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func extractComment(loc protogen.CommentSet) string {
	leading := strings.TrimSpace(string(loc.Leading))
	if leading == "" {
		return ""
	}
	// Take just the first line.
	if idx := strings.IndexByte(leading, '\n'); idx >= 0 {
		leading = leading[:idx]
	}
	return strings.TrimSpace(leading)
}

func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func loadTemplate(name string, data []byte) *template.Template {
	funcMap := template.FuncMap{
		"kebab":  toKebab,
		"lower":  strings.ToLower,
		"title":  titleCase,
		"quote":  func(s string) string { return fmt.Sprintf("%q", s) },
		"join":   strings.Join,
		"repeat": strings.Repeat,
	}
	return template.Must(template.New(name).Funcs(funcMap).Parse(string(data)))
}
