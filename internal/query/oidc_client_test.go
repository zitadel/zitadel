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
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed testdata/oidc_client_jwt.json
	testdataOidcClientJWT string
	//go:embed testdata/oidc_client_public.json
	testdataOidcClientPublic string
	//go:embed testdata/oidc_client_secret.json
	testdataOidcClientSecret string
	//go:embed testdata/oidc_client_no_settings.json
	testdataOidcClientNoSettings string
)

func TestQueries_GetOIDCClientByID(t *testing.T) {
	expQuery := regexp.QuoteMeta(oidcClientQuery)
	cols := []string{"client"}
	pubkey := `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2ufAL1b72bIy1ar+Ws6b
GohJJQFB7dfRapDqeqM8Ukp6CVdPzq/pOz1viAq50yzWZJryF+2wshFAKGF9A2/B
2Yf9bJXPZ/KbkFrYT3NTvYDkvlaSTl9mMnzrU29s48F1PTWKfB+C3aMsOEG1BufV
s63qF4nrEPjSbhljIco9FZq4XppIzhMQ0fDdA/+XygCJqvuaL0LibM1KrlUdnu71
YekhSJjEPnvOisXIk4IXywoGIOwtjxkDvNItQvaMVldr4/kb6uvbgdWwq5EwBZXq
low2kyJov38V4Uk2I8kuXpLcnrpw5Tio2ooiUE27b0vHZqBKOei9Uo88qCrn3EKx
6QIDAQAB
-----END RSA PUBLIC KEY-----
`

	tests := []struct {
		name    string
		mock    sqlExpectation
		want    *OIDCClient
		wantErr error
	}{
		{
			name:    "no rows",
			mock:    mockQueryErr(expQuery, sql.ErrNoRows, "instanceID", "clientID", true),
			wantErr: zerrors.ThrowNotFound(sql.ErrNoRows, "QUERY-wu6Ee", "Errors.App.NotFound"),
		},
		{
			name:    "internal error",
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, "instanceID", "clientID", true),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-ieR7R", "Errors.Internal"),
		},
		{
			name: "jwt client",
			mock: mockQuery(expQuery, cols, []driver.Value{testdataOidcClientJWT}, "instanceID", "clientID", true),
			want: &OIDCClient{
				InstanceID:               "230690539048009730",
				AppID:                    "236647088211886082",
				State:                    domain.AppStateActive,
				ClientID:                 "236647088211951618@tests",
				HashedSecret:             "",
				RedirectURIs:             []string{"http://localhost:9999/auth/callback"},
				ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode, domain.OIDCGrantTypeRefreshToken},
				ApplicationType:          domain.OIDCApplicationTypeWeb,
				AuthMethodType:           domain.OIDCAuthMethodTypePrivateKeyJWT,
				PostLogoutRedirectURIs:   []string{"https://example.com/logout"},
				IsDevMode:                true,
				AccessTokenType:          domain.OIDCTokenTypeJWT,
				AccessTokenRoleAssertion: true,
				IDTokenRoleAssertion:     true,
				IDTokenUserinfoAssertion: true,
				ClockSkew:                1000000000,
				AdditionalOrigins:        []string{"https://example.com"},
				ProjectID:                "236645808328409090",
				ProjectRoleAssertion:     true,
				PublicKeys:               map[string][]byte{"236647201860747266": []byte(pubkey)},
				ProjectRoleKeys:          []string{"role1", "role2"},
				Settings: &OIDCSettings{
					AccessTokenLifetime: 43200000000000,
					IdTokenLifetime:     43200000000000,
				},
			},
		},
		{
			name: "public client",
			mock: mockQuery(expQuery, cols, []driver.Value{testdataOidcClientPublic}, "instanceID", "clientID", true),
			want: &OIDCClient{
				InstanceID:               "230690539048009730",
				AppID:                    "236646457053020162",
				State:                    domain.AppStateActive,
				ClientID:                 "236646457053085698@tests",
				HashedSecret:             "",
				RedirectURIs:             []string{"http://localhost:9999/auth/callback"},
				ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				ApplicationType:          domain.OIDCApplicationTypeWeb,
				AuthMethodType:           domain.OIDCAuthMethodTypeNone,
				PostLogoutRedirectURIs:   nil,
				IsDevMode:                true,
				AccessTokenType:          domain.OIDCTokenTypeBearer,
				AccessTokenRoleAssertion: false,
				IDTokenRoleAssertion:     false,
				IDTokenUserinfoAssertion: false,
				ClockSkew:                0,
				AdditionalOrigins:        nil,
				PublicKeys:               nil,
				ProjectID:                "236645808328409090",
				ProjectRoleAssertion:     true,
				ProjectRoleKeys:          []string{"role1", "role2"},
				Settings: &OIDCSettings{
					AccessTokenLifetime: 43200000000000,
					IdTokenLifetime:     43200000000000,
				},
			},
		},
		{
			name: "secret client",
			mock: mockQuery(expQuery, cols, []driver.Value{testdataOidcClientSecret}, "instanceID", "clientID", true),
			want: &OIDCClient{
				InstanceID:               "230690539048009730",
				AppID:                    "236646858984783874",
				State:                    domain.AppStateActive,
				ClientID:                 "236646858984849410@tests",
				HashedSecret:             "$2a$14$OzZ0XEZZEtD13py/EPba2evsS6WcKZ5orVMj9pWHEGEHmLu2h3PFq",
				RedirectURIs:             []string{"http://localhost:9999/auth/callback"},
				ResponseTypes:            []domain.OIDCResponseType{0},
				GrantTypes:               []domain.OIDCGrantType{0},
				ApplicationType:          domain.OIDCApplicationTypeWeb,
				AuthMethodType:           domain.OIDCAuthMethodTypeBasic,
				PostLogoutRedirectURIs:   nil,
				IsDevMode:                true,
				AccessTokenType:          domain.OIDCTokenTypeBearer,
				AccessTokenRoleAssertion: false,
				IDTokenRoleAssertion:     false,
				IDTokenUserinfoAssertion: false,
				ClockSkew:                0,
				AdditionalOrigins:        nil,
				PublicKeys:               nil,
				ProjectID:                "236645808328409090",
				ProjectRoleAssertion:     false,
				ProjectRoleKeys:          []string{"role1", "role2"},
				Settings: &OIDCSettings{
					AccessTokenLifetime: 43200000000000,
					IdTokenLifetime:     43200000000000,
				},
			},
		},
		{
			name: "no oidc settings",
			mock: mockQuery(expQuery, cols, []driver.Value{testdataOidcClientNoSettings}, "instanceID", "clientID", true),
			want: &OIDCClient{
				InstanceID:   "239520764275982338",
				AppID:        "239520764276441090",
				State:        domain.AppStateActive,
				ClientID:     "239520764779364354@zitadel",
				HashedSecret: "",
				RedirectURIs: []string{
					"http://test2-qucuh5.localhost:9000/ui/console/auth/callback",
					"http://test.localhost.com:9000/ui/console/auth/callback"},
				ResponseTypes:   []domain.OIDCResponseType{0},
				GrantTypes:      []domain.OIDCGrantType{0},
				ApplicationType: domain.OIDCApplicationTypeUserAgent,
				AuthMethodType:  domain.OIDCAuthMethodTypeNone,
				PostLogoutRedirectURIs: []string{
					"http://test2-qucuh5.localhost:9000/ui/console/signedout",
					"http://test.localhost.com:9000/ui/console/signedout",
				},
				IsDevMode:                true,
				AccessTokenType:          domain.OIDCTokenTypeBearer,
				AccessTokenRoleAssertion: false,
				IDTokenRoleAssertion:     false,
				IDTokenUserinfoAssertion: false,
				ClockSkew:                0,
				AdditionalOrigins:        nil,
				PublicKeys:               nil,
				ProjectID:                "239520764276178946",
				ProjectRoleAssertion:     false,
				ProjectRoleKeys:          nil,
				Settings:                 nil,
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
				ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")
				got, err := q.GetOIDCClientByID(ctx, "clientID", true)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
