package crdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	workerName     = "test_worker"
	projectionName = "my_projection"
	lockTable      = "my_lock_table"
)

var (
	renewNoRowsAffectedErr = zerrors.ThrowAlreadyExists(nil, "CRDB-mmi4J", "projection already locked")
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
		instanceIDs  []string
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
					expectLock(lockTable, workerName, 2*time.Second, "instanceID"),
					expectLock(lockTable, workerName, 2*time.Second, "instanceID"),
					expectLockErr(lockTable, workerName, 2*time.Second, "instanceID", errLock),
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
				instanceIDs: []string{"instanceID"},
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 2*time.Second, "instanceID"),
					expectLock(lockTable, workerName, 2*time.Second, "instanceID"),
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				ctx:          context.Background(),
				errMock: &errsMock{
					errs:            make(chan error),
					successfulIters: 2,
				},
				instanceIDs: []string{"instanceID"},
			},
		},
		{
			name: "success with multiple",
			want: want{
				expectations: []mockExpectation{
					expectLockMultipleInstances(lockTable, workerName, 2*time.Second, "instanceID1", "instanceID2"),
					expectLockMultipleInstances(lockTable, workerName, 2*time.Second, "instanceID1", "instanceID2"),
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				ctx:          context.Background(),
				errMock: &errsMock{
					errs:            make(chan error),
					successfulIters: 2,
				},
				instanceIDs: []string{"instanceID1", "instanceID2"},
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
				lockStmt: func(values string, instances int) string {
					return fmt.Sprintf(lockStmtFormat, lockTable, values, instances)
				},
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			ctx, cancel := context.WithCancel(tt.args.ctx)

			go tt.args.errMock.handleErrs(t, cancel)

			go h.handleLock(ctx, tt.args.errMock.errs, tt.args.lockDuration, tt.args.instanceIDs...)

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
		instanceIDs  []string
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
					expectLockErr(lockTable, workerName, 1*time.Second, "instanceID", sql.ErrTxDone),
				},
				isErr: func(err error) bool {
					return errors.Is(err, sql.ErrTxDone)
				},
			},
			args: args{
				lockDuration: 1 * time.Second,
				instanceIDs:  database.TextArray[string]{"instanceID"},
			},
		},
		{
			name: "lock no rows",
			want: want{
				expectations: []mockExpectation{
					expectLockNoRows(lockTable, workerName, 2*time.Second, "instanceID"),
				},
				isErr: func(err error) bool {
					return errors.Is(err, renewNoRowsAffectedErr)
				},
			},
			args: args{
				lockDuration: 2 * time.Second,
				instanceIDs:  database.TextArray[string]{"instanceID"},
			},
		},
		{
			name: "success",
			want: want{
				expectations: []mockExpectation{
					expectLock(lockTable, workerName, 3*time.Second, "instanceID"),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
			args: args{
				lockDuration: 3 * time.Second,
				instanceIDs:  database.TextArray[string]{"instanceID"},
			},
		},
		{
			name: "success with multiple",
			want: want{
				expectations: []mockExpectation{
					expectLockMultipleInstances(lockTable, workerName, 3*time.Second, "instanceID1", "instanceID2"),
				},
				isErr: func(err error) bool {
					return errors.Is(err, nil)
				},
			},
			args: args{
				lockDuration: 3 * time.Second,
				instanceIDs:  []string{"instanceID1", "instanceID2"},
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
				lockStmt: func(values string, instances int) string {
					return fmt.Sprintf(lockStmtFormat, lockTable, values, instances)
				},
			}

			for _, expectation := range tt.want.expectations {
				expectation(mock)
			}

			err = h.renewLock(context.Background(), tt.args.lockDuration, tt.args.instanceIDs...)
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
				lockStmt: func(values string, instances int) string {
					return fmt.Sprintf(lockStmtFormat, lockTable, values, instances)
				},
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
