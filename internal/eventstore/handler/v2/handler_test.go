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

	t.Run("resumes with sort key instead of OFFSET", func(t *testing.T) {
		position := decimal.NewFromFloat(1788336088.993079)
		builder := h.eventQuery(&state{
			instanceID:    "inst",
			position:      position,
			offset:        5,
			aggregateType: "user",
			aggregateID:   "388963124960608862",
			sequence:      5,
		})
		assert.Equal(t, uint32(0), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().IsZero())
		key := builder.GetEventSortKeyAfter()
		require.NotNil(t, key)
		assert.True(t, position.Equal(key.Position))
		assert.Equal(t, uint32(5), key.InTxOrder)
		assert.Equal(t, eventstore.AggregateType("user"), key.AggregateType)
		assert.Equal(t, "388963124960608862", key.AggregateID)
		assert.Equal(t, uint64(5), key.Sequence)
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
	assert.Equal(t, uint32(1), statements[0].offset)
	assert.Equal(t, uint32(2), statements[1].offset)
	assert.Equal(t, uint32(1), statements[2].offset)
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
	assert.Equal(t, uint32(1), statements[0].offset)
	assert.Equal(t, uint32(1), statements[1].offset)
}

func TestSkipPreviouslyReducedEvents(t *testing.T) {
	pos := decimal.NewFromInt(1)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			Seq: 1,
			Pos: pos,
		},
		&eventstore.BaseEvent{
			Agg: &eventstore.Aggregate{ID: "agg-b", Type: "user"},
			Seq: 1,
			Pos: pos,
		},
	}

	idx := skipPreviouslyReducedEvents(events, &state{
		position:      pos,
		aggregateID:   "agg-a",
		aggregateType: "user",
		sequence:      1,
	})
	assert.Equal(t, 0, idx)
}
