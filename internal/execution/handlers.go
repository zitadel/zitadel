package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/execution"
)

const (
	ExecutionsHandlerTable = "projections.executions_handler"
)

type executionsHandler struct{}

func NewExecutionsHandler(
	ctx context.Context,
	config handler.Config,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(executionsHandler))
}

func (u *executionsHandler) Name() string {
	return ExecutionsHandlerTable
}

func (u *executionsHandler) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: execution.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  execution.SetEventType,
					Reduce: u.reduceInitCodeAdded,
				},
				{
					Event:  execution.SetEventV2Type,
					Reduce: u.reduceInitCodeAdded,
				},
				{
					Event:  execution.SetEventV2Type,
					Reduce: u.reduceInitCodeAdded,
				},
			},
		},
	}
}
