package v3

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Query[O object] interface {
	Where(condition Condition)
	Join(tables ...Table)
	Limit(limit uint32)
	Offset(offset uint32)
	OrderBy(columns ...Column)

	Result(ctx context.Context, client database.Querier) (*O, error)
	Results(ctx context.Context, client database.Querier) ([]O, error)

	fmt.Stringer
	statementBuilder
}

type query[O object] struct {
	*statement[O]
	joins   []join
	limit   uint32
	offset  uint32
	orderBy []Column
}

func NewQuery[O object](table Table) Query[O] {
	return &query[O]{
		statement: newStatement[O](table),
	}
}

// Result implements [Query].
func (q *query[O]) Result(ctx context.Context, client database.Querier) (*O, error) {
	var object O
	row := client.QueryRow(ctx, q.String(), q.args...)
	if err := object.Scan(row); err != nil {
		return nil, err
	}
	return &object, nil
}

// Results implements [Query].
func (q *query[O]) Results(ctx context.Context, client database.Querier) ([]O, error) {
	var objects []O
	rows, err := client.Query(ctx, q.String(), q.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var object O
		if err := object.Scan(rows); err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}

	return objects, rows.Err()
}

// Join implements [Query].
func (q *query[O]) Join(tables ...Table) {
	for _, tbl := range tables {
		cols := q.tbl.(*table).possibleJoins(tbl)
		if len(cols) == 0 {
			panic(fmt.Sprintf("table %q does not have any possible joins with table %q", q.tbl.Name(), tbl.Name()))
		}

		q.joins = append(q.joins, join{
			table:      tbl,
			conditions: make([]joinCondition, 0, len(cols)),
		})

		for colName, col := range cols {
			q.joins[len(q.joins)-1].conditions = append(q.joins[len(q.joins)-1].conditions, joinCondition{
				left:  q.tbl.(*table).columns[colName],
				right: col,
			})
		}
	}
}

func (q *query[O]) Limit(limit uint32) {
	q.limit = limit
}

func (q *query[O]) Offset(offset uint32) {
	q.offset = offset
}

func (q *query[O]) OrderBy(columns ...Column) {
	for _, allowedColumn := range q.columns {
		for _, column := range columns {
			if allowedColumn.Name() == column.Name() {
				q.orderBy = append(q.orderBy, column)
			}
		}
	}
}

// String implements [fmt.Stringer] and [Query].
func (q *query[O]) String() string {
	q.writeSelectColumns()
	q.writeFrom()
	q.writeJoins()
	q.writeCondition()
	q.writeOrderBy()
	q.writeLimit()
	q.writeOffset()
	q.writeGroupBy()
	return q.builder.String()
}

func (q *query[O]) writeSelectColumns() {
	q.builder.WriteString("SELECT ")
	for i, column := range q.columns {
		if i > 0 {
			q.builder.WriteString(", ")
		}
		q.builder.WriteString(q.tbl.Alias())
		q.builder.WriteRune('.')
		q.builder.WriteString(column.Name())
	}
}

func (q *query[O]) writeJoins() {
	for _, join := range q.joins {
		q.builder.WriteString(" JOIN ")
		q.builder.WriteString(join.table.Schema())
		q.builder.WriteRune('.')
		q.builder.WriteString(join.table.Name())
		if join.table.Alias() != "" {
			q.builder.WriteString(" AS ")
			q.builder.WriteString(join.table.Alias())
		}

		q.builder.WriteString(" ON ")
		for i, condition := range join.conditions {
			if i > 0 {
				q.builder.WriteString(" AND ")
			}
			q.builder.WriteString(condition.left.Name())
			q.builder.WriteString(" = ")
			q.builder.WriteString(condition.right.Name())
		}
	}
}

func (q *query[O]) writeOrderBy() {
	if len(q.orderBy) == 0 {
		return
	}

	q.builder.WriteString(" ORDER BY ")
	for i, order := range q.orderBy {
		if i > 0 {
			q.builder.WriteString(", ")
		}
		order.Write(q)
	}
}

func (q *query[O]) writeLimit() {
	if q.limit == 0 {
		return
	}
	q.builder.WriteString(" LIMIT ")
	q.builder.WriteString(q.appendArg(q.limit))
}

func (q *query[O]) writeOffset() {
	if q.offset == 0 {
		return
	}
	q.builder.WriteString(" OFFSET ")
	q.builder.WriteString(q.appendArg(q.offset))
}

func (q *query[O]) writeGroupBy() {
	q.builder.WriteString(" GROUP BY ")
}
