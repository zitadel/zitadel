package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

// A notifier must be able to await other projections even when there is only
// one worker. Test both statement execution and reducers which create their own
// context, as the legacy OTP notifier does.
func TestHandler_Trigger_outsideWorkerPoolAwaitsProjections(t *testing.T) {
	StartWorkerPool(1)

	for _, duringReduction := range []bool{false, true} {
		t.Run(fmt.Sprintf("during reduction=%t", duringReduction), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(authz.WithInstanceID(t.Context(), "instance-1"), 5*time.Second)
			defer cancel()

			var projected atomic.Int32
			children := make([]*Handler, 2)
			childMocks := make([]sqlmock.Sqlmock, len(children))
			for i := range children {
				children[i] = new(Handler)
				childMocks[i], _ = setupProcessEvents(t, children[i], &projection{name: fmt.Sprintf("projections.child_%d", i)}, filterStore{filter: func() ([]eventstore.Event, error) {
					projected.Add(1)
					return nil, nil
				}})
				childMocks[i].ExpectBegin()
				childMocks[i].ExpectExec("SELECT pg_advisory_xact_lock").WillReturnResult(sqlmock.NewResult(0, 1))
				childMocks[i].ExpectQuery("projections.current_states").WillReturnError(sql.ErrNoRows)
				childMocks[i].ExpectExec("INSERT INTO projections.current_states").WillReturnResult(sqlmock.NewResult(0, 1))
				childMocks[i].ExpectCommit()
			}

			event := passwordChangedEvent()
			var observedFreshData atomic.Bool
			p := panickingProjection()
			p.reducers[0].EventReducers[0].Reduce = func(event eventstore.Event) (*Statement, error) {
				if duringReduction {
					// Reducers do not receive the processing context. Use the
					// original context, without relying on worker context values.
					if err := awaitChildProjections(ctx, children); err != nil {
						return nil, err
					}
				}
				return NewStatement(event, func(statementCtx context.Context, _ Executer, _ string) error {
					if !duringReduction {
						if err := awaitChildProjections(statementCtx, children); err != nil {
							return err
						}
					}
					observedFreshData.Store(projected.Load() == int32(len(children)))
					return nil
				}), nil
			}
			h := NewHandler(ctx, &Config{RunOutsideWorkerPool: true}, p)
			mock, _ := setupProcessEvents(t, h, p, filterStore{filter: func() ([]eventstore.Event, error) {
				return []eventstore.Event{event}, nil
			}})
			expectLockedTx(mock)
			mock.ExpectExec("SAVEPOINT exec_stmt").WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("INSERT INTO projections.current_states").WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			done := make(chan error, 1)
			go func() {
				_, err := h.Trigger(ctx)
				done <- err
			}()
			select {
			case err := <-done:
				require.NoError(t, err)
			case <-ctx.Done():
				t.Fatal("notifier did not finish with a single projection worker")
			}
			assert.True(t, observedFreshData.Load(), "the notifier must await updated user data")
			assert.NoError(t, mock.ExpectationsWereMet())
			for _, childMock := range childMocks {
				assert.NoError(t, childMock.ExpectationsWereMet())
			}
		})
	}
}

func awaitChildProjections(ctx context.Context, children []*Handler) error {
	results := make(chan error, len(children))
	for _, child := range children {
		go func() {
			_, err := child.Trigger(ctx, WithAwaitRunning())
			results <- err
		}()
	}
	errs := make([]error, len(children))
	for i := range errs {
		errs[i] = <-results
	}
	return errors.Join(errs...)
}

func TestTrigger_canceledWhileWaitingForWorkerReleasesInstance(t *testing.T) {
	StartWorkerPool(1)

	for _, fieldHandler := range []bool{false, true} {
		t.Run(fmt.Sprintf("field handler=%t", fieldHandler), func(t *testing.T) {
			releaseWorker := occupyProjectionWorker(t)
			defer releaseWorker()

			h := &Handler{projection: &projection{name: "projections.waiting"}}
			trigger := func(ctx context.Context) error {
				_, err := h.Trigger(ctx, WithAwaitRunning())
				return err
			}
			if fieldHandler {
				fh := &FieldHandler{}
				fh.projection = h.projection
				h = &fh.Handler
				trigger = func(ctx context.Context) error {
					return fh.Trigger(ctx, WithAwaitRunning())
				}
			}

			ctx := authz.WithInstanceID(t.Context(), "instance-1")
			triggerCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- trigger(triggerCtx) }()

			// Cancellation must happen after the instance lock has been
			// acquired, while the only worker is still occupied.
			require.Eventually(t, func() bool {
				lock, ok := h.triggeredInstancesSync.Load("instance-1")
				return ok && len(lock.(chan bool)) == 1
			}, 5*time.Second, time.Millisecond)
			cancel()
			select {
			case err := <-done:
				require.ErrorIs(t, err, context.Canceled)
			case <-time.After(5 * time.Second):
				t.Fatal("canceled trigger kept waiting for the worker")
			}

			unlock := h.lockInstance(ctx, new(triggerConfig))
			require.NotNil(t, unlock, "a canceled queue submission must release its instance lock")
			unlock()
		})
	}
}

// occupyProjectionWorker blocks the shared single worker until the returned
// function is called, and waits for it to become available again.
func occupyProjectionWorker(t *testing.T) func() {
	t.Helper()
	release := make(chan struct{})
	finished := make(chan struct{})
	select {
	case queue <- func() {
		<-release
		close(finished)
	}:
	case <-time.After(5 * time.Second):
		t.Fatal("projection worker is unavailable")
	}
	return func() {
		close(release)
		select {
		case <-finished:
		case <-time.After(5 * time.Second):
			t.Error("projection worker did not finish")
		}
	}
}

func TestHandler_Trigger_outsideWorkerPoolRecoversPanic(t *testing.T) {
	ctx := authz.WithInstanceID(t.Context(), "instance-1")
	h := NewHandler(ctx, &Config{RunOutsideWorkerPool: true}, panickingProjection())
	// With no database client, processEvents panics before beginning a transaction.
	for range 2 {
		_, err := h.Trigger(ctx)
		requireRecoveredError(t, err, "processEvents")
	}
}
