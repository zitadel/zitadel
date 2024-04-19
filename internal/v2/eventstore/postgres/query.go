package postgres

import (
	"context"
	"database/sql"
	"slices"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

func (s *Storage) Query(ctx context.Context, query *eventstore.Query) (eventCount int, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var stmt database.Statement
	writeQuery(&stmt, query)

	if query.Tx() != nil {
		return executeQuery(ctx, query.Tx(), &stmt, query)
	}

	return executeQuery(ctx, s.client.DB, &stmt, query)
}

func executeQuery(ctx context.Context, tx database.Querier, stmt *database.Statement, reducer eventstore.Reducer) (eventCount int, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	//nolint:rowserrcheck
	// rows is checked by database.MapRowsToObject
	rows, err := tx.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return 0, err
	}

	err = database.MapRowsToObject(rows, func(scan func(dest ...any) error) error {
		e := new(eventstore.Event[eventstore.StoragePayload])

		var payload sql.Null[[]byte]

		err := scan(
			&e.CreatedAt,
			&e.Type,
			&e.Sequence,
			&e.Position.Position,
			&e.Position.InPositionOrder,
			&payload,
			&e.Creator,
			&e.Aggregate.Owner,
			&e.Aggregate.Instance,
			&e.Aggregate.Type,
			&e.Aggregate.ID,
			&e.Revision,
		)
		if err != nil {
			return err
		}
		e.Payload = unmarshalPayload(payload.V)
		eventCount++

		return reducer.Reduce(e)
	})

	return eventCount, err
}

var (
	selectColumns = `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision`
	// TODO: condition must know if it's args are named parameters or not
	// instancePlaceholder = database.Placeholder("@instance_id")
)

func writeQuery(stmt *database.Statement, query *eventstore.Query) {
	stmt.WriteString(selectColumns)
	// stmt.SetNamedArg(instancePlaceholder, query.Instance())

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
	filter.Parent().Instance().Write(stmt, "instance_id")

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

func writePosition(stmt *database.Statement, position *eventstore.PositionCondition) {
	if position == nil {
		return
	}

	max := position.Max()
	min := position.Min()

	stmt.WriteString(" AND ")

	if max != nil {
		if max.InPositionOrder > 0 {
			stmt.WriteString("((")
			database.NewNumberEquals(max.Position).Write(stmt, "position")
			stmt.WriteString(" AND ")
			database.NewNumberLess(max.InPositionOrder).Write(stmt, "in_tx_order")
			stmt.WriteRune(')')
			stmt.WriteString(" OR ")
		}
		database.NewNumberLess(max.Position).Write(stmt, "position")
		if max.InPositionOrder > 0 {
			stmt.WriteRune(')')
		}
	}

	if max != nil && min != nil {
		stmt.WriteString(" AND ")
	}

	if min != nil {
		if min.InPositionOrder > 0 {
			stmt.WriteString("((")
			database.NewNumberEquals(min.Position).Write(stmt, "position")
			stmt.WriteString(" AND ")
			database.NewNumberGreater(min.InPositionOrder).Write(stmt, "in_tx_order")
			stmt.WriteRune(')')
			stmt.WriteString(" OR ")
		}
		database.NewNumberGreater(min.Position).Write(stmt, "position")
		if min.InPositionOrder > 0 {
			stmt.WriteRune(')')
		}
	}
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
		{column: "aggregate_id", condition: filter.IDs()},
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
		{column: "event_type", condition: filter.Types()},
		{column: "created_at", condition: filter.CreatedAt()},
		{column: "sequence", condition: filter.Sequence()},
		{column: "revision", condition: filter.Revision()},
		{column: "creator", condition: filter.Creators()},
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
	stmt.WriteString(" ORDER BY position")
	if descending {
		stmt.WriteString(" DESC")
	}

	stmt.WriteString(", in_tx_order")
	if descending {
		stmt.WriteString(" DESC")
	}
}
