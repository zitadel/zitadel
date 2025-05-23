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

var (
	_ database.Pool = (*pgxPool)(nil)
)

// Acquire implements [database.Pool].
func (c *pgxPool) Acquire(ctx context.Context) (database.Client, error) {
	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxConn{Conn: conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (c *pgxPool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := c.Pool.Query(ctx, sql, args...)
	return &Rows{rows}, err
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (c *pgxPool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return c.Pool.QueryRow(ctx, sql, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *pgxPool) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := c.Pool.Exec(ctx, sql, args...)
	return err
}

// Begin implements [database.Pool].
func (c *pgxPool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.Pool.BeginTx(ctx, transactionOptionsToPgx(opts))
	if err != nil {
		return nil, err
	}
	return &pgxTx{tx}, nil
}

// Close implements [database.Pool].
func (c *pgxPool) Close(_ context.Context) error {
	c.Pool.Close()
	return nil
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
	return err
}
