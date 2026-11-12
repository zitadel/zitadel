//go:build integration

package oidc_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestServer_CreateCallback_MaxAgeZero(t *testing.T) {
	project := Instance.CreateProject(CTX, t, "", integration.ProjectName(), false, false)
	client, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)
	user := Instance.CreateUserTypeHuman(CTX, integration.Email())
	Instance.SetUserPassword(CTX, user.GetId(), integration.UserPassword, false)
	sessionResp, err := Instance.Client.SessionV2.CreateSession(CTXLoginClient, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User:     &session.CheckUser{Search: &session.CheckUser_UserId{UserId: user.GetId()}},
			Password: &session.CheckPassword{Password: integration.UserPassword},
		},
	})
	require.NoError(t, err)

	provider, err := Instance.CreateRelyingParty(CTXLoginClient, client.GetClientId(), redirectURI)
	require.NoError(t, err)
	authURL, err := url.Parse(rp.AuthURL("state", provider, rp.WithCodeChallenge(oidc.NewSHACodeChallenge(integration.CodeVerifier))))
	require.NoError(t, err)
	query := authURL.Query()
	query.Set("max_age", "0")
	authURL.RawQuery = query.Encode()
	req, err := integration.GetRequest(authURL.String(), map[string]string{
		oidc_internal.LoginClientHeader: Instance.Users.Get(integration.UserTypeLogin).ID,
	})
	require.NoError(t, err)
	location, err := integration.CheckRedirect(req)
	require.NoError(t, err)
	loginPrefix := provider.Issuer() + Instance.Config.LoginURLV2
	require.True(t, strings.HasPrefix(location.String(), loginPrefix), "unexpected login redirect: %s", location)
	authRequestID := strings.TrimPrefix(location.String(), loginPrefix)

	callbackReq := &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{Session: &oidc_pb.Session{
			SessionId: sessionResp.GetSessionId(), SessionToken: sessionResp.GetSessionToken(),
		}},
	}
	_, err = Client.CreateCallback(CTXLoginClient, callbackReq)
	require.Equal(t, codes.FailedPrecondition, status.Code(err))

	updated, err := Instance.Client.SessionV2.SetSession(CTXLoginClient, &session.SetSessionRequest{
		SessionId: sessionResp.GetSessionId(),
		Checks:    &session.Checks{Password: &session.CheckPassword{Password: integration.UserPassword}},
	})
	require.NoError(t, err)
	callbackReq.GetSession().SessionToken = updated.GetSessionToken()
	callback, err := Client.CreateCallback(CTXLoginClient, callbackReq)
	require.NoError(t, err)
	callbackURL, err := url.Parse(callback.GetCallbackUrl())
	require.NoError(t, err)
	require.NotEmpty(t, callbackURL.Query().Get("code"))
	require.Equal(t, "state", callbackURL.Query().Get("state"))
}
