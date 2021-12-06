package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
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
		return nil, errors.ThrowInvalidArgument(nil, "USER-dsg3z", "List.Query.Invalid")
	}
}

func MembershipIAMQueryToModel(q *user_pb.MembershipIAMQuery) []*user_model.UserMembershipSearchQuery {
	return []*user_model.UserMembershipSearchQuery{
		{
			Key:    user_model.UserMembershipSearchKeyMemberType,
			Method: domain.SearchMethodEquals,
			Value:  user_model.MemberTypeIam,
		},
		//TODO: q.IAM?
	}
}

func MembershipOrgQueryToModel(q *user_pb.MembershipOrgQuery) []*user_model.UserMembershipSearchQuery {
	return []*user_model.UserMembershipSearchQuery{
		{
			Key:    user_model.UserMembershipSearchKeyMemberType,
			Method: domain.SearchMethodEquals,
			Value:  user_model.MemberTypeOrganisation,
		},
		{
			Key:    user_model.UserMembershipSearchKeyObjectID,
			Method: domain.SearchMethodEquals,
			Value:  q.OrgId,
		},
	}
}

func MembershipProjectQueryToModel(q *user_pb.MembershipProjectQuery) []*user_model.UserMembershipSearchQuery {
	return []*user_model.UserMembershipSearchQuery{
		{
			Key:    user_model.UserMembershipSearchKeyMemberType,
			Method: domain.SearchMethodEquals,
			Value:  user_model.MemberTypeProject,
		},
		{
			Key:    user_model.UserMembershipSearchKeyObjectID,
			Method: domain.SearchMethodEquals,
			Value:  q.ProjectId,
		},
	}
}

func MembershipProjectGrantQueryToModel(q *user_pb.MembershipProjectGrantQuery) []*user_model.UserMembershipSearchQuery {
	return []*user_model.UserMembershipSearchQuery{
		{
			Key:    user_model.UserMembershipSearchKeyMemberType,
			Method: domain.SearchMethodEquals,
			Value:  user_model.MemberTypeProjectGrant,
		},
		{
			Key:    user_model.UserMembershipSearchKeyObjectID,
			Method: domain.SearchMethodEquals,
			Value:  q.ProjectGrantId,
		},
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
	return &user_pb.Membership{
		UserId:      membership.UserID,
		Type:        memberTypeToPb(membership),
		DisplayName: membership.DisplayName,
		Roles:       membership.Roles,
		Details: object.ToViewDetailsPb(
			membership.Sequence,
			membership.CreationDate,
			membership.ChangeDate,
			membership.ResourceOwner,
		),
	}
}

func memberTypeToPb(membership *query.Membership) user_pb.MembershipType {
	if membership.Org != nil {
		return &user_pb.Membership_OrgId{
			OrgId: membership.Org.OrgID,
		}
	} else if membership.Project != nil {
		return &user_pb.Membership_ProjectId{
			ProjectId: membership.Project.ProjectID,
		}
	} else if membership.ProjectGrant != nil {
		return &user_pb.Membership_ProjectGrantId{
			ProjectGrantId: membership.ProjectGrant.GrantID,
		}
	} else if membership.IAM != nil {
		return &user_pb.Membership_Iam{
			Iam: true,
		}
	}
	return nil
}
