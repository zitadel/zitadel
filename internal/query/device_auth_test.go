package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

const (
	expectedDeviceAuthQueryC = `SELECT` +
		` projections.device_authorizations1.client_id,` +
		` projections.device_authorizations1.scopes,` +
		` projections.device_authorizations1.expires,` +
		` projections.device_authorizations1.state,` +
		` projections.device_authorizations1.subject,` +
		` projections.device_authorizations1.user_auth_methods,` +
		` projections.device_authorizations1.auth_time` +
		` FROM projections.device_authorizations1`
	expectedDeviceAuthWhereDeviceCodeQueryC = expectedDeviceAuthQueryC +
		` WHERE projections.device_authorizations1.client_id = $1` +
		` AND projections.device_authorizations1.device_code = $2` +
		` AND projections.device_authorizations1.instance_id = $3`
	expectedDeviceAuthWhereUserCodeQueryC = expectedDeviceAuthQueryC +
		` WHERE projections.device_authorizations1.instance_id = $1` +
		` AND projections.device_authorizations1.user_code = $2`
)

var (
	expectedDeviceAuthQuery                = regexp.QuoteMeta(expectedDeviceAuthQueryC)
	expectedDeviceAuthWhereDeviceCodeQuery = regexp.QuoteMeta(expectedDeviceAuthWhereDeviceCodeQueryC)
	expectedDeviceAuthWhereUserCodeQuery   = regexp.QuoteMeta(expectedDeviceAuthWhereUserCodeQueryC)
	expectedDeviceAuthValues               = []driver.Value{
		"client-id",
		database.TextArray[string]{"a", "b", "c"},
		testNow,
		domain.DeviceAuthStateApproved,
		"subject",
		database.Array[domain.UserAuthMethodType]{4},
		time.Unix(123, 456),
	}
	expectedDeviceAuth = &DeviceAuth{
		ClientID:        "client-id",
		Scopes:          []string{"a", "b", "c"},
		Expires:         testNow,
		State:           domain.DeviceAuthStateApproved,
		Subject:         "subject",
		UserAuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
		AuthTime:        time.Unix(123, 456),
	}
)

func TestQueries_DeviceAuthByDeviceCode(t *testing.T) {
	client, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to build mock client: %v", err)
	}
	defer client.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(expectedDeviceAuthWhereDeviceCodeQuery).WillReturnRows(
		sqlmock.NewRows(deviceAuthSelectColumns).AddRow(expectedDeviceAuthValues...),
	)
	mock.ExpectCommit()
	q := Queries{
		client: &database.DB{DB: client},
	}
	got, err := q.DeviceAuthByDeviceCode(context.TODO(), "123", "456")
	require.NoError(t, err)
	assert.Equal(t, expectedDeviceAuth, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestQueries_DeviceAuthByUserCode(t *testing.T) {
	client, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to build mock client: %v", err)
	}
	defer client.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(expectedDeviceAuthWhereUserCodeQuery).WillReturnRows(
		sqlmock.NewRows(deviceAuthSelectColumns).AddRow(expectedDeviceAuthValues...),
	)
	mock.ExpectCommit()
	q := Queries{
		client: &database.DB{DB: client},
	}
	got, err := q.DeviceAuthByUserCode(context.TODO(), "789")
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
		object any
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
			object: (*DeviceAuth)(nil),
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
			object: (*DeviceAuth)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, prepareDeviceAuthQuery, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
