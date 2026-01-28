package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type projectGrant struct{}

// ProjectGrantRepository manages project grants.
func ProjectGrantRepository() domain.ProjectGrantRepository {
	return new(projectGrant)
}

func (p projectGrant) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.ProjectGrant, error) {
	opts = append(opts,
		database.WithGroupBy(p.InstanceIDColumn(), p.IDColumn()),
	)

	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.ProjectGrant](ctx, client, builder)
}

func (p projectGrant) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.ProjectGrant, error) {
	opts = append(opts,
		database.WithGroupBy(p.InstanceIDColumn(), p.IDColumn()),
	)

	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getMany[domain.ProjectGrant](ctx, client, builder)
}

const insertProjectGrantRolesStmt = `WITH added_roles AS (
	INSERT INTO zitadel.project_grant_roles (
		instance_id, grant_id, project_id, key
	)
	VALUES ($2, $3, $4, unnest($1::text[]))
) `

func (p projectGrant) Create(ctx context.Context, client database.QueryExecutor, projectGrant *domain.ProjectGrant) error {
	var (
		createdAt, updatedAt any = database.DefaultInstruction, database.DefaultInstruction
	)
	if !projectGrant.CreatedAt.IsZero() {
		createdAt = projectGrant.CreatedAt
	}
	if !projectGrant.UpdatedAt.IsZero() {
		updatedAt = projectGrant.UpdatedAt
	}

	// separate statement to add roles to project grant
	builder := database.NewStatementBuilder(insertProjectGrantRolesStmt, projectGrant.RoleKeys)

	builder.WriteString(`INSERT INTO ` + p.qualifiedTableName() + ` (instance_id, id, project_id, granting_organization_id, granted_organization_id, state, created_at, updated_at) VALUES ( `)
	builder.WriteArgs(
		projectGrant.InstanceID,
		projectGrant.ID,
		projectGrant.ProjectID,
		projectGrant.GrantingOrganizationID,
		projectGrant.GrantedOrganizationID,
		projectGrant.State,
		createdAt,
		updatedAt,
	)
	builder.WriteString(` ) RETURNING created_at, updated_at`)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).
		Scan(&projectGrant.CreatedAt, &projectGrant.UpdatedAt)
}

const queryUpdateProjectGrantRoleStmt = `SELECT instance_id, id, project_id, $1::text[] as keys from zitadel.project_grants`

const updateProjectGrantRoleStmt = `removed_roles AS (
    DELETE FROM zitadel.project_grant_roles as pgr
    USING pg
    WHERE pg.instance_id = pgr.instance_id AND pg.id = pgr.grant_id AND NOT(pgr.key = ANY(pg.keys))
), added_roles AS (
	INSERT INTO zitadel.project_grant_roles (
		instance_id, grant_id, project_id, key
	)
	SELECT instance_id, id, project_id, unnest(pg.keys) FROM pg
	ON CONFLICT (instance_id, grant_id, key) DO NOTHING
)
UPDATE zitadel.project_grants SET `

const updateProjectGrantRoleStmtWhere = ` FROM pg
WHERE pg.instance_id = project_grants.instance_id
  AND pg.id = project_grants.id`

func (p projectGrant) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, roleKeys []string, changes ...database.Change) (int64, error) {
	// if no role keys set we only have to update the project grant table
	if roleKeys == nil {
		return updateOne(ctx, client, p, condition, changes...)
	}

	// if you want to update the roles you have to have the primary key, otherwise multiple project grants get updated
	if err := checkPKCondition(p, condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(p.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(p.UpdatedAtColumn(), database.NullInstruction))
	}

	// the statement begins with getting the project grant we want to update
	builder := database.NewStatementBuilder("WITH pg AS (")
	builder.WriteString(queryUpdateProjectGrantRoleStmt)
	builder.AppendArg(roleKeys)
	writeCondition(builder, condition)
	builder.WriteString(" ), ")

	// now we add the logic to do the required changes based on the given project grant
	builder.WriteString(updateProjectGrantRoleStmt)
	if err := database.Changes(changes).Write(builder); err != nil {
		return 0, err
	}
	builder.WriteString(updateProjectGrantRoleStmtWhere)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (p projectGrant) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return deleteOne(ctx, client, p, condition)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (p projectGrant) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(p.UpdatedAtColumn(), updatedAt)
}

func (p projectGrant) SetState(state domain.ProjectGrantState) database.Change {
	return database.NewChange(p.StateColumn(), state)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (p projectGrant) PrimaryKeyCondition(instanceID, projectGrantID string) database.Condition {
	return database.And(
		p.InstanceIDCondition(instanceID),
		p.IDCondition(projectGrantID),
	)
}

func (p projectGrant) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (p projectGrant) OrganizationCondition(organizationID string) database.Condition {
	return database.Or(
		p.GrantingOrganizationIDCondition(organizationID),
		p.GrantedOrganizationIDCondition(organizationID),
	)
}

func (p projectGrant) IDCondition(projectGrantID string) database.Condition {
	return database.NewTextCondition(p.IDColumn(), database.TextOperationEqual, projectGrantID)
}

func (p projectGrant) ProjectIDCondition(projectID string) database.Condition {
	return database.NewTextCondition(p.ProjectIDColumn(), database.TextOperationEqual, projectID)
}

func (p projectGrant) GrantingOrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(p.GrantingOrganizationIDColumn(), database.TextOperationEqual, organizationID)
}

func (p projectGrant) GrantedOrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(p.GrantedOrganizationIDColumn(), database.TextOperationEqual, organizationID)
}

func (p projectGrant) StateCondition(state domain.ProjectGrantState) database.Condition {
	return database.NewTextCondition(p.StateColumn(), database.TextOperationEqual, state.String())
}

func (p projectGrant) RoleKeyCondition(op database.TextOperation, role string) database.Condition {
	return database.NewTextCondition(database.NewColumn(p.unqualifiedRolesTableName(), "key"), op, role)
}

// ExistsRoleKey creates a correlated [database.Exists] condition on project_grant_roles.
// Use this filter to make sure the project grant returned contains a specific project grant role.
// Example usage:
//
//	projectGrant, _ := projectGrantRepo.Get(ctx,
//	    database.WithCondition(
//	        database.And(
//	            projectGrantRepo.InstanceIDCondition(instanceID),
//	            projectGrantRepo.ExistsRoleKey(projectGrantRepo.RoleKeyCondition(database.TextOperationEqual, "admin")),
//	        ),
//	    ),
//	)
func (p projectGrant) ExistsRoleKey(cond database.Condition) database.Condition {
	return database.Exists(
		p.qualifiedRolesTableName(),
		database.And(
			database.NewColumnCondition(p.InstanceIDColumn(), database.NewColumn(p.unqualifiedRolesTableName(), "instance_id")),
			database.NewColumnCondition(p.IDColumn(), database.NewColumn(p.unqualifiedRolesTableName(), "grant_id")),
			cond,
		),
	)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (p projectGrant) qualifiedTableName() string {
	return "zitadel." + p.unqualifiedTableName()
}

func (projectGrant) unqualifiedTableName() string {
	return "project_grants"
}

func (p projectGrant) qualifiedRolesTableName() string {
	return "zitadel." + p.unqualifiedRolesTableName()
}

func (projectGrant) unqualifiedRolesTableName() string {
	return "project_grant_roles"
}

func (p projectGrant) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		p.InstanceIDColumn(),
		p.IDColumn(),
	}
}

func (p projectGrant) InstanceIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "instance_id")
}

func (p projectGrant) IDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "id")
}

func (p projectGrant) ProjectIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "project_id")
}

func (p projectGrant) GrantingOrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "granting_organization_id")
}

func (p projectGrant) GrantedOrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "granted_organization_id")
}

func (p projectGrant) CreatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "created_at")
}

func (p projectGrant) UpdatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "updated_at")
}

func (p projectGrant) StateColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "state")
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryProjectGrantStmt = `SELECT
	zitadel.project_grants.instance_id,	
	zitadel.project_grants.id,
	zitadel.project_grants.project_id,
	zitadel.project_grants.granting_organization_id,
	zitadel.project_grants.granted_organization_id,
	zitadel.project_grants.created_at,
	zitadel.project_grants.updated_at,
	zitadel.project_grants.state,
    ARRAY_AGG(project_grant_roles.key) FILTER (WHERE project_grant_roles.grant_id IS NOT NULL) AS role_keys
	FROM zitadel.project_grants
	LEFT JOIN zitadel.project_grant_roles ON zitadel.project_grant_roles.instance_id = zitadel.project_grants.instance_id AND zitadel.project_grant_roles.grant_id = zitadel.project_grants.id`

func (p projectGrant) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, p.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryProjectGrantStmt)
	options.Write(builder)

	return builder, nil
}
