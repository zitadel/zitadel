package event

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/eventstore"
	"github.com/zitadel/zitadel/backend/storage/repository"
)

var _ repository.InstanceRepository = (*Instance)(nil)

type Instance struct {
	*eventstore.Eventstore

	next repository.InstanceRepository
}

func NewInstance(eventstore *eventstore.Eventstore, next repository.InstanceRepository) *Instance {
	return &Instance{next: next, Eventstore: eventstore}
}

func (i *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	return i.next.ByID(ctx, id)
}

func (i *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	return i.next.ByDomain(ctx, domain)
}

func (i *Instance) SetUp(ctx context.Context, instance *repository.Instance) error {
	err := i.next.SetUp(ctx, instance)
	if err != nil {
		return err
	}

	return i.Push(ctx, instance)
}
