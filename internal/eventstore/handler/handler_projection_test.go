package handler

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
)

var (
	ErrQuery  = errors.New("query err")
	ErrFilter = errors.New("filter err")
	ErrReduce = errors.New("reduce err")
	ErrLock   = errors.New("lock failed")
	ErrUnlock = errors.New("unlock failed")
	ErrExec   = errors.New("exec error")
	ErrBulk   = errors.New("bulk err")
	ErrUpdate = errors.New("update err")
)

func TestProjectionHandler_Trigger(t *testing.T) {
	type fields struct {
		reduce     Reduce
		update     Update
		query      SearchQuery
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx       context.Context
		instances []string
	}
	type want struct {
		isErr func(err error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			"query error",
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return nil
				},
				query: testQuery(nil, 0, ErrQuery),
			},
			args{
				context.Background(),
				nil,
			},
			want{isErr: func(err error) bool {
				return errors.Is(err, ErrQuery)
			}},
		},
		{
			"no events",
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(
						eventstore.TestConfig(es_repo_mock.NewRepo(t).ExpectFilterEvents()),
					)
				},
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					5,
					nil),
			},
			args{
				context.Background(),
				nil,
			},
			want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			"process error",
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								ID:                        "id",
								Sequence:                  1,
								PreviousAggregateSequence: 0,
								CreationDate:              time.Now(),
								Type:                      "test.added",
								Version:                   "v1",
								AggregateID:               "testid",
								AggregateType:             "testAgg",
							},
						),
					),
					)
				},
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					5,
					nil),
				reduce: testReduceErr(ErrReduce),
			},
			args{
				context.Background(),
				nil,
			},
			want{
				isErr: func(err error) bool {
					return errors.Is(err, ErrReduce)
				},
			},
		},
		{
			"process ok",
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								ID:                        "id",
								Sequence:                  1,
								PreviousAggregateSequence: 0,
								CreationDate:              time.Now(),
								Type:                      "test.added",
								Version:                   "v1",
								AggregateID:               "testid",
								AggregateType:             "testAgg",
							},
						),
					),
					)
				},
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					5,
					nil),
				reduce: testReduce(newTestStatement("testAgg", 1, 0)),
				update: testUpdate(t, 1, 0, nil),
			},
			args{
				context.Background(),
				nil,
			},
			want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			"process limit exceeded ok",
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).
							ExpectFilterEvents(
								&repository.Event{
									ID:                        "id",
									Sequence:                  1,
									PreviousAggregateSequence: 0,
									CreationDate:              time.Now(),
									Type:                      "test.added",
									Version:                   "v1",
									AggregateID:               "testid",
									AggregateType:             "testAgg",
								},
							).ExpectFilterEvents(),
					),
					)
				},
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					1,
					nil),
				reduce: testReduce(newTestStatement("testAgg", 1, 0)),
				update: testUpdate(t, 1, 0, nil),
			},
			args{
				context.Background(),
				nil,
			},
			want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				Handler: Handler{
					Eventstore: tt.fields.eventstore(t),
				},
				ProjectionName: "test",
				reduce:         tt.fields.reduce,
				update:         tt.fields.update,
				searchQuery:    tt.fields.query,
			}

			err := h.Trigger(tt.args.ctx, tt.args.instances...)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

func TestProjectionHandler_Process(t *testing.T) {
	type fields struct {
		reduce Reduce
		update Update
	}
	type args struct {
		ctx    context.Context
		events []eventstore.Event
	}
	type want struct {
		isErr func(err error) bool
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "no events",
			fields: fields{},
			args:   args{},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				index: 0,
			},
		},
		{
			name: "reduce fails",
			fields: fields{
				reduce: testReduceErr(ErrReduce),
			},
			args: args{
				events: []eventstore.Event{newTestEvent("id", "description", nil)},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, ErrReduce)
				},
				index: -1,
			},
		},
		{
			name: "stmt failed",
			fields: fields{
				reduce: testReduce(newTestStatement("aggregate1", 1, 0)),
				update: testUpdate(t, 1, -1, ErrSomeStmtsFailed),
			},
			args: args{
				events: []eventstore.Event{newTestEvent("id", "description", nil)},
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, ErrSomeStmtsFailed)
				},
				index: -1,
			},
		},
		{
			name: "stmt error",
			fields: fields{
				reduce: testReduce(newTestStatement("aggregate1", 1, 0)),
				update: testUpdate(t, 1, -1, errors.New("some error")),
			},
			args: args{
				events: []eventstore.Event{newTestEvent("id", "description", nil)},
			},
			want: want{
				isErr: func(err error) bool {
					return err.Error() == "some error"
				},
				index: -1,
			},
		},
		{
			name: "stmt succeeded",
			fields: fields{
				reduce: testReduce(newTestStatement("aggregate1", 1, 0)),
				update: testUpdate(t, 1, 0, nil),
			},
			args: args{
				events: []eventstore.Event{newTestEvent("id", "description", nil)},
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				index: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewProjectionHandler(
				context.Background(),
				ProjectionHandlerConfig{
					HandlerConfig: HandlerConfig{
						Eventstore: nil,
					},
					ProjectionName: "test",
					RequeueEvery:   -1,
				},
				tt.fields.reduce,
				tt.fields.update,
				nil,
				nil,
				nil,
				nil,
			)

			index, err := h.Process(tt.args.ctx, tt.args.events...)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error %v", err)
			}
			assert.Equal(t, tt.want.index, index)
		})
	}
}

func TestProjectionHandler_FetchEvents(t *testing.T) {
	type args struct {
		ctx         context.Context
		instanceIDs []string
	}
	type want struct {
		limitExceeded bool
		isErr         func(error) bool
	}
	type fields struct {
		eventstore *eventstore.Eventstore
		query      SearchQuery
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			name: "query returns err",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				query: testQuery(nil, 0, ErrQuery),
			},
			want: want{
				limitExceeded: false,
				isErr: func(err error) bool {
					return errors.Is(err, ErrQuery)
				},
			},
		},
		{
			name: "eventstore returns err",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(eventstore.TestConfig(
					es_repo_mock.NewRepo(t).ExpectFilterEventsError(ErrFilter),
				),
				),
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					5,
					nil,
				),
			},
			want: want{
				limitExceeded: false,
				isErr: func(err error) bool {
					return errors.Is(err, ErrFilter)
				},
			},
		},
		{
			name: "no events found",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(
					eventstore.TestConfig(es_repo_mock.NewRepo(t).ExpectFilterEvents()),
				),
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					5,
					nil,
				),
			},
			want: want{
				limitExceeded: false,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "found events smaller than limit",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(eventstore.TestConfig(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							ID:                        "id",
							Sequence:                  1,
							PreviousAggregateSequence: 0,
							CreationDate:              time.Now(),
							Type:                      "test.added",
							Version:                   "v1",
							AggregateID:               "testid",
							AggregateType:             "testAgg",
						},
						&repository.Event{
							ID:                        "id",
							Sequence:                  2,
							PreviousAggregateSequence: 1,
							CreationDate:              time.Now(),
							Type:                      "test.changed",
							Version:                   "v1",
							AggregateID:               "testid",
							AggregateType:             "testAgg",
						},
					),
				),
				),
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					5,
					nil,
				),
			},
			want: want{
				limitExceeded: false,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "found events exceeds limit",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				eventstore: eventstore.NewEventstore(eventstore.TestConfig(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							ID:                        "id",
							Sequence:                  1,
							PreviousAggregateSequence: 0,
							CreationDate:              time.Now(),
							Type:                      "test.added",
							Version:                   "v1",
							AggregateID:               "testid",
							AggregateType:             "testAgg",
						},
						&repository.Event{
							ID:                        "id",
							Sequence:                  2,
							PreviousAggregateSequence: 1,
							CreationDate:              time.Now(),
							Type:                      "test.changed",
							Version:                   "v1",
							AggregateID:               "testid",
							AggregateType:             "testAgg",
						},
					),
				),
				),
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					2,
					nil,
				),
			},
			want: want{
				limitExceeded: true,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				Handler: Handler{
					Eventstore: tt.fields.eventstore,
				},
				searchQuery: tt.fields.query,
			}
			_, limitExceeded, err := h.FetchEvents(tt.args.ctx, tt.args.instanceIDs...)
			if !tt.want.isErr(err) {
				t.Errorf("ProjectionHandler.prepareBulkStmts() error = %v", err)
				return
			}
			if limitExceeded != tt.want.limitExceeded {
				t.Errorf("ProjectionHandler.prepareBulkStmts() = %v, want %v", limitExceeded, tt.want.limitExceeded)
			}
		})
	}
}

func TestProjection_subscribe(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		reduce Reduce
		update Update
		events []eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		fields fields
	}{
		{
			"panic",
			args{
				ctx: context.Background(),
			},
			fields{
				reduce: nil,
				update: nil,
				events: []eventstore.Event{
					newTestEvent("id", "", nil),
				},
			},
		},
		{
			"error",
			args{
				ctx: context.Background(),
			},
			fields{
				reduce: testReduceErr(ErrReduce),
				update: nil,
				events: []eventstore.Event{
					newTestEvent("id", "", nil),
				},
			},
		},
		{
			"not all statement",
			args{
				ctx: context.Background(),
			},
			fields{
				reduce: testReduce(newTestStatement("aggregate1", 1, 0)),
				update: testUpdate(t, 1, 0, ErrSomeStmtsFailed),
				events: []eventstore.Event{
					newTestEvent("id", "", nil),
				},
			},
		},
		{
			"single event ok",
			args{
				ctx: context.Background(),
			},
			fields{
				reduce: testReduce(newTestStatement("aggregate1", 1, 0)),
				update: testUpdate(t, 1, 1, nil),
				events: []eventstore.Event{
					newTestEvent("id", "", nil),
				},
			},
		},
		{
			"multiple events ok",
			args{
				ctx: context.Background(),
			},
			fields{
				reduce: testReduce(newTestStatement("aggregate1", 1, 0)),
				update: testUpdate(t, 2, 2, nil),
				events: []eventstore.Event{
					newTestEvent("id", "", nil),
					newTestEvent("id2", "", nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				Handler: Handler{
					EventQueue: make(chan eventstore.Event, 10),
				},
				reduce: tt.fields.reduce,
				update: tt.fields.update,
			}
			ctx, cancel := context.WithCancel(tt.args.ctx)
			go func() {
				//changed go h.subscribe(ctx) to this to be able to ignore logs easily
				t.Helper()
				h.subscribe(ctx)
			}()
			for _, event := range tt.fields.events {
				h.EventQueue <- event
			}
			time.Sleep(1 * time.Second)
			cancel()
		})
	}
}

func TestProjection_schedule(t *testing.T) {

	now := func() time.Time {
		return time.Date(2023, 1, 31, 12, 0, 0, 0, time.UTC)
	}

	type args struct {
		ctx context.Context
	}
	type fields struct {
		reduce                  Reduce
		update                  Update
		eventstore              func(t *testing.T) *eventstore.Eventstore
		lock                    *lockMock
		unlock                  *unlockMock
		query                   SearchQuery
		handleInactiveInstances bool
	}
	type want struct {
		locksCount   int
		lockCanceled bool
		unlockCount  int
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   want
	}{
		{
			"panic",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return nil
				},
			},
			want{},
		},
		{
			"filter succeeded once error",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEventsError(ErrFilter),
					),
					)
				},
				handleInactiveInstances: false,
			},
			want{
				locksCount:   0,
				lockCanceled: false,
				unlockCount:  0,
			},
		},
		{
			"filter instance ids error",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								AggregateType:             "system",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "",
								Type:                      "system.projections.scheduler.succeeded",
							}).
							ExpectInstanceIDsError(ErrFilter),
					),
					)
				},
				handleInactiveInstances: false,
			},
			want{
				locksCount:   0,
				lockCanceled: false,
				unlockCount:  0,
			},
		},
		{
			"lock error",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								AggregateType:             "system",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "",
								Type:                      "system.projections.scheduler.succeeded",
							}).ExpectInstanceIDs(nil, "instanceID1"),
					),
					)
				},
				lock: &lockMock{
					errWait:  100 * time.Millisecond,
					firstErr: ErrLock,
					canceled: make(chan bool, 1),
				},
				handleInactiveInstances: false,
			},
			want{
				locksCount:   1,
				lockCanceled: true,
				unlockCount:  0,
			},
		},
		{
			"trigger error",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								AggregateType:             "system",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "",
								Type:                      "system.projections.scheduler.succeeded",
							}).ExpectInstanceIDs(nil, "instanceID1"),
					),
					)
				},
				lock: &lockMock{
					canceled: make(chan bool, 1),
					firstErr: nil,
					errWait:  100 * time.Millisecond,
				},
				unlock:                  &unlockMock{},
				query:                   testQuery(nil, 0, ErrQuery),
				handleInactiveInstances: false,
			},
			want{
				locksCount:   1,
				lockCanceled: true,
				unlockCount:  1,
			},
		},
		{
			"only active instances are handled",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).
							ExpectFilterEvents(&repository.Event{
								AggregateType:             "system",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "",
								Type:                      "system.projections.scheduler.succeeded",
							}).
							ExpectInstanceIDs(
								[]*repository.Filter{{
									Field:     repository.FieldInstanceID,
									Operation: repository.OperationNotIn,
									Value:     database.StringArray{""},
								}, {
									Field:     repository.FieldCreationDate,
									Operation: repository.OperationGreater,
									Value:     now().Add(-2 * time.Hour),
								}},
								"206626268110651755",
							).
							ExpectFilterEvents(&repository.Event{
								AggregateType:             "quota",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "206626268110651755",
								Type:                      "quota.notificationdue",
							}),
					))
				},
				lock: &lockMock{
					canceled: make(chan bool, 1),
					firstErr: nil,
					errWait:  100 * time.Millisecond,
				},
				unlock:                  &unlockMock{},
				handleInactiveInstances: false,
				reduce:                  testReduce(newTestStatement("aggregate1", 1, 0)),
				update:                  testUpdate(t, 1, 1, nil),
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					2,
					nil,
				),
			},
			want{
				locksCount:   1,
				lockCanceled: false,
				unlockCount:  1,
			},
		},
		{
			"all instances are handled",
			args{
				ctx: context.Background(),
			},
			fields{
				eventstore: func(t *testing.T) *eventstore.Eventstore {
					return eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).
							ExpectFilterEvents(&repository.Event{
								AggregateType:             "system",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "",
								Type:                      "system.projections.scheduler.succeeded",
							}).
							ExpectInstanceIDs([]*repository.Filter{{
								Field:     repository.FieldInstanceID,
								Operation: repository.OperationNotIn,
								Value:     database.StringArray{""},
							}}, "206626268110651755").
							ExpectFilterEvents(&repository.Event{
								AggregateType:             "quota",
								Sequence:                  6,
								PreviousAggregateSequence: 5,
								InstanceID:                "206626268110651755",
								Type:                      "quota.notificationdue",
							}),
					))
				},
				lock: &lockMock{
					canceled: make(chan bool, 1),
					firstErr: nil,
					errWait:  100 * time.Millisecond,
				},
				unlock:                  &unlockMock{},
				handleInactiveInstances: true,
				reduce:                  testReduce(newTestStatement("aggregate1", 1, 0)),
				update:                  testUpdate(t, 1, 1, nil),
				query: testQuery(
					eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
						AddQuery().
						AggregateTypes("test").
						Builder(),
					2,
					nil,
				),
			},
			want{
				locksCount:   1,
				lockCanceled: false,
				unlockCount:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				Handler: Handler{
					EventQueue: make(chan eventstore.Event, 10),
					Eventstore: tt.fields.eventstore(t),
				},
				reduce:                  tt.fields.reduce,
				update:                  tt.fields.update,
				searchQuery:             tt.fields.query,
				lock:                    tt.fields.lock.lock(),
				unlock:                  tt.fields.unlock.unlock(),
				triggerProjection:       time.NewTimer(0), // immediately run an iteration
				requeueAfter:            time.Hour,        // run only one iteration
				concurrentInstances:     1,
				handleInactiveInstances: tt.fields.handleInactiveInstances,
				retries:                 0,
				nowFunc:                 now,
			}
			ctx, cancel := context.WithCancel(tt.args.ctx)
			go func() {
				//changed go h.schedule(ctx) to this to be able to ignore logs easily
				t.Helper()
				h.schedule(ctx)
			}()

			time.Sleep(time.Second)
			cancel()
			if tt.fields.lock != nil {
				tt.fields.lock.check(t, tt.want.locksCount, tt.want.lockCanceled)
			}
			if tt.fields.unlock != nil {
				tt.fields.unlock.check(t, tt.want.unlockCount)
			}
		})
	}
}

func Test_cancelOnErr(t *testing.T) {
	type args struct {
		ctx  context.Context
		errs chan error
		err  error
	}
	tests := []struct {
		name         string
		args         args
		cancelMocker *cancelMocker
	}{
		{
			name: "error occured",
			args: args{
				ctx:  context.Background(),
				errs: make(chan error),
				err:  ErrNoCondition,
			},
			cancelMocker: &cancelMocker{
				shouldBeCalled: true,
				wasCalled:      make(chan bool, 1),
			},
		},
		{
			name: "ctx done",
			args: args{
				ctx:  canceledCtx(),
				errs: make(chan error),
			},
			cancelMocker: &cancelMocker{
				shouldBeCalled: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{}
			go h.cancelOnErr(tt.args.ctx, tt.args.errs, tt.cancelMocker.mockCancel)
			if tt.args.err != nil {
				tt.args.errs <- tt.args.err
			}
			tt.cancelMocker.check(t)
		})
	}
}

func newTestStatement(aggType eventstore.AggregateType, seq, previousSeq uint64) *Statement {
	return &Statement{
		AggregateType:    aggType,
		Sequence:         seq,
		PreviousSequence: previousSeq,
	}
}

// testEvent implements the Event interface
type testEvent struct {
	eventstore.BaseEvent

	description string
	data        func() interface{}
}

func newTestEvent(id, description string, data func() interface{}) *testEvent {
	return &testEvent{
		description: description,
		data:        data,
		BaseEvent: *eventstore.NewBaseEventForPush(
			service.WithService(authz.NewMockContext("instanceID", "resourceOwner", "editorUser"), "editorService"),
			eventstore.NewAggregate(authz.NewMockContext("zitadel", "caos", "adlerhurst"), id, "test.aggregate", "v1"),
			"test.event",
		),
	}
}

func (e *testEvent) Data() interface{} {
	if e.data == nil {
		return nil
	}
	return e.data()
}

func (e *testEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func testUpdate(t *testing.T, expectedStmtCount, returnedIndex int, returnedErr error) Update {
	return func(ctx context.Context, stmts []*Statement, reduce Reduce) (int, error) {
		if expectedStmtCount != len(stmts) {
			t.Errorf("expected %d stmts got %d", expectedStmtCount, len(stmts))
		}
		return returnedIndex, returnedErr
	}
}

func testReduce(stmts *Statement) Reduce {
	return func(event eventstore.Event) (*Statement, error) {
		return stmts, nil
	}
}

func testReduceErr(err error) Reduce {
	return func(event eventstore.Event) (*Statement, error) {
		return nil, err
	}
}

func testQuery(query *eventstore.SearchQueryBuilder, limit uint64, err error) SearchQuery {
	return func(ctx context.Context, instanceIDs []string) (*eventstore.SearchQueryBuilder, uint64, error) {
		return query, limit, err
	}
}

type lockMock struct {
	callCount int
	canceled  chan bool
	mu        sync.Mutex

	firstErr error
	err      error
	errWait  time.Duration
}

func (m *lockMock) lock() Lock {
	return func(ctx context.Context, _ time.Duration, _ ...string) <-chan error {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.callCount++
		errs := make(chan error)
		go func() {
			for i := 0; ; i++ {
				select {
				case <-ctx.Done():
					m.canceled <- true
					close(errs)
					return
				case <-time.After(m.errWait):
					err := m.err
					if i == 0 {
						err = m.firstErr
					}
					errs <- err
				}
			}
		}()
		return errs
	}
}

func (m *lockMock) check(t *testing.T, callCount int, shouldBeCanceled bool) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()
	if callCount != m.callCount {
		t.Errorf("wrong call count: expected %v got: %v", callCount, m.callCount)
	}
	if shouldBeCanceled {
		select {
		case <-m.canceled:
		case <-time.After(5 * time.Second):
			t.Error("lock should be canceled but wasn't")
		}
	}
}

type unlockMock struct {
	callCount int
	err       error
	mu        sync.Mutex
}

func (m *unlockMock) unlock() Unlock {
	return func(...string) error {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.callCount++
		return m.err
	}
}

func (m *unlockMock) check(t *testing.T, callCount int) {
	t.Helper()
	m.mu.Lock()
	defer m.mu.Unlock()
	if callCount != m.callCount {
		t.Errorf("wrong call count: expected %v got: %v", callCount, m.callCount)
	}
}

func canceledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

type cancelMocker struct {
	shouldBeCalled bool
	wasCalled      chan bool
}

func (m *cancelMocker) mockCancel() {
	m.wasCalled <- true
}

func (m *cancelMocker) check(t *testing.T) {
	t.Helper()
	if m.shouldBeCalled {
		if wasCalled := <-m.wasCalled; !wasCalled {
			t.Errorf("cancel: should: %t got: %t", m.shouldBeCalled, wasCalled)
		}
	}
}
