package events_test

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	v2beta_project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_ProjectRoleReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id
	roleRepo := repository.ProjectRepository().Role()

	projectRes, err := ProjectClient.CreateProject(CTX, &v2beta_project.CreateProjectRequest{
		OrganizationId:        orgID,
		Name:                  "foobar",
		ProjectRoleAssertion:  true,
		AuthorizationRequired: true,
		ProjectAccessRequired: true,
	})
	require.NoError(t, err)
	projectID := projectRes.GetId()
	_, err = ProjectClient.AddProjectRole(CTX, &v2beta_project.AddProjectRoleRequest{
		ProjectId:   projectID,
		RoleKey:     "key",
		DisplayName: "display name",
		Group:       proto.String("group"),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)

	t.Run("add project role reduces", func(t *testing.T) {
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbRole, err := roleRepo.Get(CTX, pool, database.WithCondition(
				roleRepo.PrimaryKeyCondition(instanceID, projectID, "key"),
			))
			require.NoError(collect, err)
			require.NotNil(collect, dbRole)
			assert.Equal(collect, projectID, dbRole.ProjectID)
			assert.Equal(collect, "key", dbRole.Key)
			assert.Equal(collect, "display name", dbRole.DisplayName)
			assert.Equal(collect, gu.Ptr("group"), dbRole.RoleGroup)
			assert.NotNil(collect, dbRole.CreatedAt)
			assert.NotNil(collect, dbRole.UpdatedAt)
		}, retryDuration, tick, "project role not found within %v: %v", retryDuration, err)
	})

	t.Run("update project role reduces", func(t *testing.T) {
		_, err := ProjectClient.UpdateProjectRole(CTX, &v2beta_project.UpdateProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     "key",
			DisplayName: proto.String("new display name"),
			Group:       proto.String("new group"),
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			dbRole, err := roleRepo.Get(CTX, pool, database.WithCondition(
				roleRepo.PrimaryKeyCondition(instanceID, projectID, "key"),
			))
			require.NoError(collect, err)
			require.NotNil(collect, dbRole)
			assert.Equal(collect, "new display name", dbRole.DisplayName)
			assert.Equal(collect, gu.Ptr("new group"), dbRole.RoleGroup)
		}, retryDuration, tick, "project role not updated within %v: %v", retryDuration, err)
	})

	t.Run("remove project role reduces", func(t *testing.T) {
		_, err := ProjectClient.RemoveProjectRole(CTX, &v2beta_project.RemoveProjectRoleRequest{
			ProjectId: projectID,
			RoleKey:   "key",
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			_, err := roleRepo.Get(CTX, pool, database.WithCondition(
				roleRepo.PrimaryKeyCondition(instanceID, projectID, "key"),
			))
			require.ErrorIs(collect, err, database.NewNoRowFoundError(nil))
		}, retryDuration, tick, "project role not deleted within %v: %v", retryDuration, err)
	})
}
