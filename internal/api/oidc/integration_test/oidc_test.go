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
	"github.com/zitadel/zitadel/pkg/grpc/app"
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
	Instance *integration.Instance
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
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		User = Instance.CreateHumanUser(CTX)
		Instance.SetUserPassword(CTX, User.GetUserId(), integration.UserPassword, false)
		Instance.RegisterUserPasskey(CTX, User.GetUserId())
		CTXLOGIN = Instance.WithAuthorization(ctx, integration.UserTypeLogin)
		CTXIAM = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		return m.Run()
	}())
}

func Test_ZITADEL_API_missing_audience_scope(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID)
	sessionID, sessionToken, startTime, changeTime := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_missing_authentication(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	createResp, err := Instance.Client.SessionV2.CreateSession(CTXLOGIN, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
	})
	require.NoError(t, err)
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_missing_mfa_policy(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	org := Instance.CreateOrganization(CTXIAM, integration.OrganizationName(), integration.Email())
	userID := org.CreatedAdmins[0].GetUserId()
	Instance.SetUserPassword(CTXIAM, userID, integration.UserPassword, false)
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Instance.CreatePasswordSession(t, CTXLOGIN, userID, integration.UserPassword)
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)
	assertIDTokenClaims(t, tokens.IDTokenClaims, userID, armPassword, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	// pre check if request would succeed
	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, userID, myUserResp.GetUser().GetId())

	// require MFA
	ctxOrg := metadata.AppendToOutgoingContext(CTXIAM, "x-zitadel-orgid", org.GetOrganizationId())
	_, err = Instance.Client.Mgmt.AddCustomLoginPolicy(ctxOrg, &mgmt.AddCustomLoginPolicyRequest{
		ForceMfa: true,
	})
	require.NoError(t, err)

	// make sure policy is projected
	retryDuration := 5 * time.Second
	if ctxDeadline, ok := CTX.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		got, getErr := Instance.Client.Mgmt.GetLoginPolicy(ctxOrg, &mgmt.GetLoginPolicyRequest{})
		assert.NoError(ttt, getErr)
		assert.False(ttt, got.GetPolicy().IsDefault)

	}, retryDuration, time.Millisecond*100, "timeout waiting for login policy")

	// now it must fail
	myUserResp, err = Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_success(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())
}

func Test_ZITADEL_API_glob_redirects(t *testing.T) {
	const redirectURI = "https://my-org-1yfnjl2xj-my-app.vercel.app/api/auth/callback/zitadel"
	clientID, _ := createClientWithOpts(t, Instance, clientOpts{
		redirectURI: "https://my-org-*-my-app.vercel.app/api/auth/callback/zitadel",
		logoutURI:   "https://my-org-*-my-app.vercel.app/",
		devMode:     true,
	})
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))

	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())
}

func Test_ZITADEL_API_inactive_access_token(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	// make sure token works
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())

	// refresh token
	newTokens, err := refreshTokens(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)

	// use invalidated token
	ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err = Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_terminated_session(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	provider, err := Instance.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)
	authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
	sessionID, sessionToken, startTime, changeTime := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
	tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	// make sure token works
	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	require.Equal(t, User.GetUserId(), myUserResp.GetUser().GetId())

	// end session
	postLogoutRedirect, err := rp.EndSession(CTX, provider, tokens.IDToken, logoutRedirectURI, "state", "", nil)
	require.NoError(t, err)
	assert.Equal(t, logoutRedirectURI+"?state=state", postLogoutRedirect.String())

	// use token from terminated session
	ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
	myUserResp, err = Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.Error(t, err)
	require.Nil(t, myUserResp)
}

func Test_ZITADEL_API_terminated_session_user_disabled(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	tests := []struct {
		name    string
		disable func(userID string) error
	}{
		{
			name: "deactivated",
			disable: func(userID string) error {
				_, err := Instance.Client.UserV2.DeactivateUser(CTX, &user.DeactivateUserRequest{UserId: userID})
				return err
			},
		},
		{
			name: "locked",
			disable: func(userID string) error {
				_, err := Instance.Client.UserV2.LockUser(CTX, &user.LockUserRequest{UserId: userID})
				return err
			},
		},
		{
			name: "deleted",
			disable: func(userID string) error {
				_, err := Instance.Client.UserV2.DeleteUser(CTX, &user.DeleteUserRequest{UserId: userID})
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			disabledUser := Instance.CreateHumanUser(CTX)
			Instance.SetUserPassword(CTX, disabledUser.GetUserId(), integration.UserPassword, false)
			authRequestID := createAuthRequest(t, Instance, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess, zitadelAudienceScope)
			sessionID, sessionToken, startTime, changeTime := Instance.CreatePasswordSession(t, CTXLOGIN, disabledUser.GetUserId(), integration.UserPassword)
			linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
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
			tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
			require.NoError(t, err)
			assertTokens(t, tokens, true)
			assertIDTokenClaims(t, tokens.IDTokenClaims, disabledUser.GetUserId(), armPassword, startTime, changeTime, sessionID)

			// make sure token works
			ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
			myUserResp, err := Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
			require.NoError(t, err)
			require.Equal(t, disabledUser.GetUserId(), myUserResp.GetUser().GetId())

			// deactivate user
			err = tt.disable(disabledUser.GetUserId())
			require.NoError(t, err)

			// use token from deactivated user
			ctx = metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("%s %s", tokens.TokenType, tokens.AccessToken))
			myUserResp, err = Instance.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
			require.Error(t, err)
			require.Nil(t, myUserResp)
		})
	}
}

func createClient(t testing.TB, instance *integration.Instance) (clientID, projectID string) {
	return createClientWithOpts(t, instance, clientOpts{
		redirectURI:  redirectURI,
		logoutURI:    logoutRedirectURI,
		devMode:      false,
		LoginVersion: nil,
	})
}

func createClientLoginV2(t testing.TB, instance *integration.Instance) (clientID, projectID string) {
	return createClientWithOpts(t, instance, clientOpts{
		redirectURI:  redirectURI,
		logoutURI:    logoutRedirectURI,
		devMode:      false,
		LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: nil}}},
	})
}

type clientOpts struct {
	redirectURI  string
	logoutURI    string
	devMode      bool
	LoginVersion *app.LoginVersion
}

func createClientWithOpts(t testing.TB, instance *integration.Instance, opts clientOpts) (clientID, projectID string) {
	ctx := instance.WithAuthorization(CTX, integration.UserTypeOrgOwner)

	project := instance.CreateProject(ctx, t.(*testing.T), "", integration.ProjectName(), false, false)
	app, err := instance.CreateOIDCClientLoginVersion(ctx, opts.redirectURI, opts.logoutURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, opts.devMode, opts.LoginVersion)
	require.NoError(t, err)
	return app.GetClientId(), project.GetId()
}

func createImplicitClient(t testing.TB) string {
	app, err := Instance.CreateOIDCImplicitFlowClient(CTX, t.(*testing.T), redirectURIImplicit, nil)
	require.NoError(t, err)
	return app.GetClientId()
}

func createImplicitClientNoLoginClientHeader(t testing.TB) string {
	app, err := Instance.CreateOIDCImplicitFlowClient(CTX, t.(*testing.T), redirectURIImplicit, &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: nil}}})
	require.NoError(t, err)
	return app.GetClientId()
}

func createAuthRequest(t testing.TB, instance *integration.Instance, clientID, redirectURI string, scope ...string) string {
	_, redURL, err := instance.CreateOIDCAuthRequest(CTX, clientID, instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func createAuthRequestNoLoginClientHeader(t testing.TB, instance *integration.Instance, clientID, redirectURI string, scope ...string) string {
	_, redURL, err := instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientID, redirectURI, "", scope...)
	require.NoError(t, err)
	return redURL
}

func createAuthRequestImplicit(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Instance.CreateOIDCAuthRequestImplicit(CTX, clientID, Instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func createAuthRequestImplicitNoLoginClientHeader(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(CTX, clientID, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func assertOIDCTime(t *testing.T, actual oidc.Time, expected time.Time) {
	assertOIDCTimeRange(t, actual, expected, expected)
}

func assertOIDCTimeRange(t *testing.T, actual oidc.Time, expectedStart, expectedEnd time.Time) {
	assert.WithinRange(t, actual.AsTime(), expectedStart.Add(-10*time.Second), expectedEnd.Add(10*time.Second))
}
