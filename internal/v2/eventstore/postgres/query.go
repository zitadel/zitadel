package postgres

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/v2/database"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

var (
	selectEvents = `SELECT created_at, event_type, "sequence", "position", payload, creator, "owner", instance_id, aggregate_type, aggregate_id, revision FROM eventstore.events2 WHERE `
)

func (s *Storage) Query(ctx context.Context, f *eventstore.Filter, reducer eventstore.Reducer) (err error) {
	var stmt database.Statement
	filterQuery(&stmt, f)

	if f.Tx != nil {
		return query(ctx, f.Tx, &stmt, reducer)
	}

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
			&e.position,
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

func filterQuery(stmt *database.Statement, filter *eventstore.Filter) {
	stmt.WriteString(selectEvents)
	writeFilterClauses(stmt, filter)

	writeEventFilters(stmt, filter.EventFilters)

	writeOrdering(stmt, filter.Descending)
}

func writeFilterClauses(stmt *database.Statement, filter *eventstore.Filter) {
	filter.Instances.Write(stmt, "instance_id")
	if filter.Position != nil {
		stmt.WriteString(" AND ")
		filter.Position.Write(stmt, "position")
	}
}

func writeEventFilters(stmt *database.Statement, filters []*eventstore.EventFilter) {
	if len(filters) == 0 {
		return
	}

	stmt.Builder.WriteString(" AND ")

	if len(filters) > 1 {
		stmt.Builder.WriteRune('(')
	}
	for queryIdx, eventFilter := range filters {
		writeEventFilter(stmt, eventFilter)

		if queryIdx < len(filters)-1 {
			stmt.Builder.WriteString(" OR ")
		}
	}
	if len(filters) > 1 {
		stmt.Builder.WriteRune(')')
	}
}

func writeEventFilter(stmt *database.Statement, filter *eventstore.EventFilter) {
	filters := filter.Filters()

	if len(filters) > 1 {
		stmt.Builder.WriteRune('(')
	}

	var mustAddAnd bool

	if filter.AggregateTypes != nil {
		filter.AggregateTypes.Write(stmt, "aggregate_type")
		mustAddAnd = true
	}
	if filter.AggregateIDs != nil {
		if mustAddAnd {
			stmt.WriteString(" AND ")
		}
		filter.AggregateIDs.Write(stmt, "aggregate_id")
		mustAddAnd = true
	}
	if filter.EventTypes != nil {
		if mustAddAnd {
			stmt.WriteString(" AND ")
		}
		filter.EventTypes.Write(stmt, "event_type")
		mustAddAnd = true
	}
	if filter.Sequence != nil {
		if mustAddAnd {
			stmt.WriteString(" AND ")
		}
		filter.Sequence.Write(stmt, "sequence")
	}

	if len(filters) > 1 {
		stmt.Builder.WriteRune(')')
	}
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
