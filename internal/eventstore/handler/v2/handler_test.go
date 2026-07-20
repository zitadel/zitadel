package handler

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// TestHandler_executeStatement_rollbackSurvivesContextCancellation reproduces
// the SAVEPOINT leak fixed alongside this test. With the broken code the
// `ROLLBACK TO SAVEPOINT exec_stmt` call inherited the request ctx, so when
// the ctx had been cancelled by the time the projection statement returned
// an error the ExecContext call short-circuited with context.Canceled
// BEFORE issuing the SQL. The connection went back to the pool with
// BEGIN + SAVEPOINT still alive server-side and was killed by Postgres
// after idle_in_transaction_session_timeout (SQLSTATE 25P03). The fix
// detaches the rollback from the parent ctx with context.WithoutCancel;
// this test asserts the rollback actually reaches the database.
func TestHandler_executeStatement_rollbackSurvivesContextCancellation(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SAVEPOINT exec_stmt").
		WillReturnResult(sqlmock.NewResult(0, 0))
	// The smoking-gun expectation: this MUST fire even though the
	// projection cancels ctx before returning. Without the fix this
	// expectation is unmet and ExpectationsWereMet() reports the leak.
	mock.ExpectExec("ROLLBACK TO SAVEPOINT exec_stmt").
		WillReturnResult(sqlmock.NewResult(0, 0))
	// handleFailedStmt() does an embedded SELECT on the failure-count
	// table; we make it fail so the function short-circuits without a
	// follow-up UPDATE and executeStatement returns &executionError{}.
	// This keeps the test focused on the rollback signal and avoids
	// duplicating the failed_event SQL contract in the assertion.
	// Match the actual failure_event_get_count.sql shape rather than
	// any query — keeps the test strict against unrelated statements.
	mock.ExpectQuery(`(?i)SELECT.+failed_events`).WillReturnError(sql.ErrConnDone)
	// The deferred tx.Rollback() in t.Cleanup() below issues a ROLLBACK
	// that sqlmock would otherwise flag as unexpected.
	mock.ExpectRollback()

	tx, err := db.Begin()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	h := &Handler{
		projection:      &projection{name: "test_projection"},
		maxFailureCount: 5,
	}

	stmt := &Statement{
		Aggregate:    &eventstore.Aggregate{InstanceID: "i", ID: "a", Type: "t"},
		Sequence:     1,
		Position:     decimal.NewFromInt(1),
		CreationDate: time.Now(),
		Execute: func(ctx context.Context, _ Executer, _ string) error {
			// Mimic the real-world race: the request ctx gets cancelled
			// (h.txDuration timeout, caller cancel, pod-shutdown SIGTERM,
			// gRPC deadline) while the projection statement is in flight.
			// The projection then returns an error.
			cancel()
			return errors.New("projection logic failed mid-flight")
		},
	}

	execErr := h.executeStatement(ctx, tx, stmt)
	// With handleFailedStmt forced to short-circuit, executeStatement
	// returns the wrapped executionError. The exact error type is
	// secondary to the real assertion: ExpectationsWereMet().
	assert.Error(t, execErr)

	// Close the outer tx so the ExpectRollback expectation can settle
	// before we verify all expectations. This mirrors what the
	// production caller (processEvents) does in its deferred
	// tx.Rollback() / tx.Commit() chain.
	require.NoError(t, tx.Rollback(), "outer tx.Rollback() must succeed")

	require.NoError(t, mock.ExpectationsWereMet(),
		"ROLLBACK TO SAVEPOINT must reach the database even when the request context was cancelled during statement execution")
}
