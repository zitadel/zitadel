package spooler

import (
	"database/sql"
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

//
//func StartSpooler(c SpoolerConfig) *spooler.Spooler {
//	spoolerConfig := spooler.Config{
//		Client:          c.EsClient,
//		Locker:          &locker{dbClient: c.SQL},
//		ConcurrentTasks: c.ConcurrentTasks,
//		ViewHandlers:    handler.Register(c.Handlers, c.BulkLimit, c.View, c.EsClient),
//	}
//	spool := spoolerConfig.New()
//	spool.Start()
//	return spool
//}
