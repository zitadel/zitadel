package spooler

import (
	"database/sql"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	key_model "github.com/caos/zitadel/internal/key/model"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentWorkers     int
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore, view *view.View, client *sql.DB, repos handler.EventstoreRepos, systemDefaults sd.SystemDefaults, keyChan chan<- *key_model.KeyView) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:        es,
		Locker:            &locker{dbClient: client},
		ConcurrentWorkers: c.ConcurrentWorkers,
		ViewHandlers:      handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, repos, systemDefaults, keyChan),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
