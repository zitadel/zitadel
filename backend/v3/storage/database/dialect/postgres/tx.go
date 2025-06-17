package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type pgxTx struct{ pgx.Tx }

var _ database.Transaction = (*pgxTx)(nil)

// Commit implements [database.Transaction].
func (tx *pgxTx) Commit(ctx context.Context) error {
	return tx.Tx.Commit(ctx)
}

// Rollback implements [database.Transaction].
func (tx *pgxTx) Rollback(ctx context.Context) error {
	return tx.Tx.Rollback(ctx)
}

// End implements [database.Transaction].
func (tx *pgxTx) End(ctx context.Context, err error) error {
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
func (tx *pgxTx) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := tx.Tx.Query(ctx, sql, args...)
	return &Rows{rows}, err
}

// QueryRow implements [database.Transaction].
// Subtle: this method shadows the method (Tx).QueryRow of pgxTx.Tx.
func (tx *pgxTx) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return tx.Tx.QueryRow(ctx, sql, args...)
}

// Exec implements [database.Transaction].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (tx *pgxTx) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	res, err := tx.Tx.Exec(ctx, sql, args...)
	return res.RowsAffected(), err
}

// Begin implements [database.Transaction].
// As postgres does not support nested transactions we use savepoints to emulate them.
func (tx *pgxTx) Begin(ctx context.Context) (database.Transaction, error) {
	savepoint, err := tx.Tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTx{savepoint}, nil
}

func transactionOptionsToPgx(opts *database.TransactionOptions) pgx.TxOptions {
	if opts == nil {
		return pgx.TxOptions{}
	}

	return pgx.TxOptions{
		IsoLevel:   isolationToPgx(opts.IsolationLevel),
		AccessMode: accessModeToPgx(opts.AccessMode),
	}
}

func isolationToPgx(isolation database.IsolationLevel) pgx.TxIsoLevel {
	switch isolation {
	case database.IsolationLevelSerializable:
		return pgx.Serializable
	case database.IsolationLevelReadCommitted:
		return pgx.ReadCommitted
	default:
		return pgx.Serializable
	}
}

func accessModeToPgx(accessMode database.AccessMode) pgx.TxAccessMode {
	switch accessMode {
	case database.AccessModeReadWrite:
		return pgx.ReadWrite
	case database.AccessModeReadOnly:
		return pgx.ReadOnly
	default:
		return pgx.ReadWrite
	}
}
