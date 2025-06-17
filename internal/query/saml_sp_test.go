package query

import (
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQueries_ActiveSAMLServiceProviderByID(t *testing.T) {
	expQuery := regexp.QuoteMeta(samlSPQuery)
	cols := []string{
		"instance_id",
		"app_id",
		"state",
		"entity_id",
		"metadata",
		"metadata_url",
		"project_id",
		"project_role_assertion",
		"login_version",
		"login_base_uri",
	}

	tests := []struct {
		name    string
		mock    sqlExpectation
		want    *SAMLServiceProvider
		wantErr error
	}{
		{
			name:    "no rows",
			mock:    mockQueryErr(expQuery, sql.ErrNoRows, "instanceID", "entityID"),
			wantErr: zerrors.ThrowNotFound(sql.ErrNoRows, "QUERY-HeOcis2511", "Errors.App.NotFound"),
		},
		{
			name:    "internal error",
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, "instanceID", "entityID"),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-OyJx1Rp30z", "Errors.Internal"),
		},
		{
			name: "sp",
			mock: mockQuery(expQuery, cols, []driver.Value{
				"230690539048009730",
				"236647088211886082",
				domain.AppStateActive,
				"https://test.com/metadata",
				"metadata",
				"https://test.com/metadata",
				"236645808328409090",
				true,
				domain.LoginVersionUnspecified,
				"",
			}, "instanceID", "entityID"),
			want: &SAMLServiceProvider{
				InstanceID:           "230690539048009730",
				AppID:                "236647088211886082",
				State:                domain.AppStateActive,
				EntityID:             "https://test.com/metadata",
				Metadata:             []byte("metadata"),
				MetadataURL:          "https://test.com/metadata",
				ProjectID:            "236645808328409090",
				ProjectRoleAssertion: true,
			},
		},
		{
			name: "sp with loginversion",
			mock: mockQuery(expQuery, cols, []driver.Value{
				"230690539048009730",
				"236647088211886082",
				domain.AppStateActive,
				"https://test.com/metadata",
				"metadata",
				"https://test.com/metadata",
				"236645808328409090",
				true,
				domain.LoginVersion2,
				"https://test.com/login",
			}, "instanceID", "entityID"),
			want: &SAMLServiceProvider{
				InstanceID:           "230690539048009730",
				AppID:                "236647088211886082",
				State:                domain.AppStateActive,
				EntityID:             "https://test.com/metadata",
				Metadata:             []byte("metadata"),
				MetadataURL:          "https://test.com/metadata",
				ProjectID:            "236645808328409090",
				ProjectRoleAssertion: true,
				LoginVersion:         domain.LoginVersion2,
				LoginBaseURI: func() *url.URL {
					ret, _ := url.Parse("https://test.com/login")
					return ret
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.mock, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB: db,
					},
				}
				ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")
				got, err := q.ActiveSAMLServiceProviderByID(ctx, "entityID")
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
