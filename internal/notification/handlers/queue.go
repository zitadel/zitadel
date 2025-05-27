package handlers

import (
	"context"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/queue"
)

type Queue interface {
	Insert(ctx context.Context, args river.JobArgs, opts ...queue.InsertOpt) error
}
