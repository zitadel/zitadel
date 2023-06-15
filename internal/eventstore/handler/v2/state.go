package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type state struct {
	instanceID     string
	aggregateType  eventstore.AggregateType
	aggregateID    string
	eventTimestamp time.Time
	eventSequence  uint64
}

var (
	//go:embed state_get.sql
	currentStateStmt string
	//go:embed state_set.sql
	updateStateStmt string
	//go:embed state_lock.sql
	lockStateStmt string
	//go:embed state_set_last_run.sql
	updateStateLastRunStmt string
)

func (h *Handler) currentState(ctx context.Context, tx *sql.Tx) (currentState *state, err error) {
	currentState = &state{
		instanceID: authz.GetInstance(ctx).InstanceID(),
	}

	timestamp := new(sql.NullTime)
	aggregateType := new(sql.NullString)
	aggregateID := new(sql.NullString)
	sequence := new(sql.NullInt64)
	row := tx.QueryRowContext(ctx, currentStateStmt, currentState.instanceID, h.projection.Name())
	err = row.Scan(
		timestamp,
		aggregateType,
		aggregateID,
		sequence,
	)
	if errors.Is(err, sql.ErrNoRows) {
		err = h.lockState(ctx, tx, currentState.instanceID)
	}
	if err != nil {
		h.log().WithError(err).Debug("unable to query current state")
		return nil, err
	}

	currentState.eventTimestamp = timestamp.Time
	currentState.aggregateType = eventstore.AggregateType(aggregateType.String)
	currentState.aggregateID = aggregateID.String
	currentState.eventSequence = uint64(sequence.Int64)
	return currentState, nil
}

func (h *Handler) setState(ctx context.Context, tx *sql.Tx, updatedState *state) error {
	res, err := tx.ExecContext(ctx, updateStateStmt,
		h.projection.Name(),
		updatedState.instanceID,
		updatedState.eventTimestamp,
		updatedState.aggregateType,
		updatedState.aggregateID,
		updatedState.eventSequence,
	)
	if err != nil {
		h.log().WithError(err).Debug("unable to update state")
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		h.log().OnError(err).Error("unable to check if states are updated")
		return errs.ThrowInternal(err, "V2-FGEKi", "unable to update state")
	}
	return nil
}

func (h *Handler) updateLastUpdated(ctx context.Context, tx *sql.Tx, updatedState *state) {
	_, err := tx.ExecContext(ctx, updateStateLastRunStmt, h.projection.Name(), updatedState.instanceID)
	h.log().OnError(err).Debug("unable to update last updated")
}

func (h *Handler) lockState(ctx context.Context, tx *sql.Tx, instanceID string) error {
	res, err := tx.ExecContext(ctx, lockStateStmt,
		h.projection.Name(),
		instanceID,
	)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 || err != nil {
		return errs.ThrowInternal(err, "V2-lpiK0", "projection already locked")
	}
	return nil
}
