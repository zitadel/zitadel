package schema

import "github.com/santhosh-tekuri/jsonschema/v5"

type Schema struct {
	schema *jsonschema.Schema
}

func New(url, schema string) (*Schema, error) {
	compiled, err := jsonschema.CompileString(url, schema)
	if err != nil {
		return nil, err
	}
	return &Schema{schema: compiled}, nil
}

func (s *Schema) Schema() *jsonschema.Schema { return s.schema }

func (s *Schema) Validate(rawJSON interface{}) error {
	return s.schema.Validate(rawJSON)
}
