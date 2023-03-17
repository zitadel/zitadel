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
	orgMemberTable = table{
		name:          projection.OrgMemberProjectionTable,
		alias:         "members",
		instanceIDCol: projection.MemberInstanceID,
	}
	OrgMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: orgMemberTable,
	}
	OrgMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: orgMemberTable,
	}
	OrgMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: orgMemberTable,
	}
	OrgMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: orgMemberTable,
	}
	OrgMemberSequence = Column{
		name:  projection.MemberSequence,
		table: orgMemberTable,
	}
	OrgMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: orgMemberTable,
	}
	OrgMemberInstanceID = Column{
		name:  projection.MemberInstanceID,
		table: orgMemberTable,
	}
	OrgMemberOrgID = Column{
		name:  projection.OrgMemberOrgIDCol,
		table: orgMemberTable,
	}
	OrgMemberOwnerRemoved = Column{
		name:  projection.MemberOwnerRemoved,
		table: orgMemberTable,
	}
	OrgMemberOwnerRemovedUser = Column{
		name:  projection.MemberUserOwnerRemoved,
		table: orgMemberTable,
	}
)

type OrgMembersQuery struct {
	MembersQuery
	OrgID string
}

func (q *OrgMembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.MembersQuery.
		toQuery(query).
		Where(sq.Eq{OrgMemberOrgID.identifier(): q.OrgID})
}

func addOrgMemberWithoutOwnerRemoved(eq map[string]interface{}) {
	eq[OrgMemberOwnerRemoved.identifier()] = false
	eq[OrgMemberOwnerRemovedUser.identifier()] = false
}

func (q *Queries) OrgMembers(ctx context.Context, queries *OrgMembersQuery, withOwnerRemoved bool) (_ *Members, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareOrgMembersQuery(ctx, q.client)
	eq := sq.Eq{OrgMemberInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		addOrgMemberWithoutOwnerRemoved(eq)
		addLoginNameWithoutOwnerRemoved(eq)
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-PDAVB", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestSequence(ctx, orgsTable)
	if err != nil {
		return nil, err
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-5g4yV", "Errors.Internal")
	}
	members, err := scan(rows)
	if err != nil {
		return nil, err
	}
	members.LatestSequence = currentSequence
	return members, err
}

func prepareOrgMembersQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Members, error)) {
	return sq.Select(
			OrgMemberCreationDate.identifier(),
			OrgMemberChangeDate.identifier(),
			OrgMemberSequence.identifier(),
			OrgMemberResourceOwner.identifier(),
			OrgMemberUserID.identifier(),
			OrgMemberRoles.identifier(),
			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			countColumn.identifier(),
		).From(orgMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, OrgMemberUserID)).
			LeftJoin(join(MachineUserIDCol, OrgMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, OrgMemberUserID) + db.Timetravel(call.Took(ctx))).
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
				return nil, errors.ThrowInternal(err, "QUERY-N34NV", "Errors.Query.CloseRows")
			}

			return &Members{
				Members: members,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
