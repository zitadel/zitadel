package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.AuthorizationRepository = (*authorization)(nil)

type authorization struct{}

func AuthorizationRepository() domain.AuthorizationRepository {
	return new(authorization)
}

func (a authorization) unqualifiedTableName() string {
	return "authorizations"
}

func (a authorization) unqualifiedAuthorizationRolesTableName() string {
	return "authorization_roles"
}

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

const insertAuthorizationStmt = `INSERT INTO zitadel.authorizations (
	instance_id, id, user_id, grant_id, project_id, state
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING created_at, updated_at`

const insertAuthorizationWithRolesStmt = `WITH roles AS (
  INSERT INTO zitadel.authorization_roles
      (instance_id, authorization_id, grant_id, project_id, role_key)
      VALUES ($1, $2, $4, $5, unnest($7::text[])))` + insertAuthorizationStmt

// Create implements [domain.AuthorizationRepository].
func (a authorization) Create(ctx context.Context, client database.QueryExecutor, authorization *domain.Authorization) error {
	var builder *database.StatementBuilder
	if len(authorization.Roles) == 0 {
		builder = database.NewStatementBuilder(insertAuthorizationStmt,
			authorization.InstanceID,
			authorization.ID,
			authorization.UserID,
			authorization.GrantID,
			authorization.ProjectID,
			authorization.State,
		)
	} else {
		builder = database.NewStatementBuilder(insertAuthorizationWithRolesStmt,
			authorization.InstanceID,
			authorization.ID,
			authorization.UserID,
			authorization.GrantID,
			authorization.ProjectID,
			authorization.State,
			authorization.Roles,
		)
	}

	if err := client.QueryRow(
		ctx,
		builder.String(),
		builder.Args()...).
		Scan(
			&authorization.CreatedAt,
			&authorization.UpdatedAt,
		); err != nil {
		return err
	}
	return nil
}

const queryAuthorizationStmt = `SELECT zitadel.authorizations.instance_id,
       zitadel.authorizations.id,
       zitadel.authorizations.user_id,
       zitadel.authorizations.grant_id,
       zitadel.authorizations.project_id,
       zitadel.authorizations.state,
       zitadel.authorizations.created_at,
       zitadel.authorizations.updated_at,
       ARRAY_AGG(authorization_roles.role_key)
       FILTER (WHERE authorization_roles.authorization_id IS NOT NULL) AS roles
FROM zitadel.authorizations
         LEFT JOIN zitadel.authorization_roles
                   ON zitadel.authorizations.instance_id = zitadel.authorization_roles.instance_id
                       AND zitadel.authorizations.id = zitadel.authorization_roles.authorization_id`

// Get implements [domain.AuthorizationRepository].
func (a authorization) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Authorization, error) {
	opts = append(opts,
		database.WithGroupBy(a.InstanceIDColumn(), a.IDColumn()),
	)

	builder, err := a.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.Authorization](ctx, client, builder)
}

// List implements [domain.AuthorizationRepository].
func (a authorization) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Authorization, error) {
	opts = append(opts,
		database.WithGroupBy(a.InstanceIDColumn(), a.IDColumn()),
	)

	builder, err := a.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getMany[domain.Authorization](ctx, client, builder)
}

func (a authorization) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, a.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryAuthorizationStmt)
	options.Write(builder)

	return builder, nil
}

const queryUpdateAuthorizationRoleStmt = `SELECT instance_id, id, project_id, grant_id, $1::text[] as roles from zitadel.authorizations`

const updateAuthorizationRoleStmt = `deleted_roles AS (
    DELETE FROM zitadel.authorization_roles as azr
	USING az
        WHERE azr.instance_id = az.instance_id
            AND azr.authorization_id = az.id
            AND NOT azr.role_key = ANY ($1::text[])
        RETURNING *
), inserted_roles AS (
    INSERT INTO zitadel.authorization_roles (instance_id, authorization_id, project_id, grant_id, role_key)
        SELECT instance_id,
               id,
               project_id,
               grant_id,
               UNNEST(az.roles) AS role_key
        FROM az
        ON CONFLICT DO NOTHING
        RETURNING *
)
UPDATE zitadel.authorizations SET `

const updateAuthorizationRoleStmtWhere = ` FROM az
WHERE az.instance_id = authorizations.instance_id
  AND az.id = authorizations.id`

// Update implements [domain.AuthorizationRepository].
func (a authorization) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, roles []string, changes ...database.Change) (int64, error) {
	if roles != nil {
		return a.setRoles(ctx, client, condition, roles, changes...)
	}
	return a.update(ctx, client, condition, changes)
}

// update updates the authorization table when there are no roles to be set.
func (a authorization) update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes []database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if !condition.IsRestrictingColumn(a.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(a.InstanceIDColumn())
	}
	if !database.Changes(changes).IsOnColumn(a.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(a.UpdatedAtColumn(), database.NullInstruction))
	}

	builder := database.NewStatementBuilder("UPDATE zitadel.authorizations SET ")
	database.Changes(changes).Write(builder)
	writeCondition(builder, condition)

	stmt := builder.String()

	rowsAffected, err := client.Exec(ctx, stmt, builder.Args()...)
	return rowsAffected, err
}

// setRoles sets the roles of an authorization.
func (a authorization) setRoles(ctx context.Context, client database.QueryExecutor, condition database.Condition, roles []string, changes ...database.Change) (int64, error) {
	if err := checkPKCondition(a, condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(a.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(a.UpdatedAtColumn(), database.NullInstruction))
	}

	// get the authorization to be updated
	builder := database.NewStatementBuilder("WITH az AS (")
	builder.WriteString(queryUpdateAuthorizationRoleStmt)
	builder.AppendArg(roles)
	writeCondition(builder, condition)
	builder.WriteString(" ), ")

	// set the roles
	builder.WriteString(updateAuthorizationRoleStmt)
	database.Changes(changes).Write(builder)
	builder.WriteString(updateAuthorizationRoleStmtWhere)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// Delete implements [domain.AuthorizationRepository].
func (a authorization) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if !condition.IsRestrictingColumn(a.InstanceIDColumn()) {
		return 0, database.NewMissingConditionError(a.InstanceIDColumn())
	}

	builder := database.NewStatementBuilder("DELETE FROM zitadel.authorizations")
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.Repository]
func (a authorization) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		a.InstanceIDColumn(),
		a.IDColumn(),
	}
}

func (a authorization) IDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "id")
}

func (a authorization) UserIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "user_id")
}

func (a authorization) GrantIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "grant_id")
}

func (a authorization) ProjectIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "project_id")
}

func (a authorization) InstanceIDColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "instance_id")
}

func (a authorization) StateColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "state")
}

func (a authorization) CreatedAtColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "created_at")
}

func (a authorization) UpdatedAtColumn() database.Column {
	return database.NewColumn(a.unqualifiedTableName(), "updated_at")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// PrimaryKeyCondition implements [domain.authorizationConditions]
func (a authorization) PrimaryKeyCondition(instanceID, id string) database.Condition {
	return database.And(
		a.InstanceIDCondition(instanceID),
		a.IDCondition(id),
	)
}

// IDCondition implements [domain.authorizationConditions]
func (a authorization) IDCondition(id string) database.Condition {
	return database.NewTextCondition(a.IDColumn(), database.TextOperationEqual, id)
}

// InstanceIDCondition implements [domain.authorizationConditions]
func (a authorization) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(a.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

// ProjectIDCondition implements [domain.authorizationConditions]
func (a authorization) ProjectIDCondition(projectID string) database.Condition {
	return database.NewTextCondition(a.ProjectIDColumn(), database.TextOperationEqual, projectID)
}

// GrantIDCondition implements [domain.authorizationConditions]
func (a authorization) GrantIDCondition(grantID string) database.Condition {
	return database.NewTextCondition(a.GrantIDColumn(), database.TextOperationEqual, grantID)
}

// UserIDCondition implements [domain.authorizationConditions]
func (a authorization) UserIDCondition(userID string) database.Condition {
	return database.NewTextCondition(a.UserIDColumn(), database.TextOperationEqual, userID)
}

// RoleCondition implements [domain.authorizationConditions]
func (a authorization) RoleCondition(op database.TextOperation, role string) database.Condition {
	return database.NewTextCondition(database.NewColumn(a.unqualifiedAuthorizationRolesTableName(), "role_key"), op, role)
}

// StateCondition implements [domain.authorizationConditions]
func (a authorization) StateCondition(state domain.AuthorizationState) database.Condition {
	return database.NewTextCondition(a.StateColumn(), database.TextOperationEqual, state.String())
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// SetState implements [domain.authorizationChanges]
func (a authorization) SetState(state domain.AuthorizationState) database.Change {
	return database.NewChange(a.StateColumn(), state)
}

func (a authorization) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(a.UpdatedAtColumn(), updatedAt)
}
