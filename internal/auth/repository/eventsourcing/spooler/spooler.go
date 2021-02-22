package spooler

import (
	"database/sql"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	key_model "github.com/caos/zitadel/internal/key/model"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentWorkers     int
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig, es eventstore.Eventstore, view *view.View, client *sql.DB, systemDefaults sd.SystemDefaults, keyChan chan<- *key_model.KeyView) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:        es,
		Locker:            &locker{dbClient: client},
		ConcurrentWorkers: c.ConcurrentWorkers,
		ViewHandlers:      handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, systemDefaults, keyChan),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
