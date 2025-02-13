package logged

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/backend/storage/repository"
)

type Instance struct {
	*slog.Logger

	next repository.InstanceRepository
}

func NewInstance(logger *slog.Logger, next repository.InstanceRepository) *Instance {
	return &Instance{Logger: logger, next: next}
}

var _ repository.InstanceRepository = (*Instance)(nil)

func (i *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	i.Logger.InfoContext(ctx, "By ID Query", slog.String("id", id))
	return i.next.ByID(ctx, id)
}

func (i *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	i.Logger.InfoContext(ctx, "By Domain Query", slog.String("domain", domain))
	return i.next.ByDomain(ctx, domain)
}

func (i *Instance) SetUp(ctx context.Context, instance *repository.Instance) error {
	err := i.next.SetUp(ctx, instance)
	if err != nil {
		i.Logger.ErrorContext(ctx, "Failed to set up instance", slog.Any("instance", instance), slog.Any("cause", err))
		return err
	}
	i.Logger.InfoContext(ctx, "Instance set up", slog.Any("instance", instance))
	return nil
}
