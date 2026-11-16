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
	offset         uint32
	inTxOrderSet   bool
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
		filterOffset  = new(sql.NullInt64)
		inTxOrder     = new(sql.NullInt64)
	)

	row := tx.QueryRow(currentStateStmt, currentState.instanceID, h.projection.Name())
	err = row.Scan(
		aggregateID,
		aggregateType,
		sequence,
		timestamp,
		position,
		filterOffset,
		inTxOrder,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logging.WithError(ctx, err).Debug("unable to query current state")
		return nil, err
	}

	currentState.eventTimestamp = timestamp.Time
	currentState.cursor = eventstore.EventSortKey{
		Position:      position.Decimal,
		InstanceID:    currentState.instanceID,
		AggregateType: eventstore.AggregateType(aggregateType.String),
		AggregateID:   aggregateID.String,
		Sequence:      uint64(sequence.Int64),
	}
	h.hydrateState(currentState, *filterOffset, *inTxOrder)
	return currentState, nil
}

func (h *Handler) hydrateState(s *state, filterOffset, inTxOrder sql.NullInt64) {
	if h.filterOffsetIsCursor {
		s.cursor.InTxOrder = uint32(filterOffset.Int64)
		return
	}
	s.offset = uint32(filterOffset.Int64)
	// Valid includes 0: v2/v4 push stores 0-based in_tx_order for the first event of a tx.
	if inTxOrder.Valid {
		s.cursor.InTxOrder = uint32(inTxOrder.Int64)
		s.inTxOrderSet = true
	}
}

func (h *Handler) storedOffsets(s *state) (filterOffset any, inTxOrder any) {
	var inTx any
	if h.filterOffsetIsCursor || s.inTxOrderSet {
		inTx = s.cursor.InTxOrder
	}
	if h.filterOffsetIsCursor {
		return s.cursor.InTxOrder, inTx
	}
	return s.offset, inTx
}

func (h *Handler) setState(ctx context.Context, tx *sql.Tx, updatedState *state) error {
	filterOffset, inTxOrder := h.storedOffsets(updatedState)
	res, err := tx.Exec(updateStateStmt,
		h.projection.Name(),
		updatedState.instanceID,
		updatedState.cursor.AggregateID,
		updatedState.cursor.AggregateType,
		updatedState.cursor.Sequence,
		updatedState.eventTimestamp,
		updatedState.cursor.Position,
		filterOffset,
		inTxOrder,
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

func nextFilterOffset(prev decimal.Decimal, offset uint32, pos decimal.Decimal) (decimal.Decimal, uint32) {
	if prev.Equal(pos) {
		return pos, offset + 1
	}
	return pos, 1
}

// seedFilterOffset returns the starting previousPosition/offset for counting
// filter_offset over a fetched batch. Seed the stored values only when the
// query continued from that cursor (OFFSET or sort-key resume). A minPosition
// floor fetches with PositionAtLeast and no OFFSET, so recount from zero.
func seedFilterOffset(s *state, r resume) (previousPosition decimal.Decimal, offset uint32) {
	if s == nil {
		return decimal.Decimal{}, 0
	}
	if r.offset > 0 || r.afterSortKey {
		return s.cursor.Position, s.offset
	}
	return decimal.Decimal{}, 0
}

func advanceFilterOffsetOverEvents(prev decimal.Decimal, offset uint32, events []eventstore.Event) (decimal.Decimal, uint32) {
	for _, event := range events {
		prev, offset = nextFilterOffset(prev, offset, event.Position())
	}
	return prev, offset
}

func (s *state) applyStatement(stmt *Statement) {
	s.cursor = stmt.eventSortKey()
	s.eventTimestamp = stmt.CreationDate
	s.offset = stmt.offset
	s.inTxOrderSet = true
}

func (s *state) applyEvent(event eventstore.Event) {
	s.cursor = eventstore.EventSortKeyFromEvent(event)
	s.eventTimestamp = event.CreatedAt()
	s.inTxOrderSet = true
}
