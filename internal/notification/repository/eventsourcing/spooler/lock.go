package spooler

import (
	"context"
	"database/sql"
	"fmt"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"time"

	"github.com/cockroachdb/cockroach-go/crdb"
)

const (
	lockTable      = "notification.locks"
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

func (l *locker) Renew(lockerID, viewModel string, waitTime time.Duration) error {
	return crdb.ExecuteTx(context.Background(), l.dbClient, nil, func(tx *sql.Tx) error {
		query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, now()+$3) ON CONFLICT (%s) DO UPDATE SET %s = now()+$4, %s = $5 WHERE (locks.%s < now() OR locks.%s = $6) AND locks.%s = $7",
			lockTable, objectTypeKey, lockerIDKey, lockedUntilKey, objectTypeKey, lockedUntilKey, lockerIDKey, lockedUntilKey, lockerIDKey, objectTypeKey)

		rs, err := tx.Exec(query, viewModel, lockerID, waitTime.Seconds(), waitTime.Seconds(), lockerID, lockerID, viewModel)
		if err != nil {
			tx.Rollback()
			return err
		}
		if rows, _ := rs.RowsAffected(); rows == 0 {
			tx.Rollback()
			return caos_errs.ThrowAlreadyExists(nil, "SPOOL-lso0e", "view already locked")
		}
		return nil
	})
}
