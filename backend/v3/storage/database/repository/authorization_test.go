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

func TestCreateAuthorization(t *testing.T) {
	beforeCreate := time.Now()

	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)
	// create a project
	projectID := createProject(t, tx, instanceID, organizationID)
	require.NotNil(t, projectID)

	// create project roles
	role1 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role1")
	role2 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role2")
	userID := createHumanUser(t, tx, instanceID, organizationID)
	require.NotNil(t, userID)

	authorizationRepo := repository.AuthorizationRepository()

	// create project authorization
	existingAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := authorizationRepo.Create(t.Context(), tx, existingAuthorization)
	require.NoError(t, err)

	// create a project grant
	grantedOrganizationID := createOrganization(t, tx, instanceID)
	grantID := createProjectGrant(t, tx, instanceID, organizationID, grantedOrganizationID, projectID, []string{role1})

	tests := []struct {
		name          string
		authorization *domain.Authorization
		wantErr       error
	}{
		{
			name: "non-existent instance",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				InstanceID: "random-id",
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existent project",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  "random-id",
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existent project role",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  "random-id",
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      []string{"role3"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existent project grant",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    gu.Ptr("random-id"),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "non-existent user",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     "random-id",
				ProjectID:  projectID,
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
			},
			wantErr: new(database.ForeignKeyError),
		},
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
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: database.NewUniqueError("authorizations", "uq_authorizations_instance_id_id", nil),
		},
		{
			name: "missing ID",
			authorization: &domain.Authorization{
				UserID:     userID,
				ProjectID:  projectID,
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: database.NewCheckError("authorizations", "", nil),
		},
		{
			name: "create authorization without roles",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			name: "create authorization with roles",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      []string{role1, role2},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			name: "create authorization for project grant with valid role",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    gu.Ptr(grantID),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      []string{role1},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
		{
			name: "create authorization for project grant with unassigned role",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    gu.Ptr(grantID),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      []string{role2},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			wantErr: new(database.ForeignKeyError),
		},
		{
			name: "create authorization for project grant without roles",
			authorization: &domain.Authorization{
				ID:         integration.ID(),
				UserID:     userID,
				ProjectID:  projectID,
				GrantID:    gu.Ptr(grantID),
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      nil,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
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

	userID := createHumanUser(t, tx, instanceID, organizationID)
	require.NotNil(t, userID)

	// create authorization with roles
	authorizationRepo := repository.AuthorizationRepository()
	authorizationWithRolesUser := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1, role2},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser)
	require.NoError(t, err)

	// create authorization without roles
	authorizationWithoutRolesUser := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser)
	require.NoError(t, err)

	// create authorization for a project grant with/without roles
	grantedOrganizationID := createOrganization(t, tx, instanceID)
	grantID := createProjectGrant(t, tx, instanceID, organizationID, grantedOrganizationID, projectID, []string{role1})

	authorizationProjectGrantWithRoles := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		GrantID:    gu.Ptr(grantID),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationProjectGrantWithRoles)
	require.NoError(t, err)

	authorizationProjectGrantWithoutRoles := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		GrantID:    gu.Ptr(grantID),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationProjectGrantWithoutRoles)
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
				authorizationWithRolesUser.InstanceID,
				authorizationWithRolesUser.ID,
			),
			want: authorizationWithRolesUser,
		},
		{
			name: "get authorization without roles",
			condition: authorizationRepo.PrimaryKeyCondition(
				authorizationWithoutRolesUser.InstanceID,
				authorizationWithoutRolesUser.ID,
			),
			want: authorizationWithoutRolesUser,
		},
		{
			name: "get authorization, multiple rows, error",
			condition: database.And(
				authorizationRepo.InstanceIDCondition(authorizationWithRolesUser.InstanceID),
			),
			wantErr: new(database.MultipleRowsFoundError),
		},
		{
			name: "get authorization with role condition",
			condition: database.And(authorizationRepo.PrimaryKeyCondition(
				authorizationWithRolesUser.InstanceID,
				authorizationWithRolesUser.ID,
			),
				authorizationRepo.RoleCondition(database.TextOperationEqual, role1)),
			want: &domain.Authorization{
				ID:         authorizationWithRolesUser.ID,
				UserID:     userID,
				ProjectID:  projectID,
				InstanceID: instanceID,
				State:      domain.AuthorizationStateActive,
				Roles:      []string{role1},
				CreatedAt:  authorizationWithRolesUser.CreatedAt,
				UpdatedAt:  authorizationWithRolesUser.UpdatedAt,
			},
		},
		{
			name: "get authorization with non-existent role, error",
			condition: database.And(authorizationRepo.PrimaryKeyCondition(
				authorizationWithRolesUser.InstanceID,
				authorizationWithRolesUser.ID,
			),
				authorizationRepo.RoleCondition(database.TextOperationEqual, "random-role")),
			wantErr: new(database.NoRowFoundError),
		},
		{
			name: "get authorization with grant ID condition (with roles)",
			condition: database.And(authorizationRepo.PrimaryKeyCondition(
				authorizationProjectGrantWithRoles.InstanceID,
				authorizationProjectGrantWithRoles.ID,
			),
				authorizationRepo.GrantIDCondition(grantID)),
			want: authorizationProjectGrantWithRoles,
		},
		{
			name: "get authorization with grant ID condition (without roles)",
			condition: database.And(authorizationRepo.PrimaryKeyCondition(
				authorizationProjectGrantWithoutRoles.InstanceID,
				authorizationProjectGrantWithoutRoles.ID,
			),
				authorizationRepo.GrantIDCondition(grantID)),
			want: authorizationProjectGrantWithoutRoles,
		},
		{
			name: "get authorization for project grant with non-existent grant ID, error",
			condition: database.And(authorizationRepo.PrimaryKeyCondition(
				authorizationProjectGrantWithoutRoles.InstanceID,
				authorizationProjectGrantWithoutRoles.ID,
			),
				authorizationRepo.GrantIDCondition("random-id")),
			wantErr: new(database.NoRowFoundError),
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
	project1Role2 := createProjectRole(t, tx, instanceID, organizationID, project1ID, "project1Role2")

	project2ID := createProject(t, tx, instanceID, organizationID)
	require.NotNil(t, project2ID)
	project2Role1 := createProjectRole(t, tx, instanceID, organizationID, project2ID, "project2Role1")
	project2Role2 := createProjectRole(t, tx, instanceID, organizationID, project2ID, "project2Role2")

	user1ID := createHumanUser(t, tx, instanceID, organizationID)
	require.NotNil(t, user1ID)
	user2ID := createHumanUser(t, tx, instanceID, organizationID)
	require.NotNil(t, user2ID)

	// create authorization with roles for user1 for project1
	authorizationRepo := repository.AuthorizationRepository()
	authorizationWithRolesUser1Project1 := &domain.Authorization{
		ID:         "authorization1-user1",
		UserID:     user1ID,
		ProjectID:  project1ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project1Role1, project1Role2},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser1Project1)
	require.NoError(t, err)

	// create authorization without roles for user1 for project1
	authorizationWithoutRolesUser1Project1 := &domain.Authorization{
		ID:         "authorization2-user1",
		UserID:     user1ID,
		ProjectID:  project1ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser1Project1)
	require.NoError(t, err)

	// create authorization with roles for user1 for project2
	authorizationWithRolesUser1Project2 := &domain.Authorization{
		ID:         "authorization3-user1",
		UserID:     user1ID,
		ProjectID:  project2ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project2Role1, project2Role2},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser1Project2)
	require.NoError(t, err)

	// create authorization without roles for user1 for project2
	authorizationWithoutRolesUser1Project2 := &domain.Authorization{
		ID:        "authorization4-user1",
		UserID:    user1ID,
		ProjectID: project2ID,

		InstanceID: instanceID,
		State:      domain.AuthorizationStateInactive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser1Project2)
	require.NoError(t, err)

	// create authorization with roles for user2 for project1
	authorizationWithRolesUser2Project1 := &domain.Authorization{
		ID:        "authorization5-user2",
		UserID:    user2ID,
		ProjectID: project1ID,

		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project1Role1, project1Role2},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser2Project1)
	require.NoError(t, err)

	// create authorization without roles for user2 for project1
	authorizationWithoutRolesUser2Project1 := &domain.Authorization{
		ID:         "authorization6-user2",
		UserID:     user2ID,
		ProjectID:  project1ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateInactive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser2Project1)
	require.NoError(t, err)

	// create authorization with roles for user2 for project2
	authorizationWithRolesUser2Project2 := &domain.Authorization{
		ID:         "authorization7-user2",
		UserID:     user2ID,
		ProjectID:  project2ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project2Role1, project2Role2},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithRolesUser2Project2)
	require.NoError(t, err)

	// create authorization without roles for user2 for project2
	authorizationWithoutRolesUser2Project2 := &domain.Authorization{
		ID:         "authorization8-user2",
		UserID:     user2ID,
		ProjectID:  project2ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRolesUser2Project2)
	require.NoError(t, err)

	// create a project grant
	grantedOrganizationID := createOrganization(t, tx, instanceID)
	grantID := createProjectGrant(t, tx, instanceID, organizationID, grantedOrganizationID, project1ID, []string{project1Role1, project1Role2})

	authorizationProjectGrantWithRoles := &domain.Authorization{
		ID:         "authorization9-user1",
		UserID:     user1ID,
		ProjectID:  project1ID,
		GrantID:    gu.Ptr(grantID),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{project1Role1},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationProjectGrantWithRoles)
	require.NoError(t, err)

	authorizationProjectGrantWithoutRoles := &domain.Authorization{
		ID:         "authorization10-user2",
		UserID:     user2ID,
		ProjectID:  project1ID,
		GrantID:    gu.Ptr(grantID),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationProjectGrantWithoutRoles)
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
				authorizationProjectGrantWithoutRoles,
				authorizationWithoutRolesUser1Project1,
				authorizationWithRolesUser1Project2,
				authorizationWithoutRolesUser1Project2,
				authorizationWithRolesUser2Project1,
				authorizationWithoutRolesUser2Project1,
				authorizationWithRolesUser2Project2,
				authorizationWithoutRolesUser2Project2,
				authorizationProjectGrantWithRoles,
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
				database.WithOrderByAscending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationWithRolesUser1Project1,
				authorizationWithoutRolesUser1Project1,
				authorizationWithRolesUser1Project2,
				authorizationWithoutRolesUser1Project2,
				authorizationProjectGrantWithRoles,
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
				authorizationProjectGrantWithoutRoles,
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
				authorizationProjectGrantWithoutRoles,
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
				authorizationProjectGrantWithRoles,
				authorizationWithoutRolesUser2Project1,
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
						authorizationRepo.RoleCondition(database.TextOperationContains, project1Role1),
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
				{
					ID:         "authorization9-user1",
					UserID:     user1ID,
					ProjectID:  project1ID,
					GrantID:    gu.Ptr(grantID),
					InstanceID: instanceID,
					State:      domain.AuthorizationStateActive,
					Roles:      []string{project1Role1},
					CreatedAt:  authorizationProjectGrantWithRoles.CreatedAt,
					UpdatedAt:  authorizationProjectGrantWithRoles.UpdatedAt,
				},
			},
		},
		{
			name: "list all authorizations for project grant",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.GrantIDCondition(grantID),
					),
				),
				database.WithOrderByDescending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationProjectGrantWithRoles,
				authorizationProjectGrantWithoutRoles,
			},
		},
		{
			name: "list all authorizations for non-existent project grant",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.GrantIDCondition("random-id"),
					),
				),
			},
			want: []*domain.Authorization{},
		},
		{
			name: "list all authorizations with exists role condition",
			condition: []database.QueryOption{
				database.WithCondition(
					database.And(
						authorizationRepo.InstanceIDCondition(instanceID),
						authorizationRepo.ExistsRole(authorizationRepo.RoleCondition(database.TextOperationEqual, project2Role2)),
					),
				),
				database.WithOrderByAscending(authorizationRepo.IDColumn()),
			},
			want: []*domain.Authorization{
				authorizationWithRolesUser1Project2,
				authorizationWithRolesUser2Project2,
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
	beforeUpdate := time.Now()

	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, organizationID)
	userID := createHumanUser(t, tx, instanceID, organizationID)

	authorizationRepo := repository.AuthorizationRepository()
	activeAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := authorizationRepo.Create(t.Context(), tx, activeAuthorization)
	require.NoError(t, err)

	inactiveAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateInactive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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

			rowsAffected, err := authorizationRepo.Update(t.Context(), savepoint, tt.condition, nil, tt.changes...)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			afterUpdate := time.Now()
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)

			if tt.wantRowsAffected == 0 {
				return
			}
			// verify update
			updatedAuthorization, err := authorizationRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.NoError(t, err)
			assert.Equal(t, tt.wantState, updatedAuthorization.State)
			assert.WithinRange(t, updatedAuthorization.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}

}

func TestAuthorizationUpdate_setRoles(t *testing.T) {
	beforeUpdate := time.Now()

	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, organizationID)
	role1 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role1")
	role2 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role2")
	role3 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role3")
	role4 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role4")

	project2ID := createProject(t, tx, instanceID, organizationID)
	role5 := createProjectRole(t, tx, instanceID, organizationID, project2ID, "role5")
	role6 := createProjectRole(t, tx, instanceID, organizationID, project2ID, "role6")

	userID := createHumanUser(t, tx, instanceID, organizationID)

	authorizationRepo := repository.AuthorizationRepository()
	authorization1 := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := authorizationRepo.Create(t.Context(), tx, authorization1)
	require.NoError(t, err)

	authorization2 := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  project2ID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role5, role6},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorization2)
	require.NoError(t, err)

	tests := []struct {
		name             string
		roles            []string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:      "incomplete primary key condition - missing instance ID",
			roles:     []string{role1, role2},
			condition: authorizationRepo.IDCondition(authorization1.ID),
			wantErr:   new(database.MissingConditionError),
		},
		{
			name:      "incomplete primary key condition - missing authorization ID",
			roles:     []string{role1, role2},
			condition: authorizationRepo.InstanceIDCondition(authorization1.InstanceID),
			wantErr:   new(database.MissingConditionError),
		},
		{
			name:      "non-existent role",
			roles:     []string{"random-role"},
			condition: authorizationRepo.PrimaryKeyCondition(authorization1.InstanceID, authorization1.ID),
			wantErr:   new(database.ForeignKeyError),
		},
		{
			name:             "set single role",
			roles:            []string{role1},
			condition:        authorizationRepo.PrimaryKeyCondition(authorization1.InstanceID, authorization1.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "set multiple roles",
			roles:            []string{role1, role2},
			condition:        authorizationRepo.PrimaryKeyCondition(authorization1.InstanceID, authorization1.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "overwrite roles",
			roles:            []string{role3, role4},
			condition:        authorizationRepo.PrimaryKeyCondition(authorization1.InstanceID, authorization1.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "remove role",
			roles:            []string{role3},
			condition:        authorizationRepo.PrimaryKeyCondition(authorization1.InstanceID, authorization1.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "no changes",
			roles:            []string{role5, role6},
			condition:        authorizationRepo.PrimaryKeyCondition(authorization2.InstanceID, authorization2.ID),
			wantRowsAffected: 1,
		},
		{
			name:             "clear roles",
			roles:            []string{},
			condition:        authorizationRepo.PrimaryKeyCondition(authorization2.InstanceID, authorization2.ID),
			wantRowsAffected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			rowsAffected, err := authorizationRepo.Update(t.Context(), savepoint, tt.condition, tt.roles)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			afterUpdate := time.Now()
			require.Equal(t, tt.wantRowsAffected, rowsAffected)

			// get roles
			updatedAuthorization, err := authorizationRepo.Get(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.roles, updatedAuthorization.Roles)
			assert.WithinRange(t, updatedAuthorization.UpdatedAt, beforeUpdate, afterUpdate)
		})
	}
}

func TestDeleteAuthorization(t *testing.T) {
	tx, rollback := transactionForRollback(t)
	defer rollback()

	instanceID := createInstance(t, tx)
	organizationID := createOrganization(t, tx, instanceID)
	projectID := createProject(t, tx, instanceID, organizationID)
	role1 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role1")
	role2 := createProjectRole(t, tx, instanceID, organizationID, projectID, "role2")
	userID := createHumanUser(t, tx, instanceID, organizationID)

	authorizationRepo := repository.AuthorizationRepository()
	authorizationWithRoles := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1, role2},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := authorizationRepo.Create(t.Context(), tx, authorizationWithRoles)
	require.NoError(t, err)

	authorizationWithoutRoles := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationWithoutRoles)
	require.NoError(t, err)

	deletedAuthorization := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, deletedAuthorization)
	require.NoError(t, err)
	_, err = authorizationRepo.Delete(t.Context(), tx, authorizationRepo.PrimaryKeyCondition(
		instanceID,
		deletedAuthorization.ID,
	))
	require.NoError(t, err)

	grantedOrganizationID := createOrganization(t, tx, instanceID)
	grantID := createProjectGrant(t, tx, instanceID, organizationID, grantedOrganizationID, projectID, []string{role1})
	authorizationProjectGrant := &domain.Authorization{
		ID:         integration.ID(),
		UserID:     userID,
		ProjectID:  projectID,
		GrantID:    gu.Ptr(grantID),
		InstanceID: instanceID,
		State:      domain.AuthorizationStateActive,
		Roles:      []string{role1},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = authorizationRepo.Create(t.Context(), tx, authorizationProjectGrant)
	require.NoError(t, err)

	tests := []struct {
		name             string
		condition        database.Condition
		wantRowsAffected int64
		wantErr          error
	}{
		{
			name:      "incomplete condition",
			condition: authorizationRepo.IDCondition(authorizationWithRoles.ID),
			wantErr:   new(database.MissingConditionError),
		},
		{
			name: "delete non-existent authorization",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				"random-id",
			),
			wantRowsAffected: 0,
		},
		{
			name: "delete authorization with roles",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				authorizationWithRoles.ID,
			),
			wantRowsAffected: 1,
		},
		{
			name: "delete authorization without roles",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				authorizationWithoutRoles.ID,
			),
			wantRowsAffected: 1,
		},
		{
			name: "delete again",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				deletedAuthorization.ID,
			),
			wantRowsAffected: 0,
		},
		{
			name: "delete project grant authorization",
			condition: authorizationRepo.PrimaryKeyCondition(
				instanceID,
				authorizationProjectGrant.ID,
			),
			wantRowsAffected: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savepoint, rollback := savepointForRollback(t, tx)
			defer rollback()

			rowsAffected, err := authorizationRepo.Delete(t.Context(), savepoint, tt.condition)
			require.ErrorIs(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRowsAffected, rowsAffected)

			// verify deletion
			got, err := authorizationRepo.List(t.Context(), savepoint, database.WithCondition(tt.condition))
			require.NoError(t, err)
			require.Empty(t, got)
		})
	}
}
