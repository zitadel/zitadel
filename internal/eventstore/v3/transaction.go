package eventstore

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (es *Eventstore) WithInstanceLock(ctx context.Context, change func(context.Context) error) error {
	for {
		tx, closeTx, err := es.pushTx(ctx, es.client)
		if err != nil {
			return err
		}
		var acquired bool
		err = tx.QueryRow(ctx, "SELECT pg_try_advisory_xact_lock(hashtext('zitadel.group-lifecycle'), hashtext($1))", authz.GetInstance(ctx).InstanceID()).Scan(&acquired)
		if err != nil {
			return closeTx(err)
		}
		if acquired {
			sqlTx, ok := tx.(*sql.Transaction)
			if !ok {
				return closeTx(zerrors.ThrowInternal(nil, "EVENTSTORE-tx-type", "Unsupported transaction type"))
			}
			return closeTx(change(eventstore.ContextWithTransaction(ctx, sqlTx.Tx)))
		}
		// Do not occupy pool connections while waiting for another replica.
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		timer := time.NewTimer(20 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
