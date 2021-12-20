package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
	"github.com/lib/pq"
)

type Memberships struct {
	SearchResponse
	Memberships []*Membership
}

type Membership struct {
	UserID        string
	Roles         []string
	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	Org          *OrgMembership
	IAM          *IAMMembership
	Project      *ProjectMembership
	ProjectGrant *ProjectGrantMembership
}

type OrgMembership struct {
	OrgID string
	Name  string
}

type IAMMembership struct {
	IAMID string
	Name  string
}

type ProjectMembership struct {
	ProjectID string
	Name      string
}

type ProjectGrantMembership struct {
	ProjectID   string
	ProjectName string
	GrantID     string
}

type MembershipSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func NewMembershipUserIDQuery(userID string) (SearchQuery, error) {
	return NewTextQuery(membershipUserID.setTable(membershipAlias), userID, TextEquals)
}

func NewMembershipResourceOwnerQuery(value string) (SearchQuery, error) {
	return NewTextQuery(membershipResourceOwner.setTable(membershipAlias), value, TextEquals)
}

func NewMembershipOrgIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(membershipOrgID, value, TextEquals)
}

func NewMembershipProjectIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(membershipProjectID, value, TextEquals)
}

func NewMembershipProjectGrantIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(membershipGrantID, value, TextEquals)
}

func NewMembershipIsIAMQuery() (SearchQuery, error) {
	return NewNotNullQuery(membershipIAMID)
}

func (q *MembershipSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) Memberships(ctx context.Context, queries *MembershipSearchQuery) (*Memberships, error) {
	query, scan := prepareMembershipsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-T84X9", "Errors.Query.InvalidRequest")
	}
	latestSequence, err := q.latestSequence(ctx, orgMemberTable, iamMemberTable, projectMemberTable, projectGrantMemberTable)
	if err != nil {
		return nil, err
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-eAV2x", "Errors.Internal")
	}
	memberships, err := scan(rows)
	if err != nil {
		return nil, err
	}
	memberships.LatestSequence = latestSequence
	return memberships, nil
}

var (
	//membershipAlias is a hack to satisfy checks in the queries
	membershipAlias = table{
		name: "memberships",
	}
	membershipUserID = Column{
		name:  projection.MemberUserIDCol,
		table: membershipAlias,
	}
	membershipRoles = Column{
		name:  projection.MemberRolesCol,
		table: membershipAlias,
	}
	membershipCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: membershipAlias,
	}
	membershipChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: membershipAlias,
	}
	membershipSequence = Column{
		name:  projection.MemberSequence,
		table: membershipAlias,
	}
	membershipResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: membershipAlias,
	}
	membershipOrgID = Column{
		name:  projection.OrgMemberOrgIDCol,
		table: membershipAlias,
	}
	membershipIAMID = Column{
		name:  projection.IAMMemberIAMIDCol,
		table: membershipAlias,
	}
	membershipProjectID = Column{
		name:  projection.ProjectMemberProjectIDCol,
		table: membershipAlias,
	}
	membershipGrantID = Column{
		name:  projection.ProjectGrantMemberGrantIDCol,
		table: membershipAlias,
	}

	membershipFrom = "(" +
		prepareOrgMember() +
		" UNION ALL " +
		prepareIAMMember() +
		" UNION ALL " +
		prepareProjectMember() +
		" UNION ALL " +
		prepareProjectGrantMember() +
		") AS " + membershipAlias.identifier()
)

func prepareMembershipsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Memberships, error)) {
	return sq.Select(
			membershipUserID.identifier(),
			membershipRoles.identifier(),
			membershipCreationDate.identifier(),
			membershipChangeDate.identifier(),
			membershipSequence.identifier(),
			membershipResourceOwner.identifier(),
			membershipOrgID.identifier(),
			membershipIAMID.identifier(),
			membershipProjectID.identifier(),
			membershipGrantID.identifier(),
			ProjectColumnName.identifier(),
			OrgColumnName.identifier(),
			countColumn.identifier(),
		).From(membershipFrom).
			LeftJoin(join(ProjectColumnID, membershipProjectID)).
			LeftJoin(join(OrgColumnID, membershipOrgID)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Memberships, error) {
			memberships := make([]*Membership, 0)
			var count uint64
			for rows.Next() {

				var (
					membership  = new(Membership)
					orgID       = sql.NullString{}
					iamID       = sql.NullString{}
					projectID   = sql.NullString{}
					grantID     = sql.NullString{}
					roles       = pq.StringArray{}
					projectName = sql.NullString{}
					orgName     = sql.NullString{}
				)

				err := rows.Scan(
					&membership.UserID,
					&roles,
					&membership.CreationDate,
					&membership.ChangeDate,
					&membership.Sequence,
					&membership.ResourceOwner,
					&orgID,
					&iamID,
					&projectID,
					&grantID,
					&projectName,
					&orgName,
					&count,
				)

				if err != nil {
					return nil, err
				}

				membership.Roles = roles

				if orgID.Valid {
					membership.Org = &OrgMembership{
						OrgID: orgID.String,
						Name:  orgName.String,
					}
				} else if iamID.Valid {
					membership.IAM = &IAMMembership{
						IAMID: iamID.String,
						Name:  iamID.String,
					}
				} else if projectID.Valid && grantID.Valid {
					membership.ProjectGrant = &ProjectGrantMembership{
						ProjectID:   projectID.String,
						ProjectName: projectName.String,
						GrantID:     grantID.String,
					}
				} else if projectID.Valid {
					membership.Project = &ProjectMembership{
						ProjectID: projectID.String,
						Name:      projectName.String,
					}
				}

				memberships = append(memberships, membership)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-N34NV", "Errors.Query.CloseRows")
			}

			return &Memberships{
				Memberships: memberships,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareOrgMember() string {
	stmt, _ := sq.Select(
		OrgMemberUserID.identifier(),
		OrgMemberRoles.identifier(),
		OrgMemberCreationDate.identifier(),
		OrgMemberChangeDate.identifier(),
		OrgMemberSequence.identifier(),
		OrgMemberResourceOwner.identifier(),
		OrgMemberOrgID.identifier(),
		"NULL::STRING AS "+membershipIAMID.name,
		"NULL::STRING AS "+membershipProjectID.name,
		"NULL::STRING AS "+membershipGrantID.name,
	).From(orgMemberTable.identifier()).MustSql()
	return stmt
}

func prepareIAMMember() string {
	stmt, _ := sq.Select(
		IAMMemberUserID.identifier(),
		IAMMemberRoles.identifier(),
		IAMMemberCreationDate.identifier(),
		IAMMemberChangeDate.identifier(),
		IAMMemberSequence.identifier(),
		IAMMemberResourceOwner.identifier(),
		"NULL::STRING AS "+membershipOrgID.name,
		IAMMemberIAMID.identifier(),
		"NULL::STRING AS "+membershipProjectID.name,
		"NULL::STRING AS "+membershipGrantID.name,
	).From(iamMemberTable.identifier()).MustSql()
	return stmt
}

func prepareProjectMember() string {
	stmt, _ := sq.Select(
		ProjectMemberUserID.identifier(),
		ProjectMemberRoles.identifier(),
		ProjectMemberCreationDate.identifier(),
		ProjectMemberChangeDate.identifier(),
		ProjectMemberSequence.identifier(),
		ProjectMemberResourceOwner.identifier(),
		"NULL::STRING AS "+membershipOrgID.name,
		"NULL::STRING AS "+membershipIAMID.name,
		ProjectMemberProjectID.identifier(),
		"NULL::STRING AS "+membershipGrantID.name,
	).From(projectMemberTable.identifier()).MustSql()

	return stmt
}

func prepareProjectGrantMember() string {
	stmt, _ := sq.Select(
		ProjectGrantMemberUserID.identifier(),
		ProjectGrantMemberRoles.identifier(),
		ProjectGrantMemberCreationDate.identifier(),
		ProjectGrantMemberChangeDate.identifier(),
		ProjectGrantMemberSequence.identifier(),
		ProjectGrantMemberResourceOwner.identifier(),
		"NULL::STRING AS "+membershipOrgID.name,
		"NULL::STRING AS "+membershipIAMID.name,
		ProjectGrantMemberProjectID.identifier(),
		ProjectGrantMemberGrantID.identifier(),
	).From(projectGrantMemberTable.identifier()).MustSql()

	return stmt
}
