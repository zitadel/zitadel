package initialise

import (
	"database/sql"
	"errors"
	"testing"
)

func Test_verifySchema(t *testing.T) {
	type args struct {
		db db
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "exists fails",
			args: args{
				db: prepareDB(t, expectQueryErr("SELECT EXISTS(SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1)", sql.ErrConnDone, "test_schema")),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "doesn't exists, create fails",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1)", false, "test_schema"),
					expectExec("CREATE SCHEMA test_schema", sql.ErrTxDone),
				),
			},
			targetErr: sql.ErrTxDone,
		},
		{
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1)", false, "test_schema"),
					expectExec("CREATE SCHEMA test_schema", nil),
				),
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT schema_name FROM [SHOW SCHEMAS] WHERE schema_name = $1)", true, "test_schema"),
				),
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifySchema(tt.args.db.db, "test_schema"); !errors.Is(err, tt.targetErr) {
				t.Errorf("verifySchema() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}

func Test_verifyEvents(t *testing.T) {
	type args struct {
		db db
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "exists fails",
			args: args{
				db: prepareDB(t, expectQueryErr("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", sql.ErrConnDone, eventsTable)),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "doesn't exists, unable to begin",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", false, eventsTable),
					expectBegin(sql.ErrConnDone),
				),
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "doesn't exists, hash sharded indexes fails",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", false, eventsTable),
					expectBegin(nil),
					expectExec("SET experimental_enable_hash_sharded_indexes = on", sql.ErrNoRows),
					expectRollback(nil),
				),
			},
			targetErr: sql.ErrNoRows,
		},
		{
			name: "doesn't exists, create table fails",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", false, eventsTable),
					expectBegin(nil),
					expectExec("SET experimental_enable_hash_sharded_indexes = on", nil),
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
					expectExists("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", false, eventsTable),
					expectBegin(nil),
					expectExec("SET experimental_enable_hash_sharded_indexes = on", nil),
					expectExec(createEventsStmt, nil),
					expectCommit(nil),
				),
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT table_name FROM [SHOW TABLES] WHERE table_name = $1)", true, eventsTable),
				),
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyEvents(tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("verifyEvents() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
