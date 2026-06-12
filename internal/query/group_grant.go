package query

import (
	"context"
	"database/sql"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	groupGrantsTable = table{
		name:          projection.GroupGrantProjectionTable,
		instanceIDCol: projection.GroupGrantInstanceID,
	}
	GroupGrantColumnID = Column{
		name:  projection.GroupGrantID,
		table: groupGrantsTable,
	}
	GroupGrantColumnCreationDate = Column{
		name:  projection.GroupGrantCreationDate,
		table: groupGrantsTable,
	}
	GroupGrantColumnChangeDate = Column{
		name:  projection.GroupGrantChangeDate,
		table: groupGrantsTable,
	}
	GroupGrantColumnSequence = Column{
		name:  projection.GroupGrantSequence,
		table: groupGrantsTable,
	}
	GroupGrantColumnState = Column{
		name:  projection.GroupGrantState,
		table: groupGrantsTable,
	}
	GroupGrantColumnResourceOwner = Column{
		name:  projection.GroupGrantResourceOwner,
		table: groupGrantsTable,
	}
	GroupGrantColumnInstanceID = Column{
		name:  projection.GroupGrantInstanceID,
		table: groupGrantsTable,
	}
	GroupGrantColumnGroupID = Column{
		name:  projection.GroupGrantGroupID,
		table: groupGrantsTable,
	}
	GroupGrantColumnProjectID = Column{
		name:  projection.GroupGrantProjectID,
		table: groupGrantsTable,
	}
	GroupGrantColumnGrantID = Column{
		name:  projection.GroupGrantGrantID,
		table: groupGrantsTable,
	}
	GroupGrantColumnRoles = Column{
		name:  projection.GroupGrantRoles,
		table: groupGrantsTable,
	}
)

type GroupGrants struct {
	SearchResponse
	GroupGrants []*GroupGrant
}

type GroupGrant struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	GroupID        string
	GroupName      string
	ProjectID      string
	ProjectGrantID string
	RoleKeys       []string
}

type GroupGrantsSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func NewGroupGrantGroupIDsSearchQuery(groupIDs []string) (SearchQuery, error) {
	list := make([]interface{}, len(groupIDs))
	for i, value := range groupIDs {
		list[i] = value
	}
	return NewListQuery(GroupGrantColumnGroupID, list, ListIn)
}

func NewGroupGrantProjectIDSearchQuery(projectID string) (SearchQuery, error) {
	return NewTextQuery(GroupGrantColumnProjectID, projectID, TextEquals)
}

// SearchGroupGrants returns the group grants matching the queries.
// Each returned grant carries the name of the group that supplies it for provenance.
func (q *Queries) SearchGroupGrants(ctx context.Context, queries *GroupGrantsSearchQuery, permissionCheck domain.PermissionCheck) (_ *GroupGrants, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	permissionCheckV2 := PermissionV2(ctx, permissionCheck)

	grants, err := q.searchGroupGrants(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !permissionCheckV2 {
		groupGrantsCheckPermission(ctx, grants, permissionCheck)
	}
	return grants, nil
}

func (q *Queries) searchGroupGrants(ctx context.Context, queries *GroupGrantsSearchQuery, permissionCheckV2 bool) (grants *GroupGrants, err error) {
	query, scan := prepareGroupGrantsQuery()
	query = groupGrantPermissionCheckV2(ctx, query, permissionCheckV2)
	eq := sq.Eq{
		GroupGrantColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-gG4kVx", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		grants, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-jB6wPl", "Errors.Internal")
	}
	grants.State, err = q.latestState(ctx, groupGrantsTable)
	return grants, err
}

func groupGrantsCheckPermission(ctx context.Context, grants *GroupGrants, permissionCheck domain.PermissionCheck) {
	grants.GroupGrants = slices.DeleteFunc(grants.GroupGrants,
		func(grant *GroupGrant) bool {
			return permissionCheck(ctx, domain.PermissionGroupGrantRead, grant.ResourceOwner, grant.GroupID) != nil
		},
	)
}

func groupGrantPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, permissionCheckV2 bool) sq.SelectBuilder {
	if !permissionCheckV2 {
		return query
	}

	join, args := PermissionClause(
		ctx,
		GroupGrantColumnResourceOwner,
		domain.PermissionGroupGrantRead,
	)
	return query.JoinClause(join, args...)
}

func prepareGroupGrantsQuery() (sq.SelectBuilder, func(*sql.Rows) (*GroupGrants, error)) {
	return sq.Select(
			GroupGrantColumnID.identifier(),
			GroupGrantColumnCreationDate.identifier(),
			GroupGrantColumnChangeDate.identifier(),
			GroupGrantColumnSequence.identifier(),
			GroupGrantColumnResourceOwner.identifier(),
			GroupGrantColumnGroupID.identifier(),
			GroupColumnName.identifier(),
			GroupGrantColumnProjectID.identifier(),
			GroupGrantColumnGrantID.identifier(),
			GroupGrantColumnRoles.identifier(),
			countColumn.identifier()).
			From(groupGrantsTable.identifier()).
			LeftJoin(join(GroupColumnID, GroupGrantColumnGroupID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GroupGrants, error) {
			grants := make([]*GroupGrant, 0)
			var count uint64
			for rows.Next() {
				grant := new(GroupGrant)
				var (
					groupName sql.NullString
					grantID   sql.NullString
					roles     database.TextArray[string]
				)

				err := rows.Scan(
					&grant.ID,
					&grant.CreationDate,
					&grant.ChangeDate,
					&grant.Sequence,
					&grant.ResourceOwner,
					&grant.GroupID,
					&groupName,
					&grant.ProjectID,
					&grantID,
					&roles,
					&count,
				)
				if err != nil {
					return nil, err
				}

				grant.GroupName = groupName.String
				grant.ProjectGrantID = grantID.String
				grant.RoleKeys = roles
				grants = append(grants, grant)
			}
			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-pV3mKw", "Errors.Query.CloseRows")
			}

			return &GroupGrants{
				GroupGrants: grants,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *GroupGrantsSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}
