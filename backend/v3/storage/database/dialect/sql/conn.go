package sql

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sqlConn struct {
	*sql.Conn
}

var _ database.Client = (*sqlConn)(nil)

// Release implements [database.Client].
func (c *sqlConn) Release(_ context.Context) error {
	return c.Close()
}

// Begin implements [database.Client].
func (c *sqlConn) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToSQL(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &sqlTx{tx}, nil
}

// Query implements sql.Client.
// Subtle: this method shadows the method (*Conn).Query of pgxConn.Conn.
func (c *sqlConn) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	//nolint:rowserrcheck // Rows.Close is called by the caller
	rows, err := c.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements sql.Client.
// Subtle: this method shadows the method (*Conn).QueryRow of pgxConn.Conn.
func (c *sqlConn) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.QueryRowContext(ctx, sql, args...)}
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *sqlConn) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// Migrate implements [database.Migrator].
func (c *sqlConn) Migrate(ctx context.Context) error {
	return ErrMigrate
}
