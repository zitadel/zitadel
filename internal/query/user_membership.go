package query

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
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
}

type IAMMembership struct {
	IAMID string
}

type ProjectMembership struct {
	ProjectID string
}

type ProjectGrantMembership struct {
	ProjectID string
	GrantID   string
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

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-eAV2x", "Errors.Internal")
	}
	memberships, err := scan(rows)
	if err != nil {
		return nil, err
	}
	// TODO: memberships.LatestSequence, err = q.latestSequence(ctx, idpTable)
	return memberships, err
}

var (
	//membershipAlias is a hack to satisfy checks in the queries
	membershipAlias = table{
		name: "m",
	}
	membershipUserID = Column{
		name: projection.MemberUserIDCol,
	}
	membershipRoles = Column{
		name: projection.MemberRolesCol,
	}
	membershipCreationDate = Column{
		name: projection.MemberCreationDate,
	}
	membershipChangeDate = Column{
		name: projection.MemberChangeDate,
	}
	membershipSequence = Column{
		name: projection.MemberSequence,
	}
	membershipResourceOwner = Column{
		name: projection.MemberResourceOwner,
	}
	membershipOrgID = Column{
		name: projection.OrgMemberOrgIDCol,
	}
	membershipIAMID = Column{
		name: projection.IAMMemberIAMIDCol,
	}
	membershipProjectID = Column{
		name: projection.ProjectMemberProjectIDCol,
	}
	membershipGrantID = Column{
		name: projection.ProjectGrantMemberGrantIDCol,
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
			countColumn.identifier(),
		).From(membershipFrom),
		func(rows *sql.Rows) (*Memberships, error) {
			memberships := make([]*Membership, 0)
			var count uint64
			for rows.Next() {
				membership := new(Membership)

				orgID := sql.NullString{}
				iamID := sql.NullString{}
				projectID := sql.NullString{}
				grantID := sql.NullString{}

				err := rows.Scan(
					&membership.UserID,
					&membership.Roles,
					&membership.CreationDate,
					&membership.ChangeDate,
					&membership.Sequence,
					&membership.ResourceOwner,
					&orgID,
					&iamID,
					&projectID,
					&grantID,
					&count,
				)

				if err != nil {
					return nil, err
				}

				if orgID.Valid {
					membership.Org = &OrgMembership{
						OrgID: orgID.String,
					}
				} else if iamID.Valid {
					membership.IAM = &IAMMembership{
						IAMID: iamID.String,
					}
				} else if projectID.Valid && grantID.Valid {
					membership.ProjectGrant = &ProjectGrantMembership{
						GrantID: grantID.String,
					}
				} else if projectID.Valid {
					membership.Project = &ProjectMembership{
						ProjectID: projectID.String,
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
		"NULL::STRING AS "+membershipIAMID.identifier(),
		"NULL::STRING AS "+membershipProjectID.identifier(),
		"NULL::STRING AS "+membershipGrantID.identifier(),
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
		"NULL::STRING AS "+membershipOrgID.identifier(),
		IAMMemberIAMID.identifier(),
		"NULL::STRING AS "+membershipProjectID.identifier(),
		"NULL::STRING AS "+membershipGrantID.identifier(),
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
		"NULL::STRING AS "+membershipOrgID.identifier(),
		"NULL::STRING AS "+membershipIAMID.identifier(),
		ProjectMemberProjectID.identifier(),
		"NULL::STRING AS "+membershipGrantID.identifier(),
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
		"NULL::STRING AS "+membershipOrgID.identifier(),
		"NULL::STRING AS "+membershipIAMID.identifier(),
		ProjectGrantMemberProjectID.identifier(),
		ProjectGrantMemberGrantID.identifier(),
	).From(projectGrantMemberTable.identifier()).MustSql()

	return stmt
}
