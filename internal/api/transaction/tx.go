package transaction

import (
	"context"
	"database/sql"
)

type txKey struct{}

var key *txKey = (*txKey)(nil)

func WithTx(parent context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(parent, key, tx)
}

func FromContext(ctx context.Context) *sql.Tx {
	value := ctx.Value(key)
	if tx, ok := value.(*sql.Tx); ok {
		return tx
	}

	return nil
}
