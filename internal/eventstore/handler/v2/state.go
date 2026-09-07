package handler

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type state struct {
	instanceID     string
	eventTimestamp time.Time
	cursor         eventstore.EventSortKey
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
		inTxOrder     = new(sql.NullInt64)
	)

	row := tx.QueryRow(currentStateStmt, currentState.instanceID, h.projection.Name())
	err = row.Scan(
		aggregateID,
		aggregateType,
		sequence,
		timestamp,
		position,
		inTxOrder,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logging.WithError(ctx, err).Debug("unable to query current state")
		return nil, err
	}

	currentState.eventTimestamp = timestamp.Time
	currentState.cursor = eventstore.EventSortKey{
		Position:      position.Decimal,
		InTxOrder:     uint32(inTxOrder.Int64),
		AggregateType: eventstore.AggregateType(aggregateType.String),
		AggregateID:   aggregateID.String,
		Sequence:      uint64(sequence.Int64),
	}
	return currentState, nil
}

func (h *Handler) setState(ctx context.Context, tx *sql.Tx, updatedState *state) error {
	res, err := tx.Exec(updateStateStmt,
		h.projection.Name(),
		updatedState.instanceID,
		updatedState.cursor.AggregateID,
		updatedState.cursor.AggregateType,
		updatedState.cursor.Sequence,
		updatedState.eventTimestamp,
		updatedState.cursor.Position,
		updatedState.cursor.InTxOrder,
	)
	if err != nil {
		err = zerrors.ThrowInternal(err, "V2-WF23g2", "unable to update state")
		logging.Warn(ctx, "unable to update state", "err", err)
		return err
	}
	if affected, err := res.RowsAffected(); affected == 0 {
		err = zerrors.ThrowInternal(err, "V2-FGEKi", "unable to update state")
		logging.Error(ctx, "unable to check if states are updated", "err", err)
		return err
	}
	return nil
}

func (s *state) applyStatement(stmt *Statement) {
	s.cursor = stmt.eventSortKey()
	s.eventTimestamp = stmt.CreationDate
}

func (s *state) applyEvent(event eventstore.Event) {
	s.cursor = eventstore.EventSortKeyFromEvent(event)
	s.eventTimestamp = event.CreatedAt()
}
