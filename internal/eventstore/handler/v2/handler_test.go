package handler

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestHandler_eventQuery(t *testing.T) {
	h := &Handler{
		bulkLimit: 200,
		eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
			"user": {"user.human.added", "user.human.password.changed"},
		},
	}

	t.Run("no cursor on empty position", func(t *testing.T) {
		builder := h.eventQuery(&state{instanceID: "inst"})
		assert.Equal(t, uint32(0), builder.GetOffset())
		assert.Nil(t, builder.GetEventSortKeyAfter())
		assert.True(t, builder.GetPositionAtLeast().IsZero())
	})

	t.Run("resumes after event sort key", func(t *testing.T) {
		cursor := eventstore.EventSortKey{
			Position:      decimal.NewFromFloat(1788336088.993079),
			InTxOrder:     5,
			AggregateType: "user",
			AggregateID:   "388963124960608862",
			Sequence:      5,
		}
		builder := h.eventQuery(&state{
			instanceID: "inst",
			cursor:     cursor,
		})
		assert.Equal(t, uint32(0), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().IsZero())
		require.NotNil(t, builder.GetEventSortKeyAfter())
		assert.Equal(t, cursor, *builder.GetEventSortKeyAfter())
	})
}

func TestHandler_eventsToStatements_inTxOrder(t *testing.T) {
	h := &Handler{
		projection: &projection{name: "users14"},
	}

	pos1 := decimal.NewFromFloat(1)
	pos2 := decimal.NewFromFloat(1.5)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos1,
			InTx:      1,
		},
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.email.verified",
			Seq:       2,
			Pos:       pos1,
			InTx:      2,
		},
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-b", Type: "user"},
			EventType: "user.human.password.changed",
			Seq:       10,
			Pos:       pos2,
			InTx:      1,
		},
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events)
	require.NoError(t, err)
	require.Len(t, statements, 3)
	assert.Equal(t, uint32(1), statements[0].inTxOrder)
	assert.Equal(t, uint32(2), statements[1].inTxOrder)
	assert.Equal(t, uint32(1), statements[2].inTxOrder)
	assert.Equal(t, "agg-b", statements[2].Aggregate.ID)
}

func TestHandler_eventsToStatements_samePositionInTxOrder(t *testing.T) {
	h := &Handler{
		projection: &projection{name: "users14"},
	}
	pos := decimal.NewFromInt(1)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos,
			InTx:      1,
		},
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-b", Type: "user"},
			EventType: "user.human.password.changed",
			Seq:       1,
			Pos:       pos,
			InTx:      1,
		},
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events)
	require.NoError(t, err)
	require.Len(t, statements, 2)
	assert.Equal(t, "agg-a", statements[0].Aggregate.ID)
	assert.Equal(t, "agg-b", statements[1].Aggregate.ID)
	assert.Equal(t, uint32(1), statements[0].inTxOrder)
	assert.Equal(t, uint32(1), statements[1].inTxOrder)
}

func TestSkipPreviouslyReduced(t *testing.T) {
	pos := decimal.NewFromInt(1)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:  &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			Seq:  1,
			Pos:  pos,
			InTx: 5,
		},
		&eventstore.BaseEvent{
			Agg:  &eventstore.Aggregate{ID: "agg-b", Type: "user"},
			Seq:  1,
			Pos:  pos,
			InTx: 1,
		},
	}

	t.Run("matches last reduced event", func(t *testing.T) {
		idx := skipPreviouslyReduced(events, eventstore.EventSortKeyFromEvent(events[0]), eventstore.EventSortKeyFromEvent)
		assert.Equal(t, 0, idx)
	})

	t.Run("matches without in_tx_order so pre-backfill 0 still skips", func(t *testing.T) {
		cursor := eventstore.EventSortKeyFromEvent(events[0])
		cursor.InTxOrder = 0
		idx := skipPreviouslyReduced(events, cursor, eventstore.EventSortKeyFromEvent)
		assert.Equal(t, 0, idx)
	})

	t.Run("unknown cursor keeps the batch", func(t *testing.T) {
		idx := skipPreviouslyReduced(events, eventstore.EventSortKey{
			Position:      pos,
			AggregateType: "user",
			AggregateID:   "agg-c",
			Sequence:      1,
		}, eventstore.EventSortKeyFromEvent)
		assert.Equal(t, -1, idx)
	})
}
