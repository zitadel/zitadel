package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

const (
	ProjectGrantGrantedOrgTableAlias    = "o"
	ProjectGrantResourceOwnerTableAlias = "r"
)

var (
	projectGrantsTable = table{
		name: projection.ProjectGrantProjectionTable,
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
	ProjectGrantColumnCreator = Column{
		name:  projection.ProjectGrantColumnCreator,
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
	GrantedRoleKeys   pq.StringArray
	ResourceOwnerName string
}

type ProjectGrantSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) ProjectGrantByID(ctx context.Context, id string) (*ProjectGrant, error) {
	stmt, scan := prepareProjectGrantQuery()
	query, args, err := stmt.Where(sq.Eq{
		ProjectGrantColumnGrantID.identifier(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Nf93d", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ProjectGrantByIDAndGrantedOrg(ctx context.Context, id, grantedOrg string) (*ProjectGrant, error) {
	stmt, scan := prepareProjectGrantQuery()
	query, args, err := stmt.Where(sq.Eq{
		ProjectGrantColumnGrantID.identifier():      id,
		ProjectGrantColumnGrantedOrgID.identifier(): grantedOrg,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-MO9fs", "Errors.Query.SQLStatment")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ExistsProjectGrant(ctx context.Context, id string) (err error) {
	_, err = q.ProjectGrantByID(ctx, id)
	return err
}

func (q *Queries) SearchProjectGrants(ctx context.Context, queries *ProjectGrantSearchQueries) (projects *ProjectGrants, err error) {
	query, scan := prepareProjectGrantsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
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

func (q *Queries) SearchProjectGrantsByProjectIDAndRoleKey(ctx context.Context, projectID, roleKey string) (projects *ProjectGrants, err error) {
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
	return q.SearchProjectGrants(ctx, searchQuery)
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
		query = q.ToQuery(query)
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
			From(projectGrantsTable.identifier()).PlaceholderFormat(sq.Dollar).
			LeftJoin(join(ProjectColumnID, ProjectGrantColumnProjectID)).
			LeftJoin(join(resourceOwnerIDColumn, ProjectGrantColumnResourceOwner)).
			LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID)),
		func(row *sql.Row) (*ProjectGrant, error) {
			p := new(ProjectGrant)
			err := row.Scan(
				&p.ProjectID,
				&p.GrantID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.State,
				&p.Sequence,
				&p.ProjectName,
				&p.GrantedOrgID,
				&p.OrgName,
				&p.GrantedRoleKeys,
				&p.ResourceOwnerName,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-n98GGs", "Errors.ProjectGrant.NotFound")
				}
				return nil, errors.ThrowInternal(err, "QUERY-w9fsH", "Errors.Internal")
			}
			return p, nil
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
			From(projectGrantsTable.identifier()).PlaceholderFormat(sq.Dollar).
			LeftJoin(join(ProjectColumnID, ProjectGrantColumnProjectID)).
			LeftJoin(join(resourceOwnerIDColumn, ProjectGrantColumnResourceOwner)).
			LeftJoin(join(grantedOrgIDColumn, ProjectGrantColumnGrantedOrgID)),
		func(rows *sql.Rows) (*ProjectGrants, error) {
			projects := make([]*ProjectGrant, 0)
			var count uint64
			for rows.Next() {
				project := new(ProjectGrant)
				err := rows.Scan(
					&project.ProjectID,
					&project.GrantID,
					&project.CreationDate,
					&project.ChangeDate,
					&project.ResourceOwner,
					&project.State,
					&project.Sequence,
					&project.ProjectName,
					&project.GrantedOrgID,
					&project.OrgName,
					&project.GrantedRoleKeys,
					&project.ResourceOwnerName,
					&count,
				)
				if err != nil {
					return nil, err
				}
				projects = append(projects, project)
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
