package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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
	ProjectGrantColumnOwnerRemoved = Column{
		name:  projection.ProjectGrantColumnOwnerRemoved,
		table: projectGrantsTable,
	}
	ProjectGrantColumnGrantGrantedOrgRemoved = Column{
		name:  projection.ProjectGrantColumnGrantedOrgRemoved,
		table: projectGrantsTable,
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
	GrantedRoleKeys   database.StringArray
	ResourceOwnerName string
}

type ProjectGrantSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) ProjectGrantByID(ctx context.Context, shouldTriggerBulk bool, id string, withOwnerRemoved bool) (_ *ProjectGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.ProjectGrantProjection.Trigger(ctx)
	}

	stmt, scan := prepareProjectGrantQuery(ctx, q.client)
	eq := sq.Eq{
		ProjectGrantColumnGrantID.identifier():    id,
		ProjectGrantColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[ProjectGrantColumnOwnerRemoved.identifier()] = false
		eq[ProjectGrantColumnGrantGrantedOrgRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Nf93d", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ProjectGrantByIDAndGrantedOrg(ctx context.Context, id, grantedOrg string, withOwnerRemoved bool) (_ *ProjectGrant, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareProjectGrantQuery(ctx, q.client)
	eq := sq.Eq{
		ProjectGrantColumnGrantID.identifier():      id,
		ProjectGrantColumnGrantedOrgID.identifier(): grantedOrg,
		ProjectGrantColumnInstanceID.identifier():   authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[ProjectGrantColumnOwnerRemoved.identifier()] = false
		eq[ProjectGrantColumnGrantGrantedOrgRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-MO9fs", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) SearchProjectGrants(ctx context.Context, queries *ProjectGrantSearchQueries, withOwnerRemoved bool) (projects *ProjectGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareProjectGrantsQuery(ctx, q.client)
	eq := sq.Eq{
		ProjectGrantColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[ProjectGrantColumnOwnerRemoved.identifier()] = false
		eq[ProjectGrantColumnGrantGrantedOrgRemoved.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-N9fsg", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-PP02n", "Errors.Internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projectGrantsTable)
	return projects, err
}

func (q *Queries) SearchProjectGrantsByProjectIDAndRoleKey(ctx context.Context, projectID, roleKey string, withOwnerRemoved bool) (projects *ProjectGrants, err error) {
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
	return q.SearchProjectGrants(ctx, searchQuery, withOwnerRemoved)
}

func NewProjectGrantProjectIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnProjectID, value, TextEquals)
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

func (q *ProjectGrantSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareProjectGrantQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*ProjectGrant, error)) {
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
			LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID) + db.Timetravel(call.Took(ctx))),
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
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-n98GGs", "Errors.ProjectGrant.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-w9fsH", "Errors.Internal")
			}

			grant.ProjectName = projectName.String
			grant.ResourceOwnerName = resourceOwnerName.String
			grant.OrgName = orgName.String

			return grant, nil
		}
}

func prepareProjectGrantsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*ProjectGrants, error)) {
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
			LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID) + db.Timetravel(call.Took(ctx))),
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
				return nil, errors.ThrowInternal(err, "QUERY-K9gEE", "Errors.Query.CloseRows")
			}

			return &ProjectGrants{
				ProjectGrants: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
