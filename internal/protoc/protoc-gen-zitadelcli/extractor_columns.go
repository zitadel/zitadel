package main

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func extractResponseColumns(msg *protogen.Message, verb string) (isList bool, listFieldGoName string, listFieldProtoName string, responseUnwrapField string, cols []columnDef) {
	// For List* methods, look for the repeated result field.
	for _, field := range msg.Fields {
		if field.Desc.IsList() && field.Desc.Kind() == protoreflect.MessageKind {
			isList = true
			listFieldGoName = field.GoName
			listFieldProtoName = string(field.Desc.Name())
			cols = extractMessageColumns(field.Message)
			return
		}
	}

	// For non-list methods, extract columns from top-level response fields.
	cols = extractTopLevelColumns(msg)
	if len(cols) > 0 {
		return false, "", "", "", cols
	}

	// If no scalar columns at top level, unwrap the first non-details message field
	// (e.g. GetUserByIDResponse.user, GetProjectResponse.project).
	for _, field := range msg.Fields {
		if field.Desc.Kind() == protoreflect.MessageKind && !field.Desc.IsList() && string(field.Desc.Name()) != "details" {
			responseUnwrapField = string(field.Desc.Name())
			innerCols := extractMessageColumns(field.Message)
			prefix := "Get" + field.GoName + "()."
			for i := range innerCols {
				innerCols[i].GoAccessor = prefix + innerCols[i].GoAccessor
				// FieldPath is already relative to the inner message; no need to prefix.
			}
			return false, "", "", responseUnwrapField, innerCols
		}
	}

	return false, "", "", "", cols
}

func extractMessageColumns(msg *protogen.Message) []columnDef {
	// "Primary ID" header: the resource's own ID field, e.g. "PROJECT ID" for Project,
	// "USER ID" for User. We rename it to plain "ID" since context (listing projects/users)
	// makes the resource type obvious.
	primaryIDHeader := strings.ToUpper(strings.ReplaceAll(toKebab(msg.GoIdent.GoName), "-", " ")) + " ID"

	var (
		primaryIdCol  *columnDef  // resource's own ID — shown as "ID"
		orgIdCols     []columnDef // ORGANIZATION ID (from details.resource_owner or direct field)
		otherIdCols   []columnDef // other foreign " ID" fields (user_id, granted_organization_id, …)
		semanticCols  []columnDef // names, states, booleans, enums
		timestampCols []columnDef // timestamps — shown last as metadata
	)

	for _, field := range msg.Fields {
		col := fieldToColumn(field)
		if col == nil {
			continue
		}
		switch {
		case col.Header == "ID":
			// Bare "id" field — the resource is its own top-level entity (e.g., Organization, IDP).
			primaryIdCol = col
		case col.Header == "ORGANIZATION ID":
			// Must be checked before primaryIDHeader — for Organization, primaryIDHeader is also
			// "ORGANIZATION ID" and details.resource_owner must go to orgIdCols, not primaryIdCol.
			orgIdCols = append(orgIdCols, *col)
		case col.Header == primaryIDHeader:
			// e.g., "PROJECT ID" for Project, "USER ID" for User — rename to plain "ID".
			col.Header = "ID"
			primaryIdCol = col
		case strings.HasSuffix(col.Header, " ID"):
			otherIdCols = append(otherIdCols, *col)
		case col.IsTimestamp:
			timestampCols = append(timestampCols, *col)
		default:
			semanticCols = append(semanticCols, *col)
		}
	}

	// Detect oneof fields and add a TYPE column showing which variant is set.
	processed := map[string]bool{}
	for _, field := range msg.Fields {
		if field.Oneof == nil || field.Desc.HasOptionalKeyword() {
			continue
		}
		oneofName := field.Oneof.GoName
		if processed[oneofName] {
			continue
		}
		processed[oneofName] = true

		header := strings.ToUpper(toKebab(oneofName))
		header = strings.ReplaceAll(header, "-", " ")
		accessor := "Get" + oneofName + "()"

		var variants []oneofVariantColumn
		for _, of := range field.Oneof.Fields {
			wrapperType := msg.GoIdent.GoName + "_" + of.GoName
			displayName := toKebab(string(of.Desc.Name()))
			variants = append(variants, oneofVariantColumn{
				GoWrapperType: wrapperType,
				DisplayName:   displayName,
			})
		}
		semanticCols = append(semanticCols, columnDef{
			Header:        header,
			GoAccessor:    accessor,
			IsOneofType:   true,
			OneofVariants: variants,
		})
	}

	// Final column order:
	//   1. Primary ID ("ID") — the resource's own key
	//   2. ORGANIZATION ID — ownership context
	//   3. Other foreign IDs (user_id, granted_organization_id, …)
	//   4. Semantic fields (name, state, booleans, enums)
	//   5. Timestamps (creation_date, change_date, …) — metadata, last
	var cols []columnDef
	if primaryIdCol != nil {
		cols = append(cols, *primaryIdCol)
	}
	cols = append(cols, orgIdCols...)
	cols = append(cols, otherIdCols...)
	cols = append(cols, semanticCols...)
	cols = append(cols, timestampCols...)
	return cols
}

func extractTopLevelColumns(msg *protogen.Message) []columnDef {
	var cols []columnDef
	var detailsField *protogen.Field
	hasNestedMessage := false
	for _, field := range msg.Fields {
		if string(field.Desc.Name()) == "details" {
			detailsField = field
			continue // try other fields first
		}
		col := fieldToColumn(field)
		if col != nil {
			cols = append(cols, *col)
		} else if field.Desc.Kind() == protoreflect.MessageKind && !field.Desc.IsList() {
			// fieldToColumn couldn't handle this nested message; let extractResponseColumns unwrap it.
			hasNestedMessage = true
		}
	}
	// When the response has no scalar columns AND no unwrappable nested message field,
	// show change_date from Details as a confirmation timestamp. This covers simple
	// "action" responses (e.g. delete, verify) that only return Details.
	// If there IS a nested message (e.g. GetUserByIDResponse.user), skip this so that
	// extractResponseColumns can unwrap it into full column definitions instead.
	if len(cols) == 0 && !hasNestedMessage && detailsField != nil && detailsField.Message != nil {
		for _, sub := range detailsField.Message.Fields {
			if string(sub.Desc.Name()) == "change_date" {
				cols = append(cols, columnDef{
					Header:      "CHANGE DATE",
					GoAccessor:  "GetDetails().GetChangeDate()",
					IsTimestamp: true,
				})
				break
			}
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
		return &columnDef{Header: header, GoAccessor: accessor, FieldPath: string(desc.Name())}
	case protoreflect.EnumKind:
		return &columnDef{Header: header, GoAccessor: accessor + ".String()", FieldPath: string(desc.Name()), IsEnum: true}
	case protoreflect.MessageKind:
		// Check for well-known types.
		fullName := desc.Message().FullName()
		if fullName == "google.protobuf.Timestamp" {
			return &columnDef{Header: header, GoAccessor: accessor, FieldPath: string(desc.Name()), IsTimestamp: true}
		}
		// For nested Details, extract resource_owner.
		if string(desc.Name()) == "details" {
			for _, subField := range field.Message.Fields {
				if string(subField.Desc.Name()) == "resource_owner" {
					return &columnDef{
						Header:     "ORGANIZATION ID",
						GoAccessor: accessor + ".GetResourceOwner()",
						FieldPath:  "details.resource_owner",
					}
				}
			}
			return nil
		}
		return nil
	default:
		return nil
	}
}
