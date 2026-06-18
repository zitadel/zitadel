package command

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	oidc_pkg "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	providers "github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AllIDPWriteModel(t *testing.T) {
	type args struct {
		resourceOwner string
		instanceBool  bool
		id            string
		idpType       domain.IDPType
	}
	type res struct {
		writeModelType     interface{}
		samlWriteModelType interface{}
		err                error
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "writemodel instance oidc",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeOIDC,
			},
			res: res{
				writeModelType: &InstanceOIDCIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance jwt",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeJWT,
			},
			res: res{
				writeModelType: &InstanceJWTIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance oauth",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeOAuth,
			},
			res: res{
				writeModelType: &InstanceOAuthIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance ldap",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeLDAP,
			},
			res: res{
				writeModelType: &InstanceLDAPIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance azureAD",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeAzureAD,
			},
			res: res{
				writeModelType: &InstanceAzureADIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance github",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitHub,
			},
			res: res{
				writeModelType: &InstanceGitHubIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance github enterprise",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitHubEnterprise,
			},
			res: res{
				writeModelType: &InstanceGitHubEnterpriseIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance gitlab",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitLab,
			},
			res: res{
				writeModelType: &InstanceGitLabIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance gitlab self hosted",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGitLabSelfHosted,
			},
			res: res{
				writeModelType: &InstanceGitLabSelfHostedIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance google",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeGoogle,
			},
			res: res{
				writeModelType: &InstanceGoogleIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel instance saml",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeSAML,
			},
			res: res{
				samlWriteModelType: &InstanceSAMLIDPWriteModel{},
				err:                nil,
			},
		},
		{
			name: "writemodel instance unspecified",
			args: args{
				resourceOwner: "owner",
				instanceBool:  true,
				id:            "id",
				idpType:       domain.IDPTypeUnspecified,
			},
			res: res{
				err: zerrors.ThrowInternal(nil, "COMMAND-xw921211", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			name: "writemodel org oidc",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeOIDC,
			},
			res: res{
				writeModelType: &OrgOIDCIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org jwt",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeJWT,
			},
			res: res{
				writeModelType: &OrgJWTIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org oauth",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeOAuth,
			},
			res: res{
				writeModelType: &OrgOAuthIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org ldap",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeLDAP,
			},
			res: res{
				writeModelType: &OrgLDAPIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org azureAD",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeAzureAD,
			},
			res: res{
				writeModelType: &OrgAzureADIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org github",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitHub,
			},
			res: res{
				writeModelType: &OrgGitHubIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org github enterprise",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitHubEnterprise,
			},
			res: res{
				writeModelType: &OrgGitHubEnterpriseIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org gitlab",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitLab,
			},
			res: res{
				writeModelType: &OrgGitLabIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org gitlab self hosted",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGitLabSelfHosted,
			},
			res: res{
				writeModelType: &OrgGitLabSelfHostedIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org google",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeGoogle,
			},
			res: res{
				writeModelType: &OrgGoogleIDPWriteModel{},
				err:            nil,
			},
		},
		{
			name: "writemodel org saml",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeSAML,
			},
			res: res{
				samlWriteModelType: &OrgSAMLIDPWriteModel{},
				err:                nil,
			},
		},
		{
			name: "writemodel org unspecified",
			args: args{
				resourceOwner: "owner",
				instanceBool:  false,
				id:            "id",
				idpType:       domain.IDPTypeUnspecified,
			},
			res: res{
				err: zerrors.ThrowInternal(nil, "COMMAND-xw921111", "Errors.IDPConfig.NotExisting"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm, err := NewAllIDPWriteModel(tt.args.resourceOwner, tt.args.instanceBool, tt.args.id, tt.args.idpType)
			require.ErrorIs(t, err, tt.res.err)
			if wm != nil {
				if tt.res.writeModelType != nil {
					assert.IsType(t, tt.res.writeModelType, wm.model)
				}
				if tt.res.samlWriteModelType != nil {
					assert.IsType(t, tt.res.samlWriteModelType, wm.samlModel)
				}
			}
		})
	}
}

func TestOAuthIDPWriteModel_ToProvider_WithPKCE(t *testing.T) {
	wm := &OAuthIDPWriteModel{
		Name:                  "OAuth",
		ClientID:              "clientID",
		ClientSecret:          &crypto.CryptoValue{Algorithm: "plain", Crypted: []byte("clientSecret"), KeyID: "keyID"},
		AuthorizationEndpoint: "https://idp.example.com/authorize",
		TokenEndpoint:         "https://idp.example.com/token",
		UserEndpoint:          "https://idp.example.com/user",
		Scopes:                []string{"user"},
		IDAttribute:           "id",
		UsePKCE:               true,
	}

	provider, err := wm.ToProvider("https://zitadel.example.com/idps/callback", plainTextEncryption{}, http.DefaultClient)
	require.NoError(t, err)

	assertProviderUsesPKCE(t, provider)
}

func TestOIDCIDPWriteModel_ToProvider_WithPKCE(t *testing.T) {
	var (
		issuer             string
		discoveryRequested atomic.Bool
	)
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != oidc_pkg.DiscoveryEndpoint {
			http.NotFound(w, r)
			return
		}
		discoveryRequested.Store(true)
		w.Header().Set("Content-Type", "application/json")
		// assert, not require: FailNow must not be called from the server's handler goroutine
		assert.NoError(t, json.NewEncoder(w).Encode(&oidc_pkg.DiscoveryConfiguration{
			Issuer:                issuer,
			AuthorizationEndpoint: issuer + "/authorize",
			TokenEndpoint:         issuer + "/token",
			UserinfoEndpoint:      issuer + "/userinfo",
		}))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)
	issuer = server.URL

	wm := &OIDCIDPWriteModel{
		Name:         "OIDC",
		Issuer:       issuer,
		ClientID:     "clientID",
		ClientSecret: &crypto.CryptoValue{Algorithm: "plain", Crypted: []byte("clientSecret"), KeyID: "keyID"},
		Scopes:       []string{"openid"},
		UsePKCE:      true,
	}

	provider, err := wm.ToProvider("https://zitadel.example.com/idps/callback", plainTextEncryption{}, http.DefaultClient)
	require.NoError(t, err)
	require.True(t, discoveryRequested.Load(), "expected OIDC discovery to be requested")

	assertProviderUsesPKCE(t, provider)
}

func assertProviderUsesPKCE(t *testing.T, provider providers.Provider) {
	t.Helper()

	session, err := provider.BeginAuth(context.Background(), "state")
	require.NoError(t, err)
	auth, err := session.GetAuth(context.Background())
	require.NoError(t, err)
	redirect, ok := auth.(*providers.RedirectAuth)
	require.True(t, ok)
	authURL, err := url.Parse(redirect.RedirectURL)
	require.NoError(t, err)

	query := authURL.Query()
	assert.NotEmpty(t, query.Get("code_challenge"))
	assert.Equal(t, "S256", query.Get("code_challenge_method"))
	assert.NotEmpty(t, session.PersistentParameters()[oauth.CodeVerifier])
}

type plainTextEncryption struct{}

func (plainTextEncryption) Algorithm() string                               { return "plain" }
func (plainTextEncryption) EncryptionKeyID() string                         { return "keyID" }
func (plainTextEncryption) DecryptionKeyIDs() []string                      { return []string{"keyID"} }
func (plainTextEncryption) Encrypt(value []byte) ([]byte, error)            { return value, nil }
func (plainTextEncryption) Decrypt(hashed []byte, _ string) ([]byte, error) { return hashed, nil }
func (plainTextEncryption) DecryptString(hashed []byte, _ string) (string, error) {
	return string(hashed), nil
}
