package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQueries_SamlRequestByID(t *testing.T) {
	expQuery := regexp.QuoteMeta(fmt.Sprintf(
		samlRequestByIDQuery,
		asOfSystemTime,
	))

	cols := []string{
		projection.SamlRequestColumnID,
		projection.SamlRequestColumnCreationDate,
		projection.SamlRequestColumnLoginClient,
		projection.SamlRequestColumnIssuer,
		projection.SamlRequestColumnACS,
		projection.SamlRequestColumnRelayState,
		projection.SamlRequestColumnBinding,
	}
	type args struct {
		shouldTriggerBulk bool
		id                string
		checkLoginClient  bool
	}
	tests := []struct {
		name            string
		args            args
		expect          sqlExpectation
		permissionCheck domain.PermissionCheck
		want            *SamlRequest
		wantErr         error
	}{
		{
			name: "success, all values",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"loginClient",
				"issuer",
				"acs",
				"relayState",
				"binding",
			}, "123", "instanceID"),
			want: &SamlRequest{
				ID:           "id",
				CreationDate: testNow,
				LoginClient:  "loginClient",
				Issuer:       "issuer",
				ACS:          "acs",
				RelayState:   "relayState",
				Binding:      "binding",
			},
		},
		{
			name: "no rows",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
			},
			expect:  mockQueryScanErr(expQuery, cols, nil, "123", "instanceID"),
			wantErr: zerrors.ThrowNotFound(sql.ErrNoRows, "QUERY-Thee9", "Errors.SamlRequest.NotExisting"),
		},
		{
			name: "query error",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
			},
			expect:  mockQueryErr(expQuery, sql.ErrConnDone, "123", "instanceID"),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-Ou8ue", "Errors.Internal"),
		},
		{
			name: "wrong login client/ not permitted",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"wrongLoginClient",
				"issuer",
				"acs",
				"relayState",
				"binding",
			}, "123", "instanceID"),
			permissionCheck: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return zerrors.ThrowPermissionDenied(nil, "id", "not permitted")
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "id", "not permitted"),
		},
		{
			name: "wrong login client / permitted",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"otherLoginClient",
				"issuer",
				"acs",
				"relayState",
				"binding",
			}, "123", "instanceID"),
			permissionCheck: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return nil
			},
			want: &SamlRequest{
				ID:           "id",
				CreationDate: testNow,
				LoginClient:  "otherLoginClient",
				Issuer:       "issuer",
				ACS:          "acs",
				RelayState:   "relayState",
				Binding:      "binding",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.expect, func(db *sql.DB) {
				q := &Queries{
					checkPermission: tt.permissionCheck,
					client: &database.DB{
						DB:       db,
						Database: &prepareDB{},
					},
				}
				ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")

				got, err := q.SamlRequestByID(ctx, tt.args.shouldTriggerBulk, tt.args.id, tt.args.checkLoginClient)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
