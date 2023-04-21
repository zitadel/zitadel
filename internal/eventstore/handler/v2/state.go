package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
)

type state struct {
	InstanceID     string
	EventTimestamp time.Time
}

var (
	//go:embed current_state.sql
	currentStateStmt string
	//go:embed update_state.sql
	updateStateStmt string
)

func (h *Handler) currentState(ctx context.Context, tx *sql.Tx) (*state, error) {
	row := tx.QueryRowContext(ctx, currentStateStmt, authz.GetInstance(ctx).InstanceID(), h.projection.Name())
	if row.Err() != nil {
		logging.WithError(row.Err()).Debug("unable to query current state")
		return nil, row.Err()
	}
	currentState := new(state)
	if err := row.Scan(&currentState.InstanceID, &currentState.EventTimestamp); err != nil {
		return nil, err
	}
	return currentState, nil
}

func (h *Handler) updateState(ctx context.Context, updatedState *state, tx *sql.Tx) error {
	res, err := tx.ExecContext(ctx, updateStateStmt, h.projection.Name(), updatedState.InstanceID, updatedState.EventTimestamp)
	if err != nil {
		logging.WithError(err).Debug("unable to update state")
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		logging.OnError(err).Error("unable to check if states are updated")
		return errors.ThrowInternal(err, "V2-lpiK0", "unable to update state")
	}
	return nil
}
