package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnexpectedQueryTypeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		typeVal  any
		expected string
	}{
		{
			name:     "string type",
			id:       "id1",
			typeVal:  "test",
			expected: "ID=id1 Message=unexpected query type 'string'",
		},
		{
			name:     "int type",
			id:       "id2",
			typeVal:  42,
			expected: "ID=id2 Message=unexpected query type 'int'",
		},
		{
			name:     "struct type",
			id:       "id3",
			typeVal:  struct{}{},
			expected: "ID=id3 Message=unexpected query type 'struct {}'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewUnexpectedQueryTypeError(tt.id, tt.typeVal)
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}
