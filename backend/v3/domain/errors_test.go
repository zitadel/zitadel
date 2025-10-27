package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnexpectedQueryTypeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		typeVal  any
		expected string
	}{
		{
			name:     "string type",
			typeVal:  "test",
			expected: "Message=unexpected query type 'string'",
		},
		{
			name:     "int type",
			typeVal:  42,
			expected: "Message=unexpected query type 'int'",
		},
		{
			name:     "struct type",
			typeVal:  struct{}{},
			expected: "Message=unexpected query type 'struct {}'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUnexpectedQueryTypeError(tt.typeVal)
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}
