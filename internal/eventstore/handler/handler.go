package handler

import (
	"github.com/caos/zitadel/internal/eventstore"
)

type Handler struct {
	Eventstore *eventstore.Eventstore
	Sub        *eventstore.Subscription
	EventQueue chan eventstore.EventReader
}

func NewHandler(es *eventstore.Eventstore) *Handler {
	h := Handler{
		Eventstore: es,
		EventQueue: make(chan eventstore.EventReader, 100),
	}

	return &h
}

func (h Handler) Subscribe(aggregates ...eventstore.AggregateType) {
	h.Sub = eventstore.Subscribe(h.EventQueue, aggregates...)
}
