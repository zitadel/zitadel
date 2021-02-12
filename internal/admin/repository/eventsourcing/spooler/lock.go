package spooler

import (
	"database/sql"
	es_locker "github.com/caos/zitadel/internal/eventstore/locker"
	"time"
)

const (
	lockTable = "adminapi.locks"
)

type locker struct {
	dbClient *sql.DB
}

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) error {
	return es_locker.Renew(l.dbClient, lockTable, lockerID, viewModel, waitTime)
}
