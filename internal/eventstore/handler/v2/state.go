package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	errs "github.com/zitadel/zitadel/internal/errors"
)

type state struct {
	InstanceID     string
	EventTimestamp time.Time
	EventSequence  uint64
}

var (
	//go:embed state_get.sql
	currentStateStmt string
	//go:embed state_set.sql
	updateStateStmt string
)

func (h *Handler) currentState(ctx context.Context, tx *sql.Tx) (*state, error) {
	row := tx.QueryRowContext(ctx, currentStateStmt, authz.GetInstance(ctx).InstanceID(), h.projection.Name())
	currentState := new(state)
	err := row.Scan(&currentState.InstanceID, &currentState.EventTimestamp, &currentState.EventSequence)
	if errors.Is(err, sql.ErrNoRows) {
		initialState := &state{
			InstanceID: authz.GetInstance(ctx).InstanceID(),
		}
		err := h.setState(ctx, initialState, tx)
		if err != nil {
			return nil, err
		}
		return initialState, nil
	}
	if err != nil {
		logging.WithError(err).Debug("unable to query current state")
		return nil, err
	}
	return currentState, nil
}

func (h *Handler) setState(ctx context.Context, updatedState *state, tx *sql.Tx) error {
	res, err := tx.ExecContext(ctx, updateStateStmt, h.projection.Name(), updatedState.InstanceID, updatedState.EventTimestamp, updatedState.EventSequence)
	if err != nil {
		logging.WithError(err).Debug("unable to update state")
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		logging.OnError(err).Error("unable to check if states are updated")
		return errs.ThrowInternal(err, "V2-lpiK0", "unable to update state")
	}
	return nil
}
