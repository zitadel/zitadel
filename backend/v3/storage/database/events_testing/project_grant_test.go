//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2"
	"github.com/zitadel/zitadel/pkg/grpc/project/v2"
)

func TestServer_ProjectGrantReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id
	projectGrantRepo := repository.ProjectGrantRepository()

	projectName := integration.ProjectName()
	createProjectRes, err := ProjectClient.CreateProject(CTX, &project.CreateProjectRequest{
		OrganizationId:        orgID,
		Name:                  projectName,
		ProjectRoleAssertion:  true,
		AuthorizationRequired: true,
		ProjectAccessRequired: true,
	})
	require.NoError(t, err)
	keyName := "key"
	_, err = ProjectClient.AddProjectRole(CTX, &project.AddProjectRoleRequest{
		ProjectId:   createProjectRes.GetProjectId(),
		RoleKey:     keyName,
		DisplayName: "display",
		Group:       nil,
	})
	require.NoError(t, err)

	grantedOrgRes, err := OrgClient.AddOrganization(CTX, &org.AddOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	_, err = ProjectClient.CreateProjectGrant(CTX, &project.CreateProjectGrantRequest{
		ProjectId:             createProjectRes.GetProjectId(),
		GrantedOrganizationId: grantedOrgRes.GetOrganizationId(),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)

	t.Run("create project grant reduces", func(t *testing.T) {
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProjectGrant, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetProjectId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetOrganizationId()),
				),
			))
			require.NoError(collect, err)
			assert.Equal(collect, createProjectRes.GetProjectId(), dbProjectGrant.ProjectID)
			assert.Equal(collect, orgID, dbProjectGrant.GrantingOrganizationID)
			assert.Equal(collect, grantedOrgRes.GetOrganizationId(), dbProjectGrant.GrantedOrganizationID)
			assert.Equal(collect, domain.ProjectGrantStateActive, dbProjectGrant.State)
			assert.NotNil(collect, dbProjectGrant.CreatedAt)
			assert.NotNil(collect, dbProjectGrant.UpdatedAt)
		}, retryDuration, tick, "project grant not found within %v: %v", retryDuration, err)
	})

	t.Run("update project grant reduces", func(t *testing.T) {
		_, err := ProjectClient.UpdateProjectGrant(CTX, &project.UpdateProjectGrantRequest{
			ProjectId:             createProjectRes.GetProjectId(),
			GrantedOrganizationId: grantedOrgRes.GetOrganizationId(),
			RoleKeys:              []string{keyName},
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProjectGrant, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetProjectId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetOrganizationId()),
				),
			))
			require.NoError(collect, err)
			if !assert.Len(collect, dbProjectGrant.RoleKeys, 1) {
				return
			}
			assert.Equal(collect, dbProjectGrant.RoleKeys[0], keyName)
		}, retryDuration, tick, "project grant not updated within %v: %v", retryDuration, err)
	})

	t.Run("(de)activate project grant reduces", func(t *testing.T) {
		_, err := ProjectClient.DeactivateProjectGrant(CTX, &project.DeactivateProjectGrantRequest{
			ProjectId:             createProjectRes.GetProjectId(),
			GrantedOrganizationId: grantedOrgRes.GetOrganizationId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetProjectId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetOrganizationId()),
				),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.ProjectGrantStateInactive, dbProject.State)
		}, retryDuration, tick, "project grant not deactivated within %v: %v", retryDuration, err)

		_, err = ProjectClient.ActivateProjectGrant(CTX, &project.ActivateProjectGrantRequest{
			ProjectId:             createProjectRes.GetProjectId(),
			GrantedOrganizationId: grantedOrgRes.GetOrganizationId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetProjectId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetOrganizationId()),
				),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.ProjectGrantStateActive, dbProject.State)
		}, retryDuration, tick, "project grant not activated within %v: %v", retryDuration, err)
	})

	t.Run("delete project grant reduces", func(t *testing.T) {
		_, err := ProjectClient.DeleteProjectGrant(CTX, &project.DeleteProjectGrantRequest{
			ProjectId:             createProjectRes.GetProjectId(),
			GrantedOrganizationId: grantedOrgRes.GetOrganizationId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetProjectId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetOrganizationId()),
				),
			))
			require.ErrorIs(collect, err, database.NewNoRowFoundError(nil))
		}, retryDuration, tick, "project grant not deleted within %v: %v", retryDuration, err)
	})
}
