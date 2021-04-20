package eventstore

import (
	"context"
	"sync"
	"time"

	"github.com/caos/logging"
)

type Handler struct {
	ctx context.Context

	Eventstore *Eventstore
	Sub        *Subscription
	EventQueue chan EventReader
}

func NewHandler(eventstore *Eventstore) *Handler {
	h := Handler{
		Eventstore: eventstore,
		//TODO: how huge should the queue be?
		EventQueue: make(chan EventReader, 100),
	}
	//TODO: start handler for EventQueue

	return &h
}

func (h Handler) Subscribe(aggregates ...AggregateType) {
	h.Sub = Subscribe(h.EventQueue, aggregates...)
}

type ReadModelHandler struct {
	Handler
	RequeueAfter time.Duration
	Timer        *time.Timer

	lock       sync.Mutex
	stmts      []Statement
	pushSet    bool
	shouldPush chan *struct{}
}

func NewReadModelHandler(
	ctx context.Context,
	eventstore *Eventstore,
	requeueAfter time.Duration,
) *ReadModelHandler {
	return &ReadModelHandler{
		Handler:      *NewHandler(eventstore),
		RequeueAfter: requeueAfter,
		// first requeue is instant on startup
		Timer: time.NewTimer(0),
	}
}

func (h *ReadModelHandler) ResetTimer() {
	h.Timer.Reset(h.RequeueAfter)
}

type Update func([]Statement) error
type Reduce func(EventReader) (Statement, error)

func (h *ReadModelHandler) Process(
	ctx context.Context,
	reduce Reduce,
	update Update,
) {
	for {
		select {
		case <-ctx.Done():
			h.Sub.Unsubscribe()
			logging.Log("EVENT-XG5Og").Info("stop processing")
			return
		case event := <-h.Handler.EventQueue:
			stmt, err := reduce(event)
			if err != nil {
				logging.Log("EVENT-PTr4j").OnError(err).Warn("unable to process event")
				continue
			}

			h.lock.Lock()
			defer h.lock.Unlock()

			h.stmts = append(h.stmts, stmt)
			if !h.pushSet {
				h.pushSet = true
				h.shouldPush <- nil
			}
		case <-h.Timer.C:
			//TODO: bulk run
		default:
			//continue to lower prio select
		}

		// if not canceled and no events push is allowed
		select {
		case <-ctx.Done():
			logging.Log("EVENT-XG5Og").Info("stop processing")
			return
		case event := <-h.Handler.EventQueue:
			stmt, err := reduce(event)
			logging.Log("EVENT-PTr4j").WithError(err).Warn("unable to process event")
			if err == nil {
				h.stmts = append(h.stmts, stmt)
			}
		case <-h.Timer.C:
			//TODO: bulk run
		case <-h.shouldPush:
			h.lock.Lock()
			defer h.lock.Unlock()

			h.pushSet = false
			if err := update(h.stmts); err != nil {
				logging.Log("EVENT-EFDwe").WithError(err).Warn("unable to push")
				continue
			}

			h.stmts = []Statement{}
		}
	}
}
