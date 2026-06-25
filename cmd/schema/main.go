// Package main generates a JSON schema from the ZITADEL configuration struct.
// The generated schema is written to .artifacts/pack/zitadel-config-schema.json
// and is included in release artifacts.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"

	"github.com/zitadel/zitadel/cmd/start"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run generates the JSON schema and writes it to the output file.
func run() error {
	opts := &jsonschema.ForOptions{
		TypeSchemas:        customTypeSchemas(),
		IgnoreInvalidTypes: true,
	}

	schema, err := jsonschema.ForType(reflect.TypeOf(start.Config{}), opts)
	if err != nil {
		return fmt.Errorf("failed to infer schema: %w", err)
	}

	schema.ID = "https://zitadel.com/schemas/config.json"
	schema.Title = "ZITADEL Configuration"
	schema.Description = "Configuration schema for ZITADEL server"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	outputPath := ".artifacts/pack/zitadel-config-schema.json"
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	fmt.Printf("Schema written to %s\n", outputPath)
	return nil
}
