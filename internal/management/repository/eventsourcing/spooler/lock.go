package spooler

import (
	"database/sql"
	"time"
)

const (
	lockTable = "management.locks"
)

type locker struct {
	dbClient *sql.DB
}

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) error {
	return nil
}
