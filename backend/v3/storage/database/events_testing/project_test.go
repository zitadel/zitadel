//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/project/v2"
)

func TestServer_ProjectReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id
	projectRepo := repository.ProjectRepository()

	projectName := integration.ProjectName()
	createRes, err := ProjectClient.CreateProject(CTX, &project.CreateProjectRequest{
		OrganizationId:        orgID,
		Name:                  projectName,
		ProjectRoleAssertion:  true,
		AuthorizationRequired: true,
		ProjectAccessRequired: true,
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)

	t.Run("create project reduces", func(t *testing.T) {
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectRepo.Get(CTX, pool, database.WithCondition(
				projectRepo.PrimaryKeyCondition(instanceID, createRes.GetProjectId()),
			))
			require.NoError(collect, err)
			assert.Equal(collect, createRes.GetProjectId(), dbProject.ID)
			assert.Equal(collect, orgID, dbProject.OrganizationID)
			assert.Equal(collect, projectName, dbProject.Name)
			assert.Equal(collect, domain.ProjectStateActive, dbProject.State)
			assert.True(collect, dbProject.ShouldAssertRole)
			assert.True(collect, dbProject.IsAuthorizationRequired)
			assert.True(collect, dbProject.IsProjectAccessRequired)
			assert.NotNil(collect, dbProject.CreatedAt)
			assert.NotNil(collect, dbProject.UpdatedAt)
		}, retryDuration, tick, "project not found within %v: %v", retryDuration, err)
	})

	t.Run("update project reduces", func(t *testing.T) {
		_, err := ProjectClient.UpdateProject(CTX, &project.UpdateProjectRequest{
			ProjectId:              createRes.GetProjectId(),
			Name:                   gu.Ptr("new name"),
			ProjectRoleAssertion:   gu.Ptr(false),
			AuthorizationRequired:  gu.Ptr(false),
			ProjectAccessRequired:  gu.Ptr(false),
			PrivateLabelingSetting: gu.Ptr(project.PrivateLabelingSetting_PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectRepo.Get(CTX, pool, database.WithCondition(
				projectRepo.PrimaryKeyCondition(instanceID, createRes.GetProjectId()),
			))
			require.NoError(collect, err)
			assert.Equal(collect, "new name", dbProject.Name)
			assert.False(collect, dbProject.ShouldAssertRole)
			assert.False(collect, dbProject.IsAuthorizationRequired)
			assert.False(collect, dbProject.IsProjectAccessRequired)
			assert.Equal(collect, int16(2), dbProject.UsedLabelingSettingOwner)
		}, retryDuration, tick, "project not updated within %v: %v", retryDuration, err)
	})

	t.Run("(de)activate project reduces", func(t *testing.T) {
		_, err := ProjectClient.DeactivateProject(CTX, &project.DeactivateProjectRequest{
			ProjectId: createRes.GetProjectId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectRepo.Get(CTX, pool, database.WithCondition(
				projectRepo.PrimaryKeyCondition(instanceID, createRes.GetProjectId()),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.ProjectStateInactive, dbProject.State)
		}, retryDuration, tick, "project not deactivated within %v: %v", retryDuration, err)

		_, err = ProjectClient.ActivateProject(CTX, &project.ActivateProjectRequest{
			ProjectId: createRes.GetProjectId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbProject, err := projectRepo.Get(CTX, pool, database.WithCondition(
				projectRepo.PrimaryKeyCondition(instanceID, createRes.GetProjectId()),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.ProjectStateActive, dbProject.State)
		}, retryDuration, tick, "project not activated within %v: %v", retryDuration, err)
	})

	t.Run("delete project reduces", func(t *testing.T) {
		_, err := ProjectClient.DeleteProject(CTX, &project.DeleteProjectRequest{
			ProjectId: createRes.GetProjectId(),
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := projectRepo.Get(CTX, pool, database.WithCondition(
				projectRepo.PrimaryKeyCondition(instanceID, createRes.GetProjectId()),
			))
			require.ErrorIs(collect, err, database.NewNoRowFoundError(nil))
		}, retryDuration, tick, "project not deleted within %v: %v", retryDuration, err)
	})
}
