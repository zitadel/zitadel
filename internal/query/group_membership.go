package query

import (
	"context"
	"database/sql"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type GroupMemberships struct {
	SearchResponse
	GroupMemberships []*GroupMembership
}

type GroupMembership struct {
	UserID  string
	GroupID string
	Roles   database.TextArray[string]

	CreationDate  time.Time
	ChangeDate    time.Time
	Sequence      uint64
	ResourceOwner string

	Org          *OrgMembership
	IAM          *IAMMembership
	Project      *ProjectMembership
	ProjectGrant *ProjectGrantMembership
}

var (
	//groupMembershipAlias is a hack to satisfy checks in the queries
	groupMembershipAlias = table{
		name:          "group_members1", // Old value: members
		instanceIDCol: projection.MemberInstanceID,
	}
	groupMembershipUserID = Column{
		name:  projection.MemberUserIDCol,
		table: groupMembershipAlias,
	}
	groupMembershipRoles = Column{
		name:  projection.MemberRolesCol,
		table: groupMembershipAlias,
	}
	groupMembershipCreationDate = Column{
		name:  projection.MemberCreationDate,
		table: groupMembershipAlias,
	}
	groupMembershipChangeDate = Column{
		name:  projection.MemberChangeDate,
		table: groupMembershipAlias,
	}
	groupMembershipSequence = Column{
		name:  projection.MemberSequence,
		table: groupMembershipAlias,
	}
	groupMembershipResourceOwner = Column{
		name:  projection.MemberResourceOwner,
		table: groupMembershipAlias,
	}
	groupMembershipInstanceID = Column{
		name:  projection.MemberInstanceID,
		table: groupMembershipAlias,
	}
	groupMembershipOrgID = Column{
		name:  projection.OrgMemberOrgIDCol,
		table: groupMembershipAlias,
	}
	groupMembershipIAMID = Column{
		name:  projection.InstanceMemberIAMIDCol,
		table: groupMembershipAlias,
	}
	groupMembershipProjectID = Column{
		name:  projection.ProjectMemberProjectIDCol,
		table: groupMembershipAlias,
	}
	groupMembershipGrantID = Column{
		name:  projection.ProjectGrantMemberGrantIDCol,
		table: groupMembershipAlias,
	}
	groupMembershipGrantGrantedOrgID = Column{
		name:  projection.ProjectGrantColumnGrantedOrgID,
		table: groupMembershipAlias,
	}
	groupMembershipGroupID = Column{
		name:  projection.GroupMemberGroupIDCol,
		table: groupMembershipAlias,
	}
	// membershipGroupGrantID = Column{
	// 	name:  projection.GroupGrantIDCol,
	// 	table: membershipAlias,
	// }
)

func NewGroupMembershipUserIDQuery(userID string) (SearchQuery, error) {
	return NewTextQuery(groupMembershipUserID.setTable(groupMembershipAlias), userID, TextEquals)
}

func NewGroupMembershipGroupIDQuery(groupID string) (SearchQuery, error) {
	return NewTextQuery(groupMembershipGroupID.setTable(groupMembershipAlias), groupID, TextEquals)
}

func NewGroupMembershipOrgIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(OrgMemberOrgID, value, TextEquals)
}

func NewGroupMembershipResourceOwnersSearchQuery(ids ...string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(groupMembershipResourceOwner, list, ListIn)
}

func NewGroupMembershipGrantedOrgIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantColumnGrantedOrgID, id, TextEquals)
}

func NewGroupMembershipProjectIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectMemberProjectID, value, TextEquals)
}

func NewGroupMembershipProjectGrantIDQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ProjectGrantMemberGrantID, value, TextEquals)
}

func NewGroupMembershipIsIAMQuery() (SearchQuery, error) {
	return NewNotNullQuery(InstanceMemberIAMID)
}

func (q *Queries) GroupMemberships(ctx context.Context, queries *MembershipSearchQuery, shouldTrigger bool) (groupMemberships *GroupMemberships, err error) {
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

	query, queryArgs, scan := prepareGroupMembershipsQuery(queries)
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
		groupMemberships, err = scan(rows)
		return err
	}, stmt, queryArgs...)
	if err != nil {
		return nil, err
	}
	groupMemberships.State = latestSequence
	return groupMemberships, nil
}

func getGroupMembershipFromQuery(queries *MembershipSearchQuery) (string, []interface{}) {
	orgMembers, orgMembersArgs := prepareGroupOrgMember(queries)
	iamMembers, iamMembersArgs := prepareGroupIAMMember(queries)
	projectMembers, projectMembersArgs := prepareGroupProjectMember(queries)
	projectGrantMembers, projectGrantMembersArgs := prepareGroupProjectGrantMember(queries)
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
			") AS " + groupMembershipAlias.identifier(),
		args
}

func prepareGroupMembershipsQuery(queries *MembershipSearchQuery) (sq.SelectBuilder, []interface{}, func(*sql.Rows) (*GroupMemberships, error)) {
	query, args := getGroupMembershipFromQuery(queries)
	return sq.Select(
			membershipGroupID.identifier(),
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
			LeftJoin(join(InstanceColumnID, membershipInstanceID)).
			PlaceholderFormat(sq.Dollar),
		args,
		func(rows *sql.Rows) (*GroupMemberships, error) {
			memberships := make([]*GroupMembership, 0)
			var count uint64
			for rows.Next() {

				var (
					membership   = new(GroupMembership)
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

			return &GroupMemberships{
				GroupMemberships: memberships,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareGroupOrgMember(query *MembershipSearchQuery) (string, []interface{}) {
	builder := sq.Select(
		OrgMemberUserID.identifier(),
		OrgMemberRoles.identifier(),
		OrgMemberCreationDate.identifier(),
		OrgMemberChangeDate.identifier(),
		OrgMemberSequence.identifier(),
		OrgMemberResourceOwner.identifier(),
		OrgMemberInstanceID.identifier(),
		OrgMemberOrgID.identifier(),
		"NULL::TEXT AS "+groupMembershipIAMID.name,
		"NULL::TEXT AS "+groupMembershipProjectID.name,
		"NULL::TEXT AS "+groupMembershipGrantID.name,
	).From(orgMemberTable.identifier())

	for _, q := range query.Queries {
		if q.Col().table.name == groupMembershipAlias.name || q.Col().table.name == orgMemberTable.name {
			builder = q.toQuery(builder)
		}
	}
	return builder.MustSql()
}

func prepareGroupIAMMember(query *MembershipSearchQuery) (string, []interface{}) {
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
		if q.Col().table.name == groupMembershipAlias.name || q.Col().table.name == instanceMemberTable.name {
			builder = q.toQuery(builder)
		}
	}
	return builder.MustSql()
}

func prepareGroupProjectMember(query *MembershipSearchQuery) (string, []interface{}) {
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
		if q.Col().table.name == groupMembershipAlias.name || q.Col().table.name == projectMemberTable.name {
			builder = q.toQuery(builder)
		}
	}

	return builder.MustSql()
}

// Need to add GroupID but for that need to update project_grant table for adding groupID
func prepareGroupProjectGrantMember(query *MembershipSearchQuery) (string, []interface{}) {
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
		if q.Col().table.name == groupMembershipAlias.name || q.Col().table.name == projectMemberTable.name || q.Col().table.name == projectGrantMemberTable.name {
			builder = q.toQuery(builder)
		}
	}
	return builder.MustSql()
}
