package handler

import (
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	"time"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDurationMillisecond int
}

type handler struct {
	view          *view.View
	bulkLimit     uint64
	cycleDuration time.Duration
}

func Register(configs Configs, bulkLimit uint64, view *view.View, esClient es_api.EventstoreServiceClient) []spooler.Handler {
	return []spooler.Handler{
		&Org{handler{view, bulkLimit, configs.cycleDuration("Org")}, esClient},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return time.Duration(c.MinimumCycleDurationMillisecond) * time.Millisecond
}
