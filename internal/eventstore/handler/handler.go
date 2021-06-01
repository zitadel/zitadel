package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
)

type HandlerConfig struct {
	Eventstore *eventstore.Eventstore
}
type Handler struct {
	Eventstore *eventstore.Eventstore
	Sub        *eventstore.Subscription
	EventQueue chan eventstore.EventReader
}

func NewHandler(config HandlerConfig) Handler {
	return Handler{
		Eventstore: config.Eventstore,
		EventQueue: make(chan eventstore.EventReader, 100),
	}
}

func (h Handler) Subscribe(aggregates ...eventstore.AggregateType) {
	h.Sub = eventstore.SubscribeAggregates(h.EventQueue, aggregates...)
}

func (h Handler) SubscribeEvents(types map[eventstore.AggregateType][]eventstore.EventType) {
	h.Sub = eventstore.SubscribeEventTypes(h.EventQueue, types)
}
