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
	users14UsernameLowerCreate = "CREATE INDEX CONCURRENTLY IF NOT EXISTS users14_username_lower_idx ON projections.users14 (instance_id, LOWER(username));\n"
	users14PhoneLowerCreate    = "CREATE INDEX CONCURRENTLY IF NOT EXISTS users14_humans_phone_lower_idx ON projections.users14_humans (instance_id, LOWER(phone));\n"
)

func TestUsers14LoginEqualityIndexes_Execute(t *testing.T) {
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
			name: "table exists, indexes valid, creates both concurrently",
			expects: func(mock sqlmock.Sqlmock) {
				expectUsers14TableExists(mock, true)
				expectInvalidIndex(mock, users14UsernameLowerIdx, false)
				expectInvalidIndex(mock, users14HumansPhoneLowerIdx, false)
				mock.ExpectExec(regexp.QuoteMeta(users14UsernameLowerCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14PhoneLowerCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			name: "username index invalid, drop then create both",
			expects: func(mock sqlmock.Sqlmock) {
				expectUsers14TableExists(mock, true)
				expectInvalidIndex(mock, users14UsernameLowerIdx, true)
				mock.ExpectExec(regexp.QuoteMeta(dropInvalidIndexConcurrently + users14UsernameLowerIdx)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				expectInvalidIndex(mock, users14HumansPhoneLowerIdx, false)
				mock.ExpectExec(regexp.QuoteMeta(users14UsernameLowerCreate)).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(regexp.QuoteMeta(users14PhoneLowerCreate)).
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

			mig := &Users14LoginEqualityIndexes{
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

func expectUsers14TableExists(mock sqlmock.Sqlmock, exists bool) {
	mock.ExpectQuery(regexp.QuoteMeta(users14TableExistsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(exists))
}

func expectInvalidIndex(mock sqlmock.Sqlmock, name string, invalid bool) {
	mock.ExpectQuery(regexp.QuoteMeta(users14InvalidIndexQuery)).
		WithArgs(name).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(invalid))
}
