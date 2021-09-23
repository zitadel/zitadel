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

func prepareProjectRoleQuery() (sq.SelectBuilder, func(*sql.Row) (*ProjectRole, error)) {
	return sq.Select(
			ProjectRoleColumnProjectID.toFullColumnName(),
			ProjectRoleColumnCreationDate.toFullColumnName(),
			ProjectRoleColumnChangeDate.toFullColumnName(),
			ProjectRoleColumnResourceOwner.toFullColumnName(),
			ProjectRoleColumnSequence.toFullColumnName(),
			ProjectRoleColumnKey.toFullColumnName(),
			ProjectRoleColumnDisplayName.toFullColumnName(),
			ProjectRoleColumnGroupName.toFullColumnName()).
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
			ProjectRoleColumnProjectID.toFullColumnName(),
			ProjectRoleColumnCreationDate.toFullColumnName(),
			ProjectRoleColumnChangeDate.toFullColumnName(),
			ProjectRoleColumnResourceOwner.toFullColumnName(),
			ProjectRoleColumnSequence.toFullColumnName(),
			ProjectRoleColumnKey.toFullColumnName(),
			ProjectRoleColumnDisplayName.toFullColumnName(),
			ProjectRoleColumnGroupName.toFullColumnName(),
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
		Where(sq.Eq{ProjectRoleColumnProjectID.toColumnName(): projectID}).
		Where(sq.Eq{ProjectRoleColumnKey.toColumnName(): key}).ToSql()
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
	projects.LatestSequence, err = q.latestSequence(ctx, projection.ProjectRoleProjectionTable)
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

func (q *ProjectRoleSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func (r *ProjectRoleSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectRoleResourceOwnerSearchQuery(TextEquals, orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

type ProjectRoleColumn int32

const (
	ProjectRoleColumnCreationDate ProjectRoleColumn = iota + 1
	ProjectRoleColumnChangeDate
	ProjectRoleColumnResourceOwner
	ProjectRoleColumnSequence
	ProjectRoleColumnProjectID
	ProjectRoleColumnKey
	ProjectRoleColumnDisplayName
	ProjectRoleColumnGroupName
	ProjectRoleColumnCreatorName
)

func (c ProjectRoleColumn) toColumnName() string {
	switch c {
	case ProjectRoleColumnProjectID:
		return projection.ProjectRoleProjectIDCol
	case ProjectRoleColumnCreationDate:
		return projection.ProjectRoleCreationDateCol
	case ProjectRoleColumnChangeDate:
		return projection.ProjectRoleChangeDateCol
	case ProjectRoleColumnResourceOwner:
		return projection.ProjectRoleResourceOwnerCol
	case ProjectRoleColumnSequence:
		return projection.ProjectRoleSequenceCol
	case ProjectRoleColumnKey:
		return projection.ProjectRoleKeyCol
	case ProjectRoleColumnDisplayName:
		return projection.ProjectRoleDisplayNameCol
	case ProjectRoleColumnGroupName:
		return projection.ProjectRoleGroupNameCol
	case ProjectRoleColumnCreatorName:
		return projection.ProjectRoleCreatorCol
	default:
		return ""
	}
}

func (c ProjectRoleColumn) toFullColumnName() string {
	return c.toColumnName()
}
