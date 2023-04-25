package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/eventstore"
	handler2 "github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	query2 "github.com/zitadel/zitadel/internal/query"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration time.Duration
}

func Register(ctx context.Context, configs Configs, bulkLimit, errorCount uint64, view *view.View, es *eventstore.Eventstore, queries *query2.Queries) []*handler2.Handler {
	config := handler2.Config{
		Eventstore:      es,
		BulkLimit:       uint16(bulkLimit),
		MaxFailureCount: uint8(errorCount),
		RequeueEvery:    3 * time.Minute,
	}
	return []*handler2.Handler{
		newUser(ctx,
			configs.overwrite(config, "User"),
			view,
			queries,
		),
		newUserSession(ctx,
			configs.overwrite(config, "UserSession"),
			view,
			queries,
		),
		newToken(ctx,
			configs.overwrite(config, "Token"),
			view,
		),
		newRefreshToken(ctx,
			configs.overwrite(config, "RefreshToken"),
			view,
		),
	}
}

func (configs Configs) overwrite(config handler2.Config, viewModel string) handler2.Config {
	c, ok := configs[viewModel]
	if !ok {
		return config
	}
	if c.MinimumCycleDuration > 0 {
		config.RequeueEvery = c.MinimumCycleDuration
	}
	return config
}
