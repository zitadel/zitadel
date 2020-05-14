package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetOrgMemberRoles(ctx context.Context, _ *empty.Empty) (*OrgMemberRoles, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-wz4vc", "Not implemented")
}

func (s *Server) SearchOrgMembers(ctx context.Context, in *OrgMemberSearchRequest) (*OrgMemberSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-wkdl3", "Not implemented")
}

func (s *Server) AddOrgMember(ctx context.Context, member *AddOrgMemberRequest) (*OrgMember, error) {
	repositoryMember := addOrgMemberToModel(member)

	addedMember, err := s.orgMember.AddOrgMember(ctx, repositoryMember)
	if err != nil {
		return nil, err
	}

	return orgMemberFromModel(addedMember), nil
}

func (s *Server) ChangeOrgMember(ctx context.Context, member *ChangeOrgMemberRequest) (*OrgMember, error) {
	repositoryMember := changeOrgMemberToModel(member)
	changedMember, err := s.orgMember.ChangeOrgMember(ctx, repositoryMember)
	if err != nil {
		return nil, err
	}
	return orgMemberFromModel(changedMember), nil
}

func (s *Server) RemoveOrgMember(ctx context.Context, member *RemoveOrgMemberRequest) (*empty.Empty, error) {
	err := s.orgMember.RemoveOrgMember(ctx, member.OrgId, member.UserId)
	return &empty.Empty{}, err
}
