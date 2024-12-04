package initialise

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func Test_verifyUser(t *testing.T) {
	err := ReadStmts("cockroach") //TODO: check all dialects
	if err != nil {
		t.Errorf("unable to read stmts: %v", err)
		t.FailNow()
	}

	type args struct {
		db       db
		username string
		password string
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
					expectExec("-- replace zitadel-user with the name of the user\nCREATE USER IF NOT EXISTS \"zitadel-user\"", sql.ErrTxDone),
				),
				username: "zitadel-user",
				password: "",
			},
			targetErr: sql.ErrTxDone,
		},
		{
			name: "correct without password",
			args: args{
				db: prepareDB(t,
					expectExec("-- replace zitadel-user with the name of the user\nCREATE USER IF NOT EXISTS \"zitadel-user\"", nil),
				),
				username: "zitadel-user",
				password: "",
			},
			targetErr: nil,
		},
		{
			name: "correct with password",
			args: args{
				db: prepareDB(t,
					expectExec("-- replace zitadel-user with the name of the user\nCREATE USER IF NOT EXISTS \"zitadel-user\" WITH PASSWORD 'password'", nil),
				),
				username: "zitadel-user",
				password: "password",
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExec("-- replace zitadel-user with the name of the user\nCREATE USER IF NOT EXISTS \"zitadel-user\" WITH PASSWORD 'password'", nil),
				),
				username: "zitadel-user",
				password: "",
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyUser(tt.args.username, tt.args.password)(context.Background(), tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("VerifyGrant() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
