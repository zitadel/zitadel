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
			name:   "change to NULL",
			change: NewChangeToNull(NewColumn("table", "column")),
			want: want{
				stmt: "column = NULL",
				args: nil,
			},
		},
		{
			name: "multiple changes",
			change: NewChanges(
				NewChange(NewColumn("table", "column1"), "value1"),
				NewChangePtr[int](NewColumn("table", "column2"), nil),
				NewChange(NewColumn("table", "column3"), 123),
				NewChangeToNull(NewColumn("table", "column4")),
			),
			want: want{
				stmt: "column1 = $1, column2 = NULL, column3 = $2, column4 = NULL",
				args: []any{"value1", 123},
			},
		},
		{
			name:   "increment change",
			change: NewIncrementColumnChange(NewColumn("table", "counter"), NewColumn("table", "counter")),
			want: want{
				stmt: "counter = table.counter + 1",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var builder StatementBuilder
			err := test.change.Write(&builder)
			require.NoError(t, err)
			assert.Equal(t, test.want.stmt, builder.String())
			require.Len(t, builder.Args(), len(test.want.args))
			for i, arg := range test.want.args {
				assert.Equal(t, arg, builder.Args()[i])
			}
		})
	}
}

func TestChangeToStatement(t *testing.T) {
	type want struct {
		stmt string
		args []any
	}
	for _, test := range []struct {
		name           string
		prefillBuilder func(builder *StatementBuilder)
		change         *changeToStatement
		want           want
	}{
		{
			name: "change to statement",
			change: NewChangeToStatement(NewColumn("table", "column"), func(builder *StatementBuilder) {
				builder.WriteString("SELECT 1")
			}).(*changeToStatement),
			want: want{
				stmt: "column = (SELECT 1)",
				args: nil,
			},
		},
		{
			name: "change to statement with args",
			change: NewChangeToStatement(NewColumn("table", "column"), func(builder *StatementBuilder) {
				builder.WriteString("SELECT ")
				builder.WriteArg(42)
			}).(*changeToStatement),
			want: want{
				stmt: "column = (SELECT $1)",
				args: []any{42},
			},
		},
		{
			name: "change to statement with existing builder args",
			prefillBuilder: func(builder *StatementBuilder) {
				builder.WriteString("UPDATE table SET ")
				assert.NoError(t, NewChanges(
					NewChange(NewColumn("table", "field1"), "asdf"),
					NewChangeToNull(NewColumn("table", "field2")),
					NewChangeToColumn(NewColumn("table", "field3"), NewColumn("table", "field4")),
				).Write(builder))
				builder.WriteString(", ")
			},
			change: NewChangeToStatement(NewColumn("table", "column"), func(builder *StatementBuilder) {
				builder.WriteString("SELECT ")
				builder.WriteArg(42)
				builder.WriteString(" FROM other_table WHERE ")
				NewBooleanCondition(NewColumn("table", "id"), true).Write(builder)
			}).(*changeToStatement),
			want: want{
				stmt: "UPDATE table SET field1 = $1, field2 = NULL, field3 = table.field4, column = (SELECT $2 FROM other_table WHERE table.id = $3)",
				args: []any{"asdf", 42, true},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var builder StatementBuilder
			if test.prefillBuilder != nil {
				test.prefillBuilder(&builder)
			}
			assert.NoError(t, test.change.Write(&builder))
			assert.Equal(t, test.want.stmt, builder.String())
			assert.Equal(t, builder.Args(), test.want.args)
		})
	}
}

func TestCTEChange(t *testing.T) {
	type want struct {
		stmt string
		args []any
	}
	for _, test := range []struct {
		name           string
		prefillBuilder func(builder *StatementBuilder)
		afterCTE       func(builder *StatementBuilder)
		change         *cteChange
		want           want
	}{
		{
			name: "only CTE change",
			change: NewCTEChange(
				func(builder *StatementBuilder) {
					builder.WriteString("SELECT 1 AS value")
				},
				nil,
			).(*cteChange),
			want: want{
				stmt: "WITH cte AS (SELECT 1 AS value) ",
				args: nil,
			},
		},
		{
			name: "with existing CTE",
			prefillBuilder: func(builder *StatementBuilder) {
				builder.WriteString("existing_cte AS (SELECT ")
				builder.WriteArg(42)
				builder.WriteString(" AS reason), ")
			},
			change: NewCTEChange(
				func(builder *StatementBuilder) {
					builder.WriteString("SELECT ")
					Columns{
						NewColumn("table", "column1"),
						NewColumn("table", "column2"),
					}.WriteQualified(builder)
					builder.WriteString(", ")
					builder.WriteArgs("asdf", 123, false, NullInstruction)
				},
				nil,
			).(*cteChange),
			want: want{
				stmt: "WITH existing_cte AS (SELECT $1 AS reason), cte AS (SELECT table.column1, table.column2, $2, $3, $4, NULL) ",
				args: []any{42, "asdf", 123, false},
			},
		},
		{
			name: "CTE change with after CTE statements",
			change: NewCTEChange(
				func(builder *StatementBuilder) {
					builder.WriteString("SELECT ")
					builder.WriteArg(1)
					builder.WriteString(" AS value")
				},
				nil,
			).(*cteChange),
			afterCTE: func(builder *StatementBuilder) {
				builder.WriteString("SELECT * FROM cte ")
			},
			want: want{
				stmt: "WITH cte AS (SELECT $1 AS value) SELECT * FROM cte ",
				args: []any{1},
			},
		},
		{
			name: "CTE change with column change after CTE",
			change: NewCTEChange(
				func(builder *StatementBuilder) {
					builder.WriteString("SELECT ")
					builder.WriteArg(1)
					builder.WriteString(" AS value")
				},
				func(name string) Change {
					return NewChangeToStatement(NewColumn("table", "field"), func(builder *StatementBuilder) {
						builder.WriteString("SELECT * FROM ")
						builder.WriteString(name)
						builder.WriteString(" WHERE ")
						NewColumnCondition(NewColumn(name, "value"), NewColumn("table", "value")).Write(builder)
					})
				},
			).(*cteChange),
			afterCTE: func(builder *StatementBuilder) {
				builder.WriteString("UPDATE test SET ")
			},
			want: want{
				stmt: "WITH cte AS (SELECT $1 AS value) UPDATE test SET field = (SELECT * FROM cte WHERE cte.value = table.value)",
				args: []any{1},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var builder StatementBuilder
			builder.WriteString("WITH ")
			if test.prefillBuilder != nil {
				test.prefillBuilder(&builder)
			}

			builder.WriteString("cte AS (")
			test.change.SetName("cte")
			test.change.WriteCTE(&builder)
			builder.WriteString(") ")

			if test.afterCTE != nil {
				test.afterCTE(&builder)
			}
			assert.NoError(t, test.change.Write(&builder))

			assert.Equal(t, test.want.stmt, builder.String())
			assert.Equal(t, builder.Args(), test.want.args)
		})
	}
}
