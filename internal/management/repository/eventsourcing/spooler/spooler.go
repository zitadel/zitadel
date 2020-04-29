package spooler

import (
	"database/sql"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
)

type SpoolerConfig struct {
	BulkLimit       uint64
	ConcurrentTasks int
	View            *view.View
	Handlers        handler.Configs

	SQL *sql.DB
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:      es,
		Locker:          &locker{dbClient: c.SQL},
		ConcurrentTasks: c.ConcurrentTasks,
		ViewHandlers:    handler.Register(c.Handlers, c.BulkLimit, c.View, es),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
