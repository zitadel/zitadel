package initialise

import (
	"database/sql"
	"errors"
	"testing"
)

func Test_verifyUser(t *testing.T) {
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
			name: "exists fails",
			args: args{
				db:       prepareDB(t, expectQueryErr("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", sql.ErrConnDone, "zitadel-user")),
				username: "zitadel-user",
				password: "",
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "doesn't exists, create fails",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", false, "zitadel-user"),
					expectExec("CREATE USER $1 WITH PASSWORD $2", sql.ErrTxDone, "zitadel-user", nil),
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
					expectExists("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", false, "zitadel-user"),
					expectExec("CREATE USER $1 WITH PASSWORD $2", nil, "zitadel-user", nil),
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
					expectExists("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", false, "zitadel-user"),
					expectExec("CREATE USER $1 WITH PASSWORD $2", nil, "zitadel-user", "password"),
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
					expectExists("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", true, "zitadel-user"),
				),
				username: "zitadel-user",
				password: "",
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyUser(tt.args.username, tt.args.password)(tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("VerifyGrant() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
