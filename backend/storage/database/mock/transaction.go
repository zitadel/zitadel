package mock

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type Transaction struct {
	t *testing.T

	committed  bool
	rolledBack bool

	expectations []expecter
}

func NewTransaction(t *testing.T, opts ...TransactionOption) *Transaction {
	tx := &Transaction{t: t}
	for _, opt := range opts {
		opt(tx)
	}
	return tx
}

func (tx *Transaction) nextExpecter() expecter {
	if len(tx.expectations) == 0 {
		tx.t.Error("no more expectations on transaction")
		tx.t.FailNow()
	}

	e := tx.expectations[0]
	tx.expectations = tx.expectations[1:]
	return e
}

type TransactionOption func(tx *Transaction)

type expecter interface {
	assertArgs(ctx context.Context, stmt string, args ...any)
}

func ExpectExec(stmt string, args ...any) TransactionOption {
	return func(tx *Transaction) {
		tx.expectations = append(tx.expectations, &expectation[struct{}]{
			t:            tx.t,
			expectedStmt: stmt,
			expectedArgs: args,
		})
	}
}

func ExpectQuery(res database.Rows, stmt string, args ...any) TransactionOption {
	return func(tx *Transaction) {
		tx.expectations = append(tx.expectations, &expectation[database.Rows]{
			t:            tx.t,
			expectedStmt: stmt,
			expectedArgs: args,
			result:       res,
		})
	}
}

func ExpectQueryRow(res database.Row, stmt string, args ...any) TransactionOption {
	return func(tx *Transaction) {
		tx.expectations = append(tx.expectations, &expectation[database.Row]{
			t:            tx.t,
			expectedStmt: stmt,
			expectedArgs: args,
			result:       res,
		})
	}
}

type expectation[R any] struct {
	t *testing.T

	expectedStmt string
	expectedArgs []any

	result R
}

func (e *expectation[R]) assertArgs(ctx context.Context, stmt string, args ...any) {
	e.t.Helper()
	assert.Equal(e.t, e.expectedStmt, stmt)
	assert.Equal(e.t, e.expectedArgs, args)
}

// Commit implements [database.Transaction].
func (t *Transaction) Commit(ctx context.Context) error {
	if t.hasEnded() {
		return errors.New("transaction already committed or rolled back")
	}
	t.committed = true
	return nil
}

// End implements [database.Transaction].
func (tx *Transaction) End(ctx context.Context, err error) error {
	if tx.hasEnded() {
		return errors.New("transaction already committed or rolled back")
	}
	if err != nil {
		return tx.Rollback(ctx)
	}
	return tx.Commit(ctx)
}

// Exec implements [database.Transaction].
func (tx *Transaction) Exec(ctx context.Context, stmt string, args ...any) error {
	tx.nextExpecter().assertArgs(ctx, stmt, args...)

	return nil
}

// Query implements [database.Transaction].
func (tx *Transaction) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	e := tx.nextExpecter()
	e.assertArgs(ctx, stmt, args...)
	return e.(*expectation[database.Rows]).result, nil
}

// QueryRow implements [database.Transaction].
func (tx *Transaction) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	e := tx.nextExpecter()
	e.assertArgs(ctx, stmt, args...)
	return e.(*expectation[database.Row]).result
}

// Rollback implements [database.Transaction].
func (tx *Transaction) Rollback(ctx context.Context) error {
	if tx.hasEnded() {
		return errors.New("transaction already committed or rolled back")
	}
	tx.rolledBack = true
	return nil
}

var _ database.Transaction = (*Transaction)(nil)

func (tx *Transaction) hasEnded() bool {
	return tx.committed || tx.rolledBack
}
