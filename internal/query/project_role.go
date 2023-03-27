package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	projectRolesTable = table{
		name:          projection.ProjectRoleProjectionTable,
		instanceIDCol: projection.ProjectRoleColumnInstanceID,
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
	ProjectRoleColumnInstanceID = Column{
		name:  projection.ProjectRoleColumnInstanceID,
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
	ProjectRoleColumnOwnerRemoved = Column{
		name:  projection.ProjectRoleColumnOwnerRemoved,
		table: projectRolesTable,
	}
)

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

func (q *Queries) SearchProjectRoles(ctx context.Context, shouldTriggerBulk bool, queries *ProjectRoleSearchQueries, withOwnerRemoved bool) (projects *ProjectRoles, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		projection.ProjectRoleProjection.Trigger(ctx)
	}

	eq := sq.Eq{ProjectRoleColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[ProjectRoleColumnOwnerRemoved.identifier()] = false
	}

	query, scan := prepareProjectRolesQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-3N9ff", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5Ngd9", "Errors.Internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projectRolesTable)
	return projects, err
}

func (q *Queries) SearchGrantedProjectRoles(ctx context.Context, grantID, grantedOrg string, queries *ProjectRoleSearchQueries, withOwnerRemoved bool) (projects *ProjectRoles, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	grant, err := q.ProjectGrantByIDAndGrantedOrg(ctx, grantID, grantedOrg, withOwnerRemoved)
	if err != nil {
		return nil, err
	}
	err = queries.AppendRoleKeysQuery(grant.GrantedRoleKeys)
	if err != nil {
		return nil, err
	}

	eq := sq.Eq{ProjectRoleColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		eq[ProjectRoleColumnOwnerRemoved.identifier()] = false
	}

	query, scan := prepareProjectRolesQuery(ctx, q.client)
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-3N9ff", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5Ngd9", "Errors.Internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projectRolesTable)
	return projects, err
}

func NewProjectRoleProjectIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnProjectID, value, TextEquals)
}

func NewProjectRoleResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectRoleColumnResourceOwner, value, TextEquals)
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

func (r *ProjectRoleSearchQueries) AppendProjectIDQuery(projectID string) error {
	query, err := NewProjectRoleProjectIDSearchQuery(projectID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r *ProjectRoleSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectRoleResourceOwnerSearchQuery(orgID)
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

func (q *ProjectRoleSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareProjectRolesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*ProjectRoles, error)) {
	return sq.Select(
			ProjectRoleColumnProjectID.identifier(),
			ProjectRoleColumnCreationDate.identifier(),
			ProjectRoleColumnChangeDate.identifier(),
			ProjectRoleColumnResourceOwner.identifier(),
			ProjectRoleColumnSequence.identifier(),
			ProjectRoleColumnKey.identifier(),
			ProjectRoleColumnDisplayName.identifier(),
			ProjectRoleColumnGroupName.identifier(),
			countColumn.identifier()).
			From(projectRolesTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
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
				return nil, errors.ThrowInternal(err, "QUERY-ML0Fs", "Errors.Query.CloseRows")
			}

			return &ProjectRoles{
				ProjectRoles: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
