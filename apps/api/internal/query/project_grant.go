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
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectGrantGrantedOrgTableAlias    = "o"
	ProjectGrantResourceOwnerTableAlias = "r"
)

var (
	projectGrantsTable = table{
		name:          projection.ProjectGrantProjectionTable,
		instanceIDCol: projection.ProjectGrantColumnInstanceID,
	}
	ProjectGrantColumnCreationDate = Column{
		name:  projection.ProjectGrantColumnCreationDate,
		table: projectGrantsTable,
	}
	ProjectGrantColumnChangeDate = Column{
		name:  projection.ProjectGrantColumnChangeDate,
		table: projectGrantsTable,
	}
	ProjectGrantColumnResourceOwner = Column{
		name:  projection.ProjectGrantColumnResourceOwner,
		table: projectGrantsTable,
	}
	ProjectGrantColumnInstanceID = Column{
		name:  projection.ProjectGrantColumnInstanceID,
		table: projectGrantsTable,
	}
	ProjectGrantColumnState = Column{
		name:  projection.ProjectGrantColumnState,
		table: projectGrantsTable,
	}
	ProjectGrantColumnSequence = Column{
		name:  projection.ProjectGrantColumnSequence,
		table: projectGrantsTable,
	}
	ProjectGrantColumnProjectID = Column{
		name:  projection.ProjectGrantColumnProjectID,
		table: projectGrantsTable,
	}
	ProjectGrantColumnGrantedOrgID = Column{
		name:  projection.ProjectGrantColumnGrantedOrgID,
		table: projectGrantsTable,
	}
	ProjectGrantColumnGrantID = Column{
		name:  projection.ProjectGrantColumnGrantID,
		table: projectGrantsTable,
	}
	ProjectGrantColumnGrantedRoleKeys = Column{
		name:  projection.ProjectGrantColumnRoleKeys,
		table: projectGrantsTable,
	}
	ProjectGrantColumnGrantedOrgName = Column{
		name:  projection.OrgColumnName,
		table: orgsTable.setAlias(ProjectGrantGrantedOrgTableAlias),
	}
	ProjectGrantColumnResourceOwnerName = Column{
		name:  projection.OrgColumnName,
		table: orgsTable.setAlias(ProjectGrantResourceOwnerTableAlias),
	}
)

type ProjectGrants struct {
	SearchResponse
	ProjectGrants []*ProjectGrant
}

type ProjectGrant struct {
	ProjectID     string
	GrantID       string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.ProjectGrantState
	Sequence      uint64

	ProjectName       string
	GrantedOrgID      string
	OrgName           string
	GrantedRoleKeys   database.TextArray[string]
	ResourceOwnerName string
}

type ProjectGrantSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func projectGrantsCheckPermission(ctx context.Context, projectGrants *ProjectGrants, permissionCheck domain.PermissionCheck) {
	projectGrants.ProjectGrants = slices.DeleteFunc(projectGrants.ProjectGrants,
		func(projectGrant *ProjectGrant) bool {
			return projectGrantCheckPermission(ctx, projectGrant.ResourceOwner, projectGrant.ProjectID, projectGrant.GrantID, projectGrant.GrantedOrgID, permissionCheck) != nil
		},
	)
}

func projectGrantCheckPermission(ctx context.Context, resourceOwner, projectID, grantID, grantedOrgID string, permissionCheck domain.PermissionCheck) error {
	if err := permissionCheck(ctx, domain.PermissionProjectGrantRead, resourceOwner, grantID); err != nil {
		if err := permissionCheck(ctx, domain.PermissionProjectGrantRead, grantedOrgID, grantID); err != nil {
			if err := permissionCheck(ctx, domain.PermissionProjectGrantRead, resourceOwner, projectID); err != nil {
				return err
			}
		}
	}
	return nil
}

// TODO: add permission check on project grant level
func projectGrantPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *ProjectGrantSearchQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		ProjectGrantColumnResourceOwner,
		domain.PermissionProjectGrantRead,
		SingleOrgPermissionOption(queries.Queries),
	)
	return query.JoinClause(join, args...)
}

func (q *Queries) GetProjectGrantByIDWithPermission(ctx context.Context, shouldTriggerBulk bool, id string, permissionCheck domain.PermissionCheck) (*ProjectGrant, error) {
	projectGrant, err := q.ProjectGrantByID(ctx, shouldTriggerBulk, id)
	if err != nil {
		return nil, err
	}
	if err := projectCheckPermission(ctx, projectGrant.ResourceOwner, projectGrant.GrantID, permissionCheck); err != nil {
		return nil, err
	}
	return projectGrant, nil
}

func (q *Queries) ProjectGrantByID(ctx context.Context, shouldTriggerBulk bool, id string) (grant *ProjectGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerProjectGrantProjection")
		ctx, err = projection.ProjectGrantProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareProjectGrantQuery()
	eq := sq.Eq{
		ProjectGrantColumnGrantID.identifier():    id,
		ProjectGrantColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Nf93d", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		grant, err = scan(row)
		return err
	}, query, args...)
	return grant, err
}

func (q *Queries) ProjectGrantByIDAndGrantedOrg(ctx context.Context, id, grantedOrg string) (grant *ProjectGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectGrantQuery()
	eq := sq.Eq{
		ProjectGrantColumnGrantID.identifier():      id,
		ProjectGrantColumnGrantedOrgID.identifier(): grantedOrg,
		ProjectGrantColumnInstanceID.identifier():   authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-MO9fs", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		grant, err = scan(row)
		return err
	}, query, args...)
	return grant, err
}

func (q *Queries) SearchProjectGrants(ctx context.Context, queries *ProjectGrantSearchQueries, permissionCheck domain.PermissionCheck) (grants *ProjectGrants, err error) {
	// removed as permission v2 is not implemented yet for project grant level permissions
	// permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	projectsGrants, err := q.searchProjectGrants(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil { // && !authz.GetFeatures(ctx).PermissionCheckV2 {
		projectGrantsCheckPermission(ctx, projectsGrants, permissionCheck)
	}
	return projectsGrants, nil
}

func (q *Queries) searchProjectGrants(ctx context.Context, queries *ProjectGrantSearchQueries, permissionCheckV2 bool) (grants *ProjectGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareProjectGrantsQuery()
	query = projectGrantPermissionCheckV2(ctx, query, permissionCheckV2, queries)
	eq := sq.Eq{
		ProjectGrantColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-N9fsg", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		grants, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-PP02n", "Errors.Internal")
	}

	grants.State, err = q.latestState(ctx, projectGrantsTable)
	return grants, err
}

// SearchProjectGrantsByProjectIDAndRoleKey is used internally to remove the roles of a project grant, so no permission check necessary
func (q *Queries) SearchProjectGrantsByProjectIDAndRoleKey(ctx context.Context, projectID, roleKey string) (projects *ProjectGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	searchQuery := &ProjectGrantSearchQueries{
		SearchRequest: SearchRequest{},
		Queries:       make([]SearchQuery, 2),
	}
	searchQuery.Queries[0], err = NewProjectGrantProjectIDSearchQuery(projectID)
	if err != nil {
		return nil, err
	}
	searchQuery.Queries[1], err = NewProjectGrantRoleKeySearchQuery(roleKey)
	if err != nil {
		return nil, err
	}
	return q.SearchProjectGrants(ctx, searchQuery, nil)
}

func NewProjectGrantProjectIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnProjectID, value, TextEquals)
}

func (q *ProjectGrantSearchQueries) AppendPermissionQueries(permissions []string) error {
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		query, err := NewProjectGrantIDsSearchQuery(ids)
		if err != nil {
			return err
		}
		q.Queries = append(q.Queries, query)
	}
	return nil
}

func NewProjectGrantProjectIDsSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(ProjectGrantColumnProjectID, list, ListIn)
}

func NewProjectGrantIDsSearchQuery(values []string) (SearchQuery, error) {
	list := make([]interface{}, len(values))
	for i, value := range values {
		list[i] = value
	}
	return NewListQuery(ProjectGrantColumnGrantID, list, ListIn)
}
func NewProjectGrantProjectNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnName, value, method)
}

func NewProjectGrantRoleKeySearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnGrantedRoleKeys, value, TextListContains)
}

func NewProjectGrantResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnResourceOwner, value, TextEquals)
}

func NewProjectGrantGrantedOrgIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnGrantedOrgID, value, TextEquals)
}

func (q *ProjectGrantSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectGrantResourceOwnerSearchQuery(orgID)
	if err != nil {
		return err
	}
	q.Queries = append(q.Queries, query)
	return nil
}

func (q *ProjectGrantSearchQueries) AppendGrantedOrgQuery(orgID string) error {
	query, err := NewProjectGrantGrantedOrgIDSearchQuery(orgID)
	if err != nil {
		return err
	}
	q.Queries = append(q.Queries, query)
	return nil
}

func (q *ProjectGrantSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareProjectGrantQuery() (sq.SelectBuilder, func(*sql.Row) (*ProjectGrant, error)) {
	resourceOwnerOrgTable := orgsTable.setAlias(ProjectGrantResourceOwnerTableAlias)
	resourceOwnerIDColumn := OrgColumnID.setTable(resourceOwnerOrgTable)
	grantedOrgTable := orgsTable.setAlias(ProjectGrantGrantedOrgTableAlias)
	grantedOrgIDColumn := OrgColumnID.setTable(grantedOrgTable)
	return sq.Select(
			ProjectGrantColumnProjectID.identifier(),
			ProjectGrantColumnGrantID.identifier(),
			ProjectGrantColumnCreationDate.identifier(),
			ProjectGrantColumnChangeDate.identifier(),
			ProjectGrantColumnResourceOwner.identifier(),
			ProjectGrantColumnState.identifier(),
			ProjectGrantColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectGrantColumnGrantedOrgID.identifier(),
			ProjectGrantColumnGrantedOrgName.identifier(),
			ProjectGrantColumnGrantedRoleKeys.identifier(),
			ProjectGrantColumnResourceOwnerName.identifier()).
			From(projectGrantsTable.identifier()).
			PlaceholderFormat(sq.Dollar).
			LeftJoin(join(ProjectColumnID, ProjectGrantColumnProjectID)).
			LeftJoin(join(resourceOwnerIDColumn, ProjectGrantColumnResourceOwner)).
			LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID)),
		func(row *sql.Row) (*ProjectGrant, error) {
			grant := new(ProjectGrant)
			var (
				projectName       sql.NullString
				orgName           sql.NullString
				resourceOwnerName sql.NullString
			)
			err := row.Scan(
				&grant.ProjectID,
				&grant.GrantID,
				&grant.CreationDate,
				&grant.ChangeDate,
				&grant.ResourceOwner,
				&grant.State,
				&grant.Sequence,
				&projectName,
				&grant.GrantedOrgID,
				&orgName,
				&grant.GrantedRoleKeys,
				&resourceOwnerName,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-n98GGs", "Errors.ProjectGrant.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-w9fsH", "Errors.Internal")
			}

			grant.ProjectName = projectName.String
			grant.ResourceOwnerName = resourceOwnerName.String
			grant.OrgName = orgName.String

			return grant, nil
		}
}

func prepareProjectGrantsQuery() (sq.SelectBuilder, func(*sql.Rows) (*ProjectGrants, error)) {
	resourceOwnerOrgTable := orgsTable.setAlias(ProjectGrantResourceOwnerTableAlias)
	resourceOwnerIDColumn := OrgColumnID.setTable(resourceOwnerOrgTable)
	grantedOrgTable := orgsTable.setAlias(ProjectGrantGrantedOrgTableAlias)
	grantedOrgIDColumn := OrgColumnID.setTable(grantedOrgTable)
	return sq.Select(
			ProjectGrantColumnProjectID.identifier(),
			ProjectGrantColumnGrantID.identifier(),
			ProjectGrantColumnCreationDate.identifier(),
			ProjectGrantColumnChangeDate.identifier(),
			ProjectGrantColumnResourceOwner.identifier(),
			ProjectGrantColumnState.identifier(),
			ProjectGrantColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectGrantColumnGrantedOrgID.identifier(),
			ProjectGrantColumnGrantedOrgName.identifier(),
			ProjectGrantColumnGrantedRoleKeys.identifier(),
			ProjectGrantColumnResourceOwnerName.identifier(),
			countColumn.identifier()).
			From(projectGrantsTable.identifier()).
			PlaceholderFormat(sq.Dollar).
			LeftJoin(join(ProjectColumnID, ProjectGrantColumnProjectID)).
			LeftJoin(join(resourceOwnerIDColumn, ProjectGrantColumnResourceOwner)).
			LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID)),
		func(rows *sql.Rows) (*ProjectGrants, error) {
			projects := make([]*ProjectGrant, 0)
			var (
				count             uint64
				projectName       sql.NullString
				orgName           sql.NullString
				resourceOwnerName sql.NullString
			)
			for rows.Next() {
				grant := new(ProjectGrant)
				err := rows.Scan(
					&grant.ProjectID,
					&grant.GrantID,
					&grant.CreationDate,
					&grant.ChangeDate,
					&grant.ResourceOwner,
					&grant.State,
					&grant.Sequence,
					&projectName,
					&grant.GrantedOrgID,
					&orgName,
					&grant.GrantedRoleKeys,
					&resourceOwnerName,
					&count,
				)
				if err != nil {
					return nil, err
				}

				grant.ProjectName = projectName.String
				grant.ResourceOwnerName = resourceOwnerName.String
				grant.OrgName = orgName.String

				projects = append(projects, grant)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-K9gEE", "Errors.Query.CloseRows")
			}

			return &ProjectGrants{
				ProjectGrants: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
