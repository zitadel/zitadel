package gosql

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type sqlConn struct{ *sql.Conn }

var _ database.Client = (*sqlConn)(nil)

// Release implements [database.Client].
func (c *sqlConn) Release(_ context.Context) error {
	return c.Conn.Close()
}

// Begin implements [database.Client].
func (c *sqlConn) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.Conn.BeginTx(ctx, transactionOptionsToSql(opts))
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx}, nil
}

// Query implements sql.Client.
// Subtle: this method shadows the method (*Conn).Query of pgxConn.Conn.
func (c *sqlConn) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	return c.Conn.QueryContext(ctx, sql, args...)
}

// QueryRow implements sql.Client.
// Subtle: this method shadows the method (*Conn).QueryRow of pgxConn.Conn.
func (c *sqlConn) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return c.Conn.QueryRowContext(ctx, sql, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *sqlConn) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := c.Conn.ExecContext(ctx, sql, args...)
	return err
}
