package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type projectRole struct{}

func (p projectRole) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.ProjectRole, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.ProjectRole](ctx, client, builder)
}

func (p projectRole) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.ProjectRole, error) {
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

func (p projectRole) Create(ctx context.Context, client database.QueryExecutor, role *domain.ProjectRole) error {
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

func (p projectRole) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := checkPKCondition(p, condition); err != nil {
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

func (p projectRole) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkPKCondition(p, condition); err != nil {
		return 0, err
	}
	builder := database.NewStatementBuilder(`DELETE FROM zitadel.project_roles`)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (p projectRole) SetDisplayName(displayName string) database.Change {
	return database.NewChange(p.DisplayNameColumn(), displayName)
}

func (p projectRole) SetRoleGroup(roleGroup string) database.Change {
	return database.NewChange(p.RoleGroupColumn(), roleGroup)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (p projectRole) PrimaryKeyCondition(instanceID, projectID, key string) database.Condition {
	return database.And(
		p.InstanceIDCondition(instanceID),
		p.ProjectIDCondition(projectID),
		p.KeyCondition(key),
	)
}

func (p projectRole) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (p projectRole) ProjectIDCondition(projectID string) database.Condition {
	return database.NewTextCondition(p.ProjectIDColumn(), database.TextOperationEqual, projectID)
}

func (p projectRole) KeyCondition(key string) database.Condition {
	return database.NewTextCondition(p.KeyColumn(), database.TextOperationEqual, key)
}

func (p projectRole) DisplayNameCondition(op database.TextOperation, displayName string) database.Condition {
	return database.NewTextCondition(p.DisplayNameColumn(), op, displayName)
}

func (p projectRole) RoleGroupCondition(op database.TextOperation, roleGroup string) database.Condition {
	return database.NewTextCondition(p.RoleGroupColumn(), op, roleGroup)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (projectRole) unqualifiedTableName() string {
	return "project_roles"
}

// PrimaryKeyColumns implements the [pkRepository] interface
func (p projectRole) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		p.InstanceIDColumn(),
		p.ProjectIDColumn(),
		p.KeyColumn(),
	}
}

func (p projectRole) InstanceIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "instance_id")
}

func (p projectRole) OrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "organization_id")
}

func (p projectRole) ProjectIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "project_id")
}

func (p projectRole) CreatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "created_at")
}

func (p projectRole) UpdatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "updated_at")
}

func (p projectRole) KeyColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "key")
}

func (p projectRole) DisplayNameColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "display_name")
}

func (p projectRole) RoleGroupColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "role_group")
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryProjectRoleStmt = `SELECT
	project_roles.instance_id,	
	project_roles.organization_id,
	project_roles.project_id,
	project_roles.created_at,
	project_roles.updated_at,
	project_roles.key,
	project_roles.display_name,
	project_roles.role_group
	FROM zitadel.project_roles`

func (p projectRole) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, p.InstanceIDColumn(), p.ProjectIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryProjectRoleStmt)
	options.Write(builder)

	return builder, nil
}
