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

func (projectGrant) Role() domain.ProjectGrantRoleRepository {
	return projectGrantRole{}
}

func (p projectGrant) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.ProjectGrant, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.ProjectGrant](ctx, client, builder)
}

func (p projectGrant) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.ProjectGrant, error) {
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

func (projectGrant) Create(ctx context.Context, client database.QueryExecutor, projectGrant *domain.ProjectGrant) error {
	builder := database.NewStatementBuilder(insertProjectGrantStmt,
		projectGrant.InstanceID,
		projectGrant.ID,
		projectGrant.ProjectID,
		projectGrant.GrantingOrganizationID,
		projectGrant.GrantedOrganizationID,
		projectGrant.State,
	)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).
		Scan(&projectGrant.CreatedAt, &projectGrant.UpdatedAt)
}

func (p projectGrant) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return updateOne(ctx, client, p, condition, changes...)
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

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (projectGrant) unqualifiedTableName() string {
	return "project_grants"
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
	project_grants.instance_id,	
	project_grants.id,
	project_grants.project_id,
	project_grants.granting_organization_id,
	project_grants.granted_organization_id,
	project_grants.created_at,
	project_grants.updated_at,
	project_grants.state
	FROM zitadel.project_grants`

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
