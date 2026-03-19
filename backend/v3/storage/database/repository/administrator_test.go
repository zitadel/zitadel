package repository_test

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
)

func TestAdministratorRepository_Create(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	projectGrantID := createProjectGrant(t, tx, instanceID, orgID, grantedOrgID, projectID, nil)
	userID := createHumanUser(t, tx, instanceID, orgID)
	adminRepo := repository.AdministratorRepository()

	// Pre-create an admin for duplicate detection.
	duplicateUser := createHumanUser(t, tx, instanceID, orgID)
	err := adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID:     instanceID,
		UserID:         duplicateUser,
		Scope:          domain.AdministratorScopeOrganization,
		OrganizationID: gu.Ptr(orgID),
		Roles:          []string{"ORG_OWNER"},
	})
	require.NoError(t, err)

	tests := []struct {
		name          string
		administrator *domain.Administrator
		wantErr       error
	}{
		{
			name: "instance scope",
			administrator: &domain.Administrator{
				InstanceID: instanceID,
				UserID:     userID,
				Scope:      domain.AdministratorScopeInstance,
				Roles:      []string{"IAM_OWNER"},
			},
		},
		{
			name: "organization scope",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         createHumanUser(t, tx, instanceID, orgID),
				Scope:          domain.AdministratorScopeOrganization,
				OrganizationID: gu.Ptr(orgID),
				Roles:          []string{"ORG_OWNER"},
			},
		},
		{
			name: "project scope",
			administrator: &domain.Administrator{
				InstanceID: instanceID,
				UserID:     userID,
				Scope:      domain.AdministratorScopeProject,
				ProjectID:  gu.Ptr(projectID),
				Roles:      []string{"PROJECT_OWNER"},
			},
		},
		{
			name: "project grant scope",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         createHumanUser(t, tx, instanceID, orgID),
				Scope:          domain.AdministratorScopeProjectGrant,
				ProjectGrantID: gu.Ptr(projectGrantID),
				Roles:          []string{"PROJECT_OWNER"},
			},
		},
		{
			name: "duplicate admin",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         duplicateUser,
				Scope:          domain.AdministratorScopeOrganization,
				OrganizationID: gu.Ptr(orgID),
				Roles:          []string{"ORG_OWNER"},
			},
			wantErr: new(database.UniqueError),
		},
		{
			name: "missing user",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         integration.ID(),
				Scope:          domain.AdministratorScopeOrganization,
				OrganizationID: gu.Ptr(orgID),
				Roles:          []string{"ORG_OWNER"},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "missing organization",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         createHumanUser(t, tx, instanceID, orgID),
				Scope:          domain.AdministratorScopeOrganization,
				OrganizationID: gu.Ptr(integration.ID()),
				Roles:          []string{"ORG_OWNER"},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "missing project",
			administrator: &domain.Administrator{
				InstanceID: instanceID,
				UserID:     userID,
				Scope:      domain.AdministratorScopeProject,
				ProjectID:  gu.Ptr(integration.ID()),
				Roles:      []string{"PROJECT_OWNER"},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "missing project grant",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         createHumanUser(t, tx, instanceID, orgID),
				Scope:          domain.AdministratorScopeProjectGrant,
				ProjectGrantID: gu.Ptr(integration.ID()),
				Roles:          []string{"PROJECT_OWNER"},
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "invalid scope alignment",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         createHumanUser(t, tx, instanceID, orgID),
				Scope:          domain.AdministratorScopeProjectGrant,
				ProjectID:      gu.Ptr(projectID),
				ProjectGrantID: gu.Ptr(projectGrantID),
				Roles:          []string{"PROJECT_OWNER"},
			},
			wantErr: new(database.CheckError),
		},
		{
			name: "empty role",
			administrator: &domain.Administrator{
				InstanceID:     instanceID,
				UserID:         createHumanUser(t, tx, instanceID, orgID),
				Scope:          domain.AdministratorScopeOrganization,
				OrganizationID: gu.Ptr(orgID),
				Roles:          []string{""},
			},
			wantErr: new(database.CheckError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			err := adminRepo.Create(t.Context(), savepoint, tt.administrator)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, tt.administrator.ID)
			assert.NotZero(t, tt.administrator.CreatedAt)
			assert.NotZero(t, tt.administrator.UpdatedAt)

			got, err := adminRepo.Get(t.Context(), savepoint, database.WithCondition(
				adminRepo.PrimaryKeyCondition(instanceID, tt.administrator.ID),
			))
			require.NoError(t, err)
			assert.Equal(t, tt.administrator.InstanceID, got.InstanceID)
			assert.Equal(t, tt.administrator.UserID, got.UserID)
			assert.Equal(t, tt.administrator.Scope, got.Scope)
			assert.Equal(t, tt.administrator.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.administrator.ProjectID, got.ProjectID)
			assert.Equal(t, tt.administrator.ProjectGrantID, got.ProjectGrantID)
			assert.ElementsMatch(t, tt.administrator.Roles, got.Roles)
		})
	}
}

func TestAdministratorRepository_Get(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	projectGrantID := createProjectGrant(t, tx, instanceID, orgID, grantedOrgID, projectID, nil)
	userID := createHumanUser(t, tx, instanceID, orgID)
	adminRepo := repository.AdministratorRepository()

	instanceAdmin := &domain.Administrator{
		InstanceID: instanceID,
		UserID:     userID,
		Scope:      domain.AdministratorScopeInstance,
		Roles:      []string{"IAM_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, instanceAdmin))

	orgAdmin := &domain.Administrator{
		InstanceID:     instanceID,
		UserID:         userID,
		Scope:          domain.AdministratorScopeOrganization,
		OrganizationID: gu.Ptr(orgID),
		Roles:          []string{"ORG_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, orgAdmin))

	projectAdmin := &domain.Administrator{
		InstanceID: instanceID,
		UserID:     userID,
		Scope:      domain.AdministratorScopeProject,
		ProjectID:  gu.Ptr(projectID),
		Roles:      []string{"PROJECT_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, projectAdmin))

	projectGrantAdmin := &domain.Administrator{
		InstanceID:     instanceID,
		UserID:         userID,
		Scope:          domain.AdministratorScopeProjectGrant,
		ProjectGrantID: gu.Ptr(projectGrantID),
		Roles:          []string{"PROJECT_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, projectGrantAdmin))

	tests := []struct {
		name      string
		condition database.Condition
		wantAdmin *domain.Administrator
		wantErr   error
	}{
		{
			name:      "instance scope",
			condition: adminRepo.InstanceAdministratorCondition(instanceID, userID),
			wantAdmin: instanceAdmin,
		},
		{
			name:      "organization scope",
			condition: adminRepo.OrganizationAdministratorCondition(instanceID, orgID, userID),
			wantAdmin: orgAdmin,
		},
		{
			name:      "project scope",
			condition: adminRepo.ProjectAdministratorCondition(instanceID, projectID, userID),
			wantAdmin: projectAdmin,
		},
		{
			name:      "project grant scope",
			condition: adminRepo.ProjectGrantAdministratorCondition(instanceID, projectGrantID, userID),
			wantAdmin: projectGrantAdmin,
		},
		{
			name:      "by primary key",
			condition: adminRepo.PrimaryKeyCondition(instanceID, orgAdmin.ID),
			wantAdmin: orgAdmin,
		},
		{
			name:      "not found",
			condition: adminRepo.InstanceAdministratorCondition(instanceID, integration.ID()),
			wantErr:   database.NewNoRowFoundError(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			got, err := adminRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAdmin.ID, got.ID)
			assert.Equal(t, tt.wantAdmin.Scope, got.Scope)
			assert.Equal(t, tt.wantAdmin.UserID, got.UserID)
			assert.Equal(t, tt.wantAdmin.OrganizationID, got.OrganizationID)
			assert.Equal(t, tt.wantAdmin.ProjectID, got.ProjectID)
			assert.Equal(t, tt.wantAdmin.ProjectGrantID, got.ProjectGrantID)
			assert.ElementsMatch(t, tt.wantAdmin.Roles, got.Roles)
		})
	}
}

func TestAdministratorRepository_List(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	otherOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	projectGrantID := createProjectGrant(t, tx, instanceID, orgID, grantedOrgID, projectID, nil)
	user1 := createHumanUser(t, tx, instanceID, orgID)
	user2 := createHumanUser(t, tx, instanceID, orgID)
	adminRepo := repository.AdministratorRepository()

	// instance: user1=IAM_OWNER, user2=IAM_VIEWER
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeInstance, Roles: []string{"IAM_OWNER"},
	}))
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID, UserID: user2,
		Scope: domain.AdministratorScopeInstance, Roles: []string{"IAM_VIEWER"},
	}))
	// organization: user1 in orgID
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeOrganization, OrganizationID: gu.Ptr(orgID),
		Roles: []string{"ORG_OWNER"},
	}))
	// project: user1
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeProject, ProjectID: gu.Ptr(projectID),
		Roles: []string{"PROJECT_OWNER"},
	}))
	// project grant: user1
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeProjectGrant, ProjectGrantID: gu.Ptr(projectGrantID),
		Roles: []string{"PROJECT_OWNER"},
	}))

	tests := []struct {
		name      string
		condition database.Condition
		wantLen   int
	}{
		{
			name: "all instance admins",
			condition: database.And(
				adminRepo.InstanceIDCondition(instanceID),
				adminRepo.ScopeCondition(domain.AdministratorScopeInstance),
			),
			wantLen: 2,
		},
		{
			name: "instance admins filtered by role",
			condition: database.And(
				adminRepo.InstanceIDCondition(instanceID),
				adminRepo.ScopeCondition(domain.AdministratorScopeInstance),
				adminRepo.RoleCondition(database.TextOperationEqual, "IAM_OWNER"),
			),
			wantLen: 1,
		},
		{
			name: "org admins in org",
			condition: database.And(
				adminRepo.InstanceIDCondition(instanceID),
				adminRepo.ScopeCondition(domain.AdministratorScopeOrganization),
				adminRepo.OrganizationIDCondition(orgID),
			),
			wantLen: 1,
		},
		{
			name: "org admins in other org returns empty",
			condition: database.And(
				adminRepo.InstanceIDCondition(instanceID),
				adminRepo.ScopeCondition(domain.AdministratorScopeOrganization),
				adminRepo.OrganizationIDCondition(otherOrgID),
			),
			wantLen: 0,
		},
		{
			name: "project admins",
			condition: database.And(
				adminRepo.InstanceIDCondition(instanceID),
				adminRepo.ScopeCondition(domain.AdministratorScopeProject),
				adminRepo.ProjectIDCondition(projectID),
			),
			wantLen: 1,
		},
		{
			name: "project grant admins",
			condition: database.And(
				adminRepo.InstanceIDCondition(instanceID),
				adminRepo.ScopeCondition(domain.AdministratorScopeProjectGrant),
				adminRepo.ProjectGrantIDCondition(projectGrantID),
			),
			wantLen: 1,
		},
		{
			name:      "scope isolation: instance condition excludes org admin",
			condition: adminRepo.InstanceAdministratorCondition(instanceID, user1),
			wantLen:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			list, err := adminRepo.List(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.NoError(t, err)
			assert.Len(t, list, tt.wantLen)
		})
	}
}

func TestAdministratorRepository_Update(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	userID := createHumanUser(t, tx, instanceID, orgID)
	adminRepo := repository.AdministratorRepository()

	admin := &domain.Administrator{
		InstanceID:     instanceID,
		UserID:         userID,
		Scope:          domain.AdministratorScopeOrganization,
		OrganizationID: gu.Ptr(orgID),
		Roles:          []string{"ORG_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, admin))
	originalCreatedAt := admin.CreatedAt

	condition := adminRepo.OrganizationAdministratorCondition(instanceID, orgID, userID)

	tests := []struct {
		name    string
		changes []database.Change
		verify  func(t *testing.T, got *domain.Administrator)
	}{
		{
			name:    "SetUpdatedAt",
			changes: []database.Change{adminRepo.SetUpdatedAt(time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC))},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Equal(t, time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC), got.UpdatedAt.UTC())
				assert.Equal(t, []string{"ORG_OWNER"}, got.Roles)
				assert.Equal(t, originalCreatedAt, got.CreatedAt, "created_at must not change")
			},
		},
		{
			name:    "SetRoles replaces all",
			changes: []database.Change{adminRepo.SetRoles([]string{"AUDITOR", "VIEWER"})},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.ElementsMatch(t, []string{"AUDITOR", "VIEWER"}, got.Roles)
			},
		},
		{
			name:    "SetRoles to single role",
			changes: []database.Change{adminRepo.SetRoles([]string{"AUDITOR"})},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Equal(t, []string{"AUDITOR"}, got.Roles)
			},
		},
		{
			name:    "SetRoles to empty removes all",
			changes: []database.Change{adminRepo.SetRoles([]string{})},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Empty(t, got.Roles)
			},
		},
		{
			name:    "AddRole new",
			changes: []database.Change{adminRepo.AddRole("AUDITOR")},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.ElementsMatch(t, []string{"ORG_OWNER", "AUDITOR"}, got.Roles)
			},
		},
		{
			name:    "AddRole existing is idempotent",
			changes: []database.Change{adminRepo.AddRole("ORG_OWNER")},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Equal(t, []string{"ORG_OWNER"}, got.Roles)
			},
		},
		{
			name:    "RemoveRole existing",
			changes: []database.Change{adminRepo.RemoveRole("ORG_OWNER")},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Empty(t, got.Roles)
			},
		},
		{
			name:    "RemoveRole absent is idempotent",
			changes: []database.Change{adminRepo.RemoveRole("NONEXISTENT")},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Equal(t, []string{"ORG_OWNER"}, got.Roles)
			},
		},
		{
			name: "SetUpdatedAt and SetRoles combined",
			changes: []database.Change{
				adminRepo.SetUpdatedAt(time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)),
				adminRepo.SetRoles([]string{"AUDITOR"}),
			},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Equal(t, time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC), got.UpdatedAt.UTC())
				assert.Equal(t, []string{"AUDITOR"}, got.Roles)
			},
		},
		{
			name: "AddRole and RemoveRole combined",
			changes: []database.Change{
				adminRepo.AddRole("AUDITOR"),
				adminRepo.RemoveRole("ORG_OWNER"),
			},
			verify: func(t *testing.T, got *domain.Administrator) {
				assert.Equal(t, []string{"AUDITOR"}, got.Roles)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			_, err := adminRepo.Update(t.Context(), savepoint, condition, tt.changes...)
			require.NoError(t, err)

			got, err := adminRepo.Get(t.Context(), savepoint, database.WithCondition(condition))
			require.NoError(t, err)
			tt.verify(t, got)
		})
	}
}

func TestAdministratorRepository_Delete(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	projectGrantID := createProjectGrant(t, tx, instanceID, orgID, grantedOrgID, projectID, nil)
	user1 := createHumanUser(t, tx, instanceID, orgID)
	user2 := createHumanUser(t, tx, instanceID, orgID)
	adminRepo := repository.AdministratorRepository()

	instanceAdmin := &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeInstance, Roles: []string{"IAM_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, instanceAdmin))

	otherInstanceAdmin := &domain.Administrator{
		InstanceID: instanceID, UserID: user2,
		Scope: domain.AdministratorScopeInstance, Roles: []string{"IAM_VIEWER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, otherInstanceAdmin))

	orgAdmin := &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeOrganization, OrganizationID: gu.Ptr(orgID),
		Roles: []string{"ORG_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, orgAdmin))

	projectAdmin := &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeProject, ProjectID: gu.Ptr(projectID),
		Roles: []string{"PROJECT_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, projectAdmin))

	projectGrantAdmin := &domain.Administrator{
		InstanceID: instanceID, UserID: user1,
		Scope: domain.AdministratorScopeProjectGrant, ProjectGrantID: gu.Ptr(projectGrantID),
		Roles: []string{"PROJECT_OWNER"},
	}
	require.NoError(t, adminRepo.Create(t.Context(), tx, projectGrantAdmin))

	tests := []struct {
		name            string
		condition       database.Condition
		verifyRemaining func(t *testing.T, sp database.Transaction)
	}{
		{
			name:      "instance admin, other scopes remain",
			condition: adminRepo.InstanceAdministratorCondition(instanceID, user1),
			verifyRemaining: func(t *testing.T, sp database.Transaction) {
				_, err := adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.InstanceAdministratorCondition(instanceID, user2),
				))
				require.NoError(t, err)
				_, err = adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.OrganizationAdministratorCondition(instanceID, orgID, user1),
				))
				require.NoError(t, err)
			},
		},
		{
			name:      "organization admin, other scopes remain",
			condition: adminRepo.OrganizationAdministratorCondition(instanceID, orgID, user1),
			verifyRemaining: func(t *testing.T, sp database.Transaction) {
				_, err := adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.InstanceAdministratorCondition(instanceID, user1),
				))
				require.NoError(t, err)
				_, err = adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.ProjectAdministratorCondition(instanceID, projectID, user1),
				))
				require.NoError(t, err)
			},
		},
		{
			name:      "project admin, other scopes remain",
			condition: adminRepo.ProjectAdministratorCondition(instanceID, projectID, user1),
			verifyRemaining: func(t *testing.T, sp database.Transaction) {
				_, err := adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.InstanceAdministratorCondition(instanceID, user1),
				))
				require.NoError(t, err)
				_, err = adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.ProjectGrantAdministratorCondition(instanceID, projectGrantID, user1),
				))
				require.NoError(t, err)
			},
		},
		{
			name:      "project grant admin, other scopes remain",
			condition: adminRepo.ProjectGrantAdministratorCondition(instanceID, projectGrantID, user1),
			verifyRemaining: func(t *testing.T, sp database.Transaction) {
				_, err := adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.InstanceAdministratorCondition(instanceID, user1),
				))
				require.NoError(t, err)
				_, err = adminRepo.Get(t.Context(), sp, database.WithCondition(
					adminRepo.OrganizationAdministratorCondition(instanceID, orgID, user1),
				))
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			_, err := adminRepo.Delete(t.Context(), savepoint, tt.condition)
			require.NoError(t, err)

			_, err = adminRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, database.NewNoRowFoundError(nil))

			tt.verifyRemaining(t, savepoint)
		})
	}
}
