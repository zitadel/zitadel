package filter

import (
	"reflect"
	"strings"
	"sync"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// jsonFieldCache Cache storing JSON tag to field mappings for each reflect.Type
// keep this in memory forever, as the fields / json tags are constant at runtime and never change.
// reflect.Type => structFieldCache
var jsonFieldCache sync.Map

type structFieldCache map[string]reflect.StructField

type AttributeResolver struct {
	schema schemas.ScimSchemaType
}

func (c structFieldCache) get(name string) (reflect.StructField, error) {
	field, ok := c[strings.ToLower(name)]
	if !ok {
		return reflect.StructField{}, zerrors.ThrowInvalidArgumentf(nil, "SCIM-attr12", "SCIM Attribute not found %s", name)
	}

	return field, nil
}

func (c structFieldCache) set(field reflect.StructField) {
	jsonTag := field.Tag.Get("json")

	// Skip fields explicitly excluded
	if jsonTag == "-" {
		return
	}

	// use field name as default
	fieldName := field.Name
	if jsonTag != "" {
		// strip other options such as omitempty
		fieldName = strings.Split(jsonTag, ",")[0]
	}

	c[strings.ToLower(fieldName)] = field
}

func newAttributeResolver(schema schemas.ScimSchemaType) *AttributeResolver {
	return &AttributeResolver{
		schema: schema,
	}
}

func (r *AttributeResolver) resolveAttrPath(item reflect.Value, attrPath *AttrPath) ([]string, reflect.Value, error) {
	if err := attrPath.validateSchema(r.schema); err != nil {
		return nil, reflect.Value{}, err
	}

	segments := attrPath.Segments()
	for _, segment := range segments {
		var err error
		item, err = r.resolveField(item, segment)
		if err != nil {
			return nil, reflect.Value{}, err
		}
	}

	return segments, item, nil
}

func (r *AttributeResolver) resolveField(item reflect.Value, fieldName string) (reflect.Value, error) {
	if item.Kind() == reflect.Ptr {
		item = item.Elem()
	}

	fields, err := r.getOrBuildFieldMap(item.Type())
	if err != nil {
		return reflect.Value{}, err
	}

	field, err := fields.get(fieldName)
	if err != nil {
		return reflect.Value{}, err
	}

	resolvedField := item.FieldByName(field.Name)
	if !resolvedField.IsValid() {
		return reflect.Value{}, zerrors.ThrowInvalidArgumentf(nil, "SCIM-attr13", "SCIM Attribute not found %s", fieldName)
	}

	return resolvedField, nil
}

func (r *AttributeResolver) getOrBuildFieldMap(t reflect.Type) (structFieldCache, error) {
	if cached, ok := jsonFieldCache.Load(t); ok {
		return cached.(structFieldCache), nil
	}

	// cache miss, build json name field map
	fieldMap := make(structFieldCache, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fieldMap.set(t.Field(i))
	}

	// Cache the result for future use.
	jsonFieldCache.Store(t, fieldMap)
	return fieldMap, nil
}
