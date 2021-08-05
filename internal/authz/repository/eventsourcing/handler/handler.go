package handler

import (
	"time"

	v1 "github.com/caos/zitadel/internal/eventstore/v1"

	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/eventstore/v1/query"
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

	es v1.Eventstore
}

func (h *handler) Eventstore() v1.Eventstore {
	return h.es
}

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, systemDefaults sd.SystemDefaults) []query.Handler {
	return []query.Handler{
		newUserGrant(
			handler{view, bulkLimit, configs.cycleDuration("UserGrants"), errorCount, es},
			systemDefaults.IamID),
		newUserMembership(
			handler{view, bulkLimit, configs.cycleDuration("UserMemberships"), errorCount, es}),
		newApplication(
			handler{view, bulkLimit, configs.cycleDuration("Application"), errorCount, es}),
		newOrg(
			handler{view, bulkLimit, configs.cycleDuration("Org"), errorCount, es}),
		newFeatures(
			handler{view, bulkLimit, configs.cycleDuration("Features"), errorCount, es}),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 3 * time.Minute
	}
	return c.MinimumCycleDuration.Duration
}

func (h *handler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *handler) LockDuration() time.Duration {
	return h.cycleDuration / 3
}

func (h *handler) QueryLimit() uint64 {
	return h.bulkLimit
}
