package handler

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgconn"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type EventStore interface {
	InstanceIDs(ctx context.Context, query *eventstore.SearchQueryBuilder) ([]string, error)
	Filter(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error)
	Push(ctx context.Context, cmds ...eventstore.Command) ([]eventstore.Event, error)
}

type Config struct {
	Client     *database.DB
	Eventstore EventStore

	BulkLimit             uint16
	RequeueEvery          time.Duration
	HandleActiveInstances time.Duration
	MaxFailureCount       uint8
}

type Handler struct {
	client     *database.DB
	projection Projection

	es         EventStore
	bulkLimit  uint16
	eventTypes map[eventstore.AggregateType][]eventstore.EventType

	maxFailureCount       uint8
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
	aggregates := make(map[eventstore.AggregateType][]eventstore.EventType, len(projection.Reducers()))
	for _, reducer := range projection.Reducers() {
		eventTypes := make([]eventstore.EventType, len(reducer.EventRedusers))
		for i, eventReducer := range reducer.EventRedusers {
			eventTypes[i] = eventReducer.Event
		}
		aggregates[reducer.Aggregate] = eventTypes
	}

	handler := &Handler{
		projection:            projection,
		client:                config.Client,
		es:                    config.Eventstore,
		bulkLimit:             config.BulkLimit,
		eventTypes:            aggregates,
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
	queue := make(chan eventstore.Event, 100)
	subscription := eventstore.SubscribeEventTypes(queue, h.eventTypes)
	for {
		select {
		case <-ctx.Done():
			subscription.Unsubscribe()
			h.log().Debug("shutdown")
		case event := <-queue:
			events := checkAdditionalEvents(queue, event)
			solvedInstances := make([]string, 0, len(events))
			for _, e := range events {
				if instanceSolved(solvedInstances, e.Aggregate().InstanceID) {
					continue
				}
				ctx := authz.WithInstanceID(ctx, e.Aggregate().InstanceID)
				err := h.Trigger(ctx)
				h.log().OnError(err).Debug("trigger of queued event failed")
				if err == nil {
					solvedInstances = append(solvedInstances, e.Aggregate().InstanceID)
				}
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
	for i := 0; ; i++ {
		additionalIteration, err := h.processEvents(ctx)
		h.log().WithField("iteration", i).Debug("trigger iteration")
		if !additionalIteration || err != nil {
			return err
		}
	}
}

func (h *Handler) processEvents(ctx context.Context) (additionalIteration bool, err error) {
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

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err = crdb.ExecuteTx(ctx, h.client.DB, nil, func(tx *sql.Tx) error {
		currentState, err := h.currentState(ctx, tx)
		if err != nil {
			return err
		}

		events, err := h.es.Filter(ctx, h.eventQuery(ctx, tx, currentState))
		if err != nil {
			h.log().WithError(err).Debug("filter eventstore failed")
			return err
		}
		eventAmount := len(events)
		events = skipPreviouslyReduced(events, currentState)

		if len(events) == 0 {
			h.updateLastUpdated(ctx, tx, currentState)
			return nil
		}

		statements, err := h.eventsToStatements(tx, events, currentState)
		if len(statements) == 0 {
			return err
		}

		if err = h.execute(ctx, tx, currentState, statements); err != nil {
			return err
		}

		currentState.aggregateID = statements[len(statements)-1].AggregateID
		currentState.aggregateType = statements[len(statements)-1].AggregateType
		currentState.eventSequence = statements[len(statements)-1].Sequence
		currentState.eventTimestamp = statements[len(statements)-1].CreationDate

		if err := h.setState(ctx, tx, currentState); err != nil {
			return err
		}

		if len(statements) < len(events) {
			// retry imediatly if statements failed
			additionalIteration = true
			return nil
		}

		additionalIteration = eventAmount == int(h.bulkLimit)
		return nil
	})

	return additionalIteration, err
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
		if statement.Execute == nil {
			continue
		}
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
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		Limit(uint64(h.bulkLimit)).
		AllowTimeTravel().
		OrderAsc().
		SetTx(tx).
		InstanceID(currentState.instanceID)

	for aggregateType, eventTypes := range h.eventTypes {
		builder.
			AddQuery().
			AggregateTypes(aggregateType).
			EventTypes(eventTypes...).
			CreationDateAfter(currentState.eventTimestamp.Add(-1 * time.Microsecond)).
			Builder()
	}

	return builder
}
