package spooler

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

type testHandler struct {
	cycleDuration time.Duration
	processSleep  time.Duration
	processError  error
	queryError    error
	viewModel     string
}

func (h *testHandler) ViewModel() string {
	return h.viewModel
}
func (h *testHandler) EventQuery() (*models.SearchQuery, error) {
	return nil, h.queryError
}
func (h *testHandler) Process(*models.Event) error {
	<-time.After(h.processSleep)
	return h.processError
}
func (h *testHandler) MinimumCycleDuration() time.Duration { return h.cycleDuration }

type eventstoreStub struct {
	events []*models.Event
	err    error
}

// Health returns status OK as soon as the service started
func (es *eventstoreStub) Health(ctx context.Context) error {
	return nil
}

// Ready returns status OK as soon as all dependent services are available
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

func TestSpooler_process(t *testing.T) {
	type fields struct {
		currentHandler Handler
	}
	type args struct {
		timeout time.Duration
		events  []*models.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "process all events",
			fields: fields{
				currentHandler: &testHandler{},
			},
			args: args{
				timeout: 0,
				events:  []*models.Event{&models.Event{}, &models.Event{}},
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
				events:  []*models.Event{&models.Event{}, &models.Event{}, &models.Event{}, &models.Event{}},
			},
			wantErr: false,
		},
		{
			name: "process error",
			fields: fields{
				currentHandler: &testHandler{processSleep: 1 * time.Second, processError: fmt.Errorf("i am an error")},
			},
			args: args{
				events: []*models.Event{&models.Event{}, &models.Event{}},
			},
			wantErr: true,
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

			if err := s.process(ctx, tt.args.events); (err != nil) != tt.wantErr {
				t.Errorf("Spooler.process() error = %v, wantErr %v", err, tt.wantErr)
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
		currentHandler Handler
		err            error
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
			},
		},
		{
			"with error",
			fields{
				err:            fmt.Errorf("hodor"),
				currentHandler: &testHandler{processSleep: 500 * time.Millisecond},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &spooledHandler{
				Handler: tt.fields.currentHandler,
			}

			canceled := false
			errs := make(chan error)
			cancel := func() {
				canceled = true
			}

			go s.awaitError(cancel, errs)
			errs <- tt.fields.err

			if !canceled {
				t.Error("cancel function was not called")
			}
		})
	}
}

// TestSpooler_load checks if load terminates
func TestSpooler_load(t *testing.T) {
	type fields struct {
		currentHandler Handler
		locker         *testLocker
		lockID         string
		eventstore     eventstore.Eventstore
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"lock exists",
			fields{
				currentHandler: &testHandler{processSleep: 500 * time.Millisecond, viewModel: "testView", cycleDuration: 1 * time.Second},
				lockID:         "testID",
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, fmt.Errorf("lock already exists"), 2000*time.Millisecond),
			},
		},
		{
			"lock fails",
			fields{
				currentHandler: &testHandler{processSleep: 100 * time.Millisecond, viewModel: "testView", cycleDuration: 1 * time.Second},
				lockID:         "testID",
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, fmt.Errorf("fail"), 2000*time.Millisecond),
				eventstore:     &eventstoreStub{events: []*models.Event{&models.Event{}}},
			},
		},
		{
			"query fails",
			fields{
				currentHandler: &testHandler{processSleep: 100 * time.Millisecond, viewModel: "testView", queryError: fmt.Errorf("query fail"), cycleDuration: 1 * time.Second},
				lockID:         "testID",
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, nil, 2000*time.Millisecond),
				eventstore:     &eventstoreStub{err: fmt.Errorf("fail")},
			},
		},
		{
			"process event fails",
			fields{
				currentHandler: &testHandler{processError: fmt.Errorf("oups"), processSleep: 100 * time.Millisecond, viewModel: "testView", cycleDuration: 500 * time.Millisecond},
				lockID:         "testID",
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, nil, 1000*time.Millisecond),
				eventstore:     &eventstoreStub{events: []*models.Event{&models.Event{}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.fields.locker.finish()
			s := &spooledHandler{
				Handler:    tt.fields.currentHandler,
				locker:     tt.fields.locker.mock,
				lockID:     tt.fields.lockID,
				eventstore: tt.fields.eventstore,
			}
			s.load()
		})
	}
}

func TestSpooler_lock(t *testing.T) {
	type fields struct {
		currentHandler Handler
		locker         *testLocker
		lockID         string
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
				lockID:         "testID",
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, nil, 2000*time.Millisecond),
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
				lockID:         "testID",
				locker:         newTestLocker(t, "testID", "testView").expectRenew(t, fmt.Errorf("renew failed"), 1800*time.Millisecond),
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
				lockID:  tt.fields.lockID,
			}

			errs := make(chan error, 1)
			ctx, _ := context.WithDeadline(context.Background(), tt.args.deadline)

			locked := s.lock(ctx, errs)

			if tt.fields.expectsErr {
				err := <-errs
				if err == nil {
					t.Error("No error in error queue")
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
	l.mock.EXPECT().Renew(l.lockerID, l.viewName, gomock.Any()).DoAndReturn(
		func(_, _ string, gotten time.Duration) error {
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
