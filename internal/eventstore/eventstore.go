package eventstore

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	PushTimeout time.Duration
	maxRetries  int

	pusher   Pusher
	querier  Querier
	searcher Searcher
}

var (
	eventInterceptors map[EventType]eventTypeInterceptors
	eventTypes        []string
	aggregateTypes    []string
	eventTypeMapping  = map[EventType]AggregateType{}
)

// RegisterFilterEventMapper registers a function for mapping an eventstore event to an event
func RegisterFilterEventMapper(aggregateType AggregateType, eventType EventType, mapper func(Event) (Event, error)) {
	if mapper == nil || eventType == "" {
		return
	}

	appendEventType(eventType)
	appendAggregateType(aggregateType)

	if eventInterceptors == nil {
		eventInterceptors = make(map[EventType]eventTypeInterceptors)
	}

	interceptor := eventInterceptors[eventType]
	interceptor.eventMapper = mapper
	eventInterceptors[eventType] = interceptor
	eventTypeMapping[eventType] = aggregateType
}

type eventTypeInterceptors struct {
	eventMapper func(Event) (Event, error)
}

func NewEventstore(config *Config) *Eventstore {
	return &Eventstore{
		PushTimeout: config.PushTimeout,
		maxRetries:  int(config.MaxRetries),

		pusher:   config.Pusher,
		querier:  config.Querier,
		searcher: config.Searcher,
	}
}

// Health checks if the eventstore can properly work
// It checks if the repository can serve load
func (es *Eventstore) Health(ctx context.Context) error {
	if err := es.pusher.Health(ctx); err != nil {
		return err
	}
	return es.querier.Health(ctx)
}

// Push pushes the events in a single transaction
// an event needs at least an aggregate
func (es *Eventstore) Push(ctx context.Context, cmds ...Command) ([]Event, error) {
	return es.PushWithClient(ctx, nil, cmds...)
}

// PushWithClient pushes the events in a single transaction using the provided database client
// an event needs at least an aggregate
func (es *Eventstore) PushWithClient(ctx context.Context, client database.ContextQueryExecuter, cmds ...Command) ([]Event, error) {
	if es.PushTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, es.PushTimeout)
		defer cancel()
	}
	var (
		events []Event
		err    error
	)

	// Retry when there is a collision of the sequence as part of the primary key.
	// "duplicate key value violates unique constraint \"events2_pkey\" (SQLSTATE 23505)"
	// https://github.com/zitadel/zitadel/issues/7202
retry:
	for i := 0; i <= es.maxRetries; i++ {
		events, err = es.pusher.Push(ctx, client, cmds...)
		// if there is a transaction passed the calling function needs to retry
		if _, ok := client.(database.Tx); ok {
			break retry
		}
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			break retry
		}
		if pgErr.ConstraintName == "events2_pkey" && pgErr.SQLState() == "23505" {
			logging.WithError(err).Info("eventstore push retry")
			continue
		}
		if pgErr.SQLState() == "CR000" || pgErr.SQLState() == "40001" {
			logging.WithError(err).Info("eventstore push retry")
			continue
		}
		break retry
	}
	if err != nil {
		return nil, err
	}

	mappedEvents, err := es.mapEvents(events)
	if err != nil {
		return mappedEvents, err
	}
	es.notify(mappedEvents)
	return mappedEvents, nil
}

func AggregateTypeFromEventType(typ EventType) AggregateType {
	return eventTypeMapping[typ]
}

func (es *Eventstore) EventTypes() []string {
	return eventTypes
}

func (es *Eventstore) AggregateTypes() []string {
	return aggregateTypes
}

// FillFields implements the [Searcher] interface
func (es *Eventstore) FillFields(ctx context.Context, events ...FillFieldsEvent) error {
	return es.searcher.FillFields(ctx, events...)
}

// Search implements the [Searcher] interface
func (es *Eventstore) Search(ctx context.Context, conditions ...map[FieldType]any) ([]*SearchResult, error) {
	if len(conditions) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "V3-5Xbr1", "no search conditions")
	}

	return es.searcher.Search(ctx, conditions...)
}

// Filter filters the stored events based on the searchQuery
// and maps the events to the defined event structs
//
// Deprecated: Use [FilterToQueryReducer] instead to avoid allocations.
func (es *Eventstore) Filter(ctx context.Context, searchQuery *SearchQueryBuilder) ([]Event, error) {
	events := make([]Event, 0, searchQuery.GetLimit())
	searchQuery.ensureInstanceID(ctx)
	err := es.querier.FilterToReducer(ctx, searchQuery, func(event Event) error {
		event, err := es.mapEvent(event)
		if err != nil {
			return err
		}
		events = append(events, event)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (es *Eventstore) mapEvents(events []Event) (mappedEvents []Event, err error) {
	mappedEvents = make([]Event, len(events))
	for i, event := range events {
		mappedEvents[i], err = es.mapEventLocked(event)
		if err != nil {
			return nil, err
		}
	}
	return mappedEvents, nil
}

func (es *Eventstore) mapEvent(event Event) (Event, error) {
	return es.mapEventLocked(event)
}

func (es *Eventstore) mapEventLocked(event Event) (Event, error) {
	interceptors, ok := eventInterceptors[event.Type()]
	if !ok || interceptors.eventMapper == nil {
		return BaseEventFromRepo(event), nil
	}
	return interceptors.eventMapper(event)
}

// TODO: refactor so we can change to the following interface:
/*
type reducer interface {
	// Reduce applies an event on the object.
	Reduce(Event) error
}
*/

type reducer interface {
	//Reduce handles the events of the internal events list
	// it only appends the newly added events
	Reduce() error
	//AppendEvents appends the passed events to an internal list of events
	AppendEvents(...Event)
}

// FilterToReducer filters the events based on the search query, appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r reducer) error {
	searchQuery.ensureInstanceID(ctx)
	return es.querier.FilterToReducer(ctx, searchQuery, func(event Event) error {
		event, err := es.mapEvent(event)
		if err != nil {
			return err
		}
		r.AppendEvents(event)
		return r.Reduce()
	})
}

// LatestSequence filters the latest sequence for the given search query
func (es *Eventstore) LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (float64, error) {
	queryFactory.InstanceID(authz.GetInstance(ctx).InstanceID())

	return es.querier.LatestSequence(ctx, queryFactory)
}

// InstanceIDs returns the distinct instance ids found by the search query
// Warning: this function can have high impact on performance, only use this function during setup
func (es *Eventstore) InstanceIDs(ctx context.Context, queryFactory *SearchQueryBuilder) ([]string, error) {
	return es.querier.InstanceIDs(ctx, queryFactory)
}

func (es *Eventstore) Client() *database.DB {
	return es.querier.Client()
}

type QueryReducer interface {
	reducer
	//Query returns the SearchQueryFactory for the events needed in reducer
	Query() *SearchQueryBuilder
}

// FilterToQueryReducer filters the events based on the search query of the query function,
// appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToQueryReducer(ctx context.Context, r QueryReducer) error {
	return es.FilterToReducer(ctx, r.Query(), r)
}

type Reducer func(event Event) error

type Querier interface {
	// Health checks if the connection to the storage is available
	Health(ctx context.Context) error
	// FilterToReducer calls r for every event returned from the storage
	FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r Reducer) error
	// LatestSequence returns the latest sequence found by the search query
	LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (float64, error)
	// InstanceIDs returns the instance ids found by the search query
	InstanceIDs(ctx context.Context, queryFactory *SearchQueryBuilder) ([]string, error)
	// Client returns the underlying database connection
	Client() *database.DB
}

type Pusher interface {
	// Health checks if the connection to the storage is available
	Health(ctx context.Context) error
	// Push stores the actions
	Push(ctx context.Context, client database.ContextQueryExecuter, commands ...Command) (_ []Event, err error)
	// Client returns the underlying database connection
	Client() *database.DB
}

type FillFieldsEvent interface {
	Event
	Fields() []*FieldOperation
}

type Searcher interface {
	// Search allows to search for specific fields of objects
	// The instance id is taken from the context
	// The list of conditions are combined with AND
	// The search fields are combined with OR
	// At least one must be defined
	Search(ctx context.Context, conditions ...map[FieldType]any) (result []*SearchResult, err error)
	// FillFields is to insert the fields of previously stored events
	FillFields(ctx context.Context, events ...FillFieldsEvent) error
}

func appendEventType(typ EventType) {
	i := sort.SearchStrings(eventTypes, string(typ))
	if i < len(eventTypes) && eventTypes[i] == string(typ) {
		return
	}
	eventTypes = append(eventTypes[:i], append([]string{string(typ)}, eventTypes[i:]...)...)
}

func appendAggregateType(typ AggregateType) {
	i := sort.SearchStrings(aggregateTypes, string(typ))
	if len(aggregateTypes) > i && aggregateTypes[i] == string(typ) {
		return
	}
	aggregateTypes = append(aggregateTypes[:i], append([]string{string(typ)}, aggregateTypes[i:]...)...)
}
