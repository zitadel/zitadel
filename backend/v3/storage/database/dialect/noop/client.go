package noopdb

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Client struct{}

// Begin implements [database.Client].
func (n *Client) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	return new(Transaction), nil
}

// Exec implements [database.Client].
func (n *Client) Exec(ctx context.Context, stmt string, args ...any) (int64, error) {
	return 0, nil
}

// Migrate implements [database.Client].
func (n *Client) Migrate(ctx context.Context) error {
	return nil
}

// Query implements [database.Client].
func (n *Client) Query(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return new(rows), nil
}

// QueryRow implements [database.Client].
func (n *Client) QueryRow(ctx context.Context, stmt string, args ...any) database.Row {
	return new(row)
}

// Release implements [database.Client].
func (n *Client) Release(ctx context.Context) error {
	return nil
}

var _ database.Client = (*Client)(nil)
