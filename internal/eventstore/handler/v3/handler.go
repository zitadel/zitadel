package handler

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/handler"
)

type Handler struct {
	*Iterator
	*Subscriber
	*Pusher
}

type HandlerConfig struct {
	ProjectionName string
	Reducers       []handler.AggregateReducer

	IteratorConfig
	SubscriberConfig
	PusherConfig
}

func NewHandler(config HandlerConfig) Handler {
	config.PusherConfig.projectionName = config.ProjectionName
	config.IteratorConfig.projectionName = config.ProjectionName
	config.SubscriberConfig.projectionName = config.ProjectionName

	config.IteratorConfig.reducers = config.Reducers
	config.SubscriberConfig.reducers = config.Reducers

	pusher := NewPusher(config.PusherConfig)

	config.IteratorConfig.pusher = pusher
	config.SubscriberConfig.pusher = pusher

	iterator := NewIterator(config.IteratorConfig)
	subscriber := NewSubscriber(config.SubscriberConfig)

	return Handler{
		Pusher:     pusher,
		Iterator:   iterator,
		Subscriber: subscriber,
	}
}

func (h *Handler) Project(ctx context.Context) {
	go h.Pusher.Process(ctx)
	go h.Iterator.Process(ctx)
	go h.Subscriber.Process(ctx)
}
