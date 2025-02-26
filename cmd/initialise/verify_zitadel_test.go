package initialise

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
)

func Test_verifyEvents(t *testing.T) {
	err := ReadStmts("cockroach") //TODO: check all dialects
	if err != nil {
		t.Errorf("unable to read stmts: %v", err)
		t.FailNow()
	}

	type args struct {
		db db
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "unable to begin",
			args: args{
				db: prepareDB(t,
					expectBegin(sql.ErrConnDone),
				),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "events already exists",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectQuery(
						"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'eventstore' AND table_name like 'events%'",
						nil,
						[]string{"count"},
						[][]driver.Value{
							{1},
						},
					),
					expectCommit(nil),
				),
			},
		},
		{
			name: "events and events2 already exists",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectQuery(
						"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'eventstore' AND table_name like 'events%'",
						nil,
						[]string{"count"},
						[][]driver.Value{
							{2},
						},
					),
					expectCommit(nil),
				),
			},
		},
		{
			name: "create table fails",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectQuery(
						"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'eventstore' AND table_name like 'events%'",
						nil,
						[]string{"count"},
						[][]driver.Value{
							{0},
						},
					),
					expectExec(createEventsStmt, sql.ErrNoRows),
					expectRollback(nil),
				),
			},
			targetErr: sql.ErrNoRows,
		},
		{
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectQuery(
						"SELECT count(*) FROM information_schema.tables WHERE table_schema = 'eventstore' AND table_name like 'events%'",
						nil,
						[]string{"count"},
						[][]driver.Value{
							{0},
						},
					),
					expectExec(createEventsStmt, nil),
					expectCommit(nil),
				),
			},
			targetErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := tt.args.db.db.Conn(context.Background())
			if err != nil {
				t.Error(err)
				return
			}
			if err := createEvents(context.Background(), conn); !errors.Is(err, tt.targetErr) {
				t.Errorf("createEvents() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_verifyEncryptionKeys(t *testing.T) {
	type args struct {
		db db
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "unable to begin",
			args: args{
				db: prepareDB(t,
					expectBegin(sql.ErrConnDone),
				),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "create table fails",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectExec(createEncryptionKeysStmt, sql.ErrNoRows),
					expectRollback(nil),
				),
			},
			targetErr: sql.ErrNoRows,
		},
		{
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectBegin(nil),
					expectExec(createEncryptionKeysStmt, nil),
					expectCommit(nil),
				),
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createEncryptionKeys(context.Background(), tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("createEvents() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
