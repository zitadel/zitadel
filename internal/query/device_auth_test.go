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
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

const (
	expectedDeviceAuthQueryC = `SELECT` +
		` projections.device_authorizations.id,` +
		` projections.device_authorizations.client_id,` +
		` projections.device_authorizations.scopes,` +
		` projections.device_authorizations.expires,` +
		` projections.device_authorizations.state,` +
		` projections.device_authorizations.subject` +
		` FROM projections.device_authorizations`
	expectedDeviceAuthWhereDeviceCodeQueryC = expectedDeviceAuthQueryC +
		` WHERE projections.device_authorizations.client_id = $1` +
		` AND projections.device_authorizations.device_code = $2` +
		` AND projections.device_authorizations.instance_id = $3`
	expectedDeviceAuthWhereUserCodeQueryC = expectedDeviceAuthQueryC +
		` WHERE projections.device_authorizations.instance_id = $1` +
		` AND projections.device_authorizations.user_code = $2`
)

var (
	expectedDeviceAuthQuery                = regexp.QuoteMeta(expectedDeviceAuthQueryC)
	expectedDeviceAuthWhereDeviceCodeQuery = regexp.QuoteMeta(expectedDeviceAuthWhereDeviceCodeQueryC)
	expectedDeviceAuthWhereUserCodeQuery   = regexp.QuoteMeta(expectedDeviceAuthWhereUserCodeQueryC)
	expectedDeviceAuthValues               = []driver.Value{
		"primary-id",
		"client-id",
		database.StringArray{"a", "b", "c"},
		testNow,
		domain.DeviceAuthStateApproved,
		"subject",
	}
	expectedDeviceAuth = &domain.DeviceAuth{
		ObjectRoot: models.ObjectRoot{
			AggregateID: "primary-id",
		},
		ClientID: "client-id",
		Scopes:   []string{"a", "b", "c"},
		Expires:  testNow,
		State:    domain.DeviceAuthStateApproved,
		Subject:  "subject",
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
			object: (*domain.DeviceAuth)(nil),
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
			object: (*domain.DeviceAuth)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, prepareDeviceAuthQuery, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
