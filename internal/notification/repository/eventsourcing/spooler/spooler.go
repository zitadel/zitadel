package spooler

import (
	"database/sql"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"net/http"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/notification/repository/eventsourcing/view"
)

type SpoolerConfig struct {
	BulkLimit             uint64
	FailureCountUntilSkip uint64
	ConcurrentWorkers     int
	Handlers              handler.Configs
}

func StartSpooler(c SpoolerConfig, es v1.Eventstore, view *view.View, sql *sql.DB, command *command.CommandSide, systemDefaults sd.SystemDefaults, i18n *i18n.Translator, dir http.FileSystem) *spooler.Spooler {
	spoolerConfig := spooler.Config{
		Eventstore:        es,
		Locker:            &locker{dbClient: sql},
		ConcurrentWorkers: c.ConcurrentWorkers,
		ViewHandlers:      handler.Register(c.Handlers, c.BulkLimit, c.FailureCountUntilSkip, view, es, command, systemDefaults, i18n, dir),
	}
	spool := spoolerConfig.New()
	spool.Start()
	return spool
}
