package spooler

import (
	"context"

	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
	"github.com/zitadel/zitadel/internal/static"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentWorkers     int
	ConcurrentInstances   int
	Handlers              handler.Configs
}

func StartSpooler(ctx context.Context, c SpoolerConfig, es v1.Eventstore, esV2 *eventstore.Eventstore, view *view.View, sql *database.DB, static static.Storage) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:          es,
		EventstoreV2:        esV2,
		Locker:              &locker{dbClient: sql.DB},
		ConcurrentWorkers:   c.ConcurrentWorkers,
		ConcurrentInstances: c.ConcurrentInstances,
		ViewHandlers:        handler.Register(ctx, c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, static),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
