package cached

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/cache"
)

type Instance struct {
	cache.Cache[repository.InstanceIndex, string, *repository.Instance]
}

func NewInstance(c cache.Cache[repository.InstanceIndex, string, *repository.Instance]) *Instance {
	return &Instance{c}
}

func (i *Instance) ByID(ctx context.Context, id string) *repository.Instance {
	log.Println("cached.instance.byID")
	instance, _ := i.Cache.Get(ctx, repository.InstanceByID, id)
	return instance
}

func (i *Instance) ByDomain(ctx context.Context, domain string) *repository.Instance {
	log.Println("cached.instance.byDomain")
	instance, _ := i.Cache.Get(ctx, repository.InstanceByDomain, domain)
	return instance
}

func (i *Instance) Set(ctx context.Context, instance *repository.Instance) {
	log.Println("cached.instance.set")
	i.Cache.Set(ctx, instance)
}
