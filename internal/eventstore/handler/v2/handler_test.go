package handler

import (
	"context"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestHandler_eventQuery(t *testing.T) {
	h := &Handler{
		bulkLimit:            200,
		filterOffsetIsCursor: true,
		eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
			"user": {"user.human.added", "user.human.password.changed"},
		},
	}

	t.Run("scans every aggregate's event types separately", func(t *testing.T) {
		h := &Handler{
			bulkLimit:            200,
			filterOffsetIsCursor: true,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"user": {"user.human.added", "user.locked"},
				"org":  {"org.removed"},
			},
		}
		builder := h.eventQuery(&state{instanceID: "inst"}, decimal.Decimal{})
		assert.True(t, builder.GetScanEventTypesSeparately())
		queries := builder.GetQueries()
		require.Len(t, queries, 2)
		assert.Equal(t, []eventstore.AggregateType{"org"}, queries[0].GetAggregateTypes())
		assert.Equal(t, []eventstore.EventType{"org.removed"}, queries[0].GetEventTypes())
		assert.Equal(t, []eventstore.AggregateType{"user"}, queries[1].GetAggregateTypes())
		assert.Equal(t, []eventstore.EventType{"user.human.added", "user.locked"}, queries[1].GetEventTypes())
	})

	t.Run("aggregate without event types keeps the single query", func(t *testing.T) {
		h := &Handler{
			bulkLimit:            200,
			filterOffsetIsCursor: true,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"org": nil,
			},
		}
		builder := h.eventQuery(&state{instanceID: "inst"}, decimal.Decimal{})
		assert.False(t, builder.GetScanEventTypesSeparately())
		queries := builder.GetQueries()
		require.Len(t, queries, 1)
		assert.Equal(t, []eventstore.AggregateType{"org"}, queries[0].GetAggregateTypes())
		assert.Empty(t, queries[0].GetEventTypes())
	})

	t.Run("no cursor on empty position", func(t *testing.T) {
		builder := h.eventQuery(&state{instanceID: "inst"}, decimal.Decimal{})
		assert.Equal(t, uint32(0), builder.GetOffset())
		assert.Nil(t, builder.GetEventSortKeyAfter())
		assert.True(t, builder.GetPositionAtLeast().IsZero())
	})

	t.Run("77 stored resumes after event sort key from filter_offset", func(t *testing.T) {
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
		}, decimal.Decimal{})
		assert.Equal(t, uint32(0), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().IsZero())
		require.NotNil(t, builder.GetEventSortKeyAfter())
		assert.Equal(t, cursor, *builder.GetEventSortKeyAfter())
	})

	t.Run("77 not stored uses new-column sort key when ordinal set", func(t *testing.T) {
		h := &Handler{
			bulkLimit: 200,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"user": {"user.human.added"},
			},
		}
		cursor := eventstore.EventSortKey{
			Position:      decimal.NewFromInt(1),
			InTxOrder:     7,
			AggregateType: "user",
			AggregateID:   "agg-a",
			Sequence:      5,
		}
		builder := h.eventQuery(&state{
			instanceID:   "inst",
			cursor:       cursor,
			offset:       3,
			inTxOrderSet: true,
		}, decimal.Decimal{})
		assert.Equal(t, uint32(0), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().IsZero())
		require.NotNil(t, builder.GetEventSortKeyAfter())
		assert.Equal(t, cursor, *builder.GetEventSortKeyAfter())
	})

	t.Run("77 not stored uses OFFSET when ordinal is null", func(t *testing.T) {
		h := &Handler{
			bulkLimit: 200,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"user": {"user.human.added", "user.locked"},
				"org":  {"org.removed"},
			},
		}
		builder := h.eventQuery(&state{
			instanceID: "inst",
			offset:     4,
			cursor: eventstore.EventSortKey{
				Position:      decimal.NewFromInt(1),
				AggregateType: "user",
				AggregateID:   "agg-a",
				Sequence:      5,
			},
		}, decimal.Decimal{})
		assert.Equal(t, uint32(4), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().Equal(decimal.NewFromInt(1)))
		assert.Nil(t, builder.GetEventSortKeyAfter())
		assert.False(t, builder.GetScanEventTypesSeparately())
		queries := builder.GetQueries()
		require.Len(t, queries, 1)
	})

	t.Run("min position is an inclusive floor not a smashed cursor", func(t *testing.T) {
		cursor := eventstore.EventSortKey{
			Position:      decimal.NewFromInt(1),
			InTxOrder:     5,
			AggregateType: "user",
			AggregateID:   "agg-a",
			Sequence:      5,
		}
		minPosition := decimal.NewFromInt(10)
		builder := h.eventQuery(&state{
			instanceID: "inst",
			cursor:     cursor,
		}, minPosition)
		assert.Nil(t, builder.GetEventSortKeyAfter())
		assert.True(t, builder.GetPositionAtLeast().Equal(minPosition))
		assert.Equal(t, uint32(0), builder.GetOffset())
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

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{})
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

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{})
	require.NoError(t, err)
	require.Len(t, statements, 2)
	assert.Equal(t, "agg-a", statements[0].Aggregate.ID)
	assert.Equal(t, "agg-b", statements[1].Aggregate.ID)
	assert.Equal(t, uint32(1), statements[0].inTxOrder)
	assert.Equal(t, uint32(1), statements[1].inTxOrder)
}

func TestHandler_eventsToStatements_offsetCountDiffersFromInTxOrder(t *testing.T) {
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
			InTx:      5,
		},
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-b", Type: "user"},
			EventType: "user.human.password.changed",
			Seq:       1,
			Pos:       pos,
			InTx:      9,
		},
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{
		cursor: eventstore.EventSortKey{Position: pos},
	})
	require.NoError(t, err)
	require.Len(t, statements, 2)
	assert.Equal(t, uint32(1), statements[0].offset)
	assert.Equal(t, uint32(2), statements[1].offset)
	assert.Equal(t, uint32(5), statements[0].inTxOrder)
	assert.Equal(t, uint32(9), statements[1].inTxOrder)
	assert.NotEqual(t, statements[1].offset, statements[1].inTxOrder)
}

func TestHandler_eventsToStatements_cursorModeDoesNotCountOffset(t *testing.T) {
	h := &Handler{
		projection:           &projection{name: "users14"},
		filterOffsetIsCursor: true,
	}
	pos := decimal.NewFromInt(1)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos,
			InTx:      5,
		},
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{})
	require.NoError(t, err)
	require.Len(t, statements, 1)
	assert.Equal(t, uint32(0), statements[0].offset)
	assert.Equal(t, uint32(5), statements[0].inTxOrder)
}

func TestUpdateStateStmt_singleWriteOptsIn(t *testing.T) {
	assert.Contains(t, updateStateStmt, "set_config('zitadel.keep_in_tx_order', 'true', true)")
	assert.Contains(t, updateStateStmt, "in_tx_order")
	assert.NotContains(t, updateStateStmt, "UPDATE projections.current_states SET")
	assert.Equal(t, 1, strings.Count(updateStateStmt, "INSERT INTO projections.current_states"))
}
