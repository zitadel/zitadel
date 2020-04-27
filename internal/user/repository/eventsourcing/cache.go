package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
)

type UserCache struct {
	userCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*UserCache, error) {
	userCache, err := conf.Config.NewCache()
	logging.Log("EVENT-vDneN").OnError(err).Panic("unable to create user cache")

	return &UserCache{userCache: userCache}, nil
}

func (c *UserCache) getUser(ID string) *model.User {
	user := &model.User{ObjectRoot: models.ObjectRoot{AggregateID: ID}}
	if err := c.userCache.Get(ID, user); err != nil {
		logging.Log("EVENT-4eTZh").WithError(err).Debug("error in getting cache")
	}
	return user
}

func (c *UserCache) cacheUser(user *model.User) {
	err := c.userCache.Set(user.AggregateID, user)
	if err != nil {
		logging.Log("EVENT-ThnBb").WithError(err).Debug("error in setting project cache")
	}
}
