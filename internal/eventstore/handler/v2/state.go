package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type state struct {
	instanceID     string
	position       decimal.Decimal
	eventTimestamp time.Time
	aggregateType  eventstore.AggregateType
	aggregateID    string
	sequence       uint64
	offset         uint32
}

var (
	//go:embed state_get.sql
	currentStateStmt string
	//go:embed state_set.sql
	updateStateStmt string
)

func (h *Handler) currentState(ctx context.Context, tx *sql.Tx) (currentState *state, err error) {
	currentState = &state{
		instanceID: authz.GetInstance(ctx).InstanceID(),
	}

	var (
		aggregateID   = new(sql.NullString)
		aggregateType = new(sql.NullString)
		sequence      = new(sql.NullInt64)
		timestamp     = new(sql.NullTime)
		position      = new(decimal.NullDecimal)
		offset        = new(sql.NullInt64)
	)

	row := tx.QueryRow(currentStateStmt, currentState.instanceID, h.projection.Name())
	err = row.Scan(
		aggregateID,
		aggregateType,
		sequence,
		timestamp,
		position,
		offset,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		h.log().WithError(err).Debug("unable to query current state")
		return nil, err
	}

	currentState.aggregateID = aggregateID.String
	currentState.aggregateType = eventstore.AggregateType(aggregateType.String)
	currentState.sequence = uint64(sequence.Int64)
	currentState.eventTimestamp = timestamp.Time
	currentState.position = position.Decimal
	// psql does not provide unsigned numbers so we work around it
	currentState.offset = uint32(offset.Int64)
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
		updatedState.offset,
	)
	if err != nil {
		h.log().WithError(err).Warn("unable to update state")
		return zerrors.ThrowInternal(err, "V2-WF23g2", "unable to update state")
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		h.log().OnError(err).Error("unable to check if states are updated")
		return zerrors.ThrowInternal(err, "V2-FGEKi", "unable to update state")
	}
	return nil
}
