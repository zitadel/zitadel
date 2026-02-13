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
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// TestIntrospection_RoleAddedAfterLogin tests that when a project grant (role) is added
// to an already logged-in user, the token introspection endpoint reflects the new role
// without requiring the user to log out and back in.
// This is a regression test for https://github.com/zitadel/zitadel/issues/11177
func TestIntrospection_RoleAddedAfterLogin(t *testing.T) {
	const (
		roleKey         = "test_role"
		roleDisplayName = "Test Role"
	)

	// Create project with a role
	project := Instance.CreateProject(CTX, t, "", integration.ProjectName(), false, false)
	Instance.AddProjectRole(CTX, t, project.GetId(), roleKey, roleDisplayName, "")

	// Create API client for introspection
	api, err := Instance.CreateAPIClientJWT(CTX, project.GetId())
	require.NoError(t, err)
	keyResp, err := Instance.Client.Mgmt.AddAppKey(CTX, &management.AddAppKeyRequest{
		ProjectId:      project.GetId(),
		AppId:          api.GetAppId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: nil,
	})
	require.NoError(t, err)

	// Create OIDC client
	app, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)

	// Create resource server for introspection
	resourceServer, err := Instance.CreateResourceServerJWTProfile(CTX, keyResp.GetKeyDetails())
	require.NoError(t, err)

	// Step 1 & 2: Login with user
	scope := []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopeOfflineAccess, oidc_api.ScopeResourceOwner}
	authRequestID := createAuthRequest(t, Instance, app.GetClientId(), redirectURI, scope...)
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

	// Code exchange
	code := assertCodeResponse(t, linkResp.GetCallbackUrl())
	tokens, err := exchangeTokens(t, Instance, app.GetClientId(), code, redirectURI)
	require.NoError(t, err)
	assertTokens(t, tokens, true)
	assertIDTokenClaims(t, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

	// Step 3: Perform introspection - should NOT have the role yet
	introspectionBefore, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServer, tokens.AccessToken)
	require.NoError(t, err)
	require.True(t, introspectionBefore.Active)

	// Assert role is not present yet
	projectRolesClaim := fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, project.GetId())
	rolesBefore, ok := introspectionBefore.Claims[projectRolesClaim]
	if ok {
		rolesMap, ok := rolesBefore.(map[string]interface{})
		if ok {
			_, hasRole := rolesMap[roleKey]
			assert.False(t, hasRole, "role should not be present before granting")
		}
	}

	// Step 4: Add user grant with the role to the already logged-in user
	_, err = Instance.Client.Mgmt.AddUserGrant(CTX, &management.AddUserGrantRequest{
		UserId:    User.GetUserId(),
		ProjectId: project.GetId(),
		RoleKeys:  []string{roleKey},
	})
	require.NoError(t, err)

	// Wait a bit to ensure the grant is processed
	time.Sleep(1 * time.Second)

	// Step 5: Perform introspection again - should now have the role
	// This is the bug: the role does not appear without logging out and back in
	introspectionAfter, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServer, tokens.AccessToken)
	require.NoError(t, err)
	require.True(t, introspectionAfter.Active)

	// Assert role IS present after granting
	rolesAfter, ok := introspectionAfter.Claims[projectRolesClaim]
	require.True(t, ok, "project roles claim should be present after granting")
	rolesMapAfter, ok := rolesAfter.(map[string]interface{})
	require.True(t, ok, "project roles claim should be a map")
	roleValue, hasRole := rolesMapAfter[roleKey]
	assert.True(t, hasRole, "role should be present after granting without logout/login")
	assert.NotNil(t, roleValue, "role value should not be nil")
}
