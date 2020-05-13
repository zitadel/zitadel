package spooler

import (
	"database/sql"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentTasks       int
	Handlers              handler.Configs
}

type EventstoreRepos struct {
	UserEvents *usr_event.UserEventstore
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore, view *view.View, sql *sql.DB, eventstoreRepos handler.EventstoreRepos) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:      es,
		Locker:          &locker{dbClient: sql},
		ConcurrentTasks: c.ConcurrentTasks,
		ViewHandlers:    handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, eventstoreRepos),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
