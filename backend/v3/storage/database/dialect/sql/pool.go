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
func (p *sqlPool) Acquire(ctx context.Context) (database.Connection, error) {
	conn, err := p.Conn(ctx)
	if err != nil {
		return nil, wrapError(err)
	}
	return &sqlConn{Conn: conn}, nil
}

// Query implements [database.Pool].
// Subtle: this method shadows the method (Pool).Query of pgxPool.Pool.
func (p *sqlPool) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	//nolint:rowserrcheck // Rows.Close is called by the caller
	rows, err := p.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements [database.Pool].
// Subtle: this method shadows the method (Pool).QueryRow of pgxPool.Pool.
func (p *sqlPool) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{p.QueryRowContext(ctx, sql, args...)}
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (p *sqlPool) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := p.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// Begin implements [database.Pool].
func (p *sqlPool) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	tx, err := p.BeginTx(ctx, transactionOptionsToSQL(opts))
	if err != nil {
		return nil, wrapError(err)
	}
	return &Transaction{tx}, nil
}

// Ping implements [database.Pool].
func (p *sqlPool) Ping(ctx context.Context) error {
	return wrapError(p.PingContext(ctx))
}

// Close implements [database.Pool].
func (p *sqlPool) Close(_ context.Context) error {
	return p.DB.Close()
}

// Migrate implements [database.Migrator].
func (p *sqlPool) Migrate(ctx context.Context) error {
	return ErrMigrate
}
