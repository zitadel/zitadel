package gosql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zitadel/zitadel/backend/storage/database"
)

type sqlTx struct{ *sql.Tx }

var _ database.Transaction = (*sqlTx)(nil)

// Commit implements [database.Transaction].
func (tx *sqlTx) Commit(_ context.Context) error {
	return tx.Tx.Commit()
}

// Rollback implements [database.Transaction].
func (tx *sqlTx) Rollback(_ context.Context) error {
	return tx.Tx.Rollback()
}

// End implements [database.Transaction].
func (tx *sqlTx) End(ctx context.Context, err error) error {
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// Query implements [database.Transaction].
// Subtle: this method shadows the method (Tx).Query of pgxTx.Tx.
func (tx *sqlTx) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	return tx.Tx.QueryContext(ctx, sql, args...)
}

// QueryRow implements [database.Transaction].
// Subtle: this method shadows the method (Tx).QueryRow of pgxTx.Tx.
func (tx *sqlTx) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return tx.Tx.QueryRowContext(ctx, sql, args...)
}

// Exec implements [database.Pool].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (tx *sqlTx) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := tx.Tx.ExecContext(ctx, sql, args...)
	return err
}

// Begin implements [database.Transaction].
// it is unimplemented
func (tx *sqlTx) Begin(ctx context.Context, opts *database.TransactionOptions) (database.Transaction, error) {
	return nil, errors.New("nested transactions are not supported")
}

func transactionOptionsToSql(opts *database.TransactionOptions) *sql.TxOptions {
	if opts == nil {
		return nil
	}

	return &sql.TxOptions{
		Isolation: isolationToSql(opts.IsolationLevel),
		ReadOnly:  opts.AccessMode == database.AccessModeReadOnly,
	}
}

func isolationToSql(isolation database.IsolationLevel) sql.IsolationLevel {
	switch isolation {
	case database.IsolationLevelSerializable:
		return sql.LevelSerializable
	case database.IsolationLevelReadCommitted:
		return sql.LevelReadCommitted
	default:
		return sql.LevelSerializable
	}
}
