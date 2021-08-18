package eventstore

import (
	"sync"

	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

var (
	subscriptions = map[AggregateType][]*Subscription{}
	subsMutext    sync.Mutex
)

type Subscription struct {
	Events     chan EventReader
	aggregates []AggregateType
}

func Subscribe(aggregates ...AggregateType) *Subscription {
	events := make(chan EventReader, 100)
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

func notify(events []EventReader) {
	go v1.Notify(MapEventsToV1Events(events))
	subsMutext.Lock()
	defer subsMutext.Unlock()
	for _, event := range events {
		subs, ok := subscriptions[event.Aggregate().Typ]
		if !ok {
			continue
		}
		for _, sub := range subs {
			sub.Events <- event
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
	_, ok := <-s.Events
	if ok {
		close(s.Events)
	}
}

func MapEventsToV1Events(events []EventReader) []*models.Event {
	v1Events := make([]*models.Event, len(events))
	for i, event := range events {
		v1Events[i] = mapEventToV1Event(event)
	}
	return v1Events
}

func mapEventToV1Event(event EventReader) *models.Event {
	return &models.Event{
		Sequence:      event.Sequence(),
		CreationDate:  event.CreationDate(),
		Type:          models.EventType(event.Type()),
		AggregateType: models.AggregateType(event.Aggregate().Typ),
		AggregateID:   event.Aggregate().ID,
		ResourceOwner: event.Aggregate().ResourceOwner,
		EditorService: event.EditorService(),
		EditorUser:    event.EditorUser(),
		Data:          event.DataAsBytes(),
	}
}
