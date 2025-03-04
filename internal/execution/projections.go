package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/queue"
)

var (
	projections []*handler.Handler
	worker      *Worker
)

func Register(
	ctx context.Context,
	executionsCustomConfig projection.CustomConfig,
	workerConfig WorkerConfig,
	queries eventHandlerQueries,
	eventTypes []string,
	queue *queue.Queue,
) {
	queue.ShouldStart()
	projections = []*handler.Handler{
		NewEventHandler(ctx, projection.ApplyCustomConfig(executionsCustomConfig), eventTypes, eventstore.AggregateTypeFromEventType, queries, queue),
	}
	worker = NewWorker(workerConfig, queue)
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}

func Init(ctx context.Context) error {
	for _, p := range projections {
		if err := p.Init(ctx); err != nil {
			return err
		}
	}
	return nil
}

func ProjectInstance(ctx context.Context) error {
	for _, projection := range projections {
		_, err := projection.Trigger(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func Projections() []*handler.Handler {
	return projections
}
