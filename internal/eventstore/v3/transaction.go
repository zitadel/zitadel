package eventstore

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// cmd/setup/78.go takes the same lock with pg_advisory_xact_lock.
const instanceLockStmt = "SELECT pg_try_advisory_xact_lock(hashtextextended('zitadel.group-lifecycle:' || $1, 0))"

var (
	// maxInstanceLockWait fails with DeadlineExceeded, not ErrConcurrentModification,
	// so callers do not retry the whole wait.
	maxInstanceLockWait  = 30 * time.Second
	slowInstanceLockWait = time.Second
)

func (es *Eventstore) WithInstanceLock(ctx context.Context, change func(context.Context) error) error {
	instanceID := authz.GetInstance(ctx).InstanceID()
	start := time.Now()
	warned := false
	for attempt := 0; ; attempt++ {
		tx, closeTx, err := es.pushTx(ctx, es.client)
		if err != nil {
			return err
		}
		var acquired bool
		err = tx.QueryRow(ctx, instanceLockStmt, instanceID).Scan(&acquired)
		if err != nil {
			return closeTx(err)
		}
		if acquired {
			return runLocked(ctx, tx, closeTx, change)
		}
		// Do not occupy pool connections while waiting for another replica.
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		waited := time.Since(start)
		if waited >= maxInstanceLockWait {
			return zerrors.ThrowDeadlineExceededf(nil, "V3-lock-wait", "instance lock not acquired within %s", maxInstanceLockWait)
		}
		if !warned && waited >= slowInstanceLockWait {
			warned = true
			logging.WithFields("instance", instanceID, "waited", waited).Warn("waiting for group lifecycle instance lock")
		}
		timer := time.NewTimer(min(lockRetryDelay(attempt), maxInstanceLockWait-waited))
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func runLocked(ctx context.Context, tx database.Transaction, closeTx func(error) error, change func(context.Context) error) error {
	sqlTx, ok := tx.(*sql.Transaction)
	if !ok {
		return closeTx(zerrors.ThrowInternal(nil, "EVENTSTORE-tx-type", "Unsupported transaction type"))
	}
	returned := false
	// A panic in change must not keep the transaction and lock open until ctx ends.
	defer func() {
		if !returned {
			_ = tx.Rollback(context.WithoutCancel(ctx))
		}
	}()
	err := change(eventstore.ContextWithTransaction(ctx, sqlTx.Tx))
	returned = true
	return closeTx(err)
}

const (
	minLockRetryDelay = 5 * time.Millisecond
	maxLockRetryDelay = 100 * time.Millisecond
)

// lockRetryDelay doubles from 5ms up to 100ms and jitters each attempt so
// waiters do not retry in lockstep.
func lockRetryDelay(attempt int) time.Duration {
	delay := maxLockRetryDelay
	if attempt < 5 {
		delay = min(minLockRetryDelay<<attempt, maxLockRetryDelay)
	}
	return delay/2 + rand.N(delay/2+1)
}
