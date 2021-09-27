package query

import (
	"context"
	"database/sql"
	errs "errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	projectRolesTable = table{
		name: projection.ProjectRoleProjectionTable,
	}
	ProjectRoleColumnCreationDate = Column{
		name:  projection.ProjectRoleColumnCreationDate,
		table: projectRolesTable,
	}
	ProjectRoleColumnChangeDate = Column{
		name:  projection.ProjectRoleColumnChangeDate,
		table: projectRolesTable,
	}
	ProjectRoleColumnResourceOwner = Column{
		name:  projection.ProjectRoleColumnResourceOwner,
		table: projectRolesTable,
	}
	ProjectRoleColumnSequence = Column{
		name:  projection.ProjectRoleColumnSequence,
		table: projectRolesTable,
	}
	ProjectRoleColumnProjectID = Column{
		name:  projection.ProjectRoleColumnProjectID,
		table: projectRolesTable,
	}
	ProjectRoleColumnKey = Column{
		name:  projection.ProjectRoleColumnKey,
		table: projectRolesTable,
	}
	ProjectRoleColumnDisplayName = Column{
		name:  projection.ProjectRoleColumnDisplayName,
		table: projectRolesTable,
	}
	ProjectRoleColumnGroupName = Column{
		name:  projection.ProjectRoleColumnGroupName,
		table: projectRolesTable,
	}
	ProjectRoleColumnCreator = Column{
		name:  projection.ProjectRoleColumnCreator,
		table: projectRolesTable,
	}
)

func prepareProjectRoleQuery() (sq.SelectBuilder, func(*sql.Row) (*ProjectRole, error)) {
	return sq.Select(
			ProjectRoleColumnProjectID.identifier(),
			ProjectRoleColumnCreationDate.identifier(),
			ProjectRoleColumnChangeDate.identifier(),
			ProjectRoleColumnResourceOwner.identifier(),
			ProjectRoleColumnSequence.identifier(),
			ProjectRoleColumnKey.identifier(),
			ProjectRoleColumnDisplayName.identifier(),
			ProjectRoleColumnGroupName.identifier()).
			From(projection.ProjectRoleProjectionTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*ProjectRole, error) {
			p := new(ProjectRole)
			err := row.Scan(
				&p.ProjectID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.Sequence,
				&p.Key,
				&p.DisplayName,
				&p.Group,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-Mf0wf", "errors.project_roles.not_found")
				}
				fmt.Printf("error: %v", err.Error())
				return nil, errors.ThrowInternal(err, "QUERY-M00sf", "errors.internal")
			}
			return p, nil
		}
}

func (q *Queries) prepareProjectRolesQuery() (sq.SelectBuilder, func(*sql.Rows) (*ProjectRoles, error)) {
	return sq.Select(
			ProjectRoleColumnProjectID.identifier(),
			ProjectRoleColumnCreationDate.identifier(),
			ProjectRoleColumnChangeDate.identifier(),
			ProjectRoleColumnResourceOwner.identifier(),
			ProjectRoleColumnSequence.identifier(),
			ProjectRoleColumnKey.identifier(),
			ProjectRoleColumnDisplayName.identifier(),
			ProjectRoleColumnGroupName.identifier(),
			"COUNT(*) OVER ()").
			From(projection.ProjectRoleProjectionTable).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*ProjectRoles, error) {
			projects := make([]*ProjectRole, 0)
			var count uint64
			for rows.Next() {
				project := new(ProjectRole)
				err := rows.Scan(
					&project.ProjectID,
					&project.CreationDate,
					&project.ChangeDate,
					&project.ResourceOwner,
					&project.Sequence,
					&project.Key,
					&project.DisplayName,
					&project.Group,
					&count,
				)
				if err != nil {
					return nil, err
				}
				projects = append(projects, project)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-ML0Fs", "unable to close rows")
			}

			return &ProjectRoles{
				ProjectRoles: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *Queries) prepareProjectRoleUniqueQuery() (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select("COUNT(*) = 0").
			From(projection.ProjectRoleProjectionTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (isUnique bool, err error) {
			err = row.Scan(&isUnique)
			if err != nil {
				return false, errors.ThrowInternal(err, "QUERY-2N9fs", "errors.internal")
			}
			return isUnique, err
		}
}

func (q *Queries) ProjectRoleByID(ctx context.Context, projectID, key string) (*ProjectRole, error) {
	stmt, scan := prepareProjectRoleQuery()
	query, args, err := stmt.
		Where(sq.Eq{ProjectRoleColumnProjectID.identifier(): projectID}).
		Where(sq.Eq{ProjectRoleColumnKey.identifier(): key}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-2N0fs", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ExistsProjectRole(ctx context.Context, projectID, key string) (err error) {
	_, err = q.ProjectRoleByID(ctx, projectID, key)
	return err
}

func (q *Queries) SearchProjectRoles(ctx context.Context, queries *ProjectRoleSearchQueries) (projects *ProjectRoles, err error) {
	query, scan := q.prepareProjectRolesQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-3N9ff", "Errors.project_roless.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5Ngd9", "Errors.project_roles.internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projectRolesTable)
	return projects, err
}

func (q *Queries) SearchGrantedProjectRoles(ctx context.Context, grantID, grantedOrg string, queries *ProjectRoleSearchQueries) (projects *ProjectRoles, err error) {
	grant, err := q.ProjectGrantByIDAndGrantedOrg(ctx, grantID, grantedOrg)
	if err != nil {
		return nil, err
	}
	err = queries.AppendRoleKeysQuery(grant.GrantedRoleKeys)
	if err != nil {
		return nil, err
	}
	query, scan := q.prepareProjectRolesQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-3N9ff", "Errors.project_roless.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5Ngd9", "Errors.project_roles.internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projectRolesTable)
	return projects, err
}

type ProjectRoles struct {
	SearchResponse
	ProjectRoles []*ProjectRole
}

type ProjectRole struct {
	ProjectID     string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	Key         string
	DisplayName string
	Group       string
}

type ProjectRoleSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewProjectRoleProjectIDSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnProjectID, value, method)
}

func NewProjectRoleResourceOwnerSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnResourceOwner, value, method)
}

func NewProjectRoleKeySearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnKey, value, method)
}

func NewProjectRoleKeysSearchQuery(values []string) (SearchQuery, error) {
	list := make([]interface{}, len(values))
	for i, value := range values {
		list[i] = value
	}
	return NewListQuery(ProjectRoleColumnKey, list, ListIn)
}

func NewProjectRoleDisplayNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnDisplayName, value, method)
}

func NewProjectRoleGroupSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnGroupName, value, method)
}

func (q *ProjectRoleSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func (r *ProjectRoleSearchQueries) AppendProjectIDQuery(projectID string) error {
	query, err := NewProjectRoleProjectIDSearchQuery(TextEquals, projectID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *ProjectRoleSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectRoleResourceOwnerSearchQuery(TextEquals, orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *ProjectRoleSearchQueries) AppendRoleKeysQuery(keys []string) error {
	query, err := NewProjectRoleKeysSearchQuery(keys)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}
