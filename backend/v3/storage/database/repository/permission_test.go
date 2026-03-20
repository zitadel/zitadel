package repository_test

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestCheckPermission_DB(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	otherOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	projectGrantID := createProjectGrant(t, tx, instanceID, orgID, otherOrgID, projectID, nil)

	instanceUserID := createMachineUser(t, tx, instanceID, orgID)
	organizationUserID := createMachineUser(t, tx, instanceID, orgID)
	projectUserID := createMachineUser(t, tx, instanceID, orgID)
	projectGrantUserID := createMachineUser(t, tx, instanceID, orgID)
	noneUserID := createMachineUser(t, tx, instanceID, orgID)

	adminRoleRepo := repository.AdministratorRoleRepository()
	_, err := adminRoleRepo.AddPermissions(t.Context(), tx, instanceID, "instance_admin", "instance.read")
	require.NoError(t, err)
	_, err = adminRoleRepo.AddPermissions(t.Context(), tx, instanceID, "org_admin", "organization.read", "project.read")
	require.NoError(t, err)
	_, err = adminRoleRepo.AddPermissions(t.Context(), tx, instanceID, "project_admin", "project.read")
	require.NoError(t, err)
	_, err = adminRoleRepo.AddPermissions(t.Context(), tx, instanceID, "project_grant_admin", "project_grant.read")
	require.NoError(t, err)

	adminRepo := repository.AdministratorRepository()
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID,
		UserID:     instanceUserID,
		Scope:      domain.AdministratorScopeInstance,
		Roles:      []string{"instance_admin"},
	}))
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID:     instanceID,
		UserID:         organizationUserID,
		Scope:          domain.AdministratorScopeOrganization,
		OrganizationID: gu.Ptr(orgID),
		Roles:          []string{"org_admin"},
	}))
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID: instanceID,
		UserID:     projectUserID,
		Scope:      domain.AdministratorScopeProject,
		ProjectID:  gu.Ptr(projectID),
		Roles:      []string{"project_admin"},
	}))
	require.NoError(t, adminRepo.Create(t.Context(), tx, &domain.Administrator{
		InstanceID:     instanceID,
		UserID:         projectGrantUserID,
		Scope:          domain.AdministratorScopeProjectGrant,
		ProjectGrantID: gu.Ptr(projectGrantID),
		Roles:          []string{"project_grant_admin"},
	}))

	tests := []struct {
		name       string
		userID     string
		permission string
		opts       []repository.CheckPermissionOpt
		want       bool
		wantErr    bool
	}{
		{
			name:       "instance scope allows",
			userID:     instanceUserID,
			permission: "instance.read",
			want:       true,
		},
		{
			name:       "instance scope denies for org-only user without context",
			userID:     organizationUserID,
			permission: "instance.read",
			want:       false,
		},
		{
			name:       "organization scope allows on matching org",
			userID:     organizationUserID,
			permission: "organization.read",
			opts:       []repository.CheckPermissionOpt{repository.WithOrganizationID(orgID)},
			want:       true,
		},
		{
			name:       "organization scope denies on other org",
			userID:     organizationUserID,
			permission: "organization.read",
			opts:       []repository.CheckPermissionOpt{repository.WithOrganizationID(otherOrgID)},
			want:       false,
		},
		{
			name:       "project scope allows on matching project",
			userID:     projectUserID,
			permission: "project.read",
			opts:       []repository.CheckPermissionOpt{repository.WithProjectID(projectID)},
			want:       true,
		},
		{
			name:       "project scope denies on other project",
			userID:     projectUserID,
			permission: "project.read",
			opts:       []repository.CheckPermissionOpt{repository.WithProjectID("project-other")},
			want:       false,
		},
		{
			name:       "project grant scope allows on matching project grant",
			userID:     projectGrantUserID,
			permission: "project_grant.read",
			opts:       []repository.CheckPermissionOpt{repository.WithProjectGrantID(projectGrantID)},
			want:       true,
		},
		{
			name:       "project grant scope denies on other project grant",
			userID:     projectGrantUserID,
			permission: "project_grant.read",
			opts:       []repository.CheckPermissionOpt{repository.WithProjectGrantID("project-grant-other")},
			want:       false,
		},
		{
			name:       "organization inheritance allows project when org context provided",
			userID:     organizationUserID,
			permission: "project.read",
			opts:       []repository.CheckPermissionOpt{repository.WithOrganizationID(orgID)},
			want:       true,
		},
		{
			name:       "organization inheritance allows project when org and project context provided",
			userID:     organizationUserID,
			permission: "project.read",
			opts: []repository.CheckPermissionOpt{
				repository.WithOrganizationID(orgID),
				repository.WithProjectID(projectID),
			},
			want: true,
		},
		{
			name:       "raise if denied returns error",
			userID:     noneUserID,
			permission: "instance.read",
			opts:       []repository.CheckPermissionOpt{repository.WithRaiseIfDenied()},
			wantErr:    true,
		},
		{
			name:       "raise if denied returns error because user does not exist",
			userID:     "non-existing-user",
			permission: "instance.read",
			opts:       []repository.CheckPermissionOpt{repository.WithRaiseIfDenied()},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, err := tx.Begin(t.Context())
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, savepoint.Rollback(t.Context()))
			}()

			got, err := executeCheckPermission(t, savepoint, instanceID, tt.userID, tt.permission, tt.opts...)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "Permission denied")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func executeCheckPermission(t *testing.T, executor database.QueryExecutor, instanceID, userID, permission string, opts ...repository.CheckPermissionOpt) (bool, error) {
	t.Helper()
	condition := repository.PermissionCondition(instanceID, userID, permission, opts...)
	builder := database.NewStatementBuilder("SELECT ")
	condition.Write(builder)

	var got bool
	err := executor.QueryRow(t.Context(), builder.String(), builder.Args()...).Scan(&got)
	return got, err
}
