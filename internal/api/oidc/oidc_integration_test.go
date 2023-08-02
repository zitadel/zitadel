//go:build integration

package oidc_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX      context.Context
	CTXLOGIN context.Context
	Tester   *integration.Tester
	User     *user.AddHumanUserResponse
)

const (
	redirectURI          = "oidcintegrationtest://callback"
	redirectURIImplicit  = "http://localhost:9999/callback"
	logoutRedirectURI    = "oidcintegrationtest://logged-out"
	zitadelAudienceScope = domain.ProjectIDScope + domain.ProjectIDScopeZITADEL + domain.AudSuffix
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, _ = Tester.WithAuthorization(ctx, integration.OrgOwner), errCtx
		User = Tester.CreateHumanUser(CTX)
		Tester.SetUserPassword(CTX, User.GetUserId(), integration.UserPassword)
		Tester.RegisterUserPasskey(CTX, User.GetUserId())
		CTXLOGIN, _ = Tester.WithAuthorization(ctx, integration.Login), errCtx
		return m.Run()
	}())
}

func Test_ZITADEL_API_missing_audience_scope(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID)
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
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_missing_authentication(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	createResp, err := Tester.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Domain: Tester.Config.ExternalDomain,
	})
	require.NoError(t, err)
	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    createResp.GetSessionId(),
				SessionToken: createResp.GetSessionToken(),
			},
		},
	})
	require.NoError(t, err)

	// code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_missing_mfa(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasswordSession(t, CTX, User.GetUserId(), integration.UserPassword)
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
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPassword, startTime, changeTime)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_success(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
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
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, armPasskey, startTime, changeTime)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())
}

func Test_ZITADEL_API_inactive_access_token(t *testing.T) {
	clientID := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
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

	// make sure token works
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())

	// refresh token
	newTokens, err := refreshTokens(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)

	// use invalidated token
	ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err = Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_terminated_session(t *testing.T) {
	clientID := createClient(t)
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
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

	// make sure token works
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())

	// refresh token
	postLogoutRedirect, err := rp.EndSession(provider, tokens.IDToken, logoutRedirectURI, "state")
	require.NoError(t, err)
	assert.Equal(t, logoutRedirectURI+"?state=state", postLogoutRedirect.String())

	// use token from terminated session
	ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err = Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func createClient(t testing.TB) string {
	project, err := Tester.CreateProject(CTX)
	require.NoError(t, err)
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId())
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

func assertOIDCTime(t *testing.T, actual oidc.Time, expected time.Time) {
	assertOIDCTimeRange(t, actual, expected, expected)
}

func assertOIDCTimeRange(t *testing.T, actual oidc.Time, expectedStart, expectedEnd time.Time) {
	assert.WithinRange(t, actual.AsTime(), expectedStart.Add(-1*time.Second), expectedEnd.Add(1*time.Second))
}
