package handler

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/sloggcp"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/pseudo"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// errorContains matches the error text written to the failed events table.
type errorContains string

func (c errorContains) Match(v driver.Value) bool {
	s, ok := v.(string)
	return ok && strings.Contains(s, string(c))
}

func panickingReducer(event eventstore.Event) (*Statement, error) {
	var policies map[string]*struct{ PasswordChange bool } // the policy is missing, like in #13057
	_ = policies["password"].PasswordChange
	return NewNoOpStatement(event), nil
}

func panickingExec(context.Context, Executer, string) error {
	panic("policy is nil")
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

// requireRecoveredError asserts err was created from a recovered panic
// and that its stack trace contains the panicking frame.
func requireRecoveredError(t *testing.T, err error, wantFrame string) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, zerrors.IsRecovered(err), "recovered panics must be recovered errors")

	zErr, ok := errors.AsType[*zerrors.ZitadelError](err)
	require.True(t, ok)
	trace, ok := zErr.StackTrace()
	assert.True(t, ok, "recovered panics must carry a stack trace")
	assert.Contains(t, string(trace), wantFrame)
}

func TestHandler_eventsToStatements_recoversPanic(t *testing.T) {
	tests := []struct {
		name            string
		maxFailureCount uint8
		wantErr         bool
	}{
		{
			name:            "returns the error below the max failure count",
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
				statements, err = h.eventsToStatements(t.Context(), tx, []eventstore.Event{passwordChangedEvent()}, &state{}, decimal.Decimal{})
			})
			if tt.wantErr {
				assert.Empty(t, statements)
				assert.ErrorIs(t, err, &executionError{})
				requireRecoveredError(t, errors.Unwrap(err), "panickingReducer")
			} else {
				assert.NoError(t, err)
				// the skipped event is counted as processed, so the cursor moves past it
				require.Len(t, statements, 1)
				want := NewNoOpStatement(passwordChangedEvent())
				want.offset = 1
				assert.Equal(t, []*Statement{want}, statements)
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
			name:            "returns the error below the max failure count",
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
			statement := NewStatement(passwordChangedEvent(), panickingExec)

			var err error
			require.NotPanics(t, func() {
				err = h.executeStatement(t.Context(), tx, statement)
			})
			if tt.wantErr {
				assert.ErrorIs(t, err, &executionError{})
				requireRecoveredError(t, errors.Unwrap(err), "panickingExec")
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// TestHandler_Trigger_survivesPanic makes sure a panic during the processing of events
// neither terminates the worker nor keeps the instance locked.
func TestHandler_Trigger_survivesPanic(t *testing.T) {
	StartWorkerPool(1)

	// the handler has no database client, processEvents panics as soon as it begins the transaction.
	h := &Handler{projection: panickingProjection()}
	ctx := authz.WithInstanceID(t.Context(), "instance-1")

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
			requireRecoveredError(t, err, "processEvents")
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

// filterStore is an EventStore whose Filter is the given function; a call to any other method panics on the nil interface.
type filterStore struct {
	EventStore
	filter func() ([]eventstore.Event, error)
}

func (s filterStore) Filter(context.Context, *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
	return s.filter()
}

// fieldsEvent satisfies eventstore.FillFieldsEvent for the FieldHandler.
type fieldsEvent struct{ *eventstore.BaseEvent }

func (fieldsEvent) Fields() []*eventstore.FieldOperation { return nil }

// setupProcessEvents fills h so its processEvents runs on sqlmock and records its metrics in the returned mock.
// It fills in place, because a Handler must not be copied.
func setupProcessEvents(t *testing.T, h *Handler, p *projection, es EventStore) (sqlmock.Sqlmock, *metrics.MockMetrics) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	mockMetrics := metrics.NewMockMetrics()
	h.client = &database.DB{DB: db}
	h.projection = p
	h.es = es
	h.bulkLimit = 10
	h.metrics = newProjectionMetrics(t.Context(), mockMetrics)
	return mock, mockMetrics
}

// expectLockedTx expects the beginning of an iteration: the transaction, the acquired lock and an empty current state.
func expectLockedTx(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectQuery("pg_try_advisory_xact_lock").WillReturnRows(sqlmock.NewRows([]string{"locked"}).AddRow(true))
	mock.ExpectQuery("projections.current_states").WillReturnError(sql.ErrNoRows)
}

// TestHandler_processEvents_recoversPanic makes sure a panic inside the transaction rolls back
// and is counted as a failed iteration, whether it happens before or after the metrics defer is registered.
func TestHandler_processEvents_recoversPanic(t *testing.T) {
	tests := []struct {
		name        string
		reduce      Reduce
		filter      func() ([]eventstore.Event, error)
		wantFrame   string
		wantMetrics bool
	}{
		{
			// a reducer returning no statement and no error leaves a nil statement to execute, outside safeExecute
			name:        "nil statement panics after the metrics defer",
			reduce:      func(eventstore.Event) (*Statement, error) { return nil, nil },
			filter:      func() ([]eventstore.Event, error) { return []eventstore.Event{passwordChangedEvent()}, nil },
			wantFrame:   "executeStatement",
			wantMetrics: true,
		},
		{
			name:      "event store panics before the metrics defer",
			reduce:    panickingReducer,
			filter:    func() ([]eventstore.Event, error) { panic("filter failed") },
			wantFrame: "generateStatements",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := panickingProjection()
			p.reducers[0].EventReducers[0].Reduce = tt.reduce
			h := new(Handler)
			mock, mockMetrics := setupProcessEvents(t, h, p, filterStore{filter: tt.filter})
			expectLockedTx(mock)
			mock.ExpectRollback()

			var err error
			require.NotPanics(t, func() {
				_, err = h.processEvents(authz.WithInstanceID(t.Context(), "instance-1"), new(triggerConfig))
			})
			requireRecoveredError(t, err, tt.wantFrame)
			assert.NoError(t, mock.ExpectationsWereMet(), "the iteration must roll back and never commit")
			assert.Empty(t, mockMetrics.GetHistogramValues(ProjectionHandleTimerMetric))
			labels := mockMetrics.GetCounterLabels(ProjectionEventsProcessed)
			if !tt.wantMetrics {
				assert.Empty(t, labels, "the metrics defer was not registered yet")
				return
			}
			if assert.Len(t, labels, 1) {
				assert.False(t, labels[0][SuccessLabel].AsBool(), "a panicked iteration must count as failed")
			}
		})
	}
}

func TestFieldHandler_processEvents_recoversPanic(t *testing.T) {
	// the nil embedded EventStore panics in FillFields
	fh := new(FieldHandler)
	mock, _ := setupProcessEvents(t, &fh.Handler, panickingProjection(), filterStore{filter: func() ([]eventstore.Event, error) {
		return []eventstore.Event{fieldsEvent{passwordChangedEvent()}}, nil
	}})
	expectLockedTx(mock)
	mock.ExpectRollback()

	var err error
	require.NotPanics(t, func() {
		_, err = fh.processEvents(authz.WithInstanceID(t.Context(), "instance-1"), new(triggerConfig))
	})
	requireRecoveredError(t, err, "FieldHandler).processEvents")
	assert.NoError(t, mock.ExpectationsWereMet(), "the iteration must roll back and never commit")
}

// TestHandler_triggerInstances_skipsRecoveredPanic makes sure an instance which panics
// is not retried and does not starve the instances after it.
func TestHandler_triggerInstances_skipsRecoveredPanic(t *testing.T) {
	StartWorkerPool(1)

	tests := []struct {
		name string
		// errorsBeforePanic is the number of ordinary errors instance-1 returns before it panics
		errorsBeforePanic int
	}{
		{
			name: "panics at once",
		},
		{
			name:              "panics after an ordinary error",
			errorsBeforePanic: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := make(map[string]int)
			h := new(Handler)
			mock, _ := setupProcessEvents(t, h, panickingProjection(), nil)
			h.retryFailedAfter = time.Millisecond
			// the scheduled event carries the instance, so the reducer can fail for the first instance only
			h.triggerWithoutEvents = func(event eventstore.Event) (*Statement, error) {
				instanceID := event.(*pseudo.ScheduledEvent).InstanceIDs[0]
				calls[instanceID]++
				if instanceID != "instance-1" {
					return NewNoOpStatement(event), nil
				}
				if calls[instanceID] <= tt.errorsBeforePanic {
					return nil, errors.New("instance-1 failed")
				}
				panic("instance-1 is broken")
			}
			// instance-1 fails and finally panics, each attempt rolls back
			for range tt.errorsBeforePanic + 1 {
				expectLockedTx(mock)
				mock.ExpectRollback()
			}
			// instance-2 stores its state and commits
			expectLockedTx(mock)
			mock.ExpectExec("INSERT INTO projections.current_states").WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			done := make(chan struct{})
			go func() {
				defer close(done)
				h.triggerInstances(t.Context(), []string{"instance-1", "instance-2"})
			}()
			select {
			case <-done:
			case <-time.After(10 * time.Second):
				t.Fatal("triggerInstances did not return, the panicking instance is retried forever")
			}
			assert.Equal(t, map[string]int{"instance-1": tt.errorsBeforePanic + 1, "instance-2": 1}, calls)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// TestHandler_eventsToStatements_logLevel makes sure a recovered panic is logged at Alert
// and an ordinary reduce error at Error.
func TestHandler_eventsToStatements_logLevel(t *testing.T) {
	failingProjection := panickingProjection()
	failingProjection.reducers[0].EventReducers[0].Reduce = func(eventstore.Event) (*Statement, error) {
		return nil, errors.New("reduce error")
	}
	tests := []struct {
		name       string
		projection *projection
		errText    string
		wantLevel  slog.Level
	}{
		{
			name:       "recovered panic is logged at alert",
			projection: panickingProjection(),
			errText:    "nil pointer dereference",
			wantLevel:  sloggcp.LevelAlert,
		},
		{
			name:       "reduce error is logged at error",
			projection: failingProjection,
			errText:    "reduce error",
			wantLevel:  slog.LevelError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			ctx := logging.ToCtx(t.Context(), slog.New(slog.NewJSONHandler(&buf, nil)))
			tx, mock := beginMockTx(t)
			expectFailedEvent(mock, tt.errText)
			h := &Handler{projection: tt.projection, maxFailureCount: 1}

			_, err := h.eventsToStatements(ctx, tx, []eventstore.Event{passwordChangedEvent()}, &state{}, decimal.Decimal{})
			require.NoError(t, err)

			var record struct {
				Level string `json:"level"`
				Msg   string `json:"msg"`
			}
			require.NoError(t, json.NewDecoder(&buf).Decode(&record))
			assert.Equal(t, "reduce failed", record.Msg)
			assert.Equal(t, tt.wantLevel.String(), record.Level)
		})
	}
}
