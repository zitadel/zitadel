package main

import (
	"fmt"
	"strings"

	annotations "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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
func buildOneofInlinedGroup(oneof *protogen.Oneof, msg *protogen.Message) *inlinedOneofGroup {
	g := &inlinedOneofGroup{
		GoName:    oneof.GoName,
		ProtoName: string(oneof.Desc.Name()),
		KebabName: toKebab(string(oneof.Desc.Name())),
	}

	for _, field := range oneof.Fields {
		if field.Desc.Kind() == protoreflect.BoolKind {
			parentGoName := string(msg.GoIdent.GoName)
			wrapperType := parentGoName + "_" + field.GoName
			variantName := toKebab(string(field.Desc.Name()))
			v := inlinedVariant{
				VariantName:         variantName,
				ProtoFieldName:      string(field.Desc.Name()),
				VarPrefix:           field.GoName,
				GoWrapperType:       wrapperType,
				GoFieldName:         oneof.GoName,
				IsScalarBoolVariant: true,
				ScalarGoFieldName:   field.GoName,
			}
			v.JSONTemplate = buildJSONTemplateFiltered(msg, 0, oneof.Desc.Name(), field.Desc.Name())
			g.Variants = append(g.Variants, v)
			continue
		}
		if field.Desc.Kind() == protoreflect.StringKind {
			parentGoName := string(msg.GoIdent.GoName)
			wrapperType := parentGoName + "_" + field.GoName
			variantName := toKebab(string(field.Desc.Name()))
			v := inlinedVariant{
				VariantName:             variantName,
				ProtoFieldName:          string(field.Desc.Name()),
				VarPrefix:               field.GoName,
				GoWrapperType:           wrapperType,
				GoFieldName:             oneof.GoName,
				IsScalarStringVariant:   true,
				ScalarStringGoFieldName: field.GoName,
			}
			v.JSONTemplate = buildJSONTemplateFiltered(msg, 0, oneof.Desc.Name(), field.Desc.Name())
			g.Variants = append(g.Variants, v)
			continue
		}
		if field.Desc.Kind() != protoreflect.MessageKind {
			continue // skip other scalar oneofs (int, bytes, etc.)
		}

		parentGoName := string(msg.GoIdent.GoName)
		wrapperType := parentGoName + "_" + field.GoName
		if field.Message.GoIdent.GoName == wrapperType {
			wrapperType += "_"
		}

		variantName := toKebab(string(field.Desc.Name()))
		varPrefix := field.GoName // e.g. "Human", "Machine"

		v := inlinedVariant{
			VariantName:    variantName,
			ProtoFieldName: string(field.Desc.Name()),
			VarPrefix:      varPrefix,
			GoMsgType:     field.Message.GoIdent.GoName,
			GoWrapperType: wrapperType,
			GoFieldName:   field.GoName,
		}

		// Track used leaf names to detect collisions.
		leafNames := map[string]bool{}
		processedNestedOneofs := map[string]bool{}

		for _, varField := range field.Message.Fields {
			if isFieldDeprecated(varField) {
				continue
			}
			if varField.Desc.IsMap() || varField.Desc.IsList() {
				continue
			}
			// Skip non-optional oneof fields within the variant (too complex),
			// unless they can be expanded to flags.
			if varField.Oneof != nil && !varField.Desc.HasOptionalKeyword() {
				neoKey := varField.Oneof.GoName
				if !processedNestedOneofs[neoKey] && isExpandableOneof(varField.Oneof) {
					processedNestedOneofs[neoKey] = true
					expandNestedOneof(varField.Oneof, &v, variantName, field.Message)
				}
				continue
			}

			switch varField.Desc.Kind() {
			case protoreflect.MessageKind:
				// Depth-1: expand scalar/enum children of this sub-message.
				// Skip well-known google.protobuf.* types — they cannot be expanded to simple flags;
				// users should use --request-json for these (e.g. google.protobuf.Duration).
				subMsg := varField.Message
				if isWellKnownProtoType(subMsg) {
					continue
				}
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

		v.JSONTemplate = buildJSONTemplateFiltered(msg, 0, oneof.Desc.Name(), field.Desc.Name())
		g.Variants = append(g.Variants, v)
	}

	if len(g.Variants) == 0 {
		return nil
	}
	return g
}

// buildVariantFlag builds a variantFlagDef for a field inside an oneof variant.
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

func fieldToFlag(field *protogen.Field, prefix string) *flagDef {
	desc := field.Desc

	// Skip deprecated fields.
	if isFieldDeprecated(field) {
		return nil
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
		// Append short valid values to help text so users know what to pass
		// even without tab-completion. Strip the common prefix (e.g., "USER_FIELD_NAME_")
		// and exclude the UNSPECIFIED sentinel.
		if shortVals := enumShortValues(fd.EnumValues); len(shortVals) > 0 {
			fd.Help += " (one of: " + strings.Join(shortVals, ", ") + ")"
		}
	case desc.Kind() == protoreflect.MessageKind:
		// Skip complex nested messages — use --json for those.
		return nil
	default:
		return nil
	}

	return fd
}

// enumShortValues strips the longest common prefix from enum values and excludes
// the UNSPECIFIED sentinel, returning human-readable choices for help text.
func enumShortValues(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	// Find longest common prefix (up to a trailing underscore).
	prefix := values[0]
	for _, v := range values[1:] {
		for prefix != "" && !strings.HasPrefix(v, prefix) {
			if idx := strings.LastIndex(prefix[:len(prefix)-1], "_"); idx >= 0 {
				prefix = prefix[:idx+1]
			} else {
				prefix = ""
			}
		}
	}
	var out []string
	for _, v := range values {
		short := strings.TrimPrefix(v, prefix)
		if strings.HasSuffix(short, "UNSPECIFIED") {
			continue
		}
		out = append(out, short)
	}
	return out
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

// findIDField looks for a field that should be a positional argument.
func findIDField(msg *protogen.Message, verb, suffix string) (protoName, goName string, isOptional bool) {
	if (verb == "list" || verb == "create") && suffix == "" {
		return "", "", false
	}
	for _, field := range msg.Fields {
		// Skip fields inside a real oneof — they're variant discriminators, not standalone IDs.
		if field.Oneof != nil && !field.Desc.HasOptionalKeyword() {
			continue
		}
		name := string(field.Desc.Name())
		if field.Desc.Kind() == protoreflect.StringKind && strings.HasSuffix(name, "_id") {
			return name, field.GoName, field.Desc.HasOptionalKeyword()
		}
	}
	return "", "", false
}

// expandMessageField expands well-known message types (e.g., ListQuery) into flat flags.
func expandMessageField(field *protogen.Field, flags *[]flagDef, goImportPath string, extraImports *[]extraImport) {
	if field.Desc.Kind() != protoreflect.MessageKind {
		return
	}
	fullName := field.Desc.Message().FullName()
	if !expandableMessages[fullName] {
		return
	}
	subMsg := field.Message
	parentGoField := field.GoName
	parentGoType := subMsg.GoIdent.GoName

	// Add import for the message's package if different from the service package
	msgImportPath := string(subMsg.GoIdent.GoImportPath)
	if msgImportPath != goImportPath {
		alias := importAlias(msgImportPath)
		parentGoType = alias + "." + parentGoType
		// Check if already added
		found := false
		for _, ei := range *extraImports {
			if ei.Path == msgImportPath {
				found = true
				break
			}
		}
		if !found {
			*extraImports = append(*extraImports, extraImport{Alias: alias, Path: msgImportPath})
		}
	}

	for _, subField := range subMsg.Fields {
		if isFieldDeprecated(subField) {
			continue
		}
		if subField.Desc.Kind() == protoreflect.MessageKind {
			continue
		}
		if subField.Oneof != nil && !subField.Desc.HasOptionalKeyword() {
			continue
		}
		fd := fieldToFlag(subField, "")
		if fd == nil {
			continue
		}
		fd.ParentGoField = parentGoField
		fd.ParentGoType = parentGoType
		*flags = append(*flags, *fd)
	}
}

// isExpandableOneof checks if a nested oneof can be expanded to flags.
func isExpandableOneof(oneof *protogen.Oneof) bool {
	for _, field := range oneof.Fields {
		if field.Desc.Kind() == protoreflect.BoolKind || field.Desc.Kind() == protoreflect.StringKind {
			continue // scalar variants are fine
		}
		if field.Desc.Kind() != protoreflect.MessageKind {
			return false // complex scalar type
		}
		// Message variant: check all its fields are scalar/enum with <=3 children
		msg := field.Message
		count := 0
		for _, subField := range msg.Fields {
			if isFieldDeprecated(subField) {
				continue
			}
			if subField.Desc.Kind() == protoreflect.MessageKind {
				return false // nested message → too complex
			}
			if subField.Desc.IsMap() || subField.Desc.IsList() {
				return false
			}
			count++
		}
		if count > 3 {
			return false
		}
	}
	return true
}

// expandNestedOneof expands a nested oneof within a variant into flags.
func expandNestedOneof(oneof *protogen.Oneof, v *inlinedVariant, variantName string, msg *protogen.Message) {
	neo := nestedOneofMeta{
		GoName: oneof.GoName,
	}

	for _, field := range oneof.Fields {
		oneofVariantName := toKebab(string(field.Desc.Name()))
		neo.Variants = append(neo.Variants, oneofVariantName)

		parentGoName := string(msg.GoIdent.GoName)
		wrapperType := parentGoName + "_" + field.GoName

		if field.Desc.Kind() == protoreflect.BoolKind || field.Desc.Kind() == protoreflect.StringKind {
			// Scalar oneof variant
			flagName := oneofVariantName
			goVarSuffix := "Neo_" + field.GoName
			fd := buildVariantFlag(field, flagName, goVarSuffix, variantName, "", "")
			if fd != nil {
				fd.Required = false
				fd.NestedOneofGroup = oneof.GoName
				fd.NestedOneofVariant = field.GoName
				fd.NestedOneofWrapperType = wrapperType
				fd.NestedOneofMsgType = ""
				v.Flags = append(v.Flags, *fd)
			}
			continue
		}

		if field.Desc.Kind() != protoreflect.MessageKind {
			continue
		}

		subMsg := field.Message
		for _, subField := range subMsg.Fields {
			if isFieldDeprecated(subField) {
				continue
			}
			if subField.Desc.Kind() == protoreflect.MessageKind {
				continue
			}

			childKebab := toKebab(string(subField.Desc.Name()))
			flagName := oneofVariantName + "-" + childKebab
			if oneofVariantName == childKebab {
				flagName = oneofVariantName
			}
			goVarSuffix := "Neo_" + field.GoName + "_" + subField.GoName

			fd := buildVariantFlag(subField, flagName, goVarSuffix, variantName, "", "")
			if fd != nil {
				fd.Required = false
				fd.NestedOneofGroup = oneof.GoName
				fd.NestedOneofVariant = field.GoName
				fd.NestedOneofWrapperType = wrapperType
				fd.NestedOneofMsgType = subMsg.GoIdent.GoName
				v.Flags = append(v.Flags, *fd)
			}
		}
	}

	v.NestedOneofs = append(v.NestedOneofs, neo)
}

// extractFilterConvenienceFlags detects repeated search filter fields with IDFilter oneof variants.
func extractFilterConvenienceFlags(msg *protogen.Message, goImportPath string) ([]filterConvenienceFlag, []extraImport) {
	var flags []filterConvenienceFlag
	var imports []extraImport

	for _, field := range msg.Fields {
		if !field.Desc.IsList() || field.Desc.Kind() != protoreflect.MessageKind {
			continue
		}
		// Look for repeated *SearchFilter fields
		filterMsg := field.Message
		if !strings.HasSuffix(filterMsg.GoIdent.GoName, "SearchFilter") {
			continue
		}

		// Check for a "filter" oneof
		for _, filterField := range filterMsg.Fields {
			if filterField.Oneof == nil || filterField.Desc.HasOptionalKeyword() {
				continue
			}

			// Look for IDFilter variants
			if filterField.Desc.Kind() != protoreflect.MessageKind {
				continue
			}

			// Check if this is an IDFilter type (zitadel.filter.v2.IDFilter)
			idFilterFullName := filterField.Desc.Message().FullName()
			if idFilterFullName != "zitadel.filter.v2.IDFilter" {
				continue
			}

			oneofFieldName := filterField.GoName
			flagName := toKebab(string(filterField.Desc.Name()))
			flagName = strings.TrimSuffix(flagName, "-filter")

			// Build the wrapper type name: FilterMsg_OneofFieldName
			wrapperType := filterMsg.GoIdent.GoName + "_" + oneofFieldName

			// IDFilter import
			idFilterMsg := filterField.Message
			idFilterImportPath := string(idFilterMsg.GoIdent.GoImportPath)
			idFilterType := idFilterMsg.GoIdent.GoName
			idFilterAlias := importAlias(idFilterImportPath)

			if idFilterImportPath != goImportPath {
				found := false
				for _, ei := range imports {
					if ei.Path == idFilterImportPath {
						found = true
						break
					}
				}
				if !found {
					imports = append(imports, extraImport{Alias: idFilterAlias, Path: idFilterImportPath})
				}
			} else {
				idFilterAlias = ""
			}

			help := fmt.Sprintf("Filter by %s", strings.ReplaceAll(flagName, "-", " "))

			cf := filterConvenienceFlag{
				FlagName:            flagName,
				Help:                help,
				OneofFieldName:      oneofFieldName,
				FilterListField:     field.GoName,
				FilterMsgType:       filterMsg.GoIdent.GoName,
				WrapperType:         wrapperType,
				OneofGoFieldName:    filterField.Oneof.GoName,
				IDFilterImportAlias: idFilterAlias,
				IDFilterType:        idFilterType,
			}
			flags = append(flags, cf)
		}
	}

	return flags, imports
}
