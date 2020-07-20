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
	queryStmtFormat  = "SELECT 1 FROM %s WHERE object_type = $1 AND (locked_until < now() OR locker_id = $2) FOR UPDATE"
	insertStmtFormat = "INSERT INTO %s (object_type, locker_id, locked_until) VALUES ($1, $2, now()+$3) ON CONFLICT (object_type) DO UPDATE SET locked_until = now()+$4, locker_id = $5 WHERE locks.object_type = $6"
)

type lock struct {
	LockerID    string    `gorm:"column:locker_id;primary_key"`
	LockedUntil time.Time `gorm:"column:locked_until"`
	ViewName    string    `gorm:"column:object_type;primary_key"`
}

func Renew(dbClient *sql.DB, lockTable, lockerID, viewModel string, waitTime time.Duration) error {
	return crdb.ExecuteTx(context.Background(), dbClient, nil, func(tx *sql.Tx) error {
		query := fmt.Sprintf(queryStmtFormat, lockTable)
		insert := fmt.Sprintf(insertStmtFormat, lockTable)
		rows, err := tx.Query(query, viewModel, lockerID)
		if err != nil || !rows.Next() {
			rows.Close()
			return err
		}
		err = rows.Close()
		logging.Log("LOCKE-bmpfY").OnError(err).Debug("unable to close rows")

		rs, err := tx.Exec(insert, viewModel, lockerID, waitTime.Seconds(), waitTime.Seconds(), lockerID, viewModel)
		if err != nil {
			tx.Rollback()
			return err
		}
		if rows, _ := rs.RowsAffected(); rows == 0 {
			return caos_errs.ThrowAlreadyExists(nil, "SPOOL-lso0e", "view already locked")
		}
		return nil
	})
}
