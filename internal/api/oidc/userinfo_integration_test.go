//go:build integration

package oidc_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// TestServer_UserInfo is a top-level test which re-executes the actual
// userinfo integration test against a matrix of different feature flags.
// This ensure that the response of the different implementations remains the same.
func TestServer_UserInfo(t *testing.T) {
	iamOwnerCTX := Tester.WithAuthorization(CTX, integration.IAMOwner)
	t.Cleanup(func() {
		_, err := Tester.Client.FeatureV2.ResetInstanceFeatures(iamOwnerCTX, &feature.ResetInstanceFeaturesRequest{})
		require.NoError(t, err)
	})
	tests := []struct {
		name    string
		legacy  bool
		trigger bool
	}{
		{
			name:   "legacy enabled",
			legacy: true,
		},
		{
			name:    "legacy disabled, trigger disabled",
			legacy:  false,
			trigger: false,
		},
		{
			name:    "legacy disabled, trigger enabled",
			legacy:  false,
			trigger: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Tester.Client.FeatureV2.SetInstanceFeatures(iamOwnerCTX, &feature.SetInstanceFeaturesRequest{
				OidcLegacyIntrospection:             &tt.legacy,
				OidcTriggerIntrospectionProjections: &tt.trigger,
			})
			require.NoError(t, err)
			testServer_UserInfo(t)
		})
	}
}

// testServer_UserInfo is the actual userinfo integration test,
// which calls the userinfo endpoint with different client configurations, roles and token scopes.
func testServer_UserInfo(t *testing.T) {
	const (
		roleFoo = "foo"
		roleBar = "bar"
	)

	clientID, projectID := createClient(t)
	addProjectRolesGrants(t, User.GetUserId(), projectID, roleFoo, roleBar)

	tests := []struct {
		name       string
		prepare    func(t *testing.T, clientID string, scope []string) *oidc.Tokens[*oidc.IDTokenClaims]
		scope      []string
		assertions []func(*testing.T, *oidc.UserInfo)
		wantErr    bool
	}{
		{
			name: "invalid token",
			prepare: func(*testing.T, string, []string) *oidc.Tokens[*oidc.IDTokenClaims] {
				return &oidc.Tokens[*oidc.IDTokenClaims]{
					Token: &oauth2.Token{
						AccessToken: "DEAFBEEFDEADBEEF",
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: &oidc.IDTokenClaims{
						TokenClaims: oidc.TokenClaims{
							Subject: User.GetUserId(),
						},
					},
				}
			},
			scope: []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
			assertions: []func(*testing.T, *oidc.UserInfo){
				func(t *testing.T, ui *oidc.UserInfo) {
					assert.Nil(t, ui)
				},
			},
			wantErr: true,
		},
		{
			name:    "standard scopes",
			prepare: getTokens,
			scope:   []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
			assertions: []func(*testing.T, *oidc.UserInfo){
				assertUserinfo,
				func(t *testing.T, ui *oidc.UserInfo) {
					assertNoReservedScopes(t, ui.Claims)
				},
			},
		},
		{
			name: "project role assertion",
			prepare: func(t *testing.T, clientID string, scope []string) *oidc.Tokens[*oidc.IDTokenClaims] {
				_, err := Tester.Client.Mgmt.UpdateProject(CTX, &management.UpdateProjectRequest{
					Id:                   projectID,
					Name:                 fmt.Sprintf("project-%d", time.Now().UnixNano()),
					ProjectRoleAssertion: true,
				})
				require.NoError(t, err)
				t.Cleanup(func() {
					_, err := Tester.Client.Mgmt.UpdateProject(CTX, &management.UpdateProjectRequest{
						Id:                   projectID,
						Name:                 fmt.Sprintf("project-%d", time.Now().UnixNano()),
						ProjectRoleAssertion: false,
					})
					require.NoError(t, err)
				})
				resp, err := Tester.Client.Mgmt.GetProjectByID(CTX, &management.GetProjectByIDRequest{Id: projectID})
				require.NoError(t, err)
				require.True(t, resp.GetProject().GetProjectRoleAssertion(), "project role assertion")

				return getTokens(t, clientID, scope)
			},
			scope: []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess},
			assertions: []func(*testing.T, *oidc.UserInfo){
				assertUserinfo,
				func(t *testing.T, ui *oidc.UserInfo) {
					assertProjectRoleClaims(t, projectID, ui.Claims, true, []string{roleFoo, roleBar}, []string{Tester.Organisation.ID})
				},
			},
		},
		{
			name:    "project role scope",
			prepare: getTokens,
			scope: []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
				oidc_api.ScopeProjectRolePrefix + roleFoo,
			},
			assertions: []func(*testing.T, *oidc.UserInfo){
				assertUserinfo,
				func(t *testing.T, ui *oidc.UserInfo) {
					assertProjectRoleClaims(t, projectID, ui.Claims, true, []string{roleFoo}, []string{Tester.Organisation.ID})
				},
			},
		},
		{
			name:    "project role and audience scope",
			prepare: getTokens,
			scope: []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
				oidc_api.ScopeProjectRolePrefix + roleFoo,
				domain.ProjectIDScope + projectID + domain.AudSuffix,
			},
			assertions: []func(*testing.T, *oidc.UserInfo){
				assertUserinfo,
				func(t *testing.T, ui *oidc.UserInfo) {
					assertProjectRoleClaims(t, projectID, ui.Claims, true, []string{roleFoo}, []string{Tester.Organisation.ID})
				},
			},
		},
		{
			name: "PAT",
			prepare: func(t *testing.T, clientID string, scope []string) *oidc.Tokens[*oidc.IDTokenClaims] {
				user := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.OrgOwner)
				return &oidc.Tokens[*oidc.IDTokenClaims]{
					Token: &oauth2.Token{
						AccessToken: user.Token,
						TokenType:   oidc.BearerToken,
					},
					IDTokenClaims: &oidc.IDTokenClaims{
						TokenClaims: oidc.TokenClaims{
							Subject: user.ID,
						},
					},
				}
			},
			assertions: []func(*testing.T, *oidc.UserInfo){
				func(t *testing.T, ui *oidc.UserInfo) {
					user := Tester.Users.Get(integration.FirstInstanceUsersKey, integration.OrgOwner)
					assert.Equal(t, user.ID, ui.Subject)
					assert.Equal(t, user.PreferredLoginName, ui.PreferredUsername)
					assert.Equal(t, user.Machine.Name, ui.Name)
					assert.Equal(t, user.ResourceOwner, ui.Claims[oidc_api.ClaimResourceOwnerID])
					assert.NotEmpty(t, ui.Claims[oidc_api.ClaimResourceOwnerName])
					assert.NotEmpty(t, ui.Claims[oidc_api.ClaimResourceOwnerPrimaryDomain])
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tt.prepare(t, clientID, tt.scope)
			provider, err := Tester.CreateRelyingParty(CTX, clientID, redirectURI)
			require.NoError(t, err)
			userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			for _, assertion := range tt.assertions {
				assertion(t, userinfo)
			}
		})
	}
}

// TestServer_UserInfo_OrgIDRoles tests the [domain.OrgRoleIDScope] functionality
// it is a separate test because it is not supported in legacy mode.
func TestServer_UserInfo_OrgIDRoles(t *testing.T) {
	const (
		roleFoo = "foo"
		roleBar = "bar"
	)
	clientID, projectID := createClient(t)
	addProjectRolesGrants(t, User.GetUserId(), projectID, roleFoo, roleBar)
	grantedOrgID := addProjectOrgGrant(t, User.GetUserId(), projectID, roleFoo, roleBar)

	_, err := Tester.Client.Mgmt.UpdateProject(CTX, &management.UpdateProjectRequest{
		Id:                   projectID,
		Name:                 fmt.Sprintf("project-%d", time.Now().UnixNano()),
		ProjectRoleAssertion: true,
	})
	require.NoError(t, err)
	resp, err := Tester.Client.Mgmt.GetProjectByID(CTX, &management.GetProjectByIDRequest{Id: projectID})
	require.NoError(t, err)
	require.True(t, resp.GetProject().GetProjectRoleAssertion(), "project role assertion")

	tests := []struct {
		name           string
		scope          []string
		wantRoleOrgIDs []string
	}{
		{
			name: "default returns all role orgs",
			scope: []string{
				oidc.ScopeOpenID, oidc.ScopeOfflineAccess,
			},
			wantRoleOrgIDs: []string{Tester.Organisation.ID, grantedOrgID},
		},
		{
			name: "only granted org",
			scope: []string{
				oidc.ScopeOpenID, oidc.ScopeOfflineAccess,
				domain.OrgRoleIDScope + grantedOrgID},
			wantRoleOrgIDs: []string{grantedOrgID},
		},
		{
			name: "only own org",
			scope: []string{
				oidc.ScopeOpenID, oidc.ScopeOfflineAccess,
				domain.OrgRoleIDScope + Tester.Organisation.ID,
			},
			wantRoleOrgIDs: []string{Tester.Organisation.ID},
		},
		{
			name: "request both orgs",
			scope: []string{
				oidc.ScopeOpenID, oidc.ScopeOfflineAccess,
				domain.OrgRoleIDScope + Tester.Organisation.ID,
				domain.OrgRoleIDScope + grantedOrgID,
			},
			wantRoleOrgIDs: []string{Tester.Organisation.ID, grantedOrgID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := getTokens(t, clientID, tt.scope)
			provider, err := Tester.CreateRelyingParty(CTX, clientID, redirectURI)
			require.NoError(t, err)
			userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, tokens.IDTokenClaims.Subject, provider)
			require.NoError(t, err)
			assertProjectRoleClaims(t, projectID, userinfo.Claims, true, []string{roleFoo, roleBar}, tt.wantRoleOrgIDs)
		})
	}
}

// https://github.com/zitadel/zitadel/issues/6662
func TestServer_UserInfo_Issue6662(t *testing.T) {
	const (
		roleFoo = "foo"
		roleBar = "bar"
	)

	project, err := Tester.CreateProject(CTX)
	projectID := project.GetId()
	require.NoError(t, err)
	user, _, clientID, clientSecret, err := Tester.CreateOIDCCredentialsClient(CTX)
	require.NoError(t, err)
	addProjectRolesGrants(t, user.GetUserId(), projectID, roleFoo, roleBar)

	scope := []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
		oidc_api.ScopeProjectRolePrefix + roleFoo,
		domain.ProjectIDScope + projectID + domain.AudSuffix,
	}

	provider, err := rp.NewRelyingPartyOIDC(CTX, Tester.OIDCIssuer(), clientID, clientSecret, redirectURI, scope)
	require.NoError(t, err)
	tokens, err := rp.ClientCredentials(CTX, provider, nil)
	require.NoError(t, err)

	userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, tokens.TokenType, user.GetUserId(), provider)
	require.NoError(t, err)
	assertProjectRoleClaims(t, projectID, userinfo.Claims, false, []string{roleFoo}, []string{Tester.Organisation.ID})
}

func addProjectRolesGrants(t *testing.T, userID, projectID string, roles ...string) {
	t.Helper()
	bulkRoles := make([]*management.BulkAddProjectRolesRequest_Role, len(roles))
	for i, role := range roles {
		bulkRoles[i] = &management.BulkAddProjectRolesRequest_Role{
			Key:         role,
			DisplayName: role,
		}
	}
	_, err := Tester.Client.Mgmt.BulkAddProjectRoles(CTX, &management.BulkAddProjectRolesRequest{
		ProjectId: projectID,
		Roles:     bulkRoles,
	})
	require.NoError(t, err)
	_, err = Tester.Client.Mgmt.AddUserGrant(CTX, &management.AddUserGrantRequest{
		UserId:    userID,
		ProjectId: projectID,
		RoleKeys:  roles,
	})
	require.NoError(t, err)
}

// addProjectOrgGrant adds a new organization which will be granted on the projectID with the specified roles.
// The userID will be granted in the new organization to the project with the same roles.
func addProjectOrgGrant(t *testing.T, userID, projectID string, roles ...string) (grantedOrgID string) {
	grantedOrg := Tester.CreateOrganization(CTXIAM, fmt.Sprintf("ZITADEL_GRANTED_%d", time.Now().UnixNano()), fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()))
	projectGrant, err := Tester.Client.Mgmt.AddProjectGrant(CTX, &management.AddProjectGrantRequest{
		ProjectId:    projectID,
		GrantedOrgId: grantedOrg.GetOrganizationId(),
		RoleKeys:     roles,
	})
	require.NoError(t, err)

	ctxOrg := metadata.AppendToOutgoingContext(CTXIAM, "x-zitadel-orgid", grantedOrg.GetOrganizationId())
	_, err = Tester.Client.Mgmt.AddUserGrant(ctxOrg, &management.AddUserGrantRequest{
		UserId:         userID,
		ProjectId:      projectID,
		ProjectGrantId: projectGrant.GetGrantId(),
		RoleKeys:       roles,
	})
	require.NoError(t, err)
	return grantedOrg.GetOrganizationId()
}

func getTokens(t *testing.T, clientID string, scope []string) *oidc.Tokens[*oidc.IDTokenClaims] {
	authRequestID := createAuthRequest(t, clientID, redirectURI, scope...)
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

	return tokens
}

func assertUserinfo(t *testing.T, userinfo *oidc.UserInfo) {
	t.Helper()
	assert.Equal(t, User.GetUserId(), userinfo.Subject)
	assert.Equal(t, "Mickey", userinfo.GivenName)
	assert.Equal(t, "Mouse", userinfo.FamilyName)
	assert.Equal(t, "Mickey Mouse", userinfo.Name)
	assert.NotEmpty(t, userinfo.PreferredUsername)
	assert.Equal(t, userinfo.PreferredUsername, userinfo.Email)
	assert.False(t, bool(userinfo.EmailVerified))
	assertOIDCTime(t, userinfo.UpdatedAt, User.GetDetails().GetChangeDate().AsTime())
}

func assertNoReservedScopes(t *testing.T, claims map[string]any) {
	t.Helper()
	t.Log(claims)
	for claim := range claims {
		assert.Falsef(t, strings.HasPrefix(claim, oidc_api.ClaimPrefix), "claim %s has prefix %s", claim, oidc_api.ClaimPrefix)
	}
}

// assertProjectRoleClaims asserts the projectRoles in the claims.
// By default it searches for the [oidc_api.ClaimProjectRolesFormat] claim with a project ID,
// and optionally for the [oidc_api.ClaimProjectRoles] claim if claimProjectRole is true.
// Each claim should contain the roles expected by wantRoles and
// each role should contain the org IDs expected by wantRoleOrgIDs.
//
// In the claim map, each project role claim is expected to be a map of multiple roles and
// each role is expected to be a map of multiple Org IDs to Org Domains.
func assertProjectRoleClaims(t *testing.T, projectID string, claims map[string]any, claimProjectRole bool, wantRoles, wantRoleOrgIDs []string) {
	t.Helper()
	projectRoleClaims := []string{fmt.Sprintf(oidc_api.ClaimProjectRolesFormat, projectID)}
	if claimProjectRole {
		projectRoleClaims = append(projectRoleClaims, oidc_api.ClaimProjectRoles)
	}
	for _, claim := range projectRoleClaims {
		roleMap, ok := claims[claim].(map[string]any) // map of multiple roles
		require.Truef(t, ok, "claim %s not found or wrong type %T", claim, claims[claim])

		gotRoles := make([]string, 0, len(roleMap))
		for roleKey := range roleMap {
			role, ok := roleMap[roleKey].(map[string]any) // map of multiple org IDs to org domains
			require.Truef(t, ok, "role %s not found or wrong type %T", roleKey, roleMap[roleKey])
			gotRoles = append(gotRoles, roleKey)

			gotRoleOrgIDs := make([]string, 0, len(role))
			for orgID := range role {
				gotRoleOrgIDs = append(gotRoleOrgIDs, orgID)
			}
			assert.ElementsMatch(t, wantRoleOrgIDs, gotRoleOrgIDs)
		}
		assert.ElementsMatch(t, wantRoles, gotRoles)
	}
}
