package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
)

func TestGetProjectGrant(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	grantingOrgID := createOrganization(t, tx, instanceID)
	firstGrantedOrgID := createOrganization(t, tx, instanceID)
	secondGrantedOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, grantingOrgID)
	firstRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	secondRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	thirdRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	projectGrantRepo := repository.ProjectGrantRepository()

	firstProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		GrantingOrganizationID: grantingOrgID,
		ProjectID:              projectID,
		GrantedOrganizationID:  firstGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               []string{firstRoleKey},
	}
	err := projectGrantRepo.Create(t.Context(), tx, firstProjectGrant)
	require.NoError(t, err)
	secondProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		GrantingOrganizationID: grantingOrgID,
		ProjectID:              projectID,
		GrantedOrganizationID:  secondGrantedOrgID,
		State:                  domain.ProjectGrantStateInactive,
		RoleKeys:               []string{firstRoleKey, secondRoleKey, thirdRoleKey},
	}
	err = projectGrantRepo.Create(t.Context(), tx, secondProjectGrant)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.ProjectGrant
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: projectGrantRepo.IDCondition(firstProjectGrant.ID),
			wantErr:   database.NewMissingConditionError(projectGrantRepo.IDColumn()),
		},
		{
			name:      "not found",
			condition: projectGrantRepo.PrimaryKeyCondition(instanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: projectGrantRepo.InstanceIDCondition(instanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok first",
			condition: projectGrantRepo.PrimaryKeyCondition(instanceID, firstProjectGrant.ID),
			want:      firstProjectGrant,
		},
		{
			name:      "ok second",
			condition: projectGrantRepo.PrimaryKeyCondition(instanceID, secondProjectGrant.ID),
			want:      secondProjectGrant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectGrantRepo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assertProjectGrant(t, tt.want, got)
		})
	}
}

func TestListProjectGrants(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	firstGrantingOrgID := createOrganization(t, tx, instanceID)
	firstProjectID := createProject(t, tx, instanceID, firstGrantingOrgID)
	firstRoleKey := createProjectRole(t, tx, instanceID, firstGrantingOrgID, firstProjectID, "")
	secondRoleKey := createProjectRole(t, tx, instanceID, firstGrantingOrgID, firstProjectID, "")
	thirdRoleKey := createProjectRole(t, tx, instanceID, firstGrantingOrgID, firstProjectID, "")
	secondGrantingOrgID := createOrganization(t, tx, instanceID)
	secondProjectID := createProject(t, tx, instanceID, secondGrantingOrgID)

	firstGrantedOrgID := createOrganization(t, tx, instanceID)
	secondGrantedOrgID := createOrganization(t, tx, instanceID)
	projectGrantRepo := repository.ProjectGrantRepository()

	projectGrants := [...]*domain.ProjectGrant{
		{
			InstanceID:             instanceID,
			ID:                     "1",
			ProjectID:              firstProjectID,
			GrantedOrganizationID:  firstGrantedOrgID,
			GrantingOrganizationID: firstGrantingOrgID,
			State:                  domain.ProjectGrantStateActive,
			RoleKeys:               []string{firstRoleKey},
		},
		{
			InstanceID:             instanceID,
			ID:                     "2",
			ProjectID:              firstProjectID,
			GrantedOrganizationID:  secondGrantedOrgID,
			GrantingOrganizationID: firstGrantingOrgID,
			State:                  domain.ProjectGrantStateInactive,
			RoleKeys:               []string{secondRoleKey, thirdRoleKey},
		},
		{
			InstanceID:             instanceID,
			ID:                     "3",
			ProjectID:              secondProjectID,
			GrantedOrganizationID:  firstGrantedOrgID,
			GrantingOrganizationID: secondGrantingOrgID,
			State:                  domain.ProjectGrantStateInactive,
		},
		{
			InstanceID:             instanceID,
			ID:                     "4",
			ProjectID:              secondProjectID,
			GrantedOrganizationID:  secondGrantedOrgID,
			GrantingOrganizationID: secondGrantingOrgID,
			State:                  domain.ProjectGrantStateActive,
		},
	}
	for _, projectGrant := range projectGrants {
		err := projectGrantRepo.Create(t.Context(), tx, projectGrant)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.ProjectGrant
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: projectGrantRepo.GrantedOrganizationIDCondition(firstGrantedOrgID),
			wantErr:   database.NewMissingConditionError(projectGrantRepo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: projectGrantRepo.PrimaryKeyCondition(instanceID, "nix"),
		},
		{
			name:      "all from instance",
			condition: projectGrantRepo.InstanceIDCondition(instanceID),
			want:      projectGrants[:],
		},
		{
			name: "all granting from first org",
			condition: database.And(
				projectGrantRepo.InstanceIDCondition(instanceID),
				projectGrantRepo.GrantingOrganizationIDCondition(firstGrantingOrgID),
			),
			want: projectGrants[0:2],
		},
		{
			name: "all granting from second org",
			condition: database.And(
				projectGrantRepo.InstanceIDCondition(instanceID),
				projectGrantRepo.GrantingOrganizationIDCondition(secondGrantingOrgID),
			),
			want: projectGrants[2:4],
		},
		{
			name: "all granted from first org",
			condition: database.And(
				projectGrantRepo.InstanceIDCondition(instanceID),
				projectGrantRepo.GrantingOrganizationIDCondition(firstGrantingOrgID),
			),
			want: []*domain.ProjectGrant{projectGrants[0], projectGrants[1]},
		},
		{
			name: "all granted to second org",
			condition: database.And(
				projectGrantRepo.InstanceIDCondition(instanceID),
				projectGrantRepo.GrantedOrganizationIDCondition(secondGrantedOrgID),
			),
			want: []*domain.ProjectGrant{projectGrants[1], projectGrants[3]},
		},
		{
			name: "state active",
			condition: database.And(
				projectGrantRepo.InstanceIDCondition(instanceID),
				projectGrantRepo.StateCondition(domain.ProjectGrantStateActive),
			),
			want: []*domain.ProjectGrant{projectGrants[0], projectGrants[3]},
		},
		{
			name: "exists role key",
			condition: database.And(
				projectGrantRepo.InstanceIDCondition(instanceID),
				projectGrantRepo.ExistsRoleKey(projectGrantRepo.RoleKeyCondition(database.TextOperationEqual, firstRoleKey)),
			),
			want: []*domain.ProjectGrant{projectGrants[0]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectGrantRepo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(projectGrantRepo.PrimaryKeyColumns()...),
			)
			require.ErrorIs(t, err, tt.wantErr)

			if !assert.Len(t, got, len(tt.want)) {
				return
			}
			for i, want := range tt.want {
				assertProjectGrant(t, want, got[i])
			}
		})
	}
}

func TestCreateProjectGrant(t *testing.T) {
	beforeCreate := time.Now()

	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	grantingOrgID := createOrganization(t, tx, instanceID)
	firstGrantedOrgID := createOrganization(t, tx, instanceID)
	secondGrantedOrgID := createOrganization(t, tx, instanceID)
	thirdGrantedOrgID := createOrganization(t, tx, instanceID)
	fourthGrantedOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, grantingOrgID)
	firstRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	secondRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	thirdRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	projectGrantRepo := repository.ProjectGrantRepository()

	existingProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  firstGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
	}
	err := projectGrantRepo.Create(t.Context(), tx, existingProjectGrant)
	require.NoError(t, err)

	tests := []struct {
		name         string
		projectGrant *domain.ProjectGrant
		wantErr      error
	}{
		{
			name: "add project grant",
			projectGrant: &domain.ProjectGrant{
				InstanceID:             instanceID,
				ID:                     integration.ID(),
				ProjectID:              projectID,
				GrantingOrganizationID: grantingOrgID,
				GrantedOrganizationID:  secondGrantedOrgID,
				State:                  domain.ProjectGrantStateActive,
			},
		},
		{
			name: "non-existing instance",
			projectGrant: &domain.ProjectGrant{
				InstanceID:             "foo",
				ID:                     integration.ID(),
				ProjectID:              projectID,
				GrantingOrganizationID: grantingOrgID,
				GrantedOrganizationID:  secondGrantedOrgID,
				State:                  domain.ProjectGrantStateActive,
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			projectGrant: &domain.ProjectGrant{
				InstanceID:             instanceID,
				ID:                     integration.ID(),
				ProjectID:              projectID,
				GrantingOrganizationID: "foo",
				GrantedOrganizationID:  secondGrantedOrgID,
				State:                  domain.ProjectGrantStateActive,
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "empty id error",
			projectGrant: &domain.ProjectGrant{
				InstanceID:             instanceID,
				ID:                     "",
				ProjectID:              projectID,
				GrantingOrganizationID: grantingOrgID,
				GrantedOrganizationID:  secondGrantedOrgID,
				State:                  domain.ProjectGrantStateActive,
			},
			wantErr: new(database.CheckError),
		},
		{
			name:         "duplicate project grant",
			projectGrant: existingProjectGrant,
			wantErr:      new(database.UniqueError),
		},
		{
			name: "add project grant with role",
			projectGrant: &domain.ProjectGrant{
				InstanceID:             instanceID,
				ID:                     integration.ID(),
				ProjectID:              projectID,
				GrantingOrganizationID: grantingOrgID,
				GrantedOrganizationID:  thirdGrantedOrgID,
				State:                  domain.ProjectGrantStateActive,
				RoleKeys:               []string{firstRoleKey},
			},
		},
		{
			name: "add project grant with multiple roles",
			projectGrant: &domain.ProjectGrant{
				InstanceID:             instanceID,
				ID:                     integration.ID(),
				ProjectID:              projectID,
				GrantingOrganizationID: grantingOrgID,
				GrantedOrganizationID:  fourthGrantedOrgID,
				State:                  domain.ProjectGrantStateActive,
				RoleKeys:               []string{firstRoleKey, secondRoleKey, thirdRoleKey},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := projectGrantRepo.Create(t.Context(), savepoint, tt.projectGrant)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// check instance values
			projectGrant, err := projectGrantRepo.Get(t.Context(), tx,
				database.WithCondition(
					projectGrantRepo.PrimaryKeyCondition(tt.projectGrant.InstanceID, tt.projectGrant.ID),
				),
			)
			require.NoError(t, err)
			assertProjectGrant(t, tt.projectGrant, projectGrant)
			assert.WithinRange(t, projectGrant.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, projectGrant.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func assertProjectGrant(t assert.TestingT, expected *domain.ProjectGrant, actual *domain.ProjectGrant) {
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.ProjectID, actual.ProjectID)
	assert.Equal(t, expected.GrantingOrganizationID, actual.GrantingOrganizationID)
	assert.Equal(t, expected.GrantedOrganizationID, actual.GrantedOrganizationID)
	assert.Equal(t, expected.State, actual.State)
	assert.ElementsMatch(t, expected.RoleKeys, actual.RoleKeys)
}

func TestUpdateProjectGrant(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	grantingOrgID := createOrganization(t, tx, instanceID)
	firstGrantedOrgID := createOrganization(t, tx, instanceID)
	secondGrantedOrgID := createOrganization(t, tx, instanceID)
	thirdGrantedOrgID := createOrganization(t, tx, instanceID)
	fourthGrantedOrgID := createOrganization(t, tx, instanceID)
	fifthGrantedOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, grantingOrgID)
	firstRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	secondRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	thirdRoleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	projectGrantRepo := repository.ProjectGrantRepository()

	existingProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  firstGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
	}
	err := projectGrantRepo.Create(t.Context(), tx, existingProjectGrant)
	require.NoError(t, err)

	existingProjectGrantWithRole := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  secondGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               []string{firstRoleKey},
	}
	err = projectGrantRepo.Create(t.Context(), tx, existingProjectGrantWithRole)
	require.NoError(t, err)

	existingProjectGrantWithRoles := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  thirdGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               []string{firstRoleKey, secondRoleKey, thirdRoleKey},
	}
	err = projectGrantRepo.Create(t.Context(), tx, existingProjectGrantWithRoles)
	require.NoError(t, err)

	existingProjectGrantWithAllRoles := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  fourthGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               []string{firstRoleKey, secondRoleKey, thirdRoleKey},
	}
	err = projectGrantRepo.Create(t.Context(), tx, existingProjectGrantWithAllRoles)
	require.NoError(t, err)

	existingProjectGrantToRemoveRoles := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  fifthGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               []string{firstRoleKey, secondRoleKey, thirdRoleKey},
	}
	err = projectGrantRepo.Create(t.Context(), tx, existingProjectGrantToRemoveRoles)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		roleKeys         []string
		changes          []database.Change
		wantRowsAffected int64
		wantErr          error
		assertChanges    func(t *testing.T, project *domain.ProjectGrant)
	}{
		{
			name:             "no changes",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrant.ID),
			changes:          []database.Change{},
			wantRowsAffected: 0,
			wantErr:          database.ErrNoChanges,
		},
		{
			name:      "incomplete condition",
			condition: projectGrantRepo.InstanceIDCondition(instanceID),
			changes: []database.Change{
				projectGrantRepo.SetState(domain.ProjectGrantStateActive),
			},
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(projectGrantRepo.IDColumn()),
		},
		{
			name:      "set state",
			condition: projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrant.ID),
			changes: []database.Change{
				projectGrantRepo.SetState(domain.ProjectGrantStateInactive),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.ProjectGrant) {
				assert.Equal(t, domain.ProjectGrantStateInactive, updatedProject.State)
			},
		},
		{
			name:      "incomplete condition for set rolekey",
			condition: projectGrantRepo.InstanceIDCondition(instanceID),
			roleKeys:  []string{firstRoleKey},
			changes: []database.Change{
				projectGrantRepo.SetState(domain.ProjectGrantStateActive),
			},
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(projectGrantRepo.IDColumn()),
		},
		{
			name:             "set rolekey",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrant.ID),
			roleKeys:         []string{firstRoleKey},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProjectGrant *domain.ProjectGrant) {
				assert.WithinRange(t, updatedProjectGrant.CreatedAt, existingProjectGrant.CreatedAt, existingProjectGrant.CreatedAt)
				assert.WithinRange(t, updatedProjectGrant.UpdatedAt, existingProjectGrant.UpdatedAt, existingProjectGrant.UpdatedAt.Add(time.Second))
				assert.ElementsMatch(t, []string{firstRoleKey}, updatedProjectGrant.RoleKeys)
			},
		},
		{
			name:             "add rolekeys",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrantWithRole.ID),
			roleKeys:         []string{firstRoleKey, secondRoleKey, thirdRoleKey},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProjectGrant *domain.ProjectGrant) {
				assert.WithinRange(t, updatedProjectGrant.CreatedAt, existingProjectGrantWithRole.CreatedAt, existingProjectGrantWithRole.CreatedAt)
				assert.WithinRange(t, updatedProjectGrant.UpdatedAt, existingProjectGrantWithRole.UpdatedAt, existingProjectGrantWithRole.UpdatedAt.Add(time.Second))
				assert.ElementsMatch(t, []string{firstRoleKey, secondRoleKey, thirdRoleKey}, updatedProjectGrant.RoleKeys)
			},
		},
		{
			name:             "remove rolekeys",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrantWithRoles.ID),
			roleKeys:         []string{firstRoleKey},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProjectGrant *domain.ProjectGrant) {
				assert.WithinRange(t, updatedProjectGrant.CreatedAt, existingProjectGrantWithRoles.CreatedAt, existingProjectGrantWithRoles.CreatedAt)
				assert.WithinRange(t, updatedProjectGrant.UpdatedAt, existingProjectGrantWithRoles.UpdatedAt, existingProjectGrantWithRoles.UpdatedAt.Add(time.Second))
				assert.ElementsMatch(t, []string{firstRoleKey}, updatedProjectGrant.RoleKeys)
			},
		},
		{
			name:             "remove all rolekeys",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrantToRemoveRoles.ID),
			roleKeys:         []string{},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProjectGrant *domain.ProjectGrant) {
				assert.WithinRange(t, updatedProjectGrant.CreatedAt, existingProjectGrantToRemoveRoles.CreatedAt, existingProjectGrantToRemoveRoles.CreatedAt)
				assert.WithinRange(t, updatedProjectGrant.UpdatedAt, existingProjectGrantToRemoveRoles.UpdatedAt, existingProjectGrantToRemoveRoles.UpdatedAt.Add(time.Second))
				assert.ElementsMatch(t, []string{}, updatedProjectGrant.RoleKeys)
			},
		},
		{
			name:             "no changes role keys",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrantWithAllRoles.ID),
			roleKeys:         []string{firstRoleKey, secondRoleKey, thirdRoleKey},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProjectGrant *domain.ProjectGrant) {
				assert.WithinRange(t, updatedProjectGrant.CreatedAt, existingProjectGrantWithAllRoles.CreatedAt, existingProjectGrantWithAllRoles.CreatedAt)
				assert.WithinRange(t, updatedProjectGrant.UpdatedAt, existingProjectGrantWithAllRoles.UpdatedAt, existingProjectGrantWithAllRoles.UpdatedAt.Add(time.Second))
				assert.ElementsMatch(t, []string{firstRoleKey, secondRoleKey, thirdRoleKey}, updatedProjectGrant.RoleKeys)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := projectGrantRepo.Update(t.Context(), savepoint, tt.condition, tt.roleKeys, tt.changes...)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)

			if tt.assertChanges != nil {
				updatedProjectGrant, err := projectGrantRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
				require.NoError(t, err)
				tt.assertChanges(t, updatedProjectGrant)
			}
		})
	}
}

func TestDeleteProjectGrant(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	grantingOrgID := createOrganization(t, tx, instanceID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, grantingOrgID)
	roleKey := createProjectRole(t, tx, instanceID, grantingOrgID, projectID, "")
	projectGrantRepo := repository.ProjectGrantRepository()

	existingProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  grantedOrgID,
		State:                  domain.ProjectGrantStateActive,
		RoleKeys:               []string{roleKey},
	}
	err := projectGrantRepo.Create(t.Context(), tx, existingProjectGrant)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        projectGrantRepo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(projectGrantRepo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete project grant",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrant.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "delete project grant twice",
			condition:        projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrant.ID),
			wantRowsAffected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rowsAffected, err := projectGrantRepo.Delete(t.Context(), tx, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}
