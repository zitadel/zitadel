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
	// RequestType is the Go type name for the request (e.g. "AddOrganizationRequest").
	RequestType string
	// ResponseType is the Go type name for the response.
	ResponseType string
	// FullMethodName is the fully-qualified gRPC method name (e.g. "zitadel.user.v2.UserService/GetUserByID").
	FullMethodName string
	// Flags is the list of flags derived from request fields.
	Flags []flagDef
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

// flagDef describes a CLI flag derived from a proto field.
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
	// IsOneofSelector is true when this flag selects a oneof variant by name.
	// Two flags are generated: --<Name> (enum) and --<Name>-data (optional JSON for the variant's fields).
	IsOneofSelector bool
	// GoOneofName is the Go field name for the oneof interface on the request (e.g. "UserType").
	GoOneofName string
	// OneofVariants lists all message variants in the oneof group.
	OneofVariants []oneofVariant
}

// oneofVariant describes one message variant within a proto oneof.
type oneofVariant struct {
	// VariantName is the proto field name used as the flag value (e.g. "human", "machine").
	VariantName string
	// GoMsgType is the Go message type (e.g. "CreateUserRequest_Human").
	GoMsgType string
	// GoWrapperType is the Go oneof wrapper type (e.g. "CreateUserRequest_Human_").
	GoWrapperType string
	// GoFieldName is the field name on the wrapper struct (e.g. "Human").
	GoFieldName string
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

	// Extract flags from request message fields.
	md.Flags = extractFlags(method.Input)

	// Determine if there's a positional ID argument.
	md.IDArg, md.IDArgGoName, md.IDArgIsOptional = findIDField(method.Input, verb, suffix)
	if md.IDArg != "" {
		md.Use = cliUse + " <" + md.IDArg + ">"
		// Remove the ID field from flags since it's positional.
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

func extractFlags(msg *protogen.Message) []flagDef {
	var flags []flagDef
	processedOneofs := map[string]bool{}
	for _, field := range msg.Fields {
		// Non-optional oneof groups: emit one selector flagDef for the whole group.
		if field.Oneof != nil && !field.Desc.HasOptionalKeyword() {
			key := field.Oneof.GoName
			if processedOneofs[key] {
				continue
			}
			processedOneofs[key] = true
			fd := buildOneofSelectorFlag(field.Oneof, msg)
			if fd != nil {
				flags = append(flags, *fd)
			}
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

// buildOneofSelectorFlag generates a single IsOneofSelector flagDef for all message
// variants in a oneof group. The caller is responsible for deduplication.
func buildOneofSelectorFlag(oneof *protogen.Oneof, msg *protogen.Message) *flagDef {
	var variants []oneofVariant
	for _, field := range oneof.Fields {
		if field.Desc.Kind() != protoreflect.MessageKind {
			continue // skip rare scalar oneof alternatives
		}
		parentGoName := string(msg.GoIdent.GoName)
		wrapperType := parentGoName + "_" + field.GoName
		// Go protoc adds a trailing _ when wrapper base name equals the message type name.
		if field.Message.GoIdent.GoName == wrapperType {
			wrapperType += "_"
		}
		variants = append(variants, oneofVariant{
			VariantName:   string(field.Desc.Name()),
			GoMsgType:     field.Message.GoIdent.GoName,
			GoWrapperType: wrapperType,
			GoFieldName:   field.GoName,
		})
	}
	if len(variants) == 0 {
		return nil
	}

	var variantNames []string
	for _, v := range variants {
		variantNames = append(variantNames, v.VariantName)
	}
	kebabName := toKebab(string(oneof.Desc.Name()))
	return &flagDef{
		Name:            kebabName,
		GoName:          oneof.GoName,
		FlagType:        "string",
		FlagFunc:        "StringVar",
		DefaultValue:    `""`,
		Help:            "Select " + kebabName + ": " + strings.Join(variantNames, " or ") + " (use --" + kebabName + "-data for type-specific fields as JSON)",
		IsOneofSelector: true,
		GoOneofName:     oneof.GoName,
		OneofVariants:   variants,
		EnumValues:      variantNames,
	}
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

	// Non-optional oneof fields are handled by buildOneofSelectorFlag in extractFlags.
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
	// Heuristic: fields named "name" or ending in "_id" in Create/Add methods are typically required.
	// Full annotation parsing of google.api.field_behavior would require importing the extension,
	// which adds complexity. We rely on cobra MarkFlagRequired in the template for key fields.
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
