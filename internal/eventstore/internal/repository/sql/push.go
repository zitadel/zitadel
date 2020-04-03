package sql

import (
	"context"
	"database/sql"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/cockroachdb/cockroach-go/crdb"
)

const insertStmt = "insert into eventstore.events " +
	"(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, editor_user, editor_service, editor_tenant, resource_owner, previous_sequence) " +
	"select $1, $2, $3, $4, coalesce($5, now()), $6, $7, $8, $9, $10, " +
	// case is to set the highest sequence or NULL in previous_sequence
	"case (select exists(select event_sequence from eventstore.events where aggregate_type = $11 AND aggregate_id = $12)) " +
	"WHEN true then (select event_sequence from eventstore.events where aggregate_type = $13 AND aggregate_id = $14 order by event_sequence desc limit 1) " +
	"ELSE NULL " +
	"end " +
	"where (" +
	// exactly one event of requested aggregate must have a >= sequence (last inserted event)
	"(select count(id) from eventstore.events where event_sequence >= $15 AND aggregate_type = $16 AND aggregate_id = $17) = 1 OR " +
	// previous sequence = 0, no events must exist for the requested aggregate
	"((select count(id) from eventstore.events where aggregate_type = $18 and aggregate_id = $19) = 0 AND $20 = 0)) " +
	"RETURNING id, event_sequence, creation_date"

func (db *SQL) PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	err = crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(insertStmt)
		if err != nil {
			tx.Rollback()
			logging.Log("SQL-9ctx5").WithError(err).Warn("prepare failed")
			return errors.ThrowInternal(err, "SQL-juCgA", "prepare failed")
		}

		for _, aggregate := range aggregates {
			err = insertEvents(stmt, aggregate.Events)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})

	if _, ok := err.(*errors.CaosError); !ok && err != nil {
		err = errors.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}

func insertEvents(stmt *sql.Stmt, events []*models.Event) error {
	previousSequence := events[0].PreviousSequence
	for _, event := range events {
		event.PreviousSequence = previousSequence

		if event.Data == nil || len(event.Data) == 0 {
			//json decoder failes with EOF if json text is empty
			event.Data = []byte("{}")
		}

		rows, err := stmt.Query(event.Type, event.AggregateType, event.AggregateID, event.AggregateVersion, event.CreationDate, event.Data, event.EditorUser, event.EditorService, event.EditorOrg, event.ResourceOwner,
			event.AggregateType, event.AggregateID,
			event.AggregateType, event.AggregateID,
			event.PreviousSequence, event.AggregateType, event.AggregateID,
			event.AggregateType, event.AggregateID, event.PreviousSequence)

		if err != nil {
			logging.Log("SQL-EXA0q").WithError(err).Info("query failed")
			return errors.ThrowInternal(err, "SQL-SBP37", "unable to create event")
		}
		defer rows.Close()

		rowInserted := false
		for rows.Next() {
			rowInserted = true
			err = rows.Scan(&event.ID, &event.Sequence, &event.CreationDate)
			logging.Log("SQL-rAvLD").OnError(err).Info("unable to scan result into event")
		}

		if !rowInserted {
			return errors.ThrowAlreadyExists(nil, "SQL-GKcAa", "wrong sequence")
		}

		previousSequence = event.Sequence
	}

	return nil
}
