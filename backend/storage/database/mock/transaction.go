package mock

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type Transaction struct {
	committed  bool
	rolledBack bool
}

func NewTransaction() *Transaction {
	return new(Transaction)
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
func (t *Transaction) End(ctx context.Context, err error) error {
	if t.hasEnded() {
		return errors.New("transaction already committed or rolled back")
	}
	if err != nil {
		return t.Rollback(ctx)
	}
	return t.Commit(ctx)
}

// Exec implements [database.Transaction].
func (t *Transaction) Exec(ctx context.Context, sql string, args ...any) error {
	return nil
}

// Query implements [database.Transaction].
func (t *Transaction) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	return nil, nil
}

// QueryRow implements [database.Transaction].
func (t *Transaction) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return nil
}

// Rollback implements [database.Transaction].
func (t *Transaction) Rollback(ctx context.Context) error {
	if t.hasEnded() {
		return errors.New("transaction already committed or rolled back")
	}
	t.rolledBack = true
	return nil
}

var _ database.Transaction = (*Transaction)(nil)

func (t *Transaction) hasEnded() bool {
	return t.committed || t.rolledBack
}
