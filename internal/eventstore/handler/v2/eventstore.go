package handler

import (
	"context"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
)

type EventstoreQuerier interface {
	Query(ctx context.Context, queue chan eventstore.EventReader)
	Cancel()
}

type SubscriptionHandler struct {
	sub   *eventstore.Subscription
	types []eventstore.AggregateType
}

func NewSubscriptionHandler(types []eventstore.AggregateType) EventstoreQuerier {
	return &SubscriptionHandler{
		types: types,
	}
}

func (h *SubscriptionHandler) Query(ctx context.Context, queue chan eventstore.EventReader) {
	h.sub = eventstore.SubscribeAggregates(queue, h.types...)
}

func (h *SubscriptionHandler) Cancel() {
	h.sub.Unsubscribe()
}

type IterationHandler struct {
	Handler
	es       *eventstore.Eventstore
	interval time.Duration
	errs     chan<- error
	cancel   func()
}

type IterationHandlerConfig struct {
	Eventstore *eventstore.Eventstore
	Interval   time.Duration
	PreSteps   []PreStep
	PostSteps  []PostStep
}

func NewIterationHandler(config IterationHandlerConfig) EventstoreQuerier {
	// TODO: validate config
	return &IterationHandler{
		es:       config.Eventstore,
		interval: config.Interval,
		Handler: Handler{
			preSteps:  config.PreSteps,
			postSteps: config.PostSteps,
		},
	}
}

func (h *IterationHandler) Query(ctx context.Context, queue chan eventstore.EventReader) {
	ctx, h.cancel = context.WithCancel(ctx)

	t := time.NewTimer(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				logging.Log("V2-HDCmY").Debug("stop iterating")
				return
			case <-t.C:
				if err := h.execPreSteps(); err != nil {
					logging.Log("V2-5GsNd").WithError(err).Warn("pre step failed")
					h.execPostSteps()
					break
				}

				query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent)
				events, err := h.es.FilterEvents(ctx, query)
				if err != nil {
					h.reportErr(err)
					h.execPostSteps()
					break
				}
				for _, event := range events {
					queue <- event
				}
				h.execPostSteps()
			}
			t.Reset(h.interval)
		}
	}()
}

func (h *IterationHandler) Cancel() {
	h.cancel()
}

func (h *IterationHandler) reportErr(err error) {
	if h.errs != nil {
		h.errs <- err
		return
	}
	logging.Log("V2-i1jmN").WithError(err).Warn("error in iteration")
}
