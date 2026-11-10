package handler

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type FieldHandler struct {
	Handler
}

type fieldProjection struct {
	name string
}

// Name implements Projection.
func (f *fieldProjection) Name() string {
	return f.name
}

// Reducers implements Projection.
func (f *fieldProjection) Reducers() []AggregateReducer {
	return nil
}

var _ Projection = (*fieldProjection)(nil)

// NewFieldHandler returns a projection handler which backfills the `eventstore.fields` table with historic events which
// might have existed before they had and Field Operations defined.
// The events are filtered by the mapped aggregate types and each event type for that aggregate.
func NewFieldHandler(config *Config, name string, eventTypes map[eventstore.AggregateType][]eventstore.EventType) *FieldHandler {
	return &FieldHandler{
		Handler: Handler{
			projection:             &fieldProjection{name: name},
			client:                 config.Client,
			es:                     config.Eventstore,
			bulkLimit:              config.BulkLimit,
			eventTypes:             eventTypes,
			requeueEvery:           config.RequeueEvery,
			now:                    time.Now,
			maxFailureCount:        config.MaxFailureCount,
			retryFailedAfter:       config.RetryFailedAfter,
			triggeredInstancesSync: sync.Map{},
			triggerWithoutEvents:   config.TriggerWithoutEvents,
			txDuration:             config.TransactionDuration,
		},
	}
}

// Trigger executes the backfill job of events for the instance currently in the context.
func (h *FieldHandler) Trigger(ctx context.Context, opts ...TriggerOpt) (err error) {
	config := new(triggerConfig)
	for _, opt := range opts {
		opt(config)
	}

	cancel := h.lockInstance(ctx, config)
	if cancel == nil {
		return nil
	}
	defer cancel()

	for i := 0; ; i++ {
		var additionalIteration bool
		var wg sync.WaitGroup
		wg.Add(1)
		queue <- func() {
			additionalIteration, err = h.processEvents(ctx, config)
			wg.Done()
		}
		wg.Wait()
		logging.Debug(ctx, "trigger iteration", "iteration", i)
		if !additionalIteration || err != nil {
			return err
		}
	}
}

func (h *FieldHandler) processEvents(ctx context.Context, config *triggerConfig) (additionalIteration bool, err error) {
	defer func() {
		pgErr := new(pgconn.PgError)
		if errors.As(err, &pgErr) {
			// error returned if the row is currently locked by another connection
			if pgErr.Code == "55P03" {
				logging.WithError(ctx, err).Info("another handler is already updating this projection")
				err = nil
				additionalIteration = false
			}
		}
	}()

	txCtx := ctx
	if h.txDuration > 0 {
		var cancel, cancelTx func()
		// add 100ms to store current state if iteration takes too long
		txCtx, cancelTx = context.WithTimeout(ctx, h.txDuration+100*time.Millisecond)
		defer cancelTx()
		ctx, cancel = context.WithTimeout(ctx, h.txDuration)
		defer cancel()
	}

	tx, err := h.client.BeginTx(txCtx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil && !errors.Is(err, &executionError{}) {
			rollbackErr := tx.Rollback()
			logging.OnError(ctx, rollbackErr).Error("unable to rollback tx")
			return
		}
		commitErr := tx.Commit()
		if err == nil {
			err = commitErr
		}
	}()

	var hasLocked bool
	err = tx.QueryRowContext(ctx, "SELECT pg_try_advisory_xact_lock(hashtext($1), hashtext($2))", h.ProjectionName(), authz.GetInstance(ctx).InstanceID()).Scan(&hasLocked)
	if err != nil {
		return false, err
	}
	if !hasLocked {
		return false, zerrors.ThrowInternal(nil, "V2-xRffO", "projection already locked")
	}

	// always await currently running transactions
	config.awaitRunning = true
	currentState, err := h.currentState(ctx, tx)
	if err != nil {
		return additionalIteration, err
	}
	if !config.maxPosition.IsZero() && currentState.cursor.Position.GreaterThanOrEqual(config.maxPosition) {
		return false, nil
	}

	events, additionalIteration, err := h.fetchEvents(ctx, tx, currentState, config.minPosition)
	if err != nil {
		return additionalIteration, err
	}
	if len(events) == 0 {
		err = h.setState(ctx, tx, currentState)
		return additionalIteration, err
	}

	err = h.es.FillFields(ctx, events...)
	if err != nil {
		return false, err
	}

	err = h.setState(ctx, tx, currentState)

	return additionalIteration, err
}

func (h *FieldHandler) fetchEvents(ctx context.Context, tx *sql.Tx, currentState *state, minPosition decimal.Decimal) (_ []eventstore.FillFieldsEvent, additionalIteration bool, err error) {
	events, err := h.es.Filter(ctx, h.eventQuery(currentState, minPosition).SetTx(tx))
	if err != nil || len(events) == 0 {
		logging.OnError(ctx, err).Debug("filter eventstore failed")
		return nil, false, err
	}

	currentState.applyEvent(events[len(events)-1])
	additionalIteration = len(events) == int(h.bulkLimit)

	fillFieldsEvents := make([]eventstore.FillFieldsEvent, len(events))
	for i, event := range events {
		fillFieldsEvents[i] = event.(eventstore.FillFieldsEvent)
	}

	return fillFieldsEvents, additionalIteration, nil
}
