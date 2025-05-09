package repository

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/storage/cache"
)

type InstanceCache struct {
	cache.Cache[InstanceIndex, string, *Instance]
}

type InstanceIndex uint8

var InstanceIndices = []InstanceIndex{
	InstanceByID,
	InstanceByDomain,
}

const (
	InstanceByID InstanceIndex = iota
	InstanceByDomain
)

var _ cache.Entry[InstanceIndex, string] = (*Instance)(nil)

// Keys implements [cache.Entry].
func (i *Instance) Keys(index InstanceIndex) (key []string) {
	switch index {
	case InstanceByID:
		return []string{i.ID}
	case InstanceByDomain:
		return []string{i.Name}
	}
	return nil
}

func NewInstanceCache(c cache.Cache[InstanceIndex, string, *Instance]) *InstanceCache {
	return &InstanceCache{c}
}

func (i *InstanceCache) ByID(ctx context.Context, id string) *Instance {
	log.Println("cached.instance.byID")
	instance, _ := i.Cache.Get(ctx, InstanceByID, id)
	return instance
}

func (i *InstanceCache) ByDomain(ctx context.Context, domain string) *Instance {
	log.Println("cached.instance.byDomain")
	instance, _ := i.Cache.Get(ctx, InstanceByDomain, domain)
	return instance
}

func (i *InstanceCache) Set(ctx context.Context, instance *Instance) {
	log.Println("cached.instance.set")
	i.Cache.Set(ctx, instance)
}
