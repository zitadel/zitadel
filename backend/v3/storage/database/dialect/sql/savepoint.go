package sql

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ database.Transaction = (*sqlSavepoint)(nil)

const (
	savepointName       = "zitadel_savepoint"
	createSavepoint     = "SAVEPOINT " + savepointName
	rollbackToSavepoint = "ROLLBACK TO SAVEPOINT " + savepointName
	commitSavepoint     = "RELEASE SAVEPOINT " + savepointName
)

type sqlSavepoint struct {
	parent database.Transaction
}

// Commit implements [database.Transaction].
func (s *sqlSavepoint) Commit(ctx context.Context) error {
	_, err := s.parent.Exec(ctx, commitSavepoint)
	return wrapError(err)
}

// Rollback implements [database.Transaction].
func (s *sqlSavepoint) Rollback(ctx context.Context) error {
	_, err := s.parent.Exec(ctx, rollbackToSavepoint)
	return wrapError(err)
}

// End implements [database.Transaction].
func (s *sqlSavepoint) End(ctx context.Context, err error) error {
	if err != nil {
		rollbackErr := s.Rollback(ctx)
		if rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		return err
	}
	return s.Commit(ctx)
}

// Query implements [database.Transaction].
// Subtle: this method shadows the method (Tx).Query of pgxTx.Tx.
func (s *sqlSavepoint) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	return s.parent.Query(ctx, sql, args...)
}

// QueryRow implements [database.Transaction].
// Subtle: this method shadows the method (Tx).QueryRow of pgxTx.Tx.
func (s *sqlSavepoint) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	return s.parent.QueryRow(ctx, sql, args...)
}

// Exec implements [database.Transaction].
// Subtle: this method shadows the method (Pool).Exec of pgxPool.Pool.
func (s *sqlSavepoint) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	return s.parent.Exec(ctx, sql, args...)
}

// Begin implements [database.Transaction].
// As postgres does not support nested transactions we use savepoints to emulate them.
func (s *sqlSavepoint) Begin(ctx context.Context) (database.Transaction, error) {
	return s.parent.Begin(ctx)
}
