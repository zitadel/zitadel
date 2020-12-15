package eventstore

import "github.com/caos/zitadel/internal/eventstore/models"

type Subscription struct {
	events <-chan *models.Event
}

func (es *eventstore) Subscribe(aggregates ...models.AggregateType) *Subscription {
	//TODO: save to es
	events := make(chan *models.Event, 100)
	return &Subscription{
		events: events,
	}
}

func (es *eventstore) notify(aggregates []*models.Aggregate) {
	for _, aggregate := range aggregates {

	}
}
