package handler

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database/mock"
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

	t.Run("77 not stored uses OFFSET when ordinal set", func(t *testing.T) {
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
		assert.Equal(t, uint32(3), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().Equal(decimal.NewFromInt(1)))
		assert.Nil(t, builder.GetEventSortKeyAfter())
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
		assert.True(t, builder.GetScanEventTypesSeparately())
	})

	t.Run("OFFSET mode never scans event types separately", func(t *testing.T) {
		h := &Handler{
			bulkLimit: 200,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"user": {"user.human.added", "user.locked"},
				"org":  {"org.removed"},
			},
		}
		builder := h.eventQuery(&state{instanceID: "inst"}, decimal.Decimal{})
		assert.False(t, builder.GetScanEventTypesSeparately())
		require.Len(t, builder.GetQueries(), 1)
	})

	t.Run("OFFSET mode with in_tx_order set still uses single query", func(t *testing.T) {
		h := &Handler{
			bulkLimit: 200,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"user": {"user.human.added"},
				"org":  {"org.removed"},
			},
		}
		builder := h.eventQuery(&state{
			instanceID:   "inst",
			inTxOrderSet: true,
			offset:       2,
			cursor: eventstore.EventSortKey{
				Position:      decimal.NewFromInt(1),
				InTxOrder:     0,
				AggregateType: "user",
				AggregateID:   "agg-a",
				Sequence:      5,
			},
		}, decimal.Decimal{})
		assert.False(t, builder.GetScanEventTypesSeparately())
		assert.Equal(t, uint32(2), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().Equal(decimal.NewFromInt(1)))
		assert.Nil(t, builder.GetEventSortKeyAfter())
		require.Len(t, builder.GetQueries(), 1)
	})

	t.Run("zero in_tx_order uses OFFSET when 77 not stored", func(t *testing.T) {
		h := &Handler{
			bulkLimit: 200,
			eventTypes: map[eventstore.AggregateType][]eventstore.EventType{
				"user": {"user.human.added"},
			},
		}
		cursor := eventstore.EventSortKey{
			Position:      decimal.NewFromInt(1),
			InTxOrder:     0,
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
		assert.Equal(t, uint32(3), builder.GetOffset())
		assert.True(t, builder.GetPositionAtLeast().Equal(decimal.NewFromInt(1)))
		assert.Nil(t, builder.GetEventSortKeyAfter())
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

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{}, decimal.Decimal{})
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

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{}, decimal.Decimal{})
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
	}, decimal.Decimal{})
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

	statements, err := h.eventsToStatements(context.Background(), nil, events, &state{}, decimal.Decimal{})
	require.NoError(t, err)
	require.Len(t, statements, 1)
	assert.Equal(t, uint32(0), statements[0].offset)
	assert.Equal(t, uint32(5), statements[0].inTxOrder)
}

func TestHandler_eventsToStatements_minPositionDoesNotSeedStoredOffset(t *testing.T) {
	h := &Handler{projection: &projection{name: "users14"}}
	pos := decimal.NewFromInt(10)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos,
			InTx:      1,
		},
	}
	current := &state{
		cursor: eventstore.EventSortKey{Position: pos},
		offset: 1,
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events, current, pos)
	require.NoError(t, err)
	require.Len(t, statements, 1)
	assert.Equal(t, uint32(1), statements[0].offset)
}

func TestHandler_eventsToStatements_offsetResumeSeedsStoredOffset(t *testing.T) {
	h := &Handler{projection: &projection{name: "users14"}}
	pos := decimal.NewFromInt(10)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos,
			InTx:      4,
		},
	}
	current := &state{
		cursor: eventstore.EventSortKey{Position: pos},
		offset: 3,
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events, current, decimal.Decimal{})
	require.NoError(t, err)
	require.Len(t, statements, 1)
	assert.Equal(t, uint32(4), statements[0].offset)
}

func TestHandler_eventsToStatements_sortKeyResumeSeedsDualWriteOffset(t *testing.T) {
	h := &Handler{projection: &projection{name: "users14"}}
	pos := decimal.NewFromInt(10)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       2,
			Pos:       pos,
			InTx:      0,
		},
	}
	current := &state{
		cursor: eventstore.EventSortKey{
			Position:      pos,
			InTxOrder:     0,
			AggregateType: "user",
			AggregateID:   "agg-prev",
			Sequence:      1,
		},
		offset:       2,
		inTxOrderSet: true,
	}

	statements, err := h.eventsToStatements(context.Background(), nil, events, current, decimal.Decimal{})
	require.NoError(t, err)
	require.Len(t, statements, 1)
	assert.Equal(t, uint32(3), statements[0].offset)
}

func TestSeedFilterOffset(t *testing.T) {
	pos := decimal.NewFromInt(5)
	s := &state{cursor: eventstore.EventSortKey{Position: pos}, offset: 2}

	t.Run("minPosition does not seed", func(t *testing.T) {
		prev, offset := seedFilterOffset(s, resume{position: pos})
		assert.True(t, prev.IsZero())
		assert.Equal(t, uint32(0), offset)
	})
	t.Run("OFFSET resume seeds", func(t *testing.T) {
		prev, offset := seedFilterOffset(s, resume{position: pos, offset: 2})
		assert.True(t, prev.Equal(pos))
		assert.Equal(t, uint32(2), offset)
	})
	t.Run("sort-key resume seeds", func(t *testing.T) {
		prev, offset := seedFilterOffset(s, resume{afterSortKey: true})
		assert.True(t, prev.Equal(pos))
		assert.Equal(t, uint32(2), offset)
	})
	t.Run("nil state", func(t *testing.T) {
		prev, offset := seedFilterOffset(nil, resume{afterSortKey: true})
		assert.True(t, prev.IsZero())
		assert.Equal(t, uint32(0), offset)
	})
}

func TestAdvanceFilterOffsetOverEvents(t *testing.T) {
	pos := decimal.NewFromInt(1)
	events := []eventstore.Event{
		&eventstore.BaseEvent{Pos: pos},
		&eventstore.BaseEvent{Pos: pos},
		&eventstore.BaseEvent{Pos: decimal.NewFromInt(2)},
	}
	prev, offset := advanceFilterOffsetOverEvents(decimal.Decimal{}, 0, events)
	assert.True(t, prev.Equal(decimal.NewFromInt(2)))
	assert.Equal(t, uint32(1), offset)

	prev, offset = advanceFilterOffsetOverEvents(pos, 2, events[:1])
	assert.True(t, prev.Equal(pos))
	assert.Equal(t, uint32(3), offset)
}

func TestFieldHandler_fetchEvents_minPositionDoesNotSeedStoredOffset(t *testing.T) {
	pos := decimal.NewFromInt(10)
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       &eventstore.Aggregate{ID: "agg-a", Type: "user"},
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos,
			InTx:      1,
		},
	}
	es := &stubFilterEventstore{events: events}
	h := &FieldHandler{Handler: Handler{
		es:         es,
		bulkLimit:  200,
		eventTypes: map[eventstore.AggregateType][]eventstore.EventType{"user": {"user.human.added"}},
	}}
	current := &state{
		instanceID: "inst",
		cursor:     eventstore.EventSortKey{Position: pos},
		offset:     1,
	}

	got, more, err := h.fetchEvents(context.Background(), nil, current, pos)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.False(t, more)
	assert.Equal(t, uint32(1), current.offset)
}

type stubFilterEventstore struct {
	events []eventstore.Event
}

func (s *stubFilterEventstore) InstanceIDs(context.Context, *eventstore.SearchQueryBuilder) ([]string, error) {
	return nil, nil
}
func (s *stubFilterEventstore) FilterToQueryReducer(context.Context, eventstore.QueryReducer) error {
	return nil
}
func (s *stubFilterEventstore) Filter(context.Context, *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
	return s.events, nil
}
func (s *stubFilterEventstore) Push(context.Context, ...eventstore.Command) ([]eventstore.Event, error) {
	return nil, nil
}
func (s *stubFilterEventstore) FillFields(context.Context, ...eventstore.FillFieldsEvent) error {
	return nil
}

func TestHandler_eventsToStatements_skippedReduceAdvancesOffset(t *testing.T) {
	sqlMock := mock.NewSQLMock(t,
		mock.ExpectBegin(nil),
		mock.ExpectQuery(failureCountStmt,
			mock.WithQueryResult([]string{"failure_count"}, [][]driver.Value{{int64(0)}}),
		),
		mock.ExcpectExec(setFailedEventStmt,
			mock.WithExecRowsAffected(1),
		),
	)
	t.Cleanup(func() { sqlMock.Assert(t) })

	tx, err := sqlMock.DB.BeginTx(context.Background(), nil)
	require.NoError(t, err)

	reduceErr := errors.New("reduce failed")
	h := &Handler{
		projection: &projection{
			name: "users14",
			reducers: []AggregateReducer{{
				Aggregate: "user",
				EventReducers: []EventReducer{
					{
						Event: "user.human.added",
						Reduce: func(eventstore.Event) (*Statement, error) {
							return nil, reduceErr
						},
					},
					{
						Event: "user.human.password.changed",
						Reduce: func(event eventstore.Event) (*Statement, error) {
							return NewStatement(event, func(context.Context, Executer, string) error { return nil }), nil
						},
					},
				},
			}},
		},
		maxFailureCount: 1,
	}
	pos := decimal.NewFromInt(1)
	agg := &eventstore.Aggregate{ID: "agg-a", Type: "user", InstanceID: "inst"}
	events := []eventstore.Event{
		&eventstore.BaseEvent{
			Agg:       agg,
			EventType: "user.human.added",
			Seq:       1,
			Pos:       pos,
			InTx:      5,
		},
		&eventstore.BaseEvent{
			Agg:       agg,
			EventType: "user.human.password.changed",
			Seq:       2,
			Pos:       pos,
			InTx:      9,
		},
	}

	statements, err := h.eventsToStatements(context.Background(), tx, events, &state{}, decimal.Decimal{})
	require.NoError(t, err)
	require.Len(t, statements, 2)
	assert.Nil(t, statements[0].Execute)
	assert.NotNil(t, statements[1].Execute)
	assert.Equal(t, uint32(1), statements[0].offset)
	assert.Equal(t, uint32(2), statements[1].offset)
}

func TestUpdateStateStmt_singleWriteOptsIn(t *testing.T) {
	assert.NotContains(t, updateStateStmt, "set_config")
	assert.Contains(t, updateStateStmt, "in_tx_order")
	assert.NotContains(t, updateStateStmt, "UPDATE projections.current_states SET")
	assert.Equal(t, 1, strings.Count(updateStateStmt, "INSERT INTO projections.current_states"))
}
