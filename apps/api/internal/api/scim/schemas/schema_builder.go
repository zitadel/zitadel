package schemas

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	"golang.org/x/text/language"
)

type SchemaBuilderArgs struct {
	ID           ScimSchemaType
	Name         ScimResourceTypeSingular
	EndpointName ScimResourceTypePlural
	Description  string
	Resource     any
}

type fieldSchemaInfo struct {
	Ignore    bool
	Required  bool
	CaseExact bool
	Unique    bool
}

var (
	timeType            = reflect.TypeOf(time.Time{})
	languageTagType     = reflect.TypeOf(language.Tag{})
	httpURLType         = reflect.TypeOf(HttpURL{})
	writeOnlyStringType = reflect.TypeOf(WriteOnlyString(""))
)

func BuildSchema(args SchemaBuilderArgs) *ResourceSchema {
	return &ResourceSchema{
		Resource: &Resource{
			Schemas: []ScimSchemaType{IdSchema},
			ID:      string(args.ID),
			Meta: &ResourceMeta{
				ResourceType: SchemaResourceType,
			},
		},
		ID:          args.ID,
		Name:        args.Name,
		PluralName:  args.EndpointName,
		Description: args.Description,
		Attributes:  buildSchemaAttributes(reflect.TypeOf(args.Resource)),
	}
}

func buildSchemaAttributes(fieldType reflect.Type) []*SchemaAttribute {
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	if fieldType.Kind() != reflect.Struct {
		return nil
	}

	attributes := make([]*SchemaAttribute, 0, fieldType.NumField())
	for i := 0; i < fieldType.NumField(); i++ {
		field := fieldType.Field(i)
		attribute := buildAttribute(field)

		if attribute != nil {
			attributes = append(attributes, attribute)
		}
	}

	return attributes
}

func buildAttribute(field reflect.StructField) *SchemaAttribute {
	info := getFieldSchemaInfo(field)
	if info.Ignore {
		return nil
	}

	fieldType := getFieldType(field)
	attribute := &SchemaAttribute{
		Name:        getFieldJsonName(field),
		Description: "For details see RFC7643",
		Type:        getFieldAttributeType(fieldType),
		MultiValued: isFieldMultiValued(field),
		Required:    info.Required,
		CaseExact:   info.CaseExact,
		Mutability:  SchemaAttributeMutabilityReadWrite,
		Returned:    SchemaAttributeReturnedAlways,
		Uniqueness:  SchemaAttributeUniquenessNone,
	}

	if attribute.Type == SchemaAttributeTypeComplex {
		attribute.SubAttributes = buildSchemaAttributes(fieldType)
	}

	if fieldType == writeOnlyStringType {
		attribute.Returned = SchemaAttributeReturnedNever
		attribute.Mutability = SchemaAttributeMutabilityWriteOnly
	}

	if info.Unique {
		attribute.Uniqueness = SchemaAttributeUniquenessServer
	}

	return attribute
}

func isFieldMultiValued(field reflect.StructField) bool {
	if field.Type.Kind() != reflect.Ptr {
		return field.Type.Kind() == reflect.Slice
	}

	return field.Type.Elem().Kind() == reflect.Slice
}

func getFieldAttributeType(fieldType reflect.Type) SchemaAttributeType {
	switch fieldType.Kind() { //nolint:exhaustive
	case reflect.String:
		return SchemaAttributeTypeString
	case reflect.Bool:
		return SchemaAttributeTypeBoolean
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return SchemaAttributeTypeInteger
	case reflect.Float32, reflect.Float64:
		return SchemaAttributeTypeDecimal
	case reflect.Struct:
		switch fieldType {
		case timeType:
			return SchemaAttributeTypeDateTime
		case writeOnlyStringType, languageTagType, httpURLType:
			return SchemaAttributeTypeString
		default:
			return SchemaAttributeTypeComplex
		}
	default:
		panic(fmt.Sprintf("unsupported field type: %v", fieldType.Kind()))
	}
}

func getFieldType(field reflect.StructField) reflect.Type {
	fieldType := field.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
		fieldType = fieldType.Elem()

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
	}
	return fieldType
}

func getFieldSchemaInfo(field reflect.StructField) *fieldSchemaInfo {
	tag := field.Tag.Get("scim")
	tagOptions := strings.Split(tag, ",")
	return &fieldSchemaInfo{
		Ignore:    slices.Contains(tagOptions, "ignoreInSchema"),
		Required:  slices.Contains(tagOptions, "required"),
		CaseExact: !slices.Contains(tagOptions, "caseInsensitive"),
		Unique:    slices.Contains(tagOptions, "unique"),
	}
}

func getFieldJsonName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")

	// Skip fields explicitly excluded
	if jsonTag == "-" {
		return ""
	}

	// use field name as default
	if jsonTag == "" {
		return field.Name
	}

	// strip other options such as omitempty
	return strings.Split(jsonTag, ",")[0]
}
