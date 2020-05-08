package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
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
	ProjectEvents *proj_event.ProjectEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, eventstore eventstore.Eventstore, repos EventstoreRepos) []spooler.Handler {
	return []spooler.Handler{
		&GrantedProject{handler: handler{view, bulkLimit, configs.cycleDuration("GrantedProject"), errorCount}, eventstore: eventstore, projectEvents: repos.ProjectEvents},
		&ProjectRole{handler: handler{view, bulkLimit, configs.cycleDuration("ProjectRole"), errorCount}, projectEvents: repos.ProjectEvents},
		&ProjectMember{handler: handler{view, bulkLimit, configs.cycleDuration("ProjectMember"), errorCount}},
		&ProjectGrantMember{handler: handler{view, bulkLimit, configs.cycleDuration("ProjectGrantMember"), errorCount}},
		&Application{handler: handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount}},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return time.Duration(c.MinimumCycleDurationMillisecond) * time.Millisecond
}
