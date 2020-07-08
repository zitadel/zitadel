package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

type UserGrantCache struct {
	userGrantCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*UserGrantCache, error) {
	userGrantCache, err := conf.Config.NewCache()
	logging.Log("EVENT-8EhUZ").OnError(err).Panic("unable to create user grant cache")

	return &UserGrantCache{userGrantCache: userGrantCache}, nil
}

func (c *UserGrantCache) getUserGrant(ID string) *model.UserGrant {
	user := &model.UserGrant{ObjectRoot: models.ObjectRoot{AggregateID: ID}}
	err := c.userGrantCache.Get(ID, user)
	logging.Log("EVENT-QAd7T").OnError(err).Debug("error in getting cache")

	return user
}

func (c *UserGrantCache) cacheUserGrant(grant *model.UserGrant) {
	err := c.userGrantCache.Set(grant.AggregateID, grant)

	logging.Log("EVENT-w2KNQ").OnError(err).Debug("error in setting user grant cache")
}
