package query

import (
	"context"
	"database/sql"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Memberships struct {
	SearchResponse
	Memberships []*Membership
}

type Membership struct {
	UserID        string
	Roles         database.TextArray[string]
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
	ProjectID    string
	ProjectName  string
	GrantID      string
	GrantedOrgID string
}

type MembershipSearchQuery struct {
	SearchRequest
	Queries []SearchQuery
}

func NewMembershipUserIDQuery(userID string) (SearchQuery, error) {
	return NewTextQuery(membershipUserID.setTable(membershipAlias), userID, TextEquals)
}

func NewMembershipOrgIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgMemberOrgID, value, TextEquals)
}

func NewMembershipResourceOwnersSearchQuery(ids ...string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(membershipResourceOwner, list, ListIn)
}

func NewMembershipGrantedOrgIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnGrantedOrgID, id, TextEquals)
}

func NewMembershipProjectIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectMemberProjectID, value, TextEquals)
}

func NewMembershipProjectGrantIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantMemberGrantID, value, TextEquals)
}

func NewMembershipIsIAMQuery() (SearchQuery, error) {
	return NewNotNullQuery(InstanceMemberIAMID)
}

func (q *MembershipSearchQuery) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) Memberships(ctx context.Context, queries *MembershipSearchQuery, shouldTrigger bool) (memberships *Memberships, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTrigger {
		wg := sync.WaitGroup{}
		wg.Add(4)
		go func() {
			spanCtx, triggerSpan := tracing.NewNamedSpan(ctx, "TriggerOrgMemberProjection")
			_, _ = projection.OrgMemberProjection.Trigger(spanCtx, handler.WithAwaitRunning())
			triggerSpan.End()
			wg.Done()
		}()
		go func() {
			spanCtx, triggerSpan := tracing.NewNamedSpan(ctx, "TriggerInstanceMemberProjection")
			_, _ = projection.InstanceMemberProjection.Trigger(spanCtx, handler.WithAwaitRunning())
			triggerSpan.End()
			wg.Done()
		}()
		go func() {
			spanCtx, triggerSpan := tracing.NewNamedSpan(ctx, "TriggerProjectMemberProjection")
			_, _ = projection.ProjectMemberProjection.Trigger(spanCtx, handler.WithAwaitRunning())
			triggerSpan.End()
			wg.Done()
		}()
		go func() {
			spanCtx, triggerSpan := tracing.NewNamedSpan(ctx, "TriggerProjectGrantMemberProjection")
			_, _ = projection.ProjectGrantMemberProjection.Trigger(spanCtx, handler.WithAwaitRunning())
			triggerSpan.End()
			wg.Done()
		}()

		wg.Wait()
	}

	query, queryArgs, scan := prepareMembershipsQuery(ctx, q.client, queries)
	eq := sq.Eq{membershipInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-T84X9", "Errors.Query.InvalidRequest")
	}
	latestSequence, err := q.latestState(ctx, orgMemberTable, instanceMemberTable, projectMemberTable, projectGrantMemberTable)
	if err != nil {
		return nil, err
	}
	queryArgs = append(queryArgs, args...)

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		memberships, err = scan(rows)
		return err
	}, stmt, queryArgs...)
	if err != nil {
		return nil, err
	}
	memberships.State = latestSequence
	return memberships, nil
}

var (
	//membershipAlias is a hack to satisfy checks in the queries
	membershipAlias = table{
		name:          "members",
		instanceIDCol: projection.MemberInstanceID,
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
	membershipInstanceID = Column{
		name:  projection.MemberInstanceID,
		table: membershipAlias,
	}
	membershipOrgID = Column{
		name:  projection.OrgMemberOrgIDCol,
		table: membershipAlias,
	}
	membershipIAMID = Column{
		name:  projection.InstanceMemberIAMIDCol,
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
	membershipGrantGrantedOrgID = Column{
		name:  projection.ProjectGrantColumnGrantedOrgID,
		table: membershipAlias,
	}
)

func getMembershipFromQuery(queries *MembershipSearchQuery) (string, []interface{}) {
	orgMembers, orgMembersArgs := prepareOrgMember(queries)
	iamMembers, iamMembersArgs := prepareIAMMember(queries)
	projectMembers, projectMembersArgs := prepareProjectMember(queries)
	projectGrantMembers, projectGrantMembersArgs := prepareProjectGrantMember(queries)
	args := make([]interface{}, 0)
	args = append(append(append(append(args, orgMembersArgs...), iamMembersArgs...), projectMembersArgs...), projectGrantMembersArgs...)

	return "(" +
			orgMembers +
			" UNION ALL " +
			iamMembers +
			" UNION ALL " +
			projectMembers +
			" UNION ALL " +
			projectGrantMembers +
			") AS " + membershipAlias.identifier(),
		args
}

func prepareMembershipsQuery(ctx context.Context, db prepareDatabase, queries *MembershipSearchQuery) (sq.SelectBuilder, []interface{}, func(*sql.Rows) (*Memberships, error)) {
	query, args := getMembershipFromQuery(queries)
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
			ProjectGrantColumnGrantedOrgID.identifier(),
			ProjectColumnName.identifier(),
			OrgColumnName.identifier(),
			InstanceColumnName.identifier(),
			countColumn.identifier(),
		).From(query).
			LeftJoin(join(ProjectColumnID, membershipProjectID)).
			LeftJoin(join(OrgColumnID, membershipOrgID)).
			LeftJoin(join(ProjectGrantColumnGrantID, membershipGrantID)).
			LeftJoin(join(InstanceColumnID, membershipInstanceID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		args,
		func(rows *sql.Rows) (*Memberships, error) {
			memberships := make([]*Membership, 0)
			var count uint64
			for rows.Next() {

				var (
					membership   = new(Membership)
					orgID        = sql.NullString{}
					instanceID   = sql.NullString{}
					projectID    = sql.NullString{}
					grantID      = sql.NullString{}
					grantedOrgID = sql.NullString{}
					projectName  = sql.NullString{}
					orgName      = sql.NullString{}
					instanceName = sql.NullString{}
				)

				err := rows.Scan(
					&membership.UserID,
					&membership.Roles,
					&membership.CreationDate,
					&membership.ChangeDate,
					&membership.Sequence,
					&membership.ResourceOwner,
					&orgID,
					&instanceID,
					&projectID,
					&grantID,
					&grantedOrgID,
					&projectName,
					&orgName,
					&instanceName,
					&count,
				)

				if err != nil {
					return nil, err
				}

				if orgID.Valid {
					membership.Org = &OrgMembership{
						OrgID: orgID.String,
						Name:  orgName.String,
					}
				} else if instanceID.Valid {
					membership.IAM = &IAMMembership{
						IAMID: instanceID.String,
						Name:  instanceName.String,
					}
				} else if projectID.Valid && grantID.Valid && grantedOrgID.Valid {
					membership.ProjectGrant = &ProjectGrantMembership{
						ProjectID:    projectID.String,
						ProjectName:  projectName.String,
						GrantID:      grantID.String,
						GrantedOrgID: grantedOrgID.String,
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
				return nil, zerrors.ThrowInternal(err, "QUERY-N34NV", "Errors.Query.CloseRows")
			}

			return &Memberships{
				Memberships: memberships,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareOrgMember(query *MembershipSearchQuery) (string, []interface{}) {
	builder := sq.Select(
		OrgMemberUserID.identifier(),
		OrgMemberRoles.identifier(),
		OrgMemberCreationDate.identifier(),
		OrgMemberChangeDate.identifier(),
		OrgMemberSequence.identifier(),
		OrgMemberResourceOwner.identifier(),
		OrgMemberInstanceID.identifier(),
		OrgMemberOrgID.identifier(),
		"NULL::TEXT AS "+membershipIAMID.name,
		"NULL::TEXT AS "+membershipProjectID.name,
		"NULL::TEXT AS "+membershipGrantID.name,
	).From(orgMemberTable.identifier())

	for _, q := range query.Queries {
		if q.Col().table.name == membershipAlias.name || q.Col().table.name == orgMemberTable.name {
			builder = q.toQuery(builder)
		}
	}
	return builder.MustSql()
}

func prepareIAMMember(query *MembershipSearchQuery) (string, []interface{}) {
	builder := sq.Select(
		InstanceMemberUserID.identifier(),
		InstanceMemberRoles.identifier(),
		InstanceMemberCreationDate.identifier(),
		InstanceMemberChangeDate.identifier(),
		InstanceMemberSequence.identifier(),
		InstanceMemberResourceOwner.identifier(),
		InstanceMemberInstanceID.identifier(),
		"NULL::TEXT AS "+membershipOrgID.name,
		InstanceMemberIAMID.identifier(),
		"NULL::TEXT AS "+membershipProjectID.name,
		"NULL::TEXT AS "+membershipGrantID.name,
	).From(instanceMemberTable.identifier())

	for _, q := range query.Queries {
		if q.Col().table.name == membershipAlias.name || q.Col().table.name == instanceMemberTable.name {
			builder = q.toQuery(builder)
		}
	}
	return builder.MustSql()
}

func prepareProjectMember(query *MembershipSearchQuery) (string, []interface{}) {
	builder := sq.Select(
		ProjectMemberUserID.identifier(),
		ProjectMemberRoles.identifier(),
		ProjectMemberCreationDate.identifier(),
		ProjectMemberChangeDate.identifier(),
		ProjectMemberSequence.identifier(),
		ProjectMemberResourceOwner.identifier(),
		ProjectMemberInstanceID.identifier(),
		"NULL::TEXT AS "+membershipOrgID.name,
		"NULL::TEXT AS "+membershipIAMID.name,
		ProjectMemberProjectID.identifier(),
		"NULL::TEXT AS "+membershipGrantID.name,
	).From(projectMemberTable.identifier())

	for _, q := range query.Queries {
		if q.Col().table.name == membershipAlias.name || q.Col().table.name == projectMemberTable.name {
			builder = q.toQuery(builder)
		}
	}

	return builder.MustSql()
}

func prepareProjectGrantMember(query *MembershipSearchQuery) (string, []interface{}) {
	builder := sq.Select(
		ProjectGrantMemberUserID.identifier(),
		ProjectGrantMemberRoles.identifier(),
		ProjectGrantMemberCreationDate.identifier(),
		ProjectGrantMemberChangeDate.identifier(),
		ProjectGrantMemberSequence.identifier(),
		ProjectGrantMemberResourceOwner.identifier(),
		ProjectGrantMemberInstanceID.identifier(),
		"NULL::TEXT AS "+membershipOrgID.name,
		"NULL::TEXT AS "+membershipIAMID.name,
		ProjectGrantMemberProjectID.identifier(),
		ProjectGrantMemberGrantID.identifier(),
	).From(projectGrantMemberTable.identifier())

	for _, q := range query.Queries {
		if q.Col().table.name == membershipAlias.name || q.Col().table.name == projectMemberTable.name || q.Col().table.name == projectGrantMemberTable.name {
			builder = q.toQuery(builder)
		}
	}
	return builder.MustSql()
}
