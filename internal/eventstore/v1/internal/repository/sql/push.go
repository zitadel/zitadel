package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

const (
	insertStmt = "INSERT INTO eventstore.events " +
		"(event_type, aggregate_type, aggregate_id, aggregate_version, creation_date, event_data, editor_user, editor_service, resource_owner, previous_sequence) " +
		"SELECT $1, $2, $3, $4, COALESCE($5, now()), $6, $7, $8, $9, $10 " +
		"WHERE EXISTS (" +
		"SELECT 1 FROM eventstore.events WHERE aggregate_type = $11 AND aggregate_id = $12 HAVING MAX(event_sequence) = $13 OR ($14::BIGINT IS NULL AND COUNT(*) = 0)) " +
		"RETURNING event_sequence, creation_date"
)

func (db *SQL) PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	err = crdb.ExecuteTx(ctx, db.client, nil, func(tx *sql.Tx) error {
		stmt, err := tx.Prepare(insertStmt)
		if err != nil {
			tx.Rollback()
			logging.Log("SQL-9ctx5").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("prepare failed")
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
	events, err := filter(tx, models.FactoryFromSearchQuery(aggregate.Precondition.Query))
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
		creationDate := sql.NullTime{Time: event.CreationDate, Valid: !event.CreationDate.IsZero()}
		err := stmt.QueryRow(event.Type, event.AggregateType, event.AggregateID, event.AggregateVersion, creationDate, Data(event.Data), event.EditorUser, event.EditorService, event.ResourceOwner, previousSequence,
			event.AggregateType, event.AggregateID, previousSequence, previousSequence).Scan(&previousSequence, &event.CreationDate)

		if err != nil {
			logging.LogWithFields("SQL-5M0sd",
				"aggregate", event.AggregateType,
				"previousSequence", previousSequence,
				"aggregateId", event.AggregateID,
				"aggregateType", event.AggregateType,
				"eventType", event.Type).WithError(err).Info("query failed")
			return caos_errs.ThrowInternal(err, "SQL-SBP37", "unable to create event")
		}

		event.Sequence = uint64(previousSequence)
	}

	return nil
}
