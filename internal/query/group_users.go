package query

import (
	"context"
	"database/sql"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	groupUsersTable = table{
		name:          projection.GroupUsersProjectionTable,
		instanceIDCol: projection.GroupUsersColumnInstanceID,
	}
	GroupUsersColumnGroupID = Column{
		name:  projection.GroupUsersColumnGroupID,
		table: groupUsersTable,
	}
	GroupUsersColumnUserID = Column{
		name:  projection.GroupUsersColumnUserID,
		table: groupUsersTable,
	}
	GroupUsersColumnResourceOwner = Column{
		name:  projection.GroupUsersColumnResourceOwner,
		table: groupUsersTable,
	}
	GroupUsersColumnCreationDate = Column{
		name:  projection.GroupUsersColumnCreationDate,
		table: groupUsersTable,
	}
	GroupUsersColumnInstanceID = Column{
		name:  projection.GroupUsersColumnInstanceID,
		table: groupUsersTable,
	}
	GroupUsersColumnSequence = Column{
		name:  projection.GroupUsersColumnSequence,
		table: groupUsersTable,
	}
)

type GroupUsers struct {
	SearchResponse
	GroupUsers []*GroupUser
}

type GroupUser struct {
	GroupID       string
	GroupName     string
	ResourceOwner string
	CreationDate  time.Time
	Sequence      uint64

	// user fields
	UserID             string
	PreferredLoginName string
	DisplayName        string
	AvatarUrl          string
}

type GroupUsersSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *Queries) SearchGroupUsers(ctx context.Context, queries *GroupUsersSearchQuery, permissionCheck domain.PermissionCheck) (_ *GroupUsers, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	permissionCheckV2 := PermissionV2(ctx, permissionCheck)

	groupUsers, err := q.searchGroupUsers(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !permissionCheckV2 {
		groupUsersCheckPermission(ctx, groupUsers, permissionCheck)
	}
	return groupUsers, nil
}

func NewGroupUsersUserIDsSearchQuery(userIDs []string) (SearchQuery, error) {
	list := make([]interface{}, len(userIDs))
	for i, value := range userIDs {
		list[i] = value
	}
	return NewListQuery(GroupUsersColumnUserID, list, ListIn)
}

func NewGroupUsersGroupIDsSearchQuery(groupIDs []string) (SearchQuery, error) {
	list := make([]interface{}, len(groupIDs))
	for i, value := range groupIDs {
		list[i] = value
	}
	return NewListQuery(GroupUsersColumnGroupID, list, ListIn)
}

func NewGroupUsersOrganizationIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(GroupUsersColumnResourceOwner, id, TextEquals)
}

func (q *Queries) searchGroupUsers(ctx context.Context, queries *GroupUsersSearchQuery, permissionCheckV2 bool) (groupUsers *GroupUsers, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareGroupUsersQuery()
	query = groupUsersPermissionCheckV2(ctx, query, queries, permissionCheckV2)
	eq := sq.And{
		sq.Eq{
			GroupUsersColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		},
	}
	// unique tiebreaker keeps pagination stable when sorting values collide,
	// e.g. members added in one batch sharing the same creation date
	stmt, args, err := queries.toQuery(query).
		OrderBy(GroupUsersColumnGroupID.identifier(), GroupUsersColumnUserID.identifier()).
		Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-TTlfF6", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		groupUsers, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-M5O50l", "Errors.Internal")
	}
	groupUsers.State, err = q.latestState(ctx, groupUsersTable)
	return groupUsers, nil
}

func prepareGroupUsersQuery() (query sq.SelectBuilder, scan func(*sql.Rows) (*GroupUsers, error)) {
	return sq.Select(
			GroupUsersColumnGroupID.identifier(),
			GroupColumnName.identifier(),
			GroupUsersColumnUserID.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			LoginNameNameCol.identifier(),
			GroupUsersColumnResourceOwner.identifier(),
			HumanAvatarURLCol.identifier(),
			GroupUsersColumnCreationDate.identifier(),
			GroupUsersColumnSequence.identifier(),
			countColumn.identifier(),
		).From(groupUsersTable.identifier()).
			LeftJoin(join(GroupColumnID, GroupUsersColumnGroupID)).
			LeftJoin(join(HumanUserIDCol, GroupUsersColumnUserID)).
			LeftJoin(join(MachineUserIDCol, GroupUsersColumnUserID)).
			LeftJoin(join(LoginNameUserIDCol, GroupUsersColumnUserID)).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GroupUsers, error) {
			groupUsers := make([]*GroupUser, 0)
			var count uint64
			for rows.Next() {
				g := new(GroupUser)

				var (
					groupName          sql.NullString
					displayName        sql.NullString
					machineName        sql.NullString
					avatarURL          sql.NullString
					preferredLoginName sql.NullString
				)

				err := rows.Scan(
					&g.GroupID,
					&groupName,
					&g.UserID,
					&displayName,
					&machineName,
					&preferredLoginName,
					&g.ResourceOwner,
					&avatarURL,
					&g.CreationDate,
					&g.Sequence,
					&count,
				)
				if err != nil {
					return nil, err
				}

				g.GroupName = groupName.String
				if displayName.Valid {
					g.DisplayName = displayName.String
				} else {
					g.DisplayName = machineName.String
				}
				g.AvatarUrl = avatarURL.String
				g.PreferredLoginName = preferredLoginName.String
				groupUsers = append(groupUsers, g)
			}
			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-JuX6i5", "Errors.Query.CloseRows")
			}
			return &GroupUsers{
				GroupUsers: groupUsers,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func groupUsersCheckPermission(ctx context.Context, groupUsers *GroupUsers, permissionCheck domain.PermissionCheck) {
	groupUsers.GroupUsers = slices.DeleteFunc(groupUsers.GroupUsers,
		func(gu *GroupUser) bool {
			return permissionCheck(ctx, domain.PermissionGroupUserRead, gu.ResourceOwner, gu.GroupID) != nil
		},
	)
}

func groupUsersPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, queries *GroupUsersSearchQuery, permissionCheckV2 bool) sq.SelectBuilder {
	if !permissionCheckV2 {
		return query
	}

	join, args := PermissionClause(
		ctx,
		GroupUsersColumnResourceOwner,
		domain.PermissionGroupUserRead,
	)

	return query.JoinClause(join, args...)
}

func (q *GroupUsersSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}
