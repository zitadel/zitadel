package cache

import (
	"context"
	"sync"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/cache"
)

type Instance struct {
	mu       *sync.RWMutex
	byID     cache.Cache[string, *repository.Instance]
	byDomain cache.Cache[string, *repository.Instance]

	next repository.InstanceRepository
}

func (i *Instance) SetNext(next repository.InstanceRepository) *Instance {
	return &Instance{
		mu:       i.mu,
		byID:     i.byID,
		byDomain: i.byDomain,
		next:     next,
	}
}

// ByDomain implements repository.InstanceRepository.
func (i *Instance) ByDomain(ctx context.Context, domain string) (instance *repository.Instance, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if instance, ok := i.byDomain.Get(domain); ok {
		return instance, nil
	}

	instance, err = i.next.ByDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	i.set(instance, domain)

	return instance, nil
}

// ByID implements repository.InstanceRepository.
func (i *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if instance, ok := i.byID.Get(id); ok {
		return instance, nil
	}

	instance, err := i.next.ByID(ctx, id)
	if err != nil {
		return nil, err

	}

	i.set(instance, "")
	return instance, nil
}

// SetUp implements repository.InstanceRepository.
func (i *Instance) SetUp(ctx context.Context, instance *repository.Instance) error {
	err := i.next.SetUp(ctx, instance)
	if err != nil {
		return err
	}

	i.set(instance, "")
	return nil
}

var _ repository.InstanceRepository = (*Instance)(nil)

func (i *Instance) set(instance *repository.Instance, domain string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if domain != "" {
		i.byDomain.Set(domain, instance)
	}
	i.byID.Set(instance.ID, instance)
}
