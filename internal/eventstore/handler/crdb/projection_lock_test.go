package crdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	z_errs "github.com/caos/zitadel/internal/errors"

	"github.com/DATA-DOG/go-sqlmock"
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
					expectLock(lockTable, workerName, 2),
					expectLock(lockTable, workerName, 2),
					expectLockErr(lockTable, workerName, 2, errLock),
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
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 2),
					expectLock(lockTable, workerName, 2),
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				ctx:          context.Background(),
				errMock: &errsMock{
					errs:            make(chan error),
					successfulIters: 2,
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
			h := &StatementHandler{
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

			go h.handleLock(ctx, tt.args.errMock.errs, tt.args.lockDuration)

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
					expectLockErr(lockTable, workerName, 1, sql.ErrTxDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
			args: args{
				lockDuration: 1 * time.Second,
			},
		},
		{
			name: "lock no rows",
			want: want{
				expectations: []mockExpectation{
					expectLockNoRows(lockTable, workerName, 2),
				},
				isErr: func(err error) bool {
					return errors.As(err, &renewNoRowsAffectedErr)
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 3),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
			args: args{
				lockDuration: 3 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			h := &StatementHandler{
				projectionName: projectionName,
				client:         client,
				workerName:     workerName,
				lockStmt:       fmt.Sprintf(lockStmtFormat, lockTable),
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			err = h.renewLock(context.Background(), tt.args.lockDuration)
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
	tests := []struct {
		name string
		want want
	}{
		{
			name: "unlock fails",
			want: want{
				expectations: []mockExpectation{
					expectLockErr(lockTable, workerName, 0, sql.ErrTxDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 0),
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
			h := &StatementHandler{
				projectionName: projectionName,
				client:         client,
				workerName:     workerName,
				lockStmt:       fmt.Sprintf(lockStmtFormat, lockTable),
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			err = h.Unlock()
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
