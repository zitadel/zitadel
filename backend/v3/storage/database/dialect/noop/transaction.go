package noopdb

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Transaction struct{}

// Begin implements [database.Transaction].
func (n *Transaction) Begin(ctx context.Context) (database.Transaction, error) {
	return new(Transaction), nil
}

// Commit implements [database.Transaction].
func (n *Transaction) Commit(ctx context.Context) error {
	return nil
}

// End implements [database.Transaction].
func (n *Transaction) End(ctx context.Context, err error) error {
	return err
}

// Exec implements [database.Transaction].
func (n *Transaction) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	return 0, nil
}

// Query implements [database.Transaction].
func (n *Transaction) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return new(rows), nil
}

// QueryRow implements [database.Transaction].
func (n *Transaction) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	return new(row)
}

// Rollback implements [database.Transaction].
func (n *Transaction) Rollback(ctx context.Context) error {
	return nil
}

var _ database.Transaction = (*Transaction)(nil)
