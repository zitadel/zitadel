package database

import (
	"context"
	"database/sql"
)

type transactionKey struct{}

// ContextWithTransaction lets command-side reads reuse the lifecycle transaction.
func ContextWithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

func TransactionFromContext(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(transactionKey{}).(*sql.Tx)
	return tx
}
