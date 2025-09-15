package eventstore

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/queue"
)

//go:generate mockgen -package mock -destination ./mock/queue.mock.go github.com/zitadel/zitadel/internal/eventstore ExecutionQueue

type ExecutionQueue interface {
	// InsertManyFastTx wraps [river.Client.InsertManyFastTx] to insert all jobs in
	// a single `COPY FROM` execution, within an existing transaction.
	//
	// Opts are applied to each job before sending them to river.
	InsertManyFastTx(ctx context.Context, tx *sql.Tx, args []river.JobArgs, opts ...queue.InsertOpt) error
}
