package group

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	group_pb "github.com/zitadel/zitadel/pkg/grpc/group"
)

func GroupMembershipQueriesToQuery(queries []*group_pb.MembershipQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, 0)
	for _, query := range queries {
		qs, err := GroupMembershipQueryToQuery(query)
		if err != nil {
			return nil, err
		}
		q = append(q, qs)
	}
	return q, nil
}

func GroupMembershipQueryToQuery(req *group_pb.MembershipQuery) (query.SearchQuery, error) {
	switch q := req.Query.(type) {
	case *group_pb.MembershipQuery_OrgQuery:
		return query.NewGroupMembershipOrgIDQuery(q.OrgQuery.OrgId)
	case *group_pb.MembershipQuery_ProjectQuery:
		return query.NewGroupMembershipProjectIDQuery(q.ProjectQuery.ProjectId)
	case *group_pb.MembershipQuery_ProjectGrantQuery:
		return query.NewGroupMembershipProjectGrantIDQuery(q.ProjectGrantQuery.ProjectGrantId)
	case *group_pb.MembershipQuery_IamQuery:
		return query.NewGroupMembershipIsIAMQuery()
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-dsg3z", "Errors.List.Query.Invalid")
	}
}

func GroupMembershipsToMembershipsPb(memberships []*query.GroupMembership) []*group_pb.Membership {
	converted := make([]*group_pb.Membership, len(memberships))
	for i, membership := range memberships {
		converted[i] = GroupMembershipToMembershipPb(membership)
	}
	return converted
}

func GroupMembershipToMembershipPb(membership *query.GroupMembership) *group_pb.Membership {
	typ, name := groupMemberTypeToPb(membership)
	return &group_pb.Membership{
		GroupId:   membership.UserID,
		Type:      typ,
		GroupName: name,
		Roles:     membership.Roles,
		Details: object.ToViewDetailsPb(
			membership.Sequence,
			membership.CreationDate,
			membership.ChangeDate,
			membership.ResourceOwner,
		),
	}
}

func groupMemberTypeToPb(membership *query.GroupMembership) (group_pb.MembershipType, string) {
	if membership.Org != nil {
		return &group_pb.Membership_OrgId{
			OrgId: membership.Org.OrgID,
		}, membership.Org.Name
	} else if membership.Project != nil {
		return &group_pb.Membership_ProjectId{
			ProjectId: membership.Project.ProjectID,
		}, membership.Project.Name
	} else if membership.ProjectGrant != nil {
		return &group_pb.Membership_ProjectGrantId{
			ProjectGrantId: membership.ProjectGrant.GrantID,
		}, membership.ProjectGrant.ProjectName
	} else if membership.IAM != nil {
		return &group_pb.Membership_Iam{
			Iam: true,
		}, membership.IAM.Name
	}
	return nil, ""
}
