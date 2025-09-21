package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres/migration"
)

type pgxPool struct {
	*pgxpool.Pool
}

var _ database.Pool = (*pgxPool)(nil)

func PGxPool(pool *pgxpool.Pool) *pgxPool {
	return &pgxPool{
		Pool: pool,
	}
}

// Acquire implements [database.Pool].
func (c *pgxPool) Acquire(ctx context.Context) (database.Connection, error) {
	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return nil, wrapError(err)
	}
	return &pgxConn{Conn: conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (c *pgxPool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := c.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (c *pgxPool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.Pool.QueryRow(ctx, sql, args...)}
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *pgxPool) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected(), nil
}

// Begin implements [database.Pool].
func (c *pgxPool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToPgx(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &Transaction{tx}, nil
}

// Close implements [database.Pool].
func (c *pgxPool) Close(_ context.Context) error {
	c.Pool.Close()
	return nil
}

// Ping implements [database.Pool].
func (c *pgxPool) Ping(ctx context.Context) error {
	return wrapError(c.Pool.Ping(ctx))
}

// Migrate implements [database.Migrator].
func (c *pgxPool) Migrate(ctx context.Context) error {
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
func (c *pgxPool) MigrateTest(ctx context.Context) error {
	client, err := c.Pool.Acquire(ctx)
	if err != nil {
		return err
	}

	err = migration.Migrate(ctx, client.Conn())
	isMigrated = err == nil
	return err
}
