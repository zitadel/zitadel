//go:build integration

package admin_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestServer_ExportData_includesUserGrantsOnGrantedProject(t *testing.T) {
	roleKey := integration.RoleKey()

	selfOrgID := Instance.CreateOrganization(AdminCTX, integration.OrganizationName(), integration.Email()).OrganizationId
	foreignOrg := Instance.CreateOrganization(AdminCTX, integration.OrganizationName(), integration.Email())

	projectID := Instance.CreateProject(AdminCTX, t, foreignOrg.OrganizationId, integration.ProjectName(), false, false).Id
	Instance.AddProjectRole(AdminCTX, t, projectID, roleKey, integration.RoleDisplayName(), "")
	Instance.CreateProjectGrant(AdminCTX, t, projectID, selfOrgID, roleKey)

	humanUser := Instance.CreateHumanUserVerified(AdminCTX, selfOrgID, integration.Email(), integration.Phone())
	machineUser := Instance.CreateUserTypeMachine(AdminCTX, selfOrgID)
	Instance.CreateAuthorizationProjectGrant(t, AdminCTX, projectID, selfOrgID, humanUser.GetUserId(), roleKey)
	Instance.CreateAuthorizationProjectGrant(t, AdminCTX, projectID, selfOrgID, machineUser.GetId(), roleKey)

	var export *admin.ExportDataResponse
	require.Eventually(t, func() bool {
		resp, err := Client.ExportData(AdminCTX, &admin.ExportDataRequest{
			// Granted org before project owner org: export must not depend on org processing order.
			OrgIds:         []string{selfOrgID, foreignOrg.OrganizationId},
			WithPasswords:  true,
			WithOtp:        true,
			ResponseOutput: true,
			Timeout:        time.Minute.String(),
		})
		if err != nil {
			return false
		}
		export = resp
		return userGrantsExported(export, selfOrgID, projectID, roleKey, humanUser.GetUserId(), machineUser.GetId())
	}, time.Minute, 200*time.Millisecond)

	require.True(t, userGrantsExported(export, selfOrgID, projectID, roleKey, humanUser.GetUserId(), machineUser.GetId()))
}

func userGrantsExported(export *admin.ExportDataResponse, orgID, projectID, roleKey string, userIDs ...string) bool {
	for _, userID := range userIDs {
		if !userGrantExported(export, orgID, userID, projectID, roleKey) {
			return false
		}
	}
	return true
}

func userGrantExported(export *admin.ExportDataResponse, orgID, userID, projectID, roleKey string) bool {
	for _, org := range export.GetOrgs() {
		if org.GetOrgId() != orgID {
			continue
		}
		for _, grant := range org.GetUserGrants() {
			if grant.GetUserId() != userID || grant.GetProjectId() != projectID {
				continue
			}
			for _, key := range grant.GetRoleKeys() {
				if key == roleKey {
					return grant.GetProjectGrantId() != ""
				}
			}
		}
	}
	return false
}
