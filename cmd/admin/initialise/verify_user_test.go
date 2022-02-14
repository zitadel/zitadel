package initialise

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/caos/zitadel/internal/database"
)

func Test_verifyUser(t *testing.T) {
	type args struct {
		db     db
		config database.Config
	}
	tests := []struct {
		name      string
		args      args
		targetErr error
	}{
		{
			name: "exists fails",
			args: args{
				db: prepareDB(t, expectQueryErr("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", sql.ErrConnDone, "zitadel-user")),
				config: database.Config{
					Database: "zitadel",
					User: database.User{
						Username: "zitadel-user",
					},
				},
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
				config: database.Config{
					Database: "zitadel",
					User: database.User{
						Username: "zitadel-user",
					},
				},
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
				config: database.Config{
					Database: "zitadel",
					User: database.User{
						Username: "zitadel-user",
					},
				},
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
				config: database.Config{
					Database: "zitadel",
					User: database.User{
						Username: "zitadel-user",
						Password: "password",
					},
				},
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT username FROM [show roles] WHERE username = $1)", true, "zitadel-user"),
				),
				config: database.Config{
					Database: "zitadel",
					User: database.User{
						Username: "zitadel-user",
					},
				},
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyUser(tt.args.config)(tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("verifyGrant() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
