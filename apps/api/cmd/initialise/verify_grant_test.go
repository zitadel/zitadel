package initialise

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func Test_verifyGrant(t *testing.T) {
	type args struct {
		db       db
		database string
		username string
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
					expectExec("GRANT ALL ON DATABASE \"zitadel\" TO \"zitadel-user\"", sql.ErrTxDone),
				),
				database: "zitadel",
				username: "zitadel-user",
			},
			targetErr: sql.ErrTxDone,
		},
		{
			name: "correct",
			args: args{
				db: prepareDB(t,
					expectExec("GRANT ALL ON DATABASE \"zitadel\" TO \"zitadel-user\"", nil),
				),
				database: "zitadel",
				username: "zitadel-user",
			},
			targetErr: nil,
		},
		{
			name: "already exists",
			args: args{
				db: prepareDB(t,
					expectExec("GRANT ALL ON DATABASE \"zitadel\" TO \"zitadel-user\"", nil),
				),
				database: "zitadel",
				username: "zitadel-user",
			},
			targetErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyGrant(tt.args.database, tt.args.username)(context.Background(), tt.args.db.db); !errors.Is(err, tt.targetErr) {
				t.Errorf("VerifyGrant() error = %v, want: %v", err, tt.targetErr)
			}
			if err := tt.args.db.mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
