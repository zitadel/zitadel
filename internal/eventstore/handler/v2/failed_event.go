package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed failed_event_set.sql
	setFailedEventStmt string
	//go:embed failed_event_get_count.sql
	failureCountStmt string
)

type failure struct {
	sequence      uint64
	instance      string
	aggregateID   string
	aggregateType eventstore.AggregateType
	eventDate     time.Time
	err           error
}

func (f *failure) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("sequence", f.sequence),
		slog.String("instance", f.instance),
		slog.String("aggregate_id", f.aggregateID),
	)
}

func failureFromEvent(event eventstore.Event, err error) *failure {
	return &failure{
		sequence:      event.Sequence(),
		instance:      event.Aggregate().InstanceID,
		aggregateID:   event.Aggregate().ID,
		aggregateType: event.Aggregate().Type,
		eventDate:     event.CreatedAt(),
		err:           err,
	}
}

func failureFromStatement(statement *Statement, err error) *failure {
	return &failure{
		sequence:      statement.Sequence,
		instance:      statement.Aggregate.InstanceID,
		aggregateID:   statement.Aggregate.ID,
		aggregateType: statement.Aggregate.Type,
		eventDate:     statement.CreationDate,
		err:           err,
	}
}

func (h *Handler) handleFailedStmt(ctx context.Context, tx *sql.Tx, f *failure) (shouldContinue bool) {
	ctx = logging.With(ctx, "failure", f)
	failureCount, err := h.failureCount(tx, f)
	if err != nil {
		logging.Warn(ctx, "unable to get failure count", "err", err)
		return false
	}
	failureCount += 1
	err = h.setFailureCount(tx, failureCount, f)
	logging.OnError(ctx, err).Warn("unable to update failure count")

	return failureCount >= h.maxFailureCount
}

func (h *Handler) failureCount(tx *sql.Tx, f *failure) (count uint8, err error) {
	row := tx.QueryRow(failureCountStmt,
		h.projection.Name(),
		f.instance,
		f.aggregateType,
		f.aggregateID,
		f.sequence,
	)
	if err = row.Err(); err != nil {
		return 0, zerrors.ThrowInternal(err, "CRDB-Unnex", "unable to update failure count")
	}
	if err = row.Scan(&count); err != nil {
		return 0, zerrors.ThrowInternal(err, "CRDB-RwSMV", "unable to scan count")
	}
	return count, nil
}

func (h *Handler) setFailureCount(tx *sql.Tx, count uint8, f *failure) error {
	_, err := tx.Exec(setFailedEventStmt,
		h.projection.Name(),
		f.instance,
		f.aggregateType,
		f.aggregateID,
		f.eventDate,
		f.sequence,
		count,
		f.err.Error(),
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "CRDB-4Ht4x", "set failure count failed")
	}
	return nil
}
