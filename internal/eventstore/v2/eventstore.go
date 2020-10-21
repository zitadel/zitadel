package eventstore

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

//Event is the representation of a state change
type Event interface {
	//CheckPrevious ensures the event order if true
	// if false the previous sequence is not checked on push
	CheckPrevious() bool
	//EditorService must return the name of the service which creates the new event
	EditorService() string
	//EditorUser must return the id of the user who created the event
	EditorUser() string
	//Type must return an event type which should be unique in the aggregate
	Type() EventType
	//Data returns the payload of the event. It represent the changed fields by the event
	// valid types are:
	// * nil (no payload),
	// * json byte array
	// * struct which can be marshalled to json
	// * pointer to struct which can be marshalled to json
	Data() interface{}
}

//Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	repo             repository.Repository
	interceptorMutex sync.Mutex
	eventMapper      map[EventType]eventTypeInterceptors
}

type eventTypeInterceptors struct {
	filterMapper func(*repository.Event) (Event, error)
}

//Health checks if the eventstore can properly work
// It checks if the repository can serve load
func (es *Eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

type aggregater interface {
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
	//PreviouseSequence should return the sequence of the latest event of this aggregate
	// stored in the eventstore
	// it's set to the first event of this push transaction,
	// later events consume the sequence of the previously pushed event of the aggregate
	PreviousSequence() uint64
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
			data, err := eventData(event)
			if err != nil {
				return nil, err
			}
			events = append(events, &repository.Event{
				AggregateID:           aggregate.ID(),
				AggregateType:         repository.AggregateType(aggregate.Type()),
				ResourceOwner:         aggregate.ResourceOwner(),
				EditorService:         event.EditorService(),
				EditorUser:            event.EditorUser(),
				Type:                  repository.EventType(event.Type()),
				Version:               repository.Version(aggregate.Version()),
				PreviousEvent:         previousEvent,
				PreviousSequence:      aggregate.PreviousSequence(),
				Data:                  data,
				CheckPreviousSequence: event.CheckPrevious(),
			})
			previousEvent = events[len(events)-1]
		}
	}
	return events, nil
}

//FilterEvents filters the stored events based on the searchQuery
// and maps the events to the defined event structs
func (es *Eventstore) FilterEvents(ctx context.Context, queryFactory *SearchQueryFactory) ([]Event, error) {
	query, err := queryFactory.build()
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

type reducer interface {
	//Reduce handles the events of the internal events list
	// it only appends the newly added events
	Reduce() error
	//AppendEvents appends the passed events to an internal list of events
	AppendEvents(...Event) error
}

//FilterToReducer filters the events based on the search query, appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToReducer(ctx context.Context, searchQuery *SearchQueryFactory, r reducer) error {
	events, err := es.FilterEvents(ctx, searchQuery)
	if err != nil {
		return err
	}
	if err = r.AppendEvents(events...); err != nil {
		return err
	}

	return r.Reduce()
}

//LatestSequence filters the latest sequence for the given search query
func (es *Eventstore) LatestSequence(ctx context.Context, queryFactory *SearchQueryFactory) (uint64, error) {
	query, err := queryFactory.build()
	if err != nil {
		return 0, err
	}
	return es.repo.LatestSequence(ctx, query)
}

//RegisterFilterEventMapper registers a function for mapping an eventstore event to an event
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

func eventData(event Event) ([]byte, error) {
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
			return nil, errors.ThrowInvalidArgument(err, "V2-xG87M", "could  not marhsal data")
		}
		return dataBytes, nil
	}
	return nil, errors.ThrowInvalidArgument(nil, "V2-91NRm", "wrong type of event data")
}
