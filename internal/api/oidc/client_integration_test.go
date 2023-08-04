//go:build integration

package oidc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/client/rs"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
)

func TestOPStorage_SetUserinfoFromToken(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
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
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)
	userinfo, err := rp.Userinfo(tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
	require.NoError(t, err)
	assertUserinfo(t, userinfo)
}

func TestOPStorage_SetIntrospectionFromToken(t *testing.T) {
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
	resourceServer, err := Tester.CreateResourceServer(keyResp.GetKeyDetails())
	require.NoError(t, err)

	scope := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess}
	authRequestID := createAuthRequest(t, app.GetClientId(), redirectURI, scope...)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasskeySession(t, CTXLOGIN, User.GetUserId())
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
	introspection, err := rs.Introspect(context.Background(), resourceServer, tokens.AccessToken)
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
	assert.Equal(t, introspection.Username, introspection.Email)
	assert.False(t, bool(introspection.EmailVerified))
	assertOIDCTime(t, introspection.UpdatedAt, User.GetDetails().GetChangeDate().AsTime())
}
