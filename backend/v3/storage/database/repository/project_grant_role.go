package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type projectGrantRole struct{}

func (p projectGrantRole) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.ProjectGrantRole, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getMany[domain.ProjectGrantRole](ctx, client, builder)
}

const insertProjectGrantRoleStmt = `INSERT INTO zitadel.project_grant_roles(
	instance_id, grant_id, key, project_org_id, project_id
)
VALUES ($1, $2, $3, $4, $5)
RETURNING created_at`

func (p projectGrantRole) Add(ctx context.Context, client database.QueryExecutor, role *domain.ProjectGrantRole) error {
	builder := database.NewStatementBuilder(insertProjectGrantRoleStmt,
		role.InstanceID,
		role.ProjectID,
		role.Key,
		role.ProjectOrgID,
		role.ProjectID,
	)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).
		Scan(&role.CreatedAt)
}

func (p projectGrantRole) Remove(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return deleteOne(ctx, client, p, condition)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (p projectGrantRole) PrimaryKeyCondition(instanceID, grantID, key string) database.Condition {
	return database.And(
		p.InstanceIDCondition(instanceID),
		p.GrantIDCondition(grantID),
		p.KeyCondition(key),
	)
}

func (p projectGrantRole) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (p projectGrantRole) GrantIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.GrantIDColumn(), database.TextOperationEqual, instanceID)
}

func (p projectGrantRole) KeyCondition(key string) database.Condition {
	return database.NewTextCondition(p.KeyColumn(), database.TextOperationEqual, key)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (projectGrantRole) unqualifiedTableName() string {
	return "project_grant_roles"
}

// PrimaryKeyColumns implements the [pkRepository] interface
func (p projectGrantRole) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		p.InstanceIDColumn(),
		p.GrantIDColumn(),
		p.KeyColumn(),
	}
}

func (p projectGrantRole) InstanceIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "instance_id")
}

func (p projectGrantRole) GrantIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "grant_id")
}

func (p projectGrantRole) KeyColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "key")
}

func (p projectGrantRole) CreatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "created_at")
}

func (p projectGrantRole) ProjectOrgIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "project_org_id")
}

func (p projectGrantRole) ProjectIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "project_id")
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryProjectGrantRoleStmt = `SELECT
	project_grant_roles.instance_id,	
	project_grant_roles.grant_id,	
	project_grant_roles.key,
	project_grant_roles.created_at,
	project_grant_roles.project_org_id,
	project_grant_roles.project_id
	FROM zitadel.project_grant_roles`

func (p projectGrantRole) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, p.InstanceIDColumn(), p.ProjectOrgIDColumn(), p.ProjectIDColumn(), p.GrantIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryProjectGrantRoleStmt)
	options.Write(builder)

	return builder, nil
}
