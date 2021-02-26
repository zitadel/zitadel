package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
)

type IAMCache struct {
	iamCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*IAMCache, error) {
	iamCache, err := conf.Config.NewCache()
	logging.Log("EVENT-9siew").OnError(err).Panic("unable to create iam cache")

	return &IAMCache{iamCache: iamCache}, nil
}

func (c *IAMCache) getIAM(ID string) *model.IAM {
	user := &model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: ID}}
	if err := c.iamCache.Get(ID, user); err != nil {
		logging.Log("EVENT-slo9x").WithError(err).Debug("error in getting cache")
	}
	return user
}

func (c *IAMCache) cacheIAM(iam *model.IAM) {
	err := c.iamCache.Set(iam.AggregateID, iam)
	if err != nil {
		logging.Log("EVENT-os03w").WithError(err).Debug("error in setting iam cache")
	}
}
