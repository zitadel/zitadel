//go:build integration

package oidc_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// TestServer_CreateCallback_LoginClientHeaderGrantsNoPermission ensures that the
// x-zitadel-login-client header of the (unauthenticated) authorize request does not
// grant the named user any permission on the auth request.
// The header is chosen by whoever starts the flow, so a user without the session.link
// permission must not be able to complete an auth request by naming itself as the login client.
func TestServer_CreateCallback_LoginClientHeaderGrantsNoPermission(t *testing.T) {
	project := Instance.CreateProject(CTX, t, "", integration.ProjectName(), false, false)
	client, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)

	caller := Instance.Users.Get(integration.UserTypeNoPermission)
	callerCtx := Instance.WithAuthorization(CTX, integration.UserTypeNoPermission)
	victimID := Instance.CreateHumanUser(CTX).GetUserId()

	t.Run("create callback", func(t *testing.T) {
		// a session with only a user check, as created when the login name is entered in the login UI
		sessionResp := createSession(t, CTXLoginClient, victimID)
		_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), caller.ID, redirectURI)
		require.NoError(t, err)

		got, err := Client.CreateCallback(callerCtx, &oidc_pb.CreateCallbackRequest{
			AuthRequestId: authRequestID,
			CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
				Session: &oidc_pb.Session{
					SessionId:    sessionResp.GetSessionId(),
					SessionToken: sessionResp.GetSessionToken(),
				},
			},
		})
		if err == nil {
			t.Fatalf("auth request completed by a caller without session.link, issued tokens for subject %q (victim %q)",
				exchangeCallbackCode(t, client.GetClientId(), got.GetCallbackUrl()), victimID)
		}
		assertDenied(t, err)
	})

	t.Run("fail auth request", func(t *testing.T) {
		_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), caller.ID, redirectURI)
		require.NoError(t, err)

		_, err = Client.CreateCallback(callerCtx, &oidc_pb.CreateCallbackRequest{
			AuthRequestId: authRequestID,
			CallbackKind: &oidc_pb.CreateCallbackRequest_Error{
				Error: &oidc_pb.AuthorizationError{Error: oidc_pb.ErrorReason_ERROR_REASON_ACCESS_DENIED},
			},
		})
		assertDenied(t, err)
	})

	t.Run("login client with permission", func(t *testing.T) {
		sessionResp := createSession(t, CTXLoginClient, victimID)
		_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeLogin).ID, redirectURI)
		require.NoError(t, err)

		_, err = Client.CreateCallback(CTXLoginClient, &oidc_pb.CreateCallbackRequest{
			AuthRequestId: authRequestID,
			CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
				Session: &oidc_pb.Session{
					SessionId:    sessionResp.GetSessionId(),
					SessionToken: sessionResp.GetSessionToken(),
				},
			},
		})
		require.NoError(t, err)
	})
}

// exchangeCallbackCode redeems the code of the callback URL (PKCE, public client) and returns the subject of the ID token.
func exchangeCallbackCode(t *testing.T, clientID, callbackURL string) string {
	parsed, err := url.Parse(callbackURL)
	require.NoError(t, err)
	provider, err := Instance.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)
	tokens, err := rp.CodeExchange[*oidc.IDTokenClaims](CTX, parsed.Query().Get("code"), provider, rp.WithCodeVerifier(integration.CodeVerifier))
	require.NoError(t, err)
	return tokens.IDTokenClaims.Subject
}

// assertDenied checks the call was rejected by the permission check.
// A caller without any membership gets NotFound ("membership not found"), one with memberships PermissionDenied.
func assertDenied(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	st := status.Convert(err)
	switch st.Code() {
	case codes.PermissionDenied:
	case codes.NotFound:
		// a caller without any membership is rejected by the permission check itself;
		// assert the message so the test cannot pass on a genuinely missing request.
		assert.Contains(t, st.Message(), "membership not found", err)
	default:
		t.Errorf("expected PermissionDenied or NotFound, got %s: %v", st.Code(), err)
	}
}
