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
	projectGrantRepo := repository.ProjectGrantRepository()

	firstProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		GrantingOrganizationID: grantingOrgID,
		ProjectID:              projectID,
		GrantedOrganizationID:  firstGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
	}
	err := projectGrantRepo.Create(t.Context(), tx, firstProjectGrant)
	require.NoError(t, err)
	secondProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		GrantingOrganizationID: grantingOrgID,
		ProjectID:              projectID,
		GrantedOrganizationID:  secondGrantedOrgID,
		State:                  domain.ProjectGrantStateActive,
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
			name:      "ok",
			condition: projectGrantRepo.PrimaryKeyCondition(instanceID, firstProjectGrant.ID),
			want:      firstProjectGrant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectGrantRepo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListProjectGrants(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	firstGrantingOrgID := createOrganization(t, tx, instanceID)
	firstProjectID := createProject(t, tx, instanceID, firstGrantingOrgID)
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
		},
		{
			InstanceID:             instanceID,
			ID:                     "2",
			ProjectID:              firstProjectID,
			GrantedOrganizationID:  secondGrantedOrgID,
			GrantingOrganizationID: firstGrantingOrgID,
			State:                  domain.ProjectGrantStateInactive,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectGrantRepo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(projectGrantRepo.PrimaryKeyColumns()...),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateProjectGrant(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	grantingOrgID := createOrganization(t, tx, instanceID)
	firstGrantedOrgID := createOrganization(t, tx, instanceID)
	secondGrantedOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, grantingOrgID)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := projectGrantRepo.Create(t.Context(), savepoint, tt.projectGrant)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUpdateProjectGrant(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	grantingOrgID := createOrganization(t, tx, instanceID)
	grantedOrgID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, grantingOrgID)
	projectGrantRepo := repository.ProjectGrantRepository()

	existingProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  grantedOrgID,
		State:                  domain.ProjectGrantStateActive,
	}
	err := projectGrantRepo.Create(t.Context(), tx, existingProjectGrant)
	require.NoError(t, err)
	lastUpdatedAt := existingProjectGrant.UpdatedAt

	tests := []struct {
		name             string
		condition        database.Condition
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := projectGrantRepo.Update(t.Context(), savepoint, tt.condition, tt.changes...)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)

			if tt.assertChanges != nil {
				updatedProjectGrant, err := projectGrantRepo.Get(t.Context(), savepoint, database.WithCondition(
					projectGrantRepo.PrimaryKeyCondition(instanceID, existingProjectGrant.ID),
				))
				require.NoError(t, err)
				assert.WithinRange(t, updatedProjectGrant.CreatedAt, existingProjectGrant.CreatedAt, existingProjectGrant.CreatedAt)
				assert.WithinRange(t, updatedProjectGrant.UpdatedAt, lastUpdatedAt, lastUpdatedAt.Add(time.Second))
				lastUpdatedAt = updatedProjectGrant.UpdatedAt
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
	projectGrantRepo := repository.ProjectGrantRepository()

	existingProjectGrant := &domain.ProjectGrant{
		InstanceID:             instanceID,
		ID:                     integration.ID(),
		ProjectID:              projectID,
		GrantingOrganizationID: grantingOrgID,
		GrantedOrganizationID:  grantedOrgID,
		State:                  domain.ProjectGrantStateActive,
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
