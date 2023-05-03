package spooler

import (
	"context"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	v1 "github.com/zitadel/zitadel/internal/eventstore/v1"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/eventstore/v1/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/view/repository"
)

const (
	systemID           = "system"
	schedulerSucceeded = eventstore.EventType("system.projections.scheduler.succeeded")
	aggregateType      = eventstore.AggregateType("system")
	aggregateID        = "SYSTEM"
)

type Spooler struct {
	handlers            []query.Handler
	locker              Locker
	lockID              string
	eventstore          v1.Eventstore
	esV2                *eventstore.Eventstore
	workers             int
	queue               chan *spooledHandler
	concurrentInstances int
}

type Locker interface {
	Renew(lockerID, viewModel, instanceID string, waitTime time.Duration) error
}

type spooledHandler struct {
	query.Handler
	locker              Locker
	queuedAt            time.Time
	eventstore          v1.Eventstore
	esV2                *eventstore.Eventstore
	concurrentInstances int
	succeededOnce       bool
}

func (s *Spooler) Start() {
	defer logging.WithFields("lockerID", s.lockID, "workers", s.workers).Info("spooler started")
	if s.workers < 1 {
		return
	}

	for i := 0; i < s.workers; i++ {
		go func(workerIdx int) {
			workerID := s.lockID + "--" + strconv.Itoa(workerIdx)
			for task := range s.queue {
				go requeueTask(task, s.queue)
				task.load(workerID)
			}
		}(i)
	}
	go func() {
		for _, handler := range s.handlers {
			s.queue <- &spooledHandler{Handler: handler, locker: s.locker, queuedAt: time.Now(), eventstore: s.eventstore, esV2: s.esV2, concurrentInstances: s.concurrentInstances}
		}
	}()
}

func requeueTask(task *spooledHandler, queue chan<- *spooledHandler) {
	time.Sleep(task.MinimumCycleDuration() - time.Since(task.queuedAt))
	task.queuedAt = time.Now()
	queue <- task
}

func (s *spooledHandler) hasSucceededOnce(ctx context.Context) (bool, error) {
	events, err := s.esV2.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(schedulerSucceeded).
		EventData(map[string]interface{}{
			"name": s.ViewModel(),
		}).
		Builder(),
	)
	return len(events) > 0 && err == nil, err
}

func (s *spooledHandler) setSucceededOnce(ctx context.Context) error {
	_, err := s.esV2.Push(ctx, &handler.ProjectionSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx,
			eventstore.NewAggregate(ctx, aggregateID, aggregateType, "v1"),
			schedulerSucceeded,
		),
		Name: s.ViewModel(),
	})
	s.succeededOnce = err == nil
	return err
}

func (s *spooledHandler) load(workerID string) {
	errs := make(chan error)
	defer func() {
		close(errs)
		err := recover()

		if err != nil {
			logging.WithFields(
				"cause", err,
				"stack", string(debug.Stack()),
			).Error("reduce panicked")
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	go s.awaitError(cancel, errs, workerID)
	hasLocked := s.lock(ctx, errs, workerID)

	if <-hasLocked {
		if !s.succeededOnce {
			var err error
			s.succeededOnce, err = s.hasSucceededOnce(ctx)
			if err != nil {
				logging.WithFields("view", s.ViewModel()).OnError(err).Debug("initial lock failed for first schedule")
				errs <- err
				return
			}
		}

		instanceIDQuery := models.NewSearchQuery().SetColumn(models.Columns_InstanceIDs).AddQuery().ExcludedInstanceIDsFilter("")
		for {
			if s.succeededOnce {
				// since we have at least one successful run, we can restrict it to events not older than
				// twice the requeue time (just to be sure not to miss an event)
				instanceIDQuery = instanceIDQuery.CreationDateNewerFilter(time.Now().Add(-2 * s.MinimumCycleDuration()))
			}
			ids, err := s.eventstore.InstanceIDs(ctx, instanceIDQuery.SearchQuery())
			if err != nil {
				errs <- err
				break
			}
			for i := 0; i < len(ids); i = i + s.concurrentInstances {
				max := i + s.concurrentInstances
				if max > len(ids) {
					max = len(ids)
				}
				err = s.processInstances(ctx, workerID, ids[i:max])
				if err != nil {
					errs <- err
				}
			}
			if ctx.Err() == nil {
				if !s.succeededOnce {
					err = s.setSucceededOnce(ctx)
					logging.WithFields("view", s.ViewModel()).OnError(err).Warn("unable to push first schedule succeeded")
				}
				errs <- nil
			}
			break
		}
	}
	<-ctx.Done()
}

func (s *spooledHandler) processInstances(ctx context.Context, workerID string, ids []string) error {
	for {
		processCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		events, err := s.query(processCtx, ids)
		if err != nil {
			cancel()
			return err
		}
		if len(events) == 0 {
			cancel()
			return nil
		}
		err = s.process(processCtx, events, workerID, ids)
		cancel()
		if err != nil {
			return err
		}
		if uint64(len(events)) < s.QueryLimit() {
			// no more events to process
			return nil
		}
	}
}

func (s *spooledHandler) awaitError(cancel func(), errs chan error, workerID string) {
	select {
	case err := <-errs:
		cancel()
		logging.OnError(err).WithField("view", s.ViewModel()).WithField("worker", workerID).Debug("load canceled")
	}
}

func (s *spooledHandler) process(ctx context.Context, events []*models.Event, workerID string, instanceIDs []string) error {
	for i, event := range events {
		select {
		case <-ctx.Done():
			logging.WithFields("view", s.ViewModel(), "worker", workerID, "traceID", tracing.TraceIDFromCtx(ctx)).Debug("context canceled")
			return nil
		default:
			if err := s.Reduce(event); err != nil {
				err = s.OnError(event, err)
				if err == nil {
					continue
				}
				time.Sleep(100 * time.Millisecond)
				return s.process(ctx, events[i:], workerID, instanceIDs)
			}
		}
	}
	err := s.OnSuccess(instanceIDs)
	logging.WithFields("view", s.ViewModel(), "worker", workerID, "traceID", tracing.TraceIDFromCtx(ctx)).OnError(err).Warn("could not process on success func")
	return err
}

func (s *spooledHandler) query(ctx context.Context, instanceIDs []string) ([]*models.Event, error) {
	query, err := s.EventQuery(ctx, instanceIDs)
	if err != nil {
		return nil, err
	}
	query.Limit = s.QueryLimit()
	return s.eventstore.FilterEvents(ctx, query)
}

// lock ensures the lock on the database.
// the returned channel will be closed if ctx is done or an error occured durring lock
func (s *spooledHandler) lock(ctx context.Context, errs chan<- error, workerID string) chan bool {
	renewTimer := time.After(0)
	locked := make(chan bool)

	go func(locked chan bool) {
		var firstLock sync.Once
		defer close(locked)
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				err := s.locker.Renew(workerID, s.ViewModel(), systemID, s.LockDuration())
				firstLock.Do(func() {
					locked <- err == nil
				})
				if err == nil {
					renewTimer = time.After(s.LockDuration())
					continue
				}

				if ctx.Err() == nil {
					errs <- err
				}
				return
			}
		}
	}(locked)

	return locked
}

func HandleError(event *models.Event, failedErr error,
	latestFailedEvent func(sequence uint64, instanceID string) (*repository.FailedEvent, error),
	processFailedEvent func(*repository.FailedEvent) error,
	processSequence func(*models.Event) error,
	errorCountUntilSkip uint64) error {
	failedEvent, err := latestFailedEvent(event.Sequence, event.InstanceID)
	if err != nil {
		return err
	}
	failedEvent.FailureCount++
	failedEvent.ErrMsg = failedErr.Error()
	failedEvent.InstanceID = event.InstanceID
	failedEvent.LastFailed = time.Now()
	err = processFailedEvent(failedEvent)
	if err != nil {
		return err
	}
	if errorCountUntilSkip <= failedEvent.FailureCount {
		return processSequence(event)
	}
	return failedErr
}

func HandleSuccess(updateSpoolerRunTimestamp func([]string) error, instanceIDs []string) error {
	return updateSpoolerRunTimestamp(instanceIDs)
}
