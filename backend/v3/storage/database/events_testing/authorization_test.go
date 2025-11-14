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
	authorization_v2 "github.com/zitadel/zitadel/pkg/grpc/authorization/v2"
	project_v2beta "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

func TestServer_AuthorizationReduces(t *testing.T) {
	instanceID := Instance.ID()
	orgID := Instance.DefaultOrg.Id
	authorizationRepo := repository.AuthorizationRepository()

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)

	t.Run("user grant added reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})

		// create authorization
		createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectID,
			RoleKeys:       []string{role1, role2},
			OrganizationId: orgID,
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, []string{role1, role2}, az.Roles)
			assert.Equal(collect, domain.AuthorizationStateActive, az.State)
			assert.NotNil(collect, az.CreatedAt)
			assert.NotNil(collect, az.UpdatedAt)
		}, retryDuration, tick, "authorization not found within %v: %v", retryDuration, err)
	})

	t.Run("user grant update reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})

		// create authorization
		createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectID,
			OrganizationId: orgID,
			RoleKeys:       []string{role1, role2},
		})
		require.NoError(t, err)

		// add a new role to the project
		role3 := "role3"
		_, err = ProjectClient.AddProjectRole(CTX, &project_v2beta.AddProjectRoleRequest{
			ProjectId:   projectID,
			RoleKey:     role3,
			DisplayName: "display",
			Group:       nil,
		})
		require.NoError(t, err)

		// update roles to [role1, role2, role3]
		_, err = AuthorizationClient.UpdateAuthorization(CTX, &authorization_v2.UpdateAuthorizationRequest{
			Id:       createdAuthorization.Id,
			RoleKeys: []string{role1, role2, role3},
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			require.Len(collect, az.Roles, 3)
			assert.Equal(collect, []string{role1, role2, role3}, az.Roles)
		}, retryDuration, tick, "authorization not updated within %v: %v", retryDuration, err)
	})

	t.Run("user grant deactivate reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project
		projectID := prepareProjectAndProjectRoles(t, orgID, nil)
		// create authorization
		createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectID,
			OrganizationId: orgID,
		})
		require.NoError(t, err)

		// deactivate authorization
		_, err = AuthorizationClient.DeactivateAuthorization(CTX, &authorization_v2.DeactivateAuthorizationRequest{
			Id: createdAuthorization.Id,
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.AuthorizationStateInactive, az.State)
		}, retryDuration, tick, "authorization not deactivated within %v: %v", retryDuration, err)
	})

	t.Run("user grant activate reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project
		projectID := prepareProjectAndProjectRoles(t, orgID, nil)
		// create authorization
		createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectID,
			OrganizationId: orgID,
		})
		require.NoError(t, err)

		// deactivate authorization
		_, err = AuthorizationClient.DeactivateAuthorization(CTX, &authorization_v2.DeactivateAuthorizationRequest{
			Id: createdAuthorization.Id,
		})
		require.NoError(t, err)

		// re-activate authorization
		_, err = AuthorizationClient.ActivateAuthorization(CTX, &authorization_v2.ActivateAuthorizationRequest{
			Id: createdAuthorization.Id,
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.AuthorizationStateActive, az.State)
		}, retryDuration, tick, "authorization not activated within %v: %v", retryDuration, err)
	})

	t.Run("user grant removed reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})

		// create authorization with roles
		createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectID,
			RoleKeys:       []string{role1, role2},
			OrganizationId: orgID,
		})
		require.NoError(t, err)

		// delete authorization
		_, err = AuthorizationClient.DeleteAuthorization(CTX, &authorization_v2.DeleteAuthorizationRequest{
			Id: createdAuthorization.Id,
		})
		require.NoError(t, err)

		assert.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})

}

func prepareProjectAndProjectRoles(t *testing.T, orgID string, roles []string) string {
	project, err := ProjectClient.CreateProject(CTX, &project_v2beta.CreateProjectRequest{
		OrganizationId: orgID,
		Name:           integration.ProjectName(),
	})
	require.NoError(t, err)

	if len(roles) == 0 {
		return project.Id
	}

	for _, role := range roles {
		_, err = ProjectClient.AddProjectRole(CTX, &project_v2beta.AddProjectRoleRequest{
			ProjectId:   project.GetId(),
			RoleKey:     role,
			DisplayName: "display",
			Group:       nil,
		})
		require.NoError(t, err)
	}
	return project.Id
}
