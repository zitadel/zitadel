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
type Lock func(context.Context, chan error, time.Duration)

//Unlock releases the mutex of the projection
type Unlock func() error

//SearchQuery generates the search query to lookup for events
type SearchQuery func() (query *eventstore.SearchQueryBuilder, maxEvents uint64, err error)

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
	query SearchQuery,
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
			//lower prio select with push
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

func (h *ProjectionHandler) bulk(ctx context.Context, lock Lock, query SearchQuery, reduce Reduce, update Update, unlock Unlock) {
	ctx, cancel := context.WithCancel(ctx)
	errs := make(chan error)
	moreEvents := make(chan *struct{}, 1)
	defer func() {
		cancel()
		close(moreEvents)
		close(errs)
	}()

	lock(ctx, errs, h.RequeueAfter)
	//wait until projection is locked
	if err := <-errs; err != nil {
		logging.Log("HANDL-XDJ4i").WithError(err).Warn("initial lock failed")
		return
	}

	go func() {
		select {
		case err := <-errs:
			if err != nil {
				logging.Log("HANDL-cVop2").WithError(err).Warn("bulk canceled")
				cancel()
			}
		case <-ctx.Done():
			return
		}
	}()

	moreEvents <- nil

	//TODO: find solution without label
eventHandling:
	for {
		select {
		case <-ctx.Done():
			break eventHandling
		case <-moreEvents:
			eventQuery, maxEvents, err := query()
			if err != nil {
				logging.Log("HANDL-x6qvs").WithError(err).Warn("unable to create event query")
				return
			}
			events, err := h.Eventstore.FilterEvents(ctx, eventQuery)
			if err != nil {
				logging.Log("EVENT-ACMMS").WithError(err).Info("Unable to filter events in batch job")
				break eventHandling
			}
			for _, event := range events {
				h.processEvent(ctx, event, reduce)
			}

			if len(h.stmts) == 0 {
				break eventHandling
			}
			<-h.shouldPush
			if err = h.push(ctx, update); err != nil {
				logging.Log("EVENT-EFDwe").WithError(err).Warn("unable to push")
				break eventHandling
			}

			if len(events) < int(maxEvents) {
				break eventHandling
			}
			moreEvents <- nil
		default:
			break eventHandling
		}
	}
	err := unlock()
	logging.Log("EVENT-boPv1").OnError(err).Warn("unable to unlock")
}

func (h *ProjectionHandler) push(ctx context.Context, update Update) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.pushSet = false
	err := update(ctx, h.stmts)
	h.stmts = nil
	if err != nil {
		return err
	}

	return nil
}

func (h *ProjectionHandler) shutdown() {
	h.Sub.Unsubscribe()
	logging.Log("EVENT-XG5Og").Info("stop processing")
}
