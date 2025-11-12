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

func TestCreateAuthorization(t *testing.T) {
	beforeCreate := time.Now()

	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, organizationID)
	require.NotNil(t, projectID)
	role1 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role1")
	role2 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role2")
	// TODO: uncomment when user table is available
	//userID := createUser(t, tx, instanceID, organizationID)
	//require.NotNil(t, userID)
	userID := integration.ID()

	authorizationRepo := repository.AuthorizationRepository()

	existingAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1},
	}
	err := authorizationRepo.Create(t.Context(), tx, existingAuthorization)
	require.NoError(t, err)

	tests := []struct {
		name          string
		authorization *domain.Authorization
		wantErr       error
	}{
		{
			name: "create authorization without roles",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    integration.ID(),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
			},
		},
		{
			name: "create authorization with roles",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    integration.ID(),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      []string{role1, role2},
			},
		},
		{
			name: "non-existent instance",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    integration.ID(),
				InstanceID: "random-id",
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existent project",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  "random-id",
				GrantID:    integration.ID(),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
			},
			wantErr: new(database.ForeignKeyError),
		},
		// TODO: uncomment when user table is available
		//{
		//	name: "non-existent user",
		//	authorization: &domain.Authorization{
		//		ID:         integration.ID(),
		//		UserID:     "random-id",
		//		ProjectID:  projectID,
		//		GrantID:    integration.ID(),
		//		InstanceID: instanceID,
		//		State:      domain.AuthorizationStateActive,
		//		Roles:      nil,
		//	},
		//	wantErr: new(database.ForeignKeyError),
		//},
		{
			name: "duplicate authorization",
			authorization: &domain.Authorization{
				ID:         existingAuthorization.ID,
				UserID:     existingAuthorization.UserID,
				ProjectID:  existingAuthorization.ProjectID,
				GrantID:    existingAuthorization.GrantID,
				InstanceID: existingAuthorization.InstanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
			},
			wantErr: database.NewUniqueError("authorizations", "uq_authorizations_instance_id_id", nil),
		},
		{
			name: "missing ID",
			authorization: &domain.Authorization{
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    integration.ID(),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
			},
			wantErr: database.NewCheckError("authorizations", "", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			err := authorizationRepo.Create(t.Context(), savepoint, tt.authorization)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			afterCreate := time.Now()

			// get authorization
			got, err := authorizationRepo.Get(t.Context(), savepoint,
				database.WithCondition(
					authorizationRepo.PrimaryKeyCondition(tt.authorization.InstanceID, tt.authorization.ID),
				),
			)
			require.NoError(t, err)
			require.Equal(t, tt.authorization, got)
			assert.WithinRange(t, got.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, got.UpdatedAt, beforeCreate, afterCreate)
		})
	}
}

func TestGetAuthorization(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)

	projectID := createProject(t, tx, instanceID, organizationID)
	require.NotNil(t, projectID)
	role1 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role1")
	role2 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role2")

	// TODO: uncomment when user table is available
	//userID := createUser(t, tx, instanceID, organizationID)
	//require.NotNil(t, userID)
	userID := integration.ID()

	// create authorization with roles
	authorizationRepo := repository.AuthorizationRepository()
	authorizationWithRolesUser1 := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1, role2},
	}
	err := authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser1)
	require.NoError(t, err)

	// create authorization without roles
	authorizationWithoutRolesUser1 := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser1)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition database.Condition
		want      *domain.Authorization
		wantErr   error
	}{
		{
			name:      "incomplete condition",
			condition: authorizationRepo.IDCondition("123"),
			wantErr:   new(database.MissingConditionError),
		},
		{
			name: "non-existent authorization",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				"random-id",
			),
			wantErr: new(database.NoRowFoundError),
		},
		{
			name: "get authorization with roles",
			condition: authorizationRepo.PrimaryKeyCondition(
				authorizationWithRolesUser1.InstanceID,
				authorizationWithRolesUser1.ID,
			),
			want: authorizationWithRolesUser1,
		},
		{
			name: "get authorization without roles",
			condition: authorizationRepo.PrimaryKeyCondition(
				authorizationWithoutRolesUser1.InstanceID,
				authorizationWithoutRolesUser1.ID,
			),
			want: authorizationWithoutRolesUser1,
		},
		{
			name: "get authorization, multiple rows, error",
			condition: database.And(
				authorizationRepo.InstanceIDCondition(authorizationWithRolesUser1.InstanceID),
			),
			wantErr: new(database.MultipleRowsFoundError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authorizationRepo.Get(t.Context(), tx,
				database.WithCondition(tt.condition),
			)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestListAuthorization(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)

	project1ID := createProject(t, tx, instanceID, organizationID)
	require.NotNil(t, project1ID)
	project1Role1 := createProjectRole(t, tx, instanceID, organizationID, project1ID, "project1Role1")
	project1Role2 := createProjectRole(t, tx, instanceID, organizationID, project1ID, "project2Role2")

	project2ID := createProject(t, tx, instanceID, organizationID)
	require.NotNil(t, project2ID)
	project2Role1 := createProjectRole(t, tx, instanceID, organizationID, project2ID, "project2Role1")
	project2Role2 := createProjectRole(t, tx, instanceID, organizationID, project2ID, "project2Role2")

	// TODO: uncomment when user table is available
	//user1ID := createUser(t, tx, instanceID, organizationID)
	//require.NotNil(t, user1ID)
	user1ID := integration.ID()
	//user2ID := createUser(t, tx, instanceID, organizationID)
	//require.NotNil(t, user2ID)
	user2ID := integration.ID()

	// create authorization with roles for user1 for project1
	authorizationRepo := repository.AuthorizationRepository()
	authorizationWithRolesUser1Project1 := &domain.Authorization{
		ID:         "authorization1-user1",
		UserID:     user1ID,
		ProjectID:  project1ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project1Role1, project1Role2},
	}
	err := authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser1Project1)
	require.NoError(t, err)

	// create authorization without roles for user1 for project1
	authorizationWithoutRolesUser1Project1 := &domain.Authorization{
		ID:         "authorization2-user1",
		UserID:     user1ID,
		ProjectID:  project1ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser1Project1)
	require.NoError(t, err)

	// create authorization with roles for user1 for project2
	authorizationWithRolesUser1Project2 := &domain.Authorization{
		ID:         "authorization3-user1",
		UserID:     user1ID,
		ProjectID:  project2ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project2Role1, project2Role2},
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser1Project2)
	require.NoError(t, err)

	// create authorization without roles for user1 for project2
	authorizationWithoutRolesUser1Project2 := &domain.Authorization{
		ID:         "authorization4-user1",
		UserID:     user1ID,
		ProjectID:  project2ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateInactive,
		Roles:      nil,
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser1Project2)
	require.NoError(t, err)

	// create authorization with roles for user2 for project1
	authorizationWithRolesUser2Project1 := &domain.Authorization{
		ID:         "authorization5-user2",
		UserID:     user2ID,
		ProjectID:  project1ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project1Role1, project1Role2},
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser2Project1)
	require.NoError(t, err)

	// create authorization without roles for user2 for project1
	authorizationWithoutRolesUser2Project1 := &domain.Authorization{
		ID:         "authorization6-user2",
		UserID:     user2ID,
		ProjectID:  project1ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateInactive,
		Roles:      nil,
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser2Project1)
	require.NoError(t, err)

	// create authorization with roles for user2 for project2
	authorizationWithRolesUser2Project2 := &domain.Authorization{
		ID:         "authorization7-user2",
		UserID:     user2ID,
		ProjectID:  project2ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project2Role1, project2Role2},
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser2Project2)
	require.NoError(t, err)

	// create authorization without roles for user2 for project2
	authorizationWithoutRolesUser2Project2 := &domain.Authorization{
		ID:         "authorization8-user2",
		UserID:     user2ID,
		ProjectID:  project2ID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser2Project2)
	require.NoError(t, err)

	tests := []struct {
		name      string
		condition []database.QueryOption
		want      []*domain.Authorization
		wantErr   error
	}{
		{
			name: "incomplete condition",
			condition: []database.QueryOption{
				database.WithCondition(
					authorizationRepo.IDCondition("123"),
				),
			},
			wantErr: new(database.MissingConditionError),
		},
		{
			name: "non-existent authorization",
			condition: []database.QueryOption{
				database.WithCondition(authorizationRepo.PrimaryKeyCondition(
					instanceID,
					"1234"),
				),
			},
			want: []*domain.Authorization{},
		},
		{
			name: "list all authorizations in instance",
			condition: []database.QueryOption{
				database.WithCondition(authorizationRepo.InstanceIDCondition(
					instanceID,
				)),
				database.WithOrderByAscending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationWithRolesUser1Project1,
				authorizationWithoutRolesUser1Project1,
				authorizationWithRolesUser1Project2,
				authorizationWithoutRolesUser1Project2,
				authorizationWithRolesUser2Project1,
				authorizationWithoutRolesUser2Project1,
				authorizationWithRolesUser2Project2,
				authorizationWithoutRolesUser2Project2,
			},
		},
		{
			name: "list all authorizations for user1",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.UserIDCondition(user1ID),
					),
				),
			},
			want: []*domain.Authorization{
				authorizationWithRolesUser1Project1,
				authorizationWithoutRolesUser1Project1,
				authorizationWithRolesUser1Project2,
				authorizationWithoutRolesUser1Project2,
			},
		},
		{
			name: "list all authorizations for user2",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.UserIDCondition(user2ID),
					),
				),
			},
			want: []*domain.Authorization{
				authorizationWithRolesUser2Project1,
				authorizationWithoutRolesUser2Project1,
				authorizationWithRolesUser2Project2,
				authorizationWithoutRolesUser2Project2,
			},
		},
		{
			name: "list all authorizations for user1 and project2",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.UserIDCondition(user1ID),
						authorizationRepo.ProjectIDCondition(project2ID),
					),
				),
			},
			want: []*domain.Authorization{
				authorizationWithRolesUser1Project2,
				authorizationWithoutRolesUser1Project2,
			},
		},
		{
			name: "list all authorizations for user2 and project1, ordered by ID descending",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.UserIDCondition(user2ID),
						authorizationRepo.ProjectIDCondition(project1ID),
					),
				),
				database.WithOrderByDescending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationWithoutRolesUser2Project1,
				authorizationWithRolesUser2Project1,
			},
		},
		{
			name: "list all authorizations for project1, limited to 2",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.ProjectIDCondition(project1ID),
					),
				),
				database.WithLimit(2),
				database.WithOrderByDescending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationWithoutRolesUser2Project1,
				authorizationWithRolesUser2Project1,
			},
		},
		{
			name: "list all inactive authorizations",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.StateCondition(domain.AuthorizationStateInactive),
					),
				),
				database.WithOrderByAscending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationWithoutRolesUser1Project2,
				authorizationWithoutRolesUser2Project1,
			},
		},
		{
			name: "list all authorizations with role project1Role1",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.RolesCondition(database.TextOperationContains, project1Role1),
					),
				),
				database.WithOrderByAscending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				{
					ID:         "authorization1-user1",
					UserID:     user1ID,
					ProjectID:  project1ID,
					GrantID:    authorizationWithRolesUser1Project1.GrantID,
					InstanceID: instanceID,
					State:      domain.AuthorizationStateActive,
					Roles:      []string{project1Role1},
					CreatedAt:  authorizationWithRolesUser1Project1.CreatedAt,
					UpdatedAt:  authorizationWithRolesUser1Project1.UpdatedAt,
				},
				{
					ID:         "authorization5-user2",
					UserID:     user2ID,
					ProjectID:  project1ID,
					GrantID:    authorizationWithRolesUser2Project1.GrantID,
					InstanceID: instanceID,
					State:      domain.AuthorizationStateActive,
					Roles:      []string{project1Role1},
					CreatedAt:  authorizationWithRolesUser2Project1.CreatedAt,
					UpdatedAt:  authorizationWithRolesUser2Project1.UpdatedAt,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authorizationRepo.List(t.Context(), tx, tt.condition...)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			require.Len(t, got, len(tt.want))
			for i := range got {
				assertAuthorization(t, tt.want[i], got[i])
			}
		})
	}
}

func assertAuthorization(t *testing.T, expected, actual *domain.Authorization) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.ProjectID, actual.ProjectID)
	assert.Equal(t, expected.GrantID, actual.GrantID)
	assert.Equal(t, expected.InstanceID, actual.InstanceID)
	assert.Equal(t, expected.State, actual.State)
	assert.ElementsMatch(t, expected.Roles, actual.Roles)
	assert.Equal(t, expected.CreatedAt, actual.CreatedAt)
	assert.Equal(t, expected.UpdatedAt, actual.UpdatedAt)
}

func TestUpdateAuthorization(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, organizationID)

	authorizationRepo := repository.AuthorizationRepository()
	activeAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     integration.ID(),
		ProjectID:  projectID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
	}
	err := authorizationRepo.Create(t.Context(), tx, activeAuthorization)
	require.NoError(t, err)

	inactiveAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     integration.ID(),
		ProjectID:  projectID,
		GrantID:    integration.ID(),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateInactive,
		Roles:      nil,
	}
	err = authorizationRepo.Create(t.Context(), tx, inactiveAuthorization)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		changes          []database.Change
		wantState        domain.AuthorizationState
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name: "no changes",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				activeAuthorization.ID,
			),
			changes: nil,
			wantErr: database.ErrNoChanges,
		},
		{
			name:      "incomplete condition",
			condition: authorizationRepo.IDCondition(activeAuthorization.ID),
			changes: []database.Change{
				authorizationRepo.SetState(domain.AuthorizationStateInactive),
			},
			wantErr: new(database.MissingConditionError),
		},
		{
			name: "update state to inactive",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				activeAuthorization.ID,
			),
			changes: []database.Change{
				authorizationRepo.SetState(domain.AuthorizationStateInactive),
			},
			wantState:        domain.AuthorizationStateInactive,
			wantRowsAffected: 1,
		},
		{
			name: "update state to active",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				inactiveAuthorization.ID,
			),
			changes: []database.Change{
				authorizationRepo.SetState(domain.AuthorizationStateActive),
			},
			wantState:        domain.AuthorizationStateActive,
			wantRowsAffected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			rowsAffected, err := authorizationRepo.Update(t.Context(), savepoint, tt.condition, tt.changes...)
			require.ErrorIs(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)

			if tt.wantRowsAffected == 0 {
				return
			}
			// verify update
			updatedAuthorization, err := authorizationRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.NoError(t, err)
			assert.Equal(t, tt.wantState, updatedAuthorization.State)
		})
	}

}
