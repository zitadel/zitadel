package setup

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
)

const userGrants5InstanceUserIDCreate = "-- Lets user grant lists of an instance read grant ids ordered by user_id from the index, without sorting the joined rows.\n" +
	"CREATE INDEX CONCURRENTLY IF NOT EXISTS user_grants5_instance_user_id_idx ON projections.user_grants5 (instance_id, user_id, id);\n"

func TestUserGrants5InstanceUserIDIndex_Execute(t *testing.T) {
	tests := []struct {
		name    string
		expects func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "table missing, no create",
			expects: func(mock sqlmock.Sqlmock) {
				expectUserGrants5TableExists(mock, false)
			},
		},
		{
			name: "table exists, no invalid index, creates index",
			expects: func(mock sqlmock.Sqlmock) {
				expectUserGrants5TableExists(mock, true)
				expectInvalidIndex(mock, userGrants5InstanceUserIDIdx, false)
				mock.ExpectExec(regexp.QuoteMeta(userGrants5InstanceUserIDCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "invalid leftover index, drop then create",
			expects: func(mock sqlmock.Sqlmock) {
				expectUserGrants5TableExists(mock, true)
				expectInvalidIndex(mock, userGrants5InstanceUserIDIdx, true)
				mock.ExpectExec(regexp.QuoteMeta(dropInvalidIndexConcurrently + userGrants5InstanceUserIDIdx)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(userGrants5InstanceUserIDCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "create fails",
			expects: func(mock sqlmock.Sqlmock) {
				expectUserGrants5TableExists(mock, true)
				expectInvalidIndex(mock, userGrants5InstanceUserIDIdx, false)
				mock.ExpectExec(regexp.QuoteMeta(userGrants5InstanceUserIDCreate)).
					WillReturnError(errors.New("create failed"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer func() {
				require.NoError(t, mock.ExpectationsWereMet())
			}()
			defer db.Close()
			tt.expects(mock)

			mig := &UserGrants5InstanceUserIDIndex{
				dbClient: &database.DB{DB: db},
			}
			err = mig.Execute(context.Background(), nil)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func expectUserGrants5TableExists(mock sqlmock.Sqlmock, exists bool) {
	mock.ExpectQuery(regexp.QuoteMeta(userGrants5TableExistsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(exists))
}
