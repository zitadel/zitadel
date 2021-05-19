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
	lockErr   = errors.New("lock failed")
	unlockErr = errors.New("unlock failed")
	execErr   = errors.New("exec error")
	bulkErr   = errors.New("bulk err")
	updateErr = errors.New("update err")
)

func TestProjectionHandler_processEvent(t *testing.T) {
	type fields struct {
		stmts      []Statement
		pushSet    bool
		shouldPush chan *struct{}
	}
	type args struct {
		ctx    context.Context
		event  eventstore.EventReader
		reduce Reduce
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
			},
			args: args{
				reduce: testReduceErr(reduceErr),
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
			},
			args: args{
				reduce: testReduce(),
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
					NewNoOpStatement(1, 0),
				},
				pushSet:    false,
				shouldPush: make(chan *struct{}, 1),
			},
			args: args{
				reduce: testReduce(NewNoOpStatement(2, 1)),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
				stmts: []Statement{
					NewNoOpStatement(1, 0),
					NewNoOpStatement(2, 1),
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
			}
			err := h.processEvent(tt.args.ctx, tt.args.event, tt.args.reduce)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error %v", err)
			}
			if !reflect.DeepEqual(tt.want.stmts, h.stmts) {
				t.Errorf("unexpected stmts\n want: %v\n got: %v", tt.want.stmts, h.stmts)
			}
		})
	}
}

func TestProjectionHandler_fetchBulkStmts(t *testing.T) {
	type args struct {
		ctx    context.Context
		query  SearchQuery
		reduce Reduce
	}
	type want struct {
		shouldLimitExeeded bool
		isErr              func(error) bool
	}
	type fields struct {
		eventstore *eventstore.Eventstore
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
				ctx:    context.Background(),
				query:  testQuery(nil, 0, queryErr),
				reduce: testReduce(),
			},
			fields: fields{},
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
				ctx:    context.Background(),
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 5, nil),
				reduce: testReduce(),
			},
			fields: fields{
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
				ctx:    context.Background(),
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 5, nil),
				reduce: testReduce(),
			},
			fields: fields{
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
				ctx:    context.Background(),
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 5, nil),
				reduce: testReduce(),
			},
			fields: fields{
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
				ctx:    context.Background(),
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "test"), 2, nil),
				reduce: testReduce(),
			},
			fields: fields{
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
				shouldPush: make(chan *struct{}, 10),
			}
			gotLimitExeeded, err := h.fetchBulkStmts(tt.args.ctx, tt.args.query, tt.args.reduce)
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
	}
	type args struct {
		ctx          context.Context
		previousLock time.Duration
		update       Update
		reduce       Reduce
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
					NewNoOpStatement(1, 0),
					NewNoOpStatement(2, 1),
				},
				pushSet: true,
			},
			args: args{
				ctx:          context.Background(),
				previousLock: 200 * time.Millisecond,
				update:       testUpdate(t, 2, nil),
				reduce:       testReduce(),
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
					NewNoOpStatement(1, 0),
					NewNoOpStatement(2, 1),
				},
				pushSet: true,
			},
			args: args{
				ctx:    context.Background(),
				update: testUpdate(t, 2, errors.New("some error")),
				reduce: testReduce(),
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
			}
			if tt.args.previousLock > 0 {
				h.lockMu.Lock()
				go func() {
					<-time.After(tt.args.previousLock)
					h.lockMu.Unlock()
				}()
			}
			start := time.Now()
			if err := h.push(tt.args.ctx, tt.args.update, tt.args.reduce); !tt.want.isErr(err) {
				t.Errorf("ProjectionHandler.push() error = %v", err)
			}
			executionTime := time.Since(start)
			if tt.want.minExecution.Truncate(executionTime) > 0 {
				t.Errorf("expected execution time >= %v got %v", tt.want.minExecution, executionTime)
			}
			if h.pushSet {
				t.Error("expected push set to be false")
			}
			if len(h.stmts) != 0 {
				t.Errorf("expected stmts to be nil but was %v", h.stmts)
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
			go cancelOnErr(tt.args.ctx, tt.args.errs, tt.cancelMocker.mockCancel)
			if tt.args.err != nil {
				tt.args.errs <- tt.args.err
			}
			tt.cancelMocker.check(t)
		})
	}
}

func TestProjectionHandler_bulk(t *testing.T) {
	type args struct {
		ctx         context.Context
		executeBulk *executeBulkMock
		lock        *lockMock
		unlock      *unlockMock
	}
	type res struct {
		lockCount           int
		lockCanceled        bool
		executeBulkCount    int
		executeBulkCanceled bool
		unlockCount         int
		isErr               func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "lock fails",
			args: args{
				ctx:         context.Background(),
				executeBulk: &executeBulkMock{},
				lock: &lockMock{
					firstErr: lockErr,
					errWait:  time.Duration(500 * time.Millisecond),
				},
				unlock: &unlockMock{},
			},
			res: res{
				lockCount:        1,
				executeBulkCount: 0,
				unlockCount:      0,
				isErr: func(err error) bool {
					return errors.Is(err, lockErr)
				},
			},
		},
		{
			name: "unlock fails",
			args: args{
				ctx:         context.Background(),
				executeBulk: &executeBulkMock{},
				lock: &lockMock{
					err:     nil,
					errWait: time.Duration(500 * time.Millisecond),
				},
				unlock: &unlockMock{
					err: unlockErr,
				},
			},
			res: res{
				lockCount:        1,
				executeBulkCount: 1,
				unlockCount:      1,
				isErr: func(err error) bool {
					return errors.Is(err, unlockErr)
				},
			},
		},
		{
			name: "no error",
			args: args{
				ctx:         context.Background(),
				executeBulk: &executeBulkMock{},
				lock: &lockMock{
					err:      nil,
					errWait:  time.Duration(500 * time.Millisecond),
					canceled: make(chan bool, 1),
				},
				unlock: &unlockMock{
					err: nil,
				},
			},
			res: res{
				lockCount:        1,
				executeBulkCount: 1,
				unlockCount:      1,
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
		{
			name: "ctx canceled before lock",
			args: args{
				ctx:         canceledCtx(),
				executeBulk: &executeBulkMock{},
				lock: &lockMock{
					err:      nil,
					errWait:  time.Duration(500 * time.Millisecond),
					canceled: make(chan bool, 1),
				},
				unlock: &unlockMock{
					err: nil,
				},
			},
			res: res{
				lockCount:        1,
				lockCanceled:     true,
				executeBulkCount: 0,
				unlockCount:      0,
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
		{
			name: "2nd lock fails",
			args: args{
				ctx: context.Background(),
				executeBulk: &executeBulkMock{
					canceled:      make(chan bool, 1),
					waitForCancel: true,
				},
				lock: &lockMock{
					firstErr: nil,
					err:      lockErr,
					errWait:  time.Duration(100 * time.Millisecond),
					canceled: make(chan bool, 1),
				},
				unlock: &unlockMock{
					err: nil,
				},
			},
			res: res{
				lockCount:        1,
				lockCanceled:     true,
				executeBulkCount: 1,
				unlockCount:      1,
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
		{
			name: "bulk fails",
			args: args{
				ctx: context.Background(),
				executeBulk: &executeBulkMock{
					canceled:      make(chan bool, 1),
					err:           bulkErr,
					waitForCancel: false,
				},
				lock: &lockMock{
					firstErr: nil,
					err:      nil,
					errWait:  time.Duration(100 * time.Millisecond),
					canceled: make(chan bool, 1),
				},
				unlock: &unlockMock{
					err: nil,
				},
			},
			res: res{
				lockCount:        1,
				lockCanceled:     true,
				executeBulkCount: 1,
				unlockCount:      1,
				isErr: func(err error) bool {
					return errors.Is(err, bulkErr)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				RequeueAfter: time.Duration(0),
			}
			err := h.bulk(tt.args.ctx, tt.args.lock.lock(), tt.args.executeBulk.executeBulk(), tt.args.unlock.unlock())
			if !tt.res.isErr(err) {
				t.Errorf("unexpected error %v", err)
			}
			tt.args.lock.check(t, tt.res.lockCount, tt.res.lockCanceled)
			tt.args.executeBulk.check(t, tt.res.executeBulkCount, tt.res.executeBulkCanceled)
			tt.args.unlock.check(t, tt.res.unlockCount)
		})
	}
}

func TestProjectionHandler_prepareExecuteBulk(t *testing.T) {
	type fields struct {
		Handler       Handler
		SequenceTable string
		stmts         []Statement
		pushSet       bool
		shouldPush    chan *struct{}
	}
	type args struct {
		ctx    context.Context
		query  SearchQuery
		reduce Reduce
		update Update
	}
	type want struct {
		isErr func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "ctx done",
			args: args{
				ctx: canceledCtx(),
			},
			want: want{
				isErr: func(err error) bool {
					return err == nil
				},
			},
		},
		{
			name:   "fetch fails",
			fields: fields{},
			args: args{
				query: testQuery(nil, 10, ErrNoTable),
				ctx:   context.Background(),
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, ErrNoTable)
				},
			},
		},
		{
			name: "push fails",
			fields: fields{
				Handler: NewHandler(
					eventstore.NewEventstore(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								ID:               "id2",
								Sequence:         1,
								PreviousSequence: 0,
								CreationDate:     time.Now(),
								Type:             "test.added",
								Version:          "v1",
								AggregateID:      "testid",
								AggregateType:    "testAgg",
							},
							&repository.Event{
								ID:               "id2",
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
				),
				shouldPush: make(chan *struct{}, 1),
			},
			args: args{
				update: testUpdate(t, 2, updateErr),
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "testAgg"), 10, nil),
				reduce: testReduce(
					NewNoOpStatement(2, 1),
				),
				ctx: context.Background(),
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, updateErr)
				},
			},
		},
		{
			name: "success",
			fields: fields{
				Handler: NewHandler(
					eventstore.NewEventstore(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(
							&repository.Event{
								ID:               "id2",
								Sequence:         1,
								PreviousSequence: 0,
								CreationDate:     time.Now(),
								Type:             "test.added",
								Version:          "v1",
								AggregateID:      "testid",
								AggregateType:    "testAgg",
							},
							&repository.Event{
								ID:               "id2",
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
				),
				shouldPush: make(chan *struct{}, 1),
			},
			args: args{
				update: testUpdate(t, 4, nil),
				query:  testQuery(eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, "testAgg"), 10, nil),
				reduce: testReduce(
					NewNoOpStatement(1, 0),
					NewNoOpStatement(2, 1),
				),
				ctx: context.Background(),
			},
			want: want{
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ProjectionHandler{
				Handler:       tt.fields.Handler,
				SequenceTable: tt.fields.SequenceTable,
				lockMu:        sync.Mutex{},
				stmts:         tt.fields.stmts,
				pushSet:       tt.fields.pushSet,
				shouldPush:    tt.fields.shouldPush,
			}
			execBulk := h.prepareExecuteBulk(tt.args.query, tt.args.reduce, tt.args.update)
			err := execBulk(tt.args.ctx)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected err %v", err)
			}
		})
	}
}

func testUpdate(t *testing.T, expectedStmtCount int, returnedErr error) Update {
	return func(ctx context.Context, stmts []Statement, reduce Reduce) ([]Statement, error) {
		if expectedStmtCount != len(stmts) {
			t.Errorf("expected %d stmts got %d", expectedStmtCount, len(stmts))
		}
		return []Statement{}, returnedErr
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

type executeBulkMock struct {
	callCount     int
	err           error
	waitForCancel bool
	canceled      chan bool
}

func (m *executeBulkMock) executeBulk() executeBulk {
	return func(ctx context.Context) error {
		m.callCount++
		if m.waitForCancel {
			select {
			case <-ctx.Done():
				m.canceled <- true
				return nil
			case <-time.After(500 * time.Millisecond):
			}
		}
		return m.err
	}
}

func (m *executeBulkMock) check(t *testing.T, callCount int, shouldBeCalled bool) {
	t.Helper()
	if callCount != m.callCount {
		t.Errorf("wrong call count: expected %v got: %v", m.callCount, callCount)
	}
	if shouldBeCalled {
		select {
		case <-m.canceled:
		default:
			t.Error("bulk should be canceled but wasn't")
		}
	}
}

type lockMock struct {
	callCount int
	canceled  chan bool

	firstErr error
	err      error
	errWait  time.Duration
}

func (m *lockMock) lock() Lock {
	return func(ctx context.Context, _ time.Duration) <-chan error {
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
}

func (m *unlockMock) unlock() Unlock {
	return func() error {
		m.callCount++
		return m.err
	}
}

func (m *unlockMock) check(t *testing.T, callCount int) {
	t.Helper()
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
