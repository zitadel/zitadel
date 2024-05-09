package readmodel

import (
	"database/sql"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type QueryOpt func(opts []eventstore.QueryOpt) []eventstore.QueryOpt

func WithTx(tx *sql.Tx) QueryOpt {
	return func(opts []eventstore.QueryOpt) []eventstore.QueryOpt {
		return append(opts, eventstore.SetQueryTx(tx))
	}
}
