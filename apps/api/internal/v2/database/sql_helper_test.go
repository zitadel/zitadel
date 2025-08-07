package database

import (
	"errors"
	"reflect"
	"testing"
)

func TestCloseTx(t *testing.T) {
	type args struct {
		tx  *testTx
		err error
	}
	tests := []struct {
		name      string
		args      args
		assertErr func(t *testing.T, err error) bool
	}{
		{
			name: "exec err",
			args: args{
				tx: &testTx{
					rollback: execution{
						shouldExecute: true,
					},
				},
				err: errExec,
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errExec)
				if !is {
					t.Errorf("execution error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "exec err and rollback err",
			args: args{
				tx: &testTx{
					rollback: execution{
						err:           true,
						shouldExecute: true,
					},
				},
				err: errExec,
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errExec)
				if !is {
					t.Errorf("execution error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "commit Err",
			args: args{
				tx: &testTx{
					commit: execution{
						err:           true,
						shouldExecute: true,
					},
				},
				err: nil,
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errCommit)
				if !is {
					t.Errorf("commit error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "no err",
			args: args{
				tx: &testTx{
					commit: execution{
						shouldExecute: true,
					},
				},
				err: nil,
			},
			assertErr: func(t *testing.T, err error) bool {
				is := err == nil
				if !is {
					t.Errorf("no error expected, got: %v", err)
				}
				return is
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CloseTx(tt.args.tx, tt.args.err)
			tt.assertErr(t, err)
			tt.args.tx.assert(t)
		})
	}
}

func TestMapRows(t *testing.T) {
	type args struct {
		rows   *testRows
		mapper DestMapper[string]
	}
	var emptyString string
	tests := []struct {
		name       string
		args       args
		wantResult []*string
		assertErr  func(t *testing.T, err error) bool
	}{
		{
			name: "no rows, close err",
			args: args{
				rows: &testRows{
					closeErr: true,
				},
				mapper: nil,
			},
			wantResult: nil,
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errClose)
				if !is {
					t.Errorf("close error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "no rows, close err",
			args: args{
				rows: &testRows{
					hasErr: true,
				},
				mapper: nil,
			},
			wantResult: nil,
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errRows)
				if !is {
					t.Errorf("rows error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "scan err",
			args: args{
				rows: &testRows{
					scanErr:   true,
					nextCount: 1,
				},
				mapper: func(index int, scan func(dest ...any) error) (*string, error) {
					var s string
					if err := scan(&s); err != nil {
						return nil, err
					}
					return &s, nil
				},
			},
			wantResult: nil,
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errScan)
				if !is {
					t.Errorf("scan error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "exec err",
			args: args{
				rows: &testRows{
					nextCount: 1,
				},
				mapper: func(index int, scan func(dest ...any) error) (*string, error) {
					return nil, errExec
				},
			},
			wantResult: nil,
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errExec)
				if !is {
					t.Errorf("exec error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "exec err, close err",
			args: args{
				rows: &testRows{
					closeErr:  true,
					nextCount: 1,
				},
				mapper: func(index int, scan func(dest ...any) error) (*string, error) {
					return nil, errExec
				},
			},
			wantResult: nil,
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errExec)
				if !is {
					t.Errorf("exec error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "rows err",
			args: args{
				rows: &testRows{
					nextCount: 1,
					hasErr:    true,
				},
				mapper: func(index int, scan func(dest ...any) error) (*string, error) {
					var s string
					if err := scan(&s); err != nil {
						return nil, err
					}
					return &s, nil
				},
			},
			wantResult: nil,
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errRows)
				if !is {
					t.Errorf("rows error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "no err",
			args: args{
				rows: &testRows{
					nextCount: 1,
				},
				mapper: func(index int, scan func(dest ...any) error) (*string, error) {
					var s string
					if err := scan(&s); err != nil {
						return nil, err
					}
					return &s, nil
				},
			},
			wantResult: []*string{&emptyString},
			assertErr: func(t *testing.T, err error) bool {
				is := err == nil
				if !is {
					t.Errorf("no error expected, got: %v", err)
				}
				return is
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := MapRows(tt.args.rows, tt.args.mapper)
			tt.assertErr(t, err)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("MapRows() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestMapRowsToObject(t *testing.T) {
	type args struct {
		rows   *testRows
		mapper func(scan func(dest ...any) error) error
	}
	tests := []struct {
		name      string
		args      args
		assertErr func(t *testing.T, err error) bool
	}{
		{
			name: "no rows, close err",
			args: args{
				rows: &testRows{
					closeErr: true,
				},
				mapper: nil,
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errClose)
				if !is {
					t.Errorf("close error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "no rows, close err",
			args: args{
				rows: &testRows{
					hasErr: true,
				},
				mapper: nil,
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errRows)
				if !is {
					t.Errorf("rows error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "scan err",
			args: args{
				rows: &testRows{
					scanErr:   true,
					nextCount: 1,
				},
				mapper: func(scan func(dest ...any) error) error {
					var s string
					if err := scan(&s); err != nil {
						return err
					}
					return nil
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errScan)
				if !is {
					t.Errorf("scan error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "exec err",
			args: args{
				rows: &testRows{
					nextCount: 1,
				},
				mapper: func(scan func(dest ...any) error) error {
					return errExec
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errExec)
				if !is {
					t.Errorf("exec error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "exec err, close err",
			args: args{
				rows: &testRows{
					closeErr:  true,
					nextCount: 1,
				},
				mapper: func(scan func(dest ...any) error) error {
					return errExec
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errExec)
				if !is {
					t.Errorf("exec error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "rows err",
			args: args{
				rows: &testRows{
					nextCount: 1,
					hasErr:    true,
				},
				mapper: func(scan func(dest ...any) error) error {
					var s string
					return scan(&s)
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := errors.Is(err, errRows)
				if !is {
					t.Errorf("rows error expected, got: %v", err)
				}
				return is
			},
		},
		{
			name: "no err",
			args: args{
				rows: &testRows{
					nextCount: 1,
				},
				mapper: func(scan func(dest ...any) error) error {
					var s string
					return scan(&s)
				},
			},
			assertErr: func(t *testing.T, err error) bool {
				is := err == nil
				if !is {
					t.Errorf("no error expected, got: %v", err)
				}
				return is
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MapRowsToObject(tt.args.rows, tt.args.mapper)
			tt.assertErr(t, err)
		})
	}
}

var _ Tx = (*testTx)(nil)

type testTx struct {
	commit, rollback execution
}

type execution struct {
	err           bool
	didExecute    bool
	shouldExecute bool
}

var (
	errCommit   = errors.New("commit err")
	errRollback = errors.New("rollback err")
	errExec     = errors.New("exec err")
)

// Commit implements Tx.
func (t *testTx) Commit() error {
	t.commit.didExecute = true
	if t.commit.err {
		return errCommit
	}
	return nil
}

// Rollback implements Tx.
func (t *testTx) Rollback() error {
	t.rollback.didExecute = true
	if t.rollback.err {
		return errRollback
	}
	return nil
}

func (tx *testTx) assert(t *testing.T) {
	if tx.commit.didExecute != tx.commit.shouldExecute {
		t.Errorf("unexpected execution of commit: should %v, did: %v", tx.commit.shouldExecute, tx.commit.didExecute)
	}
	if tx.rollback.didExecute != tx.rollback.shouldExecute {
		t.Errorf("unexpected execution of rollback: should %v, did: %v", tx.rollback.shouldExecute, tx.rollback.didExecute)
	}
}

var _ Rows = (*testRows)(nil)

var (
	errClose = errors.New("err close")
	errRows  = errors.New("err rows")
	errScan  = errors.New("err scan")
)

type testRows struct {
	closeErr  bool
	scanErr   bool
	hasErr    bool
	nextCount int
}

// Close implements Rows.
func (t *testRows) Close() error {
	if t.closeErr {
		return errClose
	}
	return nil
}

// Err implements Rows.
func (t *testRows) Err() error {
	if t.hasErr {
		return errRows
	}
	if t.closeErr {
		return errClose
	}
	return nil
}

// Next implements Rows.
func (t *testRows) Next() bool {
	t.nextCount--
	return t.nextCount >= 0
}

// Scan implements Rows.
func (t *testRows) Scan(dest ...any) error {
	if t.scanErr {
		return errScan
	}
	return nil
}
