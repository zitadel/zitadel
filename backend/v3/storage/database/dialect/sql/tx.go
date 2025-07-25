package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sqlTx struct{ *sql.Tx }

var _ database.Transaction = (*sqlTx)(nil)

func SQLTx(tx *sql.Tx) *sqlTx {
	return &sqlTx{
		Tx: tx,
	}
}

// Commit implements [database.Transaction].
func (tx *sqlTx) Commit(ctx context.Context) error {
	return wrapError(tx.Tx.Commit())
}

// Rollback implements [database.Transaction].
func (tx *sqlTx) Rollback(ctx context.Context) error {
	return wrapError(tx.Tx.Rollback())
}

// End implements [database.Transaction].
func (tx *sqlTx) End(ctx context.Context, err error) error {
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
func (tx *sqlTx) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := tx.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	return &Rows{rows}, nil
}

// QueryRow implements [database.Transaction].
// Subtle: this method shadows the method (Tx).QueryRow of pgxTx.Tx.
func (tx *sqlTx) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return &Row{tx.QueryRowContext(ctx, sql, args...)}
}

// Exec implements [database.Transaction].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (tx *sqlTx) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, wrapError(err)
	}
	return res.RowsAffected()
}

// Begin implements [database.Transaction].
// As postgres does not support nested transactions we use savepoints to emulate them.
func (tx *sqlTx) Begin(ctx context.Context) (database.Transaction, error) {
	_, err := tx.ExecContext(ctx, createSavepoint)
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
