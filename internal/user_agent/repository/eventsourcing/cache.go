package eventsourcing

import (
	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/repository/eventsourcing/model"
)

type UserAgentCache struct {
	userAgentCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*UserAgentCache, error) {
	userAgentCache, err := conf.Config.NewCache()
	logging.Log("EVENT-df2s2").OnError(err).Panic("unable to create user agent cache")

	return &UserAgentCache{userAgentCache: userAgentCache}, nil
}

func (c *UserAgentCache) getUserAgent(ID string) (userAgent *model.UserAgent, sequence uint64) {
	userAgent = &model.UserAgent{ObjectRoot: models.ObjectRoot{AggregateID: ID}}
	if err := c.userAgentCache.Get(ID, userAgent); err == nil {
		sequence = userAgent.Sequence
	} else {
		logging.Log("EVENT-sd23A").WithError(err).Debug("error in getting cache")
	}
	return userAgent, sequence
}

func (c *UserAgentCache) cacheUserAgent(userAgent *model.UserAgent) {
	err := c.userAgentCache.Set(userAgent.AggregateID, userAgent)
	if err != nil {
		logging.Log("EVENT-ds275").WithError(err).Debug("error in setting user agent cache")
	}
}
