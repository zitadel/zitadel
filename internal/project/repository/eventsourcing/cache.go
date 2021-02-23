package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
)

type ProjectCache struct {
	projectCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*ProjectCache, error) {
	projectCache, err := conf.Config.NewCache()
	logging.Log("EVENT-CsHdo").OnError(err).Panic("unable to create project cache")

	return &ProjectCache{projectCache: projectCache}, nil
}

func (c *ProjectCache) getProject(ID string) (project *model.Project) {
	project = &model.Project{ObjectRoot: models.ObjectRoot{AggregateID: ID}}
	if err := c.projectCache.Get(ID, project); err != nil {
		logging.Log("EVENT-tMydV").WithError(err).Debug("error in getting cache")
	}
	return project
}

func (c *ProjectCache) cacheProject(project *model.Project) {
	err := c.projectCache.Set(project.AggregateID, project)
	if err != nil {
		logging.Log("EVENT-3wKzj").WithError(err).Debug("error in setting project cache")
	}
}
