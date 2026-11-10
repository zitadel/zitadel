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
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	v2beta_project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_ProjectGrantReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id
	projectGrantRepo := repository.ProjectGrantRepository()

	projectName := integration.ProjectName()
	createProjectRes, err := ProjectClient.CreateProject(CTX, &v2beta_project.CreateProjectRequest{
		OrganizationId:        orgID,
		Name:                  projectName,
		ProjectRoleAssertion:  true,
		AuthorizationRequired: true,
		ProjectAccessRequired: true,
	})
	require.NoError(t, err)
	keyName := "key"
	_, err = ProjectClient.AddProjectRole(CTX, &v2beta_project.AddProjectRoleRequest{
		ProjectId:   createProjectRes.GetId(),
		RoleKey:     keyName,
		DisplayName: "display",
		Group:       nil,
	})
	require.NoError(t, err)

	grantedOrgRes, err := OrgClient.CreateOrganization(CTX, &org.CreateOrganizationRequest{
		Name: integration.OrganizationName(),
	})
	require.NoError(t, err)
	_, err = ProjectClient.CreateProjectGrant(CTX, &v2beta_project.CreateProjectGrantRequest{
		ProjectId:             createProjectRes.GetId(),
		GrantedOrganizationId: grantedOrgRes.GetId(),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)

	t.Run("create project grant reduces", func(t *testing.T) {
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProjectGrant, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetId()),
				),
			))
			require.NoError(collect, err)
			assert.Equal(collect, createProjectRes.GetId(), dbProjectGrant.ProjectID)
			assert.Equal(collect, orgID, dbProjectGrant.GrantingOrganizationID)
			assert.Equal(collect, grantedOrgRes.GetId(), dbProjectGrant.GrantedOrganizationID)
			assert.Equal(collect, domain.ProjectGrantStateActive, dbProjectGrant.State)
			assert.NotNil(collect, dbProjectGrant.CreatedAt)
			assert.NotNil(collect, dbProjectGrant.UpdatedAt)
		}, retryDuration, tick, "project grant not found within %v: %v", retryDuration, err)
	})

	t.Run("update project grant reduces", func(t *testing.T) {
		_, err := ProjectClient.UpdateProjectGrant(CTX, &v2beta_project.UpdateProjectGrantRequest{
			ProjectId:             createProjectRes.GetId(),
			GrantedOrganizationId: grantedOrgRes.GetId(),
			RoleKeys:              []string{keyName},
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProjectGrant, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetId()),
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
		_, err := ProjectClient.DeactivateProjectGrant(CTX, &v2beta_project.DeactivateProjectGrantRequest{
			ProjectId:             createProjectRes.GetId(),
			GrantedOrganizationId: grantedOrgRes.GetId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetId()),
				),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.ProjectGrantStateInactive, dbProject.State)
		}, retryDuration, tick, "project grant not deactivated within %v: %v", retryDuration, err)

		_, err = ProjectClient.ActivateProjectGrant(CTX, &v2beta_project.ActivateProjectGrantRequest{
			ProjectId:             createProjectRes.GetId(),
			GrantedOrganizationId: grantedOrgRes.GetId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetId()),
				),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.ProjectGrantStateActive, dbProject.State)
		}, retryDuration, tick, "project grant not activated within %v: %v", retryDuration, err)
	})

	t.Run("delete project grant reduces", func(t *testing.T) {
		_, err := ProjectClient.DeleteProjectGrant(CTX, &v2beta_project.DeleteProjectGrantRequest{
			ProjectId:             createProjectRes.GetId(),
			GrantedOrganizationId: grantedOrgRes.GetId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := projectGrantRepo.Get(CTX, pool, database.WithCondition(
				database.And(
					projectGrantRepo.InstanceIDCondition(instanceID),
					projectGrantRepo.ProjectIDCondition(createProjectRes.GetId()),
					projectGrantRepo.GrantedOrganizationIDCondition(grantedOrgRes.GetId()),
				),
			))
			require.ErrorIs(collect, err, database.NewNoRowFoundError(nil))
		}, retryDuration, tick, "project grant not deleted within %v: %v", retryDuration, err)
	})
}
