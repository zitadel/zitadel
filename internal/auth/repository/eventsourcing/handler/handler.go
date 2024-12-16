package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	handler2 "github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/id"
	query2 "github.com/zitadel/zitadel/internal/query"
)

type Config struct {
	Client     *database.DB
	Eventstore *eventstore.Eventstore

	BulkLimit             uint64
	FailureCountUntilSkip uint64
	TransactionDuration   time.Duration
	Handlers              map[string]*ConfigOverwrites

	ActiveInstancer interface {
		ActiveInstances() []string
	}
}

type ConfigOverwrites struct {
	MinimumCycleDuration time.Duration
}

var projections []*handler.Handler

func Register(ctx context.Context, configs Config, view *view.View, queries *query2.Queries) {
	// make sure the slice does not contain old values
	projections = nil

	projections = append(projections, newUser(ctx,
		configs.overwrite("User"),
		view,
		queries,
	))

	projections = append(projections, newUserSession(ctx,
		configs.overwrite("UserSession"),
		view,
		queries,
		id.SonyFlakeGenerator(),
	))

	projections = append(projections, newToken(ctx,
		configs.overwrite("Token"),
		view,
	))

	projections = append(projections, newRefreshToken(ctx,
		configs.overwrite("RefreshToken"),
		view,
	))
}

func Start(ctx context.Context) {
	for _, projection := range projections {
		projection.Start(ctx)
	}
}

func Projections() []*handler2.Handler {
	return projections
}

func ProjectInstance(ctx context.Context) error {
	for _, projection := range projections {
		_, err := projection.Trigger(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (config Config) overwrite(viewModel string) handler2.Config {
	c := handler2.Config{
		Client:              config.Client,
		Eventstore:          config.Eventstore,
		BulkLimit:           uint16(config.BulkLimit),
		RequeueEvery:        3 * time.Minute,
		MaxFailureCount:     uint8(config.FailureCountUntilSkip),
		TransactionDuration: config.TransactionDuration,
		ActiveInstancer:     config.ActiveInstancer,
	}
	overwrite, ok := config.Handlers[viewModel]
	if !ok {
		return c
	}
	if overwrite.MinimumCycleDuration > 0 {
		c.RequeueEvery = overwrite.MinimumCycleDuration
	}
	return c
}
