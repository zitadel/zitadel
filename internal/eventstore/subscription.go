package eventstore

import (
	"context"
	"slices"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/zitadel/logging"
)

var (
	subscriptions = map[AggregateType][]*Subscription{}
	subsMutext    sync.Mutex
)

type Subscription struct {
	Trigger SubscriptionTrigger
	types   map[AggregateType][]EventType
}

type SubscriptionTrigger func(ctx context.Context, position decimal.Decimal) error

// SubscribeEventTypes subscribes for the given event types
// if no event types are provided the subscription is for all events of the aggregate
func SubscribeEventTypes(trigger SubscriptionTrigger, types map[AggregateType][]EventType) *Subscription {
	sub := &Subscription{
		Trigger: trigger,
		types:   types,
	}

	subsMutext.Lock()
	defer subsMutext.Unlock()

	for aggregate := range types {
		subscriptions[aggregate] = append(subscriptions[aggregate], sub)
	}

	return sub
}

func (es *Eventstore) notify(ctx context.Context, events []Event) <-chan bool {
	subsMutext.Lock()
	defer subsMutext.Unlock()
	var toNotify []*Subscription
	for _, event := range events {
		subs, ok := subscriptions[event.Aggregate().Type]
		if !ok {
			continue
		}
		for _, sub := range subs {
			eventTypes := sub.types[event.Aggregate().Type]
			//subscription for all events
			if len(eventTypes) == 0 {
				toNotify = append(toNotify, sub)
				continue
			}
			//subscription for certain events
			for _, eventType := range eventTypes {
				if event.Type() == eventType {
					toNotify = append(toNotify, sub)
					break
				}
			}
		}
	}
	var wg sync.WaitGroup
	for _, sub := range slices.Compact(toNotify) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := sub.Trigger(ctx, events[len(events)-1].Position())
			if err != nil {
				logging.WithError(err).Error("failed to trigger subscription")
			}
		}()
	}
	lock := make(chan bool, 1)
	go func() {
		wg.Wait()
		lock <- true
		close(lock)
	}()

	return lock
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
}
