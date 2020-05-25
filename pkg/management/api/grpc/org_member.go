package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

const (
	orgRolePrefix = "ORG_"
)

func (s *Server) GetOrgMemberRoles(ctx context.Context, _ *empty.Empty) (*OrgMemberRoles, error) {
	return &OrgMemberRoles{Roles: s.getOrgMemberRoles()}, nil
}

func (s *Server) SearchOrgMembers(ctx context.Context, in *OrgMemberSearchRequest) (*OrgMemberSearchResponse, error) {
	members, err := s.org.SearchOrgMembers(ctx, orgMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddOrgMember(ctx context.Context, member *AddOrgMemberRequest) (*OrgMember, error) {
	repositoryMember := addOrgMemberToModel(member)

	addedMember, err := s.org.AddOrgMember(ctx, repositoryMember)
	if err != nil {
		return nil, err
	}

	return orgMemberFromModel(addedMember), nil
}

func (s *Server) ChangeOrgMember(ctx context.Context, member *ChangeOrgMemberRequest) (*OrgMember, error) {
	repositoryMember := changeOrgMemberToModel(member)
	changedMember, err := s.org.ChangeOrgMember(ctx, repositoryMember)
	if err != nil {
		return nil, err
	}
	return orgMemberFromModel(changedMember), nil
}

func (s *Server) RemoveOrgMember(ctx context.Context, member *RemoveOrgMemberRequest) (*empty.Empty, error) {
	err := s.org.RemoveOrgMember(ctx, member.OrgId, member.UserId)
	return &empty.Empty{}, err
}
