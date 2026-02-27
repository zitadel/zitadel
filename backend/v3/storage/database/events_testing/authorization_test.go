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
	instance_v2 "github.com/zitadel/zitadel/pkg/grpc/instance/v2"
	org_v2beta "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	project_v2beta "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
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
		createAndEnsureAuthorization(t, instanceID, orgID, user.UserId, projectID, []string{role1, role2}, retryDuration, tick)
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
		updated, err := AuthorizationClient.UpdateAuthorization(CTX, &authorization_v2.UpdateAuthorizationRequest{
			Id:       createdAuthorization.Id,
			RoleKeys: []string{role1, role2, role3},
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			require.Len(collect, az.Roles, 3)
			assert.Equal(collect, []string{role1, role2, role3}, az.Roles)
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.CreatedAt.UTC())
			assert.Equal(collect, updated.GetChangeDate().AsTime(), az.UpdatedAt.UTC())
		}, retryDuration, tick, "authorization not updated within %v: %v", retryDuration, err)
	})

	t.Run("user grant deactivate reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		projectID := prepareProjectAndProjectRoles(t, orgID, nil)
		// create authorization
		createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectID,
			OrganizationId: orgID,
		})
		require.NoError(t, err)

		// deactivate authorization
		deactivated, err := AuthorizationClient.DeactivateAuthorization(CTX, &authorization_v2.DeactivateAuthorizationRequest{
			Id: createdAuthorization.Id,
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.AuthorizationStateInactive, az.State)
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.CreatedAt.UTC())
			assert.Equal(collect, deactivated.GetChangeDate().AsTime(), az.UpdatedAt.UTC())
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
		reactivated, err := AuthorizationClient.ActivateAuthorization(CTX, &authorization_v2.ActivateAuthorizationRequest{
			Id: createdAuthorization.Id,
		})
		require.NoError(t, err)

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, domain.AuthorizationStateActive, az.State)
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.CreatedAt.UTC())
			assert.Equal(collect, reactivated.GetChangeDate().AsTime(), az.UpdatedAt.UTC())
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

		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})

	t.Run("user removed reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})
		// create authorization with roles
		createdAuthorization := createAndEnsureAuthorization(t, instanceID, orgID, user.UserId, projectID, []string{role1, role2}, retryDuration, tick)
		// delete the user
		_, err := UserClient.DeleteUser(CTX, &user_v2.DeleteUserRequest{
			UserId: user.UserId,
		})
		require.NoError(t, err)

		// ensure authorization is deleted
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})

	t.Run("project removed reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})
		// create authorization with roles
		createdAuthorization := createAndEnsureAuthorization(t, instanceID, orgID, user.UserId, projectID, []string{role1, role2}, retryDuration, tick)
		// delete the project
		_, err := ProjectClient.DeleteProject(CTX, &project_v2beta.DeleteProjectRequest{
			Id: projectID,
		})
		require.NoError(t, err)

		// ensure authorization is deleted
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})

	t.Run("project role removed reduces", func(t *testing.T) {
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})
		// create authorization with roles
		createdAuthorization := createAndEnsureAuthorization(t, instanceID, orgID, user.UserId, projectID, []string{role1, role2}, retryDuration, tick)
		// delete a project role
		roleRemoved, err := ProjectClient.RemoveProjectRole(CTX, &project_v2beta.RemoveProjectRoleRequest{
			ProjectId: projectID,
			RoleKey:   role2,
		})
		require.NoError(t, err)

		// ensure authorization is updated
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, []string{role1}, az.Roles)
			assert.Equal(collect, domain.AuthorizationStateActive, az.State)
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.CreatedAt.UTC())
			// project role removal triggers a UserGrantCascadeChangedEvent,
			// so the RT UpdatedAt should match the role removal time
			assert.Equal(collect, roleRemoved.GetRemovalDate().AsTime(), az.UpdatedAt.UTC())
		}, retryDuration, tick, "authorization not updated within %v: %v", retryDuration, err)
	})

	t.Run("user grant added for a project grant reduces", func(t *testing.T) {
		// prepare project and project roles
		role1 := "role1"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1})

		// granted organization
		grantedOrganizationName := integration.OrganizationName()
		grantedOrganization := Instance.CreateOrganization(CTX, grantedOrganizationName, integration.Email())

		// create project grant
		_, err := ProjectClient.CreateProjectGrant(CTX, &project_v2beta.CreateProjectGrantRequest{
			ProjectId:             projectID,
			RoleKeys:              []string{role1},
			GrantedOrganizationId: grantedOrganization.OrganizationId,
		})
		require.NoError(t, err)

		// create user
		user := Instance.CreateHumanUserVerified(CTX, grantedOrganization.OrganizationId, integration.Email(), integration.Phone())
		// create authorization with roles
		createAndEnsureAuthorization(t, instanceID, grantedOrganization.OrganizationId, user.UserId, projectID, []string{role1}, retryDuration, tick)
	})

	t.Run("project grant removed reduces", func(t *testing.T) {
		// prepare project and project roles
		role1 := "role1"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1})

		// granted organization
		grantedOrganizationName := integration.OrganizationName()
		grantedOrganization := Instance.CreateOrganization(CTX, grantedOrganizationName, integration.Email())

		// create project grant
		_, err := ProjectClient.CreateProjectGrant(CTX, &project_v2beta.CreateProjectGrantRequest{
			ProjectId:             projectID,
			RoleKeys:              []string{role1},
			GrantedOrganizationId: grantedOrganization.OrganizationId,
		})
		require.NoError(t, err)
		// create user
		user := Instance.CreateHumanUserVerified(CTX, grantedOrganization.OrganizationId, integration.Email(), integration.Phone())
		// create authorization with roles
		createdAuthorization := createAndEnsureAuthorization(t, instanceID, grantedOrganization.OrganizationId, user.UserId, projectID, []string{role1}, retryDuration, tick)
		// delete project grant
		_, err = ProjectClient.DeleteProjectGrant(CTX, &project_v2beta.DeleteProjectGrantRequest{
			ProjectId:             projectID,
			GrantedOrganizationId: grantedOrganization.OrganizationId,
		})
		require.NoError(t, err)

		// ensure authorization is deleted
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})

	t.Run("project grant updated reduces", func(t *testing.T) {
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgID, []string{role1, role2})
		// granted organization
		grantedOrganizationName := integration.OrganizationName()
		grantedOrganization := Instance.CreateOrganization(CTX, grantedOrganizationName, integration.Email())
		// create project grant
		_, err := ProjectClient.CreateProjectGrant(CTX, &project_v2beta.CreateProjectGrantRequest{
			ProjectId:             projectID,
			RoleKeys:              []string{role1, role2},
			GrantedOrganizationId: grantedOrganization.OrganizationId,
		})
		require.NoError(t, err)
		// create user
		user := Instance.CreateHumanUserVerified(CTX, grantedOrganization.OrganizationId, integration.Email(), integration.Phone())
		// create authorization with roles
		createdAuthorization := createAndEnsureAuthorization(t, instanceID, grantedOrganization.OrganizationId, user.UserId, projectID, []string{role1, role2}, retryDuration, tick)
		// update project grant
		_, err = ProjectClient.UpdateProjectGrant(CTX, &project_v2beta.UpdateProjectGrantRequest{
			ProjectId:             projectID,
			GrantedOrganizationId: grantedOrganization.OrganizationId,
			RoleKeys:              []string{role1},
		})
		require.NoError(t, err)

		// ensure authorization is updated
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, []string{role1}, az.Roles)
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.CreatedAt.UTC())
			// (@grvijayan) todo: should we check the timestamp here?
			// the authorization/authorization_roles RT is updated based on cascade deletions in the project_grant_roles table,
			// so the project grant update time cannot be used here
		}, retryDuration, tick, "authorization not updated within %v: %v", retryDuration, err)
	})

	t.Run("org removed reduces", func(t *testing.T) {
		// create a new organization
		orgName := integration.OrganizationName()
		orgResp := Instance.CreateOrganization(CTX, orgName, integration.Email())
		// create user
		user := Instance.CreateHumanUserVerified(CTX, orgResp.OrganizationId, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectID := prepareProjectAndProjectRoles(t, orgResp.OrganizationId, []string{role1, role2})
		// create authorization with roles
		createdAuthorization := createAndEnsureAuthorization(t, instanceID, orgResp.OrganizationId, user.UserId, projectID, []string{role1, role2}, retryDuration, tick)
		// delete the organization
		_, err := OrgClient.DeleteOrganization(CTX, &org_v2beta.DeleteOrganizationRequest{
			Id: orgResp.OrganizationId,
		})
		require.NoError(t, err)

		// ensure authorization is deleted
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})

	t.Run("instance removed reduces", func(t *testing.T) {
		// create a new instance
		instance := integration.NewInstance(CTX)
		// create a new organization
		orgResp := instance.CreateOrganization(CTX, integration.OrganizationName(), integration.Email())
		// create a user
		user := instance.CreateHumanUserVerified(CTX, orgResp.OrganizationId, integration.Email(), integration.Phone())
		// prepare project and project roles
		role1, role2 := "role1", "role2"
		projectResp := instance.CreateProject(CTX, t, orgResp.OrganizationId, integration.ProjectName(), false, false)
		_ = instance.AddProjectRole(CTX, t, projectResp.Id, role1, "display", "")
		_ = instance.AddProjectRole(CTX, t, projectResp.Id, role2, "display", "")

		// create authorization with roles
		createdAuthorization, err := instance.Client.AuthorizationV2.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
			UserId:         user.UserId,
			ProjectId:      projectResp.Id,
			RoleKeys:       []string{role1, role2},
			OrganizationId: orgResp.OrganizationId,
		})
		require.NoError(t, err)

		// ensure authorization exists
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instance.ID(), createdAuthorization.Id),
			))
			require.NoError(collect, err)
			assert.Equal(collect, []string{role1, role2}, az.Roles)
			assert.Equal(collect, domain.AuthorizationStateActive, az.State)
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.CreatedAt.UTC())
			assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime(), az.UpdatedAt.UTC())
		}, retryDuration, tick, "authorization not found within %v: %v", retryDuration, err)

		// delete the instance
		_, err = Instance.Client.InstanceV2.DeleteInstance(CTX, &instance_v2.DeleteInstanceRequest{
			InstanceId: instance.ID(),
		})
		require.NoError(t, err)

		// ensure authorization is deleted
		require.EventuallyWithT(t, func(collect *assert.CollectT) {
			az, err := authorizationRepo.Get(CTX, pool, database.WithCondition(
				authorizationRepo.PrimaryKeyCondition(instance.ID(), createdAuthorization.Id),
			))
			require.Empty(collect, az)
			require.ErrorIs(collect, err, new(database.NoRowFoundError))
		}, retryDuration, tick, "authorization not deleted within %v: %v", retryDuration, err)
	})
}

func createAndEnsureAuthorization(t *testing.T, instanceID, orgID, userID, projectID string, roles []string, retryDuration time.Duration, tick time.Duration) *authorization_v2.CreateAuthorizationResponse {
	t.Helper()

	// create authorization
	createdAuthorization, err := AuthorizationClient.CreateAuthorization(CTX, &authorization_v2.CreateAuthorizationRequest{
		UserId:         userID,
		ProjectId:      projectID,
		RoleKeys:       roles,
		OrganizationId: orgID,
	})
	require.NoError(t, err)
	// ensure authorization exists
	authzRepo := repository.AuthorizationRepository()
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		az, err := authzRepo.Get(CTX, pool, database.WithCondition(
			authzRepo.PrimaryKeyCondition(instanceID, createdAuthorization.Id),
		))
		require.NoError(collect, err)
		assert.Equal(collect, roles, az.Roles)
		assert.Equal(collect, domain.AuthorizationStateActive, az.State)
		assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime().UTC(), az.CreatedAt.UTC())
		assert.Equal(collect, createdAuthorization.GetCreationDate().AsTime().UTC(), az.UpdatedAt.UTC())
	}, retryDuration, tick, "authorization %q not found within %v: %v", createdAuthorization.Id, retryDuration, err)
	return createdAuthorization
}

func prepareProjectAndProjectRoles(t *testing.T, orgID string, roles []string) string {
	t.Helper()

	project, err := ProjectClient.CreateProject(CTX, &project_v2beta.CreateProjectRequest{
		OrganizationId: orgID,
		Name:           integration.ProjectName(),
	})
	require.NoError(t, err)

	if len(roles) == 0 {
		return project.Id
	}

	// Wait for the project to be created in the relational db before adding roles
	projectRepo := repository.ProjectRepository()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, err := projectRepo.Get(CTX, pool, database.WithCondition(
			projectRepo.PrimaryKeyCondition(Instance.ID(), project.Id),
		))
		require.NoError(collect, err)
	}, retryDuration, tick, "project %q not found within %v: %v", project.Id, retryDuration, err)

	for _, role := range roles {
		_, err = ProjectClient.AddProjectRole(CTX, &project_v2beta.AddProjectRoleRequest{
			ProjectId:   project.GetId(),
			RoleKey:     role,
			DisplayName: "display",
		})
		require.NoError(t, err)
	}
	return project.Id
}
