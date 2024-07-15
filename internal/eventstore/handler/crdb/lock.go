package crdb

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	lockStmtFormat = "INSERT INTO %[1]s" +
		" (locker_id, locked_until, projection_name, instance_id) VALUES %[2]s" +
		" ON CONFLICT (projection_name, instance_id)" +
		" DO UPDATE SET locker_id = $1, locked_until = now()+$2::INTERVAL" +
		" WHERE %[1]s.projection_name = $3 AND %[1]s.instance_id = ANY ($%[3]d) AND (%[1]s.locker_id = $1 OR %[1]s.locked_until < now())"
)

type Locker interface {
	Lock(ctx context.Context, lockDuration time.Duration, instanceIDs ...string) <-chan error
	Unlock(instanceIDs ...string) error
}

type locker struct {
	client         *sql.DB
	lockStmt       func(values string, instances int) string
	workerName     string
	projectionName string
}

func NewLocker(client *sql.DB, lockTable, projectionName string) Locker {
	workerName, err := id.SonyFlakeGenerator().Next()
	logging.OnError(err).Panic("unable to generate lockID")
	return &locker{
		client: client,
		lockStmt: func(values string, instances int) string {
			return fmt.Sprintf(lockStmtFormat, lockTable, values, instances)
		},
		workerName:     workerName,
		projectionName: projectionName,
	}
}

func (h *locker) Lock(ctx context.Context, lockDuration time.Duration, instanceIDs ...string) <-chan error {
	errs := make(chan error)
	go h.handleLock(ctx, errs, lockDuration, instanceIDs...)
	return errs
}

func (h *locker) handleLock(ctx context.Context, errs chan error, lockDuration time.Duration, instanceIDs ...string) {
	renewLock := time.NewTimer(0)
	for {
		select {
		case <-renewLock.C:
			errs <- h.renewLock(ctx, lockDuration, instanceIDs...)
			// refresh the lock 500ms before it times out. 500ms should be enough for one transaction
			renewLock.Reset(lockDuration - (500 * time.Millisecond))
		case <-ctx.Done():
			close(errs)
			renewLock.Stop()
			return
		}
	}
}

func (h *locker) renewLock(ctx context.Context, lockDuration time.Duration, instanceIDs ...string) error {
	lockStmt, values := h.lockStatement(lockDuration, instanceIDs)
	res, err := h.client.ExecContext(ctx, lockStmt, values...)
	if err != nil {
		return zerrors.ThrowInternal(err, "CRDB-uaDoR", "unable to execute lock")
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return zerrors.ThrowAlreadyExists(nil, "CRDB-mmi4J", "projection already locked")
	}
	return nil
}

func (h *locker) Unlock(instanceIDs ...string) error {
	lockStmt, values := h.lockStatement(0, instanceIDs)
	_, err := h.client.Exec(lockStmt, values...)
	if err != nil {
		return zerrors.ThrowUnknown(err, "CRDB-JjfwO", "unlock failed")
	}
	return nil
}

func (h *locker) lockStatement(lockDuration time.Duration, instanceIDs database.TextArray[string]) (string, []interface{}) {
	valueQueries := make([]string, len(instanceIDs))
	values := make([]interface{}, len(instanceIDs)+4)
	values[0] = h.workerName
	// the unit of crdb interval is seconds (https://www.cockroachlabs.com/docs/stable/interval.html).
	values[1] = lockDuration
	values[2] = h.projectionName
	for i, instanceID := range instanceIDs {
		valueQueries[i] = "($1, now()+$2::INTERVAL, $3, $" + strconv.Itoa(i+4) + ")"
		values[i+3] = instanceID
	}
	values[len(values)-1] = instanceIDs
	return h.lockStmt(strings.Join(valueQueries, ", "), len(values)), values
}
