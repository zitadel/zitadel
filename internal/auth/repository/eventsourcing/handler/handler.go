package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	query2 "github.com/zitadel/zitadel/internal/query"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration time.Duration
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

func Register(configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, systemDefaults sd.SystemDefaults, queries *query2.Queries) []query.Handler {
	return []query.Handler{
		newUser(
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es}, queries),
		newUserSession(
			handler{view, bulkLimit, configs.cycleDuration("UserSession"), errorCount, es}),
		newToken(
			handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount, es}),
		newIDPConfig(
			handler{view, bulkLimit, configs.cycleDuration("IDPConfig"), errorCount, es}),
		newIDPProvider(
			handler{view, bulkLimit, configs.cycleDuration("IDPProvider"), errorCount, es},
			systemDefaults, queries),
		newExternalIDP(
			handler{view, bulkLimit, configs.cycleDuration("ExternalIDP"), errorCount, es},
			systemDefaults, queries),
		newRefreshToken(handler{view, bulkLimit, configs.cycleDuration("RefreshToken"), errorCount, es}),
		newOrgProjectMapping(handler{view, bulkLimit, configs.cycleDuration("OrgProjectMapping"), errorCount, es}),
	}
}

func (configs Configs) cycleDuration(viewModel string) time.Duration {
	c, ok := configs[viewModel]
	if !ok {
		return 3 * time.Minute
	}
	return c.MinimumCycleDuration
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

func withInstanceID(ctx context.Context, instanceID string) context.Context {
	return authz.WithInstanceID(ctx, instanceID)
}
