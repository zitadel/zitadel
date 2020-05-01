package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
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

func Register(configs Configs, bulkLimit uint64, view *view.View, eventstore eventstore.Eventstore) []spooler.Handler {
	return []spooler.Handler{
		&GrantedProject{handler{view, bulkLimit, configs.cycleDuration("GrantedProject")}, eventstore},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return time.Duration(c.MinimumCycleDurationMillisecond) * time.Millisecond
}
