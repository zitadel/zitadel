package spooler

import (
	"database/sql"

	"github.com/caos/zitadel/internal/config/systemdefaults"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore, view *view.View, sql *sql.DB, eventstoreRepos handler.EventstoreRepos, defaults systemdefaults.SystemDefaults) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:   es,
		Locker:       &locker{dbClient: sql},
		ViewHandlers: handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, eventstoreRepos, defaults),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
