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

func Register(
	ctx context.Context,
	executionsCustomConfig projection.CustomConfig,
	workerConfig WorkerConfig,
	queries *query.Queries,
	es *eventstore.Eventstore,
	client *database.DB,
) {
	q := NewExecutionsQueries(queries, client)
	projections = append(projections, NewExecutionsHandler(ctx, projection.ApplyCustomConfig(executionsCustomConfig), es, queries))
	worker = NewWorker(workerConfig, client, q)
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
	worker.Start(ctx)
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
