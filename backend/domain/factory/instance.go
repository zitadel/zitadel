package factory

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type InstanceBuilder struct {
	tracer *tracing.Tracer
	logger *logging.Logger
}

func (f *InstanceBuilder) BuildSetUpInstance(tx database.Transaction) func(ctx context.Context, instance *repository.Instance) error {
	return func(ctx context.Context, instance *repository.Instance) error {
		return tx.Exec(ctx, "INSERT INTO instances (id, name) VALUES ($1, $2)", instance.ID, instance.Name)
	}
}
