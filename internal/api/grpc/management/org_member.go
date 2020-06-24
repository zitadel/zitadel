package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/management/grpc"
)

func (s *Server) GetOrgMemberRoles(ctx context.Context, _ *empty.Empty) (*grpc.OrgMemberRoles, error) {
	return &grpc.OrgMemberRoles{Roles: s.org.GetOrgMemberRoles()}, nil
}

func (s *Server) SearchMyOrgMembers(ctx context.Context, in *grpc.OrgMemberSearchRequest) (*grpc.OrgMemberSearchResponse, error) {
	members, err := s.org.SearchMyOrgMembers(ctx, orgMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddMyOrgMember(ctx context.Context, member *grpc.AddOrgMemberRequest) (*grpc.OrgMember, error) {
	addedMember, err := s.org.AddMyOrgMember(ctx, addOrgMemberToModel(member))
	if err != nil {
		return nil, err
	}

	return orgMemberFromModel(addedMember), nil
}

func (s *Server) ChangeMyOrgMember(ctx context.Context, member *grpc.ChangeOrgMemberRequest) (*grpc.OrgMember, error) {
	changedMember, err := s.org.ChangeMyOrgMember(ctx, changeOrgMemberToModel(member))
	if err != nil {
		return nil, err
	}
	return orgMemberFromModel(changedMember), nil
}

func (s *Server) RemoveMyOrgMember(ctx context.Context, member *grpc.RemoveOrgMemberRequest) (*empty.Empty, error) {
	err := s.org.RemoveMyOrgMember(ctx, member.UserId)
	return &empty.Empty{}, err
}
