package eventstore

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type transactionKey struct{}

type transactionState struct {
	tx     *sql.Tx
	events []Event
}

// ContextWithTransaction binds command reads and writes to the same transaction.
func ContextWithTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(database.ContextWithTransaction(ctx, tx), transactionKey{}, &transactionState{tx: tx})
}

func transactionFromContext(ctx context.Context) *transactionState {
	state, _ := ctx.Value(transactionKey{}).(*transactionState)
	return state
}

// WithInstanceLock serializes group lifecycle operations across replicas. Reads,
// cascades, events, reservations and webhook jobs commit or roll back together.
func (es *Eventstore) WithInstanceLock(ctx context.Context, change func(context.Context) error) error {
	if transactionFromContext(ctx) != nil {
		return change(ctx)
	}
	runner, ok := es.pusher.(interface {
		WithInstanceLock(context.Context, func(context.Context) error) error
	})
	if !ok {
		return zerrors.ThrowUnimplemented(nil, "EVENTSTORE-transaction", "Event store does not support transactional group lifecycle operations")
	}
	var state *transactionState
	err := runner.WithInstanceLock(ctx, func(ctx context.Context) error {
		state = transactionFromContext(ctx)
		return change(ctx)
	})
	if err == nil && state != nil {
		es.notify(state.events)
	}
	return err
}
