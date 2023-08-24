package handler

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/jackc/pgconn"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/pseudo"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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
	RetryFailedAfter      time.Duration
	HandleActiveInstances time.Duration
	MaxFailureCount       uint8

	TriggerWithoutEvents Reduce
}

type Handler struct {
	client     *database.DB
	projection Projection

	es         EventStore
	bulkLimit  uint16
	eventTypes map[eventstore.AggregateType][]eventstore.EventType

	maxFailureCount       uint8
	retryFailedAfter      time.Duration
	requeueEvery          time.Duration
	handleActiveInstances time.Duration
	now                   nowFunc

	triggeredInstancesSync sync.Map

	triggerWithoutEvents Reduce
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
		projection:             projection,
		client:                 config.Client,
		es:                     config.Eventstore,
		bulkLimit:              config.BulkLimit,
		eventTypes:             aggregates,
		requeueEvery:           config.RequeueEvery,
		handleActiveInstances:  config.HandleActiveInstances,
		now:                    time.Now,
		maxFailureCount:        config.MaxFailureCount,
		retryFailedAfter:       config.RetryFailedAfter,
		triggeredInstancesSync: sync.Map{},
		triggerWithoutEvents:   config.TriggerWithoutEvents,
	}

	return handler
}

func (h *Handler) Start(ctx context.Context) {
	go h.schedule(ctx)
	if h.triggerWithoutEvents != nil {
		return
	}
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
			scheduledCtx := call.WithTimestamp(ctx)
			for _, instance := range instances {
				instanceCtx := authz.WithInstanceID(scheduledCtx, instance)
				_, err = h.Trigger(instanceCtx)
				instanceFailed = instanceFailed || err != nil
				h.log().WithField("instance", instance).OnError(err).Info("scheduled trigger failed")

				for ; err != nil; _, err = h.Trigger(instanceCtx) {
					instanceFailed = instanceFailed || err != nil
					h.log().WithField("instance", instance).OnError(err).Info("scheduled trigger failed")
					if err == nil {
						break
					}
					time.Sleep(h.retryFailedAfter)
				}
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
			return
		case event := <-queue:
			events := checkAdditionalEvents(queue, event)
			solvedInstances := make([]string, 0, len(events))
			queueCtx := call.WithTimestamp(ctx)
			for _, e := range events {
				if instanceSolved(solvedInstances, e.Aggregate().InstanceID) {
					continue
				}
				queueCtx = authz.WithInstanceID(queueCtx, e.Aggregate().InstanceID)
				_, err := h.Trigger(queueCtx)
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
		// AggregateTypes(instance.AggregateType).
		// EventTypes(instance.InstanceAddedEventType).
		ExcludedInstanceID("")
	if didInitialize {
		query = query.
			CreationDateAfter(h.now().Add(-1 * h.handleActiveInstances))
	}
	return h.es.InstanceIDs(ctx, query.Builder())
}

func (h *Handler) Trigger(ctx context.Context) (_ context.Context, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	cancel := h.lockInstance(ctx)
	if cancel == nil {
		return call.ResetTimestamp(ctx), nil
	}
	defer cancel()

	for i := 0; ; i++ {
		additionalIteration, err := h.processEvents(ctx)
		h.log().WithField("iteration", i).Debug("trigger iteration")
		if !additionalIteration || err != nil {
			return call.ResetTimestamp(ctx), err
		}
	}
}

func (h *Handler) lockInstance(ctx context.Context) func() {
	instanceID := authz.GetInstance(ctx).InstanceID()

	instanceMu, _ := h.triggeredInstancesSync.LoadOrStore(instanceID, new(sync.Mutex))
	if !instanceMu.(*sync.Mutex).TryLock() {
		instanceMu.(*sync.Mutex).Lock()
		defer instanceMu.(*sync.Mutex).Unlock()
		return nil
	}
	return func() {
		instanceMu.(*sync.Mutex).Unlock()
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

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var processErr error

	err = crdb.ExecuteTx(ctx, h.client.DB, nil, func(tx *sql.Tx) error {
		currentState, err := h.currentState(ctx, tx)
		if err != nil {
			return err
		}

		var statements []*Statement
		statements, additionalIteration, err = h.generateStatements(ctx, tx, currentState)
		if err != nil || len(statements) == 0 {
			return err
		}

		lastProcessedIndex, err := h.execute(ctx, tx, currentState, statements)
		if lastProcessedIndex < 0 {
			processErr = err
			return nil
		}

		currentState.position = statements[lastProcessedIndex].Position
		currentState.eventTimestamp = statements[lastProcessedIndex].CreationDate

		return h.setState(ctx, tx, currentState)
	})

	if processErr != nil {
		return false, processErr
	}

	return additionalIteration, err
}

func (h *Handler) generateStatements(ctx context.Context, tx *sql.Tx, currentState *state) (_ []*Statement, additionalIteration bool, err error) {
	if h.triggerWithoutEvents != nil {
		stmt, err := h.triggerWithoutEvents(pseudo.NewScheduledEvent(ctx, time.Now(), currentState.instanceID))
		if err != nil {
			return nil, false, err
		}
		return []*Statement{stmt}, false, nil
	}

	events, err := h.es.Filter(ctx, h.eventQuery(currentState))
	if err != nil {
		h.log().WithError(err).Debug("filter eventstore failed")
		return nil, false, err
	}
	eventAmount := len(events)
	events = skipPreviouslyReduced(events, currentState)

	if len(events) == 0 {
		h.updateLastUpdated(ctx, tx, currentState)
		return nil, false, nil
	}

	statements, err := h.eventsToStatements(tx, events, currentState)
	if len(statements) == 0 {
		return nil, false, err
	}

	additionalIteration = eventAmount == int(h.bulkLimit)
	if len(statements) < len(events) {
		// retry imediatly if statements failed
		additionalIteration = true
	}

	return statements, additionalIteration, nil
}

func skipPreviouslyReduced(events []eventstore.Event, currentState *state) []eventstore.Event {
	for i, event := range events {
		if event.Position() == currentState.position {
			return events[i+1:]
		}
	}
	return events
}

func (h *Handler) execute(ctx context.Context, tx *sql.Tx, currentState *state, statements []*Statement) (lastProcessedIndex int, err error) {
	lastProcessedIndex = -1

	for i, statement := range statements {
		if statement.Execute == nil {
			lastProcessedIndex = i
			continue
		}
		_, err := tx.Exec("SAVEPOINT exec")
		if err != nil {
			h.log().WithError(err).Debug("create savepoint failed")
			return lastProcessedIndex, err
		}
		if err = statement.Execute(tx, h.projection.Name()); err != nil {
			h.log().WithError(err).Error("statement execution failed")

			_, savepointErr := tx.Exec("ROLLBACK TO SAVEPOINT exec")
			if savepointErr != nil {
				h.log().WithError(savepointErr).Debug("rollback savepoint failed")
				return lastProcessedIndex, savepointErr
			}

			if h.handleFailedStmt(tx, currentState, failureFromStatement(statement, err)) {
				lastProcessedIndex = i
				continue
			}

			return lastProcessedIndex, err
		}
		if _, err = tx.Exec("RELEASE SAVEPOINT exec"); err != nil {
			return lastProcessedIndex, err
		}
		lastProcessedIndex = i
	}
	return lastProcessedIndex, nil
}

func (h *Handler) eventQuery(currentState *state) *eventstore.SearchQueryBuilder {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		Limit(uint64(h.bulkLimit)).
		AllowTimeTravel().
		OrderAsc().
		InstanceID(currentState.instanceID)

	if currentState.position > 0 {
		builder = builder.PositionAfter(currentState.position)
	}

	for aggregateType, eventTypes := range h.eventTypes {
		query := builder.
			AddQuery().
			AggregateTypes(aggregateType).
			EventTypes(eventTypes...)

		builder = query.Builder()
	}

	return builder
}
