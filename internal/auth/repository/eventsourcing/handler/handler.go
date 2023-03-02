package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	query2 "github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/view/repository"
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

func Register(ctx context.Context, configs Configs, bulkLimit, errorCount uint64, view *view.View, es v1.Eventstore, systemDefaults sd.SystemDefaults, queries *query2.Queries) []query.Handler {
	return []query.Handler{
		newUser(ctx,
			handler{view, bulkLimit, configs.cycleDuration("User"), errorCount, es}, queries),
		newUserSession(ctx,
			handler{view, bulkLimit, configs.cycleDuration("UserSession"), errorCount, es}, queries),
		newToken(ctx,
			handler{view, bulkLimit, configs.cycleDuration("Token"), errorCount, es}),
		newRefreshToken(ctx, handler{view, bulkLimit, configs.cycleDuration("RefreshToken"), errorCount, es}),
		newOrgProjectMapping(ctx, handler{view, bulkLimit, configs.cycleDuration("OrgProjectMapping"), errorCount, es}),
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

func newSearchQuery(sequences []*repository.CurrentSequence, aggregateTypes []models.AggregateType, instanceIDs []string) *models.SearchQuery {
	searchQuery := models.NewSearchQuery()
	for _, instanceID := range instanceIDs {
		var seq uint64
		for _, sequence := range sequences {
			if sequence.InstanceID == instanceID {
				seq = sequence.CurrentSequence
				break
			}
		}
		searchQuery.AddQuery().
			AggregateTypeFilter(aggregateTypes...).
			LatestSequenceFilter(seq).
			InstanceIDFilter(instanceID)
	}
	return searchQuery
}
