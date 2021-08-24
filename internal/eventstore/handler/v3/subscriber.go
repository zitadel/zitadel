package handler

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
)

type Subscriber struct {
	projectionName string
	pusher         *Pusher
	reducers       map[eventstore.EventType]handler.Reduce
}

type SubscriberConfig struct {
	projectionName string
	pusher         *Pusher
	reducers       []handler.AggregateReducer
}

func NewSubscriber(config SubscriberConfig) *Subscriber {
	reducers := make(map[eventstore.EventType]handler.Reduce, len(config.reducers))
	for _, aggReducer := range config.reducers {
		for _, eventReducer := range aggReducer.EventRedusers {
			reducers[eventReducer.Event] = eventReducer.Reduce
		}
	}

	return &Subscriber{
		projectionName: config.projectionName,
		pusher:         config.pusher,
		reducers:       reducers,
	}
}

func (s *Subscriber) Process(ctx context.Context) {
	queue := make(chan eventstore.EventReader)
	sub := eventstore.SubscribeAggregates(queue)
	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return
		case event := <-queue:
			reduce, ok := s.reducers[event.Type()]
			if !ok {
				s.pusher.appendStmts(NewNoOpStatement(event))
				continue
			}

			stmts, err := reduce(event)
			logging.LogWithFields("V3-ipDkK", "projection", s.projectionName, "seq", event.Sequence()).OnError(err).Fatal("reduce failed")
			s.pusher.appendStmts(stmts...)
		}
	}
}
