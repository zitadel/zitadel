package handler

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
)

type Projection struct {
	Handler
	ProjectionName      string
	reduce              Reduce
	update              Update
	searchQuery         SearchQuery
	shouldBulk          *time.Timer
	lock                Lock
	unlock              Unlock
	requeueAfter        time.Duration
	retryFailedAfter    time.Duration
	retries             int
	concurrentInstances int
}

func NewProjection(config ProjectionHandlerConfig, reduce Reduce, update Update, searchQuery SearchQuery, lock Lock, unlock Unlock) *Projection {
	p := &Projection{
		Handler:        NewHandler(config.HandlerConfig),
		ProjectionName: config.ProjectionName,
		reduce:         reduce,
		update:         update,
		searchQuery:    searchQuery,
		lock:           lock,
		unlock:         unlock,
		requeueAfter:   config.RequeueEvery,
		//// first bulk is instant on startup
		shouldBulk: time.NewTimer(0),
		//shouldPush:       time.NewTimer(0),
		retryFailedAfter:    config.RetryFailedAfter,
		retries:             int(config.Retries),
		concurrentInstances: int(config.ConcurrentInstances),
	}

	go p.subscribe()

	go p.schedule()

	return p
}

func (p *Projection) subscribe() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		err := recover()
		if err != nil {
			p.Handler.Sub.Unsubscribe()
			logging.WithFields("projection", p.ProjectionName).Errorf("subscription panicked: %v", err)
		}
		cancel()
	}()
	for firstEvent := range p.EventQueue {
		events := checkAdditionalEvents(p.EventQueue, firstEvent)

		index, err := p.Process(ctx, events...)
		if err != nil || index < len(events)-1 {
			logging.WithFields("projection", p.ProjectionName).WithError(err).Error("unable to process all events from subscription")
		}
	}
}

func (p *Projection) schedule() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		err := recover()
		if err != nil {
			logging.WithFields("projection", p.ProjectionName).Errorf("schedule panicked: %v", err)
		}
		cancel()
	}()
	for {
		select {
		case <-p.shouldBulk.C:
			ids, err := p.Eventstore.InstanceIDs(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).AddQuery().ExcludedInstanceID("").Builder())
			if err != nil {
				logging.WithFields("projection", p.ProjectionName).WithError(err).Error("instance ids")
				p.shouldBulk.Reset(p.requeueAfter)
				continue
			}
			for i := 0; i < len(ids); i = i + p.concurrentInstances {
				max := i + p.concurrentInstances
				if max > len(ids) {
					max = len(ids)
				}
				instances := ids[i:max]
				ctx, cancel := context.WithCancel(ctx)
				errs := p.lock(ctx, p.requeueAfter, instances...)
				//wait until projection is locked
				if err, ok := <-errs; err != nil || !ok {
					cancel()
					logging.WithFields("projection", p.ProjectionName).OnError(err).Warn("initial lock failed")
					continue
				}
				go p.cancelOnErr(ctx, errs, cancel)
				err = p.Trigger(ctx, instances...)
				if err != nil {
					logging.WithFields("projection", p.ProjectionName, "instanceIDs", instances).WithError(err).Error("trigger failed")
				}

				cancel()
				unlockErr := p.unlock(instances...)
				logging.WithFields("projection", p.ProjectionName).OnError(unlockErr).Warn("unable to unlock")
			}
			p.shouldBulk.Reset(p.requeueAfter)
		}
	}
}

func (h *Projection) cancelOnErr(ctx context.Context, errs <-chan error, cancel func()) {
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

func (p *Projection) Trigger(ctx context.Context, instances ...string) error {
	ids := []string{authz.GetInstance(ctx).InstanceID()}
	if len(instances) > 0 {
		ids = instances
	}
	for {
		events, hasLimitExeeded, err := p.FetchEvents(ctx, ids...)
		if err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}
		_, err = p.Process(ctx, events...)
		if err != nil {
			return err
		}
		if !hasLimitExeeded {
			return nil
		}
	}
}

//Process handles multiple events by reducing them to statements and updating the projection
func (p *Projection) Process(ctx context.Context, events ...eventstore.Event) (index int, err error) {
	statements := make([]*Statement, len(events))
	for i, event := range events {
		statements[i], err = p.reduce(event)
		if err != nil {
			return 0, err
		}
	}
	index = -1
	for retry := 0; retry <= p.retries; retry++ {
		index, err = p.update(ctx, statements[index+1:], p.reduce)
		if err != nil && !errors.Is(err, ErrSomeStmtsFailed) {
			return index, err
		}
		if err == nil {
			return index, nil
		}
		time.Sleep(p.retryFailedAfter)
	}
	return index, err
}

//FetchEvents checks the current sequences and filters for newer events
func (p *Projection) FetchEvents(ctx context.Context, instances ...string) ([]eventstore.Event, bool, error) {
	eventQuery, eventsLimit, err := p.searchQuery(ctx, instances)
	if err != nil {
		return nil, false, err
	}
	events, err := p.Eventstore.Filter(ctx, eventQuery)
	if err != nil {
		return nil, false, err
	}
	return events, int(eventsLimit) == len(events), err
}
