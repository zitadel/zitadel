package query

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	projectsTable = table{
		name:          projection.ProjectProjectionTable,
		instanceIDCol: projection.ProjectColumnInstanceID,
	}
	ProjectColumnID = Column{
		name:  projection.ProjectColumnID,
		table: projectsTable,
	}
	ProjectColumnName = Column{
		name:  projection.ProjectColumnName,
		table: projectsTable,
	}
	ProjectColumnProjectRoleAssertion = Column{
		name:  projection.ProjectColumnProjectRoleAssertion,
		table: projectsTable,
	}
	ProjectColumnProjectRoleCheck = Column{
		name:  projection.ProjectColumnProjectRoleCheck,
		table: projectsTable,
	}
	ProjectColumnHasProjectCheck = Column{
		name:  projection.ProjectColumnHasProjectCheck,
		table: projectsTable,
	}
	ProjectColumnPrivateLabelingSetting = Column{
		name:  projection.ProjectColumnPrivateLabelingSetting,
		table: projectsTable,
	}
	ProjectColumnCreationDate = Column{
		name:  projection.ProjectColumnCreationDate,
		table: projectsTable,
	}
	ProjectColumnChangeDate = Column{
		name:  projection.ProjectColumnChangeDate,
		table: projectsTable,
	}
	ProjectColumnResourceOwner = Column{
		name:  projection.ProjectColumnResourceOwner,
		table: projectsTable,
	}
	ProjectColumnInstanceID = Column{
		name:  projection.ProjectColumnInstanceID,
		table: projectsTable,
	}
	ProjectColumnSequence = Column{
		name:  projection.ProjectColumnSequence,
		table: projectsTable,
	}
	ProjectColumnState = Column{
		name:  projection.ProjectColumnState,
		table: projectsTable,
	}
)

var (
	grantedProjectsAlias = table{
		name:          "granted_projects",
		instanceIDCol: projection.ProjectColumnInstanceID,
	}
	GrantedProjectColumnID = Column{
		name:  projection.ProjectColumnID,
		table: grantedProjectsAlias,
	}
	GrantedProjectColumnCreationDate = Column{
		name:  projection.ProjectColumnCreationDate,
		table: grantedProjectsAlias,
	}
	GrantedProjectColumnChangeDate = Column{
		name:  projection.ProjectColumnChangeDate,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnResourceOwner = Column{
		name:  projection.ProjectColumnResourceOwner,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnInstanceID = Column{
		name:  projection.ProjectGrantColumnInstanceID,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnState = Column{
		name:  "project_state",
		table: grantedProjectsAlias,
	}
	GrantedProjectColumnName = Column{
		name:  "project_name",
		table: grantedProjectsAlias,
	}
	grantedProjectColumnProjectRoleAssertion = Column{
		name:  projection.ProjectColumnProjectRoleAssertion,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnProjectRoleCheck = Column{
		name:  projection.ProjectColumnProjectRoleCheck,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnHasProjectCheck = Column{
		name:  projection.ProjectColumnHasProjectCheck,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnPrivateLabelingSetting = Column{
		name:  projection.ProjectColumnPrivateLabelingSetting,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnGrantResourceOwner = Column{
		name:  "project_grant_resource_owner",
		table: grantedProjectsAlias,
	}
	grantedProjectColumnGrantID = Column{
		name:  projection.ProjectGrantColumnGrantID,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnGrantedOrganization = Column{
		name:  projection.ProjectGrantColumnGrantedOrgID,
		table: grantedProjectsAlias,
	}
	grantedProjectColumnGrantedOrganizationName = Column{
		name:  "granted_org_name",
		table: grantedProjectsAlias,
	}
	grantedProjectColumnGrantState = Column{
		name:  "project_grant_state",
		table: grantedProjectsAlias,
	}
)

type Projects struct {
	SearchResponse
	Projects []*Project
}

func projectsCheckPermission(ctx context.Context, projects *Projects, permissionCheck domain.PermissionCheck) {
	projects.Projects = slices.DeleteFunc(projects.Projects,
		func(project *Project) bool {
			return projectCheckPermission(ctx, project.ResourceOwner, project.ID, permissionCheck) != nil
		},
	)
}

func projectCheckPermission(ctx context.Context, resourceOwner string, projectID string, permissionCheck domain.PermissionCheck) error {
	return permissionCheck(ctx, domain.PermissionProjectRead, resourceOwner, projectID)
}

type Project struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.ProjectState
	Sequence      uint64

	Name                   string
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
}

type ProjectSearchQueries struct {
	SearchRequest
	Queries      []SearchQuery
	GrantQueries []SearchQuery
}

func (q *Queries) GetProjectByIDWithPermission(ctx context.Context, shouldTriggerBulk bool, id string, permissionCheck domain.PermissionCheck) (*Project, error) {
	project, err := q.ProjectByID(ctx, shouldTriggerBulk, id)
	if err != nil {
		return nil, err
	}
	if err := projectCheckPermission(ctx, project.ResourceOwner, project.ID, permissionCheck); err != nil {
		return nil, err
	}
	return project, nil
}

func (q *Queries) ProjectByID(ctx context.Context, shouldTriggerBulk bool, id string) (project *Project, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerProjectProjection")
		ctx, err = projection.ProjectProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareProjectQuery()
	eq := sq.Eq{
		ProjectColumnID.identifier():         id,
		ProjectColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2m00Q", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		project, err = scan(row)
		return err
	}, query, args...)
	return project, err
}

func (q *Queries) SearchProjects(ctx context.Context, queries *ProjectSearchQueries, permissionCheck domain.PermissionCheck) (*Projects, error) {
	projects, err := q.searchProjects(ctx, queries)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil {
		projectsCheckPermission(ctx, projects, permissionCheck)
	}
	return projects, nil
}

func (q *Queries) searchProjects(ctx context.Context, queries *ProjectSearchQueries) (projects *Projects, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareProjectsQuery()
	eq := sq.Eq{ProjectColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-fn9ew", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		projects, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-2j00f", "Errors.Internal")
	}
	projects.State, err = q.latestState(ctx, projectsTable)
	return projects, err
}

type ProjectAndGrantedProjectSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *ProjectAndGrantedProjectSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func projectPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *ProjectAndGrantedProjectSearchQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		grantedProjectColumnResourceOwner,
		domain.PermissionProjectRead,
		SingleOrgPermissionOption(queries.Queries),
		WithProjectsPermissionOption(GrantedProjectColumnID),
	)
	return query.JoinClause(join, args...)
}

func (q *Queries) SearchGrantedProjects(ctx context.Context, queries *ProjectAndGrantedProjectSearchQueries, permissionCheck domain.PermissionCheck) (*GrantedProjects, error) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	projects, err := q.searchGrantedProjects(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil { // && !authz.GetFeatures(ctx).PermissionCheckV2 {
		grantedProjectsCheckPermission(ctx, projects, permissionCheck)
	}
	return projects, nil
}

func (q *Queries) searchGrantedProjects(ctx context.Context, queries *ProjectAndGrantedProjectSearchQueries, permissionCheckV2 bool) (grantedProjects *GrantedProjects, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareGrantedProjectsQuery()
	query = projectPermissionCheckV2(ctx, query, permissionCheckV2, queries)
	eq := sq.Eq{grantedProjectColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-T84X9", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		grantedProjects, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	return grantedProjects, nil
}

func NewGrantedProjectNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(GrantedProjectColumnName, value, method)
}

func NewGrantedProjectResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(grantedProjectColumnResourceOwner, value, TextEquals)
}

func NewGrantedProjectIDSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(GrantedProjectColumnID, list, ListIn)
}

func NewGrantedProjectOrganizationIDSearchQuery(value string) (SearchQuery, error) {
	project, err := NewGrantedProjectResourceOwnerSearchQuery(value)
	if err != nil {
		return nil, err
	}
	grant, err := NewGrantedProjectGrantedOrganizationIDSearchQuery(value)
	if err != nil {
		return nil, err
	}
	return NewOrQuery(project, grant)
}

func NewGrantedProjectGrantResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(grantedProjectColumnGrantResourceOwner, value, TextEquals)
}

func NewGrantedProjectGrantedOrganizationIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(grantedProjectColumnGrantedOrganization, value, TextEquals)
}

func NewProjectNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnName, value, method)
}

func NewProjectIDSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(ProjectColumnID, list, ListIn)
}

func NewProjectResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnResourceOwner, value, TextEquals)
}

func (r *ProjectSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *ProjectSearchQueries) AppendPermissionQueries(permissions []string) error {
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		query, err := NewProjectIDSearchQuery(ids)
		if err != nil {
			return err
		}
		r.Queries = append(r.Queries, query)
	}
	return nil
}

func (q *ProjectSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareProjectQuery() (sq.SelectBuilder, func(*sql.Row) (*Project, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier()).
			From(projectsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Project, error) {
			p := new(Project)
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.State,
				&p.Sequence,
				&p.Name,
				&p.ProjectRoleAssertion,
				&p.ProjectRoleCheck,
				&p.HasProjectCheck,
				&p.PrivateLabelingSetting,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-fk2fs", "Errors.Project.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-dj2FF", "Errors.Internal")
			}
			return p, nil
		}
}

func prepareProjectsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Projects, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier(),
			countColumn.identifier()).
			From(projectsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Projects, error) {
			projects := make([]*Project, 0)
			var count uint64
			for rows.Next() {
				project := new(Project)
				err := rows.Scan(
					&project.ID,
					&project.CreationDate,
					&project.ChangeDate,
					&project.ResourceOwner,
					&project.State,
					&project.Sequence,
					&project.Name,
					&project.ProjectRoleAssertion,
					&project.ProjectRoleCheck,
					&project.HasProjectCheck,
					&project.PrivateLabelingSetting,
					&count,
				)
				if err != nil {
					return nil, err
				}
				projects = append(projects, project)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-QMXJv", "Errors.Query.CloseRows")
			}

			return &Projects{
				Projects: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

type GrantedProjects struct {
	SearchResponse
	GrantedProjects []*GrantedProject
}

func grantedProjectsCheckPermission(ctx context.Context, grantedProjects *GrantedProjects, permissionCheck domain.PermissionCheck) {
	grantedProjects.GrantedProjects = slices.DeleteFunc(grantedProjects.GrantedProjects,
		func(grantedProject *GrantedProject) bool {
			if grantedProject.GrantedOrgID != "" {
				return projectGrantCheckPermission(ctx, grantedProject.ResourceOwner, grantedProject.ProjectID, grantedProject.GrantID, grantedProject.GrantedOrgID, permissionCheck) != nil
			}
			return projectCheckPermission(ctx, grantedProject.ResourceOwner, grantedProject.ProjectID, permissionCheck) != nil
		},
	)
}

type GrantedProject struct {
	ProjectID     string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	InstanceID    string
	ProjectState  domain.ProjectState
	ProjectName   string

	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting

	GrantID           string
	GrantedOrgID      string
	OrgName           string
	ProjectGrantState domain.ProjectGrantState
}

func prepareGrantedProjectsQuery() (sq.SelectBuilder, func(*sql.Rows) (*GrantedProjects, error)) {
	return sq.Select(
			GrantedProjectColumnID.identifier(),
			GrantedProjectColumnCreationDate.identifier(),
			GrantedProjectColumnChangeDate.identifier(),
			grantedProjectColumnResourceOwner.identifier(),
			grantedProjectColumnInstanceID.identifier(),
			grantedProjectColumnState.identifier(),
			GrantedProjectColumnName.identifier(),
			grantedProjectColumnProjectRoleAssertion.identifier(),
			grantedProjectColumnProjectRoleCheck.identifier(),
			grantedProjectColumnHasProjectCheck.identifier(),
			grantedProjectColumnPrivateLabelingSetting.identifier(),
			grantedProjectColumnGrantID.identifier(),
			grantedProjectColumnGrantedOrganization.identifier(),
			grantedProjectColumnGrantedOrganizationName.identifier(),
			grantedProjectColumnGrantState.identifier(),
			countColumn.identifier(),
		).From(getProjectsAndGrantedProjectsFromQuery()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GrantedProjects, error) {
			projects := make([]*GrantedProject, 0)
			var (
				count             uint64
				grantID           = sql.NullString{}
				orgID             = sql.NullString{}
				orgName           = sql.NullString{}
				projectGrantState = sql.NullInt16{}
			)
			for rows.Next() {
				grantedProject := new(GrantedProject)
				err := rows.Scan(
					&grantedProject.ProjectID,
					&grantedProject.CreationDate,
					&grantedProject.ChangeDate,
					&grantedProject.ResourceOwner,
					&grantedProject.InstanceID,
					&grantedProject.ProjectState,
					&grantedProject.ProjectName,
					&grantedProject.ProjectRoleAssertion,
					&grantedProject.ProjectRoleCheck,
					&grantedProject.HasProjectCheck,
					&grantedProject.PrivateLabelingSetting,
					&grantID,
					&orgID,
					&orgName,
					&projectGrantState,
					&count,
				)
				if err != nil {
					return nil, err
				}
				if grantID.Valid {
					grantedProject.GrantID = grantID.String
				}
				if orgID.Valid {
					grantedProject.GrantedOrgID = orgID.String
				}
				if orgName.Valid {
					grantedProject.OrgName = orgName.String
				}
				if projectGrantState.Valid {
					grantedProject.ProjectGrantState = domain.ProjectGrantState(projectGrantState.Int16)
				}
				projects = append(projects, grantedProject)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-K9gEE", "Errors.Query.CloseRows")
			}

			return &GrantedProjects{
				GrantedProjects: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func getProjectsAndGrantedProjectsFromQuery() string {
	return "(" +
		prepareProjects() +
		" UNION ALL " +
		prepareGrantedProjects() +
		") AS " + grantedProjectsAlias.identifier()
}

func prepareProjects() string {
	builder := sq.Select(
		ProjectColumnID.identifier()+" AS "+GrantedProjectColumnID.name,
		ProjectColumnCreationDate.identifier()+" AS "+GrantedProjectColumnCreationDate.name,
		ProjectColumnChangeDate.identifier()+" AS "+GrantedProjectColumnChangeDate.name,
		ProjectColumnResourceOwner.identifier()+" AS "+grantedProjectColumnResourceOwner.name,
		ProjectColumnInstanceID.identifier()+" AS "+grantedProjectColumnInstanceID.name,
		ProjectColumnState.identifier()+" AS "+grantedProjectColumnState.name,
		ProjectColumnName.identifier()+" AS "+GrantedProjectColumnName.name,
		ProjectColumnProjectRoleAssertion.identifier()+" AS "+grantedProjectColumnProjectRoleAssertion.name,
		ProjectColumnProjectRoleCheck.identifier()+" AS "+grantedProjectColumnProjectRoleCheck.name,
		ProjectColumnHasProjectCheck.identifier()+" AS "+grantedProjectColumnHasProjectCheck.name,
		ProjectColumnPrivateLabelingSetting.identifier()+" AS "+grantedProjectColumnPrivateLabelingSetting.name,
		"NULL::TEXT AS "+grantedProjectColumnGrantResourceOwner.name,
		"NULL::TEXT AS "+grantedProjectColumnGrantID.name,
		"NULL::TEXT AS "+grantedProjectColumnGrantedOrganization.name,
		"NULL::TEXT AS "+grantedProjectColumnGrantedOrganizationName.name,
		"NULL::SMALLINT AS "+grantedProjectColumnGrantState.name,
		countColumn.identifier()).
		From(projectsTable.identifier()).
		PlaceholderFormat(sq.Dollar)

	stmt, _ := builder.MustSql()
	return stmt
}

func prepareGrantedProjects() string {
	grantedOrgTable := orgsTable.setAlias(ProjectGrantGrantedOrgTableAlias)
	grantedOrgIDColumn := OrgColumnID.setTable(grantedOrgTable)
	builder := sq.Select(
		ProjectGrantColumnProjectID.identifier()+" AS "+GrantedProjectColumnID.name,
		ProjectGrantColumnCreationDate.identifier()+" AS "+GrantedProjectColumnCreationDate.name,
		ProjectGrantColumnChangeDate.identifier()+" AS "+GrantedProjectColumnChangeDate.name,
		ProjectColumnResourceOwner.identifier()+" AS "+grantedProjectColumnResourceOwner.name,
		ProjectGrantColumnInstanceID.identifier()+" AS "+grantedProjectColumnInstanceID.name,
		ProjectColumnState.identifier()+" AS "+grantedProjectColumnState.name,
		ProjectColumnName.identifier()+" AS "+GrantedProjectColumnName.name,
		ProjectColumnProjectRoleAssertion.identifier()+" AS "+grantedProjectColumnProjectRoleAssertion.name,
		ProjectColumnProjectRoleCheck.identifier()+" AS "+grantedProjectColumnProjectRoleCheck.name,
		ProjectColumnHasProjectCheck.identifier()+" AS "+grantedProjectColumnHasProjectCheck.name,
		ProjectColumnPrivateLabelingSetting.identifier()+" AS "+grantedProjectColumnPrivateLabelingSetting.name,
		ProjectGrantColumnResourceOwner.identifier()+" AS "+grantedProjectColumnGrantResourceOwner.name,
		ProjectGrantColumnGrantID.identifier()+" AS "+grantedProjectColumnGrantID.name,
		ProjectGrantColumnGrantedOrgID.identifier()+" AS "+grantedProjectColumnGrantedOrganization.name,
		ProjectGrantColumnGrantedOrgName.identifier()+" AS "+grantedProjectColumnGrantedOrganizationName.name,
		ProjectGrantColumnState.identifier()+" AS "+grantedProjectColumnGrantState.name,
		countColumn.identifier()).
		From(projectGrantsTable.identifier()).
		PlaceholderFormat(sq.Dollar).
		LeftJoin(join(ProjectColumnID, ProjectGrantColumnProjectID)).
		LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID))

	stmt, _ := builder.MustSql()
	return stmt
}
