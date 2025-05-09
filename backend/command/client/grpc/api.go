package grpc

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/command"
	"github.com/zitadel/zitadel/backend/command/query"
	"github.com/zitadel/zitadel/backend/command/receiver"
	"github.com/zitadel/zitadel/backend/command/receiver/cache"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type api struct {
	db database.Pool

	manipulator receiver.InstanceManipulator
	reader      receiver.InstanceReader
	tracer      *tracing.Tracer
	logger      *logging.Logger
	cache       cache.Cache[receiver.InstanceIndex, string, *receiver.Instance]
}

func (a *api) CreateInstance(ctx context.Context) error {
	instance := &receiver.Instance{
		ID:   "123",
		Name: "test",
	}
	return command.Trace(
		a.tracer,
		command.SetCache(a.cache,
			command.Activity(a.logger, command.CreateInstance(a.manipulator, instance)),
			instance,
		),
	).Execute(ctx)
}

func (a *api) DeleteInstance(ctx context.Context) error {
	return command.Trace(
		a.tracer,
		command.DeleteCache(a.cache,
			command.Activity(
				a.logger,
				command.DeleteInstance(a.manipulator, &receiver.Instance{
					ID: "123",
				})),
			receiver.InstanceByID,
			"123",
		)).Execute(ctx)
}

func (a *api) InstanceByID(ctx context.Context) (*receiver.Instance, error) {
	q := query.InstanceByID(a.reader, "123")
	return q.Execute(ctx)
}
