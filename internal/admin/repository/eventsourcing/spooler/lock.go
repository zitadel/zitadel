package spooler

import (
	"database/sql"
	"time"

	es_locker "github.com/caos/zitadel/internal/eventstore/locker"
)

const (
	lockTable      = "adminapi.locks"
	lockedUntilKey = "locked_until"
	lockerIDKey    = "locker_id"
	objectTypeKey  = "object_type"
)

type locker struct {
	dbClient *sql.DB
}

type lock struct {
	LockerID    string    `gorm:"column:locker_id;primary_key"`
	LockedUntil time.Time `gorm:"column:locked_until"`
	ViewName    string    `gorm:"column:object_type;primary_key"`
}

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) (lockedUntil time.Time, isLeaseHolder bool, err error) {
	return es_locker.Renew(l.dbClient, lockTable, lockerID, viewModel, waitTime)
}
