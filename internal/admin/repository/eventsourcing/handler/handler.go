package handler

import (
	"time"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/query"
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, repos EventstoreRepos) []query.Handler {
	return []query.Handler{
		&Org{handler: handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount}},
		&IamMember{handler: handler{view, bulkLimit, configs.cycleDuration("IamMember"), errorCount}, userEvents: repos.UserEvents},
		&IdpConfig{handler: handler{view, bulkLimit, configs.cycleDuration("IdpConfig"), errorCount}},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
