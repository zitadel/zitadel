package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type Tx struct{ *sql.Tx }

var _ database.Transaction = (*Tx)(nil)

func SQLTx(tx *sql.Tx) *Tx {
	return &Tx{
		Tx: tx,
	}
}

// Commit implements [database.Transaction].
func (tx *Tx) Commit(ctx context.Context) error {
	return wrapError(tx.Tx.Commit())
}

// Rollback implements [database.Transaction].
func (tx *Tx) Rollback(ctx context.Context) error {
	return wrapError(tx.Tx.Rollback())
}

// End implements [database.Transaction].
func (tx *Tx) End(ctx context.Context, err error) error {
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

// Query implements [database.Transaction].
// Subtle: this method shadows the method (Tx).Query of pgxTx.Tx.
func (tx *Tx) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	//nolint:rowserrcheck // Rows.Close is called by the caller
	rows, err := tx.Tx.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryContext implements [database.Transaction].
// Subtle: this method shadows the method (*Tx).QueryContext of [sqlTx.Tx].
// QueryContext is for backwards compatibility, it calls Query.
func (tx *Tx) QueryContext(ctx context.Context, stmt string, args ...any) (database.Rows, error) {
	return tx.Query(ctx, stmt, args...)
}

// QueryRow implements [database.Transaction].
// Subtle: this method shadows the method (Tx).QueryRow of pgxTx.Tx.
func (tx *Tx) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{tx.Tx.QueryRowContext(ctx, sql, args...)}
}

// QueryContext implements [database.Transaction].
// Subtle: this method shadows the method (*Conn).QueryContext of sqlConn.Conn.
// QueryContext is for backwards compatibility, it calls Query.
func (c *Tx) QueryRowContext(ctx context.Context, stmt string, args ...any) database.Row {
	return c.QueryRow(ctx, stmt, args...)
}

// Exec implements [database.Transaction].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (tx *Tx) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := tx.Tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// ExecContext implements [database.Transaction].
// Subtle: this method shadows the method (*Tx).ExecContext of [sqlTx.Tx].
// ExecContext is for backwards compatibility, it calls Exec.
func (tx *Tx) ExecContext(ctx context.Context, stmt string, args ...any) (int64, error) {
	return tx.Exec(ctx, stmt, args...)
}

// Begin implements [database.Transaction].
// As postgres does not support nested transactions we use savepoints to emulate them.
func (tx *Tx) Begin(ctx context.Context) (database.Transaction, error) {
	_, err := tx.Exec(ctx, createSavepoint)
	if err != nil {
		return nil, wrapError(err)
	}
	return &sqlSavepoint{tx}, nil
}

func transactionOptionsToSQL(opts *database.TransactionOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}

	return &sql.TxOptions{
		Isolation: isolationToSQL(opts.IsolationLevel),
		ReadOnly:  accessModeToSQL(opts.AccessMode),
	}
}

func isolationToSQL(isolation database.IsolationLevel) sql.IsolationLevel {
	switch isolation {
	case database.IsolationLevelSerializable:
		return sql.LevelSerializable
	case database.IsolationLevelReadCommitted:
		return sql.LevelReadCommitted
	default:
		return sql.LevelSerializable
	}
}

func accessModeToSQL(accessMode database.AccessMode) bool {
	return accessMode == database.AccessModeReadOnly
}
