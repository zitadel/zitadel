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
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX      context.Context
	CTXLOGIN context.Context
	CTXIAM   context.Context
	Tester   *integration.Tester
	User     *user.AddHumanUserResponse
)

const (
	redirectURI          = "https://callback"
	redirectURIImplicit  = "http://localhost:9999/callback"
	logoutRedirectURI    = "https://logged-out"
	zitadelAudienceScope = domain.ProjectIDScope + domain.ProjectIDScopeZITADEL + domain.AudSuffix
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(10 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX = Tester.WithAuthorization(ctx, integration.OrgOwner)
		User = Tester.CreateHumanUser(CTX)
		Tester.SetUserPassword(CTX, User.GetUserId(), integration.UserPassword, false)
		Tester.RegisterUserPasskey(CTX, User.GetUserId())
		CTXLOGIN = Tester.WithAuthorization(ctx, integration.Login)
		CTXIAM = Tester.WithAuthorization(ctx, integration.IAMOwner)
		return m.Run()
	}())
}

func Test_ZITADEL_API_missing_audience_scope(t *testing.T) {
	clientID, _ := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID)
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_missing_authentication(t *testing.T) {
	clientID, _ := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	createResp, err := Tester.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_missing_mfa_policy(t *testing.T) {
	clientID, _ := createClient(t)
	org := Tester.CreateOrganization(CTXIAM, fmt.Sprintf("ZITADEL_API_MISSING_MFA_%d", time.Now().UnixNano()), fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()))
	userID := org.CreatedAdmins[0].GetUserId()
	Tester.SetUserPassword(CTXIAM, userID, integration.UserPassword, false)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Tester.CreatePasswordSession(t, CTXLOGIN, userID, integration.UserPassword)
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertIDTokenClaims(t, tokens.IDTokenClaims, userID, armPassword, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	// pre check if request would succeed
	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, userID, myUserResp.GetUser().GetId())

	// require MFA
	ctxOrg := metadata.AppendToOutgoingContext(CTXIAM, "x-zitadel-orgid", org.GetOrganizationId())
	_, err = Tester.Client.Mgmt.AddCustomLoginPolicy(ctxOrg, &mgmt.AddCustomLoginPolicyRequest{
		ForceMfa: true,
	})
	require.NoError(t, err)

	// make sure policy is projected
	retryDuration := 5 * time.Second
	if ctxDeadline, ok := CTX.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, getErr := Tester.Client.Mgmt.GetLoginPolicy(ctxOrg, &mgmt.GetLoginPolicyRequest{})
		assert.NoError(ttt, getErr)
		assert.False(ttt, got.GetPolicy().IsDefault)

	}, retryDuration, time.Millisecond*100, "timeout waiting for login policy")

	// now it must fail
	myUserResp, err = Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_success(t *testing.T) {
	clientID, _ := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())
}

func Test_ZITADEL_API_glob_redirects(t *testing.T) {
	const redirectURI = "https://my-org-1yfnjl2xj-my-app.vercel.app/api/auth/callback/zitadel"
	clientID, _ := createClientWithOpts(t, clientOpts{
		redirectURI: "https://my-org-*-my-app.vercel.app/api/auth/callback/zitadel",
		logoutURI:   "https://my-org-*-my-app.vercel.app/",
		devMode:     true,
	})
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())
}

func Test_ZITADEL_API_inactive_access_token(t *testing.T) {
	clientID, _ := createClient(t)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

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
	clientID, _ := createClient(t)
	provider, err := Tester.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
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
	tokens, err := exchangeTokens(t, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	// make sure token works
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())

	// end session
	postLogoutRedirect, err := rp.EndSession(CTX, provider, tokens.IDToken, logoutRedirectURI, "state")
	require.NoError(t, err)
	assert.Equal(t, logoutRedirectURI+"?state=state", postLogoutRedirect.String())

	// use token from terminated session
	ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err = Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_terminated_session_user_disabled(t *testing.T) {
	clientID, _ := createClient(t)
	tests := []struct {
		name    string
		disable func(userID string) error
	}{
		{
			name: "deactivated",
			disable: func(userID string) error {
				_, err := Tester.Client.UserV2.DeactivateUser(CTX, &user.DeactivateUserRequest{UserId: userID})
				return err
			},
		},
		{
			name: "locked",
			disable: func(userID string) error {
				_, err := Tester.Client.UserV2.LockUser(CTX, &user.LockUserRequest{UserId: userID})
				return err
			},
		},
		{
			name: "deleted",
			disable: func(userID string) error {
				_, err := Tester.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: userID})
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disabledUser := Tester.CreateHumanUser(CTX)
			Tester.SetUserPassword(CTX, disabledUser.GetUserId(), integration.UserPassword, false)
			authRequestID := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
			sessionID, sessionToken, startTime, changeTime := Tester.CreatePasswordSession(t, CTXLOGIN, disabledUser.GetUserId(), integration.UserPassword)
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
			tokens, err := exchangeTokens(t, clientID, code, redirectURI)
			require.NoError(t, err)
			assertTokens(t, tokens, true)
			assertIDTokenClaims(t, tokens.IDTokenClaims, disabledUser.GetUserId(), armPassword, startTime, changeTime, sessionID)

			// make sure token works
			ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
			myUserResp, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
			require.NoError(t, err)
			require.Equal(t, disabledUser.GetUserId(), myUserResp.GetUser().GetId())

			// deactivate user
			err = tt.disable(disabledUser.GetUserId())
			require.NoError(t, err)

			// use token from deactivated user
			ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
			myUserResp, err = Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
			require.Error(t, err)
			require.Nil(t, myUserResp)
		})
	}
}

func createClient(t testing.TB) (clientID, projectID string) {
	return createClientWithOpts(t, clientOpts{
		redirectURI: redirectURI,
		logoutURI:   logoutRedirectURI,
		devMode:     false,
	})
}

type clientOpts struct {
	redirectURI string
	logoutURI   string
	devMode     bool
}

func createClientWithOpts(t testing.TB, opts clientOpts) (clientID, projectID string) {
	project, err := Tester.CreateProject(CTX)
	require.NoError(t, err)
	app, err := Tester.CreateOIDCNativeClient(CTX, opts.redirectURI, opts.logoutURI, project.GetId(), opts.devMode)
	require.NoError(t, err)
	return app.GetClientId(), project.GetId()
}

func createImplicitClient(t testing.TB) string {
	app, err := Tester.CreateOIDCImplicitFlowClient(CTX, redirectURIImplicit)
	require.NoError(t, err)
	return app.GetClientId()
}

func createAuthRequest(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequest(CTX, clientID, Tester.Users[integration.FirstInstanceUsersKey][integration.Login].ID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func createAuthRequestImplicit(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequestImplicit(CTX, clientID, Tester.Users[integration.FirstInstanceUsersKey][integration.Login].ID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func assertOIDCTime(t *testing.T, actual oidc.Time, expected time.Time) {
	assertOIDCTimeRange(t, actual, expected, expected)
}

func assertOIDCTimeRange(t *testing.T, actual oidc.Time, expectedStart, expectedEnd time.Time) {
	assert.WithinRange(t, actual.AsTime(), expectedStart.Add(-1*time.Second), expectedEnd.Add(1*time.Second))
}
