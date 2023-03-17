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
	projectMemberTable = table{
		name:          projection.ProjectMemberProjectionTable,
		alias:         "members",
		instanceIDCol: projection.MemberInstanceID,
	}
	ProjectMemberUserID = Column{
		name:  projection.MemberUserIDCol,
		table: projectMemberTable,
	}
	ProjectMemberRoles = Column{
		name:  projection.MemberRolesCol,
		table: projectMemberTable,
	}
	ProjectMemberCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: projectMemberTable,
	}
	ProjectMemberChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: projectMemberTable,
	}
	ProjectMemberSequence = Column{
		name:  projection.MemberSequence,
		table: projectMemberTable,
	}
	ProjectMemberResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: projectMemberTable,
	}
	ProjectMemberInstanceID = Column{
		name:  projection.MemberInstanceID,
		table: projectMemberTable,
	}
	ProjectMemberProjectID = Column{
		name:  projection.ProjectMemberProjectIDCol,
		table: projectMemberTable,
	}
	ProjectMemberOwnerRemoved = Column{
		name:  projection.MemberOwnerRemoved,
		table: orgMemberTable,
	}
	ProjectMemberOwnerRemovedUser = Column{
		name:  projection.MemberUserOwnerRemoved,
		table: orgMemberTable,
	}
)

type ProjectMembersQuery struct {
	MembersQuery
	ProjectID string
}

func (q *ProjectMembersQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	return q.MembersQuery.
		toQuery(query).
		Where(sq.Eq{ProjectMemberProjectID.identifier(): q.ProjectID})
}

func addProjectMemberWithoutOwnerRemoved(eq map[string]interface{}) {
	eq[ProjectMemberOwnerRemoved.identifier()] = false
	eq[ProjectMemberOwnerRemovedUser.identifier()] = false
}

func (q *Queries) ProjectMembers(ctx context.Context, queries *ProjectMembersQuery, withOwnerRemoved bool) (_ *Members, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareProjectMembersQuery(ctx, q.client)
	eq := sq.Eq{ProjectMemberInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	if !withOwnerRemoved {
		addProjectMemberWithoutOwnerRemoved(eq)
		addLoginNameWithoutOwnerRemoved(eq)
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-T8CuT", "Errors.Query.InvalidRequest")
	}

	currentSequence, err := q.latestSequence(ctx, projectMemberTable)
	if err != nil {
		return nil, err
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-uh6pj", "Errors.Internal")
	}
	members, err := scan(rows)
	if err != nil {
		return nil, err
	}
	members.LatestSequence = currentSequence
	return members, err
}

func prepareProjectMembersQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Members, error)) {
	return sq.Select(
			ProjectMemberCreationDate.identifier(),
			ProjectMemberChangeDate.identifier(),
			ProjectMemberSequence.identifier(),
			ProjectMemberResourceOwner.identifier(),
			ProjectMemberUserID.identifier(),
			ProjectMemberRoles.identifier(),
			LoginNameNameCol.identifier(),
			HumanEmailCol.identifier(),
			HumanFirstNameCol.identifier(),
			HumanLastNameCol.identifier(),
			HumanDisplayNameCol.identifier(),
			MachineNameCol.identifier(),
			HumanAvatarURLCol.identifier(),
			countColumn.identifier(),
		).From(projectMemberTable.identifier()).
			LeftJoin(join(HumanUserIDCol, ProjectMemberUserID)).
			LeftJoin(join(MachineUserIDCol, ProjectMemberUserID)).
			LeftJoin(join(LoginNameUserIDCol, ProjectMemberUserID) + db.Timetravel(call.Took(ctx))).
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
				return nil, errors.ThrowInternal(err, "QUERY-ZJ1Ii", "Errors.Query.CloseRows")
			}

			return &Members{
				Members: members,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
