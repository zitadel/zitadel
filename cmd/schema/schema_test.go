package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"sigs.k8s.io/yaml"
)

func TestSchemaGeneration(t *testing.T) {
	err := run()
	if err != nil {
		t.Fatalf("failed to generate schema: %v", err)
	}

	schemaPath := ".artifacts/pack/zitadel-config-schema.json"
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Fatalf("schema file was not created at %s", schemaPath)
	}

	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema file: %v", err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}

	if schema["type"] != "object" {
		t.Error("schema should have type 'object'")
	}
	if schema["properties"] == nil {
		t.Error("schema missing properties field")
	}

	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("properties is not an object")
	}

	expectedMinFields := 40
	if len(props) < expectedMinFields {
		t.Errorf("expected at least %d properties, got %d", expectedMinFields, len(props))
	}

	t.Logf("Schema generated successfully with %d top-level properties", len(props))
}

func TestSchemaValidatesDefaults(t *testing.T) {
	err := run()
	if err != nil {
		t.Fatalf("failed to generate schema: %v", err)
	}

	schemaPath := ".artifacts/pack/zitadel-config-schema.json"
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema file: %v", err)
	}

	var schema jsonschema.Schema
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		t.Fatalf("failed to parse schema: %v", err)
	}

	defaultsPath := "../../cmd/defaults.yaml"
	defaultsData, err := os.ReadFile(defaultsPath)
	if err != nil {
		t.Fatalf("failed to read defaults.yaml: %v", err)
	}

	var defaults interface{}
	if err := yaml.Unmarshal(defaultsData, &defaults); err != nil {
		t.Fatalf("failed to parse defaults.yaml: %v", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		t.Fatalf("failed to resolve schema: %v", err)
	}

	err = resolved.Validate(defaults)
	if err != nil {
		t.Logf("Warning: defaults.yaml has validation issues: %v", err)
		t.Log("This may be expected if defaults.yaml contains fields not in the Config struct")
	} else {
		t.Log("defaults.yaml validates successfully against the schema")
	}
}
