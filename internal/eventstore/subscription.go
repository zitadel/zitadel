package eventstore

import (
	"sync"

	"github.com/zitadel/logging"
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

func (es *Eventstore) notify(events []Event) {
	eventReaders, err := es.mapEvents(events)
	if err != nil {
		logging.WithError(err).Debug("unable to map events")
		return
	}

	subsMutext.Lock()
	defer subsMutext.Unlock()
	for _, event := range eventReaders {
		subs, ok := subscriptions[event.Aggregate().Type]
		if !ok {
			continue
		}
		for _, sub := range subs {
			eventTypes := sub.types[event.Aggregate().Type]
			//subscription for all events
			if len(eventTypes) == 0 {
				sub.Events <- event
				continue
			}
			//subscription for certain events
			for _, eventType := range eventTypes {
				if event.Type() == eventType {
					sub.Events <- event
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

// func mapEventsToRepo(events []Event) []*repository.Event {
// 	repoEvents := make([]*repository.Event, len(events))
// 	for i, event := range events {
// 		repoEvents[i] = &repository.Event{
// 			Seq:           event.Sequence(),
// 			CreationDate:  event.CreatedAt(),
// 			Typ:           event.Type(),
// 			EditorService: "zitadel",
// 			EditorUser:    event.Creator(),
// 			Version:       repository.Version(event.Aggregate().Version),
// 			AggregateID:   event.Aggregate().ID,
// 			AggregateType: event.Aggregate().Type,
// 			ResourceOwner: sql.NullString{
// 				String: event.Aggregate().ResourceOwner,
// 				Valid:  event.Aggregate().ResourceOwner != "",
// 			},
// 			InstanceID: event.Aggregate().InstanceID,
// 			Data:       event.DataAsBytes(),
// 		}
// 	}
// 	return repoEvents
// }
