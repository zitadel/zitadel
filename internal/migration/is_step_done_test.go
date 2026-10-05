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

func TestIsStepDone(t *testing.T) {
	t.Run("nil eventstore is not done", func(t *testing.T) {
		done, err := IsStepDone(t.Context(), nil, EventstorePositionClockTimestampStep)
		assert.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("no events is not done", func(t *testing.T) {
		es := &stubEventQuerier{}
		done, err := IsStepDone(t.Context(), es, EventstorePositionClockTimestampStep)
		require.NoError(t, err)
		assert.False(t, done)
		require.NotNil(t, es.query)
		queries := es.query.GetQueries()
		require.Len(t, queries, 1)
		assert.Equal(t, []eventstore.EventType{DoneType}, queries[0].GetEventTypes())
	})

	t.Run("done event is stored", func(t *testing.T) {
		es := &stubEventQuerier{
			events: []eventstore.Event{
				&eventstore.BaseEvent{EventType: DoneType},
			},
		}
		done, err := IsStepDone(t.Context(), es, EventstorePositionClockTimestampStep)
		require.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("undefined table is not done", func(t *testing.T) {
		es := &stubEventQuerier{
			err: &pgconn.PgError{Code: "42P01"},
		}
		done, err := IsStepDone(t.Context(), es, EventstorePositionClockTimestampStep)
		require.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("other errors propagate", func(t *testing.T) {
		want := errors.New("boom")
		es := &stubEventQuerier{err: want}
		done, err := IsStepDone(t.Context(), es, EventstorePositionClockTimestampStep)
		assert.ErrorIs(t, err, want)
		assert.False(t, done)
	})
}
