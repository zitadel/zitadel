package initialise

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/caos/zitadel/internal/database"
)

func Test_verifyGrant(t *testing.T) {
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
				db: prepareDB(t, expectQueryErr("SELECT EXISTS(SELECT * FROM [SHOW GRANTS ON DATABASE zitadel] where grantee = $1 AND privilege_type = 'ALL'", sql.ErrConnDone, "zitadel-user")),
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
					expectExists("SELECT EXISTS(SELECT * FROM [SHOW GRANTS ON DATABASE zitadel] where grantee = $1 AND privilege_type = 'ALL'", false, "zitadel-user"),
					expectExec("GRANT ALL ON DATABASE zitadel TO zitadel-user", sql.ErrTxDone),
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
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT * FROM [SHOW GRANTS ON DATABASE zitadel] where grantee = $1 AND privilege_type = 'ALL'", false, "zitadel-user"),
					expectExec("GRANT ALL ON DATABASE zitadel TO zitadel-user", nil),
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
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT * FROM [SHOW GRANTS ON DATABASE zitadel] where grantee = $1 AND privilege_type = 'ALL'", true, "zitadel-user"),
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
			if err := verifyGrant(tt.args.db.db, tt.args.config); !errors.Is(err, tt.targetErr) {
				t.Errorf("verifyGrant() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
