package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgMemberRoles(ctx context.Context, _ *empty.Empty) (*management.OrgMemberRoles, error) {
	return &management.OrgMemberRoles{Roles: s.org.GetOrgMemberRoles()}, nil
}

func (s *Server) SearchMyOrgMembers(ctx context.Context, in *management.OrgMemberSearchRequest) (*management.OrgMemberSearchResponse, error) {
	members, err := s.org.SearchMyOrgMembers(ctx, orgMemberSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgMemberSearchResponseFromModel(members), nil
}

func (s *Server) AddMyOrgMember(ctx context.Context, member *management.AddOrgMemberRequest) (*management.OrgMember, error) {
	addedMember, err := s.command.AddOrgMember(ctx, addOrgMemberToDomain(ctx, member))
	if err != nil {
		return nil, err
	}

	return orgMemberFromDomain(addedMember), nil
}

func (s *Server) ChangeMyOrgMember(ctx context.Context, member *management.ChangeOrgMemberRequest) (*management.OrgMember, error) {
	changedMember, err := s.command.ChangeOrgMember(ctx, changeOrgMemberToModel(ctx, member))
	if err != nil {
		return nil, err
	}
	return orgMemberFromDomain(changedMember), nil
}

func (s *Server) RemoveMyOrgMember(ctx context.Context, member *management.RemoveOrgMemberRequest) (*empty.Empty, error) {
	err := s.command.RemoveOrgMember(ctx, authz.GetCtxData(ctx).OrgID, member.UserId)
	return &empty.Empty{}, err
}
