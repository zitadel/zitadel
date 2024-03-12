package postgres

import (
	"context"
	"database/sql"
	"slices"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

func (s *Storage) Query(ctx context.Context, query *eventstore.Query) (err error) {
	var stmt database.Statement
	writeQuery(&stmt, query)

	if query.Tx() != nil {
		return executeQuery(ctx, query.Tx(), &stmt, query.Reducer())
	}

	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return err
	}
	defer func() {
		err = database.CloseTx(tx, err)
	}()

	return executeQuery(ctx, tx, &stmt, query.Reducer())
}

func executeQuery(ctx context.Context, tx *sql.Tx, stmt *database.Statement, reducer eventstore.Reducer) error {
	rows, err := tx.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return err
	}

	return database.MapRowsToObject(rows, func(scan func(dest ...any) error) error {
		e := &event{
			aggregate: &eventstore.Aggregate{},
		}

		var payload sql.Null[[]byte]

		err := scan(
			&e.createdAt,
			&e.typ,
			&e.sequence,
			&e.position.Position,
			&e.position.InPositionOrder,
			&payload,
			&e.creator,
			&e.aggregate.Owner,
			&e.aggregate.Instance,
			&e.aggregate.Type,
			&e.aggregate.ID,
			&e.revision,
		)
		if err != nil {
			return err
		}
		e.payload = payload.V

		return reducer.Reduce(e)
	})
}

var (
	selectColumns       = `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision`
	instancePlaceholder = database.Placeholder("instance_id")
)

func writeQuery(stmt *database.Statement, query *eventstore.Query) {
	stmt.WriteString(selectColumns)
	stmt.SetNamedArg(instancePlaceholder, query.Instance())

	stmt.WriteString(" FROM (")
	writeFilters(stmt, query.Filters())
	stmt.WriteRune(')')
	writePagination(stmt, query.Pagination())
}

var from = " FROM eventstore.events2"

func writeFilters(stmt *database.Statement, filters []*eventstore.Filter) {
	if len(filters) == 0 {
		logging.Fatal("query does not contain filters")
	}

	for i, filter := range filters {
		if i > 0 {
			stmt.WriteString(" UNION ALL ")
		}
		stmt.WriteRune('(')
		stmt.WriteString(selectColumns)
		stmt.WriteString(from)

		writeFilter(stmt, filter)

		stmt.WriteString(")")
	}
}

func writeFilter(stmt *database.Statement, filter *eventstore.Filter) {
	stmt.WriteString(" WHERE ")
	database.NewTextEqual(instancePlaceholder).Write(stmt, "instance_id")

	writeAggregateFilters(stmt, filter.AggregateFilters())
	writePagination(stmt, filter.Pagination())
}

func writePagination(stmt *database.Statement, pagination *eventstore.Pagination) {
	writePosition(stmt, pagination.Position())
	writeOrdering(stmt, pagination.Desc())
	if pagination.Pagination() != nil {
		pagination.Pagination().Write(stmt)
	}
}

func writePosition(stmt *database.Statement, position *eventstore.GlobalPosition) {
	if position == nil {
		return
	}

	stmt.WriteString(" AND ")

	if position.InPositionOrder == 0 {
		database.NewNumberGreater(position.Position).Write(stmt, "position")
		return
	}

	stmt.WriteString("((")
	database.NewNumberEquals(position.Position).Write(stmt, "position")
	stmt.WriteString(" AND ")
	database.NewNumberGreater(position.InPositionOrder).Write(stmt, "in_tx_order")
	stmt.WriteRune(')')
	stmt.WriteString(" OR ")
	database.NewNumberGreater(position.Position).Write(stmt, "position")
	stmt.WriteRune(')')
}

func writeAggregateFilters(stmt *database.Statement, filters []*eventstore.AggregateFilter) {
	if len(filters) == 0 {
		return
	}

	stmt.WriteString(" AND ")
	if len(filters) > 1 {
		stmt.WriteRune('(')
	}
	for i, filter := range filters {
		if i > 0 {
			stmt.WriteString(" OR ")
		}
		writeAggregateFilter(stmt, filter)
	}
	if len(filters) > 1 {
		stmt.WriteRune(')')
	}
}

func writeAggregateFilter(stmt *database.Statement, filter *eventstore.AggregateFilter) {
	conditions := definedConditions([]*condition{
		{column: "aggregate_type", condition: filter.Type()},
		{column: "aggregate_id", condition: filter.ID()},
	})

	if len(conditions) > 1 || len(filter.Events()) > 0 {
		stmt.WriteRune('(')
	}

	writeConditions(
		stmt,
		conditions,
		" AND ",
	)
	writeEventFilters(stmt, filter.Events())

	if len(conditions) > 1 || len(filter.Events()) > 0 {
		stmt.WriteRune(')')
	}
}

func writeEventFilters(stmt *database.Statement, filters []*eventstore.EventFilter) {
	if len(filters) == 0 {
		return
	}

	stmt.WriteString(" AND ")
	if len(filters) > 1 {
		stmt.WriteRune('(')
	}

	for i, filter := range filters {
		if i > 0 {
			stmt.WriteString(" OR ")
		}
		writeEventFilter(stmt, filter)
	}

	if len(filters) > 1 {
		stmt.WriteRune(')')
	}
}

func writeEventFilter(stmt *database.Statement, filter *eventstore.EventFilter) {
	conditions := definedConditions([]*condition{
		{column: "event_type", condition: filter.Type()},
		{column: "created_at", condition: filter.CreatedAt()},
		{column: "sequence", condition: filter.Sequence()},
		{column: "revision", condition: filter.Revision()},
		{column: "creator", condition: filter.Creator()},
	})

	if len(conditions) > 1 {
		stmt.WriteRune('(')
	}

	writeConditions(
		stmt,
		conditions,
		" AND ",
	)

	if len(conditions) > 1 {
		stmt.WriteRune(')')
	}
}

type condition struct {
	column    string
	condition database.Condition
}

func writeConditions(stmt *database.Statement, conditions []*condition, sep string) {
	var i int
	for _, cond := range conditions {
		if i > 0 {
			stmt.WriteString(sep)
		}
		cond.condition.Write(stmt, cond.column)
		i++
	}
}

func definedConditions(conditions []*condition) []*condition {
	return slices.DeleteFunc(conditions, func(cond *condition) bool {
		return cond.condition == nil
	})
}

func writeOrdering(stmt *database.Statement, descending bool) {
	stmt.Builder.WriteString(" ORDER BY position")
	if descending {
		stmt.Builder.WriteString(" DESC")
	}

	stmt.Builder.WriteString(", in_tx_order")
	if descending {
		stmt.Builder.WriteString(" DESC")
	}
}
