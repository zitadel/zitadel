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
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	expectedDeviceAuthQueryC = `SELECT` +
		` projections.device_auth_requests2.client_id,` +
		` projections.device_auth_requests2.device_code,` +
		` projections.device_auth_requests2.user_code,` +
		` projections.device_auth_requests2.scopes,` +
		` projections.device_auth_requests2.audience,` +
		` projections.apps7.name,` +
		` projections.projects4.name` +
		` FROM projections.device_auth_requests2` +
		` LEFT JOIN projections.apps7_oidc_configs` +
		` ON projections.device_auth_requests2.client_id = projections.apps7_oidc_configs.client_id` +
		` AND projections.device_auth_requests2.instance_id = projections.apps7_oidc_configs.instance_id` +
		` LEFT JOIN projections.apps7 ON projections.apps7_oidc_configs.app_id = projections.apps7.id` +
		` AND projections.apps7_oidc_configs.instance_id = projections.apps7.instance_id` +
		` LEFT JOIN projections.projects4 ON projections.apps7.project_id = projections.projects4.id` +
		` AND projections.apps7.instance_id = projections.projects4.instance_id`
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
		"appName",
		"projectName",
	}
	expectedDeviceAuth = domain.NewAuthRequestDevice(
		"client-id",
		"device1",
		"user-code",
		[]string{"a", "b", "c"},
		[]string{"projectID", "clientID"},
		"appName",
		"projectName",
	)
)

func TestQueries_DeviceAuthRequestByUserCode(t *testing.T) {
	permissionDenied := zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied")

	tests := []struct {
		name            string
		checkPermission bool
		permissionErr   error
		wantQuery       bool
		wantErr         error
	}{
		{
			name:      "without permission check",
			wantQuery: true,
		},
		{
			name:            "permission check passes",
			checkPermission: true,
			wantQuery:       true,
		},
		{
			name:            "permission denied, no query",
			checkPermission: true,
			permissionErr:   permissionDenied,
			wantErr:         permissionDenied,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, mock, err := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
			if err != nil {
				t.Fatalf("failed to build mock client: %v", err)
			}
			defer client.Close()

			if tt.wantQuery {
				mock.ExpectQuery(expectedDeviceAuthWhereUserCodeQuery).WillReturnRows(
					mock.NewRows(deviceAuthSelectColumns).AddRow(expectedDeviceAuthValues...),
				)
			}
			var permissionChecked bool
			q := Queries{
				client: &database.DB{DB: client},
				checkPermission: func(_ context.Context, permission, orgID, resourceID string) error {
					permissionChecked = true
					assert.Equal(t, domain.PermissionSessionRead, permission)
					assert.Empty(t, resourceID)
					return tt.permissionErr
				},
			}
			got, err := q.DeviceAuthRequestByUserCode(context.TODO(), "789", tt.checkPermission)
			assert.Equal(t, tt.checkPermission, permissionChecked)
			require.NoError(t, mock.ExpectationsWereMet())
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, expectedDeviceAuth, got)
		})
	}
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
			assertPrepare(t, prepareDeviceAuthQuery, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
