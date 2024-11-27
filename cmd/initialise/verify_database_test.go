package initialise

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func Test_verifyDB(t *testing.T) {
	err := ReadStmts("cockroach") //TODO: check all dialects
	if err != nil {
		t.Errorf("unable to read stmts: %v", err)
		t.FailNow()
	}

	type args struct {
		db       db
		database string
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "doesn't exists, create fails",
			args: args{
				db: prepareDB(t,
					expectExec("-- replace zitadel with the name of the database\nCREATE DATABASE IF NOT EXISTS \"zitadel\"", sql.ErrTxDone),
				),
				database: "zitadel",
			},
			targetErr: sql.ErrTxDone,
		},
		{
			name: "doesn't exists, create successful",
			args: args{
				db: prepareDB(t,
					expectExec("-- replace zitadel with the name of the database\nCREATE DATABASE IF NOT EXISTS \"zitadel\"", nil),
				),
				database: "zitadel",
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExec("-- replace zitadel with the name of the database\nCREATE DATABASE IF NOT EXISTS \"zitadel\"", nil),
				),
				database: "zitadel",
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyDatabase(tt.args.database)(context.Background(), tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("verifyDB() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
