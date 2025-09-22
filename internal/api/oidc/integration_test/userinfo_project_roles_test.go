//go:build integration

package oidc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestServer_UserInfo_ProjectsRolesClaim(t *testing.T) {
	const (
		roleFoo = "foo"
		roleBar = "bar"
		roleBaz = "baz"
	)

	clientID, projectID1 := createClient(t, Instance)
	_, projectID2 := createClient(t, Instance)
	addProjectRolesGrants(t, User.GetUserId(), projectID1, roleFoo, roleBar)
	addProjectRolesGrants(t, User.GetUserId(), projectID2, roleBaz)

	tests := []struct {
		name       string
		prepare    func(t *testing.T, clientID string, scope []string) *oidc.Tokens[*oidc.IDTokenClaims]
		scope      []string
		assertions []func(*testing.T, *oidc.UserInfo)
		wantErr    bool
	}{
		{
			name:    "project role and audience scope with only current project",
			prepare: getTokens,
			scope: []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
				domain.ProjectIDScope + projectID1 + domain.AudSuffix, domain.ProjectsIDScope,
			},
			assertions: []func(*testing.T, *oidc.UserInfo){
				assertUserinfo,
				func(t *testing.T, ui *oidc.UserInfo) {
					assertProjectsRolesClaims(t, ui.Claims, []string{roleFoo, roleBar}, []string{Instance.DefaultOrg.Id})
				},
			},
		},
		{
			name:    "project role and audience scope with multiple projects as audience",
			prepare: getTokens,
			scope: []string{oidc.ScopeProfile, oidc.ScopeOpenID, oidc.ScopeEmail, oidc.ScopeOfflineAccess,
				domain.ProjectIDScope + projectID1 + domain.AudSuffix,
				domain.ProjectIDScope + projectID2 + domain.AudSuffix,
				domain.ProjectsIDScope,
			},
			assertions: []func(*testing.T, *oidc.UserInfo){
				assertUserinfo,
				func(t *testing.T, ui *oidc.UserInfo) {
					assertProjectsRolesClaims(t, ui.Claims, []string{roleFoo, roleBar, roleBaz}, []string{Instance.DefaultOrg.Id})
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tt.prepare(t, clientID, tt.scope)
			provider, err := Instance.CreateRelyingParty(CTX, clientID, redirectURI)
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

// assertProjectsRoleClaims asserts the projectRoles in the claims.
func assertProjectsRolesClaims(t *testing.T, claims map[string]any, wantRoles, wantRoleOrgIDs []string) {
	t.Helper()
	projectsRoleClaims := make([]string, 0, 2)
	projectsRoleClaims = append(projectsRoleClaims, oidc_api.ClaimProjectsRoles)
	for _, claim := range projectsRoleClaims {
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
