package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/logging"

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
}

type Handler struct {
	client     *database.DB
	projection Projection

	es         *eventstore.Eventstore
	bulkLimit  uint16
	aggregates []eventstore.AggregateType

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
				// TODO(adlerhurst): are multiple succeed writes a problem?
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
		case <-queue:

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
			rollbackErr := tx.Rollback()
			h.log().OnError(rollbackErr).Debug("unable to rollback trigger")
			return
		}
		err = tx.Commit()
	}()

	currentState, err := h.currentState(ctx, tx)
	if err != nil {
		return err
	}

	// TODO: check if bulk limit exeeded

	events, err := h.es.Filter(ctx, h.eventQuery(ctx, tx, currentState))
	if err != nil || len(events) == 0 {
		return err
	}

	statements, err := h.eventsToStatements(events)
	if err != nil {
		return err
	}

	if err = h.execute(ctx, tx, statements); err != nil {
		return err
	}

	return h.setState(ctx, &state{
		InstanceID:     events[len(events)-1].Aggregate().InstanceID,
		EventTimestamp: events[len(events)-1].CreationDate(),
	}, tx)
}

func (h *Handler) execute(ctx context.Context, tx *sql.Tx, statements []*Statement) error {
	for _, statement := range statements {
		if err := statement.Execute(tx, h.projection.Name()); err != nil {
			// TODO(adlerhurst): failed event
			return err
		}
	}
	return nil
}

func (h *Handler) eventsToStatements(events []eventstore.Event) (statements []*Statement, err error) {
	statements = make([]*Statement, len(events))
	for i, event := range events {
		statements[i], err = h.reduce(event)
		if err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func (h *Handler) reduce(event eventstore.Event) (*Statement, error) {
	for _, reducer := range h.projection.Reducers() {
		if reducer.Aggregate != event.Aggregate().Type {
			continue
		}
		for _, reduce := range reducer.EventRedusers {
			if reduce.Event != event.Type() {
				continue
			}
			return reduce.Reduce(event)
		}
	}
	return NewNoOpStatement(event), nil
}

func (h *Handler) eventQuery(ctx context.Context, tx *sql.Tx, currentState *state) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		Limit(uint64(h.bulkLimit)).
		AllowTimeTravel().
		SetTx(tx).
		AddQuery().
		AggregateTypes(h.aggregates...).
		CreationDateAfter(currentState.EventTimestamp.Add(-1 * time.Microsecond)).
		Builder()
}

func (h *Handler) log() *logging.Entry {
	return logging.WithFields("projection", h.projection.Name())
}
