package eventstore

import (
	"encoding/json"
	"sync"

	"github.com/zitadel/logging"

	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

var (
	subscriptions = map[AggregateType][]*Subscription{}
	subsMutext    sync.Mutex
)

type Subscription struct {
	Events chan Event
	types  map[AggregateType][]EventType
}

// SubscribeAggregates subscribes for all events on the given aggregates
func SubscribeAggregates(eventQueue chan Event, aggregates ...AggregateType) *Subscription {
	types := make(map[AggregateType][]EventType, len(aggregates))
	for _, aggregate := range aggregates {
		types[aggregate] = nil
	}
	sub := &Subscription{
		Events: eventQueue,
		types:  types,
	}

	subsMutext.Lock()
	defer subsMutext.Unlock()

	for _, aggregate := range aggregates {
		subscriptions[aggregate] = append(subscriptions[aggregate], sub)
	}

	return sub
}

// SubscribeEventTypes subscribes for the given event types
// if no event types are provided the subscription is for all events of the aggregate
func SubscribeEventTypes(eventQueue chan Event, types map[AggregateType][]EventType) *Subscription {
	aggregates := make([]AggregateType, len(types))
	sub := &Subscription{
		Events: eventQueue,
		types:  types,
	}

	subsMutext.Lock()
	defer subsMutext.Unlock()

	for _, aggregate := range aggregates {
		subscriptions[aggregate] = append(subscriptions[aggregate], sub)
	}

	return sub
}

func notify(events []eventstore.Event) {
	go v1.Notify(MapEventsToV1Events(events))
	subsMutext.Lock()
	defer subsMutext.Unlock()
	for _, event := range events {
		e := event.(Event)
		subs, ok := subscriptions[event.Aggregate().Type]
		if !ok {
			continue
		}
		for _, sub := range subs {
			eventTypes := sub.types[event.Aggregate().Type]
			//subscription for all events
			if len(eventTypes) == 0 {
				sub.Events <- e
				continue
			}
			//subscription for certain events
			for _, eventType := range eventTypes {
				if event.Type() == eventType {
					sub.Events <- e
					break
				}
			}
		}
	}
}

func (s *Subscription) Unsubscribe() {
	subsMutext.Lock()
	defer subsMutext.Unlock()
	for aggregate := range s.types {
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

func MapEventsToV1Events(events []eventstore.Event) []*models.Event {
	v1Events := make([]*models.Event, len(events))
	for i, event := range events {
		v1Events[i] = mapEventToV1Event(event)
	}
	return v1Events
}

func mapEventToV1Event(event eventstore.Event) *models.Event {
	payload := make(map[string]any)
	err := event.Unmarshal(&payload)
	logging.OnError(err).Debug("unmarshal failed")
	data, err := json.Marshal(payload)
	logging.OnError(err).Debug("marshal failed")

	return &models.Event{
		Sequence:      event.Sequence(),
		CreationDate:  event.CreatedAt(),
		Type:          models.EventType(event.Type()),
		AggregateType: models.AggregateType(event.Aggregate().Type),
		AggregateID:   event.Aggregate().ID,
		ResourceOwner: event.Aggregate().ResourceOwner,
		InstanceID:    event.Aggregate().InstanceID,
		EditorService: "zitadel",
		EditorUser:    event.Creator(),
		Data:          data,
	}
}
