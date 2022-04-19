package spooler

import (
	"database/sql"
	"time"

	es_locker "github.com/caos/zitadel/internal/eventstore/v1/locker"
)

const (
	lockTable = "authz.locks"
)

type locker struct {
	dbClient *sql.DB
}

func (l *locker) Renew(lockerID, viewModel, instanceID string, waitTime time.Duration) error {
	return es_locker.Renew(l.dbClient, lockTable, lockerID, viewModel, instanceID, waitTime)
}
