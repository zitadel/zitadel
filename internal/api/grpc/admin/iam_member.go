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
	addedMember, err := s.command.AddIAMMember(ctx, addIamMemberToDomain(member))
	if err != nil {
		return nil, err
	}

	return iamMemberFromDomain(addedMember), nil
}

func (s *Server) ChangeIamMember(ctx context.Context, member *admin.ChangeIamMemberRequest) (*admin.IamMember, error) {
	changedMember, err := s.command.ChangeIAMMember(ctx, changeIamMemberToDomain(member))
	if err != nil {
		return nil, err
	}
	return iamMemberFromDomain(changedMember), nil
}

func (s *Server) RemoveIamMember(ctx context.Context, member *admin.RemoveIamMemberRequest) (*empty.Empty, error) {
	err := s.command.RemoveIAMMember(ctx, member.UserId)
	return &empty.Empty{}, err
}
