package handler

import (
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/spooler"
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, eventstore eventstore.Eventstore, systemDefaults sd.SystemDefaults) []spooler.Handler {
	return []spooler.Handler{
		&UserGrant{
			handler:    handler{view, bulkLimit, configs.cycleDuration("UserGrant"), errorCount},
			eventstore: eventstore,
			iamID:      systemDefaults.IamID},
		&Application{handler: handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount}},
		&Token{handler: handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount}},
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 1 * time.Second
	}
	return c.MinimumCycleDuration.Duration
}
