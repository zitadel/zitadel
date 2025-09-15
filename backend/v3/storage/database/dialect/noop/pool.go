package noopdb

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Pool struct{}

// Acquire implements [database.Pool].
func (n *Pool) Acquire(ctx context.Context) (database.Client, error) {
	return new(Client), nil
}

// Begin implements [database.Pool].
func (n *Pool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	return new(Transaction), nil
}

// Close implements [database.Pool].
func (n *Pool) Close(ctx context.Context) error {
	return nil
}

// Exec implements [database.Pool].
func (n *Pool) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	return 0, nil
}

// Migrate implements [database.Pool].
func (n *Pool) Migrate(ctx context.Context) error {
	return nil
}

// Query implements [database.Pool].
func (n *Pool) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return new(rows), nil
}

// QueryRow implements [database.Pool].
func (n *Pool) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	return new(row)
}

var _ database.Pool = (*Pool)(nil)
