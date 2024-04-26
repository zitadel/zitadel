package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQueries_DeviceAuthByDeviceCode(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	timestamp := time.Date(2015, 12, 15, 22, 13, 45, 0, time.UTC)
	tests := []struct {
		name       string
		eventstore func(t *testing.T) *eventstore.Eventstore
		want       *DeviceAuth
		wantErr    error
	}{
		{
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "not found",
			eventstore: expectEventstore(
				expectFilter(),
			),
			wantErr: zerrors.ThrowNotFound(nil, "QUERY-eeR0e", "Errors.DeviceAuth.NotExisting"),
		},
		{
			name: "ok, initiated",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						"client1", "device1", "user-code", timestamp, []string{"foo", "bar"},
						[]string{"projectID", "clientID"},
					)),
				),
			),
			want: &DeviceAuth{
				ClientID:   "client1",
				DeviceCode: "device1",
				UserCode:   "user-code",
				Expires:    timestamp,
				Scopes:     []string{"foo", "bar"},
				Audience:   []string{"projectID", "clientID"},
				State:      domain.DeviceAuthStateInitiated,
			},
		},
		{
			name: "ok, approved",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						"client1", "device1", "user-code", timestamp, []string{"foo", "bar"},
						[]string{"projectID", "clientID"},
					)),
					eventFromEventPusher(deviceauth.NewApprovedEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						"user1", []domain.UserAuthMethodType{domain.UserAuthMethodTypePasswordless},
						timestamp,
					)),
				),
			),
			want: &DeviceAuth{
				ClientID:        "client1",
				DeviceCode:      "device1",
				UserCode:        "user-code",
				Expires:         timestamp,
				Scopes:          []string{"foo", "bar"},
				Audience:        []string{"projectID", "clientID"},
				State:           domain.DeviceAuthStateApproved,
				Subject:         "user1",
				UserAuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePasswordless},
				AuthTime:        timestamp,
			},
		},
		{
			name: "ok, denied",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						"client1", "device1", "user-code", timestamp, []string{"foo", "bar"},
						[]string{"projectID", "clientID"},
					)),
					eventFromEventPusher(deviceauth.NewCanceledEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						domain.DeviceAuthCanceledDenied,
					)),
				),
			),
			want: &DeviceAuth{
				ClientID:   "client1",
				DeviceCode: "device1",
				UserCode:   "user-code",
				Expires:    timestamp,
				Scopes:     []string{"foo", "bar"},
				Audience:   []string{"projectID", "clientID"},
				State:      domain.DeviceAuthStateDenied,
			},
		},
		{
			name: "ok, expired",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(deviceauth.NewAddedEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						"client1", "device1", "user-code", timestamp, []string{"foo", "bar"},
						[]string{"projectID", "clientID"},
					)),
					eventFromEventPusher(deviceauth.NewCanceledEvent(
						ctx,
						deviceauth.NewAggregate("device1", "instance1"),
						domain.DeviceAuthCanceledExpired,
					)),
				),
			),
			want: &DeviceAuth{
				ClientID:   "client1",
				DeviceCode: "device1",
				UserCode:   "user-code",
				Expires:    timestamp,
				Scopes:     []string{"foo", "bar"},
				Audience:   []string{"projectID", "clientID"},
				State:      domain.DeviceAuthStateExpired,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				eventstore: tt.eventstore(t),
			}
			got, err := q.DeviceAuthByDeviceCode(ctx, "device1")
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

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

	mock.ExpectBegin()
	mock.ExpectQuery(expectedDeviceAuthWhereUserCodeQuery).WillReturnRows(
		mock.NewRows(deviceAuthSelectColumns).AddRow(expectedDeviceAuthValues...),
	)
	mock.ExpectCommit()
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
