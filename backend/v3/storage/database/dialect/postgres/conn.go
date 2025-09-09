package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

type Conn struct {
	*pgxpool.Conn
}

// QueryContext implements database.Client.
func (c *Conn) QueryContext(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return c.Query(ctx, stmt, args...)
}

var _ database.Client = (*Conn)(nil)

// Release implements [database.Client].
func (c *Conn) Release(_ context.Context) error {
	c.Conn.Release()
	return nil
}

// Begin implements [database.Client].
func (c *Conn) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToPgx(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{tx}, nil
}

// Query implements sql.Client.
// Subtle: this method shadows the method (*Conn).Query of pgxConn.Conn.
func (c *Conn) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := c.Conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements sql.Client.
// Subtle: this method shadows the method (*Conn).QueryRow of pgxConn.Conn.
func (c *Conn) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.Conn.QueryRow(ctx, sql, args...)}
}

// QueryContext implements [database.Client].
// Subtle: this method shadows the method (*Conn).QueryContext of sqlConn.Conn.
// QueryContext is for backwards compatibility, it calls Query.
func (c *Conn) QueryRowContext(ctx context.Context, stmt string, args ...any) database.Row {
	return c.QueryRow(ctx, stmt, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *Conn) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.Conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected(), nil
}

// ExecContext implements database.Client.
func (c *Conn) ExecContext(ctx context.Context, stmt string, args ...any) (int64, error) {
	return c.Exec(ctx, stmt, args...)
}

// Migrate implements [database.Migrator].
func (c *Conn) Migrate(ctx context.Context) error {
	if isMigrated {
		return nil
	}
	err := migration.Migrate(ctx, c.Conn.Conn())
	isMigrated = err == nil
	return wrapError(err)
}
