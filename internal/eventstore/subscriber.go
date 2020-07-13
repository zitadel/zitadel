package eventstore

import (
	"sync"

	"github.com/caos/zitadel/internal/eventstore/models"
)

var (
	subscriberLock sync.Mutex
	subscribers    map[models.AggregateType]map[*subscriber]chan<- *models.Event = make(map[models.AggregateType]map[*subscriber]chan<- *models.Event)
)

type subscriber struct{}

func sendAggregates(aggregates []*models.Aggregate) {
	for _, aggregate := range aggregates {
		send(aggregate.Events...)
	}
}

func send(events ...*models.Event) {
	subscriberLock.Lock()
	defer subscriberLock.Unlock()
	for _, event := range events {
		subscribers, ok := subscribers[event.AggregateType]
		if !ok {
			continue
		}
		for _, channel := range subscribers {
			channel <- event
		}
	}
}

func (s *subscriber) Unsubscribe() {
	subscriberLock.Lock()
	defer subscriberLock.Unlock()

	for _, subs := range subscribers {
		delete(subs, s)
	}
}

func Subscribe(feed chan<- *models.Event, aggregateTypes ...models.AggregateType) *subscriber {
	subscriberLock.Lock()
	defer subscriberLock.Unlock()

	s := new(subscriber)

	for _, aggregateType := range aggregateTypes {
		_, ok := subscribers[aggregateType]
		if !ok {
			subscribers[aggregateType] = make(map[*subscriber]chan<- *models.Event)
		}
		subscribers[aggregateType][s] = feed
	}

	return s
}
