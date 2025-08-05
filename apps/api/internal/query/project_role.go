package query

import (
	"context"
	"database/sql"
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

func projectRolesCheckPermission(ctx context.Context, projectRoles *ProjectRoles, permissionCheck domain.PermissionCheck) {
	projectRoles.ProjectRoles = slices.DeleteFunc(projectRoles.ProjectRoles,
		func(projectRole *ProjectRole) bool {
			return projectRoleCheckPermission(ctx, projectRole.ResourceOwner, projectRole.Key, permissionCheck) != nil
		},
	)
}

func projectRoleCheckPermission(ctx context.Context, resourceOwner string, grantID string, permissionCheck domain.PermissionCheck) error {
	return permissionCheck(ctx, domain.PermissionProjectGrantRead, resourceOwner, grantID)
}

func projectRolePermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *ProjectRoleSearchQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		ProjectRoleColumnResourceOwner,
		domain.PermissionProjectRoleRead,
		SingleOrgPermissionOption(queries.Queries),
		WithProjectsPermissionOption(ProjectRoleColumnProjectID),
	)
	return query.JoinClause(join, args...)
}

func (q *Queries) SearchProjectRoles(ctx context.Context, shouldTriggerBulk bool, queries *ProjectRoleSearchQueries, permissionCheck domain.PermissionCheck) (roles *ProjectRoles, err error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	projectRoles, err := q.searchProjectRoles(ctx, shouldTriggerBulk, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		projectRolesCheckPermission(ctx, projectRoles, permissionCheck)
	}
	return projectRoles, nil
}

func (q *Queries) searchProjectRoles(ctx context.Context, shouldTriggerBulk bool, queries *ProjectRoleSearchQueries, permissionCheckV2 bool) (roles *ProjectRoles, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerProjectRoleProjection")
		ctx, err = projection.ProjectRoleProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("trigger failed")
		traceSpan.EndWithError(err)
	}

	eq := sq.Eq{ProjectRoleColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}

	query, scan := prepareProjectRolesQuery()
	query = projectRolePermissionCheckV2(ctx, query, permissionCheckV2, queries)
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-3N9ff", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		roles, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-5Ngd9", "Errors.Internal")
	}
	roles.State, err = q.latestState(ctx, projectRolesTable)
	return roles, err
}

func (q *Queries) SearchGrantedProjectRoles(ctx context.Context, grantID, grantedOrg string, queries *ProjectRoleSearchQueries) (roles *ProjectRoles, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	grant, err := q.ProjectGrantByIDAndGrantedOrg(ctx, grantID, grantedOrg)
	if err != nil {
		return nil, err
	}
	err = queries.AppendRoleKeysQuery(grant.GrantedRoleKeys)
	if err != nil {
		return nil, err
	}

	eq := sq.Eq{ProjectRoleColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}

	query, scan := prepareProjectRolesQuery()
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-3N9ff", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		roles, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-5Ngd9", "Errors.Internal")
	}

	roles.State, err = q.latestState(ctx, projectRolesTable)
	return roles, err
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

func prepareProjectRolesQuery() (sq.SelectBuilder, func(*sql.Rows) (*ProjectRoles, error)) {
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
			From(projectRolesTable.identifier()).
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
				return nil, zerrors.ThrowInternal(err, "QUERY-ML0Fs", "Errors.Query.CloseRows")
			}

			return &ProjectRoles{
				ProjectRoles: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
