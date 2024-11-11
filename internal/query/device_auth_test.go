package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	expectedDeviceAuthQueryC = `SELECT` +
		` projections.device_auth_requests2.client_id,` +
		` projections.device_auth_requests2.device_code,` +
		` projections.device_auth_requests2.user_code,` +
		` projections.device_auth_requests2.scopes,` +
		` projections.device_auth_requests2.audience` +
		` FROM projections.device_auth_requests2`
	expectedDeviceAuthWhereUserCodeQueryC = expectedDeviceAuthQueryC +
		` WHERE projections.device_auth_requests2.instance_id = $1` +
		` AND projections.device_auth_requests2.user_code = $2`
)

var (
	expectedDeviceAuthQuery              = regexp.QuoteMeta(expectedDeviceAuthQueryC)
	expectedDeviceAuthWhereUserCodeQuery = regexp.QuoteMeta(expectedDeviceAuthWhereUserCodeQueryC)
	expectedDeviceAuthValues             = []driver.Value{
		"client-id",
		"device1",
		"user-code",
		database.TextArray[string]{"a", "b", "c"},
		[]string{"projectID", "clientID"},
	}
	expectedDeviceAuth = &domain.AuthRequestDevice{
		ClientID:   "client-id",
		DeviceCode: "device1",
		UserCode:   "user-code",
		Scopes:     []string{"a", "b", "c"},
		Audience:   []string{"projectID", "clientID"},
	}
)

func TestQueries_DeviceAuthRequestByUserCode(t *testing.T) {
	client, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
	if err != nil {
		t.Fatalf("failed to build mock client: %v", err)
	}
	defer client.Close()

	mock.ExpectQuery(expectedDeviceAuthWhereUserCodeQuery).WillReturnRows(
		mock.NewRows(deviceAuthSelectColumns).AddRow(expectedDeviceAuthValues...),
	)
	q := Queries{
		client: &database.DB{DB: client},
	}
	got, err := q.DeviceAuthRequestByUserCode(context.TODO(), "789")
	require.NoError(t, err)
	assert.Equal(t, expectedDeviceAuth, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func Test_prepareDeviceAuthQuery(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name   string
		want   want
		object *domain.AuthRequestDevice
	}{
		{
			name: "success",
			want: want{
				sqlExpectations: mockQueries(
					expectedDeviceAuthQuery,
					deviceAuthSelectColumns,
					[][]driver.Value{expectedDeviceAuthValues},
				),
			},
			object: expectedDeviceAuth,
		},
		{
			name: "not found error",
			want: want{
				sqlExpectations: mockQueryErr(
					expectedDeviceAuthQuery,
					sql.ErrNoRows,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrNoRows) {
						return fmt.Errorf("err should be sql.ErrNoRows got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
		{
			name: "other error",
			want: want{
				sqlExpectations: mockQueryErr(
					expectedDeviceAuthQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, prepareDeviceAuthQuery, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
