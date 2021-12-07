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
	projectGrantMemberTable = table{
		name:  projection.ProjectGrantMemberProjectionTable,
		alias: "m",
	}
	ProjectGrantMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberSequence = Column{
		name:  projection.MemberSequence,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberProjectID = Column{
		name:  projection.ProjectGrantMemberProjectIDCol,
		table: projectGrantMemberTable,
	}
	ProjectGrantMemberGrantID = Column{
		name:  projection.ProjectGrantMemberGrantIDCol,
		table: projectGrantMemberTable,
	}
)

type ProjectGrantMembersQuery struct {
	MembersQuery
	ProjectID, GrantID string
}

func (q *ProjectGrantMembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.MembersQuery.
		toQuery(query).
		Where(sq.Eq{
			ProjectGrantMemberProjectID.identifier(): q.ProjectID,
			ProjectGrantMemberGrantID.identifier():   q.GrantID,
		})
}

func (q *Queries) ProjectGrantMembers(ctx context.Context, queries *ProjectGrantMembersQuery) (*Members, error) {
	query, scan := prepareProjectGrantMembersQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-USNwM", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestSequence(ctx, projectGrantMemberTable)
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

func prepareProjectGrantMembersQuery() (sq.SelectBuilder, func(*sql.Rows) (*Members, error)) {
	return sq.Select(
			ProjectGrantMemberCreationDate.identifier(),
			ProjectGrantMemberChangeDate.identifier(),
			ProjectGrantMemberSequence.identifier(),
			ProjectGrantMemberResourceOwner.identifier(),
			ProjectGrantMemberUserID.identifier(),
			ProjectGrantMemberRoles.identifier(),
			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvaterURLCol.identifier(),
			countColumn.identifier(),
		).From(projectGrantMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, ProjectGrantMemberUserID)).
			LeftJoin(join(MachineUserIDCol, ProjectGrantMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, ProjectGrantMemberUserID)).
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
