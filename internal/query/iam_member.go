package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

var (
	instanceMemberTable = table{
		name:          projection.InstanceMemberProjectionTable,
		alias:         "members",
		instanceIDCol: projection.MemberInstanceID,
	}
	InstanceMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: instanceMemberTable,
	}
	InstanceMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: instanceMemberTable,
	}
	InstanceMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: instanceMemberTable,
	}
	InstanceMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: instanceMemberTable,
	}
	InstanceMemberSequence = Column{
		name:  projection.MemberSequence,
		table: instanceMemberTable,
	}
	InstanceMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: instanceMemberTable,
	}
	InstanceMemberInstanceID = Column{
		name:  projection.MemberInstanceID,
		table: instanceMemberTable,
	}
	InstanceMemberIAMID = Column{
		name:  projection.InstanceMemberIAMIDCol,
		table: instanceMemberTable,
	}
	InstanceMemberOwnerRemoved = Column{
		name:  projection.MemberOwnerRemoved,
		table: instanceMemberTable,
	}
	InstanceMemberOwnerRemovedUser = Column{
		name:  projection.MemberUserOwnerRemoved,
		table: instanceMemberTable,
	}
)

type IAMMembersQuery struct {
	MembersQuery
}

func (q *IAMMembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.MembersQuery.
		toQuery(query)
}

func addIamMemberWithoutOwnerRemoved(eq map[string]interface{}) {
	eq[InstanceMemberOwnerRemoved.identifier()] = false
	eq[InstanceMemberOwnerRemovedUser.identifier()] = false
}

func (q *Queries) IAMMembers(ctx context.Context, queries *IAMMembersQuery, withOwnerRemoved bool) (_ *Members, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareInstanceMembersQuery(ctx, q.client)
	eq := sq.Eq{InstanceMemberInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		addIamMemberWithoutOwnerRemoved(eq)
		addLoginNameWithoutOwnerRemoved(eq)
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-USNwM", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestSequence(ctx, instanceMemberTable)
	if err != nil {
		return nil, err
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-Pdg1I", "Errors.Internal")
	}
	members, err := scan(rows)
	if err != nil {
		return nil, err
	}
	members.LatestSequence = currentSequence
	return members, err
}

func prepareInstanceMembersQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Members, error)) {
	return sq.Select(
			InstanceMemberCreationDate.identifier(),
			InstanceMemberChangeDate.identifier(),
			InstanceMemberSequence.identifier(),
			InstanceMemberResourceOwner.identifier(),
			InstanceMemberUserID.identifier(),
			InstanceMemberRoles.identifier(),
			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			countColumn.identifier(),
		).From(instanceMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, InstanceMemberUserID)).
			LeftJoin(join(MachineUserIDCol, InstanceMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, InstanceMemberUserID) + db.Timetravel(call.Took(ctx))).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Members, error) {
			members := make([]*Member, 0)
			var count uint64

			for rows.Next() {
				member := new(Member)

				var (
					preferredLoginName = sql.NullString{}
					email              = sql.NullString{}
					firstName          = sql.NullString{}
					lastName           = sql.NullString{}
					displayName        = sql.NullString{}
					machineName        = sql.NullString{}
					avatarURL          = sql.NullString{}
				)

				err := rows.Scan(
					&member.CreationDate,
					&member.ChangeDate,
					&member.Sequence,
					&member.ResourceOwner,
					&member.UserID,
					&member.Roles,
					&preferredLoginName,
					&email,
					&firstName,
					&lastName,
					&displayName,
					&machineName,
					&avatarURL,

					&count,
				)

				if err != nil {
					return nil, err
				}

				member.PreferredLoginName = preferredLoginName.String
				member.Email = email.String
				member.FirstName = firstName.String
				member.LastName = lastName.String
				member.AvatarURL = avatarURL.String
				if displayName.Valid {
					member.DisplayName = displayName.String
				} else {
					member.DisplayName = machineName.String
				}

				members = append(members, member)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-EqJFc", "Errors.Query.CloseRows")
			}

			return &Members{
				Members: members,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
