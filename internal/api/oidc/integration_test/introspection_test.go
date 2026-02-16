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
	// SETUP PROJECT 1:
	// - role_1
	// - app to log in
	// - api for token introspection
	project1 := Instance.CreateProject(CTX, t, "", "project-1", false, false)
	const role1Key = "role_1"
	Instance.AddProjectRole(CTX, t, project1.GetId(), role1Key, role1Key, "")

	// create app to log in with
	app, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project1.GetId(), false)
	require.NoError(t, err)

	// create api for introspection
	api, err := Instance.CreateAPIClientJWT(CTX, project1.GetId())
	require.NoError(t, err)
	keyResp, err := Instance.Client.Mgmt.AddAppKey(CTX, &management.AddAppKeyRequest{
		ProjectId:      project1.GetId(),
		AppId:          api.GetAppId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: nil,
	})
	require.NoError(t, err)
	resourceServer, err := Instance.CreateResourceServerJWTProfile(CTX, keyResp.GetKeyDetails())
	require.NoError(t, err)

	// SETUP PROJECT 2:
	// - role_2
	project2 := Instance.CreateProject(CTX, t, "", "project-2", false, false)
	const role2Key = "role_2"
	Instance.AddProjectRole(CTX, t, project2.GetId(), role2Key, role2Key, "")

	// LOGIN WITH USER
	scopes := []string{
		oidc.ScopeOpenID,
		oidc.ScopeProfile,
		"urn:zitadel:iam:org:project:roles",
		oidc_api.ScopeProjectsRoles,
		oidc_api.ScopeProjectRolePrefix + role1Key,
		oidc_api.ScopeProjectRolePrefix + role2Key,
		fmt.Sprintf("urn:zitadel:iam:org:project:id:%s:aud", project1.GetId()),
		fmt.Sprintf("urn:zitadel:iam:org:project:id:%s:aud", project2.GetId()),
	}
	authRequestID := createAuthRequest(t, Instance, app.GetClientId(), redirectURI, scopes...)
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
	introspectionBefore, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServer, tokens.AccessToken)
	require.NoError(t, err)
	assert.True(t, introspectionBefore.Active)
	projectRolesClaimName := fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, project1.GetId())
	if roles, ok := introspectionBefore.Claims[projectRolesClaimName]; ok {
		if rolesMap, ok := roles.(map[string]interface{}); ok {
			assert.Contains(t, rolesMap, role1Key, "role should not be present before granting")
		}
	}

	// ASSIGN ROLE TO USER
	_, err = Instance.Client.Mgmt.AddUserGrant(CTX, &management.AddUserGrantRequest{
		UserId:    User.GetUserId(),
		ProjectId: project1.GetId(),
		RoleKeys:  []string{role1Key},
	})
	require.NoError(t, err)
	_, err = Instance.Client.Mgmt.AddUserGrant(CTX, &management.AddUserGrantRequest{
		UserId:    User.GetUserId(),
		ProjectId: project2.GetId(),
		RoleKeys:  []string{role2Key},
	})
	require.NoError(t, err)
	// wait a bit to ensure the grant is processed
	time.Sleep(1 * time.Second)

	// INTROSPECT TOKEN AND VERIFY ROLE IS PRESENT
	introspectionAfter, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), resourceServer, tokens.AccessToken)
	require.NoError(t, err)
	require.True(t, introspectionAfter.Active)
	roles, ok := introspectionAfter.Claims[projectRolesClaimName]
	require.True(t, ok, "project roles claim should be present after granting")
	rolesMapAfter, ok := roles.(map[string]interface{})
	require.True(t, ok, "project roles claim should be a map")
	assert.Contains(t, rolesMapAfter, role1Key)
	assert.Contains(t, rolesMapAfter, role2Key)
}
