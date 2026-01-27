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

type project struct{}

// ProjectRepository manages projects and project roles.
func ProjectRepository() domain.ProjectRepository {
	return project{}
}

func (project) Role() domain.ProjectRoleRepository {
	return projectRole{}
}

func (p project) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Project, error) {
	builder, err := p.prepareGetQuery(opts)
	if err != nil {
		return nil, err
	}
	return getOne[domain.Project](ctx, client, builder)
}

func (p project) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Project, error) {
	builder, err := p.prepareListQuery(opts)
	if err != nil {
		return nil, err
	}
	return getMany[domain.Project](ctx, client, builder)
}

const insertProjectStmt = `INSERT INTO zitadel.projects(
	instance_id, organization_id, id, name, state, should_assert_role, is_authorization_required, is_project_access_required, used_labeling_setting_owner
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING created_at, updated_at`

func (project) Create(ctx context.Context, client database.QueryExecutor, project *domain.Project) error {
	builder := database.NewStatementBuilder(insertProjectStmt,
		project.InstanceID,
		project.OrganizationID,
		project.ID,
		project.Name,
		project.State,
		project.ShouldAssertRole,
		project.IsAuthorizationRequired,
		project.IsProjectAccessRequired,
		project.UsedLabelingSettingOwner,
	)
	return client.QueryRow(ctx, builder.String(), builder.Args()...).
		Scan(&project.CreatedAt, &project.UpdatedAt)
}

func (p project) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return updateOne(ctx, client, p, condition, changes...)
}

func (p project) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	return deleteOne(ctx, client, p, condition)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (p project) SetUpdatedAt(updatedAt time.Time) database.Change {
	return database.NewChange(p.UpdatedAtColumn(), updatedAt)
}

func (p project) SetName(name string) database.Change {
	return database.NewChange(p.NameColumn(), name)
}

func (p project) SetState(state domain.ProjectState) database.Change {
	return database.NewChange(p.StateColumn(), state)
}

func (p project) SetShouldAssertRole(shouldAssertRole bool) database.Change {
	return database.NewChange(p.ShouldAssertRoleColumn(), shouldAssertRole)
}

func (p project) SetIsAuthorizationRequired(isAuthorizationRequired bool) database.Change {
	return database.NewChange(p.IsAuthorizationRequiredColumn(), isAuthorizationRequired)
}

func (p project) SetIsProjectAccessRequired(isProjectAccessRequired bool) database.Change {
	return database.NewChange(p.IsProjectAccessRequiredColumn(), isProjectAccessRequired)
}

func (p project) SetUsedLabelingSettingOwner(usedLabelingSettingOwner int16) database.Change {
	return database.NewChange(p.UsedLabelingSettingOwnerColumn(), usedLabelingSettingOwner)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (p project) PrimaryKeyCondition(instanceID, projectID string) database.Condition {
	return database.And(
		p.InstanceIDCondition(instanceID),
		p.IDCondition(projectID),
	)
}

func (p project) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (p project) OrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(p.OrganizationIDColumn(), database.TextOperationEqual, organizationID)
}

func (p project) GrantedOrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(p.GrantedOrganizationIDColumn(), database.TextOperationEqual, organizationID)
}

func (p project) IDCondition(projectID string) database.Condition {
	return database.NewTextCondition(p.IDColumn(), database.TextOperationEqual, projectID)
}

func (p project) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(p.NameColumn(), op, name)
}

func (p project) StateCondition(state domain.ProjectState) database.Condition {
	return database.NewTextCondition(p.StateColumn(), database.TextOperationEqual, state.String())
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (p project) qualifiedTableName() string {
	return "zitadel." + p.unqualifiedTableName()
}

func (project) unqualifiedTableName() string {
	return "projects"
}

// PrimaryKeyColumns implements the [pkRepository] interface
func (p project) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		p.InstanceIDColumn(),
		p.IDColumn(),
	}
}

func (p project) InstanceIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "instance_id")
}

func (p project) OrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "organization_id")
}

func (p project) GrantedOrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "granted_organization_id")
}

func (p project) IDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "id")
}

func (p project) CreatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "created_at")
}

func (p project) UpdatedAtColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "updated_at")
}

func (p project) NameColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "name")
}

func (p project) StateColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "state")
}

func (p project) ShouldAssertRoleColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "should_assert_role")
}

func (p project) IsAuthorizationRequiredColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "is_authorization_required")
}

func (p project) IsProjectAccessRequiredColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "is_project_access_required")
}

func (p project) UsedLabelingSettingOwnerColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "used_labeling_setting_owner")
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryProjectStmt = `SELECT
	projects.instance_id,	
	projects.organization_id,
	projects.id,
	projects.created_at,
	projects.updated_at,
	projects.name,
	projects.state,
	projects.should_assert_role,
	projects.is_authorization_required,
	projects.is_project_access_required,
	projects.used_labeling_setting_owner,
	NULL::TEXT AS grant_id, 
	NULL::TEXT AS granted_organization_id,  
	NULL::zitadel.project_grant_state AS grant_state    
	FROM zitadel.projects`

func (p project) prepareGetQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, p.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryProjectStmt)
	options.Write(builder)

	return builder, nil
}

const queryGrantedProjectStmt = `SELECT
	projects.instance_id,	
	projects.organization_id,
	projects.id as id,
	grants.created_at,
	grants.updated_at,
	projects.name,
	projects.state,
	projects.should_assert_role,
	projects.is_authorization_required,
	projects.is_project_access_required,
	projects.used_labeling_setting_owner,
	grants.id AS grant_id,
	grants.granted_organization_id,
	grants.state as grant_state
	FROM zitadel.projects 
	INNER JOIN zitadel.project_grants as grants
	    ON grants.instance_id = projects.instance_id 
	    AND grants.granting_organization_id = projects.organization_id
	    AND grants.project_id = projects.id`

const queryProjectsUnionStmt = `SELECT
	projects.instance_id,	
	projects.organization_id,
	projects.id as id,
	projects.created_at,
	projects.updated_at,
	projects.name,
	projects.state,
	projects.should_assert_role,
	projects.is_authorization_required,
	projects.is_project_access_required,
	projects.used_labeling_setting_owner,
	projects.grant_id,
	projects.granted_organization_id,
	projects.grant_state	
	FROM (` +
	queryProjectStmt +
	` UNION ALL ` +
	queryGrantedProjectStmt +
	`) AS projects`

func (p project) prepareListQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, p.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(
		queryProjectsUnionStmt,
	)
	options.Write(builder)

	return builder, nil
}
