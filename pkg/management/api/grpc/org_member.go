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

func (s *Server) AddOrgMember(ctx context.Context, member *AddOrgMemberRequest) (*empty.Empty, error) {
	repositoryMember := addOrgMemberToModel(member)

	_, err := s.orgMember.AddOrgMember(ctx, repositoryMember)
	return &empty.Empty{}, err
}

func (s *Server) ChangeOrgMember(ctx context.Context, member *ChangeOrgMemberRequest) (*empty.Empty, error) {
	repositoryMember := changeOrgMemberToModel(member)
	_, err := s.orgMember.ChangeOrgMember(ctx, repositoryMember)
	return &empty.Empty{}, err
}

func (s *Server) RemoveOrgMember(ctx context.Context, member *RemoveOrgMemberRequest) (*empty.Empty, error) {
	err := s.orgMember.RemoveOrgMember(ctx, member.OrgId, member.UserId)
	return &empty.Empty{}, err
}
