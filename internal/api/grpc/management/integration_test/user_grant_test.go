//go:build integration

package management_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

// PoC for GHSA-v859-c572-qh5p
//
// Updating a project grant to drop roles cascades into the user grants of that
// grant via removeRoleFromUserGrant (internal/command/user_grant.go), which
// removes the dropped roles in a single call. Because that loop mutates the
// slice while iterating forward, two *adjacent* removed roles cause the second to be
// skipped: the user keeps authorisation that should have been revoked.
//
// Here a user grant holds ["admin","viewer","editor"]; the project grant is
// updated to keep only ["editor"], so "admin" and "viewer" should both be
// stripped. The test asserts the CORRECT result (["editor"]); against faulty
// implementations the grant ends up ["viewer","editor"], so it fails.
func TestServer_UpdateProjectGrant_cascadeRoleRemovalSkipsAdjacent(t *testing.T) {
	const (
		roleAdmin  = "admin"
		roleViewer = "viewer"
		roleEditor = "editor"
	)

	// Org that the project will be granted to.
	grantedOrg := Instance.CreateOrganization(IAMOwnerCTX, integration.OrganizationName(), integration.Email())
	grantedOrgID := grantedOrg.GetOrganizationId()
	grantedOrgCTX := integration.SetOrgID(IAMOwnerCTX, grantedOrgID)

	// Project (owned by the default org) with three roles, in order.
	projectResp, err := Client.AddProject(OrgCTX, &mgmt_pb.AddProjectRequest{
		Name: integration.ProjectName(),
	})
	require.NoError(t, err)
	projectID := projectResp.GetId()

	for _, role := range []string{roleAdmin, roleViewer, roleEditor} {
		_, err = Client.AddProjectRole(OrgCTX, &mgmt_pb.AddProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     role,
			DisplayName: role,
		})
		require.NoError(t, err)
	}

	// Grant the project (with all three roles) to the granted org.
	grantResp, err := Client.AddProjectGrant(OrgCTX, &mgmt_pb.AddProjectGrantRequest{
		ProjectId:    projectID,
		GrantedOrgId: grantedOrgID,
		RoleKeys:     []string{roleAdmin, roleViewer, roleEditor},
	})
	require.NoError(t, err)
	grantID := grantResp.GetGrantId()

	// A user in the granted org, holding all three roles on the project grant.
	userResp := Instance.CreateHumanUserVerified(IAMOwnerCTX, grantedOrgID, integration.Email(), "+41791234567")
	userID := userResp.GetUserId()

	userGrantResp, err := Client.AddUserGrant(grantedOrgCTX, &mgmt_pb.AddUserGrantRequest{
		UserId:         userID,
		ProjectId:      projectID,
		ProjectGrantId: grantID,
		RoleKeys:       []string{roleAdmin, roleViewer, roleEditor},
	})
	require.NoError(t, err)
	userGrantID := userGrantResp.GetUserGrantId()

	// UpdateProjectGrant only cascades to user grants the query side already
	// returns, so wait for the projection to catch up to all three roles.
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		grant, err := Client.GetUserGrantByID(grantedOrgCTX, &mgmt_pb.GetUserGrantByIDRequest{
			UserId:  userID,
			GrantId: userGrantID,
		})
		assert.NoError(ct, err)
		assert.ElementsMatch(ct, []string{roleAdmin, roleViewer, roleEditor}, grant.GetUserGrant().GetRoleKeys())
	}, retryDuration, tick)

	// Drop "admin" and "viewer" from the project grant, keeping only "editor".
	// This cascades removeRoleFromUserGrant(["admin","viewer"]) into the user grant.
	_, err = Client.UpdateProjectGrant(OrgCTX, &mgmt_pb.UpdateProjectGrantRequest{
		ProjectId: projectID,
		GrantId:   grantID,
		RoleKeys:  []string{roleEditor},
	})
	require.NoError(t, err)

	// Wait until the cascade has been applied and projected. "admin" (the first
	// removed role) is dropped by both the buggy and fixed code, so its absence
	// is a stable signal that the cascade has converged.
	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	var finalRoleKeys []string
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		grant, err := Client.GetUserGrantByID(grantedOrgCTX, &mgmt_pb.GetUserGrantByIDRequest{
			UserId:  userID,
			GrantId: userGrantID,
		})
		assert.NoError(ct, err)
		finalRoleKeys = grant.GetUserGrant().GetRoleKeys()
		assert.NotContains(ct, finalRoleKeys, roleAdmin)
	}, retryDuration, tick)

	// Correct expectation: only "editor" survives. With the off-by-one bug the
	// user grant is ["viewer","editor"] — "viewer" was skipped during removal.
	assert.ElementsMatch(t, []string{roleEditor}, finalRoleKeys)
}
