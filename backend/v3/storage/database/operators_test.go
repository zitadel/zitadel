package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_writeOperation(t *testing.T) {
	type want struct {
		shouldPanic bool
		stmt        string
		args        []any
	}
	tests := []struct {
		name  string
		write func(builder *StatementBuilder)
		// col   Column
		// op    string
		// value any
		want want
	}{
		{
			name: "unsupported operation panics",
			write: func(builder *StatementBuilder) {
				writeOperation[string](builder, NewColumn("table", "column"), "", "value")
			},
			want: want{
				shouldPanic: true,
			},
		},
		{
			name: "unsupported value type panics",
			write: func(builder *StatementBuilder) {
				writeOperation[string](builder, NewColumn("table", "column"), " = ", struct{}{})
			},
			want: want{
				shouldPanic: true,
			},
		},
		{
			name: "unsupported wrapped value type panics",
			write: func(builder *StatementBuilder) {
				writeOperation[string](builder, NewColumn("table", "column"), " = ", SHA256Value(123))
			},
			want: want{
				shouldPanic: true,
			},
		},
		{
			name: "text equal",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, NewColumn("table", "column"), TextOperationEqual, "value")
			},
			want: want{
				stmt: "table.column = $1",
				args: []any{"value"},
			},
		},
		{
			name: "text not equal",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, NewColumn("table", "column"), TextOperationNotEqual, "value")
			},
			want: want{
				stmt: "table.column <> $1",
				args: []any{"value"},
			},
		},
		{
			name: "text starts with",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, NewColumn("table", "column"), TextOperationStartsWith, "value")
			},
			want: want{
				stmt: "table.column LIKE $1 || '%'",
				args: []any{"value"},
			},
		},
		{
			name: "text ends with",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, NewColumn("table", "column"), TextOperationEndsWith, "value")
			},
			want: want{
				stmt: "table.column LIKE '%' || $1",
				args: []any{"value"},
			},
		},
		{
			name: "text contains",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, NewColumn("table", "column"), TextOperationContains, "value")
			},
			want: want{
				stmt: "table.column LIKE '%' || $1 || '%'",
				args: []any{"value"},
			},
		},
		{
			name: "text equal with wrapped value",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, LowerColumn(NewColumn("table", "column")), TextOperationEqual, LowerValue("value"))
			},
			want: want{
				stmt: "LOWER(table.column) = LOWER($1)",
				args: []any{"value"},
			},
		},
		{
			name: "text not equal with wrapped value",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, LowerColumn(NewColumn("table", "column")), TextOperationNotEqual, LowerValue("value"))
			},
			want: want{
				stmt: "LOWER(table.column) <> LOWER($1)",
				args: []any{"value"},
			},
		},
		{
			name: "text starts with with wrapped value",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, LowerColumn(NewColumn("table", "column")), TextOperationStartsWith, LowerValue("value"))
			},
			want: want{
				stmt: "LOWER(table.column) LIKE LOWER($1) || '%'",
				args: []any{"value"},
			},
		},
		{
			name: "text ends with with wrapped value",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, LowerColumn(NewColumn("table", "column")), TextOperationEndsWith, LowerValue("value"))
			},
			want: want{
				stmt: "LOWER(table.column) LIKE '%' || LOWER($1)",
				args: []any{"value"},
			},
		},
		{
			name: "text contains with wrapped value",
			write: func(builder *StatementBuilder) {
				writeTextOperation[string](builder, LowerColumn(NewColumn("table", "column")), TextOperationContains, LowerValue("value"))
			},
			want: want{
				stmt: "LOWER(table.column) LIKE '%' || LOWER($1) || '%'",
				args: []any{"value"},
			},
		},

		{
			name: "number equal",
			write: func(builder *StatementBuilder) {
				writeNumberOperation[int](builder, NewColumn("table", "column"), NumberOperationEqual, 123)
			},
			want: want{
				stmt: "table.column = $1",
				args: []any{123},
			},
		},
		{
			name: "number not equal",
			write: func(builder *StatementBuilder) {
				writeNumberOperation[int](builder, NewColumn("table", "column"), NumberOperationNotEqual, 123)
			},
			want: want{
				stmt: "table.column <> $1",
				args: []any{123},
			},
		},
		{
			name: "number less than",
			write: func(builder *StatementBuilder) {
				writeNumberOperation[int](builder, NewColumn("table", "column"), NumberOperationLessThan, 123)
			},
			want: want{
				stmt: "table.column < $1",
				args: []any{123},
			},
		},
		{
			name: "number less than or equal",
			write: func(builder *StatementBuilder) {
				writeNumberOperation[int](builder, NewColumn("table", "column"), NumberOperationAtLeast, 123)
			},
			want: want{
				stmt: "table.column <= $1",
				args: []any{123},
			},
		},
		{
			name: "number greater than",
			write: func(builder *StatementBuilder) {
				writeNumberOperation[int](builder, NewColumn("table", "column"), NumberOperationGreaterThan, 123)
			},
			want: want{
				stmt: "table.column > $1",
				args: []any{123},
			},
		},
		{
			name: "number greater than or equal",
			write: func(builder *StatementBuilder) {
				writeNumberOperation[int](builder, NewColumn("table", "column"), NumberOperationAtMost, 123)
			},
			want: want{
				stmt: "table.column >= $1",
				args: []any{123},
			},
		},

		{
			name: "boolean is true",
			write: func(builder *StatementBuilder) {
				writeBooleanOperation[bool](builder, NewColumn("table", "column"), true)
			},
			want: want{
				stmt: "table.column = $1",
				args: []any{true},
			},
		},
		{
			name: "boolean is false",
			write: func(builder *StatementBuilder) {
				writeBooleanOperation[bool](builder, NewColumn("table", "column"), false)
			},
			want: want{
				stmt: "table.column = $1",
				args: []any{false},
			},
		},

		{
			name: "bytes equal",
			write: func(builder *StatementBuilder) {
				writeBytesOperation[[]byte](builder, NewColumn("table", "column"), BytesOperationEqual, []byte{0x01, 0x02, 0x03})
			},
			want: want{
				stmt: "table.column = $1",
				args: []any{[]byte{0x01, 0x02, 0x03}},
			},
		},
		{
			name: "bytes not equal",
			write: func(builder *StatementBuilder) {
				writeBytesOperation[[]byte](builder, NewColumn("table", "column"), BytesOperationNotEqual, []byte{0x01, 0x02, 0x03})
			},
			want: want{
				stmt: "table.column <> $1",
				args: []any{[]byte{0x01, 0x02, 0x03}},
			},
		},
		{
			name: "bytes equal with wrapped value",
			write: func(builder *StatementBuilder) {
				writeBytesOperation[[]byte](builder, SHA256Column(NewColumn("table", "column")), BytesOperationEqual, SHA256Value([]byte{0x01, 0x02, 0x03}))
			},
			want: want{
				stmt: "SHA256(table.column) = SHA256($1)",
				args: []any{[]byte{0x01, 0x02, 0x03}},
			},
		},
		{
			name: "bytes not equal with wrapped value",
			write: func(builder *StatementBuilder) {
				writeBytesOperation[[]byte](builder, SHA256Column(NewColumn("table", "column")), BytesOperationNotEqual, SHA256Value([]byte{0x01, 0x02, 0x03}))
			},
			want: want{
				stmt: "SHA256(table.column) <> SHA256($1)",
				args: []any{[]byte{0x01, 0x02, 0x03}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(t, tt.want.shouldPanic, r != nil)
			}()
			var builder StatementBuilder
			tt.write(&builder)

			assert.Equal(t, tt.want.stmt, builder.String())
			assert.Equal(t, tt.want.args, builder.Args())
		})
	}
}

func TestTextOperation_Value(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		op   TextOperation
		want TextOperation
	}{
		{
			name: "TextOperationStartsWith returns self",
			op:   TextOperationStartsWith,
			want: TextOperationStartsWith,
		},
		{
			name: "TextOperationContains returns self",
			op:   TextOperationContains,
			want: TextOperationContains,
		},
		{
			name: "TextOperationEndsWith returns self",
			op:   TextOperationEndsWith,
			want: TextOperationEndsWith,
		},
		{
			name: "TextOperationEqual returns self",
			op:   TextOperationEqual,
			want: TextOperationEqual,
		},
		{
			name: "TextOperationNotEqual returns self",
			op:   TextOperationNotEqual,
			want: TextOperationNotEqual,
		},
		{
			name: "TextOperationContainsIgnoreCase returns Contains",
			op:   TextOperationContainsIgnoreCase,
			want: TextOperationContains,
		},
		{
			name: "TextOperationEndsWithIgnoreCase returns EndsWith",
			op:   TextOperationEndsWithIgnoreCase,
			want: TextOperationEndsWith,
		},
		{
			name: "TextOperationStartsWithIgnoreCase returns StartsWith",
			op:   TextOperationStartsWithIgnoreCase,
			want: TextOperationStartsWith,
		},
		{
			name: "TextOperationEqualIgnoreCase returns Equal",
			op:   TextOperationEqualIgnoreCase,
			want: TextOperationEqual,
		},
		{
			name: "TextOperationNotEqualIgnoreCase returns NotEqual",
			op:   TextOperationNotEqualIgnoreCase,
			want: TextOperationNotEqual,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.op.Value()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTextOperation_IsIgnoreCaseOperation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		op   TextOperation
		want bool
	}{
		{
			name: "TextOperationContainsIgnoreCase returns true",
			op:   TextOperationContainsIgnoreCase,
			want: true,
		},
		{
			name: "TextOperationEndsWithIgnoreCase returns true",
			op:   TextOperationEndsWithIgnoreCase,
			want: true,
		},
		{
			name: "TextOperationStartsWithIgnoreCase returns true",
			op:   TextOperationStartsWithIgnoreCase,
			want: true,
		},
		{
			name: "TextOperationEqualIgnoreCase returns true",
			op:   TextOperationEqualIgnoreCase,
			want: true,
		},
		{
			name: "TextOperationNotEqualIgnoreCase returns true",
			op:   TextOperationNotEqualIgnoreCase,
			want: true,
		},
		{
			name: "TextOperationContains returns false",
			op:   TextOperationContains,
			want: false,
		},
		{
			name: "TextOperationEndsWith returns false",
			op:   TextOperationEndsWith,
			want: false,
		},
		{
			name: "TextOperationStartsWith returns false",
			op:   TextOperationStartsWith,
			want: false,
		},
		{
			name: "TextOperationEqual returns false",
			op:   TextOperationEqual,
			want: false,
		},
		{
			name: "TextOperationNotEqual returns false",
			op:   TextOperationNotEqual,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.op.IsIgnoreCaseOperation()
			assert.Equal(t, tt.want, got)
		})
	}
}
