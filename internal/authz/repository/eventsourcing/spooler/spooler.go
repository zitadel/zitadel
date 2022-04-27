package spooler

import (
	"database/sql"

	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"

	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"

	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/handler"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"

	"github.com/zitadel/zitadel/internal/eventstore/v1/spooler"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentWorkers     int
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig, es v1.Eventstore, view *view.View, sql *sql.DB, systemDefaults sd.SystemDefaults) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:        es,
		Locker:            &locker{dbClient: sql},
		ConcurrentWorkers: c.ConcurrentWorkers,
		ViewHandlers:      handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, systemDefaults),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
