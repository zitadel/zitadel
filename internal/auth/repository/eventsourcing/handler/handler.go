package handler

import (
	"time"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration types.Duration
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, repos EventstoreRepos) []spooler.Handler {
	return []spooler.Handler{
		&User{handler: handler{view, bulkLimit, configs.cycleDuration("User"), errorCount}},
		&UserSession{handler: handler{view, bulkLimit, configs.cycleDuration("UserSession"), errorCount}, userEvents: repos.UserEvents},
		&Token{handler: handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount}},
		&Key{handler: handler{view, bulkLimit, configs.cycleDuration("Key"), errorCount}},
		&Application{handler: handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount}},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}
