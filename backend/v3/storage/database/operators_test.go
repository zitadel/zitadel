package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteTextOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		col      Column
		op       TextOperation
		value    string
		expected string
		args     []any
	}{
		{
			name:     "Equal",
			col:      NewColumn("test", "col"),
			op:       TextOperationEqual,
			value:    "value",
			expected: "test.col = $1",
			args:     []any{"value"},
		},
		{
			name:     "NotEqual",
			col:      NewColumn("test", "col"),
			op:       TextOperationNotEqual,
			value:    "value",
			expected: "test.col <> $1",
			args:     []any{"value"},
		},
		{
			name:     "EqualIgnoreCase",
			col:      NewColumn("test", "col"),
			op:       TextOperationEqualIgnoreCase,
			value:    "value",
			expected: "LOWER(test.col) LIKE LOWER($1)",
			args:     []any{"value"},
		},
		{
			name:     "NotEqualIgnoreCase",
			col:      NewColumn("test", "col"),
			op:       TextOperationNotEqualIgnoreCase,
			value:    "value",
			expected: "LOWER(test.col) NOT LIKE LOWER($1)",
			args:     []any{"value"},
		},
		{
			name:     "StartsWith",
			col:      NewColumn("test", "col"),
			op:       TextOperationStartsWith,
			value:    "value",
			expected: "test.col LIKE $1 || '%'",
			args:     []any{"value"},
		},
		{
			name:     "StartsWithIgnoreCase",
			col:      NewColumn("test", "col"),
			op:       TextOperationStartsWithIgnoreCase,
			value:    "value",
			expected: "LOWER(test.col) LIKE LOWER($1) || '%'",
			args:     []any{"value"},
		},
		{
			name:     "Contains",
			col:      NewColumn("test", "col"),
			op:       TextOperationContains,
			value:    "value",
			expected: "test.col LIKE '%' || $1 || '%'",
			args:     []any{"value"},
		},
		{
			name:     "ContainsIgnoreCase",
			col:      NewColumn("test", "col"),
			op:       TextOperationContainsWithIgnoreCase,
			value:    "value",
			expected: "test.col ILIKE '%' || $1 || '%'",
			args:     []any{"value"},
		},
		{
			name:     "EndsWith",
			col:      NewColumn("test", "col"),
			op:       TextOperationEndsWith,
			value:    "value",
			expected: "test.col LIKE '%' || $1",
			args:     []any{"value"},
		},
		{
			name:     "EndsWithIgnoreCase",
			col:      NewColumn("test", "col"),
			op:       TextOperationEndsWithIgnoreCase,
			value:    "value",
			expected: "LOWER(test.col) LIKE '%' || LOWER($1)",
			args:     []any{"value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			builder := &StatementBuilder{}
			writeTextOperation(builder, tt.col, tt.op, tt.value)

			assert.Equal(t, tt.expected, builder.String())
			assert.Equal(t, tt.args, builder.Args())
		})
	}

	t.Run("panic on invalid operation", func(t *testing.T) {
		t.Parallel()
		defer func() {
			require.NotNil(t, recover())
		}()
		builder := &StatementBuilder{}
		writeTextOperation(builder, NewColumn("test", "col"), TextOperation(99), "value")
	})
}
