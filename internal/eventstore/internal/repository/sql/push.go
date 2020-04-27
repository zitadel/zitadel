package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/cockroachdb/cockroach-go/crdb"
)

const insertStmt = "insert into eventstore.events " +
	"(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, editor_user, editor_service, resource_owner, previous_sequence) " +
	"select $1, $2, $3, $4, coalesce($5, now()), $6, $7, $8, $9, " +
	// case is to set the highest sequence or NULL in previous_sequence
	"case (select exists(select event_sequence from eventstore.events where aggregate_type = $10 AND aggregate_id = $11)) " +
	"WHEN true then (select event_sequence from eventstore.events where aggregate_type = $12 AND aggregate_id = $13 order by event_sequence desc limit 1) " +
	"ELSE NULL " +
	"end " +
	"where (" +
	// exactly one event of requested aggregate must have a >= sequence (last inserted event)
	"(select count(id) from eventstore.events where event_sequence >= COALESCE($14, 0) AND aggregate_type = $15 AND aggregate_id = $16) = 1 OR " +
	// previous sequence = 0, no events must exist for the requested aggregate
	"((select count(id) from eventstore.events where aggregate_type = $17 and aggregate_id = $18) = 0 AND COALESCE($19, 0) = 0)) " +
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
		return caos_errs.ThrowPreconditionFailed(err, "SQL-s6hqU", "validation failed")
	}
	return nil
}

func insertEvents(stmt *sql.Stmt, previousSequence Sequence, events []*models.Event) error {
	for _, event := range events {
		rows, err := stmt.Query(event.Type, event.AggregateType, event.AggregateID, event.AggregateVersion, event.CreationDate, Data(event.Data), event.EditorUser, event.EditorService, event.ResourceOwner,
			event.AggregateType, event.AggregateID,
			event.AggregateType, event.AggregateID,
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
			return caos_errs.ThrowAlreadyExists(nil, "SQL-GKcAa", "wrong sequence")
		}

		event.Sequence = uint64(previousSequence)
	}

	return nil
}
