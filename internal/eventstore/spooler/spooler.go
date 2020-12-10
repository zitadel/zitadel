package spooler

import (
	"context"
	"strconv"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/view/repository"

	"time"
)

const (
	defaultLockTime = time.Second * 30
)

type Spooler struct {
	handlers   []query.Handler
	locker     Locker
	lockID     string
	eventstore eventstore.Eventstore
	workers    int
}

type Locker interface {
	Renew(lockerID, viewModel string, waitTime time.Duration) (time.Time, bool, error)
}

type worker struct {
	locker     Locker
	lockID     string
	eventstore eventstore.Eventstore
	taskQueue  chan *task
}

func (w *worker) start() {
	for t := range w.taskQueue {
		lockedUntil, isLeaseHolder, err := w.lock(t, defaultLockTime)
		if err != nil {
			go w.queueTask(t, time.Now().Add(t.MinimumCycleDuration()))
			continue
		} else if !isLeaseHolder {
			go w.queueTask(t, lockedUntil)
			continue
		}
		ctx, cancel := context.WithCancel(context.Background())
		errs := make(chan error)
		defer close(errs)

		go t.execute(ctx, w.eventstore, errs)
		for {
			renewTimer := time.After(lockedUntil.Sub(time.Now().Add(5 * time.Second)))
			select {
			case <-errs:
				cancel()
				go w.queueTask(t, time.Now().Add(defaultLockTime))
				return
			case <-renewTimer:
				lockedUntil, err = w.renewLease(lockedUntil, t)
				if err != nil {
					cancel()
					go w.queueTask(t, time.Now().Add(defaultLockTime))
					return
				}
			}
		}
	}
}

func (w *worker) renewLease(lockedUntil time.Time, task *task) (time.Time, error) {
	lockedUntil, isLeaseHolder, err := w.lock(task, defaultLockTime)
	if err != nil || !isLeaseHolder {
		return time.Time{}, errors.ThrowInternal(err, "SPOOL-pRLxA", "error in renew lease")
	}
	return lockedUntil, nil
}

func (w *worker) lock(task *task, lockDuration time.Duration) (time.Time, bool, error) {
	return w.locker.Renew(w.lockID, task.ViewModel(), lockDuration)
}

func (w *worker) queueTask(task *task, lockedUntil time.Time) {
	time.Sleep(lockedUntil.Sub(time.Now()))
	logging.Log("SPOOL-VQZko").Debug("requeue task")
	w.taskQueue <- task
}

type task struct {
	query.Handler
}

func (t *task) execute(ctx context.Context, es eventstore.Eventstore, errs chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := t.fillView(ctx, es); err != nil {
				logging.Log("SPOOL-RuKm3").WithError(err).Info("fill view failed")
				errs <- err
				return
			}
		}
	}
}

func (t *task) fillView(ctx context.Context, es eventstore.Eventstore) error {
	minDuration := time.After(t.MinimumCycleDuration())

	events, err := t.query(ctx, es)
	if err != nil || len(events) == 0 {
		return err
	}
	if err = t.process(ctx, events); err != nil {
		return err
	}

	<-minDuration
	return nil
}

func (t *task) query(ctx context.Context, es eventstore.Eventstore) ([]*models.Event, error) {
	query, err := t.Handler.EventQuery()
	if err != nil {
		return nil, err
	}
	factory := models.FactoryFromSearchQuery(query)
	sequence, err := es.LatestSequence(ctx, factory)
	logging.Log("SPOOL-7SciK").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Debug("unable to query latest sequence")
	var processedSequence uint64
	for _, filter := range query.Filters {
		if filter.GetField() == models.Field_LatestSequence {
			processedSequence = filter.GetValue().(uint64)
		}
	}
	if sequence == 0 || processedSequence >= sequence {
		return nil, nil
	}

	query.Limit = t.QueryLimit()

	return es.FilterEvents(ctx, query)
}

func (t *task) process(ctx context.Context, events []*models.Event) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := t.Handler.Reduce(event); err != nil {
				return t.OnError(event, err)
			}
		}
	}
	return nil
}

func (s *Spooler) Start() {
	defer logging.LogWithFields("SPOOL-N0V1g", "lockerID", s.lockID, "workers", s.workers).Info("spooler started")
	if s.workers < 1 {
		return
	}

	tasks := make(chan *task)

	for i := 0; i < s.workers; i++ {
		w := &worker{locker: s.locker, lockID: s.lockID + "--" + strconv.Itoa(i), eventstore: s.eventstore, taskQueue: tasks}
		go w.start()
	}

	go func() {
		for _, handler := range s.handlers {
			tasks <- &task{Handler: handler}
		}
	}()
}

func HandleError(
	event *models.Event,
	failedErr error,
	latestFailedEvent func(sequence uint64) (*repository.FailedEvent, error),
	processFailedEvent func(*repository.FailedEvent) error,
	processSequence func(uint64) error,
	errorCountUntilSkip uint64,
) error {
	failedEvent, err := latestFailedEvent(event.Sequence)
	if err != nil {
		return err
	}
	failedEvent.FailureCount++
	failedEvent.ErrMsg = failedErr.Error()
	err = processFailedEvent(failedEvent)
	if err != nil {
		return err
	}
	if errorCountUntilSkip <= failedEvent.FailureCount {
		return processSequence(event.Sequence, event.CreationDate)
	}
	return nil
}

func HandleSuccess(updateSpoolerRunTimestamp func() error) error {
	return updateSpoolerRunTimestamp()
}
