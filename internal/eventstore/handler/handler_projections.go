package handler

import (
	"context"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
)

//Update updates the projection with the given statements
type Update func(context.Context, []Statement, Reduce) error

//Reduce reduces the given event to a statement
//which is used to update the projection
type Reduce func(eventstore.EventReader) ([]Statement, error)

//Lock is used for mutex handling if needed on the projection
type Lock func(context.Context, chan error, time.Duration)

//Unlock releases the mutex of the projection
type Unlock func() error

//SearchQuery generates the search query to lookup for events
type SearchQuery func() (query *eventstore.SearchQueryBuilder, queryLimit uint64, err error)

type ProjectionHandler struct {
	Handler
	RequeueAfter  time.Duration
	Timer         *time.Timer
	SequenceTable string

	lockMu     sync.Mutex
	stmts      []Statement
	pushSet    bool
	shouldPush chan *struct{}

	reduce Reduce
	update Update
	lock   Lock
	unlock Unlock
	query  SearchQuery
}

func NewProjectionHandler(
	eventstore *eventstore.Eventstore,
	requeueAfter time.Duration,
	reduce Reduce,
	update Update,
	lock Lock,
	unlock Unlock,
	query SearchQuery,
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

//Process waits for several conditions:
// if context is canceled the function gracefully shuts down
// if an event occures it reduces the event
// if the internal timer expires the handler will check
// for unprocessed events on eventstore
func (h *ProjectionHandler) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if h.pushSet {
				h.push(context.Background())
			}
			h.shutdown()
			return
		case event := <-h.Handler.EventQueue:
			h.processEvent(ctx, event)
		case <-h.Timer.C:
			h.bulk(ctx)
			h.ResetTimer()
		default:
			//lower prio select with push
			select {
			case <-ctx.Done():
				if h.pushSet {
					h.push(context.Background())
				}
				h.shutdown()
				return
			case event := <-h.Handler.EventQueue:
				h.processEvent(ctx, event)
			case <-h.Timer.C:
				// lock, h.prepareBulk(query, reduce, update), unlock
				h.bulk(ctx)
				h.ResetTimer()
			case <-h.shouldPush:
				h.push(ctx)
				h.ResetTimer()
			}
		}
	}
}

func (h *ProjectionHandler) processEvent(ctx context.Context, event eventstore.EventReader) error {
	stmts, err := h.reduce(event)
	if err != nil {
		logging.Log("EVENT-PTr4j").WithError(err).Warn("unable to process event")
		return err
	}

	h.lockMu.Lock()
	defer h.lockMu.Unlock()

	h.stmts = append(h.stmts, stmts...)

	if !h.pushSet {
		h.pushSet = true
		h.shouldPush <- nil
	}
	return nil
}

func (h *ProjectionHandler) bulk(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	errs := make(chan error)
	defer func() {
		cancel()
		close(errs)
	}()

	h.lock(ctx, errs, h.RequeueAfter)
	//wait until projection is locked
	if err := <-errs; err != nil {
		logging.Log("HANDL-XDJ4i").WithError(err).Warn("initial lock failed")
		return
	}

	go cancelOnErr(ctx, errs, cancel)

	h.executeBulk(ctx)

	err := h.unlock()
	logging.Log("EVENT-boPv1").OnError(err).Warn("unable to unlock")
}

func cancelOnErr(ctx context.Context, errs chan error, cancel func()) {
	select {
	case err := <-errs:
		if err != nil {
			logging.Log("HANDL-cVop2").WithError(err).Warn("bulk canceled")
			cancel()
			return
		}
	case <-ctx.Done():
		return
	}
}

func (h *ProjectionHandler) executeBulk(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			hasLimitExeeded, err := h.prepareBulkStmts(ctx)
			if err != nil || len(h.stmts) == 0 {
				return
			}

			<-h.shouldPush
			if err = h.push(ctx); err != nil {
				logging.Log("EVENT-EFDwe").WithError(err).Warn("unable to push")
				return
			}

			if !hasLimitExeeded {
				return
			}
		}
	}
}

func (h *ProjectionHandler) prepareBulkStmts(ctx context.Context) (limitExeeded bool, err error) {
	eventQuery, eventsLimit, err := h.query()
	if err != nil {
		logging.Log("HANDL-x6qvs").WithError(err).Warn("unable to create event query")
		return false, err
	}

	events, err := h.Eventstore.FilterEvents(ctx, eventQuery)
	if err != nil {
		logging.Log("EVENT-ACMMS").WithError(err).Info("Unable to filter events in batch job")
		return false, err
	}
	for _, event := range events {
		h.processEvent(ctx, event)
	}

	return len(events) == int(eventsLimit), nil
}

func (h *ProjectionHandler) push(ctx context.Context) error {
	h.lockMu.Lock()
	defer h.lockMu.Unlock()

	h.pushSet = false
	err := h.update(ctx, h.stmts, h.reduce)
	h.stmts = nil

	return err
}

func (h *ProjectionHandler) shutdown() {
	h.Sub.Unsubscribe()
	logging.Log("EVENT-XG5Og").Info("stop processing")
}
