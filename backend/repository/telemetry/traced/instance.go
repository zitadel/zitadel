package traced

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

var _ repository.InstanceRepository = (*Instance)(nil)

type Instance struct {
	*tracing.Tracer

	next repository.InstanceRepository
}

func NewInstance(tracer *tracing.Tracer, next repository.InstanceRepository) *Instance {
	return &Instance{Tracer: tracer, next: next}
}

func (i *Instance) SetNext(next repository.InstanceRepository) *Instance {
	return &Instance{Tracer: i.Tracer, next: next}
}

// ByDomain implements [repository.InstanceRepository].
func (i *Instance) ByDomain(ctx context.Context, domain string) (instance *repository.Instance, err error) {
	i.Tracer.Decorate(ctx, func(ctx context.Context) error {
		instance, err = i.next.ByDomain(ctx, domain)
		return err
	})

	return instance, err
}

// ByID implements [repository.InstanceRepository].
func (i *Instance) ByID(ctx context.Context, id string) (instance *repository.Instance, err error) {
	i.Tracer.Decorate(ctx, func(ctx context.Context) error {
		instance, err = i.next.ByID(ctx, id)
		return err
	})

	return instance, err
}

// SetUp implements [repository.InstanceRepository].
func (i *Instance) SetUp(ctx context.Context, instance *repository.Instance) (err error) {
	i.Tracer.Decorate(ctx, func(ctx context.Context) error {
		err = i.next.SetUp(ctx, instance)
		return err
	})

	return err
}
