package eventstore

import (
	"context"
	"sync"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type Event interface {
	//CheckPrevious ensures the event order if true
	// if false the previous sequence is not checked on push
	CheckPrevious() bool

	EditorService() string
	EditorUser() string
	Type() EventType
	Data() interface{}
	PreviousSequence() uint64
}

type eventAppender interface {
	//AppendEvents appends the passed events to an internal list of events
	AppendEvents(...Event) error
}

type reducer interface {
	//Reduce handles the events of the internal events list
	// it only appends the newly added events
	Reduce() error
}
type aggregater interface {
	eventAppender
	reducer
	//ID returns the aggreagte id
	ID() string
	//Type returns the aggregate type
	Type() AggregateType
	//Events returns the events which will be pushed
	Events() []Event
	//ResourceOwner returns the organisation id which manages this aggregate
	ResourceOwner() string
	//Version represents the semantic version of the aggregate
	Version() Version
}
type readModeler interface {
	eventAppender
	reducer
}

type Eventstore struct {
	repo             repository.Repository
	interceptorMutex sync.Mutex
	eventMapper      map[EventType]eventTypeInterceptors
}

type eventTypeInterceptors struct {
	pushMapper   func(Event) (*repository.Event, error)
	filterMapper func(*repository.Event) (Event, error)
}

//Health checks if the eventstore can properly work
// It checks if the repository can serve load
func (es *Eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

//PushAggregates maps the events of all aggregates to an eventstore event
// based on the pushMapper
func (es *Eventstore) PushAggregates(ctx context.Context, aggregates ...aggregater) ([]Event, error) {
	events, err := es.aggregatesToEvents(aggregates)
	if err != nil {
		return nil, err
	}

	err = es.repo.Push(ctx, events...)
	if err != nil {
		return nil, err
	}

	return es.mapEvents(events)
}

func (es *Eventstore) aggregatesToEvents(aggregates []aggregater) ([]*repository.Event, error) {
	events := make([]*repository.Event, 0, len(aggregates))
	for _, aggregate := range aggregates {
		var previousEvent *repository.Event
		for _, event := range aggregate.Events() {
			//TODO: map event.Data() into json
			var data []byte
			events = append(events, &repository.Event{
				AggregateID:      aggregate.ID(),
				AggregateType:    repository.AggregateType(aggregate.Type()),
				ResourceOwner:    aggregate.ResourceOwner(),
				EditorService:    event.EditorService(),
				EditorUser:       event.EditorUser(),
				Type:             repository.EventType(event.Type()),
				Version:          repository.Version(aggregate.Version()),
				PreviousEvent:    previousEvent,
				Data:             data,
				PreviousSequence: event.PreviousSequence(),
			})
			previousEvent = events[len(events)-1]
		}
	}
	return events, nil
}

//FilterEvents filters the stored events based on the searchQuery
// and maps the events to the defined event structs
func (es *Eventstore) FilterEvents(ctx context.Context, queryFactory *SearchQueryFactory) ([]Event, error) {
	query, err := queryFactory.Build()
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
		interceptors, ok := es.eventMapper[EventType(event.Type)]
		if !ok || interceptors.filterMapper == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "V2-usujB", "event mapper not defined")
		}
		mappedEvents[i], err = interceptors.filterMapper(event)
		if err != nil {
			return nil, err
		}
	}

	return mappedEvents, nil
}

//FilterToAggregate filters the events based on the searchQuery, appends all events to the aggregate and reduces the aggregate
func (es *Eventstore) FilterToAggregate(ctx context.Context, searchQuery *SearchQueryFactory, aggregate aggregater) (err error) {
	events, err := es.FilterEvents(ctx, searchQuery)
	if err != nil {
		return err
	}
	if err = aggregate.AppendEvents(events...); err != nil {
		return err
	}

	return aggregate.Reduce()
}

//FilterToReadModel filters the events based on the searchQuery, appends all events to the readModel and reduces the readModel
func (es *Eventstore) FilterToReadModel(ctx context.Context, searchQuery *SearchQueryFactory, readModel readModeler) (err error) {
	events, err := es.FilterEvents(ctx, searchQuery)
	if err != nil {
		return err
	}
	if err = readModel.AppendEvents(events...); err != nil {
		return err
	}

	return readModel.Reduce()
}

func (es *Eventstore) LatestSequence(ctx context.Context, searchQuery *SearchQueryFactory) (uint64, error) {
	return 0, nil
}

//RegisterPushEventMapper registers a function for mapping an eventstore event to an event
func (es *Eventstore) RegisterFilterEventMapper(eventType EventType, mapper func(*repository.Event) (Event, error)) error {
	if eventType == "" || mapper == nil {
		return errors.ThrowInvalidArgument(nil, "V2-IPpUR", "eventType and mapper must be filled")
	}

	es.interceptorMutex.Lock()
	defer es.interceptorMutex.Unlock()

	interceptor := es.eventMapper[eventType]
	interceptor.filterMapper = mapper
	es.eventMapper[eventType] = interceptor

	return nil
}

//RegisterPushEventMapper registers a function for mapping an event to an eventstore event
func (es *Eventstore) RegisterPushEventMapper(eventType EventType, mapper func(Event) (*repository.Event, error)) error {
	if eventType == "" || mapper == nil {
		return errors.ThrowInvalidArgument(nil, "V2-Kexpp", "eventType and mapper must be filled")
	}

	es.interceptorMutex.Lock()
	defer es.interceptorMutex.Unlock()

	interceptor := es.eventMapper[eventType]
	interceptor.pushMapper = mapper
	es.eventMapper[eventType] = interceptor

	return nil
}
