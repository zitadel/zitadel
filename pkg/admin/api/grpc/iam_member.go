package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetIamMemberRoles(ctx context.Context, _ *empty.Empty) (*IamMemberRoles, error) {
	return &IamMemberRoles{Roles: s.iam.GetIamMemberRoles()}, nil
}

func (s *Server) SearchIamMembers(ctx context.Context, in *IamMemberSearchRequest) (*IamMemberSearchResponse, error) {
	members, err := s.iam.SearchIamMembers(ctx, iamMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return iamMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddIamMember(ctx context.Context, member *AddIamMemberRequest) (*IamMember, error) {
	addedMember, err := s.iam.AddIamMember(ctx, addIamMemberToModel(member))
	if err != nil {
		return nil, err
	}

	return iamMemberFromModel(addedMember), nil
}

func (s *Server) ChangeIamMember(ctx context.Context, member *ChangeIamMemberRequest) (*IamMember, error) {
	changedMember, err := s.iam.ChangeIamMember(ctx, changeIamMemberToModel(member))
	if err != nil {
		return nil, err
	}
	return iamMemberFromModel(changedMember), nil
}

func (s *Server) RemoveIamMember(ctx context.Context, member *RemoveIamMemberRequest) (*empty.Empty, error) {
	err := s.iam.RemoveIamMember(ctx, member.UserId)
	return &empty.Empty{}, err
}
