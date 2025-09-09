package eventstore

import (
	"context"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/queue"
)

//go:generate mockgen -package mock -destination ./mock/queue.mock.go github.com/zitadel/zitadel/internal/eventstore ExecutionQueue

type ExecutionQueue interface {
	// InsertManyFastTx wraps [river.Client.InsertManyFastTx] to insert all jobs in
	// a single `COPY FROM` execution, within an existing transaction.
	//
	// Opts are applied to each job before sending them to river.
	InsertManyFastTx(ctx context.Context, tx database.Transaction, args []river.JobArgs, opts ...queue.InsertOpt) error
}
