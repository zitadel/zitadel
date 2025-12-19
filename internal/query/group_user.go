package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GroupUser struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	UserID             string
	GroupID            string
	GroupName          string
	GroupDescription   string
	Roles              database.TextArray[string]
	Attributes         database.TextArray[string]
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
	UserType           domain.UserType
}

type GroupUsers struct {
	SearchResponse
	GroupUsers []*GroupUser
}

var (
	groupUserTable = table{
		name:          projection.GroupUserProjectionTable,
		alias:         "groupusers",
		instanceIDCol: projection.GroupUserInstanceID,
	}
	GroupUserUserID = Column{
		name:  projection.GroupUserUserIDCol,
		table: groupUserTable,
	}
	GroupUserGroupID = Column{
		name:  projection.GroupUserGroupIDCol,
		table: groupUserTable,
	}
	GroupUserCreationDate = Column{
		name:  projection.GroupUserCreationDate,
		table: groupUserTable,
	}
	GroupUserChangeDate = Column{
		name:  projection.GroupUserChangeDate,
		table: groupUserTable,
	}
	GroupUserSequence = Column{
		name:  projection.GroupUserSequence,
		table: groupUserTable,
	}
	GroupUserResourceOwner = Column{
		name:  projection.GroupUserResourceOwner,
		table: groupUserTable,
	}
	GroupUserInstanceID = Column{
		name:  projection.GroupUserInstanceID,
		table: groupUserTable,
	}
	GroupUserAttributes = Column{
		name:  projection.GroupUserAttributes,
		table: groupUserTable,
	}
)

type GroupUsersQuery struct {
	SearchRequest
	Queries []SearchQuery
	GroupID string
}

func (q *GroupUsersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query.Where(sq.Eq{GroupUserGroupID.identifier(): q.GroupID})
}

func NewGroupUserEmailSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanEmailCol, value, method)
}

func NewGroupUserFirstNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanFirstNameCol, value, method)
}

func NewGroupUserLastNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(HumanLastNameCol, value, method)
}

func NewGroupUserUserIDSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(GroupUserUserID, value, TextEquals)
}

func NewGroupUserAttributesSearchQuery(value []string) (SearchQuery, error) {
	return NewListContains(GroupUserAttributes, value)
}

func NewGroupUserResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(GroupUserResourceOwner, value, TextEquals)
}

func (q *Queries) GroupUsers(ctx context.Context, queries *GroupUsersQuery) (groupUsers *GroupUsers, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareGroupUsersQuery()
	eq := sq.Eq{GroupUserInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-S8UL0", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestState(ctx, groupUserTable)
	if err != nil {
		return nil, err
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		groupUsers, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-vk7pj", "Errors.Internal")
	}

	groupUsers.State = currentSequence
	return groupUsers, err
}

func prepareGroupUsersQuery() (sq.SelectBuilder, func(*sql.Rows) (*GroupUsers, error)) {
	return sq.Select(
			GroupUserCreationDate.identifier(),
			GroupUserChangeDate.identifier(),
			GroupUserSequence.identifier(),
			GroupUserResourceOwner.identifier(),
			GroupUserUserID.identifier(),
			GroupUserGroupID.identifier(),
			GroupUserAttributes.identifier(),

			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			UserTypeCol.identifier(),
			countColumn.identifier(),
		).From(groupUserTable.identifier()).
			LeftJoin(join(HumanUserIDCol, GroupUserUserID)).
			LeftJoin(join(MachineUserIDCol, GroupUserUserID)).
			LeftJoin(join(UserIDCol, GroupUserUserID)).
			LeftJoin(join(LoginNameUserIDCol, GroupUserUserID)).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GroupUsers, error) {
			groupUsers := make([]*GroupUser, 0)
			var count uint64

			for rows.Next() {
				groupUser := new(GroupUser)

				var (
					preferredLoginName = sql.NullString{}
					email              = sql.NullString{}
					firstName          = sql.NullString{}
					lastName           = sql.NullString{}
					displayName        = sql.NullString{}
					machineName        = sql.NullString{}
					avatarURL          = sql.NullString{}
					userType           = sql.NullInt32{}
				)

				err := rows.Scan(
					&groupUser.CreationDate,
					&groupUser.ChangeDate,
					&groupUser.Sequence,
					&groupUser.ResourceOwner,
					&groupUser.UserID,
					&groupUser.GroupID,
					&groupUser.Attributes,
					&preferredLoginName,
					&email,
					&firstName,
					&lastName,
					&displayName,
					&machineName,
					&avatarURL,
					&userType,

					&count,
				)

				if err != nil {
					return nil, err
				}

				groupUser.PreferredLoginName = preferredLoginName.String
				groupUser.Email = email.String
				groupUser.FirstName = firstName.String
				groupUser.LastName = lastName.String
				groupUser.AvatarURL = avatarURL.String
				if displayName.Valid {
					groupUser.DisplayName = displayName.String
				} else {
					groupUser.DisplayName = machineName.String
				}
				groupUser.UserType = domain.UserType(userType.Int32)

				groupUsers = append(groupUsers, groupUser)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-XK2Jj", "Errors.Query.CloseRows")
			}

			return &GroupUsers{
				GroupUsers: groupUsers,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
