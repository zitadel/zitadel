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

type GroupMember struct {
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	UserID             string
	GroupID            string
	GroupName          string
	GroupDescription   string
	Roles              database.TextArray[string]
	PreferredLoginName string
	Email              string
	FirstName          string
	LastName           string
	DisplayName        string
	AvatarURL          string
	UserType           domain.UserType
}

type GroupMembers struct {
	SearchResponse
	GroupMembers []*GroupMember
}

var (
	groupMemberTable = table{
		name:          projection.GroupMemberProjectionTable,
		alias:         "members",
		instanceIDCol: projection.MemberInstanceID,
	}
	GroupMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: groupMemberTable,
	}
	GroupMemberGroupID = Column{
		name:  projection.GroupMemberGroupIDCol,
		table: groupMemberTable,
	}
	GroupMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: groupMemberTable,
	}
	GroupMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: groupMemberTable,
	}
	GroupMemberSequence = Column{
		name:  projection.MemberSequence,
		table: groupMemberTable,
	}
	GroupMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: groupMemberTable,
	}
	GroupMemberInstanceID = Column{
		name:  projection.MemberInstanceID,
		table: groupMemberTable,
	}
)

type GroupMembersQuery struct {
	MembersQuery
	GroupID string
}

func (q *GroupMembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.MembersQuery.
		toQuery(query).
		Where(sq.Eq{GroupMemberGroupID.identifier(): q.GroupID})
}

func (q *Queries) GroupMembers(ctx context.Context, queries *GroupMembersQuery) (groupMembers *GroupMembers, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareGroupMembersQuery()
	eq := sq.Eq{GroupMemberInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-S8UL0", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestState(ctx, groupMemberTable)
	if err != nil {
		return nil, err
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		groupMembers, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-vk7pj", "Errors.Internal")
	}

	groupMembers.State = currentSequence
	return groupMembers, err
}

func prepareGroupMembersQuery() (sq.SelectBuilder, func(*sql.Rows) (*GroupMembers, error)) {
	return sq.Select(
			GroupMemberCreationDate.identifier(),
			GroupMemberChangeDate.identifier(),
			GroupMemberSequence.identifier(),
			GroupMemberResourceOwner.identifier(),
			GroupMemberUserID.identifier(),
			GroupMemberGroupID.identifier(),

			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			UserTypeCol.identifier(),
			countColumn.identifier(),
		).From(groupMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, GroupMemberUserID)).
			LeftJoin(join(MachineUserIDCol, GroupMemberUserID)).
			LeftJoin(join(UserIDCol, GroupMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, GroupMemberUserID)).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*GroupMembers, error) {
			groupMembers := make([]*GroupMember, 0)
			var count uint64

			for rows.Next() {
				groupMember := new(GroupMember)

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
					&groupMember.CreationDate,
					&groupMember.ChangeDate,
					&groupMember.Sequence,
					&groupMember.ResourceOwner,
					&groupMember.UserID,
					&groupMember.GroupID,
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

				groupMember.PreferredLoginName = preferredLoginName.String
				groupMember.Email = email.String
				groupMember.FirstName = firstName.String
				groupMember.LastName = lastName.String
				groupMember.AvatarURL = avatarURL.String
				if displayName.Valid {
					groupMember.DisplayName = displayName.String
				} else {
					groupMember.DisplayName = machineName.String
				}
				groupMember.UserType = domain.UserType(userType.Int32)

				groupMembers = append(groupMembers, groupMember)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-XK2Jj", "Errors.Query.CloseRows")
			}

			return &GroupMembers{
				GroupMembers: groupMembers,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
