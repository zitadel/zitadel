package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/query/projection"
)

var (
	executionProjection  *handler.Handler
	conditionProjections []*handler.Handler
)

func Register(
	ctx context.Context,
	executionsCustomConfig projection.CustomConfig,
) {
	executionProjection = NewExecutionsHandler(ctx, projection.ApplyCustomConfig(executionsCustomConfig), conditionProjections)
}

func Start(ctx context.Context) {
	executionProjection.Start(ctx)
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
