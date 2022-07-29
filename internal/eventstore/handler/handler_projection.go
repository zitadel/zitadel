package handler

import (
	"context"
	"errors"
	"runtime/debug"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const systemID = "system"

type ProjectionHandlerConfig struct {
	HandlerConfig
	ProjectionName      string
	RequeueEvery        time.Duration
	RetryFailedAfter    time.Duration
	Retries             uint
	ConcurrentInstances uint
}

//Update updates the projection with the given statements
type Update func(context.Context, []*Statement, Reduce) (index int, err error)

//Reduce reduces the given event to a statement
//which is used to update the projection
type Reduce func(eventstore.Event) (*Statement, error)

//SearchQuery generates the search query to lookup for events
type SearchQuery func(ctx context.Context, instanceIDs []string) (query *eventstore.SearchQueryBuilder, queryLimit uint64, err error)

//Lock is used for mutex handling if needed on the projection
type Lock func(context.Context, time.Duration, ...string) <-chan error

//Unlock releases the mutex of the projection
type Unlock func(...string) error

type ProjectionHandler struct {
	Handler
	ProjectionName      string
	reduce              Reduce
	update              Update
	searchQuery         SearchQuery
	triggerProjection   *time.Timer
	lock                Lock
	unlock              Unlock
	requeueAfter        time.Duration
	retryFailedAfter    time.Duration
	retries             int
	concurrentInstances int
}

func NewProjectionHandler(
	ctx context.Context,
	config ProjectionHandlerConfig,
	reduce Reduce,
	update Update,
	query SearchQuery,
	lock Lock,
	unlock Unlock,
) *ProjectionHandler {
	concurrentInstances := int(config.ConcurrentInstances)
	if concurrentInstances < 1 {
		concurrentInstances = 1
	}
	h := &ProjectionHandler{
		Handler:             NewHandler(config.HandlerConfig),
		ProjectionName:      config.ProjectionName,
		reduce:              reduce,
		update:              update,
		searchQuery:         query,
		lock:                lock,
		unlock:              unlock,
		requeueAfter:        config.RequeueEvery,
		triggerProjection:   time.NewTimer(0), // first trigger is instant on startup
		retryFailedAfter:    config.RetryFailedAfter,
		retries:             int(config.Retries),
		concurrentInstances: concurrentInstances,
	}

	go h.subscribe(ctx)

	go h.schedule(ctx)

	return h
}

//Trigger handles all events for the provided instances (or current instance from context if non specified)
//by calling FetchEvents and Process until the amount of events is smaller than the BulkLimit
func (h *ProjectionHandler) Trigger(ctx context.Context, instances ...string) error {
	ids := []string{authz.GetInstance(ctx).InstanceID()}
	if len(instances) > 0 {
		ids = instances
	}
	for {
		events, hasLimitExceeded, err := h.FetchEvents(ctx, ids...)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}
		_, err = h.Process(ctx, events...)
		if err != nil {
			return err
		}
		if !hasLimitExceeded {
			return nil
		}
	}
}

//Process handles multiple events by reducing them to statements and updating the projection
func (h *ProjectionHandler) Process(ctx context.Context, events ...eventstore.Event) (index int, err error) {
	if len(events) == 0 {
		return 0, nil
	}
	index = -1
	statements := make([]*Statement, len(events))
	for i, event := range events {
		statements[i], err = h.reduce(event)
		if err != nil {
			return index, err
		}
	}
	for retry := 0; retry <= h.retries; retry++ {
		index, err = h.update(ctx, statements[index+1:], h.reduce)
		if err != nil && !errors.Is(err, ErrSomeStmtsFailed) {
			return index, err
		}
		if err == nil {
			return index, nil
		}
		time.Sleep(h.retryFailedAfter)
	}
	return index, err
}

//FetchEvents checks the current sequences and filters for newer events
func (h *ProjectionHandler) FetchEvents(ctx context.Context, instances ...string) ([]eventstore.Event, bool, error) {
	eventQuery, eventsLimit, err := h.searchQuery(ctx, instances)
	if err != nil {
		return nil, false, err
	}
	events, err := h.Eventstore.Filter(ctx, eventQuery)
	if err != nil {
		return nil, false, err
	}
	return events, int(eventsLimit) == len(events), err
}

func (h *ProjectionHandler) subscribe(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		err := recover()
		if err != nil {
			h.Handler.Unsubscribe()
			logging.WithFields("projection", h.ProjectionName).Errorf("subscription panicked: %v", err)
		}
		cancel()
	}()
	for firstEvent := range h.EventQueue {
		events := checkAdditionalEvents(h.EventQueue, firstEvent)

		index, err := h.Process(ctx, events...)
		if err != nil || index < len(events)-1 {
			logging.WithFields("projection", h.ProjectionName).WithError(err).Warn("unable to process all events from subscription")
		}
	}
}

func (h *ProjectionHandler) schedule(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		err := recover()
		if err != nil {
			logging.WithFields("projection", h.ProjectionName, "cause", err, "stack", string(debug.Stack())).Error("schedule panicked")
		}
		cancel()
	}()
	for range h.triggerProjection.C {
		ids, err := h.Eventstore.InstanceIDs(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).AddQuery().ExcludedInstanceID("").Builder())
		if err != nil {
			logging.WithFields("projection", h.ProjectionName).WithError(err).Error("instance ids")
			h.triggerProjection.Reset(h.requeueAfter)
			continue
		}
		for i := 0; i < len(ids); i = i + h.concurrentInstances {
			max := i + h.concurrentInstances
			if max > len(ids) {
				max = len(ids)
			}
			instances := ids[i:max]
			lockCtx, cancelLock := context.WithCancel(ctx)
			errs := h.lock(lockCtx, h.requeueAfter, instances...)
			//wait until projection is locked
			if err, ok := <-errs; err != nil || !ok {
				cancelLock()
				logging.WithFields("projection", h.ProjectionName).OnError(err).Warn("initial lock failed")
				continue
			}
			go h.cancelOnErr(lockCtx, errs, cancelLock)
			err = h.Trigger(lockCtx, instances...)
			if err != nil {
				logging.WithFields("projection", h.ProjectionName, "instanceIDs", instances).WithError(err).Error("trigger failed")
			}

			cancelLock()
			unlockErr := h.unlock(instances...)
			logging.WithFields("projection", h.ProjectionName).OnError(unlockErr).Warn("unable to unlock")
		}
		h.triggerProjection.Reset(h.requeueAfter)
	}
}

func (h *ProjectionHandler) cancelOnErr(ctx context.Context, errs <-chan error, cancel func()) {
	for {
		select {
		case err := <-errs:
			if err != nil {
				logging.WithFields("projection", h.ProjectionName).WithError(err).Warn("bulk canceled")
				cancel()
				return
			}
		case <-ctx.Done():
			cancel()
			return
		}

	}
}

func checkAdditionalEvents(eventQueue chan eventstore.Event, event eventstore.Event) []eventstore.Event {
	events := make([]eventstore.Event, 1)
	events[0] = event
	for {
		select {
		case event := <-eventQueue:
			events = append(events, event)
		default:
			return events
		}
	}
}
