package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

type Pool struct {
	*pgxpool.Pool
}

var _ database.Pool = (*Pool)(nil)

func PGxPool(pool *pgxpool.Pool) *Pool {
	return &Pool{
		Pool: pool,
	}
}

// Acquire implements [database.Pool].
func (c *Pool) Acquire(ctx context.Context) (database.Client, error) {
	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Conn{Conn: conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (c *Pool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := c.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryContext implements [database.Pool.]
func (c *Pool) QueryContext(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return c.Query(ctx, stmt, args...)
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (c *Pool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.Pool.QueryRow(ctx, sql, args...)}
}

// QueryContext implements [database.Pool].
// Subtle: this method shadows the method (*Conn).QueryContext of sqlConn.Conn.
// QueryContext is for backwards compatibility, it calls Query.
func (c *Pool) QueryRowContext(ctx context.Context, stmt string, args ...any) database.Row {
	return c.QueryRow(ctx, stmt, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *Pool) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected(), nil
}

// ExecContext implements [database.Pool.]
func (c *Pool) ExecContext(ctx context.Context, stmt string, args ...any) (int64, error) {
	return c.Exec(ctx, stmt, args...)
}

// Begin implements [database.Pool].
func (c *Pool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToPgx(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &Tx{tx}, nil
}

// Close implements [database.Pool].
func (c *Pool) Close(_ context.Context) error {
	c.Pool.Close()
	return nil
}

// Migrate implements [database.Migrator].
func (c *Pool) Migrate(ctx context.Context) error {
	if isMigrated {
		return nil
	}

	client, err := c.Pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = migration.Migrate(ctx, client.Conn())
	isMigrated = err == nil
	return wrapError(err)
}

// Migrate implements [database.PoolTest].
func (c *Pool) MigrateTest(ctx context.Context) error {
	client, err := c.Pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = migration.Migrate(ctx, client.Conn())
	isMigrated = err == nil
	return err
}
