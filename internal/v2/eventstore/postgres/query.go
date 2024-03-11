package postgres

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var (
	selectColumns = `SELECT created_at, event_type, "sequence", "position", in_tx_order, payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision`
	from          = " FROM eventstore.events2"
)

func (s *Storage) Query(ctx context.Context, instance string, reducer eventstore.Reducer, filters ...*eventstore.Filter) (err error) {
	var stmt database.Statement
	writeFilters(&stmt, instance, filters)

	// if f.Tx != nil {
	// 	return query(ctx, f.Tx, &stmt, reducer)
	// }

	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return err
	}
	defer func() {
		err = database.CloseTx(tx, err)
	}()

	return query(ctx, tx, &stmt, reducer)
}

func query(ctx context.Context, tx *sql.Tx, stmt *database.Statement, reducer eventstore.Reducer) error {
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

func writeFilters(stmt *database.Statement, instance string, filters []*eventstore.Filter) {
	for i, filter := range filters {
		if i > 0 {
			stmt.WriteString(" UNION ALL ")
		}
		writeFilter(stmt, instance, filter)
	}

}

func writeFilter(stmt *database.Statement, instance string, filter *eventstore.Filter) {
	stmt.WriteString(selectColumns)
	stmt.WriteString(" FROM (")

	stmt.WriteString(selectColumns)
	stmt.WriteString(from)

	stmt.WriteString(" WHERE ")
	// TODO: implement WriteNamed
	database.NewTextEqual(instance).WriteNamed(stmt, "instance_id")

	writeAggregateFilters(stmt, filter.AggregateFilters())

	if filter.Pagination() != nil && filter.Pagination().Position() != nil {
		stmt.WriteString(" AND ")
		filter.Pagination().Position().Position().Write(stmt, "position")
		if txOrder := filter.Pagination().Position().InTxOrder(); txOrder != nil {
			stmt.WriteString(" AND ")
			txOrder.Write(stmt, "in_tx_order")
		}
	}

	writeOrdering(stmt, filter.Desc())
	if filter.Pagination() != nil && filter.Pagination().Pagination() != nil {
		filter.Pagination().Pagination().Write(stmt)
	}
	stmt.WriteString(") res")
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
	conditions := map[string]database.Condition{
		"aggregate_type": filter.Type(),
		"aggregate_id":   filter.ID(),
	}

	deleteUnset(conditions)
	if len(conditions) > 1 || len(filter.Events()) > 0 {
		stmt.WriteRune('(')
	}
	writeConditions(stmt, conditions, " AND ")
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
	conditions := map[string]database.Condition{
		"event_type": filter.Type(),
		"created_at": filter.CreatedAt(),
		"sequence":   filter.Sequence(),
		"revision":   filter.Revision(),
		"creator":    filter.Creator(),
	}
	deleteUnset(conditions)
	if len(conditions) > 1 {
		stmt.WriteRune('(')
	}
	writeConditions(stmt, conditions, " AND ")
	if len(conditions) > 1 {
		stmt.WriteRune(')')
	}
}

func writeConditions(stmt *database.Statement, conditions map[string]database.Condition, sep string) {
	var i int
	for columnName, condition := range conditions {
		if i > 0 {
			stmt.WriteString(sep)
		}
		condition.Write(stmt, columnName)
		i++
	}
}

func deleteUnset(conditions map[string]database.Condition) {
	for columnName, condition := range conditions {
		if condition != nil {
			continue
		}
		delete(conditions, columnName)
	}
}

// func writeFilterClauses(stmt *database.Statement, filter *eventstore.Filter) {
// 	filter.Instances.Write(stmt, "instance_id")
// 	if filter.Pagination.Position != nil {
// 		stmt.WriteString(" AND ")
// 		filter.Position.Write(stmt, "position")
// 	}

// 	writePagination(stmt, filter.Pagination.Pagination)
// }

// func writeAggregatesFilter(stmt *database.Statement, filter *eventstore.Filter) {

// }

// func writeEventsFilter(stmt *database.Statement, filter *eventstore.Filter) {

// }

// func writeEventFilter(stmt *database.Statement, filter *eventstore.EventFilter) {
// 	if filter.Type() != nil {
// 		filter.Type().Write(stmt, "event_type")
// 	}
// 	if filter.CreatedAt() != nil {
// 		filter.CreatedAt().Write(stmt, "created_at")
// 	}
// 	if position := filter.Position(); position != nil {
// 		if position.Position() != nil {
// 			position.Position().Write(stmt, "position")
// 		}
// 		if position.InTxOrder() != nil {
// 			position.InTxOrder().Write(stmt, "in_tx_order")
// 		}
// 	}

// }

// func writePagination(stmt *database.Statement, pagination *database.Pagination) {
// 	if pagination == nil {
// 		return
// 	}

// 	pagination.Write(stmt)
// }

// func writeEventFilters(stmt *database.Statement, filters []*eventstore.EventFilter) {
// 	if len(filters) == 0 {
// 		return
// 	}

// 	stmt.Builder.WriteString(" AND ")

// 	if len(filters) > 1 {
// 		stmt.Builder.WriteRune('(')
// 	}
// 	for queryIdx, eventFilter := range filters {
// 		writeEventFilter(stmt, eventFilter)

// 		if queryIdx < len(filters)-1 {
// 			stmt.Builder.WriteString(" OR ")
// 		}
// 	}
// 	if len(filters) > 1 {
// 		stmt.Builder.WriteRune(')')
// 	}
// }

// func writeEventFilter(stmt *database.Statement, filter *eventstore.EventFilter) {
// 	filters := filter.Filters()

// 	if len(filters) > 1 {
// 		stmt.Builder.WriteRune('(')
// 	}

// 	var mustAddAnd bool

// 	if filter.AggregateTypes != nil {
// 		filter.AggregateTypes.Write(stmt, "aggregate_type")
// 		mustAddAnd = true
// 	}
// 	if filter.AggregateIDs != nil {
// 		if mustAddAnd {
// 			stmt.WriteString(" AND ")
// 		}
// 		filter.AggregateIDs.Write(stmt, "aggregate_id")
// 		mustAddAnd = true
// 	}
// 	if filter.EventTypes != nil {
// 		if mustAddAnd {
// 			stmt.WriteString(" AND ")
// 		}
// 		filter.EventTypes.Write(stmt, "event_type")
// 		mustAddAnd = true
// 	}
// 	if filter.Sequence != nil {
// 		if mustAddAnd {
// 			stmt.WriteString(" AND ")
// 		}
// 		filter.Sequence.Write(stmt, "sequence")
// 	}

// 	if len(filters) > 1 {
// 		stmt.Builder.WriteRune(')')
// 	}
// }

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

// func writeEventFilter(stmt *database.Statement, filter database.Filter) {
// 	var columnName string
// 	switch filter.(type) {
// 	case *eventstore.AggregateTypesFilter:
// 		columnName = "aggregate_type"
// 	case *eventstore.EventTypesFilter:
// 		columnName = "event_type"
// 	case *eventstore.AggregateIDsFilter:
// 		columnName = "aggregate_id"
// 	case *eventstore.SequenceFilter[eventstore.SequenceEqualsType],
// 		*eventstore.SequenceFilter[eventstore.SequenceAtLeastType],
// 		*eventstore.SequenceFilter[eventstore.SequenceGreaterType],
// 		*eventstore.SequenceFilter[eventstore.SequenceAtMostType],
// 		*eventstore.SequenceFilter[eventstore.SequenceLessType],
// 		*eventstore.SequenceFilter[eventstore.SequenceBetweenType]:
// 		columnName = "sequence"
// 	default:
// 		logging.WithFields("ext_type", fmt.Sprintf("%T", filter)).
// 			Panic("event filter not implemented")
// 	}

// 	filter.Write(stmt, columnName)
// }
