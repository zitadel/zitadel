package handler

import "github.com/caos/zitadel/internal/eventstore"

type SubscriptionHandler interface {
	Handler
	Subscribe(...eventstore.AggregateType)
	SubscribeEvents(map[eventstore.AggregateType][]eventstore.EventType)
}

func NewSubscriptionHandler() SubscriptionHandler {
	return &subscriptionHandler{
		queue: make(chan eventstore.EventReader, 100),
	}
}

type subscriptionHandler struct {
	sub   *eventstore.Subscription
	queue chan eventstore.EventReader
}

func (h *subscriptionHandler) Subscribe(types ...eventstore.AggregateType) {
	h.sub = eventstore.SubscribeAggregates(h.queue, types...)
}

func (h *subscriptionHandler) SubscribeEvents(types map[eventstore.AggregateType][]eventstore.EventType) {
	h.sub = eventstore.SubscribeEventTypes(h.queue, types)
}
