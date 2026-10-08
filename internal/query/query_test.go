package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
)

type expect func(mockRepository *mock.MockRepository)

func expectEventstore(expects ...expect) func(*testing.T) *eventstore.Eventstore {
	return func(t *testing.T) *eventstore.Eventstore {
		m := mock.NewRepo(t)
		for _, e := range expects {
			e(m)
		}
		es := eventstore.NewEventstore(
			&eventstore.Config{
				Querier: m.MockQuerier,
				Pusher:  m.MockPusher,
			},
		)
		return es
	}
}

func expectFilter(events ...eventstore.Event) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	}
}
func expectFilterError(err error) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEventsError(err)
	}
}

func eventFromEventPusher(event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		InstanceID:    event.Aggregate().InstanceID,
		ID:            "",
		Seq:           0,
		CreationDate:  time.Time{},
		Typ:           event.Type(),
		Data:          data,
		EditorUser:    event.Creator(),
		Version:       event.Aggregate().Version,
		AggregateID:   event.Aggregate().ID,
		AggregateType: event.Aggregate().Type,
		ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
		Constraints:   event.UniqueConstraints(),
	}
}

// testTriggerKey identifies the projection name carried by the contexts
// returned from testTriggerer.Trigger.
type testTriggerKey struct{}

// testTriggerer is a [triggerer] which records how it was called
// and delegates to fn.
type testTriggerer struct {
	name string
	fn   func(ctx context.Context) (context.Context, error)

	calls   atomic.Int32
	optsLen atomic.Int32
}

func (t *testTriggerer) ProjectionName() string {
	return t.name
}

func (t *testTriggerer) Trigger(ctx context.Context, opts ...handler.TriggerOpt) (context.Context, error) {
	t.calls.Add(1)
	t.optsLen.Store(int32(len(opts)))
	return t.fn(ctx)
}

// newTestTriggerer returns a triggerer which returns a fresh context tagged
// with its own name, so the caller can tell which trigger produced it.
// The passed context is deliberately not used as parent: triggerBatch wraps
// it in a tracing span, so it is never identical to the one passed to triggerBatch.
func newTestTriggerer(name string, err error) *testTriggerer {
	return &testTriggerer{
		name: name,
		fn: func(context.Context) (context.Context, error) {
			return context.WithValue(context.Background(), testTriggerKey{}, name), err
		},
	}
}

func Test_triggerBatch(t *testing.T) {
	t.Run("no handlers", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), testTriggerKey{}, "input")

		gotCtx, err := triggerBatch(ctx)
		assert.NoError(t, err)
		assert.Equal(t, ctx, gotCtx)
	})

	t.Run("single handler", func(t *testing.T) {
		h := newTestTriggerer("p1", nil)

		gotCtx, err := triggerBatch(context.Background(), h)
		require.NoError(t, err)
		assert.Equal(t, "p1", gotCtx.Value(testTriggerKey{}))
		assert.Equal(t, int32(1), h.calls.Load())
		// handler.WithAwaitRunning() must be passed to Trigger.
		assert.Equal(t, int32(1), h.optsLen.Load())
	})

	t.Run("all succeed", func(t *testing.T) {
		handlers := []*testTriggerer{
			newTestTriggerer("p1", nil),
			newTestTriggerer("p2", nil),
			newTestTriggerer("p3", nil),
		}

		gotCtx, err := triggerBatch(context.Background(), handlers[0], handlers[1], handlers[2])
		require.NoError(t, err)
		// Which context is returned depends on which Trigger finishes last.
		assert.Contains(t, []any{"p1", "p2", "p3"}, gotCtx.Value(testTriggerKey{}))

		for _, h := range handlers {
			assert.Equal(t, int32(1), h.calls.Load(), h.name)
			assert.Equal(t, int32(1), h.optsLen.Load(), h.name)
		}
	})

	t.Run("errors joined", func(t *testing.T) {
		errP1 := errors.New("p1 failed")
		errP3 := errors.New("p3 failed")

		gotCtx, err := triggerBatch(context.Background(),
			newTestTriggerer("p1", errP1),
			newTestTriggerer("p2", nil),
			newTestTriggerer("p3", errP3),
		)
		require.Error(t, err)
		assert.ErrorIs(t, err, errP1)
		assert.ErrorIs(t, err, errP3)
		assert.NotNil(t, gotCtx)

		// The nil error of p2 is dropped by errors.Join.
		joined, ok := err.(interface{ Unwrap() []error })
		require.True(t, ok, "expected a joined error, got %T", err)
		assert.Len(t, joined.Unwrap(), 2)
	})

	t.Run("runs concurrently", func(t *testing.T) {
		const amount = 3

		arrived := make(chan struct{}, amount)
		release := make(chan struct{})
		go func() {
			for range amount {
				<-arrived
			}
			close(release)
		}()

		// Every Trigger blocks until all of them have been entered.
		// A sequential implementation never gets past the first one.
		fn := func(ctx context.Context) (context.Context, error) {
			arrived <- struct{}{}
			select {
			case <-release:
				return ctx, nil
			case <-time.After(5 * time.Second):
				return ctx, errors.New("handlers did not run concurrently")
			}
		}

		handlers := make([]triggerer, amount)
		for i := range handlers {
			handlers[i] = &testTriggerer{name: fmt.Sprintf("p%d", i), fn: fn}
		}

		_, err := triggerBatch(context.Background(), handlers...)
		assert.NoError(t, err)
	})
}

func Test_cleanStaticQueries(t *testing.T) {
	query := `select
	foo,
	bar
from table;`
	want := "select foo, bar from table;"
	cleanStaticQueries(&query)
	assert.Equal(t, want, query)
}
