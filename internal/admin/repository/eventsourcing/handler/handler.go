package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/eventstore"
	handler2 "github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/static"
)

type Configs map[string]*Config

type Config struct {
	MinimumCycleDuration time.Duration
}

func Register(ctx context.Context, configs Configs, bulkLimit, errorCount uint64, view *view.View, static static.Storage, es *eventstore.Eventstore) []*handler2.Handler {
	if static == nil {
		return nil
	}
	config := handler2.Config{
		Eventstore:      es,
		BulkLimit:       uint16(bulkLimit),
		MaxFailureCount: uint8(errorCount),
		RequeueEvery:    3 * time.Minute,
	}
	return []*handler2.Handler{
		newStyling(ctx,
			configs.overwrite(config, "Styling"),
			static,
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
