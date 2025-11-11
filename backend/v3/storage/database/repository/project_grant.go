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

const insertProjectGrantStmt = `INSERT INTO zitadel.project_grants(
	instance_id, id, project_id, granting_organization_id, granted_organization_id, state
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING created_at, updated_at`

const insertProjectGrantWithRolesStmt = `WITH added_roles AS (
	INSERT INTO zitadel.project_grant_roles (
		instance_id, grant_id, project_id, key
	)
	VALUES ($1, $2, $3, unnest($7::text[]))
) ` + insertProjectGrantStmt

func (p projectGrant) Create(ctx context.Context, client database.QueryExecutor, projectGrant *domain.ProjectGrant) error {
	var builder *database.StatementBuilder
	if len(projectGrant.RoleKeys) == 0 {
		builder = database.NewStatementBuilder(insertProjectGrantStmt,
			projectGrant.InstanceID,
			projectGrant.ID,
			projectGrant.ProjectID,
			projectGrant.GrantingOrganizationID,
			projectGrant.GrantedOrganizationID,
			projectGrant.State,
		)
	} else {
		builder = database.NewStatementBuilder(insertProjectGrantWithRolesStmt,
			projectGrant.InstanceID,
			projectGrant.ID,
			projectGrant.ProjectID,
			projectGrant.GrantingOrganizationID,
			projectGrant.GrantedOrganizationID,
			projectGrant.State,
			projectGrant.RoleKeys,
		)
	}

	if err := client.QueryRow(ctx, builder.String(), builder.Args()...).
		Scan(&projectGrant.CreatedAt, &projectGrant.UpdatedAt); err != nil {
		return err
	}
	return nil
}

const queryUpdateProjectGrantRoleStmt = `SELECT instance_id, id, project_id, $1::TEXT[] AS keys from zitadel.project_grants`

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
UPDATE zitadel.project_grants
SET updated_at = now()
FROM pg
WHERE pg.instance_id = project_grants.instance_id
  AND pg.id = project_grants.id`

func (p projectGrant) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return updateOne(ctx, client, p, condition, changes...)
}

func (p projectGrant) SetRoleKeys(ctx context.Context, client database.QueryExecutor, condition database.Condition, roleKeys []string) (int64, error) {
	if err := checkPKCondition(p, condition); err != nil {
		return 0, err
	}

    // the statement begins with getting the project grant we want to update
    builder := database.NewStatementBuilder("WITH pg AS (SELECT instance_id, id, project_id, unnest($1::text[]) as key from zitadel.project_grants", roleKeys)
	writeCondition(builder, condition)
    builder.WriteString(" ), ")
    // now we add the logic to do the required changes based on the given project grant
    builder.WriteString(updateProjectGrantRoleStmt)
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
