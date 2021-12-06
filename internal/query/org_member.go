package query

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/lib/pq"
)

var (
	orgMemberTable = table{
		name: projection.OrgMemberProjectionTable,
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
	OrgMemberOrgID = Column{
		name:  projection.OrgMemberOrgIDCol,
		table: orgMemberTable,
	}
)

type OrgMembersQuery struct {
	SearchRequest
}

func OrgMembers(ctx context.Context, orgID string, queries *OrgMembersQuery) (*Members, error) {
	return nil, nil
}

func prepareOrgMembersQuery() (sq.SelectBuilder, func(*sql.Rows) (*Members, error)) {
	return sq.Select(
			OrgMemberCreationDate.identifier(),
			OrgMemberChangeDate.identifier(),
			OrgMemberSequence.identifier(),
			OrgMemberResourceOwner.identifier(),
			OrgMemberUserID.identifier(),
			OrgMemberRoles.identifier(),
			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFistNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvaterURLCol.identifier(),
			countColumn.identifier(),
		).From(orgMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, OrgMemberUserID)).
			LeftJoin(join(MachineUserIDCol, OrgMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, OrgMemberUserID)).Where(
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
