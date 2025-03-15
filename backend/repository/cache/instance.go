package cache

import (
	"context"
	"sync"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/cache"
	"github.com/zitadel/zitadel/backend/storage/cache/gomap"
)

type Instance struct {
	mu       *sync.RWMutex
	byID     cache.Cache[string, *repository.Instance]
	byDomain cache.Cache[string, *repository.Instance]
}

func NewInstance() *Instance {
	return &Instance{
		mu:       &sync.RWMutex{},
		byID:     gomap.New[string, *repository.Instance](),
		byDomain: gomap.New[string, *repository.Instance](),
	}
}

func (i *Instance) Set(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
	i.set(instance, "")
	return instance, nil
}

func (i *Instance) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	instance, _ := i.byID.Get(id)
	return instance, nil
}

func (i *Instance) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	instance, _ := i.byDomain.Get(domain)
	return instance, nil
}

func (i *Instance) set(instance *repository.Instance, domain string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if domain != "" {
		i.byDomain.Set(domain, instance)
	}
	i.byID.Set(instance.ID, instance)
}
