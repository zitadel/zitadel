//go:build integration

package oidc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/text/language"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func TestOPStorage_SetUserinfoFromToken(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	// test actual userinfo
	provider, err := Tester.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)
	userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
	require.NoError(t, err)
	assertUserinfo(t, userinfo)
}

func TestServer_Introspect(t *testing.T) {
	project, err := Tester.CreateProject(CTX)
	require.NoError(t, err)
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId())
	require.NoError(t, err)
	api, err := Tester.CreateAPIClient(CTX, project.GetId())
	require.NoError(t, err)
	keyResp, err := Tester.Client.Mgmt.AddAppKey(CTX, &management.AddAppKeyRequest{
		ProjectId:      project.GetId(),
		AppId:          api.GetAppId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: nil,
	})
	require.NoError(t, err)
	resourceServer, err := Tester.CreateResourceServer(CTX, keyResp.GetKeyDetails())
	require.NoError(t, err)

	scope := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess, oidc_api.ScopeResourceOwner}
	authRequestID := createAuthRequest(t, app.GetClientId(), redirectURI, scope...)
	sessionID, sessionToken, startTime, changeTime := Tester.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, app.GetClientId(), code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	// test actual introspection
	introspection, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServer, tokens.AccessToken)
	require.NoError(t, err)
	assertIntrospection(t, introspection,
		Tester.OIDCIssuer(), app.GetClientId(),
		scope, []string{app.GetClientId(), api.GetClientId(), project.GetId()},
		tokens.Expiry, tokens.Expiry.Add(-12*time.Hour))
}

func assertUserinfo(t *testing.T, userinfo *oidc.UserInfo) {
	assert.Equal(t, User.GetUserId(), userinfo.Subject)
	assert.Equal(t, "Mickey", userinfo.GivenName)
	assert.Equal(t, "Mouse", userinfo.FamilyName)
	assert.Equal(t, "Mickey Mouse", userinfo.Name)
	assert.NotEmpty(t, userinfo.PreferredUsername)
	assert.Equal(t, userinfo.PreferredUsername, userinfo.Email)
	assert.False(t, bool(userinfo.EmailVerified))
	assertOIDCTime(t, userinfo.UpdatedAt, User.GetDetails().GetChangeDate().AsTime())
}

func assertIntrospection(
	t *testing.T,
	introspection *oidc.IntrospectionResponse,
	issuer, clientID string,
	scope, audience []string,
	expiration, creation time.Time,
) {
	assert.True(t, introspection.Active)
	assert.Equal(t, scope, []string(introspection.Scope))
	assert.Equal(t, clientID, introspection.ClientID)
	assert.Equal(t, oidc.BearerToken, introspection.TokenType)
	assertOIDCTime(t, introspection.Expiration, expiration)
	assertOIDCTime(t, introspection.IssuedAt, creation)
	assertOIDCTime(t, introspection.NotBefore, creation)
	assert.Equal(t, User.GetUserId(), introspection.Subject)
	assert.ElementsMatch(t, audience, introspection.Audience)
	assert.Equal(t, issuer, introspection.Issuer)
	assert.NotEmpty(t, introspection.JWTID)
	assert.NotEmpty(t, introspection.Username)
	assert.Equal(t, introspection.Username, introspection.PreferredUsername)
	assert.Equal(t, "Mickey", introspection.GivenName)
	assert.Equal(t, "Mouse", introspection.FamilyName)
	assert.Equal(t, "Mickey Mouse", introspection.Name)
	assert.Equal(t, oidc.Gender("male"), introspection.Gender)
	assert.Equal(t, oidc.NewLocale(language.Dutch), introspection.Locale)
	assert.Equal(t, introspection.Username, introspection.Email)
	assert.False(t, bool(introspection.EmailVerified))
	assertOIDCTime(t, introspection.UpdatedAt, User.GetDetails().GetChangeDate().AsTime())

	require.NotNil(t, introspection.Claims)
	assert.Equal(t, User.Details.ResourceOwner, introspection.Claims[oidc_api.ClaimResourceOwner+"id"])
	assert.NotEmpty(t, introspection.Claims[oidc_api.ClaimResourceOwner+"name"])
	assert.NotEmpty(t, introspection.Claims[oidc_api.ClaimResourceOwner+"primary_domain"])
}

// TestServer_VerifyClient tests verification by running code flow tests
// with clients that have different authentication methods.
func TestServer_VerifyClient(t *testing.T) {
	sessionID, sessionToken, startTime, changeTime := Tester.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	project, err := Tester.CreateProject(CTX)
	require.NoError(t, err)

	inactiveClient, err := Tester.CreateOIDCInactivateClient(CTX, redirectURI, logoutRedirectURI, project.GetId())
	require.NoError(t, err)
	nativeClient, err := Tester.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId())
	require.NoError(t, err)
	basicWebClient, err := Tester.CreateOIDCWebClientBasic(CTX, redirectURI, logoutRedirectURI, project.GetId())
	require.NoError(t, err)
	jwtWebClient, keyData, err := Tester.CreateOIDCWebClientJWT(CTX, redirectURI, logoutRedirectURI, project.GetId())
	require.NoError(t, err)

	type clientDetails struct {
		authReqClientID string
		clientID        string
		clientSecret    string
		keyData         []byte
	}
	tests := []struct {
		name    string
		client  clientDetails
		wantErr bool
	}{
		{
			name: "empty client ID error",
			client: clientDetails{
				authReqClientID: nativeClient.GetClientId(),
			},
			wantErr: true,
		},
		{
			name: "client not found error",
			client: clientDetails{
				authReqClientID: nativeClient.GetClientId(),
				clientID:        "foo",
			},
			wantErr: true,
		},
		{
			name: "client inactive error",
			client: clientDetails{
				authReqClientID: nativeClient.GetClientId(),
				clientID:        inactiveClient.GetClientId(),
			},
			wantErr: true,
		},
		{
			name: "native client success",
			client: clientDetails{
				authReqClientID: nativeClient.GetClientId(),
				clientID:        nativeClient.GetClientId(),
			},
		},
		{
			name: "web client basic secret empty error",
			client: clientDetails{
				authReqClientID: basicWebClient.GetClientId(),
				clientID:        basicWebClient.GetClientId(),
				clientSecret:    "",
			},
			wantErr: true,
		},
		{
			name: "web client basic secret invalid error",
			client: clientDetails{
				authReqClientID: basicWebClient.GetClientId(),
				clientID:        basicWebClient.GetClientId(),
				clientSecret:    "wrong",
			},
			wantErr: true,
		},
		{
			name: "web client basic secret success",
			client: clientDetails{
				authReqClientID: basicWebClient.GetClientId(),
				clientID:        basicWebClient.GetClientId(),
				clientSecret:    basicWebClient.GetClientSecret(),
			},
		},
		{
			name: "web client JWT profile empty assertion error",
			client: clientDetails{
				authReqClientID: jwtWebClient.GetClientId(),
				clientID:        jwtWebClient.GetClientId(),
			},
			wantErr: true,
		},
		{
			name: "web client JWT profile invalid assertion error",
			client: clientDetails{
				authReqClientID: jwtWebClient.GetClientId(),
				clientID:        jwtWebClient.GetClientId(),
				keyData:         createInvalidKeyData(t, jwtWebClient),
			},
			wantErr: true,
		},
		{
			name: "web client JWT profile success",
			client: clientDetails{
				authReqClientID: jwtWebClient.GetClientId(),
				clientID:        jwtWebClient.GetClientId(),
				keyData:         keyData,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("\n\n%s\n\n", tt.client.keyData)

			authRequestID, err := Tester.CreateOIDCAuthRequest(CTX, tt.client.authReqClientID, Tester.Users[integration.FirstInstanceUsersKey][integration.Login].ID, redirectURI, oidc.ScopeOpenID)
			require.NoError(t, err)
			linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
				AuthRequestId: authRequestID,
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    sessionID,
						SessionToken: sessionToken,
					},
				},
			})
			require.NoError(t, err)

			// use a new RP so we can inject different credentials
			var options []rp.Option
			if tt.client.keyData != nil {
				options = append(options, rp.WithJWTProfile(rp.SignerFromKeyFile(tt.client.keyData)))
			}
			provider, err := rp.NewRelyingPartyOIDC(CTX, Tester.OIDCIssuer(), tt.client.clientID, tt.client.clientSecret, redirectURI, []string{oidc.ScopeOpenID}, options...)
			require.NoError(t, err)

			// test code exchange
			code := assertCodeResponse(t, linkResp.GetCallbackUrl())
			codeOpts := codeExchangeOptions(t, provider)
			tokens, err := rp.CodeExchange[*oidc.IDTokenClaims](context.Background(), code, provider, codeOpts...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assertTokens(t, tokens, false)
			assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)
		})
	}
}

func codeExchangeOptions(t testing.TB, provider rp.RelyingParty) []rp.CodeExchangeOpt {
	codeOpts := []rp.CodeExchangeOpt{rp.WithCodeVerifier(integration.CodeVerifier)}
	if signer := provider.Signer(); signer != nil {
		assertion, err := client.SignedJWTProfileAssertion(provider.OAuthConfig().ClientID, []string{provider.Issuer()}, time.Hour, provider.Signer())
		require.NoError(t, err)
		codeOpts = append(codeOpts, rp.WithClientAssertionJWT(assertion))
	}
	return codeOpts
}

func createInvalidKeyData(t testing.TB, client *management.AddOIDCAppResponse) []byte {
	key := domain.ApplicationKey{
		Type:          domain.AuthNKeyTypeJSON,
		KeyID:         "1",
		PrivateKey:    []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAxHd087RoEm9ywVWZ/H+tDWxQsmVvhfRz4jAq/RfU+OWXNH4J\njMMSHdFs0Q+WP98nNXRyc7fgbMb8NdmlB2yD4qLYapN5SDaBc5dh/3EnyFt53oSs\njTlKnQUPAeJr2qh/NY046CfyUyQMM4JR5OiQFo4TssfWnqdcgamGt0AEnk2lvbMZ\nKQdAqNS9lDzYbjMGavEQPTZE35mFXFQXjaooZXq+TIa7hbaq7/idH7cHNbLcPLgj\nfPQA8q+DYvnvhXlmq0LPQZH3Oiixf+SF2vRwrBzT2mqGD2OiOkUmhuPwyqEiiBHt\nfxklRtRU6WfLa1Gcb1PsV0uoBGpV3KybIl/GlwIDAQABAoIBAEQjDduLgOCL6Gem\n0X3hpdnW6/HC/jed/Sa//9jBECq2LYeWAqff64ON40hqOHi0YvvGA/+gEOSI6mWe\nsv5tIxxRz+6+cLybsq+tG96kluCE4TJMHy/nY7orS/YiWbd+4odnEApr+D3fbZ/b\nnZ1fDsHTyn8hkYx6jLmnWsJpIHDp7zxD76y7k2Bbg6DZrCGiVxngiLJk23dvz79W\np03lHLM7XE92aFwXQmhfxHGxrbuoB/9eY4ai5IHp36H4fw0vL6NXdNQAo/bhe0p9\nAYB7y0ZumF8Hg0Z/BmMeEzLy6HrYB+VE8cO93pNjhSyH+p2yDB/BlUyTiRLQAoM0\nVTmOZXECgYEA7NGlzpKNhyQEJihVqt0MW0LhKIO/xbBn+XgYfX6GpqPa/ucnMx5/\nVezpl3gK8IU4wPUhAyXXAHJiqNBcEeyxrw0MXLujDVMJgYaLysCLJdvMVgoY08mS\nK5IQivpbozpf4+0y3mOnA+Sy1kbfxv2X8xiWLODRQW3f3q/xoklwOR8CgYEA1GEe\nfaibOFTQAYcIVj77KXtBfYZsX3EGAyfAN9O7cKHq5oaxVstwnF47WxpuVtoKZxCZ\nbNm9D5WvQ9b+Ztpioe42tzwE7Bff/Osj868GcDdRPK7nFlh9N2yVn/D514dOYVwR\n4MBr1KrJzgRWt4QqS4H+to1GzudDTSNlG7gnK4kCgYBUi6AbOHzoYzZL/RhgcJwp\ntJ23nhmH1Su5h2OO4e3mbhcP66w19sxU+8iFN+kH5zfUw26utgKk+TE5vXExQQRK\nT2k7bg2PAzcgk80ybD0BHhA8I0yrx4m0nmfjhe/TPVLgh10iwgbtP+eM0i6v1vc5\nZWyvxu9N4ZEL6lpkqr0y1wKBgG/NAIQd8jhhTW7Aav8cAJQBsqQl038avJOEpYe+\nCnpsgoAAf/K0/f8TDCQVceh+t+MxtdK7fO9rWOxZjWsPo8Si5mLnUaAHoX4/OpnZ\nlYYVWMqdOEFnK+O1Yb7k2GFBdV2DXlX2dc1qavntBsls5ecB89id3pyk2aUN8Pf6\npYQhAoGAMGtrHFely9wyaxI0RTCyfmJbWZHGVGkv6ELK8wneJjdjl82XOBUGCg5q\naRCrTZ3dPitKwrUa6ibJCIFCIziiriBmjDvTHzkMvoJEap2TVxYNDR6IfINVsQ57\nlOsiC4A2uGq4Lbfld+gjoplJ5GX6qXtTgZ6m7eo0y7U6zm2tkN0=\n-----END RSA PRIVATE KEY-----\n"),
		ApplicationID: client.GetAppId(),
		ClientID:      client.GetClientId(),
	}
	data, err := key.Detail()
	require.NoError(t, err)
	return data
}

func TestServer_CreateAccessToken_ClientCredentials(t *testing.T) {
	clientID, clientSecret, err := Tester.CreateOIDCCredentialsClient(CTX)
	require.NoError(t, err)

	type clientDetails struct {
		clientID     string
		clientSecret string
		keyData      []byte
	}
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		wantErr      bool
	}{
		{
			name:         "missing client ID error",
			clientID:     "",
			clientSecret: clientSecret,
			wantErr:      true,
		},
		{
			name:         "client not found error",
			clientID:     "foo",
			clientSecret: clientSecret,
			wantErr:      true,
		},
		{
			name: "machine user without secret error",
			clientID: func() string {
				name := gofakeit.Username()
				_, err := Tester.Client.Mgmt.AddMachineUser(CTX, &management.AddMachineUserRequest{
					Name:            name,
					UserName:        name,
					AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
				})
				require.NoError(t, err)
				return name
			}(),
			clientSecret: clientSecret,
			wantErr:      true,
		},
		{
			name:         "wrong secret error",
			clientID:     clientID,
			clientSecret: "bar",
			wantErr:      true,
		},
		{
			name:         "success",
			clientID:     clientID,
			clientSecret: clientSecret,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := rp.NewRelyingPartyOIDC(CTX, Tester.OIDCIssuer(), tt.clientID, tt.clientSecret, redirectURI, []string{oidc.ScopeOpenID})
			require.NoError(t, err)
			tokens, err := rp.ClientCredentials(CTX, provider, nil)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, tokens)
			assert.NotEmpty(t, tokens.AccessToken)
		})
	}
}
