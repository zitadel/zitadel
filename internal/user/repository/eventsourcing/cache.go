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
	logging.Log("EVENT-vJG2j").OnError(err).Panic("unable to create user cache")

	return &UserCache{userCache: userCache}, nil
}

func (c *UserCache) getUser(id string) *model.User {
	user := &model.User{ObjectRoot: models.ObjectRoot{AggregateID: id}}
	if err := c.userCache.Get(id, user); err != nil {
		logging.Log("EVENT-AtS0S").WithError(err).Debug("error in getting cache")
	}
	return user
}

func (c *UserCache) cacheUser(user *model.User) {
	err := c.userCache.Set(user.AggregateID, user)
	if err != nil {
		logging.Log("EVENT-0V2gX").WithError(err).Debug("error in setting project cache")
	}
}

type CacheUser struct{
	
}