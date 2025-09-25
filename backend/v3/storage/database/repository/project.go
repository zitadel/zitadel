package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

// -------------------------------------------------------------
// repository
// -------------------------------------------------------------

type project struct{}

func ProjectRepository() domain.ProjectRepository {
	return project{}
}

func (project) Role() domain.ProjectRoleRepository {
	return projectRoles{}
}

func (p project) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.Project, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return p.getOne(ctx, client, builder)
}

func (p project) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.Project, error) {
	builder, err := p.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return p.getMany(ctx, client, builder)
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
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	// TBD: do we support update operations of multiple projects?
	// In other words: should we require projectIDColumn as well?
	if err := p.checkRestrictingColumns(condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(p.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(p.UpdatedAtColumn(), database.NullInstruction))
	}
	builder := database.NewStatementBuilder(`UPDATE zitadel.projects SET `)
	database.Changes(changes).Write(builder)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (p project) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	// TBD: do we support delete operations of multiple projects?
	// In other words: should we require projectIDColumn as well?
	if err := p.checkRestrictingColumns(condition); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder(`DELETE FROM zitadel.projects`)
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

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

func (p project) InstanceIDCondition(instanceID string) database.Condition {
	return database.NewTextCondition(p.InstanceIDColumn(), database.TextOperationEqual, instanceID)
}

func (p project) OrganizationIDCondition(organizationID string) database.Condition {
	return database.NewTextCondition(p.OrganizationIDColumn(), database.TextOperationEqual, organizationID)
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

func (project) unqualifiedTableName() string {
	return "projects"
}

func (p project) InstanceIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "instance_id")
}

func (p project) OrganizationIDColumn() database.Column {
	return database.NewColumn(p.unqualifiedTableName(), "organization_id")
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
	projects.used_labeling_setting_owner
	FROM zitadel.projects`

func (p project) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := p.checkRestrictingColumns(options.Condition); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryProjectStmt)
	options.Write(builder)

	return builder, nil
}

func (p project) checkRestrictingColumns(condition database.Condition) error {
	return checkRestrictingColumns(
		condition,
		p.InstanceIDColumn(),
		p.OrganizationIDColumn(),
	)
}

func (project) getOne(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.Project, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	var project domain.Project
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(&project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (project) getMany(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.Project, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	var projects []*domain.Project
	if err := rows.(database.CollectableRows).Collect(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}
