package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	projections []*handler.Handler
	worker      *Worker
)

func Create(
	ctx context.Context,
	executionsCustomConfig projection.CustomConfig,
	queries *query.Queries,
	es *eventstore.Eventstore,
) {
	projections = []*handler.Handler{
		NewExecutionsHandler(ctx, projection.ApplyCustomConfig(executionsCustomConfig), es, queries),
	}
}

func Register(
	ctx context.Context,
	executionsCustomConfig projection.CustomConfig,
	workerConfig WorkerConfig,
	queries *query.Queries,
	es *eventstore.Eventstore,
	client *database.DB,
) {
	Create(ctx, executionsCustomConfig, queries, es)
	q := NewExecutionsQueries(queries, client)
	worker = NewWorker(workerConfig, client, q)
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
	worker.Start(ctx)
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
