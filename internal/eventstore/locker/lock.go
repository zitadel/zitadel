package locker

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

const (
	insertStmtFormat = "INSERT INTO %s" +
		" (locker_id, locked_until, object_type) VALUES ($1, now()+$2, $3)" +
		" ON CONFLICT (object_type)" +
		" DO UPDATE SET locker_id = $4, locked_until = now()+$5" +
		" WHERE locks.object_type = $6 AND (locks.locker_id = $7 AND locks.locked_until >= now() OR locks.locked_until < now())"
)

type lock struct {
	LockerID    string    `gorm:"column:locker_id;primary_key"`
	LockedUntil time.Time `gorm:"column:locked_until"`
	ViewName    string    `gorm:"column:object_type;primary_key"`
}

func Renew(dbClient *sql.DB, lockTable, lockerID, viewModel string, waitTime time.Duration) error {
	return crdb.ExecuteTx(context.Background(), dbClient, nil, func(tx *sql.Tx) error {
		insert := fmt.Sprintf(insertStmtFormat, lockTable)
		result, err := tx.Exec(insert,
			lockerID, waitTime.Seconds(), viewModel,
			lockerID, waitTime.Seconds(),
			viewModel, lockerID)

		if err != nil {
			tx.Rollback()
			return err
		}
		if rows, _ := result.RowsAffected(); rows == 0 {
			return caos_errs.ThrowAlreadyExists(nil, "SPOOL-lso0e", "view already locked")
		}
		logging.LogWithFields("LOCKE-lOgbg", "view", viewModel, "locker", lockerID).WithField("ts", time.Now().Format(time.StampMicro)).Debug("locker changed")
		return nil
	})
}
