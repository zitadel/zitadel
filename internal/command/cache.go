package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector"
)

type Caches struct {
	milestones cache.Cache[milestoneIndex, string, *MilestonesReached]
}

func startCaches(background context.Context, connectors connector.Connectors) (_ *Caches, err error) {
	caches := new(Caches)
	caches.milestones, err = connector.StartCache[milestoneIndex, string, *MilestonesReached](background, []milestoneIndex{milestoneIndexInstanceID}, cache.PurposeMilestones, connectors.Config.Milestones, connectors)
	if err != nil {
		return nil, err
	}
	return caches, nil
}
