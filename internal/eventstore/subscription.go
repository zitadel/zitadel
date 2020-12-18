package eventstore

import (
	"sync"

	"github.com/caos/zitadel/internal/eventstore/models"
)

var (
	subscriptions map[models.AggregateType][]*Subscription = map[models.AggregateType][]*Subscription{}
	subsMutext    sync.Mutex
)

type Subscription struct {
	Events     chan *models.Event
	aggregates []models.AggregateType
}

func (es *eventstore) Subscribe(aggregates ...models.AggregateType) *Subscription {
	events := make(chan *models.Event, 100)
	sub := &Subscription{
		Events:     events,
		aggregates: aggregates,
	}

	subsMutext.Lock()
	defer subsMutext.Unlock()

	for _, aggregate := range aggregates {
		_, ok := subscriptions[aggregate]
		if !ok {
			subscriptions[aggregate] = make([]*Subscription, 0, 1)
		}
		subscriptions[aggregate] = append(subscriptions[aggregate], sub)
	}

	return sub
}

func notify(aggregates []*models.Aggregate) {
	subsMutext.Lock()
	defer subsMutext.Unlock()
	for _, aggregate := range aggregates {
		subs, ok := subscriptions[aggregate.Type()]
		if !ok {
			continue
		}
		for _, sub := range subs {
			for _, event := range aggregate.Events {
				sub.Events <- event
			}
		}
	}
}

func (s *Subscription) Unsubscribe() {
	subsMutext.Lock()
	defer subsMutext.Unlock()
	for _, aggregate := range s.aggregates {
		subs, ok := subscriptions[aggregate]
		if !ok {
			continue
		}
		for i := len(subs) - 1; i >= 0; i-- {
			if subs[i] == s {
				subs[i] = subs[len(subs)-1]
				subs[len(subs)-1] = nil
				subs = subs[:len(subs)-1]
			}
		}
	}
	close(s.Events)
}
