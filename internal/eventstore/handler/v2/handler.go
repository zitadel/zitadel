package handler

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/migration"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/pseudo"
)

type EventStore interface {
	InstanceIDs(ctx context.Context, query *eventstore.SearchQueryBuilder) ([]string, error)
	FilterToQueryReducer(ctx context.Context, reducer eventstore.QueryReducer) error
	Filter(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error)
	Push(ctx context.Context, cmds ...eventstore.Command) ([]eventstore.Event, error)
	FillFields(ctx context.Context, events ...eventstore.FillFieldsEvent) error
}

type Config struct {
	Client     *database.DB
	Eventstore EventStore

	BulkLimit           uint16
	RequeueEvery        time.Duration
	RetryFailedAfter    time.Duration
	TransactionDuration time.Duration
	MaxFailureCount     uint8

	TriggerWithoutEvents Reduce

	ActiveInstancer interface {
		ActiveInstances() []string
	}
}

type Handler struct {
	client     *database.DB
	projection Projection

	es         EventStore
	bulkLimit  uint16
	eventTypes map[eventstore.AggregateType][]eventstore.EventType

	maxFailureCount  uint8
	retryFailedAfter time.Duration
	requeueEvery     time.Duration
	txDuration       time.Duration
	now              nowFunc

	triggeredInstancesSync sync.Map

	triggerWithoutEvents Reduce
	cacheInvalidations   []func(ctx context.Context, aggregates []*eventstore.Aggregate)

	queryInstances func() ([]string, error)
}

var _ migration.Migration = (*Handler)(nil)

// Execute implements migration.Migration.
func (h *Handler) Execute(ctx context.Context, startedEvent eventstore.Event) error {
	start := time.Now()
	logging.WithFields("projection", h.ProjectionName()).Info("projection starts prefilling")
	logTicker := time.NewTicker(30 * time.Second)
	go func() {
		for range logTicker.C {
			logging.WithFields("projection", h.ProjectionName()).Info("projection is prefilling")
		}
	}()

	instanceIDs, err := h.existingInstances(ctx)
	if err != nil {
		return err
	}

	// default amount of workers is 10
	workerCount := 10

	if h.client.DB.Stats().MaxOpenConnections > 0 {
		workerCount = h.client.DB.Stats().MaxOpenConnections / 4
	}
	// ensure that at least one worker is active
	if workerCount == 0 {
		workerCount = 1
	}
	// spawn less workers if not all workers needed
	if workerCount > len(instanceIDs) {
		workerCount = len(instanceIDs)
	}

	instances := make(chan string, workerCount)
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go h.executeInstances(ctx, instances, startedEvent, &wg)
	}

	for _, instance := range instanceIDs {
		instances <- instance
	}

	close(instances)
	wg.Wait()

	logTicker.Stop()
	logging.WithFields("projection", h.ProjectionName(), "took", time.Since(start)).Info("projections ended prefilling")
	return nil
}

func (h *Handler) executeInstances(ctx context.Context, instances <-chan string, startedEvent eventstore.Event, wg *sync.WaitGroup) {
	for instance := range instances {
		h.triggerInstances(ctx, []string{instance}, WithMaxPosition(startedEvent.Position()))
	}
	wg.Done()
}

// String implements migration.Migration.
func (h *Handler) String() string {
	return h.ProjectionName()
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
		now:                    time.Now,
		maxFailureCount:        config.MaxFailureCount,
		retryFailedAfter:       config.RetryFailedAfter,
		triggeredInstancesSync: sync.Map{},
		triggerWithoutEvents:   config.TriggerWithoutEvents,
		txDuration:             config.TransactionDuration,
		queryInstances: func() ([]string, error) {
			if config.ActiveInstancer != nil {
				return config.ActiveInstancer.ActiveInstances(), nil
			}
			return nil, nil
		},
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

type checkInit struct {
	didInit        bool
	projectionName string
}

// AppendEvents implements eventstore.QueryReducer.
func (ci *checkInit) AppendEvents(...eventstore.Event) {
	ci.didInit = true
}

// Query implements eventstore.QueryReducer.
func (ci *checkInit) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		Limit(1).
		InstanceID("").
		AddQuery().
		AggregateTypes(migration.SystemAggregate).
		AggregateIDs(migration.SystemAggregateID).
		EventTypes(migration.DoneType).
		EventData(map[string]interface{}{
			"name": ci.projectionName,
		}).
		Builder()
}

// Reduce implements eventstore.QueryReducer.
func (*checkInit) Reduce() error {
	return nil
}

var _ eventstore.QueryReducer = (*checkInit)(nil)

func (h *Handler) didInitialize(ctx context.Context) bool {
	initiated := checkInit{
		projectionName: h.ProjectionName(),
	}
	err := h.es.FilterToQueryReducer(ctx, &initiated)
	if err != nil {
		return false
	}
	return initiated.didInit
}

func (h *Handler) schedule(ctx context.Context) {
	//  start the projection and its configured `RequeueEvery`
	reset := randomizeStart(0, h.requeueEvery.Seconds())
	if !h.didInitialize(ctx) {
		reset = randomizeStart(0, 0.5)
	}
	t := time.NewTimer(reset)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			instances, err := h.queryInstances()
			h.log().OnError(err).Debug("unable to query instances")

			h.triggerInstances(call.WithTimestamp(ctx), instances)
			t.Reset(h.requeueEvery)
		}
	}
}

func (h *Handler) triggerInstances(ctx context.Context, instances []string, triggerOpts ...TriggerOpt) {
	for _, instance := range instances {
		instanceCtx := authz.WithInstanceID(ctx, instance)

		// simple implementation of do while
		_, err := h.Trigger(instanceCtx, triggerOpts...)
		h.log().WithField("instance", instance).OnError(err).Debug("trigger failed")
		time.Sleep(h.retryFailedAfter)
		// retry if trigger failed
		for ; err != nil; _, err = h.Trigger(instanceCtx, triggerOpts...) {
			time.Sleep(h.retryFailedAfter)
			h.log().WithField("instance", instance).OnError(err).Debug("trigger failed")
		}
	}
}

func randomizeStart(min, maxSeconds float64) time.Duration {
	d := min + rand.Float64()*(maxSeconds-min)
	return time.Duration(d*1000) * time.Millisecond
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

type existingInstances []string

// AppendEvents implements eventstore.QueryReducer.
func (ai *existingInstances) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch event.Type() {
		case instance.InstanceAddedEventType:
			*ai = append(*ai, event.Aggregate().InstanceID)
		case instance.InstanceRemovedEventType:
			*ai = slices.DeleteFunc(*ai, func(s string) bool {
				return s == event.Aggregate().InstanceID
			})
		}
	}
}

// Query implements eventstore.QueryReducer.
func (*existingInstances) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.InstanceAddedEventType,
			instance.InstanceRemovedEventType,
		).
		Builder()
}

// Reduce implements eventstore.QueryReducer.
// reduce is not used as events are reduced during AppendEvents
func (*existingInstances) Reduce() error {
	return nil
}

var _ eventstore.QueryReducer = (*existingInstances)(nil)

func (h *Handler) existingInstances(ctx context.Context) ([]string, error) {
	ai := existingInstances{}
	if err := h.es.FilterToQueryReducer(ctx, &ai); err != nil {
		return nil, err
	}

	return ai, nil
}

type triggerConfig struct {
	awaitRunning bool
	maxPosition  float64
}

type TriggerOpt func(conf *triggerConfig)

func WithAwaitRunning() TriggerOpt {
	return func(conf *triggerConfig) {
		conf.awaitRunning = true
	}
}

func WithMaxPosition(position float64) TriggerOpt {
	return func(conf *triggerConfig) {
		conf.maxPosition = position
	}
}

func (h *Handler) Trigger(ctx context.Context, opts ...TriggerOpt) (_ context.Context, err error) {
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
		h.log().OnError(err).Info("process events failed")
		h.log().WithField("iteration", i).Debug("trigger iteration")
		if !additionalIteration || err != nil {
			return call.ResetTimestamp(ctx), err
		}
	}
}

// RegisterCacheInvalidation registers a function to be called when a cache needs to be invalidated.
// In order to avoid race conditions, this method must be called before [Handler.Start] is called.
func (h *Handler) RegisterCacheInvalidation(invalidate func(ctx context.Context, aggregates []*eventstore.Aggregate)) {
	h.cacheInvalidations = append(h.cacheInvalidations, invalidate)
}

// lockInstance tries to lock the instance.
// If the instance is already locked from another process no cancel function is returned
// the instance can be skipped then
// If the instance is locked, an unlock deferrable function is returned
func (h *Handler) lockInstance(ctx context.Context, config *triggerConfig) func() {
	instanceID := authz.GetInstance(ctx).InstanceID()

	// Check that the instance has a lock
	instanceLock, _ := h.triggeredInstancesSync.LoadOrStore(instanceID, make(chan bool, 1))

	// in case we don't want to wait for a running trigger / lock (e.g. spooler),
	// we can directly return if we cannot lock
	if !config.awaitRunning {
		select {
		case instanceLock.(chan bool) <- true:
			return func() {
				<-instanceLock.(chan bool)
			}
		default:
			return nil
		}
	}

	// in case we want to wait for a running trigger / lock (e.g. query),
	// we try to lock as long as the context is not cancelled
	select {
	case instanceLock.(chan bool) <- true:
		return func() {
			<-instanceLock.(chan bool)
		}
	case <-ctx.Done():
		return nil
	}
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
	}()

	currentState, err := h.currentState(ctx, tx, config)
	if err != nil {
		if errors.Is(err, errJustUpdated) {
			return false, nil
		}
		return additionalIteration, err
	}
	// stop execution if currentState.eventTimestamp >= config.maxCreatedAt
	if config.maxPosition != 0 && currentState.position >= config.maxPosition {
		return false, nil
	}

	var statements []*Statement
	statements, additionalIteration, err = h.generateStatements(ctx, tx, currentState)
	if err != nil {
		return additionalIteration, err
	}

	defer func() {
		commitErr := tx.Commit()
		if err == nil {
			err = commitErr
		}
		if err == nil && currentState.aggregateID != "" && len(statements) > 0 {
			h.invalidateCaches(ctx, aggregatesFromStatements(statements))
		}
	}()

	if len(statements) == 0 {
		err = h.setState(tx, currentState)
		return additionalIteration, err
	}

	lastProcessedIndex, err := h.executeStatements(ctx, tx, statements)
	h.log().OnError(err).WithField("lastProcessedIndex", lastProcessedIndex).Debug("execution of statements failed")
	if lastProcessedIndex < 0 {
		return false, err
	}

	currentState.position = statements[lastProcessedIndex].Position
	currentState.offset = statements[lastProcessedIndex].offset
	currentState.aggregateID = statements[lastProcessedIndex].Aggregate.ID
	currentState.aggregateType = statements[lastProcessedIndex].Aggregate.Type
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

	events, err := h.es.Filter(ctx, h.eventQuery(currentState).SetTx(tx))
	if err != nil {
		h.log().WithError(err).Debug("filter eventstore failed")
		return nil, false, err
	}
	eventAmount := len(events)

	statements, err := h.eventsToStatements(tx, events, currentState)
	if err != nil || len(statements) == 0 {
		return nil, false, err
	}

	idx := skipPreviouslyReducedStatements(statements, currentState)
	if idx+1 == len(statements) {
		currentState.position = statements[len(statements)-1].Position
		currentState.offset = statements[len(statements)-1].offset
		currentState.aggregateID = statements[len(statements)-1].Aggregate.ID
		currentState.aggregateType = statements[len(statements)-1].Aggregate.Type
		currentState.sequence = statements[len(statements)-1].Sequence
		currentState.eventTimestamp = statements[len(statements)-1].CreationDate

		return nil, false, nil
	}
	statements = statements[idx+1:]

	additionalIteration = eventAmount == int(h.bulkLimit)
	if len(statements) < len(events) {
		// retry immediately if statements failed
		additionalIteration = true
	}

	return statements, additionalIteration, nil
}

func skipPreviouslyReducedStatements(statements []*Statement, currentState *state) int {
	for i, statement := range statements {
		if statement.Position == currentState.position &&
			statement.Aggregate.ID == currentState.aggregateID &&
			statement.Aggregate.Type == currentState.aggregateType &&
			statement.Sequence == currentState.sequence {
			return i
		}
	}
	return -1
}

func (h *Handler) executeStatements(ctx context.Context, tx *sql.Tx, statements []*Statement) (lastProcessedIndex int, err error) {
	lastProcessedIndex = -1

	for i, statement := range statements {
		select {
		case <-ctx.Done():
			break
		default:
			err := h.executeStatement(ctx, tx, statement)
			if err != nil {
				return lastProcessedIndex, err
			}
			lastProcessedIndex = i
		}
	}
	return lastProcessedIndex, nil
}

func (h *Handler) executeStatement(ctx context.Context, tx *sql.Tx, statement *Statement) (err error) {
	if statement.Execute == nil {
		return nil
	}

	_, err = tx.ExecContext(ctx, "SAVEPOINT exec_stmt")
	if err != nil {
		h.log().WithError(err).Debug("create savepoint failed")
		return err
	}

	if err = statement.Execute(tx, h.projection.Name()); err != nil {
		h.log().WithError(err).Error("statement execution failed")

		_, rollbackErr := tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT exec_stmt")
		h.log().OnError(rollbackErr).Error("rollback to savepoint failed")

		shouldContinue := h.handleFailedStmt(tx, failureFromStatement(statement, err))
		if shouldContinue {
			return nil
		}

		return &executionError{parent: err}
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
		// decrease position by 10 because builder.PositionAfter filters for position > and we need position >=
		builder = builder.PositionAfter(math.Float64frombits(math.Float64bits(currentState.position) - 10))
		if currentState.offset > 0 {
			builder = builder.Offset(currentState.offset)
		}
	}

	for aggregateType, eventTypes := range h.eventTypes {
		builder = builder.
			AddQuery().
			AggregateTypes(aggregateType).
			EventTypes(eventTypes...).
			Builder()
	}

	return builder
}

// ProjectionName returns the name of the underlying projection.
func (h *Handler) ProjectionName() string {
	return h.projection.Name()
}

func (h *Handler) invalidateCaches(ctx context.Context, aggregates []*eventstore.Aggregate) {
	if len(h.cacheInvalidations) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(h.cacheInvalidations))

	for _, invalidate := range h.cacheInvalidations {
		go func(invalidate func(context.Context, []*eventstore.Aggregate)) {
			defer wg.Done()
			invalidate(ctx, aggregates)
		}(invalidate)
	}
	wg.Wait()
}

// aggregatesFromStatements returns the unique aggregates from statements.
// Duplicate aggregates are omitted.
func aggregatesFromStatements(statements []*Statement) []*eventstore.Aggregate {
	aggregates := make([]*eventstore.Aggregate, 0, len(statements))
	for _, statement := range statements {
		if !slices.ContainsFunc(aggregates, func(aggregate *eventstore.Aggregate) bool {
			return *statement.Aggregate == *aggregate
		}) {
			aggregates = append(aggregates, statement.Aggregate)
		}
	}
	return aggregates
}
