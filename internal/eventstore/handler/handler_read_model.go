package handler

import (
	"context"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
)

type ReadModelHandler struct {
	Handler
	RequeueAfter  time.Duration
	Timer         *time.Timer
	SequenceTable string

	lock       sync.Mutex
	stmts      []Statement
	pushSet    bool
	shouldPush chan *struct{}
}

func NewReadModelHandler(
	eventstore *eventstore.Eventstore,
	requeueAfter time.Duration,
) *ReadModelHandler {
	return &ReadModelHandler{
		Handler:      *NewHandler(eventstore),
		RequeueAfter: requeueAfter,
		// first requeue is instant on startup
		Timer:      time.NewTimer(0),
		shouldPush: make(chan *struct{}, 1),
	}
}

func (h *ReadModelHandler) ResetTimer() {
	h.Timer.Reset(h.RequeueAfter)
}

type Update func(context.Context, []Statement) error
type Reduce func(eventstore.EventReader) ([]Statement, error)
type Lock func() error
type Unlock func() error

func (h *ReadModelHandler) Process(
	ctx context.Context,
	reduce Reduce,
	update Update,
	lock Lock,
	unlock Unlock,
	query *eventstore.SearchQueryBuilder,
) {
	for {
		select {
		case <-ctx.Done():
			h.shutdown()
			return
		case event := <-h.Handler.EventQueue:
			h.processEvent(ctx, event, reduce)
		case <-h.Timer.C:
			h.bulk(ctx, lock, query, reduce, update, unlock)
			h.ResetTimer()
		default:
			//continue to lower prio select
		}

		// allow push
		select {
		case <-ctx.Done():
			h.shutdown()
			return
		case event := <-h.Handler.EventQueue:
			h.processEvent(ctx, event, reduce)
		case <-h.Timer.C:
			h.bulk(ctx, lock, query, reduce, update, unlock)
			h.ResetTimer()
		case <-h.shouldPush:
			h.push(ctx, update)
			h.ResetTimer()
		}
	}
}

func (h *ReadModelHandler) processEvent(ctx context.Context, event eventstore.EventReader, reduce Reduce) {
	stmts, err := reduce(event)
	if err != nil {
		logging.Log("EVENT-PTr4j").WithError(err).Warn("unable to process event")
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	h.stmts = append(h.stmts, stmts...)

	if !h.pushSet {
		h.pushSet = true
		h.shouldPush <- nil
	}
}

func (h *ReadModelHandler) bulk(ctx context.Context, lock Lock, query *eventstore.SearchQueryBuilder, reduce Reduce, update Update, unlock Unlock) {
	start := time.Now()
	if err := lock(); err != nil {
		logging.Log("EVENT-G76ye").WithError(err).Info("unable to lock")
		return
	}

	events, err := h.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-ACMMS").WithError(err).Info("Unable to filter events in batch job")
		return
	}
	for _, event := range events {
		h.processEvent(ctx, event, reduce)
	}
	<-h.shouldPush
	h.push(ctx, update)

	err = unlock()
	logging.Log("EVENT-boPv1").OnError(err).Warn("unable to unlock")
	logging.LogWithFields("HANDL-MDgQc", "start", start, "end", time.Now(), "diff", time.Now().Sub(start)).Warn("bulk")
}

func (h *ReadModelHandler) push(ctx context.Context, update Update) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.pushSet = false
	start := time.Now()
	if err := update(ctx, h.stmts); err != nil {
		logging.Log("EVENT-EFDwe").WithError(err).Warn("unable to push")
		return
	}
	logging.LogWithFields("HANDL-j5vuD", "start", start, "end", time.Now(), "diff", time.Now().Sub(start)).Warn("update")

	h.stmts = nil
}

func (h *ReadModelHandler) shutdown() {
	h.Sub.Unsubscribe()
	logging.Log("EVENT-XG5Og").Info("stop processing")
}
