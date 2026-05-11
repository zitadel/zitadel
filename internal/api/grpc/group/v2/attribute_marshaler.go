package group

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// attributeMarshaler wraps protojson marshaler to flatten attribute maps in JSON output.
// Transforms:
type attributeMarshaler struct {
	protojson.MarshalOptions
}

func (m *attributeMarshaler) Marshal(msg proto.Message) ([]byte, error) {
	// Marshal with default protojson
	data, err := m.MarshalOptions.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Post-process to flatten attributes
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return data, nil // Return original if not an object
	}

	flattenAttributes(obj)

	return json.Marshal(obj)
}

// flattenAttributes recursively finds and flattens attribute maps
func flattenAttributes(obj map[string]interface{}) {
	for key, val := range obj {
		switch v := val.(type) {
		case map[string]interface{}:
			// Check if this looks like an attribute map
			if key == "attributes" && isAttributeMap(v) {
				obj[key] = flattenAttributeMap(v)
			} else {
				flattenAttributes(v)
			}
		case []interface{}:
			// Recurse into arrays
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					flattenAttributes(m)
				}
			}
		}
	}
}

// isAttributeMap checks if map contains AttributeValue wrappers
func isAttributeMap(m map[string]interface{}) bool {
	for _, val := range m {
		if attrVal, ok := val.(map[string]interface{}); ok {
			if _, hasSingle := attrVal["single"]; hasSingle {
				return true
			}
			if _, hasMultiple := attrVal["multiple"]; hasMultiple {
				return true
			}
		}
	}
	return false
}

// flattenAttributeMap transforms attribute values from wrapped to flat format
func flattenAttributeMap(attrs map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(attrs))
	for key, val := range attrs {
		attrVal, ok := val.(map[string]interface{})
		if !ok {
			result[key] = val
			continue
		}

		// Check for single value
		if single, ok := attrVal["single"].(string); ok {
			result[key] = single
			continue
		}

		// Check for multiple values
		if multiple, ok := attrVal["multiple"].(map[string]interface{}); ok {
			if values, ok := multiple["values"].([]interface{}); ok {
				result[key] = values
				continue
			}
		}

		// Fallback - keep original
		result[key] = val
	}
	return result
}
