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
	return nil, errors.ThrowUnimplemented(nil, "GRPC-Moe56", "Not implemented")
}

func (s *Server) ChangeOrgMember(ctx context.Context, member *ChangeOrgMemberRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-eod34", "Not implemented")
}

func (s *Server) RemoveOrgMember(ctx context.Context, member *RemoveOrgMemberRequest) (*empty.Empty, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-poeSw", "Not implemented")
}
