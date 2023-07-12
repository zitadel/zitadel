//go:build integration

package oidc_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/client/rs"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX      context.Context
	CTXLOGIN context.Context
	Tester   *integration.Tester
	User     *user.AddHumanUserResponse
)

const (
	redirectURI         = "oidcIntegrationTest://callback"
	redirectURIImplicit = "http://localhost:9999/callback"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, _ = Tester.WithAuthorization(ctx, integration.OrgOwner), errCtx
		User = Tester.CreateHumanUser(CTX)
		Tester.RegisterUserPasskey(CTX, User.GetUserId())
		CTXLOGIN, _ = Tester.WithAuthorization(ctx, integration.Login), errCtx
		return m.Run()
	}())
}

func createClient(t testing.TB) string {
	project, err := Tester.CreateProject(CTX)
	require.NoError(t, err)
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI, project.GetId())
	require.NoError(t, err)
	return app.GetClientId()
}

func createImplicitClient(t testing.TB) string {
	app, err := Tester.CreateOIDCImplicitFlowClient(CTX, redirectURIImplicit)
	require.NoError(t, err)
	return app.GetClientId()
}

func createAuthRequest(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequest(clientID, Tester.Users[integration.FirstInstanceUsersKey][integration.Login].ID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func createAuthRequestImplicit(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequestImplicit(clientID, Tester.Users[integration.FirstInstanceUsersKey][integration.Login].ID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func TestOPStorage_CreateAuthRequest(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, redirectURI)
	require.Contains(t, id, command.IDPrefixV2)
}

func TestOPStorage_CreateAccessToken_code(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI)
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

	// test code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, startTime, changeTime)

	// callback on a succeeded request must fail
	linkResp, err = Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.Error(t, err)

	// exchange with a used code must fail
	_, err = exchangeTokens(t, clientID, code)
	require.Error(t, err)
}

func TestOPStorage_CreateAccessToken_implicit(t *testing.T) {
	clientID := createImplicitClient(t)
	authRequestID := createAuthRequestImplicit(t, clientID, redirectURIImplicit)
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

	// test implicit callback
	callback, err := url.Parse(linkResp.GetCallbackUrl())
	require.NoError(t, err)
	values, err := url.ParseQuery(callback.Fragment)
	require.NoError(t, err)
	accessToken := values.Get("access_token")
	idToken := values.Get("id_token")
	refreshToken := values.Get("refresh_token")
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, idToken)
	assert.Empty(t, refreshToken)
	assert.NotEmpty(t, values.Get("expires_in"))
	assert.Equal(t, oidc.BearerToken, values.Get("token_type"))
	assert.Equal(t, "state", values.Get("state"))

	// check id_token / claims
	provider, err := Tester.CreateRelyingParty(clientID, redirectURIImplicit)
	require.NoError(t, err)
	claims, err := rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), accessToken, idToken, provider.IDTokenVerifier())
	require.NoError(t, err)
	assertIDTokenClaims(t, claims, startTime, changeTime)

	// callback on a succeeded request must fail
	linkResp, err = Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.Error(t, err)
}

func TestOPStorage_CreateAccessAndRefreshTokens_code(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
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

	// test code exchange (expect refresh token to be returned)
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, startTime, changeTime)
}

func TestOPStorage_CreateAccessAndRefreshTokens_refresh(t *testing.T) {
	clientID := createClient(t)
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
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
	assertIDTokenClaims(t, tokens.IDTokenClaims, startTime, changeTime)

	// test actual refresh grant
	newTokens, err := refreshTokens(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
	idToken, _ := newTokens.Extra("id_token").(string)
	assert.NotEmpty(t, idToken)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
	claims, err := rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), newTokens.AccessToken, idToken, provider.IDTokenVerifier())
	require.NoError(t, err)
	// auth time must still be the initial
	assertIDTokenClaims(t, claims, startTime, changeTime)

	// refresh with an old refresh_token must fail
	_, err = rp.RefreshAccessToken(provider, tokens.RefreshToken, "", "")
	require.Error(t, err)
}

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
	assertIDTokenClaims(t, tokens.IDTokenClaims, startTime, changeTime)

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
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI, project.GetId())
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
	assertIDTokenClaims(t, tokens.IDTokenClaims, startTime, changeTime)

	// test actual introspection
	introspection, err := rs.Introspect(context.Background(), resourceServer, tokens.AccessToken)
	require.NoError(t, err)
	assertIntrospection(t, introspection,
		Tester.OIDCIssuer(), app.GetClientId(),
		scope, []string{app.GetClientId(), api.GetClientId(), project.GetId()},
		tokens.Expiry, tokens.Expiry.Add(-12*time.Hour))
}

func exchangeTokens(t testing.TB, clientID, code string) (*oidc.Tokens[*oidc.IDTokenClaims], error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	codeVerifier := "codeVerifier"
	return rp.CodeExchange[*oidc.IDTokenClaims](context.Background(), code, provider, rp.WithCodeVerifier(codeVerifier))
}

func refreshTokens(t testing.TB, clientID, refreshToken string) (*oauth2.Token, error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	return rp.RefreshAccessToken(provider, refreshToken, "", "")
}

func assertCodeResponse(t *testing.T, callback string) string {
	callbackURL, err := url.Parse(callback)
	require.NoError(t, err)
	code := callbackURL.Query().Get("code")
	require.NotEmpty(t, code)
	assert.Equal(t, "state", callbackURL.Query().Get("state"))
	return code
}

func assertTokens(t *testing.T, tokens *oidc.Tokens[*oidc.IDTokenClaims], requireRefreshToken bool) {
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.IDToken)
	if requireRefreshToken {
		assert.NotEmpty(t, tokens.RefreshToken)
	} else {
		assert.Empty(t, tokens.RefreshToken)
	}
}

func assertIDTokenClaims(t *testing.T, claims *oidc.IDTokenClaims, sessionStart, sessionChange time.Time) {
	assert.Equal(t, User.GetUserId(), claims.Subject)
	assert.Equal(t, []string{oidc_api.UserPresence, oidc_api.MFA}, claims.AuthenticationMethodsReferences)
	assertOIDCTimeRange(t, claims.AuthTime, sessionStart, sessionChange)
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
	assert.Equal(t, audience, []string(introspection.Audience))
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

func assertOIDCTime(t *testing.T, actual oidc.Time, expected time.Time) {
	assertOIDCTimeRange(t, actual, expected, expected)
}

func assertOIDCTimeRange(t *testing.T, actual oidc.Time, expectedStart, expectedEnd time.Time) {
	assert.WithinRange(t, actual.AsTime(), expectedStart.Add(-1*time.Second), expectedEnd.Add(1*time.Second))
}
