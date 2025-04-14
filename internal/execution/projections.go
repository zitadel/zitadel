package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/queue"
)

var (
	projections []*handler.Handler
)

func Register(
	ctx context.Context,
	executionsCustomConfig projection.CustomConfig,
	workerConfig WorkerConfig,
	queries *query.Queries,
	eventTypes []string,
	queue *queue.Queue,
) {
	queue.ShouldStart()
	projections = []*handler.Handler{
		NewEventHandler(ctx, projection.ApplyCustomConfig(executionsCustomConfig), eventTypes, eventstore.AggregateTypeFromEventType, queries, queue),
	}
	queue.AddWorkers(NewWorker(workerConfig))
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}
