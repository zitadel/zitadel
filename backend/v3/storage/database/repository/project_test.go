package repository_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
)

func TestGetProject(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectRepo := repository.ProjectRepository()

	firstProject := &domain.Project{
		InstanceID:               instanceID,
		OrganizationID:           orgID,
		ID:                       gofakeit.UUID(),
		Name:                     gofakeit.Name(),
		State:                    domain.ProjectStateActive,
		ShouldAssertRole:         true,
		IsAuthorizationRequired:  true,
		IsProjectAccessRequired:  true,
		UsedLabelingSettingOwner: 1,
	}
	err := projectRepo.Create(t.Context(), tx, firstProject)
	require.NoError(t, err)
	secondProject := &domain.Project{
		InstanceID:               instanceID,
		OrganizationID:           orgID,
		ID:                       gofakeit.UUID(),
		Name:                     gofakeit.Name(),
		State:                    domain.ProjectStateActive,
		ShouldAssertRole:         true,
		IsAuthorizationRequired:  true,
		IsProjectAccessRequired:  true,
		UsedLabelingSettingOwner: 1,
	}
	err = projectRepo.Create(t.Context(), tx, secondProject)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.Project
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: projectRepo.IDCondition(firstProject.ID),
			wantErr:   database.NewMissingConditionError(projectRepo.IDColumn()),
		},
		{
			name:      "not found",
			condition: projectRepo.PrimaryKeyCondition(instanceID, "nix"),
			wantErr:   database.NewNoRowFoundError(nil),
		},
		{
			name:      "too many",
			condition: projectRepo.InstanceIDCondition(instanceID),
			wantErr:   database.NewMultipleRowsFoundError(nil),
		},
		{
			name:      "ok",
			condition: projectRepo.PrimaryKeyCondition(instanceID, firstProject.ID),
			want:      firstProject,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectRepo.Get(t.Context(), tx, database.WithCondition(tt.condition))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListProjects(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	firstOrgID := createOrganization(t, tx, instanceID)
	secondOrgID := createOrganization(t, tx, instanceID)
	projectRepo := repository.ProjectRepository()

	projects := [...]*domain.Project{
		{
			InstanceID:               instanceID,
			OrganizationID:           firstOrgID,
			ID:                       "1",
			Name:                     "spanac",
			State:                    domain.ProjectStateActive,
			ShouldAssertRole:         true,
			IsAuthorizationRequired:  true,
			IsProjectAccessRequired:  true,
			UsedLabelingSettingOwner: 1,
		},
		{
			InstanceID:               instanceID,
			OrganizationID:           firstOrgID,
			ID:                       "2",
			Name:                     "foobar",
			State:                    domain.ProjectStateInactive,
			ShouldAssertRole:         true,
			IsAuthorizationRequired:  true,
			IsProjectAccessRequired:  true,
			UsedLabelingSettingOwner: 1,
		},
		{
			InstanceID:               instanceID,
			OrganizationID:           secondOrgID,
			ID:                       "3",
			Name:                     "foobaz",
			State:                    domain.ProjectStateActive,
			ShouldAssertRole:         true,
			IsAuthorizationRequired:  true,
			IsProjectAccessRequired:  true,
			UsedLabelingSettingOwner: 1,
		},
		{
			InstanceID:               instanceID,
			OrganizationID:           secondOrgID,
			ID:                       "4",
			Name:                     "bazqux",
			State:                    domain.ProjectStateInactive,
			ShouldAssertRole:         true,
			IsAuthorizationRequired:  true,
			IsProjectAccessRequired:  true,
			UsedLabelingSettingOwner: 1,
		},
	}
	for _, project := range projects {
		err := projectRepo.Create(t.Context(), tx, project)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		condition database.Condition
		want      []*domain.Project
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: projectRepo.OrganizationIDCondition(firstOrgID),
			wantErr:   database.NewMissingConditionError(projectRepo.IDColumn()),
		},
		{
			name:      "no results, ok",
			condition: projectRepo.PrimaryKeyCondition(instanceID, "nix"),
		},
		{
			name:      "all from instance",
			condition: projectRepo.InstanceIDCondition(instanceID),
			want:      projects[:],
		},
		{
			name: "all from first org",
			condition: database.And(
				projectRepo.InstanceIDCondition(instanceID),
				projectRepo.OrganizationIDCondition(firstOrgID),
			),
			want: projects[0:2],
		},
		{
			name: "name starts with 'foo'",
			condition: database.And(
				projectRepo.InstanceIDCondition(instanceID),
				projectRepo.NameCondition(database.TextOperationStartsWith, "foo"),
			),
			want: projects[1:3],
		},
		{
			name: "state active",
			condition: database.And(
				projectRepo.InstanceIDCondition(instanceID),
				projectRepo.StateCondition(domain.ProjectStateActive),
			),
			want: []*domain.Project{projects[0], projects[2]},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := projectRepo.List(t.Context(), tx,
				database.WithCondition(tt.condition),
				database.WithOrderByAscending(projectRepo.PrimaryKeyColumns()...),
			)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateProject(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectRepo := repository.ProjectRepository()

	existingProject := &domain.Project{
		InstanceID:               instanceID,
		OrganizationID:           orgID,
		ID:                       gofakeit.UUID(),
		Name:                     gofakeit.Name(),
		State:                    domain.ProjectStateActive,
		ShouldAssertRole:         true,
		IsAuthorizationRequired:  true,
		IsProjectAccessRequired:  true,
		UsedLabelingSettingOwner: 1,
	}
	err := projectRepo.Create(t.Context(), tx, existingProject)
	require.NoError(t, err)

	tests := []struct {
		name    string
		project *domain.Project
		wantErr error
	}{
		{
			name: "add project",
			project: &domain.Project{
				InstanceID:               instanceID,
				OrganizationID:           orgID,
				ID:                       gofakeit.UUID(),
				Name:                     gofakeit.Name(),
				State:                    domain.ProjectStateActive,
				ShouldAssertRole:         true,
				IsAuthorizationRequired:  true,
				IsProjectAccessRequired:  true,
				UsedLabelingSettingOwner: 1,
			},
		},
		{
			name: "non-existing instance",
			project: &domain.Project{
				InstanceID:               "foo",
				OrganizationID:           orgID,
				ID:                       gofakeit.UUID(),
				Name:                     gofakeit.Name(),
				State:                    domain.ProjectStateActive,
				ShouldAssertRole:         true,
				IsAuthorizationRequired:  true,
				IsProjectAccessRequired:  true,
				UsedLabelingSettingOwner: 1,
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existing org",
			project: &domain.Project{
				InstanceID:               instanceID,
				OrganizationID:           "foo",
				ID:                       gofakeit.UUID(),
				Name:                     gofakeit.Name(),
				State:                    domain.ProjectStateActive,
				ShouldAssertRole:         true,
				IsAuthorizationRequired:  true,
				IsProjectAccessRequired:  true,
				UsedLabelingSettingOwner: 1,
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "empty id error",
			project: &domain.Project{
				InstanceID:               instanceID,
				OrganizationID:           orgID,
				ID:                       "",
				Name:                     gofakeit.Name(),
				State:                    domain.ProjectStateActive,
				ShouldAssertRole:         true,
				IsAuthorizationRequired:  true,
				IsProjectAccessRequired:  true,
				UsedLabelingSettingOwner: 1,
			},
			wantErr: new(database.CheckError),
		},
		{
			name: "empty name error",
			project: &domain.Project{
				InstanceID:               instanceID,
				OrganizationID:           orgID,
				ID:                       gofakeit.UUID(),
				Name:                     "",
				State:                    domain.ProjectStateActive,
				ShouldAssertRole:         true,
				IsAuthorizationRequired:  true,
				IsProjectAccessRequired:  true,
				UsedLabelingSettingOwner: 1,
			},
			wantErr: new(database.CheckError),
		},
		{
			name:    "duplicate project",
			project: existingProject,
			wantErr: new(database.UniqueError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			err := projectRepo.Create(t.Context(), savepoint, tt.project)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectRepo := repository.ProjectRepository()

	existingProject := &domain.Project{
		InstanceID:               instanceID,
		OrganizationID:           orgID,
		ID:                       gofakeit.UUID(),
		Name:                     gofakeit.Name(),
		State:                    domain.ProjectStateActive,
		ShouldAssertRole:         true,
		IsAuthorizationRequired:  true,
		IsProjectAccessRequired:  true,
		UsedLabelingSettingOwner: 1,
	}
	err := projectRepo.Create(t.Context(), tx, existingProject)
	require.NoError(t, err)
	lastUpdatedAt := existingProject.UpdatedAt

	tests := []struct {
		name             string
		condition        database.Condition
		changes          []database.Change
		wantRowsAffected int64
		wantErr          error
		assertChanges    func(t *testing.T, project *domain.Project)
	}{
		{
			name:             "no changes",
			condition:        projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes:          []database.Change{},
			wantRowsAffected: 0,
			wantErr:          database.ErrNoChanges,
		},
		{
			name:      "incomplete condition",
			condition: projectRepo.InstanceIDCondition(instanceID),
			changes: []database.Change{
				projectRepo.SetName("new name"),
			},
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(projectRepo.IDColumn()),
		},
		{
			name:      "set name",
			condition: projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes: []database.Change{
				projectRepo.SetName("new name"),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.Project) {
				assert.Equal(t, "new name", updatedProject.Name)
			},
		},
		{
			name:      "set state",
			condition: projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes: []database.Change{
				projectRepo.SetState(domain.ProjectStateInactive),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.Project) {
				assert.Equal(t, domain.ProjectStateInactive, updatedProject.State)
			},
		},
		{
			name:      "set should_assert_role",
			condition: projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes: []database.Change{
				projectRepo.SetShouldAssertRole(false),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.Project) {
				assert.False(t, updatedProject.ShouldAssertRole)
			},
		},
		{
			name:      "set is_authorization_required",
			condition: projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes: []database.Change{
				projectRepo.SetIsAuthorizationRequired(false),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.Project) {
				assert.False(t, updatedProject.IsAuthorizationRequired)
			},
		},
		{
			name:      "set is_project_access_required",
			condition: projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes: []database.Change{
				projectRepo.SetIsProjectAccessRequired(false),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.Project) {
				assert.False(t, updatedProject.IsProjectAccessRequired)
			},
		},
		{
			name:      "set used_labeling_setting_owner",
			condition: projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			changes: []database.Change{
				projectRepo.SetUsedLabelingSettingOwner(2),
			},
			wantRowsAffected: 1,
			assertChanges: func(t *testing.T, updatedProject *domain.Project) {
				assert.Equal(t, int16(2), updatedProject.UsedLabelingSettingOwner)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := projectRepo.Update(t.Context(), savepoint, tt.condition, tt.changes...)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)

			if tt.assertChanges != nil {
				updatedProject, err := projectRepo.Get(t.Context(), savepoint, database.WithCondition(
					projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
				))
				require.NoError(t, err)
				assert.WithinRange(t, updatedProject.CreatedAt, existingProject.CreatedAt, existingProject.CreatedAt)
				assert.WithinRange(t, updatedProject.UpdatedAt, lastUpdatedAt, lastUpdatedAt.Add(time.Second))
				lastUpdatedAt = updatedProject.UpdatedAt
				tt.assertChanges(t, updatedProject)
			}
		})
	}
}

func TestDeleteProject(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()
	instanceID := createInstance(t, tx)
	orgID := createOrganization(t, tx, instanceID)
	projectRepo := repository.ProjectRepository()

	existingProject := &domain.Project{
		InstanceID:               instanceID,
		OrganizationID:           orgID,
		ID:                       gofakeit.UUID(),
		Name:                     gofakeit.Name(),
		State:                    domain.ProjectStateActive,
		ShouldAssertRole:         true,
		IsAuthorizationRequired:  true,
		IsProjectAccessRequired:  true,
		UsedLabelingSettingOwner: 1,
	}
	err := projectRepo.Create(t.Context(), tx, existingProject)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:             "incomplete condition",
			condition:        projectRepo.InstanceIDCondition(instanceID),
			wantRowsAffected: 0,
			wantErr:          database.NewMissingConditionError(projectRepo.IDColumn()),
		},
		{
			name:             "not found",
			condition:        projectRepo.PrimaryKeyCondition(instanceID, "foo"),
			wantRowsAffected: 0,
		},
		{
			name:             "delete project",
			condition:        projectRepo.PrimaryKeyCondition(instanceID, existingProject.ID),
			wantRowsAffected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()
			rowsAffected, err := projectRepo.Delete(t.Context(), savepoint, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)
		})
	}
}
