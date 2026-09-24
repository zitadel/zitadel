package main

import (
	"reflect"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
)

// customTypeSchemas returns custom JSON schema definitions for Go types
// that need special handling beyond the default inference.
func customTypeSchemas() map[reflect.Type]*jsonschema.Schema {
	return map[reflect.Type]*jsonschema.Schema{
		reflect.TypeOf(time.Duration(0)): {
			Type:        "string",
			Title:       "Duration",
			Description: "A duration string like '1h', '30m', '5s', '100ms'",
			Pattern:     `^-?([0-9]+(\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$`,
			Examples:    []any{"1h", "30m", "5s", "100ms", "1h30m"},
		},
	}
}
