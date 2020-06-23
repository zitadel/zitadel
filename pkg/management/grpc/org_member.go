package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetOrgMemberRoles(ctx context.Context, _ *empty.Empty) (*OrgMemberRoles, error) {
	return &OrgMemberRoles{Roles: s.org.GetOrgMemberRoles()}, nil
}

func (s *Server) SearchMyOrgMembers(ctx context.Context, in *OrgMemberSearchRequest) (*OrgMemberSearchResponse, error) {
	members, err := s.org.SearchMyOrgMembers(ctx, orgMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddMyOrgMember(ctx context.Context, member *AddOrgMemberRequest) (*OrgMember, error) {
	addedMember, err := s.org.AddMyOrgMember(ctx, addOrgMemberToModel(member))
	if err != nil {
		return nil, err
	}

	return orgMemberFromModel(addedMember), nil
}

func (s *Server) ChangeMyOrgMember(ctx context.Context, member *ChangeOrgMemberRequest) (*OrgMember, error) {
	changedMember, err := s.org.ChangeMyOrgMember(ctx, changeOrgMemberToModel(member))
	if err != nil {
		return nil, err
	}
	return orgMemberFromModel(changedMember), nil
}

func (s *Server) RemoveMyOrgMember(ctx context.Context, member *RemoveOrgMemberRequest) (*empty.Empty, error) {
	err := s.org.RemoveMyOrgMember(ctx, member.UserId)
	return &empty.Empty{}, err
}
