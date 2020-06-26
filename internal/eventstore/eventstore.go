package eventstore

import (
	"context"
	"sync"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/internal/repository"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Eventstore interface {
	AggregateCreator() *models.AggregateCreator
	Health(ctx context.Context) error
	PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) error
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
}

var _ Eventstore = (*eventstore)(nil)

type eventstore struct {
	repo             repository.Repository
	aggregateCreator *models.AggregateCreator
	subscriberLock   sync.Mutex
	subscribers      map[models.EventType][]*subscriber
	subscriberIdx    int
}

func (es *eventstore) AggregateCreator() *models.AggregateCreator {
	return es.aggregateCreator
}

func (es *eventstore) PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	for _, aggregate := range aggregates {
		if len(aggregate.Events) == 0 {
			return errors.ThrowInvalidArgument(nil, "EVENT-cNhIj", "no events in aggregate")
		}
		for _, event := range aggregate.Events {
			if err = event.Validate(); err != nil {
				return errors.ThrowInvalidArgument(err, "EVENT-tzIhl", "validate event failed")
			}
		}
	}
	err = es.repo.PushAggregates(ctx, aggregates...)
	if err != nil {
		return err
	}
	go es.sendAggregates(aggregates)
	return nil
}

func (es *eventstore) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) ([]*models.Event, error) {
	if err := searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, searchQuery)
}

func (es *eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

func (es *eventstore) sendAggregates(aggregates []*models.Aggregate) {
	for _, aggregate := range aggregates {
		es.send(aggregate.Events...)
	}
}

func (es *eventstore) send(events ...*models.Event) {
	es.subscriberLock.Lock()
	defer es.subscriberLock.Unlock()
	for _, event := range events {
		subscribers, ok := es.subscribers[event.Type]
		if !ok {
			return
		}
		for _, subscriber := range subscribers {
			subscriber.feed <- event
		}
	}
}

type subscriber struct {
	index      int
	eventTypes []models.EventType
	feed       chan<- *models.Event
}

func (es *eventstore) Unsubscribe(s *subscriber) {
	es.subscriberLock.Lock()
	defer es.subscriberLock.Unlock()

	for _, eventType := range s.eventTypes {
		_, ok := es.subscribers[eventType]
		if !ok {
			continue
		}
		for i, subscriber := range es.subscribers[eventType] {
			if subscriber.index != s.index {
				continue
			}
			es.subscribers[eventType][i] = es.subscribers[eventType][len(es.subscribers[eventType])-1]
			es.subscribers[eventType] = es.subscribers[eventType][:len(es.subscribers[eventType])-1]
		}
	}
}

func (es *eventstore) Subscribe(feed chan<- *models.Event, eventTypes ...models.EventType) *subscriber {
	es.subscriberLock.Lock()
	defer es.subscriberLock.Unlock()

	es.subscriberIdx++
	s := &subscriber{
		index:      es.subscriberIdx,
		eventTypes: eventTypes,
		feed:       feed,
	}

	for _, eventType := range eventTypes {
		_, ok := es.subscribers[eventType]
		if !ok {
			es.subscribers[eventType] = make([]*subscriber, 0)
		}
		es.subscribers[eventType] = append(es.subscribers[eventType], s)
	}

	return s
}
