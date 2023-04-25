package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Config struct {
	Client     *database.DB
	Eventstore *eventstore.Eventstore

	BulkLimit             uint16
	RequeueEvery          time.Duration
	HandleActiveInstances time.Duration
	MaxFailureCount       uint8
}

type Handler struct {
	client     *database.DB
	projection Projection

	es         *eventstore.Eventstore
	bulkLimit  uint16
	aggregates []eventstore.AggregateType

	maxFailureCount uint8

	requeueEvery          time.Duration
	handleActiveInstances time.Duration
	now                   nowFunc
}

// nowFunc makes [time.Now] mockable
type nowFunc func() time.Time

type Projection interface {
	Name() string
	Reducers() []AggregateReducer
}

func NewHandler(
	ctx context.Context,
	config *Config,
	projection Projection,
) *Handler {
	aggregates := make([]eventstore.AggregateType, len(projection.Reducers()))
	for i, reducer := range projection.Reducers() {
		aggregates[i] = reducer.Aggregate
	}

	handler := &Handler{
		projection:            projection,
		client:                config.Client,
		es:                    config.Eventstore,
		bulkLimit:             config.BulkLimit,
		aggregates:            aggregates,
		requeueEvery:          config.RequeueEvery,
		handleActiveInstances: config.HandleActiveInstances,
		now:                   time.Now,
		maxFailureCount:       config.MaxFailureCount,
	}

	return handler
}

func (h *Handler) Start(ctx context.Context) {
	go h.schedule(ctx)
	go h.subscribe(ctx)
}

func (h *Handler) schedule(ctx context.Context) {
	// if there was no run before trigger instantly
	t := time.NewTimer(0)
	didInitialize := h.didProjectionInitialize(ctx)
	if didInitialize {
		t.Reset(h.requeueEvery)
	}

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			instances, err := h.queryInstances(ctx, didInitialize)
			h.log().OnError(err).Debug("unable to query instances")

			var instanceFailed bool
			for _, instance := range instances {
				instanceCtx := authz.WithInstanceID(ctx, instance)
				err = h.Trigger(instanceCtx)
				instanceFailed = instanceFailed || err != nil
				h.log().WithField("instance", instance).OnError(err).Info("scheduled trigger failed")
			}

			if !didInitialize && !instanceFailed {
				err = h.setSucceededOnce(ctx)
				h.log().OnError(err).Debug("unable to set succeeded once")
				didInitialize = err != nil
			}

			t.Reset(h.requeueEvery)
		}
	}
}

func (h *Handler) subscribe(ctx context.Context) {
	queue := make(chan eventstore.Event)
	subscription := eventstore.SubscribeAggregates(queue, h.aggregates...)
	for {
		select {
		case <-ctx.Done():
			subscription.Unsubscribe()
			h.log().Debug("shutdown")
		case event := <-queue:
			events := checkAdditionalEvents(queue, event)
			solvedInstances := make([]string, 0, len(events))
			for _, event := range events {
				if instanceSolved(solvedInstances, event.Aggregate().InstanceID) {
					continue
				}
				ctx := authz.WithInstanceID(ctx, event.Aggregate().InstanceID)
				err := h.Trigger(ctx)
				h.log().OnError(err).Debug("trigger of queued event failed")
			}
		}
	}
}

func instanceSolved(solvedInstances []string, instanceID string) bool {
	for _, solvedInstance := range solvedInstances {
		if solvedInstance == instanceID {
			return true
		}
	}
	return false
}

func checkAdditionalEvents(eventQueue chan eventstore.Event, event eventstore.Event) []eventstore.Event {
	events := make([]eventstore.Event, 1)
	events[0] = event
	for {
		wait := time.NewTimer(1 * time.Millisecond)
		select {
		case event := <-eventQueue:
			events = append(events, event)
		case <-wait.C:
			return events
		}
	}
}

func (h *Handler) queryInstances(ctx context.Context, didInitialize bool) ([]string, error) {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
		AllowTimeTravel().
		AddQuery().
		ExcludedInstanceID("")
	if didInitialize {
		query = query.
			CreationDateAfter(h.now().Add(-1 * h.handleActiveInstances))
	}
	return h.es.InstanceIDs(ctx, query.Builder())
}

func (h *Handler) Trigger(ctx context.Context) (err error) {
	tx, err := h.client.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			h.log().OnError(err).Debug("commit should still work")
		}
		commitErr := tx.Commit()
		h.log().OnError(commitErr).Debug("commit failed")
		if err == nil {
			err = commitErr
		}
	}()

	currentState, shouldSkip, err := h.currentState(ctx, tx)
	if err != nil || shouldSkip {
		return err
	}

	var hasChanged bool
	for {
		events, err := h.es.Filter(ctx, h.eventQuery(ctx, tx, currentState))
		if err != nil {
			h.log().WithError(err).Debug("filter eventstore failed")
			return err
		}
		events = skipPreviouslyReduced(events, currentState)
		if len(events) == 0 {
			break
		}

		statements, err := h.eventsToStatements(tx, events, currentState)
		if err != nil {
			return err
		}

		if err = h.execute(ctx, tx, currentState, statements); err != nil {
			return err
		}

		hasChanged = true
		currentState.aggregateID = events[len(events)-1].Aggregate().ID
		currentState.aggregateType = events[len(events)-1].Aggregate().Type
		currentState.eventSequence = events[len(events)-1].Sequence()
		currentState.eventTimestamp = events[len(events)-1].CreationDate()

		if len(events) < int(h.bulkLimit) {
			break
		}
	}

	if !hasChanged {
		h.updateLastUpdated(ctx, tx, currentState)
		return nil
	}

	return h.setState(ctx, tx, currentState)
}

func skipPreviouslyReduced(events []eventstore.Event, currentState *state) []eventstore.Event {
	for i, event := range events {
		if currentState.aggregateID == event.Aggregate().ID &&
			currentState.aggregateType == event.Aggregate().Type &&
			currentState.eventSequence == event.Sequence() {
			return events[i+1:]
		}
	}
	return events
}

func (h *Handler) execute(ctx context.Context, tx *sql.Tx, currentState *state, statements []*Statement) error {
	for _, statement := range statements {
		_, err := tx.Exec("SAVEPOINT exec")
		if err != nil {
			h.log().WithError(err).Debug("create savepoint failed")
			return err
		}
		if err := statement.Execute(tx, h.projection.Name()); err != nil {
			h.log().WithError(err).Error("statement execution failed")

			_, savepointErr := tx.Exec("ROLLBACK TO SAVEPOINT exec")
			if savepointErr != nil {
				h.log().WithError(savepointErr).Debug("rollback savepoint failed")
				return savepointErr
			}

			if h.handleFailedStmt(tx, currentState, failureFromStatement(statement, err)) {
				continue
			}

			return err
		}
		if _, err = tx.Exec("RELEASE SAVEPOINT exec"); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) eventQuery(ctx context.Context, tx *sql.Tx, currentState *state) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		Limit(uint64(h.bulkLimit)).
		AllowTimeTravel().
		OrderAsc().
		SetTx(tx).
		AddQuery().
		AggregateTypes(h.aggregates...).
		CreationDateAfter(currentState.eventTimestamp.Add(-1 * time.Microsecond)).
		Builder()
}
