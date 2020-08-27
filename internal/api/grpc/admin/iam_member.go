package admin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetIamMemberRoles(ctx context.Context, _ *empty.Empty) (*admin.IamMemberRoles, error) {
	return &admin.IamMemberRoles{Roles: s.iam.GetIAMMemberRoles()}, nil
}

func (s *Server) SearchIamMembers(ctx context.Context, in *admin.IamMemberSearchRequest) (*admin.IamMemberSearchResponse, error) {
	members, err := s.iam.SearchIAMMembers(ctx, iamMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return iamMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddIamMember(ctx context.Context, member *admin.AddIamMemberRequest) (*admin.IamMember, error) {
	addedMember, err := s.iam.AddIAMMember(ctx, addIamMemberToModel(member))
	if err != nil {
		return nil, err
	}

	return iamMemberFromModel(addedMember), nil
}

func (s *Server) ChangeIamMember(ctx context.Context, member *admin.ChangeIamMemberRequest) (*admin.IamMember, error) {
	changedMember, err := s.iam.ChangeIAMMember(ctx, changeIamMemberToModel(member))
	if err != nil {
		return nil, err
	}
	return iamMemberFromModel(changedMember), nil
}

func (s *Server) RemoveIamMember(ctx context.Context, member *admin.RemoveIamMemberRequest) (*empty.Empty, error) {
	err := s.iam.RemoveIAMMember(ctx, member.UserId)
	return &empty.Empty{}, err
}
