package locker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

const (
	queryStmtFmt     = "select locker_id = $1, locked_until from %s where view_name = $2"
	insertStmtFormat = "INSERT INTO %s" +
		" (locker_id, locked_until, view_name) VALUES ($1, now()+$2::INTERVAL, $3)" +
		" ON CONFLICT (view_name)" +
		" DO UPDATE SET locker_id = $1, locked_until = now()+$2::INTERVAL" +
		" WHERE locks.view_name = $3 AND (locks.locker_id = $1 OR locks.locked_until < now()) " +
		" RETURNING locked_until"
	millisecondsAsSeconds = int64(time.Second / time.Millisecond)
)

type lock struct {
	LockerID    string    `gorm:"column:locker_id;primary_key"`
	LockedUntil time.Time `gorm:"column:locked_until"`
	ViewName    string    `gorm:"column:view_name;primary_key"`
}

func Renew(dbClient *sql.DB, lockTable, lockerID, viewModel string, waitTime time.Duration) (lockedUntil time.Time, isLeaseHolder bool, err error) {
	err = dbClient.QueryRow(fmt.Sprintf(queryStmtFmt, lockTable), lockerID, viewModel).Scan(&isLeaseHolder, &lockedUntil)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return time.Time{}, false, caos_errs.ThrowInternal(nil, "SPOOL-qcRP8", "query row failed")
	}
	if !isLeaseHolder && lockedUntil.After(time.Now()) {
		return lockedUntil, false, nil
	}
	err = crdb.ExecuteTx(context.Background(), dbClient, nil, func(tx *sql.Tx) error {
		insert := fmt.Sprintf(insertStmtFormat, lockTable)
		err = tx.QueryRow(insert, lockerID, waitTime.Milliseconds()/millisecondsAsSeconds, viewModel).
			Scan(&lockedUntil)

		if err != nil {
			tx.Rollback()
			return err
		}

		isLeaseHolder = true
		logging.LogWithFields("LOCKE-lOgbg", "view", viewModel, "locker", lockerID).Debug("locker changed")
		return nil
	})

	return lockedUntil, isLeaseHolder, err
}
