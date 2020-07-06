package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

const insertStmt = "INSERT INTO eventstore.events " +
	"(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, editor_user, editor_service, resource_owner, previous_sequence) " +
	"SELECT $1, $2, $3, $4, COALESCE($5, now()), $6, $7, $8, $9, $10 " +
	"WHERE EXISTS (SELECT 1 WHERE " +
	// exactly one event of requested aggregate must have the given previous sequence as sequence (last inserted event)
	"EXISTS (SELECT 1 FROM eventstore.events WHERE event_sequence = COALESCE($11, 0) AND aggregate_type = $12 AND aggregate_id = $13) OR " +
	// if previous sequence = 0, no events must exist for the requested aggregate
	"NOT EXISTS (SELECT 1 FROM eventstore.events WHERE aggregate_type = $14 AND aggregate_id = $15) AND COALESCE($16, 0) = 0) " +
	"RETURNING id, event_sequence, creation_date"

func (db *SQL) PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	err = crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(insertStmt)
		if err != nil {
			tx.Rollback()
			logging.Log("SQL-9ctx5").WithError(err).Warn("prepare failed")
			return caos_errs.ThrowInternal(err, "SQL-juCgA", "prepare failed")
		}

		for _, aggregate := range aggregates {
			err = precondtion(tx, aggregate)
			if err != nil {
				tx.Rollback()
				return err
			}
			err = insertEvents(stmt, Sequence(aggregate.PreviousSequence), aggregate.Events)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})

	if err != nil && !errors.Is(err, &caos_errs.CaosError{}) {
		err = caos_errs.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}

func precondtion(tx *sql.Tx, aggregate *models.Aggregate) error {
	if aggregate.Precondition == nil {
		return nil
	}
	events, err := filter(tx, aggregate.Precondition.Query)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(err, "SQL-oBPxB", "filter failed")
	}
	err = aggregate.Precondition.Validation(events...)
	if err != nil {
		return err
	}
	return nil
}

func insertEvents(stmt *sql.Stmt, previousSequence Sequence, events []*models.Event) error {
	for _, event := range events {
		rows, err := stmt.Query(event.Type, event.AggregateType, event.AggregateID, event.AggregateVersion, event.CreationDate, Data(event.Data), event.EditorUser, event.EditorService, event.ResourceOwner, previousSequence,
			previousSequence, event.AggregateType, event.AggregateID,
			event.AggregateType, event.AggregateID, previousSequence)

		if err != nil {
			logging.Log("SQL-EXA0q").WithError(err).Info("query failed")
			return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
		}
		defer rows.Close()

		rowInserted := false
		for rows.Next() {
			rowInserted = true
			err = rows.Scan(&event.ID, &previousSequence, &event.CreationDate)
			logging.Log("SQL-rAvLD").OnError(err).Info("unable to scan result into event")
		}

		if !rowInserted {
			logging.LogWithFields("SQL-5aATu", "aggregate", event.AggregateType, "id", event.AggregateID).Info("wrong sequence")
			return caos_errs.ThrowAlreadyExists(nil, "SQL-GKcAa", "wrong sequence")
		}

		event.Sequence = uint64(previousSequence)
	}

	return nil
}
