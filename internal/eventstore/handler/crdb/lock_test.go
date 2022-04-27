package crdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	z_errs "github.com/zitadel/zitadel/internal/errors"
)

const (
	workerName     = "test_worker"
	projectionName = "my_projection"
	lockTable      = "my_lock_table"
)

var (
	renewNoRowsAffectedErr = z_errs.ThrowAlreadyExists(nil, "CRDB-mmi4J", "projection already locked")
	errLock                = errors.New("lock err")
)

func TestStatementHandler_handleLock(t *testing.T) {
	type want struct {
		expectations []mockExpectation
	}
	type args struct {
		lockDuration time.Duration
		ctx          context.Context
		errMock      *errsMock
		instanceID   string
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "lock fails",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 2, "instanceID"),
					expectLock(lockTable, workerName, 2, "instanceID"),
					expectLockErr(lockTable, workerName, 2, "instanceID", errLock),
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				ctx:          context.Background(),
				errMock: &errsMock{
					errs:            make(chan error),
					successfulIters: 2,
					shouldErr:       true,
				},
				instanceID: "instanceID",
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 2, "instanceID"),
					expectLock(lockTable, workerName, 2, "instanceID"),
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				ctx:          context.Background(),
				errMock: &errsMock{
					errs:            make(chan error),
					successfulIters: 2,
				},
				instanceID: "instanceID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			h := &locker{
				projectionName: projectionName,
				client:         client,
				workerName:     workerName,
				lockStmt:       fmt.Sprintf(lockStmtFormat, lockTable),
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			ctx, cancel := context.WithCancel(tt.args.ctx)

			go tt.args.errMock.handleErrs(t, cancel)

			go h.handleLock(ctx, tt.args.errMock.errs, tt.args.lockDuration, tt.args.instanceID)

			<-ctx.Done()

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func TestStatementHandler_renewLock(t *testing.T) {
	type want struct {
		expectations []mockExpectation
		isErr        func(err error) bool
	}
	type args struct {
		lockDuration time.Duration
		instanceID   string
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "lock fails",
			want: want{
				expectations: []mockExpectation{
					expectLockErr(lockTable, workerName, 1, "instanceID", sql.ErrTxDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
			args: args{
				lockDuration: 1 * time.Second,
				instanceID:   "instanceID",
			},
		},
		{
			name: "lock no rows",
			want: want{
				expectations: []mockExpectation{
					expectLockNoRows(lockTable, workerName, 2, "instanceID"),
				},
				isErr: func(err error) bool {
					return errors.As(err, &renewNoRowsAffectedErr)
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				instanceID:   "instanceID",
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 3, "instanceID"),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
			args: args{
				lockDuration: 3 * time.Second,
				instanceID:   "instanceID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			h := &locker{
				projectionName: projectionName,
				client:         client,
				workerName:     workerName,
				lockStmt:       fmt.Sprintf(lockStmtFormat, lockTable),
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			err = h.renewLock(context.Background(), tt.args.lockDuration, tt.args.instanceID)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error = %v", err)
			}

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

func TestStatementHandler_Unlock(t *testing.T) {
	type want struct {
		expectations []mockExpectation
		isErr        func(err error) bool
	}
	type args struct {
		instanceID string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "unlock fails",
			args: args{
				instanceID: "instanceID",
			},
			want: want{
				expectations: []mockExpectation{
					expectLockErr(lockTable, workerName, 0, "instanceID", sql.ErrTxDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
		},
		{
			name: "success",
			args: args{
				instanceID: "instanceID",
			},
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 0, "instanceID"),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			h := &locker{
				projectionName: projectionName,
				client:         client,
				workerName:     workerName,
				lockStmt:       fmt.Sprintf(lockStmtFormat, lockTable),
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			err = h.Unlock(tt.args.instanceID)
			if !tt.want.isErr(err) {
				t.Errorf("unexpected error = %v", err)
			}

			mock.MatchExpectationsInOrder(true)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectations not met: %v", err)
			}
		})
	}
}

type errsMock struct {
	errs            chan error
	successfulIters int
	shouldErr       bool
}

func (m *errsMock) handleErrs(t *testing.T, cancel func()) {
	for i := 0; i < m.successfulIters; i++ {
		if err := <-m.errs; err != nil {
			t.Errorf("unexpected err in iteration %d: %v", i, err)
			cancel()
			return
		}
	}
	if m.shouldErr {
		if err := <-m.errs; err == nil {
			t.Error("error must not be nil")
		}
	}
	cancel()
}
