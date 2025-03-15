package cache

import (
	"context"
	"sync"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/storage/cache"
)

type Instance struct {
	mu       *sync.RWMutex
	byID     cache.Cache[string, *repository.Instance]
	byDomain cache.Cache[string, *repository.Instance]
}

func SetUpInstance(
	cache *Instance,
	handle handler.Handle[*repository.Instance, *repository.Instance],
) handler.Handle[*repository.Instance, *repository.Instance] {
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		instance, err := handle(ctx, instance)
		if err != nil {
			return nil, err
		}

		cache.set(instance, "")
		return instance, nil
	}
}

func SetUpInstanceWithout(cache *Instance) handler.Handle[*repository.Instance, *repository.Instance] {
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		cache.set(instance, "")
		return instance, nil
	}
}

func SetUpInstanceDecorated(
	cache *Instance,
	handle handler.Handle[*repository.Instance, *repository.Instance],
	decorator handler.Decorate[*repository.Instance, *repository.Instance],
) handler.Handle[*repository.Instance, *repository.Instance] {
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		instance, err := handle(ctx, instance)
		if err != nil {
			return nil, err
		}

		return decorator(ctx, instance, func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
			cache.set(instance, "")
			return instance, nil
		})
	}
}

func ForInstanceByID(cache *Instance, handle handler.Handle[string, *repository.Instance]) handler.Handle[string, *repository.Instance] {
	return func(ctx context.Context, id string) (*repository.Instance, error) {
		cache.mu.RLock()

		instance, ok := cache.byID.Get(id)
		cache.mu.RUnlock()
		if ok {
			return instance, nil
		}

		instance, err := handle(ctx, id)
		if err != nil {
			return nil, err

		}

		cache.set(instance, "")
		return instance, nil
	}
}

func ForInstanceByDomain(cache *Instance, handle handler.Handle[string, *repository.Instance]) handler.Handle[string, *repository.Instance] {
	return func(ctx context.Context, domain string) (*repository.Instance, error) {
		cache.mu.RLock()

		instance, ok := cache.byDomain.Get(domain)
		cache.mu.RUnlock()
		if ok {
			return instance, nil
		}

		instance, err := handle(ctx, domain)
		if err != nil {
			return nil, err
		}

		cache.set(instance, domain)
		return instance, nil
	}
}

func (i *Instance) set(instance *repository.Instance, domain string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if domain != "" {
		i.byDomain.Set(domain, instance)
	}
	i.byID.Set(instance.ID, instance)
}
