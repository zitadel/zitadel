package schemas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRelaxedBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{name: "true", input: true, expected: "true"},
		{name: "false", expected: "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := NewRelaxedBool(tt.input)
			bytes, err := json.Marshal(value)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(bytes))
		})
	}
}

func TestRelaxedBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		wantErr  bool
	}{
		{name: "valid true", input: "true", expected: true},
		{name: "valid false", input: "false"},
		{name: "quoted true", input: `"true"`, expected: true},
		{name: "quoted pascal case true", input: `"True"`, expected: true},
		{name: "quoted upper case true", input: `"TRUE"`, expected: true},
		{name: "quoted false", input: `"false"`},
		{name: "quoted pascal case false", input: `"False"`},
		{name: "quoted upper case false", input: `"FALSE"`},
		{name: "invalid value", input: "invalid", wantErr: true},
		{name: "number value", input: "1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := new(RelaxedBool)
			err := json.Unmarshal([]byte(tt.input), &value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, value.Bool())
			}
		})
	}
}
