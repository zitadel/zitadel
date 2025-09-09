package sql

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Pool struct {
	*sql.DB
}

var _ database.Pool = (*Pool)(nil)

func SQLPool(db *sql.DB) *Pool {
	return &Pool{
		DB: db,
	}
}

// Acquire implements [database.Pool].
func (c *Pool) Acquire(ctx context.Context) (database.Client, error) {
	conn, err := c.Conn(ctx)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Conn{Conn: conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (c *Pool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	//nolint:rowserrcheck // Rows.Close is called by the caller
	rows, err := c.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryContext implements [database.Pool].
// Subtle: this method shadows the method (*DB).QueryContext of [sqlPool.DB].
// QueryContext is for backwards compatibility, it calls Query.
func (c *Pool) QueryContext(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return c.Query(ctx, stmt, args...)
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (c *Pool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.DB.QueryRowContext(ctx, sql, args...)}
}

// QueryContext implements database.Client.
// Subtle: this method shadows the method (*Conn).QueryContext of sqlConn.Conn.
// QueryContext is for backwards compatibility, it calls Query.
func (c *Pool) QueryRowContext(ctx context.Context, stmt string, args ...any) database.Row {
	return c.QueryRow(ctx, stmt, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *Pool) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// ExecContext implements [database.Pool].
// Subtle: this method shadows the method (*DB).ExecContext of [sqlPool.DB].
// ExecContext is for backwards compatibility, it calls Exec.
func (c *Pool) ExecContext(ctx context.Context, stmt string, args ...any) (int64, error) {
	return c.Exec(ctx, stmt, args...)
}

// Begin implements [database.Pool].
func (c *Pool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToSQL(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{tx}, nil
}

// Close implements [database.Pool].
func (c *Pool) Close(_ context.Context) error {
	return c.DB.Close()
}

// Ping implements [database.Pool].
// Subtle: this method shadows the method (*DB).Ping of sqlPool.DB.
func (c *Pool) Ping(ctx context.Context) error {
	return c.PingContext(ctx)
}

// Migrate implements [database.Migrator].
func (c *Pool) Migrate(ctx context.Context) error {
	return ErrMigrate
}
