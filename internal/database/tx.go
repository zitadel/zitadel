package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zitadel/logging"
)

func ExecuteWithRetries(ctx context.Context, db *DB, maxTransactionRetries uint8, fn func(tx *sql.Tx) error) (err error) {
	conn, err := db.DB.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// execute transactions at least once
	if maxTransactionRetries == 0 {
		maxTransactionRetries = 1
	}

	for i := uint8(0); i < maxTransactionRetries; i++ {
		var shouldRetry bool
		shouldRetry, err = execute(ctx, db, conn, fn)
		logging.WithFields("willRetry", shouldRetry, "count", i).OnError(err).Debug("exec failed")
		if !shouldRetry {
			break
		}
	}
	return err
}

func execute(ctx context.Context, db *DB, conn *sql.Conn, fn func(tx *sql.Tx) error) (shouldRetry bool, err error) {
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}

	if err = fn(tx); err != nil {
		rollbackErr := tx.Rollback()
		logging.OnError(rollbackErr).Debug("rollback of failed tx failed")
		code := errCode(err)
		return db.IsRetryable(code), err
	}

	err = tx.Commit()
	// retry if commit failed
	return err != nil, err
}

func errCode(err error) string {
	var sqlErr errWithSQLState
	if errors.As(err, &sqlErr) {
		return sqlErr.SQLState()
	}

	return ""
}

// errWithSQLState is implemented by pgx (pgconn.PgError) and lib/pq
type errWithSQLState interface {
	SQLState() string
}
