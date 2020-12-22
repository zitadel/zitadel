package spooler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/query"
	"github.com/caos/zitadel/internal/eventstore/spooler/mock"
	"github.com/caos/zitadel/internal/view/repository"
	"github.com/golang/mock/gomock"
)

type testHandler struct {
	cycleDuration time.Duration
	processSleep  time.Duration
	processError  error
	queryError    error
	viewModel     string
	bulkLimit     uint64
	maxErrCount   int
}

func (h *testHandler) AggregateTypes() []models.AggregateType {
	return nil
}

func (h *testHandler) CurrentSequence(event *models.Event) (uint64, error) {
	return 0, nil
}

func (h *testHandler) Eventstore() eventstore.Eventstore {
	return nil
}

func (h *testHandler) ViewModel() string {
	return h.viewModel
}

func (h *testHandler) EventQuery() (*models.SearchQuery, error) {
	if h.queryError != nil {
		return nil, h.queryError
	}
	return &models.SearchQuery{}, nil
}

func (h *testHandler) Reduce(*models.Event) error {
	<-time.After(h.processSleep)
	return h.processError
}

func (h *testHandler) OnError(event *models.Event, err error) error {
	if h.maxErrCount == 2 {
		return nil
	}
	h.maxErrCount++
	return err
}

func (h *testHandler) OnSuccess() error {
	return nil
}

func (h *testHandler) MinimumCycleDuration() time.Duration {
	return h.cycleDuration
}

func (h *testHandler) LockDuration() time.Duration {
	return h.cycleDuration / 2
}

func (h *testHandler) QueryLimit() uint64 {
	return h.bulkLimit
}

type eventstoreStub struct {
	events []*models.Event
	err    error
}

func (es *eventstoreStub) Subscribe(...models.AggregateType) *eventstore.Subscription { return nil }

func (es *eventstoreStub) Health(ctx context.Context) error {
	return nil
}

func (es *eventstoreStub) AggregateCreator() *models.AggregateCreator {
	return nil
}

func (es *eventstoreStub) FilterEvents(ctx context.Context, in *models.SearchQuery) ([]*models.Event, error) {
	if es.err != nil {
		return nil, es.err
	}
	return es.events, nil
}
func (es *eventstoreStub) PushAggregates(ctx context.Context, in ...*models.Aggregate) error {
	return nil
}

func (es *eventstoreStub) LatestSequence(ctx context.Context, in *models.SearchQueryFactory) (uint64, error) {
	return 0, nil
}

func TestSpooler_process(t *testing.T) {
	type fields struct {
		currentHandler *testHandler
	}
	type args struct {
		timeout time.Duration
		events  []*models.Event
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantRetries int
	}{
		{
			name: "process all events",
			fields: fields{
				currentHandler: &testHandler{},
			},
			args: args{
				timeout: 0,
				events:  []*models.Event{{}, {}},
			},
			wantErr: false,
		},
		{
			name: "deadline exeeded",
			fields: fields{
				currentHandler: &testHandler{processSleep: 501 * time.Millisecond},
			},
			args: args{
				timeout: 1 * time.Second,
				events:  []*models.Event{{}, {}, {}, {}},
			},
			wantErr: false,
		},
		{
			name: "process error",
			fields: fields{
				currentHandler: &testHandler{processSleep: 1 * time.Second, processError: fmt.Errorf("i am an error")},
			},
			args: args{
				events: []*models.Event{{}, {}},
			},
			wantErr:     false,
			wantRetries: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &spooledHandler{
				Handler: tt.fields.currentHandler,
			}

			ctx := context.Background()
			var start time.Time
			if tt.args.timeout > 0 {
				ctx, _ = context.WithTimeout(ctx, tt.args.timeout)
				start = time.Now()
			}

			if err := s.process(ctx, tt.args.events, "test"); (err != nil) != tt.wantErr {
				t.Errorf("Spooler.process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.currentHandler.maxErrCount != tt.wantRetries {
				t.Errorf("Spooler.process() wrong retry count got: %d want %d", tt.fields.currentHandler.maxErrCount, tt.wantRetries)
			}

			elapsed := time.Since(start).Round(1 * time.Second)
			if tt.args.timeout != 0 && elapsed != tt.args.timeout {
				t.Errorf("wrong timeout wanted %v elapsed %v since %v", tt.args.timeout, elapsed, time.Since(start))
			}
		})
	}
}

func TestSpooler_awaitError(t *testing.T) {
	type fields struct {
		currentHandler query.Handler
		err            error
		canceled       bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"no error",
			fields{
				err:            nil,
				currentHandler: &testHandler{processSleep: 500 * time.Millisecond},
				canceled:       false,
			},
		},
		{
			"with error",
			fields{
				err:            fmt.Errorf("hodor"),
				currentHandler: &testHandler{processSleep: 500 * time.Millisecond},
				canceled:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &spooledHandler{
				Handler: tt.fields.currentHandler,
			}
			errs := make(chan error)
			ctx, cancel := context.WithCancel(context.Background())

			go s.awaitError(cancel, errs, "test")
			errs <- tt.fields.err

			if ctx.Err() == nil {
				t.Error("cancel function was not called")
			}
		})
	}
}

// TestSpooler_load checks if load terminates
func TestSpooler_load(t *testing.T) {
	type fields struct {
		currentHandler query.Handler
		locker         *testLocker
		eventstore     eventstore.Eventstore
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"lock exists",
			fields{
				currentHandler: &testHandler{processSleep: 500 * time.Millisecond, viewModel: "testView1", cycleDuration: 1 * time.Second, bulkLimit: 10},
				locker:         newTestLocker(t, "testID", "testView1").expectRenew(t, fmt.Errorf("lock already exists"), 500*time.Millisecond),
			},
		},
		{
			"lock fails",
			fields{
				currentHandler: &testHandler{processSleep: 100 * time.Millisecond, viewModel: "testView2", cycleDuration: 1 * time.Second, bulkLimit: 10},
				locker:         newTestLocker(t, "testID", "testView2").expectRenew(t, fmt.Errorf("fail"), 500*time.Millisecond),
				eventstore:     &eventstoreStub{events: []*models.Event{{}}},
			},
		},
		{
			"query fails",
			fields{
				currentHandler: &testHandler{processSleep: 100 * time.Millisecond, viewModel: "testView3", queryError: fmt.Errorf("query fail"), cycleDuration: 1 * time.Second, bulkLimit: 10},
				locker:         newTestLocker(t, "testID", "testView3").expectRenew(t, nil, 500*time.Millisecond),
				eventstore:     &eventstoreStub{err: fmt.Errorf("fail")},
			},
		},
		{
			"process event fails",
			fields{
				currentHandler: &testHandler{processError: fmt.Errorf("oups"), processSleep: 100 * time.Millisecond, viewModel: "testView4", cycleDuration: 500 * time.Millisecond, bulkLimit: 10},
				locker: newTestLocker(t, "testID", "testView4").
					expectRenew(t, nil, 250*time.Millisecond).
					expectRenew(t, nil, 250*time.Millisecond).
					expectRenew(t, nil, 250*time.Millisecond),
				eventstore: &eventstoreStub{events: []*models.Event{{}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.fields.locker.finish()
			s := &spooledHandler{
				Handler:    tt.fields.currentHandler,
				locker:     tt.fields.locker.mock,
				eventstore: tt.fields.eventstore,
			}
			s.load("test-worker")
		})
	}
}

func TestSpooler_lock(t *testing.T) {
	type fields struct {
		currentHandler query.Handler
		locker         *testLocker
		expectsErr     bool
	}
	type args struct {
		deadline time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"renew correct",
			fields{
				currentHandler: &testHandler{cycleDuration: 1 * time.Second, viewModel: "testView"},
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, nil, 500*time.Millisecond),
				expectsErr:     false,
			},
			args{
				deadline: time.Now().Add(1 * time.Second),
			},
		},
		{
			"renew fails",
			fields{
				currentHandler: &testHandler{cycleDuration: 900 * time.Millisecond, viewModel: "testView"},
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, fmt.Errorf("renew failed"), 450*time.Millisecond),
				expectsErr:     true,
			},
			args{
				deadline: time.Now().Add(5 * time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.fields.locker.finish()
			s := &spooledHandler{
				Handler: tt.fields.currentHandler,
				locker:  tt.fields.locker.mock,
			}

			errs := make(chan error, 1)
			defer close(errs)
			ctx, _ := context.WithDeadline(context.Background(), tt.args.deadline)

			locked := s.lock(ctx, errs, "test-worker")

			if tt.fields.expectsErr {
				lock := <-locked
				err := <-errs
				if err == nil {
					t.Error("No error in error queue")
				}
				if lock {
					t.Error("lock should have failed")
				}
			} else {
				lock := <-locked
				if !lock {
					t.Error("lock should be true")
				}
			}
		})
	}
}

type testLocker struct {
	mock     *mock.MockLocker
	lockerID string
	viewName string
	ctrl     *gomock.Controller
}

func newTestLocker(t *testing.T, lockerID, viewName string) *testLocker {
	ctrl := gomock.NewController(t)
	return &testLocker{mock.NewMockLocker(ctrl), lockerID, viewName, ctrl}
}

func (l *testLocker) expectRenew(t *testing.T, err error, waitTime time.Duration) *testLocker {
	t.Helper()
	l.mock.EXPECT().Renew(gomock.Any(), l.viewName, gomock.Any()).DoAndReturn(
		func(_, _ string, gotten time.Duration) error {
			t.Helper()
			if waitTime-gotten != 0 {
				t.Errorf("expected waittime %v got %v", waitTime, gotten)
			}
			return err
		}).Times(1)

	return l
}

func (l *testLocker) finish() {
	l.ctrl.Finish()
}

func TestHandleError(t *testing.T) {
	type args struct {
		event               *models.Event
		failedErr           error
		latestFailedEvent   func(sequence uint64) (*repository.FailedEvent, error)
		errorCountUntilSkip uint64
	}
	type res struct {
		wantErr               bool
		shouldProcessSequence bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "should process sequence already too high",
			args: args{
				event:     &models.Event{Sequence: 30000000},
				failedErr: errors.ThrowInternal(nil, "SPOOL-Wk53B", "this was wrong"),
				latestFailedEvent: func(s uint64) (*repository.FailedEvent, error) {
					return &repository.FailedEvent{
						ErrMsg:         "blub",
						FailedSequence: s - 1,
						FailureCount:   6,
						ViewName:       "super.table",
					}, nil
				},
				errorCountUntilSkip: 5,
			},
			res: res{
				shouldProcessSequence: true,
			},
		},
		{
			name: "should process sequence after this event too high",
			args: args{
				event:     &models.Event{Sequence: 30000000},
				failedErr: errors.ThrowInternal(nil, "SPOOL-Wk53B", "this was wrong"),
				latestFailedEvent: func(s uint64) (*repository.FailedEvent, error) {
					return &repository.FailedEvent{
						ErrMsg:         "blub",
						FailedSequence: s - 1,
						FailureCount:   5,
						ViewName:       "super.table",
					}, nil
				},
				errorCountUntilSkip: 6,
			},
			res: res{
				shouldProcessSequence: true,
			},
		},
		{
			name: "should not process sequence",
			args: args{
				event:     &models.Event{Sequence: 30000000},
				failedErr: errors.ThrowInternal(nil, "SPOOL-Wk53B", "this was wrong"),
				latestFailedEvent: func(s uint64) (*repository.FailedEvent, error) {
					return &repository.FailedEvent{
						ErrMsg:         "blub",
						FailedSequence: s - 1,
						FailureCount:   3,
						ViewName:       "super.table",
					}, nil
				},
				errorCountUntilSkip: 5,
			},
			res: res{
				shouldProcessSequence: false,
				wantErr:               true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processedSequence := false
			err := HandleError(
				tt.args.event,
				tt.args.failedErr,
				tt.args.latestFailedEvent,
				func(*repository.FailedEvent) error {
					return nil
				},
				func(*models.Event) error {
					processedSequence = true
					return nil
				},
				tt.args.errorCountUntilSkip)

			if (err != nil) != tt.res.wantErr {
				t.Errorf("HandleError() error = %v, wantErr %v", err, tt.res.wantErr)
			}
			if tt.res.shouldProcessSequence != processedSequence {
				t.Error("should not process sequence")
			}
		})
	}
}
