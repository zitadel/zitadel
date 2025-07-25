package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

type pgxConn struct {
	*pgxpool.Conn
}

var _ database.Client = (*pgxConn)(nil)

// Release implements [database.Client].
func (c *pgxConn) Release(_ context.Context) error {
	c.Conn.Release()
	return nil
}

// Begin implements [database.Client].
func (c *pgxConn) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToPgx(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &pgxTx{tx}, nil
}

// Query implements sql.Client.
// Subtle: this method shadows the method (*Conn).Query of pgxConn.Conn.
func (c *pgxConn) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {

	rows, err := c.Conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements sql.Client.
// Subtle: this method shadows the method (*Conn).QueryRow of pgxConn.Conn.
func (c *pgxConn) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.Conn.QueryRow(ctx, sql, args...)}
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *pgxConn) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.Conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected(), nil
}

// Migrate implements [database.Migrator].
func (c *pgxConn) Migrate(ctx context.Context) error {
	if isMigrated {
		return nil
	}
	err := migration.Migrate(ctx, c.Conn.Conn())
	isMigrated = err == nil
	return wrapError(err)
}
