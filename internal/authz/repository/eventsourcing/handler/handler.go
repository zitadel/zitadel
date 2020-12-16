package handler

import (
	"time"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/query"
	iam_events "github.com/caos/zitadel/internal/iam/repository/eventsourcing"

	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
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

	es eventstore.Eventstore
}

func (h *handler) Eventstore() eventstore.Eventstore {
	return h.es
}

type EventstoreRepos struct {
	IAMEvents *iam_events.IAMEventstore
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es eventstore.Eventstore, repos EventstoreRepos, systemDefaults sd.SystemDefaults) []query.Handler {
	return []query.Handler{
		&UserGrant{
			handler:    handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount, es},
			eventstore: es,
			iamID:      systemDefaults.IamID,
			iamEvents:  repos.IAMEvents,
		},
		&Application{handler: handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount, es}},
		&Org{handler: handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}},
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
