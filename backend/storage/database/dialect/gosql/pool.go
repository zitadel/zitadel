package gosql

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type sqlPool struct{ *sql.DB }

var _ database.Pool = (*sqlPool)(nil)

// Acquire implements [database.Pool].
func (c *sqlPool) Acquire(ctx context.Context) (database.Client, error) {
	conn, err := c.DB.Conn(ctx)
	if err != nil {
		return nil, err
	}
	return &sqlConn{conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (c *sqlPool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	return c.DB.QueryContext(ctx, sql, args...)
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (c *sqlPool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return c.DB.QueryRowContext(ctx, sql, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *sqlPool) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := c.DB.ExecContext(ctx, sql, args...)
	return err
}

// Begin implements [database.Pool].
func (c *sqlPool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.DB.BeginTx(ctx, transactionOptionsToSql(opts))
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx}, nil
}

// Close implements [database.Pool].
func (c *sqlPool) Close(_ context.Context) error {
	return c.DB.Close()
}
