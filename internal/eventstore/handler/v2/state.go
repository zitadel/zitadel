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
	position       float64
	eventTimestamp time.Time
	aggregateType  eventstore.AggregateType
	aggregateID    string
	sequence       uint64
}

var (
	//go:embed state_get.sql
	currentStateStmt string
	//go:embed state_get_await.sql
	currentStateAwaitStmt string
	//go:embed state_set.sql
	updateStateStmt string
	//go:embed state_lock.sql
	lockStateStmt string
	//go:embed state_set_last_run.sql
	updateStateLastRunStmt string

	errJustUpdated = errors.New("projection was just updated")
)

func (h *Handler) currentState(ctx context.Context, tx *sql.Tx, config *triggerConfig) (currentState *state, err error) {
	currentState = &state{
		instanceID: authz.GetInstance(ctx).InstanceID(),
	}

	var (
		aggregateID   = new(sql.NullString)
		aggregateType = new(sql.NullString)
		sequence      = new(sql.NullInt64)
		timestamp     = new(sql.NullTime)
		position      = new(sql.NullFloat64)
	)

	stateQuery := currentStateStmt
	if config.awaitRunning {
		stateQuery = currentStateAwaitStmt
	}

	row := tx.QueryRow(stateQuery, currentState.instanceID, h.projection.Name())
	err = row.Scan(
		aggregateID,
		aggregateType,
		sequence,
		timestamp,
		position,
	)
	if errors.Is(err, sql.ErrNoRows) {
		err = h.lockState(tx, currentState.instanceID)
	}
	if err != nil {
		h.log().WithError(err).Debug("unable to query current state")
		return nil, err
	}

	currentState.aggregateID = aggregateID.String
	currentState.aggregateType = eventstore.AggregateType(aggregateType.String)
	currentState.sequence = uint64(sequence.Int64)
	currentState.eventTimestamp = timestamp.Time
	currentState.position = position.Float64
	return currentState, nil
}

func (h *Handler) setState(tx *sql.Tx, updatedState *state) error {
	res, err := tx.Exec(updateStateStmt,
		h.projection.Name(),
		updatedState.instanceID,
		updatedState.aggregateID,
		updatedState.aggregateType,
		updatedState.sequence,
		updatedState.eventTimestamp,
		updatedState.position,
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

func (h *Handler) lockState(tx *sql.Tx, instanceID string) error {
	res, err := tx.Exec(lockStateStmt,
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
