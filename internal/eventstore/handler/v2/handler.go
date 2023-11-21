package handler

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/jackc/pgconn"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/pseudo"
)

type EventStore interface {
	InstanceIDs(ctx context.Context, maxAge time.Duration, forceLoad bool, query *eventstore.SearchQueryBuilder) ([]string, error)
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
	TransactionDuration   time.Duration
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
	txDuration            time.Duration
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
		eventTypes := make([]eventstore.EventType, len(reducer.EventReducers))
		for i, eventReducer := range reducer.EventReducers {
			eventTypes[i] = eventReducer.Event
		}
		if _, ok := aggregates[reducer.Aggregate]; ok {
			aggregates[reducer.Aggregate] = append(aggregates[reducer.Aggregate], eventTypes...)
			continue
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
		txDuration:             config.TransactionDuration,
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

				// simple implementation of do while
				_, err = h.Trigger(instanceCtx)
				instanceFailed = instanceFailed || err != nil
				h.log().WithField("instance", instance).OnError(err).Info("scheduled trigger failed")
				// retry if trigger failed
				for ; err != nil; _, err = h.Trigger(instanceCtx) {
					time.Sleep(h.retryFailedAfter)
					instanceFailed = instanceFailed || err != nil
					h.log().WithField("instance", instance).OnError(err).Info("scheduled trigger failed")
					if err == nil {
						break
					}
				}
			}

			if !didInitialize && !instanceFailed {
				err = h.setSucceededOnce(ctx)
				h.log().OnError(err).Debug("unable to set succeeded once")
				didInitialize = err == nil
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
		AwaitOpenTransactions().
		AllowTimeTravel().
		ExcludedInstanceID("")
	if didInitialize {
		query = query.
			CreationDateAfter(h.now().Add(-1 * h.handleActiveInstances))
	}
	return h.es.InstanceIDs(ctx, h.requeueEvery, !didInitialize, query)
}

type triggerConfig struct {
	awaitRunning bool
}

type triggerOpt func(conf *triggerConfig)

func WithAwaitRunning() triggerOpt {
	return func(conf *triggerConfig) {
		conf.awaitRunning = true
	}
}

func (h *Handler) Trigger(ctx context.Context, opts ...triggerOpt) (_ context.Context, err error) {
	config := new(triggerConfig)
	for _, opt := range opts {
		opt(config)
	}

	cancel := h.lockInstance(ctx, config)
	if cancel == nil {
		return call.ResetTimestamp(ctx), nil
	}
	defer cancel()

	for i := 0; ; i++ {
		additionalIteration, err := h.processEvents(ctx, config)
		h.log().OnError(err).Warn("process events failed")
		h.log().WithField("iteration", i).Debug("trigger iteration")
		if !additionalIteration || err != nil {
			return call.ResetTimestamp(ctx), err
		}
	}
}

// lockInstances tries to lock the instance.
// If the instance is already locked from another process no cancel function is returned
// the instance can be skipped then
// If the instance is locked, an unlock deferable function is returned
func (h *Handler) lockInstance(ctx context.Context, config *triggerConfig) func() {
	instanceID := authz.GetInstance(ctx).InstanceID()

	// Check that the instance has a mutex to lock
	instanceMu, _ := h.triggeredInstancesSync.LoadOrStore(instanceID, new(sync.Mutex))
	unlock := func() {
		instanceMu.(*sync.Mutex).Unlock()
	}
	if !instanceMu.(*sync.Mutex).TryLock() {
		instanceMu.(*sync.Mutex).Lock()
		if config.awaitRunning {
			return unlock
		}
		defer unlock()
		return nil
	}
	return unlock
}

func (h *Handler) processEvents(ctx context.Context, config *triggerConfig) (additionalIteration bool, err error) {
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

	if h.txDuration > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, h.txDuration)
		defer cancel()
	}

	tx, err := h.client.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			h.log().OnError(rollbackErr).Debug("unable to rollback tx")
			return
		}
		err = tx.Commit()
	}()

	currentState, err := h.currentState(ctx, tx, config)
	if err != nil {
		if errors.Is(err, errJustUpdated) {
			return false, nil
		}
		return additionalIteration, err
	}

	var statements []*Statement
	statements, additionalIteration, err = h.generateStatements(ctx, tx, currentState)
	if err != nil || len(statements) == 0 {
		return additionalIteration, err
	}

	lastProcessedIndex, err := h.executeStatements(ctx, tx, currentState, statements)
	if lastProcessedIndex < 0 {
		return false, err
	}

	currentState.position = statements[lastProcessedIndex].Position
	currentState.aggregateID = statements[lastProcessedIndex].AggregateID
	currentState.aggregateType = statements[lastProcessedIndex].AggregateType
	currentState.sequence = statements[lastProcessedIndex].Sequence
	currentState.eventTimestamp = statements[lastProcessedIndex].CreationDate
	err = h.setState(tx, currentState)

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
		if event.Position() == currentState.position &&
			event.Aggregate().ID == currentState.aggregateID &&
			event.Aggregate().Type == currentState.aggregateType &&
			event.Sequence() == currentState.sequence {
			return events[i+1:]
		}
	}
	return events
}

func (h *Handler) executeStatements(ctx context.Context, tx *sql.Tx, currentState *state, statements []*Statement) (lastProcessedIndex int, err error) {
	lastProcessedIndex = -1

	for i, statement := range statements {
		select {
		case <-ctx.Done():
			break
		default:
			err := h.executeStatement(ctx, tx, currentState, statement)
			if err != nil {
				return lastProcessedIndex, err
			}
			lastProcessedIndex = i
		}
	}
	return lastProcessedIndex, nil
}

func (h *Handler) executeStatement(ctx context.Context, tx *sql.Tx, currentState *state, statement *Statement) (err error) {
	if statement.Execute == nil {
		return nil
	}

	_, err = tx.Exec("SAVEPOINT exec")
	if err != nil {
		h.log().WithError(err).Debug("create savepoint failed")
		return err
	}
	var shouldContinue bool
	defer func() {
		_, err = tx.Exec("RELEASE SAVEPOINT exec")
	}()

	if err = statement.Execute(tx, h.projection.Name()); err != nil {
		h.log().WithError(err).Error("statement execution failed")

		shouldContinue = h.handleFailedStmt(tx, currentState, failureFromStatement(statement, err))
		if shouldContinue {
			return nil
		}

		return err
	}

	return nil
}

func (h *Handler) eventQuery(currentState *state) *eventstore.SearchQueryBuilder {
	builder := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		Limit(uint64(h.bulkLimit)).
		AllowTimeTravel().
		OrderAsc().
		InstanceID(currentState.instanceID)

	if currentState.position > 0 {
		builder = builder.PositionAfter(math.Float64frombits(math.Float64bits(currentState.position) - 10))
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

// ProjectionName returns the name of the unlying projection.
func (h *Handler) ProjectionName() string {
	return h.projection.Name()
}
