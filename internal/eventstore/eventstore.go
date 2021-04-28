package eventstore

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

//Eventstore abstracts all functions needed to store valid events
// and filters the stored events
type Eventstore struct {
	repo              repository.Repository
	interceptorMutex  sync.Mutex
	eventInterceptors map[EventType]eventTypeInterceptors
}

type eventTypeInterceptors struct {
	eventMapper func(*repository.Event) (EventReader, error)
}

func NewEventstore(repo repository.Repository) *Eventstore {
	return &Eventstore{
		repo:              repo,
		eventInterceptors: map[EventType]eventTypeInterceptors{},
		interceptorMutex:  sync.Mutex{},
	}
}

//Health checks if the eventstore can properly work
// It checks if the repository can serve load
func (es *Eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

//PushEvents pushes the events in a single transaction
// an event needs at least an aggregate
func (es *Eventstore) PushEvents(ctx context.Context, pushEvents ...EventPusher) ([]EventReader, error) {
	events, assets, constraints, err := eventsToRepository(pushEvents)
	if err != nil {
		return nil, err
	}
	err = es.repo.Push(ctx, events, assets, constraints...)
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

func eventsToRepository(pushEvents []EventPusher) (events []*repository.Event, assets []*repository.Asset, constraints []*repository.UniqueConstraint, err error) {
	events = make([]*repository.Event, len(pushEvents))
	for i, event := range pushEvents {
		data, err := EventData(event)
		if err != nil {
			return nil, nil, nil, err
		}
		events[i] = &repository.Event{
			AggregateID:   event.Aggregate().ID,
			AggregateType: repository.AggregateType(event.Aggregate().Typ),
			ResourceOwner: event.Aggregate().ResourceOwner,
			EditorService: event.EditorService(),
			EditorUser:    event.EditorUser(),
			Type:          repository.EventType(event.Type()),
			Version:       repository.Version(event.Aggregate().Version),
			Data:          data,
		}
		if len(event.UniqueConstraints()) > 0 {
			constraints = append(constraints, uniqueConstraintsToRepository(event.UniqueConstraints())...)
		}
		if len(event.Assets()) > 0 {
			assets = append(assets, assetsToRepository(event.Assets())...)
		}
	}

	return events, assets, constraints, nil
}

func uniqueConstraintsToRepository(constraints []*EventUniqueConstraint) (uniqueConstraints []*repository.UniqueConstraint) {
	uniqueConstraints = make([]*repository.UniqueConstraint, len(constraints))
	for i, constraint := range constraints {
		uniqueConstraints[i] = &repository.UniqueConstraint{
			UniqueType:   constraint.UniqueType,
			UniqueField:  constraint.UniqueField,
			Action:       uniqueConstraintActionToRepository(constraint.Action),
			ErrorMessage: constraint.ErrorMessage,
		}
	}
	return uniqueConstraints
}

func assetsToRepository(assets []*Asset) (result []*repository.Asset) {
	result = make([]*repository.Asset, len(assets))
	for i, asset := range assets {
		result[i] = &repository.Asset{
			ID:     asset.ID,
			Asset:  asset.Asset,
			Action: assetActionToRepository(asset.Action),
		}
	}
	return result
}

//FilterEvents filters the stored events based on the searchQuery
// and maps the events to the defined event structs
func (es *Eventstore) FilterEvents(ctx context.Context, queryFactory *SearchQueryBuilder) ([]EventReader, error) {
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

func (es *Eventstore) mapEvents(events []*repository.Event) (mappedEvents []EventReader, err error) {
	mappedEvents = make([]EventReader, len(events))

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
	AppendEvents(...EventReader)
}

//FilterToReducer filters the events based on the search query, appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToReducer(ctx context.Context, searchQuery *SearchQueryBuilder, r reducer) error {
	events, err := es.FilterEvents(ctx, searchQuery)
	if err != nil {
		return err
	}

	r.AppendEvents(events...)

	return r.Reduce()
}

//LatestSequence filters the latest sequence for the given search query
func (es *Eventstore) LatestSequence(ctx context.Context, queryFactory *SearchQueryBuilder) (uint64, error) {
	query, err := queryFactory.build()
	if err != nil {
		return 0, err
	}
	return es.repo.LatestSequence(ctx, query)
}

type queryReducer interface {
	reducer
	//Query returns the SearchQueryFactory for the events needed in reducer
	Query() *SearchQueryBuilder
}

//FilterToQueryReducer filters the events based on the search query of the query function,
// appends all events to the reducer and calls it's reduce function
func (es *Eventstore) FilterToQueryReducer(ctx context.Context, r queryReducer) error {
	events, err := es.FilterEvents(ctx, r.Query())
	if err != nil {
		return err
	}
	r.AppendEvents(events...)

	return r.Reduce()
}

//RegisterFilterEventMapper registers a function for mapping an eventstore event to an event
func (es *Eventstore) RegisterFilterEventMapper(eventType EventType, mapper func(*repository.Event) (EventReader, error)) *Eventstore {
	if mapper == nil || eventType == "" {
		return es
	}
	es.interceptorMutex.Lock()
	defer es.interceptorMutex.Unlock()

	interceptor := es.eventInterceptors[eventType]
	interceptor.eventMapper = mapper
	es.eventInterceptors[eventType] = interceptor

	return es
}

func EventData(event EventPusher) ([]byte, error) {
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
	default:
		return repository.UniqueConstraintAdd
	}
}

func assetActionToRepository(action AssetAction) repository.AssetAction {
	switch action {
	case AssetAdd:
		return repository.AssetAdded
	case AssetRemove:
		return repository.AssetRemoved
	default:
		return repository.AssetAdded
	}
}
