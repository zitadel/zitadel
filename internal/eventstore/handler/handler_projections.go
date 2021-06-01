package handler

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
)

type ProjectionHandlerConfig struct {
	HandlerConfig
	ProjectionName string
	RequeueEvery   time.Duration
}

//Update updates the projection with the given statements
type Update func(context.Context, []Statement, Reduce) (unexecutedStmts []Statement, err error)

//Reduce reduces the given event to a statement
//which is used to update the projection
type Reduce func(eventstore.EventReader) ([]Statement, error)

//Lock is used for mutex handling if needed on the projection
type Lock func(context.Context, time.Duration) <-chan error

//Unlock releases the mutex of the projection
type Unlock func() error

//SearchQuery generates the search query to lookup for events
type SearchQuery func() (query *eventstore.SearchQueryBuilder, queryLimit uint64, err error)

type ProjectionHandler struct {
	Handler
	RequeueAfter time.Duration
	Timer        *time.Timer

	ProjectionName string

	lockMu     sync.Mutex
	stmts      []Statement
	pushSet    bool
	shouldPush chan *struct{}
}

func NewProjectionHandler(config ProjectionHandlerConfig) *ProjectionHandler {
	h := &ProjectionHandler{
		Handler:        NewHandler(config.HandlerConfig),
		ProjectionName: config.ProjectionName,
		RequeueAfter:   config.RequeueEvery,
		// first bulk is instant on startup
		Timer:      time.NewTimer(1 * time.Second),
		shouldPush: make(chan *struct{}, 1),
	}

	if config.RequeueEvery <= 0 {
		if !h.Timer.Stop() {
			<-h.Timer.C
		}
		logging.LogWithFields("HANDL-mC9Xx", "projection", h.ProjectionName).Debug("starting handler without requeue")
		return h
	}
	logging.LogWithFields("HANDL-fAC5O", "projection", h.ProjectionName).Debug("starting handler")
	return h
}

func (h *ProjectionHandler) ResetTimer() {
	if h.RequeueAfter > 0 {
		h.Timer.Reset(h.RequeueAfter)
	}
}

//Process waits for several conditions:
// if context is canceled the function gracefully shuts down
// if an event occures it reduces the event
// if the internal timer expires the handler will check
// for unprocessed events on eventstore
func (h *ProjectionHandler) Process(
	ctx context.Context,
	reduce Reduce,
	update Update,
	lock Lock,
	unlock Unlock,
	query SearchQuery,
) {
	execBulk := h.prepareExecuteBulk(query, reduce, update)
	for {
		select {
		case <-ctx.Done():
			if h.pushSet {
				h.push(context.Background(), update, reduce)
			}
			h.shutdown()
			return
		case event := <-h.Handler.EventQueue:
			h.processEvent(ctx, event, reduce)
		case <-h.Timer.C:
			h.bulk(ctx, lock, execBulk, unlock)
			h.ResetTimer()
		default:
			//lower prio select with push
			select {
			case <-ctx.Done():
				if h.pushSet {
					h.push(context.Background(), update, reduce)
				}
				h.shutdown()
				return
			case event := <-h.Handler.EventQueue:
				h.processEvent(ctx, event, reduce)
			case <-h.Timer.C:
				h.bulk(ctx, lock, execBulk, unlock)
				h.ResetTimer()
			case <-h.shouldPush:
				h.push(ctx, update, reduce)
				h.ResetTimer()
			}
		}
	}
}

func (h *ProjectionHandler) processEvent(
	ctx context.Context,
	event eventstore.EventReader,
	reduce Reduce,
) error {
	stmts, err := reduce(event)
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

func (h *ProjectionHandler) bulk(
	ctx context.Context,
	lock Lock,
	executeBulk executeBulk,
	unlock Unlock,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := lock(ctx, h.RequeueAfter)
	//wait until projection is locked or ctx canceled
	if err, ok := <-errs; err != nil || !ok {
		logging.Log("HANDL-XDJ4i").OnError(err).Warn("initial lock failed")
		return err
	}
	go cancelOnErr(ctx, errs, cancel)

	execErr := executeBulk(ctx)

	unlockErr := unlock()
	logging.Log("EVENT-boPv1").OnError(unlockErr).Warn("unable to unlock")

	if execErr != nil {
		return execErr
	}

	return unlockErr
}

func cancelOnErr(ctx context.Context, errs <-chan error, cancel func()) {
	for {
		select {
		case err := <-errs:
			if err != nil {
				logging.Log("HANDL-cVop2").WithError(err).Warn("bulk canceled")
				cancel()
				return
			}
		case <-ctx.Done():
			cancel()
			return
		}

	}
}

type executeBulk func(ctx context.Context) error

func (h *ProjectionHandler) prepareExecuteBulk(
	query SearchQuery,
	reduce Reduce,
	update Update,
) executeBulk {
	return func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				hasLimitExeeded, err := h.fetchBulkStmts(ctx, query, reduce)
				if err != nil || len(h.stmts) == 0 {
					return err
				}

				<-h.shouldPush
				if err = h.push(ctx, update, reduce); err != nil {
					logging.Log("EVENT-EFDwe").WithError(err).Warn("unable to push")
					return err
				}

				if !hasLimitExeeded {
					return nil
				}
			}
		}
	}
}

func (h *ProjectionHandler) fetchBulkStmts(
	ctx context.Context,
	query SearchQuery,
	reduce Reduce,
) (limitExeeded bool, err error) {
	eventQuery, eventsLimit, err := query()
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
		err = h.processEvent(ctx, event, reduce)
		logging.Log("HANDL-QRz8p").OnError(err).Warn("unable to reduce event")
	}

	return len(events) == int(eventsLimit), nil
}

func (h *ProjectionHandler) push(
	ctx context.Context,
	update Update,
	reduce Reduce,
) (err error) {
	h.lockMu.Lock()
	defer h.lockMu.Unlock()

	sort.Slice(h.stmts, func(i, j int) bool {
		return h.stmts[i].Sequence < h.stmts[j].Sequence
	})
	h.stmts, err = update(ctx, h.stmts, reduce)

	h.pushSet = len(h.stmts) > 0
	if h.pushSet {
		h.shouldPush <- nil
	}

	return err
}

func (h *ProjectionHandler) shutdown() {
	h.lockMu.Lock()
	defer h.lockMu.Unlock()
	h.Sub.Unsubscribe()
	h.Timer.Stop()
	close(h.shouldPush)
	logging.Log("EVENT-XG5Og").Info("stop processing")
}
