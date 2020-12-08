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
	queryStmtFmt = "with l as (select * from %s where view_name = $1) " +
	"select case when (select count(*) from l) = 0 then true else (select case when (locker_id = $2 or locked_until < now()) then true else false end from l) end"
	insertStmtFormat = "INSERT INTO %s" +
		" (locker_id, locked_until, view_name) VALUES ($1, now()+$2::INTERVAL, $3)" +
		" ON CONFLICT (view_name)" +
		" DO UPDATE SET locker_id = $4, locked_until = now()+$5::INTERVAL" +
		" WHERE locks.view_name = $6 AND (locks.locker_id = $7 OR locks.locked_until < now())"
	millisecondsAsSeconds = int64(time.Second / time.Millisecond)
)

type lock struct {
	LockerID    string    `gorm:"column:locker_id;primary_key"`
	LockedUntil time.Time `gorm:"column:locked_until"`
	ViewName    string    `gorm:"column:view_name;primary_key"`
}

func Renew(dbClient *sql.DB, lockTable, lockerID, viewModel string, waitTime time.Duration) error {
	var shouldInsert bool
	err := dbClient.QueryRow(fmt.Sprintf(queryStmtFmt, lockTable), viewModel, lockerID).Scan(&shouldInsert)
	if err !=nil{
		return caos_errs.ThrowInternal(nil, "SPOOL-qcRP8", "query row failed")
	}
	if !shouldInsert { 
		return caos_errs.ThrowAlreadyExists(nil, "SPOOL-7WTO6", "view already locked")
	}
	return crdb.ExecuteTx(context.Background(), dbClient, nil, func(tx *sql.Tx) error {
		insert := fmt.Sprintf(insertStmtFormat, lockTable)
		result, err := tx.Exec(insert,
			lockerID, waitTime.Milliseconds()/millisecondsAsSeconds, viewModel,
			lockerID, waitTime.Milliseconds()/millisecondsAsSeconds,
			viewModel, lockerID)

		if err != nil {
			tx.Rollback()
			return err
		}

		if rows, _ := result.RowsAffected(); rows == 0 {
			return caos_errs.ThrowAlreadyExists(nil, "SPOOL-lso0e", "view already locked")
		}
		logging.LogWithFields("LOCKE-lOgbg", "view", viewModel, "locker", lockerID).Debug("locker changed")
		return nil
	})
}
