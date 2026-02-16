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
	"google.golang.org/grpc/metadata"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// This is a regression test for https://github.com/zitadel/zitadel/issues/11177
func TestIntrospection_RoleAddedAfterLogin(t *testing.T) {
	createIntrospectionApi := func(projectId string) rs.ResourceServer {
		api, err := Instance.CreateAPIClientJWT(CTX, projectId)
		require.NoError(t, err)
		keyResp, err := Instance.Client.Mgmt.AddAppKey(CTX, &management.AddAppKeyRequest{
			ProjectId:      projectId,
			AppId:          api.GetAppId(),
			Type:           authn.KeyType_KEY_TYPE_JSON,
			ExpirationDate: nil,
		})
		require.NoError(t, err)
		resourceServer, err := Instance.CreateResourceServerJWTProfile(CTX, keyResp.GetKeyDetails())
		require.NoError(t, err)
		return resourceServer
	}

	type TestCase struct {
		Name                  string
		RoleGrantedProjectId  string
		RoleGrantingProjectId string
		IntrospectionServer   rs.ResourceServer
		Role                  string
		ProjectGrantId        string // For cross-org project grants
	}
	testCases := []TestCase{
		func() TestCase {
			project := Instance.CreateProject(CTX, t, "", "single-project", false, false)
			projectId := project.GetId()
			const role = "role_1"
			Instance.AddProjectRole(CTX, t, projectId, role, role, "")
			return TestCase{
				Name:                  "single project",
				RoleGrantedProjectId:  projectId,
				RoleGrantingProjectId: projectId,
				IntrospectionServer:   createIntrospectionApi(projectId),
				Role:                  role,
			}
		}(),
		func() TestCase {
			project1 := Instance.CreateProject(CTX, t, "", "single-org-project-1", false, false)
			project1Id := project1.GetId()
			project2 := Instance.CreateProject(CTX, t, "", "single-org-project-2", false, false)
			project2Id := project2.GetId()
			const role = "role_2"
			Instance.AddProjectRole(CTX, t, project2.GetId(), role, role, "")
			return TestCase{
				Name:                  "project2 owns role, project 1 is granted",
				RoleGrantedProjectId:  project1Id,
				RoleGrantingProjectId: project2Id,
				IntrospectionServer:   createIntrospectionApi(project1Id),
				Role:                  role,
			}
		}(),
		func() TestCase {
			project1 := Instance.CreateProject(CTX, t, "", "cross-org-project-1", false, false)
			project1Id := project1.GetId()
			grantedOrg := Instance.CreateOrganization(CTXIAM, integration.OrganizationName(), integration.Email())
			grantedOrgID := grantedOrg.GetOrganizationId()
			ctxOrg := metadata.AppendToOutgoingContext(CTXIAM, "x-zitadel-orgid", grantedOrgID)
			project3 := Instance.CreateProject(ctxOrg, t, grantedOrgID, "cross-org-project-3", false, false)
			project3Id := project3.GetId()
			const role = "role_3"
			Instance.AddProjectRole(ctxOrg, t, project3Id, role, role, "")
			projectGrant, err := Instance.Client.Mgmt.AddProjectGrant(ctxOrg, &management.AddProjectGrantRequest{
				ProjectId:    project3Id,
				GrantedOrgId: Instance.DefaultOrg.GetId(),
				RoleKeys:     []string{role},
			})
			require.NoError(t, err)
			return TestCase{
				Name:                  "project3 (other org) owns role, project 1 is granted",
				RoleGrantedProjectId:  project1Id,
				RoleGrantingProjectId: project3Id,
				IntrospectionServer:   createIntrospectionApi(project1Id),
				Role:                  role,
				ProjectGrantId:        projectGrant.GetGrantId(),
			}
		}(),
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(tt *testing.T) {
			// create app to log in with
			app, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, testCase.RoleGrantedProjectId, false)
			require.NoError(t, err)

			// login with user
			scopes := []string{
				oidc.ScopeOpenID,
				oidc.ScopeProfile,
				"urn:zitadel:iam:org:project:roles",
				oidc_api.ScopeProjectsRoles,
				oidc_api.ScopeProjectRolePrefix + testCase.Role,
				fmt.Sprintf("urn:zitadel:iam:org:project:id:%s:aud", testCase.RoleGrantingProjectId),
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
			require.NoError(tt, err)
			code := assertCodeResponse(t, linkResp.GetCallbackUrl())
			tokens, err := exchangeTokens(t, Instance, app.GetClientId(), code, redirectURI)
			require.NoError(tt, err)
			assertTokens(tt, tokens, false)
			assertIDTokenClaims(tt, tokens.IDTokenClaims, User.GetUserId(), armPasskey, startTime, changeTime, sessionID)

			// introspect token and verify role is not present
			introspectionBefore, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), testCase.IntrospectionServer, tokens.AccessToken)
			require.NoError(tt, err)
			assert.True(tt, introspectionBefore.Active)
			projectRolesClaimName := fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, testCase.RoleGrantingProjectId)
			if roles, ok := introspectionBefore.Claims[projectRolesClaimName]; ok {
				if rolesMap, ok := roles.(map[string]interface{}); ok {
					assert.Contains(tt, rolesMap, testCase.Role, "role should not be present before granting")
				}
			}

			// assign role to user
			_, err = Instance.Client.Mgmt.AddUserGrant(CTX, &management.AddUserGrantRequest{
				UserId:         User.GetUserId(),
				ProjectId:      testCase.RoleGrantingProjectId,
				ProjectGrantId: testCase.ProjectGrantId,
				RoleKeys:       []string{testCase.Role},
			})
			require.NoError(t, err)
			// wait a bit to ensure the grant is processed
			time.Sleep(1 * time.Second)

			// introspect token and verify role is present
			introspectionAfter, err := rs.Introspect[*oidc.IntrospectionResponse](context.Background(), testCase.IntrospectionServer, tokens.AccessToken)
			require.NoError(t, err)
			require.True(tt, introspectionAfter.Active)
			roles, ok := introspectionAfter.Claims[projectRolesClaimName]
			require.True(tt, ok, "project roles claim should be present after granting")
			rolesMapAfter, ok := roles.(map[string]interface{})
			require.True(tt, ok, "project roles claim should be a map")

			assert.Contains(tt, rolesMapAfter, testCase.Role)
		})
	}
}
