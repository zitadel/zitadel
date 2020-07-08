package spooler

import (
	"database/sql"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
	"net/http"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentTasks       int
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore, view *view.View, sql *sql.DB, eventstoreRepos handler.EventstoreRepos, systemDefaults sd.SystemDefaults, i18n *i18n.Translator, dir http.FileSystem) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:      es,
		Locker:          &locker{dbClient: sql},
		ConcurrentTasks: c.ConcurrentTasks,
		ViewHandlers:    handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, eventstoreRepos, systemDefaults, i18n, dir),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
