package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func MembershipQueriesToModel(queries []*user_pb.MembershipQuery) (_ []*user_model.UserMembershipSearchQuery, err error) {
	q := make([]*user_model.UserMembershipSearchQuery, 0)
	for _, query := range queries {
		qs, err := MembershipQueryToModel(query)
		if err != nil {
			return nil, err
		}
		q = append(q, qs...)
	}
	return q, nil
}

func MembershipQueryToModel(query *user_pb.MembershipQuery) ([]*user_model.UserMembershipSearchQuery, error) {
	switch q := query.Query.(type) {
	case *user_pb.MembershipQuery_Org:
		return MembershipOrgQueryToModel(q.Org), nil
	case *user_pb.MembershipQuery_Project:
		return MembershipProjectQueryToModel(q.Project), nil
	case *user_pb.MembershipQuery_ProjectGrant:
		return MembershipProjectGrantQueryToModel(q.ProjectGrant), nil
	case *user_pb.MembershipQuery_Iam:
		return MembershipIAMQueryToModel(q.Iam), nil
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

func MembershipsToMembershipsPb(memberships []*user_model.UserMembershipView) []*user_pb.Membership {
	converted := make([]*user_pb.Membership, len(memberships))
	for i, membership := range memberships {
		converted[i] = MembershipToMembershipPb(membership)
	}
	return converted
}

func MembershipToMembershipPb(membership *user_model.UserMembershipView) *user_pb.Membership {
	return &user_pb.Membership{
		UserId:      membership.UserID,
		Type:        memberTypeToPb(membership),
		DisplayName: membership.DisplayName,
		Roles:       membership.Roles,
		Details: object.ToDetailsPb(
			membership.Sequence,
			membership.CreationDate,
			membership.ChangeDate,
			membership.ResourceOwner,
		),
	}
}

func memberTypeToPb(membership *user_model.UserMembershipView) user_pb.MembershipType {
	switch membership.MemberType {
	case user_model.MemberTypeOrganisation:
		return &user_pb.Membership_OrgId{
			OrgId: membership.AggregateID,
		}
	case user_model.MemberTypeProject:
		return &user_pb.Membership_ProjectId{
			ProjectId: membership.AggregateID,
		}
	case user_model.MemberTypeProjectGrant:
		return &user_pb.Membership_ProjectGrantId{
			ProjectGrantId: membership.ObjectID,
		}
	case user_model.MemberTypeIam:
		return &user_pb.Membership_Iam{
			Iam: true, //TODO: ?
		}
	default:
		return nil //TODO: ?
	}
}
