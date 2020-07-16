package spooler

import (
	"context"
	"strconv"
	"strings"
	"sync"

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

	isWorking bool
	working   sync.Mutex
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

				if handler.isWorking {
					continue
				}
				handler.working.Lock()
				handler.isWorking = true
				handler.working.Unlock()

				handler.lockID = strings.Split(handler.lockID, "--")[0] + "--" + strconv.Itoa(taskIdx)
				handler.load()

				handler.working.Lock()
				handler.isWorking = false
				handler.working.Unlock()
			}
		}(i)
	}
	for _, handler := range s.handlers {
		handler := &spooledHandler{Handler: handler, locker: s.locker, lockID: s.lockID, queuedAt: time.Now(), eventstore: s.eventstore}
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
					// we only need to break. an error is already written by the lock-routine to the errs channel
					break
				}
			}
		}()
		events, err := s.query(ctx)
		if err != nil {
			errs <- err
		} else {
			logging.LogWithFields("SPOOL-aqLrD", "eventCount", len(events), "lock", s.lockID, "view", s.ViewModel()).Debug("will load")
			errs <- s.process(ctx, events)
		}
	}
	<-ctx.Done()
}

func (s *spooledHandler) awaitError(cancel func(), errs chan error) {
	select {
	case err := <-errs:
		cancel()
		logging.Log("SPOOL-K2lst").OnError(err).WithField("view", s.ViewModel()).WithField("lock", s.lockID).Debug("load canceled")
	}
}

func (s *spooledHandler) process(ctx context.Context, events []*models.Event) error {
	for _, event := range events {
		select {
		case <-ctx.Done():
			logging.Log("SPOOL-FTKwH").WithField("view", s.ViewModel()).WithField("lock", s.lockID).Debug("context canceled")
			return nil
		default:
			logging.LogWithFields("SPOOL-Q0YGW", "seq", event.Sequence, "lock", s.lockID, "view", s.ViewModel()).Debug("reduce")
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
	locked := make(chan bool)

	go func(locked chan bool) {
		for {
			select {
			case <-ctx.Done():
				// logging.LogWithFields("SPOOL-ymvYF", "lock", s.lockID, "view", s.ViewModel()).Debug("lock stoped")
				return
			case <-renewTimer:
				// logging.LogWithFields("SPOOL-kwuFD", "lock", s.lockID, "view", s.ViewModel()).Debug("will lock")
				err := s.locker.Renew(s.lockID, s.ViewModel(), s.MinimumCycleDuration()*2)
				// logging.LogWithFields("SPOOL-wOURR", "lock", s.lockID, "view", s.ViewModel()).WithError(err).Debug("locked")
				if err == nil {
					// logging.LogWithFields("SPOOL-FxNmG", "lock", s.lockID, "view", s.ViewModel()).WithError(err).Debug("will write to chan")
					locked <- true
					// logging.LogWithFields("SPOOL-Acr90", "lock", s.lockID, "view", s.ViewModel(), "renewDuration", renewDuration).WithError(err).Debug("written to chan")
					renewTimer = time.After(renewDuration)
					continue
				}

				if ctx.Err() == nil {
					errs <- err
				}

				// logging.LogWithFields("SPOOL-3GH43", "lock", s.lockID, "view", s.ViewModel()).WithError(err).Debug("will write false to chan")
				locked <- false
				// logging.LogWithFields("SPOOL-jq6dm", "lock", s.lockID, "view", s.ViewModel()).WithError(err).Debug("written false to chan")
				return
			}
		}
	}(locked)

	return locked
}
