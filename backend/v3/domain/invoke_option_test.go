package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ database.QueryExecutor = (*nonBeginnerDB)(nil)
var _ database.Transaction = (*transactionDB)(nil)
var _ database.Beginner = (*beginnerDB)(nil)
var _ database.QueryExecutor = (*beginnerDB)(nil)

type beginnerDB struct {
	errToReturn error
	errOnBegin  error
}

func (n *beginnerDB) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	if n.errToReturn != nil {
		return nil, n.errToReturn
	}
	if n.errOnBegin != nil {
		return nil, n.errOnBegin
	}
	return &transactionDB{}, nil
}

// Exec implements [database.QueryExecutor].
func (n *beginnerDB) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	return 0, n.errToReturn
}

// Query implements [database.QueryExecutor].
func (n *beginnerDB) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return nil, n.errToReturn
}

// QueryRow implements [database.QueryExecutor].
func (n *beginnerDB) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	return nil
}

type transactionDB struct{}

// Exec implements [database.QueryExecutor].
func (n *transactionDB) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	return 0, nil
}

// Query implements [database.QueryExecutor].
func (n *transactionDB) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return nil, nil
}

// QueryRow implements [database.QueryExecutor].
func (n *transactionDB) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	return nil
}

// Begin implements [database.Transaction].
func (n *transactionDB) Begin(ctx context.Context) (database.Transaction, error) {
	return n, nil
}

// Commit implements [database.Transaction].
func (n *transactionDB) Commit(ctx context.Context) error {
	return nil
}

// End implements [database.Transaction].
func (n *transactionDB) End(ctx context.Context, err error) error {
	return nil
}

// Rollback implements [database.Transaction].
func (n *transactionDB) Rollback(ctx context.Context) error {
	return nil
}

type nonBeginnerDB struct{}

// Exec implements [database.QueryExecutor].
func (n *nonBeginnerDB) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	return 0, nil
}

// Query implements [database.QueryExecutor].
func (n *nonBeginnerDB) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return nil, nil
}

// QueryRow implements [database.QueryExecutor].
func (n *nonBeginnerDB) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	return nil
}

func TestStartTransactionFromDB(t *testing.T) {
	t.Parallel()

	txErr := errors.New("tx error")

	tt := []struct {
		testName                 string
		inputDB                  database.QueryExecutor
		expectedError            error
		expectedValidTransaction bool
	}{
		{
			testName:      "when input DB doesn't implement database.Beginner should return internal error",
			inputDB:       &nonBeginnerDB{},
			expectedError: zerrors.CreateZitadelError(zerrors.KindInternal, nil, "DOM-LqxZbk", "database doesn't implement database.Beginner", 1),
		},
		{
			testName:      "when transaction Begin fails should return internal error",
			inputDB:       &beginnerDB{errOnBegin: txErr},
			expectedError: zerrors.CreateZitadelError(zerrors.KindInternal, txErr, "DOM-sAAd3V", "failed starting transaction", 1),
		},
		{
			testName:                 "when transaction Begin succeeds should return transaction",
			inputDB:                  &beginnerDB{},
			expectedValidTransaction: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			invokeOpts := InvokeOpts{}

			// Test
			tx, err := invokeOpts.StartTransactionFromDB(t.Context(), tc.inputDB, nil)

			// Verify
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedValidTransaction, tx != nil)
		})
	}
}
