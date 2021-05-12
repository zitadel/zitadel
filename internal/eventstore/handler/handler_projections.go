package handler

import (
	"context"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
)

type ProjectionHandler struct {
	Handler
	RequeueAfter  time.Duration
	Timer         *time.Timer
	SequenceTable string

	lock       sync.Mutex
	stmts      []Statement
	pushSet    bool
	shouldPush chan *struct{}
}

func NewProjectionHandler(
	eventstore *eventstore.Eventstore,
	requeueAfter time.Duration,
) *ProjectionHandler {
	return &ProjectionHandler{
		Handler:      NewHandler(eventstore),
		RequeueAfter: requeueAfter,
		// first bulk is instant on startup
		Timer:      time.NewTimer(0),
		shouldPush: make(chan *struct{}, 1),
	}
}

func (h *ProjectionHandler) ResetTimer() {
	h.Timer.Reset(h.RequeueAfter)
}

//Update updates the projection with the given statements
type Update func(context.Context, []Statement) error

//Reduce reduces the given event to a statement
//which is used to update the projection
type Reduce func(eventstore.EventReader) ([]Statement, error)

//Lock is used for mutex handling if needed on the projection
type Lock func() error

//Unlock releases the mutex of the projection
type Unlock func() error

//Process waits for several conditions:
// if context is canceled the function gracefully shuts down
// if an event occures it reduces the event
// if the internal timer expires the handler will check
// 	for unprocessed events on eventstore
func (h *ProjectionHandler) Process(
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
			if h.pushSet {
				h.push(context.Background(), update)
			}
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
			if h.pushSet {
				h.push(context.Background(), update)
			}
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

func (h *ProjectionHandler) processEvent(ctx context.Context, event eventstore.EventReader, reduce Reduce) {
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

func (h *ProjectionHandler) bulk(ctx context.Context, lock Lock, query *eventstore.SearchQueryBuilder, reduce Reduce, update Update, unlock Unlock) {
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

func (h *ProjectionHandler) push(ctx context.Context, update Update) {
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

func (h *ProjectionHandler) shutdown() {
	h.Sub.Unsubscribe()
	logging.Log("EVENT-XG5Og").Info("stop processing")
}
