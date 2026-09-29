package setup

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
)

const (
	users14InstanceResourceOwnerCreate = "CREATE INDEX CONCURRENTLY IF NOT EXISTS users14_instance_resource_owner_idx ON projections.users14 (instance_id, resource_owner, id) INCLUDE (type);\n"
	users14ResourceOwnerDrop           = "DROP INDEX CONCURRENTLY IF EXISTS projections.users14_resource_owner_idx;\n"
	users14UsernameDrop                = "DROP INDEX CONCURRENTLY IF EXISTS projections.users14_username_idx;\n"
)

func TestUsers14InstanceResourceOwnerIndex_Execute(t *testing.T) {
	tests := []struct {
		name    string
		expects func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "table missing, no drop or create",
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(users14TableExistsQuery)).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
		},
		{
			name: "table exists, index valid, creates then drops old indexes",
			expects: func(mock sqlmock.Sqlmock) {
				expectUsers14TableExists(mock, true)
				expectInvalidIndex(mock, users14InstanceResourceOwnerIdx, false)
				mock.ExpectExec(regexp.QuoteMeta(users14InstanceResourceOwnerCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14ResourceOwnerDrop)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14UsernameDrop)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "new index invalid, drop then create then drop old indexes",
			expects: func(mock sqlmock.Sqlmock) {
				expectUsers14TableExists(mock, true)
				expectInvalidIndex(mock, users14InstanceResourceOwnerIdx, true)
				mock.ExpectExec(regexp.QuoteMeta(dropInvalidIndexConcurrently + users14InstanceResourceOwnerIdx)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14InstanceResourceOwnerCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14ResourceOwnerDrop)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14UsernameDrop)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
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

			mig := &Users14InstanceResourceOwnerIndex{
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
