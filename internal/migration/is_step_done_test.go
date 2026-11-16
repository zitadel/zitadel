package migration

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type stubEventQuerier struct {
	events []eventstore.Event
	err    error
	query  *eventstore.SearchQueryBuilder
}

func (s *stubEventQuerier) FilterToQueryReducer(_ context.Context, reducer eventstore.QueryReducer) error {
	s.query = reducer.Query()
	if s.err != nil {
		return s.err
	}
	reducer.AppendEvents(s.events...)
	return reducer.Reduce()
}

func setupStepEvent(typ eventstore.EventType, name string) *SetupStep {
	agg := eventstore.NewAggregate(context.Background(), SystemAggregateID, SystemAggregate, "v1")
	return &SetupStep{
		BaseEvent: eventstore.BaseEvent{EventType: typ, Agg: agg},
		Name:      name,
	}
}

func TestIsStepDone(t *testing.T) {
	name := EventstorePositionClockTimestampStep

	t.Run("nil eventstore is not done", func(t *testing.T) {
		done, err := IsStepDone(t.Context(), nil, name)
		assert.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("no events is not done", func(t *testing.T) {
		es := &stubEventQuerier{}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.False(t, done)
		require.NotNil(t, es.query)
		queries := es.query.GetQueries()
		require.Len(t, queries, 1)
		assert.ElementsMatch(t, []eventstore.EventType{StartedType, DoneType, repeatableDoneType, failedType}, queries[0].GetEventTypes())
	})

	t.Run("done event is stored", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(DoneType, name),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("started without done or failed errors", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(StartedType, name),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "started without a done or failed event")
		assert.False(t, done)
	})

	t.Run("cleanup Failed after Started is not done", func(t *testing.T) {
		// `zitadel setup cleanup` pushes Failed so setup can move on.
		// Cursor mode requires a Done event; treating Failed as cursor would
		// brick clusters that aborted a 77 that never executed. Cleanup after
		// a committed 77 Execute is an operator mistake — mark the step Done.
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(StartedType, name),
				setupStepEvent(failedType, name),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("started after failed is stuck", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(StartedType, name),
				setupStepEvent(failedType, name),
				setupStepEvent(StartedType, name),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "started without a done or failed event")
		assert.False(t, done)
	})

	t.Run("started after failed then done is stored", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(StartedType, name),
				setupStepEvent(failedType, name),
				setupStepEvent(StartedType, name),
				setupStepEvent(DoneType, name),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("started after failed then failed is not done", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(StartedType, name),
				setupStepEvent(failedType, name),
				setupStepEvent(StartedType, name),
				setupStepEvent(failedType, name),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("other step name is ignored", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				setupStepEvent(DoneType, "other_step"),
			},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("undefined table is not done", func(t *testing.T) {
		es := &stubEventQuerier{
			err: &pgconn.PgError{Code: "42P01"},
		}
		done, err := IsStepDone(t.Context(), es, name)
		require.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("other errors propagate", func(t *testing.T) {
		want := errors.New("boom")
		es := &stubEventQuerier{err: want}
		done, err := IsStepDone(t.Context(), es, name)
		assert.ErrorIs(t, err, want)
		assert.False(t, done)
	})
}

func TestEventstorePositionClockTimestampStepName(t *testing.T) {
	assert.Equal(t, "77_eventstore_position_clock_timestamp", EventstorePositionClockTimestampStep)
}
