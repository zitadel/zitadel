package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// errorContains matches the error text written to the failed events table.
type errorContains string

func (c errorContains) Match(v driver.Value) bool {
	s, ok := v.(string)
	return ok && strings.Contains(s, string(c))
}

func panickingReducer(event eventstore.Event) (*Statement, error) {
	var policy *struct{ PasswordChange bool }
	_ = policy.PasswordChange // nil pointer dereference
	return NewNoOpStatement(event), nil
}

func panickingProjection() *projection {
	return &projection{
		name: "projections.notifications",
		reducers: []AggregateReducer{{
			Aggregate: "user",
			EventReducers: []EventReducer{{
				Event:  "user.human.password.changed",
				Reduce: panickingReducer,
			}},
		}},
	}
}

func passwordChangedEvent() *eventstore.BaseEvent {
	return &eventstore.BaseEvent{
		Agg:       &eventstore.Aggregate{ID: "user-1", Type: "user", InstanceID: "instance-1"},
		EventType: "user.human.password.changed",
		Seq:       7,
	}
}

func requireRecoveredError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, zerrors.IsInternal(err), "recovered panics must be internal errors")
	assert.Contains(t, err.Error(), zerrors.IDRecover)
}

func TestHandler_reduce_recoversPanic(t *testing.T) {
	h := &Handler{projection: panickingProjection()}

	var (
		statement *Statement
		err       error
	)
	require.NotPanics(t, func() {
		statement, err = h.reduce(passwordChangedEvent())
	})
	assert.Nil(t, statement)
	requireRecoveredError(t, err)
	assert.Contains(t, err.Error(), "nil pointer dereference")
}

func TestHandler_eventsToStatements_recoversPanic(t *testing.T) {
	tests := []struct {
		name            string
		maxFailureCount uint8
		wantErr         bool
	}{
		{
			name:            "stops until the max failure count is reached",
			maxFailureCount: 5,
			wantErr:         true,
		},
		{
			name:            "skips the event when the max failure count is reached",
			maxFailureCount: 1,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, mock := beginMockTx(t)
			expectFailedEvent(mock, "nil pointer dereference")

			h := &Handler{projection: panickingProjection(), maxFailureCount: tt.maxFailureCount}

			var (
				statements []*Statement
				err        error
			)
			require.NotPanics(t, func() {
				statements, err = h.eventsToStatements(context.Background(), tx, []eventstore.Event{passwordChangedEvent()})
			})
			assert.Empty(t, statements)
			if tt.wantErr {
				assert.ErrorIs(t, err, &executionError{})
				requireRecoveredError(t, errors.Unwrap(err))
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestHandler_executeStatement_recoversPanic(t *testing.T) {
	tests := []struct {
		name            string
		maxFailureCount uint8
		wantErr         bool
	}{
		{
			name:            "stops until the max failure count is reached",
			maxFailureCount: 5,
			wantErr:         true,
		},
		{
			name:            "skips the statement when the max failure count is reached",
			maxFailureCount: 1,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, mock := beginMockTx(t)
			mock.ExpectExec("SAVEPOINT exec_stmt").WillReturnResult(sqlmock.NewResult(0, 0))
			mock.ExpectExec("ROLLBACK TO SAVEPOINT exec_stmt").WillReturnResult(sqlmock.NewResult(0, 0))
			expectFailedEvent(mock, "policy is nil")

			h := &Handler{projection: panickingProjection(), maxFailureCount: tt.maxFailureCount}
			statement := NewStatement(passwordChangedEvent(), func(context.Context, Executer, string) error {
				panic("policy is nil")
			})

			var err error
			require.NotPanics(t, func() {
				err = h.executeStatement(context.Background(), tx, statement)
			})
			if tt.wantErr {
				assert.ErrorIs(t, err, &executionError{})
				requireRecoveredError(t, errors.Unwrap(err))
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_safeProcessEvents(t *testing.T) {
	t.Run("recovers a panic", func(t *testing.T) {
		var (
			additionalIteration bool
			err                 error
		)
		require.NotPanics(t, func() {
			additionalIteration, err = safeProcessEvents(context.Background(), new(triggerConfig), func(context.Context, *triggerConfig) (bool, error) {
				panic(errors.New("boom"))
			})
		})
		assert.False(t, additionalIteration)
		requireRecoveredError(t, err)
		assert.Contains(t, err.Error(), "boom")
	})

	t.Run("passes the result through", func(t *testing.T) {
		wantErr := errors.New("failed")
		additionalIteration, err := safeProcessEvents(context.Background(), new(triggerConfig), func(context.Context, *triggerConfig) (bool, error) {
			return true, wantErr
		})
		assert.True(t, additionalIteration)
		assert.ErrorIs(t, err, wantErr)
	})
}

// TestHandler_Trigger_survivesPanic makes sure a panic during the processing of events
// neither terminates the worker nor keeps the instance locked.
func TestHandler_Trigger_survivesPanic(t *testing.T) {
	StartWorkerPool(1)

	// the handler has no database client, processEvents panics as soon as it begins the transaction.
	h := &Handler{projection: panickingProjection()}
	ctx := authz.WithInstanceID(context.Background(), "instance-1")

	// the second call would block forever if the first call killed the worker
	// and would return without an error if the first call did not release the instance lock.
	for range 2 {
		done := make(chan error, 1)
		go func() {
			_, err := h.Trigger(ctx)
			done <- err
		}()

		select {
		case err := <-done:
			requireRecoveredError(t, err)
		case <-time.After(10 * time.Second):
			t.Fatal("trigger did not return, the worker did not survive the panic")
		}
	}
}

func beginMockTx(t *testing.T) (*sql.Tx, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	mock.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)
	return tx, mock
}

// expectFailedEvent expects the failure count of the event to be read and increased to 1.
func expectFailedEvent(mock sqlmock.Sqlmock, errText string) {
	mock.ExpectQuery("SELECT COALESCE").
		WithArgs("projections.notifications", "instance-1", "user", "user-1", uint64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"failure_count"}).AddRow(0))
	mock.ExpectExec("INSERT INTO projections.failed_events2").
		WithArgs("projections.notifications", "instance-1", "user", "user-1", sqlmock.AnyArg(), uint64(7), uint8(1), errorContains(errText)).
		WillReturnResult(sqlmock.NewResult(0, 1))
}
