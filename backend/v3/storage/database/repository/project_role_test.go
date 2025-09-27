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
)

func TestGetProjectRole(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	roleRepo := repository.ProjectRepository().Role()

	firstRole := &domain.ProjectRole{
		InstanceID:     instanceID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		Key:            "key1",
		DisplayName:    "Role 1",
		RoleGroup:      gu.Ptr("group1"),
	}
	err := roleRepo.Create(t.Context(), tx, firstRole)
	require.NoError(t, err)

	secondRole := &domain.ProjectRole{
		InstanceID:     instanceID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		Key:            "key2",
		DisplayName:    "Role 2",
		RoleGroup:      gu.Ptr("group2"),
	}
	err = roleRepo.Create(t.Context(), tx, secondRole)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		wantRole  *domain.ProjectRole
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: roleRepo.KeyCondition(firstRole.Key),
			wantErr:   database.NewMissingConditionError(roleRepo.InstanceIDColumn()),
		},
		{
			name:      "not found",
			condition: roleRepo.PrimaryKeyCondition(instanceID, projectID, "foo"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: database.And(roleRepo.InstanceIDCondition(instanceID), roleRepo.ProjectIDCondition(projectID)),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok",
			condition: roleRepo.PrimaryKeyCondition(instanceID, projectID, firstRole.Key),
			wantRole:  firstRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRole, err := roleRepo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRole, gotRole)
		})
	}
}

func TestListProjectRoles(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	firstProjectID := createProject(t, tx, instanceID, orgID)
	secondProjectID := createProject(t, tx, instanceID, orgID)
	roleRepo := repository.ProjectRepository().Role()

	roles := [...]*domain.ProjectRole{
		{
			InstanceID:     instanceID,
			ProjectID:      firstProjectID,
			OrganizationID: orgID,
			Key:            "key1",
			DisplayName:    "Role 1",
			RoleGroup:      gu.Ptr("group1"),
		},
		{
			InstanceID:     instanceID,
			ProjectID:      firstProjectID,
			OrganizationID: orgID,
			Key:            "key2",
			DisplayName:    "Role 2",
			RoleGroup:      gu.Ptr("group2"),
		},
		{
			InstanceID:     instanceID,
			ProjectID:      secondProjectID,
			OrganizationID: orgID,
			Key:            "key3",
			DisplayName:    "Role 3",
			RoleGroup:      gu.Ptr("group3"),
		},
		{
			InstanceID:     instanceID,
			ProjectID:      secondProjectID,
			OrganizationID: orgID,
			Key:            "key4",
			DisplayName:    "foobar",
			RoleGroup:      gu.Ptr("group3"),
		},
		{
			InstanceID:     instanceID,
			ProjectID:      secondProjectID,
			OrganizationID: orgID,
			Key:            "key5",
			DisplayName:    "foobaz",
			RoleGroup:      gu.Ptr("group4"),
		},
	}
	for _, r := range roles {
		err := roleRepo.Create(t.Context(), tx, r)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		wantRoles []*domain.ProjectRole
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: roleRepo.KeyCondition("key1"),
			wantErr:   database.NewMissingConditionError(roleRepo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: roleRepo.PrimaryKeyCondition(instanceID, firstProjectID, "foo"),
		},
		{
			name: "all from project 1",
			condition: database.And(
				roleRepo.InstanceIDCondition(instanceID),
				roleRepo.ProjectIDCondition(firstProjectID),
			),
			wantRoles: roles[0:2],
		},
		{
			name: "group 3 from project 2",
			condition: database.And(
				roleRepo.InstanceIDCondition(instanceID),
				roleRepo.ProjectIDCondition(secondProjectID),
				roleRepo.RoleGroupCondition(database.TextOperationEqual, "group3"),
			),
			wantRoles: roles[2:4],
		},
		{
			name: "name starts with 'foo' from project 2",
			condition: database.And(
				roleRepo.InstanceIDCondition(instanceID),
				roleRepo.ProjectIDCondition(secondProjectID),
				roleRepo.DisplayNameCondition(database.TextOperationStartsWith, "foo"),
			),
			wantRoles: roles[3:5],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoles, err := roleRepo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(roleRepo.PrimaryKeyColumns()...),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRoles, gotRoles)
		})
	}
}

func TestCreateProjectRole(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	roleRepo := repository.ProjectRepository().Role()

	existingRole := &domain.ProjectRole{
		InstanceID:     instanceID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		Key:            "key1",
		DisplayName:    "Role 1",
		RoleGroup:      gu.Ptr("group1"),
	}
	err := roleRepo.Create(t.Context(), tx, existingRole)
	require.NoError(t, err)

	tests := []struct {
		name    string
		role    *domain.ProjectRole
		wantErr error
	}{
		{
			name: "add role",
			role: &domain.ProjectRole{
				InstanceID:     instanceID,
				ProjectID:      projectID,
				OrganizationID: orgID,
				Key:            "key2",
				DisplayName:    "Role 2",
				RoleGroup:      gu.Ptr("group2"),
			},
		},
		{
			name: "add role, no group",
			role: &domain.ProjectRole{
				InstanceID:     instanceID,
				ProjectID:      projectID,
				OrganizationID: orgID,
				Key:            "key2",
				DisplayName:    "Role 2",
			},
		},
		{
			name: "non-existing instance",
			role: &domain.ProjectRole{
				InstanceID:     "foo",
				ProjectID:      projectID,
				OrganizationID: orgID,
				Key:            "key2",
				DisplayName:    "Role 2",
				RoleGroup:      gu.Ptr("group2"),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing project",
			role: &domain.ProjectRole{
				InstanceID:     instanceID,
				ProjectID:      "foo",
				OrganizationID: orgID,
				Key:            "key2",
				DisplayName:    "Role 2",
				RoleGroup:      gu.Ptr("group2"),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing organization",
			role: &domain.ProjectRole{
				InstanceID:     instanceID,
				ProjectID:      projectID,
				OrganizationID: "foo",
				Key:            "key2",
				DisplayName:    "Role 2",
				RoleGroup:      gu.Ptr("group2"),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "empty key error",
			role: &domain.ProjectRole{
				InstanceID:     instanceID,
				ProjectID:      projectID,
				OrganizationID: orgID,
				Key:            "",
				DisplayName:    "Role 2",
				RoleGroup:      gu.Ptr("group2"),
			},
			wantErr: new(database.CheckError),
		},
		{
			name: "empty display name error",
			role: &domain.ProjectRole{
				InstanceID:     instanceID,
				ProjectID:      projectID,
				OrganizationID: orgID,
				Key:            "key2",
				DisplayName:    "",
				RoleGroup:      gu.Ptr("group2"),
			},
			wantErr: new(database.CheckError),
		},
		{
			name:    "duplicate key",
			role:    existingRole,
			wantErr: new(database.UniqueError),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			err := roleRepo.Create(t.Context(), savepoint, tt.role)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUpdateProjectRoles(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	roleRepo := repository.ProjectRepository().Role()

	existingRole := &domain.ProjectRole{
		InstanceID:     instanceID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		Key:            "key1",
		DisplayName:    "Role 1",
		RoleGroup:      gu.Ptr("group1"),
	}
	err := roleRepo.Create(t.Context(), tx, existingRole)
	require.NoError(t, err)
	lastUpdatedAt := existingRole.UpdatedAt

	tests := []struct {
		name             string
		condition        database.Condition
		changes          []database.Change
		wantRowsAffected int64
		wantErr          error
		assertChanges    func(t *testing.T, project *domain.ProjectRole)
	}{
		{
			name:      "no changes",
			condition: roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.ProjectID, existingRole.Key),
			wantErr:   database.ErrNoChanges,
		},
		{
			name:      "incomplete condition",
			condition: roleRepo.KeyCondition(existingRole.Key),
			changes: []database.Change{
				roleRepo.SetDisplayName("Role 1 Updated"),
			},
			wantErr: database.NewMissingConditionError(roleRepo.InstanceIDColumn()),
		},
		{
			name:      "update display name",
			condition: roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.ProjectID, existingRole.Key),
			changes: []database.Change{
				roleRepo.SetDisplayName("Role 1 Updated"),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, role *domain.ProjectRole) {
				assert.Equal(t, "Role 1 Updated", role.DisplayName)
			},
		},
		{
			name:      "update role group",
			condition: roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.ProjectID, existingRole.Key),
			changes: []database.Change{
				roleRepo.SetRoleGroup("group1 Updated"),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, role *domain.ProjectRole) {
				assert.Equal(t, gu.Ptr("group1 Updated"), role.RoleGroup)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			gotRowsAffected, err := roleRepo.Update(t.Context(), savepoint, tt.condition, tt.changes...)
			assert.Equal(t, tt.wantRowsAffected, gotRowsAffected)
			assert.ErrorIs(t, err, tt.wantErr)

			if tt.assertChanges != nil {
				updatedRole, err := roleRepo.Get(t.Context(), savepoint, database.WithCondition(
					roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.ProjectID, existingRole.Key),
				))
				require.NoError(t, err)
				assert.WithinRange(t, updatedRole.CreatedAt, existingRole.CreatedAt, existingRole.CreatedAt)
				assert.WithinRange(t, updatedRole.UpdatedAt, lastUpdatedAt, lastUpdatedAt.Add(time.Second))
				lastUpdatedAt = updatedRole.UpdatedAt
				tt.assertChanges(t, updatedRole)
			}
		})
	}
}

func TestDeleteProjectRole(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	roleRepo := repository.ProjectRepository().Role()

	existingRole := &domain.ProjectRole{
		InstanceID:     instanceID,
		ProjectID:      projectID,
		OrganizationID: orgID,
		Key:            "key1",
		DisplayName:    "Role 1",
		RoleGroup:      gu.Ptr("group1"),
	}
	err := roleRepo.Create(t.Context(), tx, existingRole)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:      "incomplete condition",
			condition: roleRepo.KeyCondition(existingRole.Key),
			wantErr:   database.NewMissingConditionError(roleRepo.InstanceIDColumn()),
		},
		{
			name:             "not found",
			condition:        roleRepo.PrimaryKeyCondition(instanceID, projectID, "baz"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete role",
			condition:        roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.ProjectID, existingRole.Key),
			wantRowsAffected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			gotRowsAffected, err := roleRepo.Delete(t.Context(), savepoint, tt.condition)
			assert.Equal(t, tt.wantRowsAffected, gotRowsAffected)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
