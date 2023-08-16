package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

// Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	repo              repository.Repository
	interceptorMutex  sync.Mutex
	eventInterceptors map[EventType]eventTypeInterceptors
	eventTypes        []string
	aggregateTypes    []string
	PushTimeout       time.Duration
}

type eventTypeInterceptors struct {
	eventMapper func(*repository.Event) (Event, error)
}

func NewEventstore(config *Config) *Eventstore {
	return &Eventstore{
		repo:              config.repo,
		eventInterceptors: map[EventType]eventTypeInterceptors{},
		interceptorMutex:  sync.Mutex{},
		PushTimeout:       config.PushTimeout,
	}
}

// Health checks if the eventstore can properly work
// It checks if the repository can serve load
func (es *Eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

// Push pushes the events in a single transaction
// an event needs at least an aggregate
func (es *Eventstore) Push(ctx context.Context, cmds ...Command) ([]Event, error) {
	events, constraints, err := commandsToRepository(authz.GetInstance(ctx).InstanceID(), cmds)
	if err != nil {
		return nil, err
	}

	if es.PushTimeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, es.PushTimeout)
		defer cancel()
	}

	err = es.repo.Push(ctx, events, constraints...)
	if err != nil {
		return nil, err
	}

	eventReaders, err := es.mapEvents(events)
	if err != nil {
		return nil, err
	}

	go notify(eventReaders)
	return eventReaders, nil
}

func (es *Eventstore) NewInstance(ctx context.Context, instanceID string) error {
	return es.repo.CreateInstance(ctx, instanceID)
}

func (es *Eventstore) EventTypes() []string {
	return es.eventTypes
}

func (es *Eventstore) AggregateTypes() []string {
	return es.aggregateTypes
}

func commandsToRepository(instanceID string, cmds []Command) (events []*repository.Event, constraints []*repository.UniqueConstraint, err error) {
	events = make([]*repository.Event, len(cmds))
	for i, cmd := range cmds {
		data, err := EventData(cmd)
		if err != nil {
			return nil, nil, err
		}
		if cmd.Aggregate().ID == "" {
			return nil, nil, errors.ThrowInvalidArgument(nil, "V2-Afdfe", "aggregate id must not be empty")
		}
		if cmd.Aggregate().Type == "" {
			return nil, nil, errors.ThrowInvalidArgument(nil, "V2-Dfg32", "aggregate type must not be empty")
		}
		if cmd.Type() == "" {
			return nil, nil, errors.ThrowInvalidArgument(nil, "V2-Drg34", "event type must not be empty")
		}
		if cmd.Aggregate().Version == "" {
			return nil, nil, errors.ThrowInvalidArgument(nil, "V2-Dgfg4", "aggregate version must not be empty")
		}
		events[i] = &repository.Event{
			AggregateID:   cmd.Aggregate().ID,
			AggregateType: repository.AggregateType(cmd.Aggregate().Type),
			ResourceOwner: sql.NullString{String: cmd.Aggregate().ResourceOwner, Valid: cmd.Aggregate().ResourceOwner != ""},
			InstanceID:    instanceID,
			EditorService: cmd.EditorService(),
			EditorUser:    cmd.EditorUser(),
			Type:          repository.EventType(cmd.Type()),
			Version:       repository.Version(cmd.Aggregate().Version),
			Data:          data,
		}
		if len(cmd.UniqueConstraints()) > 0 {
			constraints = append(constraints, uniqueConstraintsToRepository(instanceID, cmd.UniqueConstraints())...)
		}
	}

	return events, constraints, nil
}

func uniqueConstraintsToRepository(instanceID string, constraints []*EventUniqueConstraint) (uniqueConstraints []*repository.UniqueConstraint) {
	uniqueConstraints = make([]*repository.UniqueConstraint, len(constraints))
	for i, constraint := range constraints {
		var id string
		if !constraint.IsGlobal {
			id = instanceID
		}
		uniqueConstraints[i] = &repository.UniqueConstraint{
			UniqueType:   constraint.UniqueType,
			UniqueField:  constraint.UniqueField,
			InstanceID:   id,
			Action:       uniqueConstraintActionToRepository(constraint.Action),
			ErrorMessage: constraint.ErrorMessage,
		}
	}
	return uniqueConstraints
}

// Filter filters the stored events based on the searchQuery
// and maps the events to the defined event structs
func (es *Eventstore) Filter(ctx context.Context, queryFactory *SearchQueryBuilder) ([]Event, error) {
	query, err := queryFactory.build(authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	events, err := es.repo.Filter(ctx, query)
	if err != nil {
		return nil, err
	}

	return es.mapEvents(events)
}

func (es *Eventstore) mapEvents(events []*repository.Event) (mappedEvents []Event, err error) {
	mappedEvents = make([]Event, len(events))

	es.interceptorMutex.Lock()
	defer es.interceptorMutex.Unlock()

	for i, event := range events {
		interceptors, ok := es.eventInterceptors[EventType(event.Type)]
		if !ok || interceptors.eventMapper == nil {
			mappedEvents[i] = BaseEventFromRepo(event)
			//TODO: return error if unable to map event
			continue
			// return nil, errors.ThrowPreconditionFailed(nil, "V2-usujB", "event mapper not defined")
		}
		mappedEvents[i], err = interceptors.eventMapper(event)
		if err != nil {
			return nil, err
		}
	}

	return mappedEvents, nil
}

type reducer interface {
	//Reduce handles the events of the internal events list
	// it only appends the newly added events
	Reduce() error
	//AppendEvents appends the passed events to an internal list of events
	AppendEvents(...Event)
}

// FilterToReducer filters the events based on the search query, appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r reducer) error {
	events, err := es.Filter(ctx, searchQuery)
	if err != nil {
		return err
	}

	r.AppendEvents(events...)

	return r.Reduce()
}

// LatestSequence filters the latest sequence for the given search query
func (es *Eventstore) LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (uint64, error) {
	query, err := queryFactory.build(authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return 0, err
	}
	return es.repo.LatestSequence(ctx, query)
}

// InstanceIDs returns the instance ids found by the search query
func (es *Eventstore) InstanceIDs(ctx context.Context, queryFactory *SearchQueryBuilder) ([]string, error) {
	query, err := queryFactory.build(authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return es.repo.InstanceIDs(ctx, query)
}

type QueryReducer interface {
	reducer
	//Query returns the SearchQueryFactory for the events needed in reducer
	Query() *SearchQueryBuilder
}

// FilterToQueryReducer filters the events based on the search query of the query function,
// appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToQueryReducer(ctx context.Context, r QueryReducer) error {
	events, err := es.Filter(ctx, r.Query())
	if err != nil {
		return err
	}
	r.AppendEvents(events...)

	return r.Reduce()
}

// RegisterFilterEventMapper registers a function for mapping an eventstore event to an event
func (es *Eventstore) RegisterFilterEventMapper(aggregateType AggregateType, eventType EventType, mapper func(*repository.Event) (Event, error)) *Eventstore {
	if mapper == nil || eventType == "" {
		return es
	}
	es.interceptorMutex.Lock()
	defer es.interceptorMutex.Unlock()

	es.appendEventType(eventType)
	es.appendAggregateType(aggregateType)

	interceptor := es.eventInterceptors[eventType]
	interceptor.eventMapper = mapper
	es.eventInterceptors[eventType] = interceptor

	return es
}

func (es *Eventstore) appendEventType(typ EventType) {
	i := sort.SearchStrings(es.eventTypes, string(typ))
	if i < len(es.eventTypes) && es.eventTypes[i] == string(typ) {
		return
	}
	es.eventTypes = append(es.eventTypes[:i], append([]string{string(typ)}, es.eventTypes[i:]...)...)
}

func (es *Eventstore) appendAggregateType(typ AggregateType) {
	i := sort.SearchStrings(es.aggregateTypes, string(typ))
	if len(es.aggregateTypes) > i && es.aggregateTypes[i] == string(typ) {
		return
	}
	es.aggregateTypes = append(es.aggregateTypes[:i], append([]string{string(typ)}, es.aggregateTypes[i:]...)...)
}

func EventData(event Command) ([]byte, error) {
	switch data := event.Data().(type) {
	case nil:
		return nil, nil
	case []byte:
		if json.Valid(data) {
			return data, nil
		}
		return nil, errors.ThrowInvalidArgument(nil, "V2-6SbbS", "data bytes are not json")
	}
	dataType := reflect.TypeOf(event.Data())
	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}
	if dataType.Kind() == reflect.Struct {
		dataBytes, err := json.Marshal(event.Data())
		if err != nil {
			return nil, errors.ThrowInvalidArgument(err, "V2-xG87M", "could  not marshal data")
		}
		return dataBytes, nil
	}
	return nil, errors.ThrowInvalidArgument(nil, "V2-91NRm", "wrong type of event data")
}

func uniqueConstraintActionToRepository(action UniqueConstraintAction) repository.UniqueConstraintAction {
	switch action {
	case UniqueConstraintAdd:
		return repository.UniqueConstraintAdd
	case UniqueConstraintRemove:
		return repository.UniqueConstraintRemoved
	case UniqueConstraintInstanceRemove:
		return repository.UniqueConstraintInstanceRemoved
	default:
		return repository.UniqueConstraintAdd
	}
}

type BaseEventSetter[T any] interface {
	Event
	SetBaseEvent(*BaseEvent)
	*T
}

func GenericEventMapper[T any, PT BaseEventSetter[T]](event *repository.Event) (Event, error) {
	e := PT(new(T))
	e.SetBaseEvent(BaseEventFromRepo(event))

	if len(event.Data) == 0 {
		return e, nil
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "V2-Thai6", "unable to unmarshal event")
	}

	return e, nil
}
