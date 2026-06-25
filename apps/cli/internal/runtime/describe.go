package runtime

import (
	"encoding/json"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// JSONSchema represents a subset of JSON Schema Draft-07 sufficient for
// describing protobuf request/response messages to AI agents.
type JSONSchema struct {
	Type        string                `json:"type,omitempty"`
	Description string                `json:"description,omitempty"`
	Properties  map[string]*JSONSchema `json:"properties,omitempty"`
	Required    []string              `json:"required,omitempty"`
	Items       *JSONSchema           `json:"items,omitempty"`
	Enum        []string              `json:"enum,omitempty"`
	Format      string                `json:"format,omitempty"`
	OneOf       []OneOfEntry          `json:"oneOf,omitempty"`
}

// OneOfEntry is a single variant in a oneOf constraint.
type OneOfEntry struct {
	Description string     `json:"description,omitempty"`
	Properties  map[string]*JSONSchema `json:"properties,omitempty"`
}

// DescribeOutput is the envelope for describing a command to AI agents.
type DescribeOutput struct {
	Group       string      `json:"group"`
	Verb        string      `json:"verb"`
	Method      string      `json:"method"`
	Short       string      `json:"short"`
	Long        string      `json:"long,omitempty"`
	Example     string      `json:"example,omitempty"`
	Deprecated  bool        `json:"deprecated,omitempty"`
	RequestSchema *JSONSchema `json:"request_schema"`
	ResponseSchema *JSONSchema `json:"response_schema,omitempty"`
}

// DescribeGroupOutput lists commands in a group.
type DescribeGroupOutput struct {
	Group    string   `json:"group"`
	Commands []string `json:"commands"`
}

// DescribeAllOutput lists all groups and global flags.
type DescribeAllOutput struct {
	GlobalFlags []GlobalFlagInfo `json:"global_flags"`
	Groups      map[string][]string `json:"groups"`
}

// GlobalFlagInfo describes a global flag for schema introspection.
type GlobalFlagInfo struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Help       string   `json:"help"`
	EnumValues []string `json:"enum_values,omitempty"`
}

// BuildDescribeOutput generates a full describe response for a CommandSpec,
// including real JSON schemas derived from protoreflect.
func BuildDescribeOutput(spec CommandSpec) (*DescribeOutput, error) {
	_, reqDesc, err := resolveMethod(spec.FullMethodName)
	if err != nil {
		// Fallback: return metadata without schemas.
		return &DescribeOutput{
			Group:  spec.Group,
			Verb:   spec.Verb,
			Method: spec.FullMethodName,
			Short:  spec.Short,
			Long:   spec.Long,
		}, nil
	}

	methodDesc, _, _ := resolveMethod(spec.FullMethodName)
	reqSchema := messageToSchema(reqDesc, 0)
	respSchema := messageToSchema(methodDesc.Output(), 0)

	return &DescribeOutput{
		Group:          spec.Group,
		Verb:           spec.Verb,
		Method:         spec.FullMethodName,
		Short:          spec.Short,
		Long:           spec.Long,
		Example:        spec.Example,
		Deprecated:     spec.Deprecated,
		RequestSchema:  reqSchema,
		ResponseSchema: respSchema,
	}, nil
}

// BuildDescribeAll generates the top-level describe output from all registered specs.
func BuildDescribeAll() *DescribeAllOutput {
	allSpecs := AllSpecs()
	groups := make(map[string][]string)
	for _, s := range allSpecs {
		cmd := s.Group + " " + s.Verb
		groups[s.Group] = appendUnique(groups[s.Group], cmd)
	}
	return &DescribeAllOutput{
		GlobalFlags: globalFlagsInfo(),
		Groups:      groups,
	}
}

// BuildDescribeGroup generates a group-level describe output.
func BuildDescribeGroup(groupName string) (*DescribeGroupOutput, []CommandSpec) {
	allSpecs := AllSpecs()
	var commands []string
	var groupSpecs []CommandSpec
	for _, s := range allSpecs {
		if s.Group == groupName {
			commands = append(commands, s.Verb)
			groupSpecs = append(groupSpecs, s)
		}
	}
	return &DescribeGroupOutput{
		Group:    groupName,
		Commands: commands,
	}, groupSpecs
}

// MarshalJSON helper for describe outputs.
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// messageToSchema converts a protobuf message descriptor to a JSON Schema.
func messageToSchema(msgDesc protoreflect.MessageDescriptor, depth int) *JSONSchema {
	if depth > 4 {
		return &JSONSchema{Type: "object"}
	}

	schema := &JSONSchema{
		Type:       "object",
		Properties: make(map[string]*JSONSchema),
	}

	// Handle oneofs.
	processedOneofs := make(map[string]bool)

	for i := 0; i < msgDesc.Fields().Len(); i++ {
		fd := msgDesc.Fields().Get(i)
		fieldName := string(fd.JSONName())

		// Track oneofs.
		if oneofDesc := fd.ContainingOneof(); oneofDesc != nil && !fd.HasOptionalKeyword() {
			oneofName := string(oneofDesc.Name())
			if !processedOneofs[oneofName] {
				processedOneofs[oneofName] = true
				// Add a pseudo-property showing the oneof constraint.
				var entries []OneOfEntry
				for j := 0; j < oneofDesc.Fields().Len(); j++ {
					of := oneofDesc.Fields().Get(j)
					entry := OneOfEntry{
						Description: string(of.Name()),
						Properties:  map[string]*JSONSchema{string(of.JSONName()): fieldToSchema(of, depth)},
					}
					entries = append(entries, entry)
				}
				schema.Properties["_oneof_"+oneofName] = &JSONSchema{
					Description: "Choose one of: " + strings.Join(oneofFieldNames(oneofDesc), ", "),
					OneOf:       entries,
				}
			}
			continue
		}

		schema.Properties[fieldName] = fieldToSchema(fd, depth)
	}

	return schema
}

// fieldToSchema converts a single proto field to a JSON Schema property.
func fieldToSchema(fd protoreflect.FieldDescriptor, depth int) *JSONSchema {
	if fd.IsList() {
		return &JSONSchema{
			Type:  "array",
			Items: scalarSchema(fd, depth),
		}
	}
	if fd.IsMap() {
		return &JSONSchema{
			Type: "object",
			Description: "map field",
		}
	}
	return scalarSchema(fd, depth)
}

func scalarSchema(fd protoreflect.FieldDescriptor, depth int) *JSONSchema {
	switch fd.Kind() {
	case protoreflect.StringKind:
		return &JSONSchema{Type: "string"}
	case protoreflect.BoolKind:
		return &JSONSchema{Type: "boolean"}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Uint32Kind,
		protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind:
		return &JSONSchema{Type: "integer"}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return &JSONSchema{Type: "number"}
	case protoreflect.BytesKind:
		return &JSONSchema{Type: "string", Format: "byte"}
	case protoreflect.EnumKind:
		var values []string
		for j := 0; j < fd.Enum().Values().Len(); j++ {
			values = append(values, string(fd.Enum().Values().Get(j).Name()))
		}
		return &JSONSchema{Type: "string", Enum: values}
	case protoreflect.MessageKind:
		fullName := fd.Message().FullName()
		switch fullName {
		case "google.protobuf.Timestamp":
			return &JSONSchema{Type: "string", Format: "date-time"}
		case "google.protobuf.Duration":
			return &JSONSchema{Type: "string", Format: "duration", Description: "e.g. '3600s'"}
		case "google.protobuf.Struct":
			return &JSONSchema{Type: "object"}
		default:
			return messageToSchema(fd.Message(), depth+1)
		}
	default:
		return &JSONSchema{Type: "string"}
	}
}

func oneofFieldNames(od protoreflect.OneofDescriptor) []string {
	names := make([]string, od.Fields().Len())
	for i := 0; i < od.Fields().Len(); i++ {
		names[i] = string(od.Fields().Get(i).Name())
	}
	return names
}

func globalFlagsInfo() []GlobalFlagInfo {
	return []GlobalFlagInfo{
		{Name: "from-json", Type: "bool", Help: "Read request body as JSON from stdin. When set, required flags are not enforced."},
		{Name: "request-json", Type: "string", Help: "Provide request body as inline JSON string. Alternative to --from-json with stdin."},
		{Name: "dry-run", Type: "bool", Help: "Print the request as JSON without calling the API."},
		{Name: "output", Type: "string", Help: "Output format: table or json (auto-detected from TTY).", EnumValues: []string{"table", "json"}},
		{Name: "context", Type: "string", Help: "Override the active context."},
	}
}

func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}
