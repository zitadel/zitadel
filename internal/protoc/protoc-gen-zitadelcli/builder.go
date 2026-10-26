package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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

	var allExtraImports []extraImport
	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue // Skip streaming methods.
		}

		md, methodImports := buildMethodData(method, service, goImportPath)
		if md == nil {
			continue
		}
		md.Example = buildExample(cfg.resourceName, md)
		sd.Methods = append(sd.Methods, *md)
		allExtraImports = append(allExtraImports, methodImports...)
	}

	// Deduplicate extra imports
	importSet := map[string]extraImport{}
	for _, ei := range allExtraImports {
		importSet[ei.Path] = ei
	}
	for _, ei := range importSet {
		sd.ExtraImports = append(sd.ExtraImports, ei)
	}

	return sd
}

func buildMethodData(method *protogen.Method, service *protogen.Service, goImportPath string) (*methodData, []extraImport) {
	rpcName := string(method.Desc.Name())

	// Skip methods deprecated via the standard proto deprecated option.
	if isFieldDeprecated(method) {
		return nil, nil
	}

	// Detect OpenAPI-level deprecation (used by ZITADEL instead of standard proto deprecated).
	deprecated := isMethodDeprecatedOpenAPI(method)

	verb, suffix := rpcNameToVerbAndSuffix(rpcName, string(service.Desc.Name()))
	if verb == "" {
		return nil, nil // Unknown method pattern — skip.
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
		Deprecated:     deprecated,
	}
	if deprecated {
		md.Short = "[DEPRECATED] " + md.Short
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

	// Expand well-known message fields (e.g., ListQuery → offset, limit, asc flags)
	var extraImports []extraImport
	for _, field := range method.Input.Fields {
		if field.Oneof != nil && !field.Desc.HasOptionalKeyword() {
			continue
		}
		if field.Desc.Kind() == protoreflect.MessageKind && !field.Desc.IsList() && !field.Desc.IsMap() {
			expandMessageField(field, &md.Flags, goImportPath, &extraImports)
		}
	}

	// Build JSON template
	md.JSONTemplate = buildJSONTemplate(method.Input, 0)

	// Extract filter convenience flags
	filterFlags, filterImports := extractFilterConvenienceFlags(method.Input, goImportPath)
	md.FilterConvenience = filterFlags
	extraImports = append(extraImports, filterImports...)

	// Extract response columns.
	md.IsListMethod, md.ListFieldGoName, md.ListFieldProtoName, md.ResponseUnwrapField, md.ResponseColumns = extractResponseColumns(method.Output, verb)

	return md, extraImports
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

// buildJSONTemplate builds a JSON template showing all fields of a request message.
func buildJSONTemplate(msg *protogen.Message, depth int) string {
	if depth > 3 {
		return "{}"
	}
	var parts []string
	for _, field := range msg.Fields {
		name := jsonFieldName(field)
		switch {
		case field.Desc.IsMap():
			continue
		case field.Desc.IsList() && field.Desc.Kind() == protoreflect.MessageKind:
			parts = append(parts, fmt.Sprintf("%q: [%s]", name, buildJSONTemplate(field.Message, depth+1)))
		case field.Desc.IsList():
			parts = append(parts, fmt.Sprintf("%q: []", name))
		case field.Desc.Kind() == protoreflect.MessageKind:
			fullName := field.Desc.Message().FullName()
			if fullName == "google.protobuf.Timestamp" {
				parts = append(parts, fmt.Sprintf("%q: %q", name, "2006-01-02T15:04:05Z"))
			} else if fullName == "google.protobuf.Duration" {
				parts = append(parts, fmt.Sprintf("%q: %q", name, "3600s"))
			} else {
				parts = append(parts, fmt.Sprintf("%q: %s", name, buildJSONTemplate(field.Message, depth+1)))
			}
		case field.Desc.Kind() == protoreflect.BoolKind:
			parts = append(parts, fmt.Sprintf("%q: false", name))
		case field.Desc.Kind() == protoreflect.EnumKind:
			if field.Desc.Enum().Values().Len() > 0 {
				parts = append(parts, fmt.Sprintf("%q: %q", name, field.Desc.Enum().Values().Get(0).Name()))
			}
		case field.Desc.Kind() == protoreflect.StringKind:
			parts = append(parts, fmt.Sprintf("%q: %q", name, ""))
		default:
			parts = append(parts, fmt.Sprintf("%q: 0", name))
		}
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// buildJSONTemplateFiltered is like buildJSONTemplate but filters oneof fields at depth 0.
func buildJSONTemplateFiltered(msg *protogen.Message, depth int, oneofFilter, chosenField protoreflect.Name) string {
	if depth > 3 {
		return "{}"
	}
	var parts []string
	for _, field := range msg.Fields {
		// At depth 0, filter oneof fields
		if depth == 0 && oneofFilter != "" && field.Oneof != nil && !field.Desc.HasOptionalKeyword() {
			if field.Oneof.Desc.Name() != oneofFilter {
				continue
			}
			if field.Desc.Name() != chosenField {
				continue
			}
		}
		name := jsonFieldName(field)
		switch {
		case field.Desc.IsMap():
			continue
		case field.Desc.IsList() && field.Desc.Kind() == protoreflect.MessageKind:
			parts = append(parts, fmt.Sprintf("%q: [%s]", name, buildJSONTemplate(field.Message, depth+1)))
		case field.Desc.IsList():
			parts = append(parts, fmt.Sprintf("%q: []", name))
		case field.Desc.Kind() == protoreflect.MessageKind:
			fullName := field.Desc.Message().FullName()
			if fullName == "google.protobuf.Timestamp" {
				parts = append(parts, fmt.Sprintf("%q: %q", name, "2006-01-02T15:04:05Z"))
			} else if fullName == "google.protobuf.Duration" {
				parts = append(parts, fmt.Sprintf("%q: %q", name, "3600s"))
			} else {
				parts = append(parts, fmt.Sprintf("%q: %s", name, buildJSONTemplate(field.Message, depth+1)))
			}
		case field.Desc.Kind() == protoreflect.BoolKind:
			parts = append(parts, fmt.Sprintf("%q: false", name))
		case field.Desc.Kind() == protoreflect.EnumKind:
			if field.Desc.Enum().Values().Len() > 0 {
				parts = append(parts, fmt.Sprintf("%q: %q", name, field.Desc.Enum().Values().Get(0).Name()))
			}
		case field.Desc.Kind() == protoreflect.StringKind:
			parts = append(parts, fmt.Sprintf("%q: %q", name, ""))
		default:
			parts = append(parts, fmt.Sprintf("%q: 0", name))
		}
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
