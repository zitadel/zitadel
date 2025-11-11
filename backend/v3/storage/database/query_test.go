package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryOptions(t *testing.T) {
	type want struct {
		stmt string
		args []any
	}
	for _, test := range []struct {
		name    string
		options []QueryOption
		want    want
	}{
		{
			name: "no options",
			want: want{
				stmt: "",
				args: nil,
			},
		},
		{
			name: "limit option",
			options: []QueryOption{
				WithLimit(10),
			},
			want: want{
				stmt: " LIMIT $1",
				args: []any{uint32(10)},
			},
		},
		{
			name: "offset option",
			options: []QueryOption{
				WithOffset(5),
			},
			want: want{
				stmt: " OFFSET $1",
				args: []any{uint32(5)},
			},
		},
		{
			name: "order by asc option",
			options: []QueryOption{
				WithOrderByAscending(NewColumn("table", "column")),
			},
			want: want{
				stmt: " ORDER BY table.column",
				args: nil,
			},
		},
		{
			name: "order by desc option",
			options: []QueryOption{
				WithOrderByDescending(NewColumn("table", "column")),
			},
			want: want{
				stmt: " ORDER BY table.column DESC",
				args: nil,
			},
		},
		{
			name: "order by option",
			options: []QueryOption{
				WithOrderBy(OrderDirectionAsc, NewColumn("table", "column1"), NewColumn("table", "column2")),
			},
			want: want{
				stmt: " ORDER BY table.column1, table.column2",
				args: nil,
			},
		},
		{
			name: "condition option",
			options: []QueryOption{
				WithCondition(NewBooleanCondition(NewColumn("table", "column"), true)),
			},
			want: want{
				stmt: " WHERE table.column = $1",
				args: []any{true},
			},
		},
		{
			name: "group by option",
			options: []QueryOption{
				WithGroupBy(NewColumn("table", "column")),
			},
			want: want{
				stmt: " GROUP BY table.column",
				args: nil,
			},
		},
		{
			name: "left join option",
			options: []QueryOption{
				WithLeftJoin("other_table", NewColumnCondition(NewColumn("table", "id"), NewColumn("other_table", "table_id"))),
			},
			want: want{
				stmt: " LEFT JOIN other_table ON table.id = other_table.table_id",
				args: nil,
			},
		},
		{
			name: "permission check option",
			options: []QueryOption{
				WithPermissionCheck("permission"),
			},
			want: want{
				stmt: "",
				args: nil,
			},
		},
		{
			name: "with lock",
			options: []QueryOption{
				WithResultLock(),
			},
			want: want{
				stmt: " FOR UPDATE",
				args: nil,
			},
		},
		{
			name: "multiple options",
			options: []QueryOption{
				WithLeftJoin("other_table", NewColumnCondition(NewColumn("table", "id"), NewColumn("other_table", "table_id"))),
				WithCondition(NewNumberCondition(NewColumn("table", "column"), NumberOperationEqual, 123)),
				WithOrderByDescending(NewColumn("table", "column")),
				WithLimit(10),
				WithOffset(5),
				WithResultLock(),
			},
			want: want{
				stmt: " LEFT JOIN other_table ON table.id = other_table.table_id WHERE table.column = $1 ORDER BY table.column DESC LIMIT $2 OFFSET $3 FOR UPDATE",
				args: []any{123, uint32(10), uint32(5)},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var b StatementBuilder
			var opts QueryOpts
			for _, option := range test.options {
				option(&opts)
			}
			opts.Write(&b)
			assert.Equal(t, test.want.stmt, b.String())
			require.Len(t, b.Args(), len(test.want.args))
			for i := range test.want.args {
				assert.Equal(t, test.want.args[i], b.Args()[i])
			}
		})
	}

}
