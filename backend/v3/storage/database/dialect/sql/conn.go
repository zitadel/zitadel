package sql

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Conn struct {
	*sql.Conn
}

var _ database.Client = (*Conn)(nil)

// Release implements [database.Client].
func (c *Conn) Release(_ context.Context) error {
	return c.Close()
}

// Begin implements [database.Client].
func (c *Conn) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToSQL(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{tx}, nil
}

// Query implements [database.Client].
// Subtle: this method shadows the method (*Conn).Query of pgxConn.Conn.
func (c *Conn) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	//nolint:rowserrcheck // Rows.Close is called by the caller
	rows, err := c.Conn.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryContext implements [database.Client].
// Subtle: this method shadows the method (*Conn).QueryContext of sqlConn.Conn.
// QueryContext is for backwards compatibility, it calls Query.
func (c *Conn) QueryContext(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return c.Query(ctx, stmt, args...)
}

// QueryRow implements sql.Client.
// Subtle: this method shadows the method (*Conn).QueryRow of pgxConn.Conn.
func (c *Conn) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.Conn.QueryRowContext(ctx, sql, args...)}
}

// QueryContext implements database.Client.
// Subtle: this method shadows the method (*Conn).QueryContext of sqlConn.Conn.
// QueryContext is for backwards compatibility, it calls Query.
func (c *Conn) QueryRowContext(ctx context.Context, stmt string, args ...any) database.Row {
	return c.QueryRow(ctx, stmt, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *Conn) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.Conn.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// ExecContext implements database.Client.
// Subtle: this method shadows the method (*Conn).ExecContext of sqlConn.Conn.
// ExecContext is for backwards compatibility, it calls Exec.
func (c *Conn) ExecContext(ctx context.Context, stmt string, args ...any) (int64, error) {
	return c.Exec(ctx, stmt, args...)
}

// Migrate implements [database.Migrator].
func (c *Conn) Migrate(ctx context.Context) error {
	return ErrMigrate
}
