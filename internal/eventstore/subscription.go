package eventstore

import "github.com/caos/zitadel/internal/eventstore/models"

type Subscription struct {
	es         *eventstore
	Events     chan *models.Event
	aggregates []models.AggregateType
}

func (es *eventstore) Subscribe(aggregates ...models.AggregateType) *Subscription {
	events := make(chan *models.Event, 100)
	sub := &Subscription{
		es:         es,
		Events:     events,
		aggregates: aggregates,
	}

	es.subsMutext.Lock()
	for _, aggregate := range aggregates {
		_, ok := es.subscriptions[aggregate]
		if !ok {
			es.subscriptions[aggregate] = make([]*Subscription, 1)
		}
		es.subscriptions[aggregate] = append(es.subscriptions[aggregate], sub)
	}

	return sub
}

func (es *eventstore) notify(aggregates []*models.Aggregate) {
	es.subsMutext.Lock()
	defer es.subsMutext.Unlock()
	for _, aggregate := range aggregates {
		subscriptions, ok := es.subscriptions[aggregate.Type()]
		if !ok {
			continue
		}
		for _, subsctiption := range subscriptions {
			for _, event := range aggregate.Events {
				subsctiption.Events <- event
			}
		}
	}
}

func (s *Subscription) Unsubscribe() {
	s.es.subsMutext.Lock()
	defer s.es.subsMutext.Unlock()
	for _, aggregate := range s.aggregates {
		subs, ok := s.es.subscriptions[aggregate]
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
}
