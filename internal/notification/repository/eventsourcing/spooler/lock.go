package spooler

import (
	"database/sql"
	"time"

	es_locker "github.com/caos/zitadel/internal/eventstore/locker"
)

const (
	lockTable = "notification.locks"
)

type locker struct {
	dbClient *sql.DB
}

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) (lockedUntil time.Time, isLeaseHolder bool, err error) {
	return es_locker.Renew(l.dbClient, lockTable, lockerID, viewModel, waitTime)
}
