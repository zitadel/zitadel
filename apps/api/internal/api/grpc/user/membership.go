package user

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func MembershipQueriesToQuery(queries []*user_pb.MembershipQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, 0)
	for _, query := range queries {
		qs, err := MembershipQueryToQuery(query)
		if err != nil {
			return nil, err
		}
		q = append(q, qs)
	}
	return q, nil
}

func MembershipQueryToQuery(req *user_pb.MembershipQuery) (query.SearchQuery, error) {
	switch q := req.Query.(type) {
	case *user_pb.MembershipQuery_OrgQuery:
		return query.NewMembershipOrgIDQuery(q.OrgQuery.OrgId)
	case *user_pb.MembershipQuery_ProjectQuery:
		return query.NewMembershipProjectIDQuery(q.ProjectQuery.ProjectId)
	case *user_pb.MembershipQuery_ProjectGrantQuery:
		return query.NewMembershipProjectGrantIDQuery(q.ProjectGrantQuery.ProjectGrantId)
	case *user_pb.MembershipQuery_IamQuery:
		return query.NewMembershipIsIAMQuery()
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-dsg3z", "Errors.List.Query.Invalid")
	}
}

func MembershipsToMembershipsPb(memberships []*query.Membership) []*user_pb.Membership {
	converted := make([]*user_pb.Membership, len(memberships))
	for i, membership := range memberships {
		converted[i] = MembershipToMembershipPb(membership)
	}
	return converted
}

func MembershipToMembershipPb(membership *query.Membership) *user_pb.Membership {
	typ, name := memberTypeToPb(membership)
	return &user_pb.Membership{
		UserId:      membership.UserID,
		Type:        typ,
		DisplayName: name,
		Roles:       membership.Roles,
		Details: object.ToViewDetailsPb(
			membership.Sequence,
			membership.CreationDate,
			membership.ChangeDate,
			membership.ResourceOwner,
		),
	}
}

func memberTypeToPb(membership *query.Membership) (user_pb.MembershipType, string) {
	if membership.Org != nil {
		return &user_pb.Membership_OrgId{
			OrgId: membership.Org.OrgID,
		}, membership.Org.Name
	} else if membership.Project != nil {
		return &user_pb.Membership_ProjectId{
			ProjectId: membership.Project.ProjectID,
		}, membership.Project.Name
	} else if membership.ProjectGrant != nil {
		return &user_pb.Membership_ProjectGrantId{
			ProjectGrantId: membership.ProjectGrant.GrantID,
		}, membership.ProjectGrant.ProjectName
	} else if membership.IAM != nil {
		return &user_pb.Membership_Iam{
			Iam: true,
		}, membership.IAM.Name
	}
	return nil, ""
}
