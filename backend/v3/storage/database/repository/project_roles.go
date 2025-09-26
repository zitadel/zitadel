package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type projectRoles struct{}

func (p projectRoles) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.ProjectRole, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.ProjectRole](ctx, client, builder)
}

func (p projectRoles) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.ProjectRole, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getMany[domain.ProjectRole](ctx, client, builder)
}

const insertProjectRoleStmt = `INSERT INTO zitadel.project_roles(
	instance_id, organization_id, project_id, key, display_name, role_group
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING created_at, updated_at`

func (p projectRoles) Create(ctx context.Context, client database.QueryExecutor, role *domain.ProjectRole) error {
	builder := database.NewStatementBuilder(insertProjectRoleStmt,
		role.InstanceID,
		role.OrganizationID,
		role.ProjectID,
		role.Key,
		role.DisplayName,
		role.RoleGroup,
	)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).
		Scan(&role.CreatedAt, &role.UpdatedAt)
}

func (p projectRoles) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := p.checkRestrictingColumns(condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(p.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(p.UpdatedAtColumn(), database.NullInstruction))
	}
	builder := database.NewStatementBuilder(`UPDATE zitadel.project_roles SET `)
	database.Changes(changes).Write(builder)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (p projectRoles) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := p.checkRestrictingColumns(condition); err != nil {
		return 0, err
	}
	builder := database.NewStatementBuilder(`DELETE FROM zitadel.project_roles`)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (p projectRoles) SetDisplayName(displayName string) database.Change {
	return database.NewChange(p.DisplayNameColumn(), displayName)
}

func (p projectRoles) SetRoleGroup(roleGroup string) database.Change {
	return database.NewChange(p.RoleGroupColumn(), roleGroup)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (p projectRoles) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (p projectRoles) OrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(p.OrganizationIDColumn(), database.TextOperationEqual, organizationID)
}

func (p projectRoles) ProjectIDCondition(projectID string) database.Condition {
	return database.NewTextCondition(p.ProjectIDColumn(), database.TextOperationEqual, projectID)
}

func (p projectRoles) KeyCondition(key string) database.Condition {
	return database.NewTextCondition(p.KeyColumn(), database.TextOperationEqual, key)
}

func (p projectRoles) DisplayNameCondition(op database.TextOperation, displayName string) database.Condition {
	return database.NewTextCondition(p.DisplayNameColumn(), op, displayName)
}

func (p projectRoles) RoleGroupCondition(op database.TextOperation, roleGroup string) database.Condition {
	return database.NewTextCondition(p.RoleGroupColumn(), op, roleGroup)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (projectRoles) unqualifiedTableName() string {
	return "project_roles"
}

func (p projectRoles) InstanceIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "instance_id")
}

func (p projectRoles) OrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "organization_id")
}

func (p projectRoles) ProjectIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "project_id")
}

func (p projectRoles) CreatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "created_at")
}

func (p projectRoles) UpdatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "updated_at")
}

func (p projectRoles) KeyColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "key")
}

func (p projectRoles) DisplayNameColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "display_name")
}

func (p projectRoles) RoleGroupColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "role_group")
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryProjectRoleStmt = `SELECT
	projects.instance_id,	
	projects.organization_id,
	projects.project_id,
	projects.created_at,
	projects.updated_at,
	projects.key,
	projects.display_name,
	projects.role_group
	FROM zitadel.project_roles`

func (p projectRoles) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := p.checkRestrictingColumns(options.Condition); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryProjectRoleStmt)
	options.Write(builder)

	return builder, nil
}

func (p projectRoles) checkRestrictingColumns(condition database.Condition) error {
	return checkRestrictingColumns(
		condition,
		p.InstanceIDColumn(),
		p.OrganizationIDColumn(),
		p.ProjectIDColumn(),
	)
}
