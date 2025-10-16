package database

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeWrite(t *testing.T) {
	type want struct {
		stmt string
		args []any
	}
	for _, test := range []struct {
		name   string
		change Change
		want   want
	}{
		{
			name:   "change",
			change: NewChange(NewColumn("table", "column"), "value"),
			want: want{
				stmt: "column = $1",
				args: []any{"value"},
			},
		},
		{
			name:   "change ptr to null",
			change: NewChangePtr[int](NewColumn("table", "column"), nil),
			want: want{
				stmt: "column = NULL",
				args: nil,
			},
		},
		{
			name:   "change ptr to value",
			change: NewChangePtr(NewColumn("table", "column"), gu.Ptr(42)),
			want: want{
				stmt: "column = $1",
				args: []any{42},
			},
		},
		{
			name: "multiple changes",
			change: NewChanges(
				NewChange(NewColumn("table", "column1"), "value1"),
				NewChangePtr[int](NewColumn("table", "column2"), nil),
				NewChange(NewColumn("table", "column3"), 123),
			),
			want: want{
				stmt: "column1 = $1, column2 = NULL, column3 = $2",
				args: []any{"value1", 123},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var builder StatementBuilder
			test.change.Write(&builder)
			assert.Equal(t, test.want.stmt, builder.String())
			require.Len(t, builder.Args(), len(test.want.args))
			for i, arg := range test.want.args {
				assert.Equal(t, arg, builder.Args()[i])
			}
		})
	}
}
