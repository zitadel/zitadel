package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	iamMemberTable = table{
		name:  projection.IAMMemberProjectionTable,
		alias: "members",
	}
	IAMMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: iamMemberTable,
	}
	IAMMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: iamMemberTable,
	}
	IAMMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: iamMemberTable,
	}
	IAMMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: iamMemberTable,
	}
	IAMMemberSequence = Column{
		name:  projection.MemberSequence,
		table: iamMemberTable,
	}
	IAMMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: iamMemberTable,
	}
	IAMMemberIAMID = Column{
		name:  projection.IAMMemberIAMIDCol,
		table: iamMemberTable,
	}
)

type IAMMembersQuery struct {
	MembersQuery
}

func (q *IAMMembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.MembersQuery.
		toQuery(query)
}

func (q *Queries) IAMMembers(ctx context.Context, queries *IAMMembersQuery) (*Members, error) {
	query, scan := prepareIAMMembersQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-USNwM", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestSequence(ctx, iamMemberTable)
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

func prepareIAMMembersQuery() (sq.SelectBuilder, func(*sql.Rows) (*Members, error)) {
	return sq.Select(
			IAMMemberCreationDate.identifier(),
			IAMMemberChangeDate.identifier(),
			IAMMemberSequence.identifier(),
			IAMMemberResourceOwner.identifier(),
			IAMMemberUserID.identifier(),
			IAMMemberRoles.identifier(),
			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			countColumn.identifier(),
		).From(iamMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, IAMMemberUserID)).
			LeftJoin(join(MachineUserIDCol, IAMMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, IAMMemberUserID)).
			Where(
				sq.Eq{LoginNameIsPrimaryCol.identifier(): true},
			).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Members, error) {
			members := make([]*Member, 0)
			var count uint64

			for rows.Next() {
				member := new(Member)
				roles := pq.StringArray{}

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
					&roles,
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

				member.Roles = roles
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
