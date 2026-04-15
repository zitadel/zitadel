package projection

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestAdministratorRelationalReducers(t *testing.T) {
	handler := new(relationalTablesProjection)
	rawTx, tx := getTransactions(t)
	t.Cleanup(func() {
		require.NoError(t, rawTx.Rollback())
	})

	instanceID, _, orgID, projectID, grantID := seedAdministratorRelationalState(t, tx)
	adminRepo := repository.AdministratorRepository()
	userRepo := repository.UserRepository()

	var userSeq int
	createUser := func(t *testing.T) string {
		t.Helper()
		userSeq++
		userID := fmt.Sprintf("admin-test-user-%d", userSeq)
		err := userRepo.Create(t.Context(), tx, &repoDomain.User{
			InstanceID:     instanceID,
			OrganizationID: orgID,
			ID:             userID,
			Username:       userID,
			State:          repoDomain.UserStateActive,
			Machine: &repoDomain.MachineUser{
				Name: "machine",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		require.NoError(t, err)
		return userID
	}

	t.Run("'instance member added' reducer should create instance administrator", func(t *testing.T) {
		userID := createUser(t)

		added := instance.NewMemberAddedEvent(t.Context(), &instance.NewAggregate(instanceID).Aggregate, userID, "IAM_OWNER")
		require.True(t, callReduce(t, rawTx, handler, added))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.InstanceAdministratorCondition(instanceID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"IAM_OWNER"}, admin.Roles)
	})

	t.Run("'instance member changed' reducer should update instance administrator roles", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID: instanceID,
			UserID:     userID,
			Scope:      repoDomain.AdministratorScopeInstance,
			Roles:      []string{"IAM_OWNER"},
		})
		require.NoError(t, err)

		changed := instance.NewMemberChangedEvent(t.Context(), &instance.NewAggregate(instanceID).Aggregate, userID, "IAM_OWNER_VIEWER")
		require.True(t, callReduce(t, rawTx, handler, changed))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.InstanceAdministratorCondition(instanceID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"IAM_OWNER_VIEWER"}, admin.Roles)
	})

	t.Run("'instance member removed' reducer should remove instance administrator", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID: instanceID,
			UserID:     userID,
			Scope:      repoDomain.AdministratorScopeInstance,
			Roles:      []string{"IAM_OWNER"},
		})
		require.NoError(t, err)

		removed := instance.NewMemberRemovedEvent(t.Context(), &instance.NewAggregate(instanceID).Aggregate, userID)
		require.True(t, callReduce(t, rawTx, handler, removed))

		_, err = adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.InstanceAdministratorCondition(instanceID, userID),
		))
		require.ErrorIs(t, err, database.NewNoRowFoundError(nil))
	})

	t.Run("'org member added' reducer should create organization administrator", func(t *testing.T) {
		userID := createUser(t)

		orgAggregate := org.NewAggregate(orgID)
		orgAggregate.InstanceID = instanceID

		added := org.NewMemberAddedEvent(t.Context(), &orgAggregate.Aggregate, userID, "ORG_OWNER")
		require.True(t, callReduce(t, rawTx, handler, added))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.OrganizationAdministratorCondition(instanceID, orgID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"ORG_OWNER"}, admin.Roles)
	})

	t.Run("'org member changed' reducer should update organization administrator roles", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID:     instanceID,
			UserID:         userID,
			Scope:          repoDomain.AdministratorScopeOrganization,
			OrganizationID: &orgID,
			Roles:          []string{"ORG_OWNER"},
		})
		require.NoError(t, err)

		orgAggregate := org.NewAggregate(orgID)
		orgAggregate.InstanceID = instanceID

		changed := org.NewMemberChangedEvent(t.Context(), &orgAggregate.Aggregate, userID, "ORG_OWNER_VIEWER")
		require.True(t, callReduce(t, rawTx, handler, changed))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.OrganizationAdministratorCondition(instanceID, orgID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"ORG_OWNER_VIEWER"}, admin.Roles)
	})

	t.Run("'org member removed' reducer should remove organization administrator", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID:     instanceID,
			UserID:         userID,
			Scope:          repoDomain.AdministratorScopeOrganization,
			OrganizationID: &orgID,
			Roles:          []string{"ORG_OWNER"},
		})
		require.NoError(t, err)

		orgAggregate := org.NewAggregate(orgID)
		orgAggregate.InstanceID = instanceID

		removed := org.NewMemberRemovedEvent(t.Context(), &orgAggregate.Aggregate, userID)
		require.True(t, callReduce(t, rawTx, handler, removed))

		_, err = adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.OrganizationAdministratorCondition(instanceID, orgID, userID),
		))
		require.ErrorIs(t, err, database.NewNoRowFoundError(nil))
	})

	t.Run("'project member added' reducer should create project administrator", func(t *testing.T) {
		userID := createUser(t)

		projectAggregate := project.NewAggregate(projectID, orgID)
		projectAggregate.InstanceID = instanceID

		added := project.NewProjectMemberAddedEvent(t.Context(), &projectAggregate.Aggregate, userID, "PROJECT_OWNER")
		require.True(t, callReduce(t, rawTx, handler, added))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.ProjectAdministratorCondition(instanceID, projectID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"PROJECT_OWNER"}, admin.Roles)
	})

	t.Run("'project member changed' reducer should update project administrator roles", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID: instanceID,
			UserID:     userID,
			Scope:      repoDomain.AdministratorScopeProject,
			ProjectID:  &projectID,
			Roles:      []string{"PROJECT_OWNER"},
		})
		require.NoError(t, err)

		projectAggregate := project.NewAggregate(projectID, orgID)
		projectAggregate.InstanceID = instanceID

		changed := project.NewProjectMemberChangedEvent(t.Context(), &projectAggregate.Aggregate, userID, "PROJECT_OWNER_VIEWER")
		require.True(t, callReduce(t, rawTx, handler, changed))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.ProjectAdministratorCondition(instanceID, projectID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"PROJECT_OWNER_VIEWER"}, admin.Roles)
	})

	t.Run("'project member removed' reducer should remove project administrator", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID: instanceID,
			UserID:     userID,
			Scope:      repoDomain.AdministratorScopeProject,
			ProjectID:  &projectID,
			Roles:      []string{"PROJECT_OWNER"},
		})
		require.NoError(t, err)

		projectAggregate := project.NewAggregate(projectID, orgID)
		projectAggregate.InstanceID = instanceID

		removed := project.NewProjectMemberRemovedEvent(t.Context(), &projectAggregate.Aggregate, userID)
		require.True(t, callReduce(t, rawTx, handler, removed))

		_, err = adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.ProjectAdministratorCondition(instanceID, projectID, userID),
		))
		require.ErrorIs(t, err, database.NewNoRowFoundError(nil))
	})

	t.Run("'project grant member added' reducer should create project grant administrator", func(t *testing.T) {
		userID := createUser(t)

		projectAggregate := project.NewAggregate(projectID, orgID)
		projectAggregate.InstanceID = instanceID

		added := project.NewProjectGrantMemberAddedEvent(t.Context(), &projectAggregate.Aggregate, userID, grantID, "PROJECT_OWNER")
		require.True(t, callReduce(t, rawTx, handler, added))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.ProjectGrantAdministratorCondition(instanceID, grantID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"PROJECT_OWNER"}, admin.Roles)
	})

	t.Run("'project grant member changed' reducer should update project grant administrator roles", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID:     instanceID,
			UserID:         userID,
			Scope:          repoDomain.AdministratorScopeProjectGrant,
			ProjectGrantID: &grantID,
			Roles:          []string{"PROJECT_OWNER"},
		})
		require.NoError(t, err)

		projectAggregate := project.NewAggregate(projectID, orgID)
		projectAggregate.InstanceID = instanceID

		changed := project.NewProjectGrantMemberChangedEvent(t.Context(), &projectAggregate.Aggregate, userID, grantID, "PROJECT_OWNER_VIEWER")
		require.True(t, callReduce(t, rawTx, handler, changed))

		admin, err := adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.ProjectGrantAdministratorCondition(instanceID, grantID, userID),
		))
		require.NoError(t, err)
		assert.Equal(t, []string{"PROJECT_OWNER_VIEWER"}, admin.Roles)
	})

	t.Run("'project grant member cascade removed' reducer should remove project grant administrator", func(t *testing.T) {
		userID := createUser(t)
		err := adminRepo.Create(t.Context(), tx, &repoDomain.Administrator{
			InstanceID:     instanceID,
			UserID:         userID,
			Scope:          repoDomain.AdministratorScopeProjectGrant,
			ProjectGrantID: &grantID,
			Roles:          []string{"PROJECT_OWNER"},
		})
		require.NoError(t, err)

		projectAggregate := project.NewAggregate(projectID, orgID)
		projectAggregate.InstanceID = instanceID

		removed := project.NewProjectGrantMemberCascadeRemovedEvent(t.Context(), &projectAggregate.Aggregate, userID, grantID)
		require.True(t, callReduce(t, rawTx, handler, removed))

		_, err = adminRepo.Get(t.Context(), tx, database.WithCondition(
			adminRepo.ProjectGrantAdministratorCondition(instanceID, grantID, userID),
		))
		require.ErrorIs(t, err, database.NewNoRowFoundError(nil))
	})
}

func seedAdministratorRelationalState(t *testing.T, tx database.QueryExecutor) (instanceID, userID, orgID, projectID, projectGrantID string) {
	t.Helper()

	instanceRepo := repository.InstanceRepository()
	orgRepo := repository.OrganizationRepository()
	projectRepo := repository.ProjectRepository()
	projectGrantRepo := repository.ProjectGrantRepository()
	userRepo := repository.UserRepository()

	now := time.Now().UnixNano()
	instanceID = fmt.Sprintf("instance-%d", now)
	err := instanceRepo.Create(t.Context(), tx, &repoDomain.Instance{
		ID:              instanceID,
		Name:            "instance",
		DefaultOrgID:    "default-org",
		IAMProjectID:    "iam-project",
		ConsoleClientID: "console-client",
		ConsoleAppID:    "console-app",
		DefaultLanguage: "en",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	require.NoError(t, err)

	orgID = fmt.Sprintf("org-%d", now+1)
	err = orgRepo.Create(t.Context(), tx, &repoDomain.Organization{
		InstanceID: instanceID,
		ID:         orgID,
		Name:       "org",
		State:      repoDomain.OrgStateActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	require.NoError(t, err)

	projectID = fmt.Sprintf("project-%d", now+2)
	err = projectRepo.Create(t.Context(), tx, &repoDomain.Project{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             projectID,
		Name:           "project",
		State:          repoDomain.ProjectStateActive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	require.NoError(t, err)

	projectGrantID = fmt.Sprintf("grant-%d", now+3)
	err = projectGrantRepo.Create(t.Context(), tx, &repoDomain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     projectGrantID,
		ProjectID:              projectID,
		GrantingOrganizationID: orgID,
		GrantedOrganizationID:  orgID,
		State:                  repoDomain.ProjectGrantStateActive,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	})
	require.NoError(t, err)

	userID = fmt.Sprintf("user-%d", now+4)
	err = userRepo.Create(t.Context(), tx, &repoDomain.User{
		InstanceID:     instanceID,
		OrganizationID: orgID,
		ID:             userID,
		Username:       userID,
		State:          repoDomain.UserStateActive,
		Machine: &repoDomain.MachineUser{
			Name: "machine",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	require.NoError(t, err)

	return instanceID, userID, orgID, projectID, projectGrantID
}
