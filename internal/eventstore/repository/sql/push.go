package sql

import (
	"context"

	"github.com/caos/utils/logging"
	"github.com/cockroachdb/cockroach-go/crdb"

	"database/sql"

	caos_errs "github.com/caos/utils/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
)

func (db *SQL) PushEvents(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	err = crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare("insert into eventstore.events " +
			"(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, modifier_user, modifier_service, modifier_tenant, resource_owner, previous_sequence) " +
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
			"RETURNING id, event_sequence, creation_date")
		if err != nil {
			tx.Rollback()
			logging.Log("SQL-9ctx5").WithError(err).Warn("prepare failed")
			return caos_errs.ThrowInternal(err, "SQL-juCgA", "prepare failed")
		}

		for _, aggregate := range aggregates {
			previousSequence := aggregate.LatestSequence
			for _, event := range aggregate.Events {
				event.PreviousSequence = previousSequence
				event.AggregateType = aggregate.Typ
				event.AggregateID = aggregate.ID
				event.AggregateVersion = aggregate.Version

				if event.Data == nil || len(event.Data) == 0 {
					//json decoder failes with EOF if json text is empty
					event.Data = []byte("{}")
				}

				rows, err := stmt.Query(event.Typ, event.AggregateType, event.AggregateID, event.AggregateVersion, event.CreationDate, event.Data, event.ModifierUser, event.ModifierService, event.ModifierTenant, event.ResourceOwner,
					event.AggregateType, event.AggregateID,
					event.AggregateType, event.AggregateID,
					event.PreviousSequence, event.AggregateType, event.AggregateID,
					event.AggregateType, event.AggregateID, event.PreviousSequence)

				if err != nil {
					logging.Log("SQL-EXA0q").WithError(err).Info("query failed")
					tx.Rollback()
					return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
				}
				defer rows.Close()

				rowInserted := false
				for rows.Next() {
					rowInserted = true
					err = rows.Scan(&event.ID, &event.Sequence, &event.CreationDate)
					logging.Log("SQL-rAvLD").OnError(err).Info("unable to scan result into event")
				}

				if !rowInserted {
					tx.Rollback()
					return caos_errs.ThrowAlreadyExists(nil, "SQL-GKcAa", "wrong sequence")
				}

				previousSequence = event.Sequence
			}
		}
		return nil
	})

	if _, ok := err.(*caos_errs.CaosError); !ok && err != nil {
		err = caos_errs.ThrowInternal(err, "SQL-DjgtG", "unable to store events")
	}

	return err
}
