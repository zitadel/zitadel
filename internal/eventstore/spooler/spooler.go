package spooler

import (
	"context"
	"strconv"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/view/repository"

	"time"
)

type Spooler struct {
	handlers        []Handler
	locker          Locker
	lockID          string
	eventstore      eventstore.Eventstore
	concurrentTasks int
	queue           chan *spooledHandler
}

type Handler interface {
	query.Handler
	MinimumCycleDuration() time.Duration
}

type Locker interface {
	Renew(lockerID, viewModel string, waitTime time.Duration) error
}

type spooledHandler struct {
	Handler
	locker     Locker
	lockID     string
	queuedAt   time.Time
	eventstore eventstore.Eventstore
}

func (s *Spooler) Start() {
	defer logging.LogWithFields("SPOOL-N0V1g", "lockerID", s.lockID, "workers", s.concurrentTasks).Info("spooler started")
	if s.concurrentTasks < 1 {
		return
	}

	for i := 0; i < s.concurrentTasks; i++ {
		go func(taskIdx int) {
			for handler := range s.queue {
				go func(handler *spooledHandler, queue chan<- *spooledHandler) {
					time.Sleep(handler.MinimumCycleDuration() - time.Since(handler.queuedAt))
					handler.queuedAt = time.Now()
					queue <- handler
				}(handler, s.queue)

				handler.lockID = strings.Split(handler.lockID, "--")[0] + "--" + strconv.Itoa(taskIdx)
				handler.load()
			}
		}(i)
	}
	for _, handler := range s.handlers {
		handler := &spooledHandler{handler, s.locker, s.lockID, time.Now(), s.eventstore}
		s.queue <- handler
	}
}

func (s *spooledHandler) load() {
	errs := make(chan error)
	defer close(errs)
	ctx, cancel := context.WithCancel(context.Background())
	go s.awaitError(cancel, errs)
	hasLocked := s.lock(ctx, errs)

	if <-hasLocked {
		go func() {
			for l := range hasLocked {
				if !l {
					ctx.Done()
				}
			}
		}()
		events, err := s.query(ctx)
		if err != nil {
			errs <- err
		} else {
			logging.LogWithFields("SPOOL-aqLrD", "eventCount", len(events), "lock", s.lockID, "ts", time.Now().Format("15-01-2018 15:04:05.000000"), "view", s.ViewModel()).Debug("i will load")
			errs <- s.process(ctx, events)
		}
	}
	<-ctx.Done()
}

func (s *spooledHandler) awaitError(cancel func(), errs chan error) {
	select {
	case err := <-errs:
		cancel()
		logging.Log("SPOOL-K2lst").OnError(err).WithField("view", s.ViewModel()).Debug("load canceled")
	}
}

func (s *spooledHandler) process(ctx context.Context, events []*models.Event) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			logging.Log("SPOOL-FTKwH").WithField("view", s.ViewModel()).Debug("context canceled")
			return nil
		default:
			if err := s.Reduce(event); err != nil {
				return s.OnError(event, err)
			}
		}
	}
	return nil
}

func HandleError(event *models.Event, failedErr error,
	latestFailedEvent func(sequence uint64) (*repository.FailedEvent, error),
	processFailedEvent func(*repository.FailedEvent) error,
	processSequence func(uint64) error, errorCountUntilSkip uint64) error {
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
		return processSequence(event.Sequence)
	}
	return nil
}

func (s *spooledHandler) query(ctx context.Context) ([]*models.Event, error) {
	query, err := s.EventQuery()
	if err != nil {
		return nil, err
	}
	return s.eventstore.FilterEvents(ctx, query)
}

func (s *spooledHandler) lock(ctx context.Context, errs chan<- error) chan bool {
	renewTimer := time.After(0)
	renewDuration := s.MinimumCycleDuration() - 50*time.Millisecond
	locked := make(chan bool, 1)

	go func(locked chan bool) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-renewTimer:
				err := s.locker.Renew(s.lockID, s.ViewModel(), s.MinimumCycleDuration()*2)
				logging.LogWithFields("SPOOL-C6FcM", "lock", s.lockID, "ts", time.Now().Format("15-01-2018 15:04:05.000000"), "view", s.ViewModel()).Debug("renewed")
				if err == nil {
					locked <- true
					logging.LogWithFields("SPOOL-NqFkJ", "lock", s.lockID, "ts", time.Now().Format("15-01-2018 15:04:05.000000"), "view", s.ViewModel()).Debug("wont reach")
					renewTimer = time.After(renewDuration)
					continue
				}

				if ctx.Err() == nil {
					errs <- err
				}

				locked <- false
				return
			}
		}
	}(locked)

	return locked
}
