package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zitadel/logging"

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

	for queryIdx, query := range filter.EventQueries {
		stmt.Builder.WriteString(" AND ")

		if len(filter.EventQueries) > 1 {
			stmt.Builder.WriteRune('(')
		}

		for extIdx, ext := range query.Exts {
			extToFilter(stmt, ext)
			if extIdx < len(query.Exts)-1 {
				stmt.Builder.WriteString(" AND ")
			}
		}

		if len(filter.EventQueries) > 1 {
			stmt.Builder.WriteRune(')')
		}

		if queryIdx < len(filter.EventQueries)-1 {
			stmt.Builder.WriteString(" OR ")
		}
	}

	writeOrdering(stmt, filter.Descending)
}

func writeFilterClauses(stmt *database.Statement, filter *eventstore.Filter) {
	writeInstanceFilter(stmt, filter.Instances)
}

func writeInstanceFilter(stmt *database.Statement, instances []string) {
	if len(instances) == 1 {
		database.NewTextEqual(instances[0]).Write(stmt, "instance_id")
		return
	}
	database.NewListContains(instances).Write(stmt, "instance_id")
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

func extToFilter(stmt *database.Statement, ext eventstore.EventQueryExt) {
	switch filter := ext.(type) {
	case *eventstore.AggregateTypesFilter:
		if len(filter.Types()) == 1 {
			database.NewTextEqual(filter.Types()[0]).Write(stmt, "aggregate_type")
			return
		}
		database.NewListContains(filter.Types()).Write(stmt, "aggregate_type")
	case *eventstore.EventTypesFilter:
		if len(filter.Types()) == 1 {
			database.NewTextEqual(filter.Types()[0]).Write(stmt, "event_type")
			return
		}
		database.NewListContains(filter.Types()).Write(stmt, "event_type")
	case *eventstore.AggregateIDsFilter:
		if len(filter.IDs()) == 1 {
			database.NewTextEqual(filter.IDs()[0]).Write(stmt, "aggregate_id")
			return
		}
		database.NewListContains(filter.IDs()).Write(stmt, "aggregate_id")
	case *eventstore.SequenceFilter[eventstore.SequenceBetweenType]:
		filter.Filter().Write(stmt, "sequence")
	case *eventstore.SequenceFilter[eventstore.SequenceEqualsType]:
		filter.Filter().Write(stmt, "sequence")
	default:
		logging.WithFields("ext_type", fmt.Sprintf("%T", ext)).Panic("event filter not implemented")
	}
}
