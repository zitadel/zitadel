package handler

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"

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
		h.log().OnError(err).Info("process events failed")
		h.log().WithField("iteration", i).Debug("trigger iteration")
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
				h.log().Debug("state already locked")
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
			h.log().OnError(rollbackErr).Debug("unable to rollback tx")
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
	// stop execution if currentState.eventTimestamp >= config.maxCreatedAt
	if !config.maxPosition.IsZero() && currentState.position.GreaterThanOrEqual(config.maxPosition) {
		return false, nil
	}

	if config.minPosition.GreaterThan(decimal.NewFromInt(0)) {
		currentState.position = config.minPosition
		currentState.offset = 0
	}

	events, additionalIteration, err := h.fetchEvents(ctx, tx, currentState)
	if err != nil {
		return additionalIteration, err
	}
	if len(events) == 0 {
		err = h.setState(tx, currentState)
		return additionalIteration, err
	}

	err = h.es.FillFields(ctx, events...)
	if err != nil {
		return false, err
	}

	err = h.setState(tx, currentState)

	return additionalIteration, err
}

func (h *FieldHandler) fetchEvents(ctx context.Context, tx *sql.Tx, currentState *state) (_ []eventstore.FillFieldsEvent, additionalIteration bool, err error) {
	events, err := h.es.Filter(ctx, h.eventQuery(currentState).SetTx(tx))
	if err != nil || len(events) == 0 {
		h.log().OnError(err).Debug("filter eventstore failed")
		return nil, false, err
	}
	eventAmount := len(events)

	idx, offset := skipPreviouslyReducedEvents(events, currentState)

	if currentState.position.Equal(events[len(events)-1].Position()) {
		offset += currentState.offset
	}
	currentState.position = events[len(events)-1].Position()
	currentState.offset = offset
	currentState.aggregateID = events[len(events)-1].Aggregate().ID
	currentState.aggregateType = events[len(events)-1].Aggregate().Type
	currentState.sequence = events[len(events)-1].Sequence()
	currentState.eventTimestamp = events[len(events)-1].CreatedAt()

	if idx+1 == len(events) {
		return nil, false, nil
	}
	events = events[idx+1:]

	additionalIteration = eventAmount == int(h.bulkLimit)

	fillFieldsEvents := make([]eventstore.FillFieldsEvent, len(events))
	highestPosition := events[len(events)-1].Position()
	for i, event := range events {
		if event.Position().Equal(highestPosition) {
			offset++
		}
		fillFieldsEvents[i] = event.(eventstore.FillFieldsEvent)
	}

	return fillFieldsEvents, additionalIteration, nil
}

func skipPreviouslyReducedEvents(events []eventstore.Event, currentState *state) (index int, offset uint32) {
	var position decimal.Decimal
	for i, event := range events {
		if !event.Position().Equal(position) {
			offset = 0
			position = event.Position()
		}
		offset++
		if event.Position().Equal(currentState.position) &&
			event.Aggregate().ID == currentState.aggregateID &&
			event.Aggregate().Type == currentState.aggregateType &&
			event.Sequence() == currentState.sequence {
			return i, offset
		}
	}
	return -1, 0
}
