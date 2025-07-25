package sql

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sqlPool struct {
	*sql.DB
}

var _ database.Pool = (*sqlPool)(nil)

func SQLPool(db *sql.DB) *sqlPool {
	return &sqlPool{
		DB: db,
	}
}

// Acquire implements [database.Pool].
func (c *sqlPool) Acquire(ctx context.Context) (database.Client, error) {
	conn, err := c.Conn(ctx)
	if err != nil {
		return nil, wrapError(err)
	}
	return &sqlConn{Conn: conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (c *sqlPool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	//nolint:sqlclosecheck // Rows.Close is called by the caller
	rows, err := c.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (c *sqlPool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{c.QueryRowContext(ctx, sql, args...)}
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (c *sqlPool) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := c.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// Begin implements [database.Pool].
func (c *sqlPool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := c.BeginTx(ctx, transactionOptionsToSQL(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &sqlTx{tx}, nil
}

// Close implements [database.Pool].
func (c *sqlPool) Close(_ context.Context) error {
	return c.DB.Close()
}

// Migrate implements [database.Migrator].
func (c *sqlPool) Migrate(ctx context.Context) error {
	return ErrMigrate
}
