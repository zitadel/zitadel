package spooler

import (
	"database/sql"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentTasks       int
	Handlers              handler.Configs
	EventstoreRepos       handler.EventstoreRepos
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore, sqlClient *sql.DB, view *view.View) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:      es,
		Locker:          &locker{dbClient: sqlClient},
		ConcurrentTasks: c.ConcurrentTasks,
		ViewHandlers:    handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, c.EventstoreRepos),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
