package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestListProjectGrantRoles(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	firstGrantedOrgID := createOrganization(t, tx, instanceID)
	firstGrantID := createProjectGrant(t, tx, instanceID, orgID, projectID, firstGrantedOrgID)
	secondGrantedOrgID := createOrganization(t, tx, instanceID)
	secondGrantID := createProjectGrant(t, tx, instanceID, orgID, projectID, secondGrantedOrgID)
	roleRepo := repository.ProjectGrantRepository().Role()

	roles := [...]*domain.ProjectGrantRole{
		{
			InstanceID:   instanceID,
			GrantID:      firstGrantID,
			ProjectOrgID: orgID,
			ProjectID:    projectID,
			Key:          "key1",
		},
		{
			InstanceID:   instanceID,
			GrantID:      firstGrantID,
			ProjectOrgID: orgID,
			ProjectID:    projectID,
			Key:          "key2",
		},
		{
			InstanceID:   instanceID,
			GrantID:      secondGrantID,
			ProjectOrgID: orgID,
			ProjectID:    projectID,
			Key:          "key1",
		},
		{
			InstanceID:   instanceID,
			GrantID:      secondGrantID,
			ProjectOrgID: orgID,
			ProjectID:    projectID,
			Key:          "key2",
		},
	}
	for _, r := range roles {
		err := roleRepo.Add(t.Context(), tx, r)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		wantRoles []*domain.ProjectGrantRole
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: roleRepo.KeyCondition("key1"),
			wantErr:   database.NewMissingConditionError(roleRepo.InstanceIDColumn()),
		},
		{
			name:      "no results, ok",
			condition: roleRepo.PrimaryKeyCondition(instanceID, firstGrantID, "foo"),
		},
		{
			name: "all from project grant 1",
			condition: database.And(
				roleRepo.InstanceIDCondition(instanceID),
				roleRepo.GrantIDCondition(firstGrantID),
			),
			wantRoles: roles[0:2],
		},
		{
			name: "key 1 from project grant 2",
			condition: database.And(
				roleRepo.InstanceIDCondition(instanceID),
				roleRepo.GrantIDCondition(secondGrantID),
				roleRepo.KeyCondition("key2"),
			),
			wantRoles: []*domain.ProjectGrantRole{roles[3]},
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

func TestAddProjectRole(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	grantID := createProjectGrant(t, tx, instanceID, orgID, projectID, grantedOrgID)

	roleRepo := repository.ProjectGrantRepository().Role()

	existingRole := &domain.ProjectGrantRole{
		InstanceID:   instanceID,
		GrantID:      grantID,
		ProjectOrgID: orgID,
		ProjectID:    projectID,
		Key:          "key1",
	}
	err := roleRepo.Add(t.Context(), tx, existingRole)
	require.NoError(t, err)

	tests := []struct {
		name    string
		role    *domain.ProjectGrantRole
		wantErr error
	}{
		{
			name: "add role",
			role: &domain.ProjectGrantRole{
				InstanceID:   instanceID,
				GrantID:      grantID,
				ProjectOrgID: orgID,
				ProjectID:    projectID,
				Key:          "key2",
			},
		},
		{
			name: "non-existing instance",
			role: &domain.ProjectGrantRole{
				InstanceID:   "foo",
				GrantID:      grantID,
				ProjectOrgID: orgID,
				ProjectID:    projectID,
				Key:          "key3",
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing project",
			role: &domain.ProjectGrantRole{
				InstanceID:   instanceID,
				GrantID:      grantID,
				ProjectOrgID: orgID,
				ProjectID:    "foo",
				Key:          "key3",
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing organization",
			role: &domain.ProjectGrantRole{
				InstanceID:   instanceID,
				GrantID:      grantID,
				ProjectOrgID: "foo",
				ProjectID:    projectID,
				Key:          "key3",
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "empty key error",
			role: &domain.ProjectGrantRole{
				InstanceID:   instanceID,
				GrantID:      grantID,
				ProjectOrgID: orgID,
				ProjectID:    projectID,
				Key:          "",
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

			err := roleRepo.Add(t.Context(), savepoint, tt.role)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRemoveProjectRole(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, orgID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	grantID := createProjectGrant(t, tx, instanceID, orgID, projectID, grantedOrgID)
	roleRepo := repository.ProjectGrantRepository().Role()

	existingRole := &domain.ProjectGrantRole{
		InstanceID:   instanceID,
		GrantID:      grantID,
		ProjectOrgID: orgID,
		ProjectID:    projectID,
		Key:          "key1",
	}
	err := roleRepo.Add(t.Context(), tx, existingRole)
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
			condition:        roleRepo.PrimaryKeyCondition(instanceID, grantID, "baz"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete role",
			condition:        roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.GrantID, existingRole.Key),
			wantRowsAffected: 1,
		},
		{
			name:             "delete role twice",
			condition:        roleRepo.PrimaryKeyCondition(existingRole.InstanceID, existingRole.GrantID, existingRole.Key),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRowsAffected, err := roleRepo.Remove(t.Context(), tx, tt.condition)
			assert.Equal(t, tt.wantRowsAffected, gotRowsAffected)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
