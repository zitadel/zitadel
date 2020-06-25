package admin

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/admin/grpc"
)

func (s *Server) GetIamMemberRoles(ctx context.Context, _ *empty.Empty) (*grpc.IamMemberRoles, error) {
	return &grpc.IamMemberRoles{Roles: s.iam.GetIamMemberRoles()}, nil
}

func (s *Server) SearchIamMembers(ctx context.Context, in *grpc.IamMemberSearchRequest) (*grpc.IamMemberSearchResponse, error) {
	members, err := s.iam.SearchIamMembers(ctx, iamMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return iamMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddIamMember(ctx context.Context, member *grpc.AddIamMemberRequest) (*grpc.IamMember, error) {
	addedMember, err := s.iam.AddIamMember(ctx, addIamMemberToModel(member))
	if err != nil {
		return nil, err
	}

	return iamMemberFromModel(addedMember), nil
}

func (s *Server) ChangeIamMember(ctx context.Context, member *grpc.ChangeIamMemberRequest) (*grpc.IamMember, error) {
	changedMember, err := s.iam.ChangeIamMember(ctx, changeIamMemberToModel(member))
	if err != nil {
		return nil, err
	}
	return iamMemberFromModel(changedMember), nil
}

func (s *Server) RemoveIamMember(ctx context.Context, member *grpc.RemoveIamMemberRequest) (*empty.Empty, error) {
	err := s.iam.RemoveIamMember(ctx, member.UserId)
	return &empty.Empty{}, err
}
