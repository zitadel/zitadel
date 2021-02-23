package spooler

import (
	"database/sql"
	es_locker "github.com/caos/zitadel/internal/eventstore/v1/locker"
	"time"
)

const (
	lockTable = "authz.locks"
)

type locker struct {
	dbClient *sql.DB
}

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) error {
	return es_locker.Renew(l.dbClient, lockTable, lockerID, viewModel, waitTime)
}
