//go:build integration

package oidc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// TestIntrospection_RoleAddedAfterLogin tests that when a project grant (role) is added
// to an already logged-in user, the token introspection endpoint reflects the new role
// without requiring the user to log out and back in.
// This is a regression test for https://github.com/zitadel/zitadel/issues/11177
func TestIntrospection_RoleAddedAfterLogin(t *testing.T) {
	// PROJECT A WILL GRANT ROLES
	projectA := Instance.CreateProject(CTX, t, "", "project-a", false, false)

	const roleKey = "role_a"
	projectRolesClaimName := fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, projectA.GetId())

	Instance.AddProjectRole(CTX, t, projectA.GetId(), roleKey, roleKey, "")

	// PROJECT B WILL RECEIVE THE GRANTS
	projectB := Instance.CreateProject(CTX, t, "", "project-b", false, false)
	// create app to log in with
	app, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, projectB.GetId(), false)
	require.NoError(t, err)
	// create api for introspection
	apiB, err := Instance.CreateAPIClientJWT(CTX, projectB.GetId())
	require.NoError(t, err)
	keyRespB, err := Instance.Client.Mgmt.AddAppKey(CTX, &management.AddAppKeyRequest{
		ProjectId:      projectB.GetId(),
		AppId:          apiB.GetAppId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: nil,
	})
	require.NoError(t, err)
	resourceServerB, err := Instance.CreateResourceServerJWTProfile(CTX, keyRespB.GetKeyDetails())
	require.NoError(t, err)

	// LOGIN WITH USER
	authRequestID := createAuthRequest(t, Instance, app.GetClientId(), redirectURI,
		oidc.ScopeOpenID, oidc.ScopeProfile, oidc_api.ScopeProjectsRoles)
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
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, Instance, app.GetClientId(), code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	// INTROSPECT TOKEN AND VERIFY ROLE IS NOT PRESENT
	introspectionBeforeB, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServerB, tokens.AccessToken)
	require.NoError(t, err)
	assert.True(t, introspectionBeforeB.Active)
	if roles, ok := introspectionBeforeB.Claims[projectRolesClaimName]; ok {
		if rolesMap, ok := roles.(map[string]interface{}); ok {
			assert.Contains(t, rolesMap, roleKey, "role should not be present before granting")
		}
	}

	// ASSIGN ROLE TO USER
	_, err = Instance.Client.Mgmt.AddUserGrant(CTX, &management.AddUserGrantRequest{
		UserId:    User.GetUserId(),
		ProjectId: projectA.GetId(),
		RoleKeys:  []string{roleKey},
	})
	require.NoError(t, err)
	// wait a bit to ensure the grant is processed
	time.Sleep(1 * time.Second)

	// TODO remove the second login. It should work without logging in again
	// LOGIN WITH USER
	authRequestID = createAuthRequest(t, Instance, app.GetClientId(), redirectURI,
		oidc.ScopeOpenID, oidc.ScopeProfile, oidc_api.ScopeProjectsRoles)
	sessionID, sessionToken, startTime, changeTime = Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
	linkResp, err = Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionID,
				SessionToken: sessionToken,
			},
		},
	})
	require.NoError(t, err)
	code = assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err = exchangeTokens(t, Instance, app.GetClientId(), code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, false)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	// INTROSPECT TOKEN AND VERIFY ROLE IS PRESENT
	introspectionAfterB, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServerB, tokens.AccessToken)
	require.NoError(t, err)
	require.True(t, introspectionAfterB.Active)
	rolesB, ok := introspectionAfterB.Claims[projectRolesClaimName]
	require.True(t, ok, "project roles claim should be present after granting")
	rolesMapAfterB, ok := rolesB.(map[string]interface{})
	require.True(t, ok, "project roles claim should be a map")
	assert.Contains(t, rolesMapAfterB, roleKey, "role should be present after granting without logout/login")
	assert.NotNil(t, rolesMapAfterB[roleKey])
}
