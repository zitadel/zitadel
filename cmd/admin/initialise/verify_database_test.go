package initialise

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/caos/zitadel/internal/database"
)

func Test_verifyDB(t *testing.T) {
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
				db: prepareDB(t, expectQueryErr("SELECT EXISTS(SELECT database_name FROM [show databases] WHERE database_name = $1)", sql.ErrConnDone, "zitadel")),
				config: database.Config{
					Database: "zitadel",
				},
			},
			targetErr: sql.ErrConnDone,
		},
		{
			name: "doesn't exists, create fails",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT database_name FROM [show databases] WHERE database_name = $1)", false, "zitadel"),
					expectExec("CREATE DATABASE zitadel", sql.ErrTxDone),
				),
				config: database.Config{
					Database: "zitadel",
				},
			},
			targetErr: sql.ErrTxDone,
		},
		{
			name: "doesn't exists, create successful",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT database_name FROM [show databases] WHERE database_name = $1)", false, "zitadel"),
					expectExec("CREATE DATABASE zitadel", nil),
				),
				config: database.Config{
					Database: "zitadel",
				},
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExists("SELECT EXISTS(SELECT database_name FROM [show databases] WHERE database_name = $1)", true, "zitadel"),
				),
				config: database.Config{
					Database: "zitadel",
				},
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyDatabase(tt.args.config)(tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("verifyDB() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
