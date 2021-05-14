package handler

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/caos/zitadel/internal/eventstore/repository/mock"
)

var (
	queryErr  = errors.New("query err")
	filterErr = errors.New("filter err")
	reduceErr = errors.New("reduce err")
)

func TestProjectionHandler_processEvent(t *testing.T) {
	type fields struct {
		stmts      []Statement
		pushSet    bool
		shouldPush chan *struct{}
		reduce     Reduce
	}
	type args struct {
		ctx   context.Context
		event eventstore.EventReader
	}
	type want struct {
		isErr func(err error) bool
		stmts []Statement
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "reduce fails",
			fields: fields{
				stmts:      nil,
				pushSet:    false,
				shouldPush: nil,
				reduce:     testReduceErr(reduceErr),
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, reduceErr)
				},
				stmts: nil,
			},
		},
		{
			name: "no stmts",
			fields: fields{
				stmts:      nil,
				pushSet:    false,
				shouldPush: make(chan *struct{}, 1),
				reduce:     testReduce(),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				stmts: nil,
			},
		},
		{
			name: "existing stmts",
			fields: fields{
				stmts: []Statement{
					NewNoOpStatement("my_table", 1, 0),
				},
				pushSet:    false,
				shouldPush: make(chan *struct{}, 1),
				reduce:     testReduce(NewNoOpStatement("my_table", 2, 1)),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				stmts: []Statement{
					NewNoOpStatement("my_table", 1, 0),
					NewNoOpStatement("my_table", 2, 1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				lockMu:     sync.Mutex{},
				stmts:      tt.fields.stmts,
				pushSet:    tt.fields.pushSet,
				shouldPush: tt.fields.shouldPush,
				reduce:     tt.fields.reduce,
			}
			err := h.processEvent(tt.args.ctx, tt.args.event)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error %v", err)
			}
			if !reflect.DeepEqual(tt.want.stmts, h.stmts) {
				t.Errorf("unexpected stmts\n want: %v\n got: %v", tt.want.stmts, h.stmts)
			}
		})
	}
}

func TestProjectionHandler_prepareBulkStmts(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type want struct {
		shouldLimitExeeded bool
		isErr              func(error) bool
	}
	type fields struct {
		eventstore *eventstore.Eventstore
		query      SearchQuery
		reduce     Reduce
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
				query:  testQuery(nil, 0, queryErr),
				reduce: testReduce(),
			},
			want: want{
				shouldLimitExeeded: false,
				isErr: func(err error) bool {
					return errors.Is(err, queryErr)
				},
			},
		},
		{
			name: "eventstore returns err",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 5, nil),
				reduce: testReduce(),
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEventsError(filterErr),
				),
			},
			want: want{
				shouldLimitExeeded: false,
				isErr: func(err error) bool {
					return errors.Is(err, filterErr)
				},
			},
		},
		{
			name: "no events found",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 5, nil),
				reduce: testReduce(),
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(),
				),
			},
			want: want{
				shouldLimitExeeded: false,
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
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 5, nil),
				reduce: testReduce(),
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							ID:               "id",
							Sequence:         1,
							PreviousSequence: 0,
							CreationDate:     time.Now(),
							Type:             "test.added",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
						&repository.Event{
							ID:               "id",
							Sequence:         2,
							PreviousSequence: 1,
							CreationDate:     time.Now(),
							Type:             "test.changed",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
					),
				),
			},
			want: want{
				shouldLimitExeeded: false,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name: "found events exeed limit",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 2, nil),
				reduce: testReduce(),
				eventstore: eventstore.NewEventstore(
					es_repo_mock.NewRepo(t).ExpectFilterEvents(
						&repository.Event{
							ID:               "id",
							Sequence:         1,
							PreviousSequence: 0,
							CreationDate:     time.Now(),
							Type:             "test.added",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
						&repository.Event{
							ID:               "id",
							Sequence:         2,
							PreviousSequence: 1,
							CreationDate:     time.Now(),
							Type:             "test.changed",
							Version:          "v1",
							AggregateID:      "testid",
							AggregateType:    "testAgg",
						},
					),
				),
			},
			want: want{
				shouldLimitExeeded: true,
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				lockMu: sync.Mutex{},
				Handler: Handler{
					Eventstore: tt.fields.eventstore,
				},
				query:      tt.fields.query,
				reduce:     tt.fields.reduce,
				shouldPush: make(chan *struct{}, 10),
			}
			gotLimitExeeded, err := h.prepareBulkStmts(tt.args.ctx)
			if !tt.want.isErr(err) {
				t.Errorf("ProjectionHandler.prepareBulkStmts() error = %v", err)
				return
			}
			if gotLimitExeeded != tt.want.shouldLimitExeeded {
				t.Errorf("ProjectionHandler.prepareBulkStmts() = %v, want %v", gotLimitExeeded, tt.want.shouldLimitExeeded)
			}
		})
	}
}

func TestProjectionHandler_push(t *testing.T) {
	type fields struct {
		stmts   []Statement
		pushSet bool
		update  Update
		reduce  Reduce
	}
	type args struct {
		ctx          context.Context
		previousLock time.Duration
	}
	type want struct {
		isErr        func(err error) bool
		minExecution time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "previous lock",
			fields: fields{
				stmts: []Statement{
					NewNoOpStatement("my_table", 1, 0),
					NewNoOpStatement("my_table", 2, 1),
				},
				pushSet: true,
				update:  testUpdate(t, 2, nil),
				reduce:  testReduce(),
			},
			args: args{
				ctx:          context.Background(),
				previousLock: 200 * time.Millisecond,
			},
			want: want{
				isErr:        func(err error) bool { return err == nil },
				minExecution: 200 * time.Millisecond,
			},
		},
		{
			name: "error in update",
			fields: fields{
				stmts: []Statement{
					NewNoOpStatement("my_table", 1, 0),
					NewNoOpStatement("my_table", 2, 1),
				},
				pushSet: true,
				update:  testUpdate(t, 2, errors.New("some error")),
				reduce:  testReduce(),
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				isErr: func(err error) bool { return err.Error() == "some error" },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				lockMu:  sync.Mutex{},
				stmts:   tt.fields.stmts,
				pushSet: tt.fields.pushSet,
				update:  tt.fields.update,
				reduce:  tt.fields.reduce,
			}
			if tt.args.previousLock > 0 {
				h.lockMu.Lock()
				go func() {
					<-time.After(tt.args.previousLock)
					h.lockMu.Unlock()
				}()
			}
			start := time.Now()
			if err := h.push(tt.args.ctx); !tt.want.isErr(err) {
				t.Errorf("ProjectionHandler.push() error = %v", err)
			}
			executionTime := time.Since(start)
			if tt.want.minExecution.Truncate(executionTime) > 0 {
				t.Errorf("expected execution time >= %v got %v", tt.want.minExecution, executionTime)
			}
			if h.pushSet {
				t.Error("expected push set to be false")
			}
			if h.stmts != nil {
				t.Errorf("expected stmts to be nil but was %v", h.stmts)
			}
		})
	}
}

func testUpdate(t *testing.T, expectedStmtCount int, returnedErr error) Update {
	return func(ctx context.Context, stmts []Statement, reduce Reduce) error {
		if expectedStmtCount != len(stmts) {
			t.Errorf("expected %d stmts got %d", expectedStmtCount, len(stmts))
		}
		return returnedErr
	}
}

func testReduce(stmts ...Statement) Reduce {
	return func(event eventstore.EventReader) ([]Statement, error) {
		return stmts, nil
	}
}

func testReduceErr(err error) Reduce {
	return func(event eventstore.EventReader) ([]Statement, error) {
		return nil, err
	}
}

func testQuery(query *eventstore.SearchQueryBuilder, limit uint64, err error) SearchQuery {
	return func() (*eventstore.SearchQueryBuilder, uint64, error) {
		return query, limit, err
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
	if wasCalled := <-m.wasCalled; m.shouldBeCalled != wasCalled {
		t.Errorf("cancel: should: %t got: %t", m.shouldBeCalled, wasCalled)
	}
}

func Test_cancelOnErr(t *testing.T) {
	type args struct {
		ctx  context.Context
		errs chan error
	}
	tests := []struct {
		name         string
		args         args
		cancelMocker *cancelMocker
	}{
		{
			name: "nil error occured",
			args: args{
				ctx:  canceledCtx(),
				errs: make(chan error),
			},
			cancelMocker: &cancelMocker{
				shouldBeCalled: false,
			},
		},
		{
			name: "ctx done",
			args: args{
				ctx:  context.Background(),
				errs: make(chan error),
			},
			cancelMocker: &cancelMocker{
				shouldBeCalled: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go cancelOnErr(tt.args.ctx, tt.args.errs, tt.cancelMocker.mockCancel)
			if tt.cancelMocker.shouldBeCalled {
				tt.args.errs <- errors.New("cancel")
			}
			tt.cancelMocker.check(t)
		})
	}
}
