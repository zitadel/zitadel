package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"time"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDurationMillisecond int
}

type handler struct {
	view                *view.View
	bulkLimit           uint64
	cycleDuration       time.Duration
	errorCountUntilSkip uint64
}

type EventstoreRepos struct {
	UserEvents *usr_event.UserEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, eventstore eventstore.Eventstore, repos EventstoreRepos) []spooler.Handler {
	return []spooler.Handler{
		&NotifyUser{handler: handler{view, bulkLimit, configs.cycleDuration("User"), errorCount}},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return time.Duration(c.MinimumCycleDurationMillisecond) * time.Millisecond
}
