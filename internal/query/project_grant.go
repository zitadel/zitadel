package query

import (
	"context"
	"database/sql"
	errs "errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

func prepareProjectGrantQuery() (sq.SelectBuilder, func(*sql.Row) (*ProjectGrant, error)) {
	joins := []JoinData{
		{
			projection.ProjectProjectionTable, ProjectGrantColumnProjectID.toColumnName(), ProjectColumnID.toColumnName(),
		},
	}
	return sq.Select(
			ProjectGrantColumnProjectID.toColumnName(),
			ProjectGrantColumnGrantID.toColumnName(),
			ProjectGrantColumnCreationDate.toColumnName(),
			ProjectGrantColumnChangeDate.toColumnName(),
			ProjectGrantColumnResourceOwner.toColumnName(),
			ProjectGrantColumnState.toColumnName(),
			ProjectGrantColumnSequence.toColumnName(),
			ProjectColumnName.toColumnName(),
			ProjectGrantColumnOrgID.toColumnName(),
			"granted_org_name",
			ProjectGrantColumnGrantedRoleKeys.toColumnName(),
			"resource_owner_name").
			From(projection.ProjectGrantProjectionTable).PlaceholderFormat(sq.Dollar).
			LeftJoin(GenerateJoinQuery(joins)),
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
					return nil, errors.ThrowNotFound(err, "QUERY-n98GGs", "errors.project_grants.not_found")
				}
				fmt.Printf("error: ", err.Error())
				return nil, errors.ThrowInternal(err, "QUERY-M00sf", "errors.internal")
			}
			return p, nil
		}
}

func (q *Queries) prepareProjectGrantsQuery() (sq.SelectBuilder, func(*sql.Rows) (*ProjectGrants, error)) {
	joins := []JoinData{
		{
			projection.ProjectProjectionTable, ProjectGrantColumnProjectID.toColumnName(), ProjectColumnID.toColumnName(),
		},
	}
	return sq.Select(
			ProjectGrantColumnProjectID.toColumnName(),
			ProjectGrantColumnGrantID.toColumnName(),
			ProjectGrantColumnCreationDate.toColumnName(),
			ProjectGrantColumnChangeDate.toColumnName(),
			ProjectGrantColumnResourceOwner.toColumnName(),
			ProjectGrantColumnState.toColumnName(),
			ProjectGrantColumnSequence.toColumnName(),
			ProjectColumnName.toColumnName(),
			ProjectGrantColumnOrgID.toColumnName(),
			"granted_org_name",
			ProjectGrantColumnGrantedRoleKeys.toColumnName(),
			"resource_owner_name",
			"COUNT(grant_id) OVER ()").
			From(projection.ProjectGrantProjectionTable).PlaceholderFormat(sq.Dollar).
			LeftJoin(GenerateJoinQuery(joins)),
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
				return nil, errors.ThrowInternal(err, "QUERY-K9gEE", "unable to close rows")
			}

			return &ProjectGrants{
				ProjectGrants: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *Queries) prepareProjectGrantUniqueQuery() (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select("COUNT(*) = 0").
			From(projection.ProjectGrantProjectionTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (isUnique bool, err error) {
			err = row.Scan(&isUnique)
			if err != nil {
				return false, errors.ThrowInternal(err, "QUERY-j92fg", "errors.internal")
			}
			return isUnique, err
		}
}

func (q *Queries) ProjectGrantByID(ctx context.Context, id string) (*ProjectGrant, error) {
	stmt, scan := prepareProjectGrantQuery()
	query, args, err := stmt.Where(sq.Eq{
		ProjectGrantColumnGrantID.toColumnName(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-MO9fs", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ExistsProjectGrant(ctx context.Context, id string) (err error) {
	_, err = q.ProjectGrantByID(ctx, id)
	return err
}

func (q *Queries) SearchProjectGrants(ctx context.Context, queries *ProjectGrantSearchQueries) (projects *ProjectGrants, err error) {
	query, scan := q.prepareProjectGrantsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-N9fsg", "Errors.project_grants.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-PP02n", "Errors.project_grants.internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projection.ProjectGrantProjectionTable)
	return projects, err
}

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

func NewProjectGrantProjectIDSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnProjectID, value, method)
}

func NewProjectGrantProjectNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnProjectName, value, method)
}

func NewProjectGrantRoleKeySearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnGrantedRoleKeys, value, method)
}

func NewProjectGrantResourceOwnerSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnResourceOwner, value, method)
}

func (q *ProjectGrantSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func (r *ProjectGrantSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectGrantResourceOwnerSearchQuery(TextEquals, orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

type ProjectGrantColumn int32

const (
	ProjectGrantColumnCreationDate ProjectGrantColumn = iota + 1
	ProjectGrantColumnChangeDate
	ProjectGrantColumnResourceOwner
	ProjectGrantColumnState
	ProjectGrantColumnSequence
	ProjectGrantColumnProjectID
	ProjectGrantColumnProjectName
	ProjectGrantColumnOrgID
	ProjectGrantColumnOrgName
	ProjectGrantColumnGrantID
	ProjectGrantColumnGrantedRoleKeys
	ProjectGrantColumnResourceOwnerName
	ProjectGrantColumnCreatorName
)

func (c ProjectGrantColumn) toColumnName() string {
	switch c {
	case ProjectGrantColumnProjectID:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantProjectIDCol
	case ProjectGrantColumnGrantID:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantIDCol
	case ProjectGrantColumnCreationDate:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantCreationDateCol
	case ProjectGrantColumnChangeDate:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantChangeDateCol
	case ProjectGrantColumnResourceOwner:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantResourceOwnerCol
	case ProjectGrantColumnState:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantStateCol
	case ProjectGrantColumnSequence:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantSequenceCol
	case ProjectGrantColumnOrgID:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantGrantedOrgIDCol
	case ProjectGrantColumnGrantedRoleKeys:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantRoleKeysCol
	case ProjectGrantColumnCreatorName:
		return projection.ProjectGrantProjectionTable + "." + projection.ProjectGrantCreatorCol
	default:
		return ""
	}
}
