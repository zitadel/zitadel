package query

import (
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
)

func TestQueries_GetIntrospectionClientByID(t *testing.T) {
	pubkeys := database.Map[[]byte]{
		"key1": {1, 2, 3},
		"key2": {4, 5, 6},
	}
	encPubkeys, err := pubkeys.Value()
	require.NoError(t, err)

	expQuery := regexp.QuoteMeta(introspectionClientByIDQuery)
	type args struct {
		clientID string
		getKeys  bool
	}
	tests := []struct {
		name    string
		args    args
		mock    sqlExpectation
		want    *IntrospectionClient
		wantErr error
	}{
		{
			name: "query error",
			args: args{
				clientID: "clientID",
				getKeys:  false,
			},
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, "instanceID", "clientID", false),
			wantErr: sql.ErrConnDone,
		},
		{
			name: "success, secret",
			args: args{
				clientID: "clientID",
				getKeys:  false,
			},
			mock: mockQuery(expQuery,
				[]string{"app_id", "client_id", "client_secret", "app_type", "project_id", "resource_owner", "project_role_assertion", "public_keys"},
				[]driver.Value{"appID", "clientID", "secret", "oidc", "projectID", "orgID", true, nil},
				"instanceID", "clientID", false),
			want: &IntrospectionClient{
				AppID:                "appID",
				ClientID:             "clientID",
				HashedSecret:         "secret",
				AppType:              AppTypeOIDC,
				ProjectID:            "projectID",
				ResourceOwner:        "orgID",
				ProjectRoleAssertion: true,
				PublicKeys:           nil,
			},
		},
		{
			name: "success, keys",
			args: args{
				clientID: "clientID",
				getKeys:  true,
			},
			mock: mockQuery(expQuery,
				[]string{"app_id", "client_id", "client_secret", "app_type", "project_id", "resource_owner", "project_role_assertion", "public_keys"},
				[]driver.Value{"appID", "clientID", "", "oidc", "projectID", "orgID", true, encPubkeys},
				"instanceID", "clientID", true),
			want: &IntrospectionClient{
				AppID:                "appID",
				ClientID:             "clientID",
				HashedSecret:         "",
				AppType:              AppTypeOIDC,
				ProjectID:            "projectID",
				ResourceOwner:        "orgID",
				ProjectRoleAssertion: true,
				PublicKeys:           pubkeys,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.mock, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB:       db,
						Database: &prepareDB{},
					},
				}
				ctx := authz.NewMockContext("instanceID", "orgID", "userID")
				got, err := q.GetIntrospectionClientByID(ctx, tt.args.clientID, tt.args.getKeys)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
