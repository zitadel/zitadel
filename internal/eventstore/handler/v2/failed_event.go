package handler

import (
	"database/sql"
	_ "embed"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	//go:embed failed_event_set.sql
	setFailedEventStmt string
	//go:embed failed_event_get_count.sql
	failureCountStmt string
)

type failure struct {
	sequence  uint64
	instance  string
	aggregate string
	eventDate time.Time
	err       error
}

func failureFromEvent(event eventstore.Event, err error) *failure {
	return &failure{
		sequence:  event.Sequence(),
		instance:  event.Aggregate().InstanceID,
		aggregate: event.Aggregate().ID,
		eventDate: event.CreationDate(),
		err:       err,
	}
}

func failureFromStatement(statement *Statement, err error) *failure {
	return &failure{
		sequence:  statement.Sequence,
		instance:  statement.InstanceID,
		aggregate: statement.AggregateID,
		eventDate: statement.CreationDate,
		err:       err,
	}
}

// TODO: if we use the tx here the insert will be reverted
func (h *Handler) handleFailedStmt(tx *sql.Tx, currentState *state, f *failure) (shouldContinue bool) {
	if currentState.eventTimestamp.After(f.eventDate) || currentState.eventTimestamp.Equal(f.eventDate) {
		return true
	}
	failureCount, err := h.failureCount(tx, f)
	if err != nil {
		h.logFailure(f).WithError(err).Warn("unable to get failure count")
		return false
	}
	failureCount += 1
	err = h.setFailureCount(tx, failureCount, f)
	h.logFailure(f).OnError(err).Warn("unable to update failure count")

	return failureCount >= h.maxFailureCount
}

func (h *Handler) failureCount(tx *sql.Tx, f *failure) (count uint8, err error) {
	row := tx.QueryRow(failureCountStmt, h.projection.Name(), f.sequence, f.instance)
	if err = row.Err(); err != nil {
		return 0, errors.ThrowInternal(err, "CRDB-Unnex", "unable to update failure count")
	}
	if err = row.Scan(&count); err != nil {
		return 0, errors.ThrowInternal(err, "CRDB-RwSMV", "unable to scan count")
	}
	return count, nil
}

func (h *Handler) setFailureCount(tx *sql.Tx, count uint8, f *failure) error {
	_, dbErr := tx.Exec(setFailedEventStmt, h.projection.Name(), f.sequence, count, f.err.Error(), f.instance)
	if dbErr != nil {
		return errors.ThrowInternal(dbErr, "CRDB-4Ht4x", "set failure count failed")
	}
	return nil
}
