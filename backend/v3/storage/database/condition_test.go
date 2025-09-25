package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	type want struct {
		stmt string
		args []any
	}
	for _, test := range []struct {
		name string
		cond Condition
		want want
	}{
		{
			name: "no and condition",
			cond: And(),
			want: want{
				stmt: "",
				args: nil,
			},
		},
		{
			name: "one and condition",
			cond: And(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
			),
			want: want{
				stmt: "table.column1 = other_table.column2",
				args: nil,
			},
		},
		{
			name: "multiple and condition",
			cond: And(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
				NewColumnCondition(NewColumn("table", "column3"), NewColumn("other_table", "column4")),
			),
			want: want{
				stmt: "(table.column1 = other_table.column2 AND table.column3 = other_table.column4)",
				args: nil,
			},
		},
		{
			name: "no or condition",
			cond: Or(),
			want: want{
				stmt: "",
				args: nil,
			},
		},
		{
			name: "one or condition",
			cond: Or(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
			),
			want: want{
				stmt: "table.column1 = other_table.column2",
				args: nil,
			},
		},
		{
			name: "multiple or condition",
			cond: Or(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
				NewColumnCondition(NewColumn("table", "column3"), NewColumn("other_table", "column4")),
			),
			want: want{
				stmt: "(table.column1 = other_table.column2 OR table.column3 = other_table.column4)",
				args: nil,
			},
		},
		{
			name: "is null condition",
			cond: IsNull(NewColumn("table", "column1")),
			want: want{
				stmt: "table.column1 IS NULL",
				args: nil,
			},
		},
		{
			name: "is not null condition",
			cond: IsNotNull(NewColumn("table", "column1")),
			want: want{
				stmt: "table.column1 IS NOT NULL",
				args: nil,
			},
		},
		{
			name: "text condition",
			cond: NewTextCondition(NewColumn("table", "column1"), TextOperationEqual, "some text"),
			want: want{
				stmt: "table.column1 = $1",
				args: []any{"some text"},
			},
		},
		{
			name: "text ignore case condition",
			cond: NewTextIgnoreCaseCondition(NewColumn("table", "column1"), TextOperationNotEqual, "some TEXT"),
			want: want{
				stmt: "LOWER(table.column1) <> LOWER($1)",
				args: []any{"some TEXT"},
			},
		},
		{
			name: "number condition",
			cond: NewNumberCondition(NewColumn("table", "column1"), NumberOperationEqual, 42),
			want: want{
				stmt: "table.column1 = $1",
				args: []any{42},
			},
		},
		{
			name: "boolean condition",
			cond: NewBooleanCondition(NewColumn("table", "column1"), true),
			want: want{
				stmt: "table.column1 = $1",
				args: []any{true},
			},
		},
		{
			name: "bytes condition",
			cond: NewBytesCondition[[]byte](NewColumn("table", "column1"), BytesOperationEqual, []byte{0x01, 0x02, 0x03}),
			want: want{
				stmt: "table.column1 = $1",
				args: []any{[]byte{0x01, 0x02, 0x03}},
			},
		},
		{
			name: "column condition",
			cond: NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
			want: want{
				stmt: "table.column1 = other_table.column2",
				args: nil,
			},
		},
		{
			name: "exists condition",
			cond: Exists("table", And(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
				NewColumnCondition(NewColumn("table", "column3"), NewColumn("other_table", "column4")),
			)),
			want: want{
				stmt: " EXISTS (SELECT 1 FROM table WHERE (table.column1 = other_table.column2 AND table.column3 = other_table.column4))",
				args: nil,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var builder StatementBuilder
			test.cond.Write(&builder)
			assert.Equal(t, test.want.stmt, builder.String())
			require.Len(t, builder.Args(), len(test.want.args))
			for i, arg := range test.want.args {
				assert.Equal(t, arg, builder.Args()[i])
			}
		})
	}
}

func TestIsRestrictingColumn(t *testing.T) {
	for _, test := range []struct {
		name string
		col  Column
		cond Condition
		want bool
	}{
		{
			name: "and with restricting column",
			col:  NewColumn("table", "column1"),
			cond: And(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column3")),
			),
			want: true,
		},
		{
			name: "and without restricting column",
			col:  NewColumn("table", "column1"),
			cond: And(
				NewColumnCondition(NewColumn("table", "column2"), NewColumn("other_table", "column3")),
				IsNull(NewColumn("table", "column4")),
				IsNotNull(NewColumn("table", "column5")),
			),
			want: false,
		},
		{
			name: "or with restricting column",
			col:  NewColumn("table", "column1"),
			cond: Or(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column3")),
			),
			want: true,
		},
		{
			name: "or without restricting column",
			col:  NewColumn("table", "column1"),
			cond: Or(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column3")),
				IsNotNull(NewColumn("table", "column4")),
				IsNull(NewColumn("table", "column5")),
			),
			want: false,
		},
		{
			name: "is null never restricts",
			col:  NewColumn("table", "column1"),
			cond: IsNull(NewColumn("table", "column1")),
			want: false,
		},
		{
			name: "is not null never restricts",
			col:  NewColumn("table", "column1"),
			cond: IsNotNull(NewColumn("table", "column1")),
			want: false,
		},
		{
			name: "exists with restricting column",
			col:  NewColumn("table", "column1"),
			cond: Exists("table", And(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column2")),
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column3")),
			)),
			want: true,
		},
		{
			name: "exists without restricting column",
			col:  NewColumn("table", "column1"),
			cond: Exists("table", Or(
				NewColumnCondition(NewColumn("table", "column1"), NewColumn("other_table", "column3")),
				IsNotNull(NewColumn("table", "column4")),
				IsNull(NewColumn("table", "column5")),
			)),
			want: false,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			isRestricting := test.cond.IsRestrictingColumn(test.col)
			assert.Equal(t, test.want, isRestricting)
		})
	}
}
